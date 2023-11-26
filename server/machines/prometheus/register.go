package prometheus

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
	"github.com/sirupsen/logrus"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	sysID := uuid.Nil
	userUUID := uuid.Nil

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	connPayload, err := utils.Cast[models.ConnectionPayload](data)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	promConn, err := utils.Cast[PromConn](connPayload.MetaData)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	promCred, err := utils.Cast[PromCred](connPayload.CredentialSecret)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}
	logrus.Debug(promCred, "PROM CRED")
	promClient := models.NewPrometheusClient()

	err = promClient.Validate(ctx, promConn.URL) // change this to accept credentials either basicauth or API Key.

	if err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrPrometheusScan(err)}).Build(), nil
	}
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
