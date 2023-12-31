package channels

// import (
// 	"sync"

// 	"github.com/layer5io/meshery/server/models/machines"
// 	"github.com/layer5io/meshkit/database"
// 	"github.com/layer5io/meshkit/logger"
// 	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
// 	"github.com/layer5io/meshsync/pkg/model"
// )

// var (
// 	autoRegistrationSingleton *AutoRegistrationChannel
// 	mutex                     sync.RWMutex
// 	connectionCompFilter      = &v1alpha1.ComponentFilter{
// 		Name:       "Connection",
// 		APIVersion: "meshery.layer5.io/v1alpha1",
// 		Greedy:     true,
// 	}
// )

// type MeshSyncDataQueue struct {
// 	Ch chan model.KubernetesResource
// }

// func newMeshSyncDataQueue() MeshSyncDataQueue {
// 	return MeshSyncDataQueue{
// 		Ch: make(chan model.KubernetesResource, 10),
// 	}
// }

// func (arh *AutoRegistrationChannel) Send(obj model.KubernetesResource) {
// 	arh.Queue.Ch <- obj
// }

// type AutoRegistrationChannel struct {
// 	DbHandler         *database.Handler
// 	Queue             MeshSyncDataQueue
// 	Log               logger.Handler
// 	SMInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
// }

// func SetAutoRegistrationHelperSingleton(dbHandler *database.Handler, smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker) {
// 	defer mutex.Unlock()
// 	mutex.Lock()
// 	if autoRegistrationSingleton == nil {
// 		autoRegistrationSingleton = &AutoRegistrationChannel{
// 			DbHandler:         dbHandler,
// 			Queue:             newMeshSyncDataQueue(),
// 			SMInstanceTracker: smInstanceTracker,
// 		}

// 	}
// }

// func GetAutoRegistrationHelperSingleton() *AutoRegistrationChannel {
// 	defer mutex.RUnlock()
// 	mutex.RLock()
// 	return autoRegistrationSingleton
// }
