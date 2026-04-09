package models

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/meshery/schemas/models/core"

	"maps"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/broker"
	channelBroker "github.com/meshery/meshkit/broker/channel"
	"github.com/meshery/meshkit/broker/nats"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	libmeshsync "github.com/meshery/meshsync/pkg/lib/meshsync"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/spf13/viper"
)

const (
	ChartRepo                     = "https://meshery.github.io/meshery.io/charts"
	MesheryServerBrokerConnection = "meshery-server"
)

type MesheryControllerStatusAndVersion struct {
	Status  controllers.MesheryControllerStatus
	Version string
}

type MesheryController int

const (
	MesheryBroker MesheryController = iota
	Meshsync
	MesheryOperator
)

type MesheryControllersHelper struct {
	// contextID stores the current context ID being processed (retained for backward compatibility but maps are used for multi-cluster)
	contextID string
	// ctxControllerHandlers maps contextID to its controller handlers
	ctxControllerHandlers map[string]map[MesheryController]controllers.IMesheryController

	// ctxOperatorStatus maps contextID to its operator status
	ctxOperatorStatus map[string]controllers.MesheryControllerStatus

	// ctxMeshsyncDataHandler maps contextID to its meshsync data handler
	ctxMeshsyncDataHandler map[string]*MeshsyncDataHandler

	mu sync.RWMutex

	log          logger.Handler
	oprDepConfig controllers.OperatorDeploymentConfig
	dbHandler    *database.Handler

	meshsyncDeploymentMode schemasConnection.MeshsyncDeploymentMode

	// event broadcasting dependencies
	eventBroadcaster *Broadcast
	provider         Provider
	systemID         *core.Uuid
}

func (mch *MesheryControllersHelper) GetControllerHandlersForContext(ctxID string) map[MesheryController]controllers.IMesheryController {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	return mch.ctxControllerHandlers[ctxID]
}

func (mch *MesheryControllersHelper) GetControllerHandlersForEachContext() map[MesheryController]controllers.IMesheryController {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	return mch.ctxControllerHandlers[mch.contextID]
}

func (mch *MesheryControllersHelper) GetMeshSyncDataHandlerForContext(ctxID string) *MeshsyncDataHandler {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	return mch.ctxMeshsyncDataHandler[ctxID]
}

// Deprecated: use GetMeshSyncDataHandlerForContext instead
func (mch *MesheryControllersHelper) GetMeshSyncDataHandlersForEachContext() *MeshsyncDataHandler {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	return mch.ctxMeshsyncDataHandler[mch.contextID]
}

func (mch *MesheryControllersHelper) GetOperatorsStatusForContext(ctxID string) controllers.MesheryControllerStatus {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	status, ok := mch.ctxOperatorStatus[ctxID]
	if !ok {
		return controllers.Unknown
	}
	return status
}

// Deprecated: use GetOperatorsStatusForContext instead
func (mch *MesheryControllersHelper) GetOperatorsStatusMap() controllers.MesheryControllerStatus {
	mch.mu.RLock()
	defer mch.mu.RUnlock()
	status, ok := mch.ctxOperatorStatus[mch.contextID]
	if !ok {
		return controllers.Unknown
	}
	return status
}

func NewMesheryControllersHelper(
	log logger.Handler,
	operatorDepConfig controllers.OperatorDeploymentConfig,
	dbHandler *database.Handler,
	eventBroadcaster *Broadcast,
	provider Provider,
	systemID *core.Uuid,
) *MesheryControllersHelper {
	return &MesheryControllersHelper{
		ctxControllerHandlers: make(map[string]map[MesheryController]controllers.IMesheryController),
		ctxOperatorStatus:     make(map[string]controllers.MesheryControllerStatus),
		ctxMeshsyncDataHandler: make(map[string]*MeshsyncDataHandler),
		log:                   log,
		oprDepConfig:          operatorDepConfig,
		dbHandler:              dbHandler,
		meshsyncDeploymentMode: schemasConnection.MeshsyncDeploymentModeOperator,
		eventBroadcaster:       eventBroadcaster,
		provider:               provider,
		systemID:               systemID,
	}
}

func (mch *MesheryControllersHelper) SetMeshsyncDeploymentMode(value schemasConnection.MeshsyncDeploymentMode) *MesheryControllersHelper {
	mch.meshsyncDeploymentMode = value
	return mch
}

