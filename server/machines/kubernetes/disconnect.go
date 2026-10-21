package kubernetes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

type DisconnectAction struct{}

func (da *DisconnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil

}
func (da *DisconnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	userUUID, sysID := machines.ExtractActorIDs(ctx)

	eventBuilder := events.NewEvent().WithCategory("connection").WithAction("update").WithDescription("Failed to interact with the connection.")
	if userUUID != core.Uuid(uuid.Nil) {
		eventBuilder.ActedUpon(userUUID)
	}
	machines.ApplyEventAttribution(eventBuilder, userUUID, sysID)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID
	if machinectx.MesheryCtrlsHelper != nil {
		go func() {
			machinectx.MesheryCtrlsHelper.
				UpdateOperatorsStatusMap(machinectx.OperatorTracker).
				UndeployDeployedOperators(machinectx.OperatorTracker, contextID).
				RemoveCtxControllerHandler(ctx, contextID)
			machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}()
	}

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
