package kubernetes

import (
	"context"

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
	
	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("register").FromSystem(*sysID).FromUser(userUUID) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.NOTFOUND)

	// peform error handling and event publishing
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from notfound state", connection)

	return machines.NoOp, eventBuilder.Build(), err
}

func(ia *NotFoundAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