// initializes Meshsync data handler for the contexts for whom it has not been
// initialized yet. Apart from updating the map, it also runs the handler after
// updating the map. The presence of a handler for a context in a map indicate that
// the meshsync data for that context is properly being handled
func (mch *MesheryControllersHelper) AddMeshsynDataHandlers(ctx context.Context, k8scontext K8sContext, userID, mesheryInstanceID uuid.UUID, provider Provider) *MesheryControllersHelper {
	ctxID := k8scontext.ID
	mch.mu.Lock()
	mch.contextID = ctxID
	msHandler := mch.ctxMeshsyncDataHandler[ctxID]
	mch.mu.Unlock()

	if msHandler == nil {
		var brokerHandler broker.Handler
		var stopFunc func()

		switch mch.meshsyncDeploymentMode {
		case schemasConnection.MeshsyncDeploymentModeOperator:
			brokerHandler = mch.meshsynDataHandlersNatsBroker(k8scontext, userID)
		case schemasConnection.MeshsyncDeploymentModeEmbedded:
			brokerHandler = channelBroker.NewChannelBrokerHandler()
			// use a standalone context here context.Background(), as
			// meshsync run must be stopped only when meshsync data handler is deregistered
			// and ctx which is passed from above, could be closed earlier
			stop, err := mch.meshsynDataHandlersStartLibMeshsyncRun(context.Background(), brokerHandler, k8scontext, userID)
			if err != nil {
				mch.log.Error(err)
				mch.emitErrorEvent("Failed to start MeshSync library run", err, map[string]any{
					"k8sContextID":           ctxID,
					"k8sContextName":         k8scontext.Name,
					"connectionID":           k8scontext.ConnectionID,
					"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
				}, userID)
				return mch
			}
			stopFunc = stop
		default:
			mch.log.Warnf(
				"MesheryControllersHelper unsupported meshsyncDeploymentMode %s",
				mch.meshsyncDeploymentMode,
			)
			mch.emitWarningEvent("Unsupported MeshSync deployment mode", nil, map[string]any{
				"k8sContextID":           ctxID,
				"k8sContextName":         k8scontext.Name,
				"connectionID":           k8scontext.ConnectionID,
				"meshsyncDeploymentMode": string(mch.meshsyncDeploymentMode),
			}, userID)
			return mch
		}

		if brokerHandler == nil {
			mch.log.Warnf("MesheryControllersHelper::AddMeshsynDataHandlers brokerHandler is nil")
			mch.emitWarningEvent("MeshSync data handler broker is nil", nil, map[string]any{
				"k8sContextID":   ctxID,
				"k8sContextName": k8scontext.Name,
				"connectionID":   k8scontext.ConnectionID,
			}, userID)
			return mch
		}
		token, _ := ctx.Value(TokenCtxKey).(string)
		msDataHandler := NewMeshsyncDataHandler(brokerHandler, *mch.dbHandler, mch.log, provider, userID, uuid.FromStringOrNil(k8scontext.ConnectionID), mesheryInstanceID, token, stopFunc)

		mch.mu.Lock()
		// Double-check pattern to prevent duplicate initialization
		if existingHandler := mch.ctxMeshsyncDataHandler[ctxID]; existingHandler != nil {
			mch.mu.Unlock()
			return mch
		}

		err := msDataHandler.Run()
		if err != nil {
			mch.mu.Unlock()
			mch.log.Warn(err)
			mch.log.Info(fmt.Sprintf("Unable to connect MeshSync for Kubernetes context (%s) due to: %s", ctxID, err.Error()))
			mch.emitErrorEvent("Unable to connect MeshSync", err, map[string]any{
				"k8sContextID":   ctxID,
				"k8sContextName": k8scontext.Name,
				"connectionID":   k8scontext.ConnectionID,
			}, userID)
			return mch
		}
		mch.ctxMeshsyncDataHandler[ctxID] = msDataHandler
		mch.mu.Unlock()
		mch.log.Info(fmt.Sprintf("MeshSync connected for Kubernetes context (%s)", ctxID))
	}

	// Emit success event for successful MeshSync data handler attachment
	var description string
	if mch.meshsyncDeploymentMode != "" {
		description = fmt.Sprintf("MeshSync connected in %s mode", string(mch.meshsyncDeploymentMode))
	} else {
		description = "MeshSync connected"
	}
	mch.emitEvent(
		description,
		events.Informational,
		map[string]any{
			"k8sContextID":           k8scontext.ID,
			"k8sContextName":         k8scontext.Name,
			"connectionID":           k8scontext.ConnectionID,
			"meshsyncDeploymentMode": string(mch.meshsyncDeploymentMode),
		},
		userID,
	)

	return mch
}

