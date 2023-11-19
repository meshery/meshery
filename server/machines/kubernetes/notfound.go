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

type NotFoundAction struct {}

func(ia *NotFoundAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {

	return machines.NoOp, nil, nil
}
func(ia *NotFoundAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	
	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.NOTFOUND)

	if err != nil {
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Connection with the context \"%s\" at \"%s\" is unreachable.", machinectx.K8sContext.Name, machinectx.K8sContext.Server)).Build(), err
	}

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from notfound state", connection)

	return machines.NoOp, eventBuilder.Build(), err
}

func(ia *NotFoundAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
