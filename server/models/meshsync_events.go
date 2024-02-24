package models

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"gorm.io/gorm"
)

const (
	MeshsyncStoreUpdatesSubject = "meshery-server.meshsync.store"
	MeshsyncRequestSubject      = "meshery.meshsync.request"
)

// TODO: Create proper error codes for the functionalities this struct implements
type MeshsyncDataHandler struct {
	broker       broker.Handler
	dbHandler    database.Handler
	log          logger.Handler
	Provider     Provider
	UserID       uuid.UUID
	ConnectionID uuid.UUID
	InstanceID   uuid.UUID
	Token        string
}

func NewMeshsyncDataHandler(broker broker.Handler, dbHandler database.Handler, log logger.Handler, provider Provider, userID, connID, instanceID uuid.UUID, token string) *MeshsyncDataHandler {
	return &MeshsyncDataHandler{
		broker:       broker,
		dbHandler:    dbHandler,
		log:          log,
		Provider:     provider,
		UserID:       userID,
		ConnectionID: connID,
		InstanceID:   instanceID,
		Token:        token,
	}
}

func (mh *MeshsyncDataHandler) Run() error {
	storeSubscriptionStatusChan := make(chan bool)
	// this subscription is independent of whether or not the stale data in the database have been cleaned up
	go mh.subscribeToMeshsyncEvents()

	go mh.subsribeToStoreUpdates(storeSubscriptionStatusChan)
	// to make sure that we don't ask for data before we start listening
	if <-storeSubscriptionStatusChan {
		// err := mh.removeStaleObjects()
		// if err != nil {
		//  return err
		// }
		err := mh.requestMeshsyncStore()
		if err != nil {
			return err
		}
	}

	return nil
}

func (mh *MeshsyncDataHandler) subscribeToMeshsyncEvents() {
	eventsChan := make(chan *broker.Message)
	err := mh.ListenToMeshSyncEvents(eventsChan)
	if err != nil {
		mh.log.Error(ErrBrokerSubscription(err))
		return
	}
	mh.log.Info("subscribed to meshery broker for meshsync events")

	for event := range eventsChan {
		if event.EventType == broker.ErrorEvent {
			// TODO: Handle errors accordingly
			mh.log.Error(event.Object.(error))
			continue
		}

		// handle the events
		err := mh.meshsyncEventsAccumulator(event)
		if err != nil {
			mh.log.Error(err)
			continue
		}
	}
}

func (mh *MeshsyncDataHandler) ListenToMeshSyncEvents(out chan *broker.Message) error {
	err := mh.broker.SubscribeWithChannel("meshery.meshsync.core", "", out)
	if err != nil {
		return err
	}
	return nil
}

func (mh *MeshsyncDataHandler) subsribeToStoreUpdates(statusChan chan bool) {
	storeChan := make(chan *broker.Message)
	mh.log.Info("subscribing to store updates from meshsync on NATS subject: ", MeshsyncStoreUpdatesSubject)
	err := mh.broker.SubscribeWithChannel(MeshsyncStoreUpdatesSubject, "", storeChan)
	if err != nil {
		mh.log.Error(ErrBrokerSubscription(err))
		statusChan <- false
		return
	}

	statusChan <- true

	for storeUpdate := range storeChan {
		if storeUpdate.EventType == broker.ErrorEvent {
			mh.log.Error(storeUpdate.Object.(error))
			continue
		}

		objectsSlice := storeUpdate.Object.([]interface{})

		for _, object := range objectsSlice {
			obj, err := mh.Unmarshal(object)
			if err != nil {
				continue
			}

			err = mh.persistStoreUpdate(&obj)
			if err != nil {
				mh.log.Error(err)
				continue
			}
		}
	}
}

func (mh *MeshsyncDataHandler) Unmarshal(object interface{}) (meshsyncmodel.KubernetesResource, error) {
	objectJSON, _ := utils.Marshal(object)
	obj := meshsyncmodel.KubernetesResource{}
	err := utils.Unmarshal(objectJSON, &obj)
	if err != nil {
		mh.log.Error(ErrUnmarshal(err, objectJSON))
		return obj, ErrUnmarshal(err, objectJSON)
	}
	return obj, nil
}