func (mch *MesheryControllersHelper) meshsynDataHandlersNatsBroker(
	k8scontext K8sContext,
	userID core.Uuid,
) broker.Handler {
	ctxID := k8scontext.ID
	mch.mu.RLock()
	controllerHandlers := mch.ctxControllerHandlers[ctxID]
	mch.mu.RUnlock()

	if controllerHandlers == nil || controllerHandlers[MesheryBroker] == nil {
		return nil
	}

	// brokerStatus := controllerHandlers[MesheryBroker].GetStatus()
	// do something if broker is being deployed , maybe try again after sometime
	brokerEndpoint, err := controllerHandlers[MesheryBroker].GetPublicEndpoint()
	if brokerEndpoint == "" {
		if err != nil {
			mch.log.Warn(err)
			mch.emitWarningEvent("Failed to get Meshery Broker endpoint", err, map[string]any{
				"k8sContextID":   ctxID,
				"k8sContextName": k8scontext.Name,
				"connectionID":   k8scontext.ConnectionID,
			}, userID)
		}
		mch.log.Info(
			fmt.Sprintf("Meshery Broker unreachable for Kubernetes context (%v)", ctxID),
		)
		mch.emitWarningEvent("Meshery Broker unreachable", nil, map[string]any{
			"k8sContextID":   ctxID,
			"k8sContextName": k8scontext.Name,
			"connectionID":   k8scontext.ConnectionID,
		}, userID)
		return nil
	}
	brokerHandler, err := nats.New(nats.Options{
		// URLS: []string{"localhost:4222"},
		URLS:           []string{brokerEndpoint},
		ConnectionName: MesheryServerBrokerConnection,
		Username:       "",
		Password:       "",
		ReconnectWait:  2 * time.Second,
		MaxReconnect:   60,
	})

	if err != nil {
		mch.log.Warn(err)
		mch.log.Info(fmt.Sprintf("MeshSync not configured for Kubernetes context (%v) due to '%v'", ctxID, err.Error()))
		mch.emitWarningEvent("Failed to connect to Meshery Broker", err, map[string]any{
			"k8sContextID":   ctxID,
			"k8sContextName": k8scontext.Name,
			"connectionID":   k8scontext.ConnectionID,
			"brokerEndpoint": brokerEndpoint,
		}, userID)
		return nil
	}
	mch.log.Info(fmt.Sprintf("Connected to Meshery Broker (%v) for Kubernetes context (%v)", brokerEndpoint, ctxID))
	return brokerHandler
}

// meshsynDataHandlersStartLibMeshsyncRun starts the libmeshsync run for the given context.
// returns stop function to stop goroutine
func (mch *MesheryControllersHelper) meshsynDataHandlersStartLibMeshsyncRun(
	ctx context.Context,
	brokerHandler broker.Handler,
	k8sContext K8sContext,
	userID core.Uuid,
) (func(), error) {
	kubeConfig, err := k8sContext.GenerateKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("MesheryControllersHelper::meshsynDataHandlersStartLibMeshsyncRun error generating kubeconfig from context: %v", err)
	}

	cancelCtx, stopFunc := context.WithCancel(ctx)

	go func() {
		if err := libmeshsync.Run(
			mch.log,
			libmeshsync.WithOutputMode("broker"),
			libmeshsync.WithBrokerHandler(brokerHandler),
			libmeshsync.WithKubeConfig(kubeConfig),
			libmeshsync.WithContext(cancelCtx),
		); err != nil {
			meshsyncErr := fmt.Errorf("MesheryControllersHelper::meshsynDataHandlersStartLibMeshsyncRun error running meshsync lib: %v", err)
			mch.log.Error(meshsyncErr)
			mch.emitErrorEvent("Error running MeshSync library", meshsyncErr, map[string]any{
				"k8sContextID":           k8sContext.ID,
				"k8sContextName":         k8sContext.Name,
				"connectionID":           k8sContext.ConnectionID,
				"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
			}, userID)
		}
	}()

	return stopFunc, nil
}

