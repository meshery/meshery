package models

import (
	"sync"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

var (
	once              sync.Once
	registrationQueue *MeshSyncRegistrationQueue
)

// MeshSyncRegistrationQueue holds a channel for queuing Kubernetes resources that are
// eligible to be registered as Meshery connections.
//
// TODO: MeshSync additionally publishes these resources to a broker topic/subject
// (meshsync.registerconnection.queue?); Meshery Server subscribes and performs the
// necessary action.
type MeshSyncRegistrationQueue struct {
	RegChan chan MeshSyncRegistrationData
}

type MeshSyncRegistrationData struct {
	MeshsyncDataHandler MeshsyncDataHandler
	Obj                 meshsyncmodel.KubernetesResource
}

func InitMeshSyncRegistrationQueue() {
	initQueue()
}

func GetMeshSyncRegistrationQueue() *MeshSyncRegistrationQueue {
	initQueue()
	return registrationQueue
}

// initQueue initializes the registration queue exactly once using sync.Once
func initQueue() {
	once.Do(func() {
		registrationQueue = &MeshSyncRegistrationQueue{
			RegChan: make(chan MeshSyncRegistrationData, 10),
		}
	})
}

func (mrq *MeshSyncRegistrationQueue) Send(data MeshSyncRegistrationData) {
	if mrq == nil || mrq.RegChan == nil {
		return
	}
	mrq.RegChan <- data
}
