package kubernetes

import (
	"context"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshkit/models/events"
)

type DisconnectAction struct{}

func (da *DisconnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil

}
// Execute cleanly tears down MeshSync data handlers and disconnects from the Kubernetes cluster.
func (da *DisconnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	userUUID, systemID := machines.ResolveUserAndSystemID(ctx)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(systemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID
	go func() {
		if machinectx.MesheryCtrlsHelper != nil {
			machinectx.MesheryCtrlsHelper.
				UpdateOperatorsStatusMap(machinectx.OperatorTracker).
				UndeployDeployedOperators(machinectx.OperatorTracker, contextID).
				RemoveCtxControllerHandler(ctx, contextID)
			machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}
	}()

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		// machinectx.MesheryCtrlsHelper.UpdateOperatorsStatusMap(machinectx.OperatorTracker)
	})

	return machines.NoOp, nil, nil
}

func (da *DisconnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