func (mch *MesheryControllersHelper) RemoveMeshSyncDataHandler(ctx context.Context, contextID string) {
	mch.mu.Lock()
	defer mch.mu.Unlock()
	if msHandler, ok := mch.ctxMeshsyncDataHandler[contextID]; ok && msHandler != nil {
		mch.log.Infof("MesheryControllersHelper::RemoveMeshSyncDataHandler for contextID = %s", contextID)
		msHandler.Stop()
		delete(mch.ctxMeshsyncDataHandler, contextID)
	}
}

func (mch *MesheryControllersHelper) ResyncMeshsync(ctx context.Context, contextID string) error {
	mch.mu.RLock()
	msHandler := mch.ctxMeshsyncDataHandler[contextID]
	mch.mu.RUnlock()
	if msHandler != nil {
		return msHandler.Resync()
	}
	return nil
}

// attach a MesheryController for each context if
// 1. the config is valid
// 2. if it is not already attached
func (mch *MesheryControllersHelper) AddCtxControllerHandlers(ctx K8sContext, userID uuid.UUID) *MesheryControllersHelper {
	mch.mu.Lock()
	mch.contextID = ctx.ID
	mch.mu.Unlock()

	cfg, err := ctx.GenerateKubeConfig()
	if err != nil {
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to generate kubeconfig", err, map[string]any{
			"k8sContextID":   ctx.ID,
			"k8sContextName": ctx.Name,
			"connectionID":   ctx.ConnectionID,
		}, userID)
		return mch
	}

	client, err := mesherykube.New(cfg)
	// means that the config is invalid
	if err != nil {
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to create Kubernetes client", err, map[string]any{
			"k8sContextID":   ctx.ID,
			"k8sContextName": ctx.Name,
			"connectionID":   ctx.ConnectionID,
		}, userID)
		return mch
	}

	mch.mu.Lock()
	mch.ctxControllerHandlers[ctx.ID] = map[MesheryController]controllers.IMesheryController{
		MesheryBroker:   controllers.NewMesheryBrokerHandler(client),
		MesheryOperator: controllers.NewMesheryOperatorHandler(client, mch.oprDepConfig),
		Meshsync:        controllers.NewMeshsyncHandler(client),
	}
	mch.mu.Unlock()

	return mch
}

func (mch *MesheryControllersHelper) RemoveCtxControllerHandler(ctx context.Context, contextID string) {
	mch.mu.Lock()
	defer mch.mu.Unlock()
	delete(mch.ctxControllerHandlers, contextID)
}

// update the status of MesheryOperator in all the contexts
// for whom MesheryControllers are attached
// should be called after AddCtxControllerHandlers
func (mch *MesheryControllersHelper) UpdateOperatorsStatusMap(ot *OperatorTracker, contextID string) *MesheryControllersHelper {
	if mch.meshsyncDeploymentMode != schemasConnection.MeshsyncDeploymentModeOperator {
		return mch
	}

	if ot.IsUndeployed(contextID) {
		mch.mu.Lock()
		mch.ctxOperatorStatus[contextID] = controllers.Undeployed
		mch.mu.Unlock()
	} else {
		mch.mu.RLock()
		handlers := mch.ctxControllerHandlers[contextID]
		mch.mu.RUnlock()
		if handlers != nil {
			operatorHandler, ok := handlers[MesheryOperator]
			if ok {
				status := operatorHandler.GetStatus()
				mch.mu.Lock()
				mch.ctxOperatorStatus[contextID] = status
				mch.mu.Unlock()
			}
		}
	}

	return mch
}

type OperatorTracker struct {
	ctxIDtoDeploymentStatus map[string]bool
	mx                      sync.Mutex
	DisableOperator         bool
}

func NewOperatorTracker(disabled bool) *OperatorTracker {
	return &OperatorTracker{
		ctxIDtoDeploymentStatus: make(map[string]bool),
		mx:                      sync.Mutex{},
		DisableOperator:         disabled,
	}
}

func (ot *OperatorTracker) Undeployed(ctxID string, undeployed bool) {
	if ot.DisableOperator { //no-op when operator is disabled
		return
	}
	ot.mx.Lock()
	defer ot.mx.Unlock()
	if ot.ctxIDtoDeploymentStatus == nil {
		ot.ctxIDtoDeploymentStatus = make(map[string]bool)
	}
	ot.ctxIDtoDeploymentStatus[ctxID] = undeployed
}

