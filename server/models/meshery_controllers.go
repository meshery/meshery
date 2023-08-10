package models

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/controllers"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

const (
	ChartRepo                     = "https://meshery.github.io/meshery.io/charts"
	MesheryServerBrokerConnection = "meshery-server"
)

type MesheryController int

const (
	MesheryBroker MesheryController = iota
	Meshsync
	MesheryOperator
)

type MesheryControllersHelper struct {
	//  maps each context with the controller handlers
	// this map will be used as the source of truth
	ctxControllerHandlersMap map[string]map[MesheryController]controllers.IMesheryController
	// maps each context with it's operator status
	ctxOperatorStatusMap map[string]controllers.MesheryControllerStatus
	// maps each context with a meshsync data handler
	ctxMeshsyncDataHandlerMap map[string]MeshsyncDataHandler

	mu sync.Mutex

	log          logger.Handler
	oprDepConfig controllers.OperatorDeploymentConfig
	dbHandler    *database.Handler
}

func (mch *MesheryControllersHelper) GetControllerHandlersForEachContext() map[string]map[MesheryController]controllers.IMesheryController {
	return mch.ctxControllerHandlersMap
}

func (mch *MesheryControllersHelper) GetMeshSyncDataHandlersForEachContext() map[string]MeshsyncDataHandler {
	return mch.ctxMeshsyncDataHandlerMap
}

func (mch *MesheryControllersHelper) GetOperatorsStatusMap() map[string]controllers.MesheryControllerStatus {
	return mch.ctxOperatorStatusMap
}

func NewMesheryControllersHelper(log logger.Handler, operatorDepConfig controllers.OperatorDeploymentConfig, dbHandler *database.Handler) *MesheryControllersHelper {
	return &MesheryControllersHelper{
		ctxControllerHandlersMap:  make(map[string]map[MesheryController]controllers.IMesheryController),
		log:                       log,
		oprDepConfig:              operatorDepConfig,
		ctxOperatorStatusMap:      make(map[string]controllers.MesheryControllerStatus),
		ctxMeshsyncDataHandlerMap: make(map[string]MeshsyncDataHandler),
		dbHandler:                 dbHandler,
	}
}

// initializes Meshsync data handler for the contexts for whom it has not been
// initialized yet. Apart from updating the map, it also runs the handler after
// updating the map. The presence of a handler for a context in a map indicate that
// the meshsync data for that context is properly being handled
func (mch *MesheryControllersHelper) UpdateMeshsynDataHandlers() *MesheryControllersHelper {
	// only checking those contexts whose MesheryConrollers are active
	mch.mu.Lock()
	defer mch.mu.Unlock()
	for ctxID, controllerHandlers := range mch.ctxControllerHandlersMap {
		if _, ok := mch.ctxMeshsyncDataHandlerMap[ctxID]; !ok {
			// brokerStatus := controllerHandlers[MesheryBroker].GetStatus()
			// do something if broker is being deployed , maybe try again after sometime
			brokerEndpoint, err := controllerHandlers[MesheryBroker].GetPublicEndpoint()
			if brokerEndpoint == "" {
				if err != nil {
					mch.log.Warn(err)
				}
				mch.log.Info(fmt.Sprintf("skipping meshsync data handler setup for contextId: %v as its public endpoint could not be found", ctxID))
				continue
			}
			mch.log.Info(fmt.Sprintf("found meshery-broker endpoint: %s for contextId: %s", brokerEndpoint, ctxID))
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
				mch.log.Info(fmt.Sprintf("skipping meshsync data handler setup for contextId: %v due to: %v", ctxID, err.Error()))
				continue
			}
			mch.log.Info(fmt.Sprintf("broker connection successfully established for contextId: %v with meshery-broker at: %v", ctxID, brokerEndpoint))
			msDataHandler := NewMeshsyncDataHandler(brokerHandler, *mch.dbHandler, mch.log)
			err = msDataHandler.Run()
			if err != nil {
				mch.log.Warn(err)
				mch.log.Info(fmt.Sprintf("skipping meshsync data handler setup for contextId: %s due to: %s", ctxID, err.Error()))
				continue
			}
			mch.ctxMeshsyncDataHandlerMap[ctxID] = *msDataHandler
			mch.log.Info(fmt.Sprintf("meshsync data handler successfully setup for contextId: %s", ctxID))
		}
	}

	return mch
}

