package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"
)

type ConnectAction struct{}

// ExecuteOnEntry is a no-op; ConnectAction performs all of its work in Execute.
//
// TODO: Execute On Entry and Exit should not return next eventtype i suppose, look again.
func (ca *ConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// Execute establishes connection to the Kubernetes cluster and registers MeshSync data handlers.
// isConnectionMeantToBeConnected reports whether the connection is in a state
// intended to be active and connected. If the connection transitioned to
// disconnected, deleted, or ignored, queued retries must not execute.
func isConnectionMeantToBeConnected(conn *connections.Connection) bool {
	if conn == nil {
		return false
	}
	switch conn.Status {
	case connections.DISCONNECTED, connections.DELETED, connections.NOTFOUND, connections.IGNORED:
		return false
	default:
		return true
	}
}

func (ca *ConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	userUUID, systemID := machines.ResolveUserAndSystemID(ctx)
	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(systemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	if provider == nil {
		errProvider := ErrConnectAction(fmt.Errorf("failed to retrieve provider from context"))
		eventBuilder.WithMetadata(map[string]interface{}{"error": errProvider})
		return machines.NoOp, eventBuilder.Build(), errProvider
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

	// The layered document (per-connection override over the server-wide
	// default) outranks the materialized meshsync_deployment_mode cache, so a
	// connection that inherits its mode picks up a changed server-wide default
	// on reconnect instead of replaying the mode it last ran in.
	//
	// TODO:
	// maybe not call to viper here and propagate default value from above,
	// f.e. when machine is created
	meshsyncDeploymentMode := connections.ResolveDeploymentMode(
		mergedControllersConfig,
		connection.Metadata,
		connections.MeshsyncDeploymentModeFromString(viper.GetString("MESHSYNC_DEFAULT_DEPLOYMENT_MODE")),
	).Mode

	machinectx.setupLock.Lock()
	machinectx.disconnected = false
	if machinectx.setupInProgress {
		machinectx.retryPending = true
		machinectx.setupLock.Unlock()
		if machinectx.log != nil {
			machinectx.log.Infof("Kubernetes connection setup already in progress for %s, queued retry after current setup finishes", machinectx.K8sContext.Name)
		}
		return machines.NoOp, nil, nil
	}
	machinectx.setupInProgress = true
	machinectx.setupLock.Unlock()

	go func() {
		workerExitedCleanly := false
		defer func() {
			if !workerExitedCleanly {
				machinectx.setupLock.Lock()
				machinectx.setupInProgress = false
				machinectx.setupLock.Unlock()
			}
		}()

		for {
			machinectx.setupLock.Lock()
			if machinectx.disconnected {
				machinectx.setupInProgress = false
				workerExitedCleanly = true
				machinectx.setupLock.Unlock()
				return
			}

			// SetControllersConfig first: AddCtxControllerHandlers constructs the
			// operator controller handler with the Helm chart version this document
			// resolves to (operator.version), and captures it there.
			ctrlHelper := machinectx.MesheryCtrlsHelper.
				SetControllersConfig(mergedControllersConfig).
				SetMeshsyncDeploymentMode(meshsyncDeploymentMode).
				AddCtxControllerHandlers(machinectx.K8sContext).
				UpdateOperatorsStatusMap(machinectx.OperatorTracker).
				DeployUndeployedOperators(machinectx.OperatorTracker, machinectx.K8sContext.ID)

			ctrlHelper.AddMeshsyncDataHandlers(ctx, machinectx.K8sContext, userUUID, systemID, provider)
			machinectx.setupLock.Unlock()

			// Operator mode: best-effort apply of the explicitly-set
			// configuration onto the cluster's MeshSync/Broker custom resources
			// and the MeshSync deployment. An empty resolved configuration is
			// applied too: it withdraws previously-applied fields (cleared
			// while the connection was down) and heals drift. Targets that do
			// not exist yet (the operator is still coming up) are skipped and
			// re-applied on the next connect or configuration change.
			if meshsyncDeploymentMode == connections.MeshsyncDeploymentModeOperator && controllersConfigResolved {
				machinectx.setupLock.Lock()
				isDisc := machinectx.disconnected
				machinectx.setupLock.Unlock()
				if isDisc {
					return
				}

				kubeClient, errClient := machinectx.K8sContext.GenerateKubeHandler()
				if errClient != nil {
					machinectx.log.Error(ErrConnectAction(errClient))
				} else {
					// Detached context: the connect request's context ends with the
					// HTTP request, while this apply runs alongside the async
					// operator deployment.
					applyCtx, cancelApply := context.WithTimeout(context.Background(), 2*time.Minute)
					if _, errApply := models.ApplyControllersConfigToCluster(applyCtx, machinectx.log, kubeClient, mergedControllersConfig); errApply != nil {
						machinectx.log.Error(errApply)
					}
					cancelApply()
				}
			}

			// Atomic check and exit under setupLock
			machinectx.setupLock.Lock()
			if machinectx.disconnected || !machinectx.retryPending {
				machinectx.retryPending = false
				machinectx.setupInProgress = false
				workerExitedCleanly = true
				machinectx.setupLock.Unlock()
				return
			}
			machinectx.retryPending = false
			machinectx.setupLock.Unlock()

			// A retry was queued. Before re-running setup, verify that the
			// connection is still meant to be connected. If it transitioned away
			// (e.g. Disconnect or Delete), abort the queued setup.
			latestConn, _, errConn := provider.GetConnectionByID(token, connectionID)
			if errConn != nil || !isConnectionMeantToBeConnected(latestConn) {
				if machinectx.log != nil {
					machinectx.log.Infof("Skipping queued Kubernetes connection setup for %s: connection is no longer connected", machinectx.K8sContext.Name)
				}
				machinectx.setupLock.Lock()
				machinectx.retryPending = false
				machinectx.setupInProgress = false
				workerExitedCleanly = true
				machinectx.setupLock.Unlock()
				return
			}
		}
	}()
	return machines.NoOp, nil, nil
}

func (ca *ConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
