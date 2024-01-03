package machines

// import (
// 	"context"
// 	"strings"
// 	"sync"

// 	"github.com/gofrs/uuid"
// 	"github.com/layer5io/meshery/server/machines/helpers"
// 	"github.com/layer5io/meshkit/database"
// 	"github.com/layer5io/meshkit/logger"
// 	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
// 	"github.com/layer5io/meshsync/pkg/model"
// )

// var (
//     autoRegistrationSingleton *AutoRegistrationHelper
//     mutex                        sync.Mutex
//     connectionCompFilter      = &v1alpha1.ComponentFilter{
//         Name:       "Connection",
//         APIVersion: "meshery.layer5.io/v1alpha1",
//         Greedy:     true,
//     }
// )

// type MeshSyncDataQueue struct {
//     ch chan model.KubernetesResource
// }

// func newMeshSyncDataQueue() MeshSyncDataQueue {
//     return MeshSyncDataQueue{
//         ch: make(chan model.KubernetesResource, 10),
//     }
// }

// // func (mdq *MeshSynMeshSyncDataQueue) Send(object model.KubernetesResource) {
// //  mdq.ch <- object
// // }
// // func (mdq *MeshSyncDataQueue)

// func (arh *AutoRegistrationHelper) Send(obj model.KubernetesResource) {
//     arh.queue.ch <- obj
// }

// type AutoRegistrationHelper struct {
//     dbHandler *database.Handler
//     queue     MeshSyncDataQueue
//     log       logger.Handler
// }

// func SetAutoRegistrationHelperSingleton(dbHandler *database.Handler) *AutoRegistrationHelper {
//     defer mutex.Unlock()
//     mutex.Lock()
//     if autoRegistrationSingleton == nil {
//         autoRegistrationSingleton = &AutoRegistrationHelper{
//             dbHandler: dbHandler,
//             queue:     newMeshSyncDataQueue(),
//         }
//         go autoRegistrationSingleton.processRegistration()
//     }
//     return autoRegistrationSingleton
// }


// func(arh *AutoRegistrationHelper) processRegistration() {
//     arh = autoRegistrationSingleton
// 	if arh == nil {
// 		return
// 	}

// 	for obj := range arh.queue.ch {
// 		go func(obj *model.KubernetesResource) {
// 			// Ideally iterate all Connection defs, extract fingerprint composite key and try to match with the given obj,
// 			// For all connections that match the fingerprint and autoRegsiter is set to true, try to do auto registration.

// 			// connectionDefs, _, _ := v1alpha1.GetMeshModelComponents(&arh.dbHandler, *connectionCompFilter)
// 			// for _, connectionDef := range connectionDefs {
// 			// 	capabilities, err := utils.Cast[map[string]interface{}](connectionDef.Metadata["capabilities"])
// 			// 	if err != nil {
// 			// 		arh.log.Error(err)
// 			// 		continue
// 			// 	}
// 			// 	autoRegister, ok := capabilities["autoRegister"].(bool)
// 			// 	if ok && autoRegister {
// 			// 		fmt.Println("TEST:: inside for loop extracted capabilities.autoRegister")
// 			// 		// ch
// 			// 	}
// 			// }

// 			// For now, the auto-registration for Prometheus/Grafana is hard-coded.
// 			connType := getTypeOfConnection(obj)
// 			if connType != "" {
// 				id, _ := uuid.NewV4() // id should be hash of somehting.
// 				machineInst, err := helpers.InitializeMachineWithContext(nil, context.TODO(), id, arh.SMInstanceTracker, arh.Log, nil, machines.DISCOVERED, connType, nil)
// 				if err != nil {
// 					arh.log.Error(ErrAutoRegister(err, connType))
// 				}
// 				machineInst.Provider = nil // set provider somehow
// 			}
// 		}(&obj)
// 	}
// }


// func getTypeOfConnection(obj *model.KubernetesResource) string {
// 	if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
// 		return "grafana"
// 	} else if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
// 		return "prometheus"
// 	}
// 	return ""
// }