func (ot *OperatorTracker) IsUndeployed(ctxID string) bool {
	if ot.DisableOperator { //Return true everytime so that operators stay in undeployed state across all contexts
		return true
	}
	ot.mx.Lock()
	defer ot.mx.Unlock()
	if ot.ctxIDtoDeploymentStatus == nil {
		ot.ctxIDtoDeploymentStatus = make(map[string]bool)
		return false
	}
	return ot.ctxIDtoDeploymentStatus[ctxID]
}

// looks at the status of Meshery Operator for each cluster and takes necessary action.
// it will deploy the operator only when it is in NotDeployed state
func (mch *MesheryControllersHelper) DeployUndeployedOperators(ot *OperatorTracker, contextID string) *MesheryControllersHelper {
	if ot.DisableOperator { //Return true everytime so that operators stay in undeployed state across all contexts
		return mch
	}
	if mch.meshsyncDeploymentMode != schemasConnection.MeshsyncDeploymentModeOperator {
		return mch
	}

	mch.mu.RLock()
	oprStatus := mch.ctxOperatorStatus[contextID]
	handlers := mch.ctxControllerHandlers[contextID]
	mch.mu.RUnlock()

	if oprStatus == controllers.NotDeployed {
		if handlers != nil {
			operatorHandler, ok := handlers[MesheryOperator]
			if ok {
				err := operatorHandler.Deploy(false)

				if err != nil {
					mch.log.Error(err)
					mch.emitErrorEvent("Failed to deploy Meshery Operator", err, map[string]any{
						"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
						"operatorStatus":         oprStatus,
						"contextID":              contextID,
					}, uuid.Nil)
				}
			}
		}
	}

	return mch
}

func (mch *MesheryControllersHelper) UndeployDeployedOperators(ot *OperatorTracker, contextID string) *MesheryControllersHelper {
	mch.mu.RLock()
	oprStatus := mch.ctxOperatorStatus[contextID]
	handlers := mch.ctxControllerHandlers[contextID]
	mch.mu.RUnlock()

	if oprStatus != controllers.Undeployed {
		if handlers != nil {
			operatorHandler, ok := handlers[MesheryOperator]
			if ok {
				err := operatorHandler.Undeploy()

				if err != nil {
					mch.log.Error(err)
					mch.emitErrorEvent("Failed to undeploy Meshery Operator", err, map[string]any{
						"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
						"operatorStatus":         oprStatus,
						"contextID":              contextID,
					}, uuid.Nil)
				}
			}
		}
	}

	return mch
}

func NewOperatorDeploymentConfig(adapterTracker AdaptersTrackerInterface) controllers.OperatorDeploymentConfig {
	// get meshery release version
	mesheryReleaseVersion := viper.GetString("BUILD")
	if mesheryReleaseVersion == "" || mesheryReleaseVersion == "Not Set" || mesheryReleaseVersion == "edge-latest" {
		_, latestRelease, err := CheckLatestVersion(mesheryReleaseVersion)
		// if unable to fetch latest release tag, meshkit helm functions handle
		// this automatically fetch the latest one
		if err != nil {
			mesheryReleaseVersion = ""
		} else {
			mesheryReleaseVersion = latestRelease
		}
	}

	return controllers.OperatorDeploymentConfig{
		MesheryReleaseVersion: mesheryReleaseVersion,
		GetHelmOverrides: func(delete bool) map[string]interface{} {
			return setOverrideValues(delete, adapterTracker)
		},
		HelmChartRepo: ChartRepo,
	}
}

// checkLatestVersion takes in the current server version compares it with the target
// and returns the (isOutdated, latestVersion, error)
func CheckLatestVersion(serverVersion string) (*bool, string, error) {
	// Inform user of the latest release version
	versions, err := utils.GetLatestReleaseTagsSorted("meshery", "meshery")
	latestVersion := versions[len(versions)-1]
	isOutdated := false
	if err != nil {
		return nil, "", ErrCreateOperatorDeploymentConfig(err)
	}
	// Compare current running Meshery server version to the latest available Meshery release on GitHub.
	if latestVersion != serverVersion {
		isOutdated = true
		return &isOutdated, latestVersion, nil
	}

	return &isOutdated, latestVersion, nil
}

