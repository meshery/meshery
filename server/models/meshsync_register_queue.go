package models

import (
	"sync"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	once              sync.Once
	registrationQueue *MeshSyncRegistrationQueue
)

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
