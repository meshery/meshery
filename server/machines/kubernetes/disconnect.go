package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
)

type DisconnectAction struct {}

func(da *DisconnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(uuid.Nil).FromUser(uuid.Nil) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	err = AssignClientSetToContext(machinectx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}
	
	return machines.NoOp, eventBuilder.Build(), nil

}
func(da *DisconnectAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error)  {
	eventBuilder := events.NewEvent().ActedUpon(uuid.Nil).WithCategory("connection").WithAction("register").FromSystem(uuid.Nil).FromUser(uuid.Nil) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}

	k8sContexts := []models.K8sContext{machinectx.K8sContext}
	machinectx.MesheryCtrlsHelper.UndeployDeployedOperators(machinectx.OperatorTracker)

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		machinectx.MesheryCtrlsHelper.UpdateCtxControllerHandlers(k8sContexts)
	})
	// ctx = context.WithValue(ctx, models.MesheryControllerHandlersKey, machinectx. MesheryCtrlsHelper.GetControllerHandlersForEachContext()) // not required verify once

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.DISCONNECTED)

	// peform error handling and event publishing
	if err != nil {
		return machines.NoOp, eventBuilder.Build(), err
	}

	fmt.Println("connection inside disconnect.go", connection, statusCode)

	return machines.NoOp, eventBuilder.Build(), nil
}

func(da *DisconnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
