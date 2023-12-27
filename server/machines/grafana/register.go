package grafana

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	connPayload, err := utils.Cast[models.ConnectionPayload](data)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	metadata, err := utils.Cast[map[string]interface{}](connPayload.MetaData)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaConn, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.GrafanaConn](metadata)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaCred, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.GrafanaCred](connPayload.CredentialSecret)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	grafanaClient := models.NewGrafanaClient()
	err = grafanaClient.Validate(ctx, grafanaConn.URL, grafanaCred.APIKeyOrBasicAuth)
	if err != nil && !connPayload.SkipCredentialVerification {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrGrafanaScan(err)}).Build(), models.ErrGrafanaScan(err)
	}
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
