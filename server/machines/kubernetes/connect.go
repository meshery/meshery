package kubernetes

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/spf13/viper"
)

type ConnectAction struct{}

// Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok || user == nil {
		err := machines.ErrMissingUserContext()
		return machines.NoOp, nil, err
	}
	sysID, ok := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	if !ok || sysID == nil {
		err := machines.ErrMissingSystemIDContext()
		return machines.NoOp, nil, err
	}
	userUUID := user.ID
	provider, ok := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if !ok || provider == nil {
		err := ErrConnectAction(fmt.Errorf("failed to retrieve provider from context"))
		return machines.NoOp, nil, err
	}

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	token, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok {
		errToken := ErrConnectAction(fmt.Errorf("failed to retrieve user token"))
		eventBuilder.WithMetadata(map[string]interface{}{"error": errToken})
		return machines.NoOp, eventBuilder.Build(), errToken

	}
	connectionID := uuid.FromStringOrNil(machinectx.K8sContext.ConnectionID)
	if connectionID == uuid.Nil {
		errConnection := ErrConnectAction(fmt.Errorf("k8sCtx.ConnectionID is empty or invalid"))
		eventBuilder.WithMetadata(map[string]interface{}{"error": errConnection})
		return machines.NoOp, eventBuilder.Build(), errConnection
	}

	connection, _, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		errConnection := ErrConnectAction(err)
		eventBuilder.WithMetadata(map[string]interface{}{"error": errConnection})
		return machines.NoOp, eventBuilder.Build(), errConnection

	}
	if connection.Kind != "kubernetes" {
		errConnection := ErrConnectAction(fmt.Errorf("connection is not of kind kubernetes"))
		eventBuilder.WithMetadata(map[string]interface{}{"error": errConnection})
		return machines.NoOp, eventBuilder.Build(), errConnection
	}

	meshsyncDeploymentMode := schemasConnection.MeshsyncDeploymentModeFromMetadata(connection.Metadata)
	if meshsyncDeploymentMode == schemasConnection.MeshsyncDeploymentModeUndefined {
		// TODO:
		// maybe not call to viper here and propagate default value from above,
		// f.e. when machine is created
		meshsyncDeploymentMode = schemasConnection.MeshsyncDeploymentModeFromString(
			viper.GetString("MESHSYNC_DEFAULT_DEPLOYMENT_MODE"),
		)
		if meshsyncDeploymentMode == schemasConnection.MeshsyncDeploymentModeUndefined {
			meshsyncDeploymentMode = schemasConnection.MeshsyncDeploymentModeDefault
		}
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
