package model

import (
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	"gorm.io/gorm"
)

const (
	meshsyncHelmChart = "meshery-meshsync"
	meshsyncYaml      = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func RunMeshSync(client *mesherykube.Client, delete bool) error {
	err := installUsingHelm(client, delete, meshsyncHelmChart)
	if err != nil {
		return ErrInstallUsingHelm(err)
	}
	return nil
}

func recordMeshSyncData(eventtype broker.EventType, handler *database.Handler, object *meshsyncmodel.Object) error {
	if handler == nil {
		return ErrEmptyHandler
	}

	handler.Lock()
	defer handler.Unlock()

	switch eventtype {
	case broker.Add, broker.Update:
		result := handler.Create(object)
		if result.Error != nil {
			result = handler.Session(&gorm.Session{FullSaveAssociations: true}).Updates(object)
			if result.Error != nil {
				return ErrCreateData(result.Error)
			}
		}
	case broker.Delete:
		result := handler.Delete(object)
		if result.Error != nil {
			return ErrDeleteData(result.Error)
		}
	case broker.ErrorEvent:
		return nil
	}

	return nil
}
