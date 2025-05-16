package models

import (
	"sync"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	once              sync.Once
	registrationQueue *MeshSyncRegistrationQueue
)

// Configure MeshSync to additionally publish the resources that can be registered as connection to other broker topic/subject (meshsync.registerconnection.queue?).
// Meshery Server subscribes to that topic/subject and performs the necessary action.
type MeshSyncRegistrationQueue struct {
	RegChan chan MeshSyncRegistrationData
}

type MeshSyncRegistrationData struct {
	MeshsyncDataHandler MeshsyncDataHandler
	Obj                 meshsyncmodel.KubernetesResource
}

func InitMeshSyncRegistrationQueue() {
	once.Do(func() {
		registrationQueue = &MeshSyncRegistrationQueue{
			RegChan: make(chan MeshSyncRegistrationData, 10),
		}
	})
}

func GetMeshSyncRegistrationQueue() *MeshSyncRegistrationQueue {
	return registrationQueue
}

func (mrq *MeshSyncRegistrationQueue) Send(data MeshSyncRegistrationData) {
	mrq.RegChan <- data
}
