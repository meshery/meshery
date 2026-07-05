package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
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
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

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

	// Resolve the layered controllers configuration for this connection:
	// per-connection override -> server-wide defaults (Settings) -> built-in
	// defaults. The merged document is stashed on the controllers helper so
	// the embedded meshsync run picks up its knobs, and is applied to the
	// cluster after the operator machinery attaches (operator mode).
	mergedControllersConfig, _, errResolve := machinectx.MesheryCtrlsHelper.ResolveControllersConfigForConnection(connection.Metadata)
	controllersConfigResolved := errResolve == nil
	if errResolve != nil {
		// The defaults store failed (a malformed override alone degrades
		// inside the resolver and still yields the Settings defaults):
		// connect with layer-free defaults and skip the cluster apply
		// below, since the intended configuration is unknown.
		machinectx.log.Error(errResolve)
		mergedControllersConfig = nil
	}

	meshsyncDeploymentMode := connections.MeshsyncDeploymentModeFromMetadata(connection.Metadata)
	if meshsyncDeploymentMode == connections.MeshsyncDeploymentModeUndefined {
		// Fall back to the server-wide default configured in Settings.
		meshsyncDeploymentMode = connections.DeploymentModeFromControllersConfig(mergedControllersConfig)
	}
	if meshsyncDeploymentMode == connections.MeshsyncDeploymentModeUndefined {
		// TODO:
		// maybe not call to viper here and propagate default value from above,
		// f.e. when machine is created
		meshsyncDeploymentMode = connections.MeshsyncDeploymentModeFromString(
			viper.GetString("MESHSYNC_DEFAULT_DEPLOYMENT_MODE"),
		)
		if meshsyncDeploymentMode == connections.MeshsyncDeploymentModeUndefined {
			meshsyncDeploymentMode = connections.MeshsyncDeploymentModeDefault
		}
	}

	go func() {
		ctrlHelper := machinectx.MesheryCtrlsHelper.
			AddCtxControllerHandlers(machinectx.K8sContext).
			SetMeshsyncDeploymentMode(meshsyncDeploymentMode).
			SetControllersConfig(mergedControllersConfig).
			UpdateOperatorsStatusMap(machinectx.OperatorTracker).
			DeployUndeployedOperators(machinectx.OperatorTracker)
		ctrlHelper.AddMeshsynDataHandlers(ctx, machinectx.K8sContext, userUUID, *sysID, provider)

		// Operator mode: best-effort apply of the explicitly-set
		// configuration onto the cluster's MeshSync/Broker custom resources
		// and the MeshSync deployment. An empty resolved configuration is
		// applied too: it withdraws previously-applied fields (cleared
		// while the connection was down) and heals drift. Targets that do
		// not exist yet (the operator is still coming up) are skipped and
		// re-applied on the next connect or configuration change.
		if meshsyncDeploymentMode == connections.MeshsyncDeploymentModeOperator && controllersConfigResolved {
			kubeClient, errClient := machinectx.K8sContext.GenerateKubeHandler()
			if errClient != nil {
				machinectx.log.Error(ErrConnectAction(errClient))
				return
			}
			// Detached context: the connect request's context ends with the
			// HTTP request, while this apply runs alongside the async
			// operator deployment.
			applyCtx, cancelApply := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancelApply()
			if _, errApply := models.ApplyControllersConfigToCluster(applyCtx, machinectx.log, kubeClient, mergedControllersConfig); errApply != nil {
				machinectx.log.Error(errApply)
			}
		}
	}()
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
