package grafana

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
	"github.com/sirupsen/logrus"
)

type ConnectAction struct{}

func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {

	return machines.NoOp, nil, nil
}
func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {

	sysID := uuid.Nil
	userUUID := uuid.Nil
	token, _ := ctx.Value(models.TokenCtxKey).(string)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := utils.Cast[MachineCtx](machineCtx)
	if err != nil {
		err := machines.ErrAssertMachineCtx(fmt.Errorf("asserting of context %v failed", machinectx))
		return machines.NoOp, eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": err,
		}).Build(), err
	}

	grafanaConn := machinectx.GrafanaConn
	secret := make(map[string]interface{}, 0)
	secret["auth"] = machinectx.GrafanaCred.APIKey
	credential, err := machinectx.provider.SaveUserCredential(token, &models.Credential{
		Name:   machinectx.GrafanaCred.Name,
		UserID: &userUUID,
		Type:   grafana,
		Secret: secret,
	})

	if err != nil {
		_err := models.ErrPersistCredential(err)
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to persist credential information for the connection %s", grafanaConn.Name)).
			WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), _err
	}

	connection, err := machinectx.provider.SaveConnection(&models.ConnectionPayload{
		Kind:    grafana,
		Type:    "observability",
		SubType: "monitoring",
		Status:  connections.CONNECTED,
		Name:    machinectx.GrafanaConn.Name,
		MetaData: map[string]interface{}{
			"name": grafanaConn.Name,
			"url":  grafanaConn.URL,
		},
		CredentialID: &credential.ID,
	}, token, false)
	if err != nil {
		_err := models.ErrPersistConnection(err)
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to perisit the \"%s\" connection details", grafanaConn.Name)).WithMetadata(map[string]interface{}{"error": _err}).Build(), _err
	}
	logrus.Debug(connection, "grafana connected connection")
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
