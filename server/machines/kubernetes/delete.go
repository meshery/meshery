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

	token, _ := ctx.Value(models.TokenCtxKey).(string)

	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.DELETED)

	if err != nil {
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Delete operation succeeded but failed to update the record for the connection with context \"%s\" at %s", machinectx.K8sContext.Name, machinectx.K8sContext.Server)).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	go models.FlushMeshSyncData(ctx, machinectx.K8sContext, machinectx.Provider, machinectx.EventBroadcaster, user.ID, sysID)

	machinectx.log.Debug("deleted connection: ", connection, "HTTP status", statusCode)

	return machines.NoOp, nil, nil
}

func (da *DeleteAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
