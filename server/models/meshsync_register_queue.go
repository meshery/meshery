package models

import (
	"sync"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var(
    mutex                        sync.Mutex
    RegistrationQueue *MeshSyncRegistrationQueue

)

type MeshSyncRegistrationQueue struct {
    RegistrationChan chan meshsyncmodel.KubernetesResource
}

func InitMeshSyncRegistrationQueue() {
    sync.OnceFunc(func() {
        RegistrationQueue =  &MeshSyncRegistrationQueue{
        RegistrationChan: make(chan meshsyncmodel.KubernetesResource, 10),
    }
    })
}

func (arh *MeshSyncRegistrationQueue) Send(obj meshsyncmodel.KubernetesResource) {
    arh.RegistrationChan <- obj
}