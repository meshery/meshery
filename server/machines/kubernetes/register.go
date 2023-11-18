package kubernetes

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshery/server/models/meshmodel/core"
	"github.com/layer5io/meshkit/models/events"
)

type RegisterAction struct {}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func(ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	
	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("register").FromSystem(*sysID).FromUser(userUUID) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	err = AssignClientSetToContext(machinectx, eventBuilder)
	if err != nil {
		return machines.NotFound, eventBuilder.Build(), err
	}
	
	return machines.NoOp, eventBuilder.Build(), nil
}

func(ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(*sysID).FromUser(userUUID) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id

	
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	machinectx.log.Debug("executing ping test for connection", machinectx.K8sContext.ConnectionID)
	err = machinectx.K8sContext.PingTest()
	
	if err != nil {
		machinectx.log.Error(err)	
		// peform error handling and event publishing
		return machines.NotFound, nil, err
	}
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	context := []*models.K8sContext{&machinectx.K8sContext}

	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.REGISTERED)
	machinectx.K8sCompRegHelper.UpdateContexts(context).RegisterComponents(context, []models.K8sRegistrationFunction{core.RegisterK8sMeshModelComponents}, machinectx.RegistryManager, machinectx.EventBroadcaster, machinectx.Provider, user.ID, true)
	// peform error handling and event publishing
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from registered state", connection)

	return machines.Connect, nil, nil
}

func(ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
