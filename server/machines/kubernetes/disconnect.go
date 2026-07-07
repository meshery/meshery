package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

type DisconnectAction struct{}

func (da *DisconnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil

}
func (da *DisconnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok || user == nil {
		err := fmt.Errorf("user missing from context")
		return machines.NoOp, events.NewEvent().WithCategory("connection").WithAction("update").WithDescription(err.Error()).WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	sysID, ok := ctx.Value(models.SystemIDKey).(*core.Uuid)
	if !ok || sysID == nil {
		err := fmt.Errorf("system ID missing from context")
		return machines.NoOp, events.NewEvent().WithCategory("connection").WithAction("update").WithDescription(err.Error()).WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	userUUID := user.ID

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	contextID := machinectx.K8sContext.ID
	go func() {
		machinectx.MesheryCtrlsHelper.
			UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			UndeployDeployedOperators(machinectx.OperatorTracker).
			RemoveCtxControllerHandler(ctx, contextID)
		machinectx.MesheryCtrlsHelper.RemoveMeshSyncDataHandler(ctx, contextID)

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
