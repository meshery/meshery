package grafana

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	sysID := uuid.Nil
	userUUID := uuid.Nil

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := utils.Cast[MachineCtx](machineCtx)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	err = machinectx.GrafanaClient.Validate(ctx, machinectx.GrafanaConn.URL, machinectx.GrafanaCred.APIKey)
	if err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": handlers.ErrGrafanaScan(err)}).Build(), nil
	}
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
