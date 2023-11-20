package kubernetes

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
)

type IgnoreAction struct {}

func(ia *IgnoreAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {

	return machines.NoOp, nil, nil
}
func(ia *IgnoreAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	
	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.IGNORED)

	if err != nil {
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to change status of the connection \"%s\" to ignore.", machinectx.K8sContext.Name)).Build(), err
	}

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from ignored state", connection)

	return machines.NoOp, eventBuilder.Build(), err
}

func(ia *IgnoreAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