// attach a MesheryController for each context if
// 1. the config is valid
// 2. if it is not already attached
func (mch *MesheryControllersHelper) UpdateCtxControllerHandlers(ctxs []K8sContext) *MesheryControllersHelper {
	mch.mu.Lock()
	defer mch.mu.Unlock()
	// resetting this value as a specific controller handler instance does not have any significance opposed to
	// a MeshsyncDataHandler instance where it signifies whether or not a listener is attached
	mch.ctxControllerHandlersMap = make(map[string]map[MesheryController]controllers.IMesheryController)
	for _, ctx := range ctxs {
		ctxID := ctx.ID
		cfg, _ := ctx.GenerateKubeConfig()
		client, err := mesherykube.New(cfg)
		// means that the config is invalid
		if err != nil {
			// invalid configs are not added to the map
			continue
		}
		mch.ctxControllerHandlersMap[ctxID] = map[MesheryController]controllers.IMesheryController{
			MesheryBroker:   controllers.NewMesheryBrokerHandler(client),
			MesheryOperator: controllers.NewMesheryOperatorHandler(client, mch.oprDepConfig),
			Meshsync:        controllers.NewMeshsyncHandler(client),
		}
	}

	return mch
}

// update the status of MesheryOperator in all the contexts
// for whom MesheryControllers are attached
// should be called after UpdateCtxControllerHandlers
func (mch *MesheryControllersHelper) UpdateOperatorsStatusMap(ot *OperatorTracker) *MesheryControllersHelper {
	mch.mu.Lock()
	defer mch.mu.Unlock()
	mch.ctxOperatorStatusMap = make(map[string]controllers.MesheryControllerStatus)
	for ctxID, ctrlHandler := range mch.ctxControllerHandlersMap {
		if ot.IsUndeployed(ctxID) {
			mch.ctxOperatorStatusMap[ctxID] = controllers.Undeployed
		} else {
			mch.ctxOperatorStatusMap[ctxID] = ctrlHandler[MesheryOperator].GetStatus()
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
func (mch *MesheryControllersHelper) DeployUndeployedOperators(ot *OperatorTracker) *MesheryControllersHelper {
	if ot.DisableOperator { //Return true everytime so that operators stay in undeployed state across all contexts
		return mch
	}
	go func(mch *MesheryControllersHelper) {
		mch.mu.Lock()
		defer mch.mu.Unlock()
		for ctxID, ctrlHandler := range mch.ctxControllerHandlersMap {
			if oprStatus, ok := mch.ctxOperatorStatusMap[ctxID]; ok {
				if oprStatus == controllers.NotDeployed {
					err := ctrlHandler[MesheryOperator].Deploy(false)
					if err != nil {
						mch.log.Error(err)
					}
				}
			}
		}
	}(mch)

	return mch
}

func NewOperatorDeploymentConfig(adapterTracker AdaptersTrackerInterface) controllers.OperatorDeploymentConfig {
	// get meshery release version
	mesheryReleaseVersion := viper.GetString("BUILD")
	if mesheryReleaseVersion == "" || mesheryReleaseVersion == "Not Set" || mesheryReleaseVersion == "edge-latest" {
		_, latestRelease, err := checkLatestVersion(mesheryReleaseVersion)
		// if unable to fetch latest release tag, meshkit helm functions handle
		// this automatically fetch the latest one
		if err != nil {
			// mch.log.Error(fmt.Errorf("Couldn't check release tag: %s. Will use latest version", err))
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
func checkLatestVersion(serverVersion string) (*bool, string, error) {
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

