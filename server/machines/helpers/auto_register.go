package helpers

import (
	"context"
	"fmt"
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

var (
	once sync.Once
	autoRegistrationSingleton *AutoRegistrationHelper
	connectionCompFilter = &v1alpha1.ComponentFilter{
		Name:       "Connection",
		APIVersion: "meshery.layer5.io/v1alpha1",
		Greedy:     true,
	}
)

type AutoRegistrationHelper struct {
	dbHandler         *database.Handler
	log               logger.Handler
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
}

func InitRegistrationHelperSingleton(dbHandler *database.Handler, log logger.Handler, smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker) {
	once.Do(func() {
		autoRegistrationSingleton = &AutoRegistrationHelper{
			smInstanceTracker: smInstanceTracker,
			log: log,
			dbHandler: dbHandler,
		}
		go autoRegistrationSingleton.processRegistration()
	})
}

func GetAutoRegistrationSingleton() *AutoRegistrationHelper {
	return autoRegistrationSingleton
}

func (arh *AutoRegistrationHelper) processRegistration() {
	if arh == nil {
		return
	}

	regChan := models.GetMeshSyncRegistrationQueue().RegChan

	for regData := range regChan {
		go func(data *models.MeshSyncRegistrationData) {
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
			connType := getTypeOfConnection(&data.Obj)
			if connType != "" {
				id, _ := uuid.NewV4() // id should be hash of somehting.
				ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: data.MeshsyncDataHandler.UserID.String()})
				
				machineInst, err := InitializeMachineWithContext(nil, ctx, id, data.MeshsyncDataHandler.UserID, arh.smInstanceTracker, arh.log, data.MeshsyncDataHandler.Provider, machines.DISCOVERED, connType, nil)
				if err != nil {
					arh.log.Error(ErrAutoRegister(err, connType))
				}
				
				machineInst.SendEvent(ctx, machines.Register, nil)
				machineInst.SendEvent(ctx, machines.Connect, nil)
				
			}
		}(&regData)
	}
}

func getTypeOfConnection(obj *model.KubernetesResource) string {
	if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") && obj.Kind == "Service" {
		fmt.Println("\n\n\n\n", "TEST:::: ---------", obj, "\n\n\n\n")
		return "grafana"
	} else if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") && obj.Kind == "Service" {
		fmt.Println("\n\n\n\n", "TEST:::: ---------", obj, "\n\n\n\n")
		return "prometheus"
	}
	return ""
}
