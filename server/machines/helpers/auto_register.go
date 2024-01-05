package helpers

import (
	"context"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshsync/pkg/model"
)

var AutoRegistrationSingleton *AutoRegistrationHelper

var (
    mutex                        sync.Mutex
    connectionCompFilter      = &v1alpha1.ComponentFilter{
        Name:       "Connection",
        APIVersion: "meshery.layer5.io/v1alpha1",
        Greedy:     true,
    }
)

type AutoRegistrationHelper struct {
    dbHandler *database.Handler
    log       logger.Handler
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
	// meshsyncDataHandler *models.MeshsyncDataHandler
}

func NewRegistrationHelper(dbHandler *database.Handler, log logger.Handler, snInstanceTracker *machines.ConnectionToStateMachineInstanceTracker) *AutoRegistrationHelper {
    defer mutex.Unlock()
    mutex.Lock()
    if AutoRegistrationSingleton == nil {
        AutoRegistrationSingleton = &AutoRegistrationHelper{
            dbHandler: dbHandler,
			log: log,
			// meshsyncDataHandler: mdh,
        }
        go AutoRegistrationSingleton.processRegistration()
    }
    return AutoRegistrationSingleton
}


func(arh *AutoRegistrationHelper) processRegistration() {
	if arh == nil {
		return
	}

	for obj := range models.RegistrationQueue.RegistrationChan {
		go func(obj *model.KubernetesResource) {
			// Ideally iterate all Connection defs, extract fingerprint composite key and try to match with the given obj,
			// For all connections that match the fingerprint and autoRegsiter is set to true, try to do auto registration.

			// connectionDefs, _, _ := v1alpha1.GetMeshModelComponents(&arh.dbHandler, *connectionCompFilter)
			// for _, connectionDef := range connectionDefs {
			// 	capabilities, err := utils.Cast[map[string]interface{}](connectionDef.Metadata["capabilities"])
			// 	if err != nil {
			// 		arh.log.Error(err)
			// 		continue
			// 	}
			// 	autoRegister, ok := capabilities["autoRegister"].(bool)
			// 	if ok && autoRegister {
			// 		fmt.Println("TEST:: inside for loop extracted capabilities.autoRegister")
			// 		// ch
			// 	}
			// }

			// For now, the auto-registration for Prometheus/Grafana is hard-coded.
			connType := getTypeOfConnection(obj)
			if connType != "" {
				id, _ := uuid.NewV4() // id should be hash of somehting.
				userID := uuid.Nil // use proper user id
				machineInst, err := InitializeMachineWithContext(nil, context.TODO(), id, userID, arh.smInstanceTracker, arh.log, "", connType, nil)
				if err != nil {
					// arh.log.Error(ErrAutoRegister(err, connType))
				}
				machineInst.Provider = nil // set provider somehow
			}
		}(&obj)
	}
}


func getTypeOfConnection(obj *model.KubernetesResource) string {
	if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
		return "grafana"
	} else if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
		return "prometheus"
	}
	return ""
}