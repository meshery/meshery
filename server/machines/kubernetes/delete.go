package kubernetes

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
)

type DeleteAction struct{}

func (da *DeleteAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (da *DeleteAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	k8sContexts := []models.K8sContext{machinectx.K8sContext}
	machinectx.MesheryCtrlsHelper.UndeployDeployedOperators(machinectx.OperatorTracker)

	_ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	context.AfterFunc(_ctx, func() {
		machinectx.MesheryCtrlsHelper.UpdateCtxControllerHandlers(k8sContexts)
	})

	go models.FlushMeshSyncData(ctx, machinectx.K8sContext, machinectx.Provider, machinectx.EventBroadcaster, user.ID, sysID)

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
