package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta1/component"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MeshsyncStoreUpdatesSubject = "meshery-server.meshsync.store"
	MeshsyncRequestSubject      = "meshery.meshsync.request"
)

type meshsyncInternalEvent struct {
	EventType broker.EventType
	Object    meshsyncmodel.KubernetesResource
}
type metadataCacheEntry struct {
	data  map[string]any
	model string
}

// TODO: Create proper error codes for the functionalities this struct implements
type MeshsyncDataHandler struct {
	broker        broker.Handler
	dbHandler     database.Handler
	log           logger.Handler
	Provider      Provider
	UserID        uuid.UUID
	ConnectionID  uuid.UUID
	InstanceID    uuid.UUID
	Token         string
	StopFunc      func()
	metadataCache *lru.Cache[string, metadataCacheEntry]
	eventBuffer   chan meshsyncInternalEvent
}

func NewMeshsyncDataHandler(broker broker.Handler, dbHandler database.Handler, log logger.Handler, provider Provider, userID, connID, instanceID uuid.UUID, token string, stopFunc func()) *MeshsyncDataHandler {
	// TODO: make the 200 configurable so the user can define the cache
	cache, _ := lru.New[string, metadataCacheEntry](200)

	return &MeshsyncDataHandler{
		broker:        broker,
		dbHandler:     dbHandler,
		log:           log,
		Provider:      provider,
		UserID:        userID,
		ConnectionID:  connID,
		InstanceID:    instanceID,
		Token:         token,
		StopFunc:      stopFunc,
		metadataCache: cache,
		// TODO: the event buffer needs to be configurable
		eventBuffer: make(chan meshsyncInternalEvent, 2000),
	}
}

func (mh *MeshsyncDataHandler) GetBrokerHandler() broker.Handler {
	return mh.broker
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

	// Start the batch event processor
	go mh.eventProcessor()

	for event := range eventsChan {
		if event.EventType == broker.ErrorEvent {
			// TODO: Handle errors accordingly
			mh.log.Error(event.Object.(error))
			continue
		}

		obj, err := mh.Unmarshal(event.Object)
		if err != nil {
			mh.log.Error(err)
			continue
		}

		// Push the event to the buffer for batch processing
		mh.eventBuffer <- meshsyncInternalEvent{
			EventType: event.EventType,
			Object:    obj,
		}
	}
}

// eventProcessor drains the eventBuffer and processes events in batches
func (mh *MeshsyncDataHandler) eventProcessor() {
	// process batch every 500ms or when buffer is full
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	const maxBatchSize = 100
	var count int

	// hot pots for events
	// updates and adds are the same
	addsUpdates := make([]meshsyncmodel.KubernetesResource, 0, maxBatchSize)
	deletes := make([]meshsyncmodel.KubernetesResource, 0, maxBatchSize)

	processAll := func() {
		if count == 0 {
			return
		}

		// High-performance batch processing with individual fallback
		mh.processAddUpdateBatch(addsUpdates)
		mh.processDeleteBatch(deletes)

		// reuse the buffers
		addsUpdates = addsUpdates[:0]
		deletes = deletes[:0]
		count = 0
	}

	for {
		select {
		case event, ok := <-mh.eventBuffer:
			if !ok {
				processAll()
				return
			}

			// set object directly before buffering
			compMetadata, model := mh.getComponentMetadata(event.Object.APIVersion, event.Object.Kind)
			event.Object.ComponentMetadata = utils.MergeMaps(event.Object.ComponentMetadata, compMetadata)
			event.Object.Model = model

			switch event.EventType {
			case broker.Add, broker.Update:
				addsUpdates = append(addsUpdates, event.Object)
			case broker.Delete:
				deletes = append(deletes, event.Object)
			}

			count++
			if count >= maxBatchSize {
				processAll()
			}
		case <-ticker.C:
			processAll()
		}
	}
}