// setOverrideValues detects the currently insalled adapters and sets appropriate
// overrides so as to not uninstall them. It also sets override values for
// operator so that it can be enabled or disabled depending on the need
func setOverrideValues(delete bool, adapterTracker AdaptersTrackerInterface) map[string]interface{} {
	installedAdapters := make([]string, 0)
	adapters := adapterTracker.GetAdapters(context.TODO())

	for _, adapter := range adapters {
		if adapter.Name != "" {
			installedAdapters = append(installedAdapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{}{
		"fullnameOverride": "meshery-operator",
		"meshery": map[string]interface{}{
			"enabled": false,
		},
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-operator": map[string]interface{}{
			"enabled": true,
		},
	}

	for _, adapter := range installedAdapters {
		if _, ok := overrideValues[adapter]; ok {
			overrideValues[adapter] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	if delete {
		overrideValues["meshery-operator"] = map[string]interface{}{
			"enabled": false,
		}
	}

	return overrideValues
}

// setOverrideValues detects the currently insalled adapters and sets appropriate
// overrides so as to not uninstall them.
func SetOverrideValuesForMesheryDeploy(adapters []Adapter, adapter Adapter, install bool) map[string]interface{} {
	installedAdapters := make([]string, 0)
	for _, adapter := range adapters {
		if adapter.Name != "" {
			installedAdapters = append(installedAdapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{}{
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
	}

	for _, adapter := range installedAdapters {
		if _, ok := overrideValues[adapter]; ok {
			overrideValues[adapter] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	// based on deploy/undeploy action change the status of adapter override
	if _, ok := overrideValues[strings.Split(adapter.Location, ":")[0]]; ok {
		overrideValues[strings.Split(adapter.Location, ":")[0]] = map[string]interface{}{
			"enabled": install,
		}
	}

	return overrideValues
}

// General helper method to emit events for system-level operations
func (mch *MesheryControllersHelper) emitEvent(description string, severity events.EventSeverity, metadata map[string]any, userID core.Uuid) {
	if mch.eventBroadcaster != nil && mch.systemID != nil {
		actedUpon := controllerEventActedUpon(userID, metadata)

		eventBuilder := events.NewEvent().
			FromSystem(*mch.systemID).
			WithCategory("connection").
			WithAction("update").
			WithSeverity(severity).
			WithDescription(description).
			WithMetadata(metadata)

		if userID != uuid.Nil {
			eventBuilder = eventBuilder.FromUser(userID)
		}

		if actedUpon != uuid.Nil {
			eventBuilder = eventBuilder.ActedUpon(actedUpon)
		}

		event := eventBuilder.Build()

		if mch.provider != nil {
			if shouldPersistControllerEvent(userID, actedUpon) {
				if err := mch.provider.PersistSystemEvent(*event); err != nil {
					mch.log.Error(fmt.Errorf("failed to persist event: %w", err))
				}
			} else {
				mch.log.Debug("skipping remote persistence for controller event without valid user or resource association")
			}
		}
		go mch.eventBroadcaster.Publish(userID, event)
	}
}

func controllerEventActedUpon(userID core.Uuid, metadata map[string]any) core.Uuid {
	if metadata != nil {
		switch connectionID := metadata["connectionID"].(type) {
		case string:
			if parsedID := uuid.FromStringOrNil(connectionID); parsedID != uuid.Nil {
				return parsedID
			}
		case core.Uuid:
			if connectionID != uuid.Nil {
				return connectionID
			}
		case *core.Uuid:
			if connectionID != nil && *connectionID != uuid.Nil {
				return *connectionID
			}
		}
	}

	if userID != uuid.Nil {
		return userID
	}

	return uuid.Nil
}

func shouldPersistControllerEvent(userID, actedUpon core.Uuid) bool {
	return userID != uuid.Nil && actedUpon != uuid.Nil
}

// Common helper for both error and warning events with error information
func (mch *MesheryControllersHelper) emitEventWithError(description string, severity events.EventSeverity, err error, metadata map[string]any, userID core.Uuid) {
	eventMetadata := make(map[string]any)

	if metadata != nil {
		maps.Copy(eventMetadata, metadata)
	}

	if err != nil {
		eventMetadata["error"] = err.Error()
	}

	mch.emitEvent(description, severity, eventMetadata, userID)
}

// Helper method to emit error events
func (mch *MesheryControllersHelper) emitErrorEvent(description string, err error, metadata map[string]any, userID core.Uuid) {
	mch.emitEventWithError(description, events.Error, err, metadata, userID)
}

// Helper method to emit warning events
func (mch *MesheryControllersHelper) emitWarningEvent(description string, err error, metadata map[string]any, userID core.Uuid) {
	mch.emitEventWithError(description, events.Warning, err, metadata, userID)
}
