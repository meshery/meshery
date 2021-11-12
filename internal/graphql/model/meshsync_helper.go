package model

import (
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	"gorm.io/gorm"
)

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