func (mh *MeshsyncDataHandler) processAddUpdateBatch(objs []meshsyncmodel.KubernetesResource) {
	if len(objs) == 0 {
		return
	}

	// Try Bulk Upsert for gain perfromance
	// the bulk works on batches with 100 objs
	err := mh.dbHandler.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(objs, 100).Error

	// Fallback to one-by-one with SavePoints to ensure "others continue"
	if err != nil {
		mh.log.Warn(err)
		var succeeded []meshsyncmodel.KubernetesResource
		txErr := mh.dbHandler.Transaction(func(tx *gorm.DB) error {
			for i, obj := range objs {
				sp := fmt.Sprintf("sp_add_%d_%d", time.Now().UnixNano(), i)
				tx.SavePoint(sp)
				if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&obj).Error; err != nil {
					mh.log.Error(err)
					tx.RollbackTo(sp)
				} else {
					succeeded = append(succeeded, obj)
				}
			}
			// always commit what succeeded
			return nil
		})
		if txErr != nil {
			mh.log.Error(fmt.Errorf("fallback transaction failed to commit: %w", txErr))
		}
		if len(succeeded) == 0 {
			mh.log.Warnf("processAddBatch: all %d objects failed in fallback", len(objs))
		}
		// Only send to the registration queue if the transaction itself committed.
		if txErr == nil {
			regQueue := GetMeshSyncRegistrationQueue()
			for _, obj := range succeeded {
				regQueue.Send(MeshSyncRegistrationData{MeshsyncDataHandler: *mh, Obj: obj})
			}
		}
		return
	}
}
func (mh *MeshsyncDataHandler) processDeleteBatch(objs []meshsyncmodel.KubernetesResource) {
	if len(objs) == 0 {
		return
	}
	// Try bulk delete
	err := mh.dbHandler.Delete(&objs).Error
	// Fallback to one-by-one with SavePoints to ensure "others continue"
	if err != nil {
		mh.log.Warn(err)
		mh.dbHandler.Transaction(func(tx *gorm.DB) error {
			for _, obj := range objs {
				sp := fmt.Sprintf("sp_del_%d", time.Now().UnixNano())
				tx.SavePoint(sp)
				if err := tx.Delete(&obj).Error; err != nil {
					mh.log.Error(err)
					tx.RollbackTo(sp)
				}
			}
			return nil
		})
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

		objectsSlice, ok := storeUpdate.Object.([]interface{})
		if !ok {
			continue
		}

		for _, object := range objectsSlice {
			obj, err := mh.Unmarshal(object)
			if err != nil {
				continue
			}

			// Pipe store updates through the same high-performance batch pipeline
			mh.eventBuffer <- meshsyncInternalEvent{
				EventType: broker.Add,
				Object:    obj,
			}
		}
	}
}

func (mh *MeshsyncDataHandler) Unmarshal(object interface{}) (meshsyncmodel.KubernetesResource, error) {
	var data []byte
	var err error

	switch v := object.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		jsonStr, err := utils.Marshal(object)
		if err != nil {
			return meshsyncmodel.KubernetesResource{}, ErrUnmarshal(err, fmt.Sprintf("%v", object))
		}
		data = []byte(jsonStr)
	}

	obj := meshsyncmodel.KubernetesResource{}
	err = encoding.Unmarshal(data, &obj)
	if err != nil {
		return obj, ErrUnmarshal(err, string(data))
	}
	return obj, nil
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
func (mh *MeshsyncDataHandler) getComponentMetadata(apiVersion string, kind string) (data map[string]interface{}, model string) {
	cacheKey := apiVersion + "/" + kind
	if entry, ok := mh.metadataCache.Get(cacheKey); ok {
		return entry.data, entry.model
	}

	componentDef := component.ComponentDefinition{}
	// Query the database for the complete component definition
	result := mh.dbHandler.Model(component.ComponentDefinition{}).Preload("Model").
		Where("component->>'version' = ? AND component->>'kind' = ?", apiVersion, kind).
		First(&componentDef)

	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			mh.log.Error(ErrDBRead(result.Error))
		}
		// Provide a default or fallback component definition
		componentDef = component.ComponentDefinition{
			Styles: &K8sMeshModelMetadata.Styles,
		}
	}

	data, _ = utils.MarshalAndUnmarshal[component.ComponentDefinition, map[string]interface{}](componentDef)
	if componentDef.Model != nil {
		model = componentDef.Model.Name
	}

	// Cache the result (including fallbacks) to avoid repeated DB misses
	mh.metadataCache.Add(cacheKey, metadataCacheEntry{data, model})
	return
}

func (mh *MeshsyncDataHandler) Resync() error {
	if mh.broker.Info() == broker.NotConnected {
		mh.log.Warnf("Resync meshsync: broker is not connected")
		return nil
	}
	err := mh.broker.Publish(MeshsyncRequestSubject, &broker.Message{
		Request: &broker.RequestObject{
			Entity: broker.ReSyncDiscoveryEntity,
		},
	})
	if err != nil {
		return ErrMeshsyncDataHandler(err)
	}
	return nil
}

func (mh *MeshsyncDataHandler) Stop() {
	if mh.StopFunc != nil {
		mh.StopFunc()
	}
}