// derives the state of the cluster from the events and persists it in the database
func (mh *MeshsyncDataHandler) meshsyncEventsAccumulator(event *broker.Message) error {
	mh.dbHandler.Lock()
	defer mh.dbHandler.Unlock()

	obj, err := mh.Unmarshal(event.Object)
	if err != nil {
		return err
	}

	regQueue := GetMeshSyncRegistrationQueue()
	switch event.EventType {
	case broker.Add:
		compMetadata := mh.getComponentMetadata(obj.APIVersion, obj.Kind)
		obj.ComponentMetadata = utils.MergeMaps(obj.ComponentMetadata, compMetadata)
		result := mh.dbHandler.Create(&obj)
		go regQueue.Send(MeshSyncRegistrationData{MeshsyncDataHandler: *mh, Obj: obj})
		// Try to update object if Create fails. If MeshSync is restarted, on initial sync the discovered data will have eventType as ADD, but the database would already have the data, leading to conflicts hence try to update the object in such cases.
		if result.Error != nil {
			result = mh.dbHandler.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&obj)
			if result.Error != nil {
				return ErrDBPut(result.Error)
			}
		}
	case broker.Update:
		compMetadata := mh.getComponentMetadata(obj.APIVersion, obj.Kind)
		obj.ComponentMetadata = utils.MergeMaps(obj.ComponentMetadata, compMetadata)

		result := mh.dbHandler.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&obj)
		if result.Error != nil {
			return ErrDBPut(result.Error)
		}
		return nil
	case broker.Delete:
		result := mh.dbHandler.Delete(&obj)
		if result.Error != nil {
			return ErrDBDelete(result.Error, "")
		}
	}

	mh.log.Info("Updated database in response to ", event.EventType, " event of object: ", obj.KubernetesResourceMeta.Name, " in namespace: ", obj.KubernetesResourceMeta.Namespace, " of kind: ", obj.Kind)

	return nil
}

func (mh *MeshsyncDataHandler) persistStoreUpdate(object *meshsyncmodel.KubernetesResource) error {
	mh.dbHandler.Lock()
	defer mh.dbHandler.Unlock()
	compMetadata := mh.getComponentMetadata(object.APIVersion, object.Kind)
	object.ComponentMetadata = utils.MergeMaps(object.ComponentMetadata, compMetadata)
	result := mh.dbHandler.Create(object)
	regQueue := GetMeshSyncRegistrationQueue()

	go regQueue.Send(MeshSyncRegistrationData{MeshsyncDataHandler: *mh, Obj: *object})

	if result.Error != nil {
		result = mh.dbHandler.Session(&gorm.Session{FullSaveAssociations: true}).Updates(object)
		if result.Error != nil {
			return ErrDBPut(result.Error)
		}
		mh.log.Info("Updated object: ", object.KubernetesResourceMeta.Name, "/", object.KubernetesResourceMeta.Namespace, " of kind: ", object.Kind, " in the database")
		return nil
	}
	mh.log.Info("Added object: ", object.KubernetesResourceMeta.Name, "/", object.KubernetesResourceMeta.Namespace, " of kind: ", object.Kind, " to the database")

	return nil
}

func RemoveStaleObjects(dbHandler database.Handler) error {
	dbHandler.Lock()
	defer dbHandler.Unlock()

	// Clear stale meshsync data
	err := dbHandler.Migrator().DropTable(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	)
	if err != nil {
		return err
	}
	err = dbHandler.Migrator().CreateTable(
		&meshsyncmodel.KubernetesKeyValue{},
		&meshsyncmodel.KubernetesResource{},
		&meshsyncmodel.KubernetesResourceSpec{},
		&meshsyncmodel.KubernetesResourceStatus{},
		&meshsyncmodel.KubernetesResourceObjectMeta{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (mh *MeshsyncDataHandler) requestMeshsyncStore() error {
	err := mh.broker.Publish(MeshsyncRequestSubject, &broker.Message{
		Request: &broker.RequestObject{
			Entity: "informer-store",
			// TODO: Name of the Reply subject should be taken from some sort of configuration
			Payload: struct{ Reply string }{Reply: MeshsyncStoreUpdatesSubject},
		}})
	if err != nil {
		return ErrRequestMeshsyncStore(err)
	}
	return nil
}

// Returns metadata for the component identified by apiVersion and kind.
// If the component does not exist in the registry, default metadata for k8s component is returned.
func (mh *MeshsyncDataHandler) getComponentMetadata(apiVersion string, kind string) map[string]interface{} {
	var metadata map[string]interface{}

	result := mh.dbHandler.Model(v1alpha1.ComponentDefinitionDB{}).Select("metadata").
		Where("api_version = ? and kind = ?", apiVersion, kind).Scan(&metadata)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			mh.log.Error(ErrResultNotFound(result.Error))
			metadata = K8sMeshModelMetadata
		} else {
			mh.log.Error(ErrDBRead(result.Error))
			metadata = K8sMeshModelMetadata
		}
	}

	return metadata
}
