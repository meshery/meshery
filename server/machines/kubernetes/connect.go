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

type ConnectAction struct{}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	k8sContexts := []models.K8sContext{machinectx.K8sContext}
	ctrlHelper := machinectx.MesheryCtrlsHelper.UpdateCtxControllerHandlers(k8sContexts).
		UpdateOperatorsStatusMap(machinectx.OperatorTracker).DeployUndeployedOperators(machinectx.OperatorTracker)
	ctrlHelper.UpdateMeshsynDataHandlers()

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := machinectx.Provider.UpdateConnectionStatusByID(token, uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID), connections.CONNECTED)

	if err != nil {
		return machines.NoOp, eventBuilder.WithDescription(fmt.Sprintf("Connect operation succeeded but failed to update the record for the connection with context \"%s\" at %s", machinectx.K8sContext.Name, machinectx.K8sContext.Server)).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	machinectx.log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)
	machinectx.log.Debug("exiting execute func from disconnected state", connection)
	
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
