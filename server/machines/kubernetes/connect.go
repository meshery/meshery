package kubernetes

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"
)

type ConnectAction struct{}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	meshsyncDeploymentMode := models.MeshsyncDeploymentModeOperator
	//  TODO:
	// this viper value check here is a temporal thing
	// meshsync deployment mode will be propagated from connection entity
	if viper.GetBool("TMP_MESHSYNC_AS_A_LIBRARY_MODE") {
		meshsyncDeploymentMode = models.MeshsyncDeploymentModeLibrary
	}

	go func() {
		ctrlHelper := machinectx.MesheryCtrlsHelper.
			AddCtxControllerHandlers(machinectx.K8sContext).
			SetMeshsyncDeploymentMode(meshsyncDeploymentMode).
			UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			DeployUndeployedOperators(machinectx.OperatorTracker)
		ctrlHelper.AddMeshsynDataHandlers(ctx, machinectx.K8sContext, userUUID, *sysID, provider)
	}()
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
