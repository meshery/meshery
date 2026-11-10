package kubernetes

import (
	"fmt"
	"sync"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

// adapterTracker is initialized lazily on first access so that the ADAPTER_URLS
// env-var read happens AFTER main() has called viper.AutomaticEnv(). A package-
// level eager init would run at import time, before viper sees env vars, and
// silently produce an empty tracker for env-only configurations.
var (
	adapterTrackerOnce sync.Once
	adapterTrackerInst *helpers.AdaptersTracker
)

func getAdapterTracker() *helpers.AdaptersTracker {
	adapterTrackerOnce.Do(func() {
		adapterURLs := utils.SplitAndTrim(viper.GetString("ADAPTER_URLS"), ", \t\n\r")
		adapterTrackerInst = helpers.NewAdaptersTracker(adapterURLs)
	})
	return adapterTrackerInst
}

func GenerateClientSetAction(k8sContext *models.K8sContext, eventBuilder *events.EventBuilder, log logger.Handler) (*kubernetes.Client, error) {
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	eventMetadata := map[string]interface{}{}

	handler, err := models.GenerateK8sClientSet(k8sContext, eventBuilder, eventMetadata, log)
	if handler == nil {
		return nil, err
	}

	return handler, nil
}

func GetMachineCtx(machinectx interface{}, eb *events.EventBuilder) (*MachineCtx, error) {
	machineCtx, ok := machinectx.(*MachineCtx)
	if !ok {
		err := machines.ErrAssertMachineCtx(fmt.Errorf("asserting of context %v failed", machinectx))
		if eb != nil {
			eb.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			})
		}
		return nil, err
	}
	return machineCtx, nil
}

func AssignClientSetToContext(machinectx *MachineCtx, eventBuilder *events.EventBuilder) error {
	if machinectx.clientset != nil {
		return nil
	}

	k8sContext := machinectx.K8sContext
	eventBuilder.ActedUpon(uuid.FromStringOrNil(k8sContext.ConnectionID))

	handler, err := GenerateClientSetAction(&k8sContext, eventBuilder, machinectx.log)
	if err != nil {
		// perofmr event publishinh and err handling
		return err
	}
	machinectx.clientset = handler
	return nil
}

func AssignControllerHandlers(machinectx *MachineCtx, systemID *core.Uuid, provider models.Provider) {
	machinectx.MesheryCtrlsHelper = models.NewMesheryControllersHelper(
		machinectx.log,
		models.NewOperatorDeploymentConfig(getAdapterTracker()),
		models.GetDBInstance(),
		machinectx.EventBroadcaster,
		provider,
		systemID,
	)
}
