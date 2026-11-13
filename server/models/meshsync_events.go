package models

import (
	"fmt"
	"strings"
	"sync"

	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1beta3/component"
	"gorm.io/gorm"
)

const (
	MeshsyncStoreUpdatesSubject = "meshery-server.meshsync.store"
	MeshsyncRequestSubject      = "meshery.meshsync.request"
	// meshsyncEventsSubject carries the MeshSync resource events the handler
	// persists. Both MeshSync subjects are subscribed with an empty queue
	// name, so pausing one handler while another connects must unsubscribe
	// them, or both handlers receive - and persist - the same message.
	meshsyncEventsSubject = "meshery.meshsync.core"
)

type MeshsyncDataHandler struct {
	broker       broker.Handler
	dbHandler    database.Handler
	log          logger.Handler
	Provider     Provider
	UserID       core.Uuid
	ConnectionID core.Uuid
	InstanceID   core.Uuid
	Token        string
	StopFunc     func()
	stopCh       chan struct{}
	stopOnce     *sync.Once
	listenerWg   *sync.WaitGroup
	// eventsChan and storeChan carry the broker's MeshSync messages to the
	// listener goroutines. Kept on the handler so PauseSubscriptions and
	// ResumeSubscriptions can unsubscribe and re-subscribe the same channels.
	// Each is written once by its listener goroutine, and Run waits for both
	// listeners to subscribe before it returns, so a pause or resume - which
	// only happens after Run - always sees the stored channel.
	eventsChan chan *broker.Message
	storeChan  chan *broker.Message
	// paused reports whether the MeshSync subscriptions are paused: the
	// listener goroutines are still running, but the broker delivers nothing
	// to them until ResumeSubscriptions re-subscribes. PauseSubscriptions and
	// ResumeSubscriptions are called only from the broker supervisor's
	// rebuild path, one goroutine, so no lock guards it - and none can be
	// added: MeshsyncDataHandler is copied by value into
	// MeshSyncRegistrationData, and a lock value in it trips the vet
	// copylocks check on those copies.
	paused bool
}

func NewMeshsyncDataHandler(broker broker.Handler, dbHandler database.Handler, log logger.Handler, provider Provider, userID, connID, instanceID core.Uuid, token string, stopFunc func()) *MeshsyncDataHandler {
	return &MeshsyncDataHandler{
		broker:       broker,
		dbHandler:    dbHandler,
		log:          log,
		Provider:     provider,
		UserID:       userID,
		ConnectionID: connID,
		InstanceID:   instanceID,
		Token:        token,
		StopFunc:     stopFunc,
		stopCh:       make(chan struct{}),
		stopOnce:     &sync.Once{},
		listenerWg:   &sync.WaitGroup{},
	}
}

func (mh *MeshsyncDataHandler) GetBrokerHandler() broker.Handler {
	return mh.broker
}

// IsConnected reports whether Meshery currently holds a live broker connection
// for this connection — i.e. it is actually able to receive MeshSync data. This
// is the authoritative "connected" signal, as opposed to re-probing the broker's
// monitoring endpoint (which false-negatives when that endpoint isn't reachable
// from Meshery even though the data path is up).
func (mh *MeshsyncDataHandler) IsConnected() bool {
	return mh != nil && mh.broker != nil && mh.broker.IsConnected()
}

func (mh *MeshsyncDataHandler) Run() error {
	storeSubscriptionStatusChan := make(chan bool)
	eventsSubscribed := make(chan struct{})
	if mh.listenerWg == nil {
		mh.listenerWg = &sync.WaitGroup{}
	}
	// this subscription is independent of whether or not the stale data in the database have been cleaned up
	mh.listenerWg.Add(1)
	go mh.subscribeToMeshsyncEvents(eventsSubscribed)

	mh.listenerWg.Add(1)
	go mh.subsribeToStoreUpdates(storeSubscriptionStatusChan)
	// Do not return before the events listener has subscribed: the broker
	// supervisor may pause this handler as soon as it is attached, and a
	// subscription that lands after the pause would keep the paused handler
	// receiving, and be doubled by the resume.
	<-eventsSubscribed
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

// subscribeToMeshsyncEvents closes subscribed once its subscription attempt
// is over, whether or not it succeeded.
func (mh *MeshsyncDataHandler) subscribeToMeshsyncEvents(subscribed chan<- struct{}) {
	defer mh.listenerWg.Done()

	eventsChan := make(chan *broker.Message, 1)
	err := mh.ListenToMeshSyncEvents(eventsChan)
	if err != nil {
		close(subscribed)
		mh.log.Error(ErrBrokerSubscription(err))
		return
	}
	// The channel lives on the handler so ResumeSubscriptions can
	// re-subscribe it after a pause. It is stored only once subscribed: a
	// failed subscription has no listener, and re-subscribing its channel
	// would block the broker on a channel nobody reads.
	mh.eventsChan = eventsChan
	close(subscribed)
	mh.log.Info("subscribed to meshery broker for meshsync events")

	for {
		select {
		case <-mh.stopSignal():
			return
		case event, ok := <-eventsChan:
			if !ok {
				return
			}
			if event.EventType == broker.ErrorEvent {
				if err, ok := event.Object.(error); ok {
					mh.log.Error(ErrMeshsyncEvent(err))
				} else {
					mh.log.Error(ErrMeshsyncEvent(fmt.Errorf("received meshsync error event with non-error object: %T", event.Object)))
				}
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
}

func (mh *MeshsyncDataHandler) ListenToMeshSyncEvents(out chan *broker.Message) error {
	err := mh.broker.SubscribeWithChannel(meshsyncEventsSubject, "", out)
	if err != nil {
		return err
	}
	return nil
}

func (mh *MeshsyncDataHandler) subsribeToStoreUpdates(statusChan chan bool) {
	defer mh.listenerWg.Done()

	// The channel lives on the handler so ResumeSubscriptions can
	// re-subscribe it after a pause.
	storeChan := make(chan *broker.Message, 1)
	mh.storeChan = storeChan
	mh.log.Info("subscribing to store updates from meshsync on NATS subject: ", MeshsyncStoreUpdatesSubject)
	err := mh.broker.SubscribeWithChannel(MeshsyncStoreUpdatesSubject, "", storeChan)
	if err != nil {
		mh.log.Error(ErrBrokerSubscription(err))
		statusChan <- false
		return
	}

	statusChan <- true

	for {
		select {
		case <-mh.stopSignal():
			return
		case storeUpdate, ok := <-storeChan:
			if !ok {
				return
			}
			if storeUpdate.EventType == broker.ErrorEvent {
				if err, ok := storeUpdate.Object.(error); ok {
					mh.log.Error(ErrMeshsyncStoreUpdates(err))
				} else {
					mh.log.Error(ErrMeshsyncStoreUpdates(fmt.Errorf("received store update error event with non-error object: %T", storeUpdate.Object)))
				}
				continue
			}

			objectsSlice, ok := storeUpdate.Object.([]interface{})
			if !ok {
				mh.log.Error(ErrMeshsyncStoreUpdates(fmt.Errorf("received store update event with unexpected object type: %T", storeUpdate.Object)))
				continue
			}

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
}

func (mh *MeshsyncDataHandler) stopSignal() <-chan struct{} {
	if mh == nil || mh.stopCh == nil {
		closedCh := make(chan struct{})
		close(closedCh)
		return closedCh
	}

	return mh.stopCh
}

func (mh *MeshsyncDataHandler) Unmarshal(object interface{}) (meshsyncmodel.KubernetesResource, error) {
	objectJSON, _ := utils.Marshal(object)
	obj := meshsyncmodel.KubernetesResource{}
	err := encoding.Unmarshal([]byte(objectJSON), &obj)
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
		compMetadata, model := mh.getComponentMetadata(obj.APIVersion, obj.Kind)
		obj.ComponentMetadata = utils.MergeMaps(obj.ComponentMetadata, compMetadata)
		obj.Model = model
		result := mh.dbHandler.Create(&obj)
		go regQueue.Send(MeshSyncRegistrationData{MeshsyncDataHandler: *mh, Obj: obj})
		// Try to update object if Create fails. If MeshSync is restarted, on initial sync the discovered data will have eventType as ADD, but the database would already have the data, leading to conflicts hence try to update the object in such cases.
		if result.Error != nil {
			if err := mh.replaceMeshsyncResource(&obj); err != nil {
				return ErrDBPut(err)
			}
		}
	case broker.Update:
		compMetadata, model := mh.getComponentMetadata(obj.APIVersion, obj.Kind)
		obj.ComponentMetadata = utils.MergeMaps(obj.ComponentMetadata, compMetadata)
		obj.Model = model
		if err := mh.replaceMeshsyncResource(&obj); err != nil {
			return ErrDBPut(err)
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
	compMetadata, model := mh.getComponentMetadata(object.APIVersion, object.Kind)
	object.ComponentMetadata = utils.MergeMaps(object.ComponentMetadata, compMetadata)
	object.Model = model
	result := mh.dbHandler.Create(object)
	regQueue := GetMeshSyncRegistrationQueue()

	go regQueue.Send(MeshSyncRegistrationData{MeshsyncDataHandler: *mh, Obj: *object})

	if result.Error != nil {
		if err := mh.replaceMeshsyncResource(object); err != nil {
			return ErrDBPut(err)
		}
		mh.log.Info("Updated object: ", object.KubernetesResourceMeta.Name, "/", object.KubernetesResourceMeta.Namespace, " of kind: ", object.Kind, " in the database")
		return nil
	}
	mh.log.Info("Added object: ", object.KubernetesResourceMeta.Name, "/", object.KubernetesResourceMeta.Namespace, " of kind: ", object.Kind, " to the database")

	return nil
}

// replaceMeshsyncResource replaces the metadata rows and resource atomically. GORM's
// association updates insert new key/value rows but do not remove old values.
func (mh *MeshsyncDataHandler) replaceMeshsyncResource(obj *meshsyncmodel.KubernetesResource) error {
	meshsyncmodel.SetID(obj)
	if obj.ID == "" || obj.KubernetesResourceMeta == nil {
		return fmt.Errorf("cannot update MeshSync resource without an identity and metadata")
	}
	return mh.dbHandler.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", obj.ID).Delete(&meshsyncmodel.KubernetesKeyValue{}).Error; err != nil {
			return err
		}
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(obj).Error
	})
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
	componentDef := component.ComponentDefinition{} // Retrieve the entire component
	defer func() {
		data, _ = utils.MarshalAndUnmarshal[component.ComponentDefinition, map[string]interface{}](componentDef)
		if componentDef.Model != nil {
			model = componentDef.Model.Name
		}
	}()

	// Query the database for the complete component definition.
	result := mh.dbHandler.Model(component.ComponentDefinition{}).
		Where("component->>'version' = ? AND component->>'kind' = ?", apiVersion, kind).
		First(&componentDef)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			mh.log.Debug("No component definition found for apiVersion: ", apiVersion, " kind: ", kind, ". Using default metadata.")
		} else if isMissingTableError(result.Error, component.ComponentDefinition{}.TableName()) {
			mh.log.Debug("Component registry table unavailable while looking up metadata for apiVersion: ", apiVersion, " kind: ", kind, ". Using default metadata.")
		} else {
			mh.log.Error(ErrDBRead(result.Error))
		}
		// Provide a default or fallback component definition
		componentDef = component.ComponentDefinition{
			Styles: &K8sMeshModelMetadata.Styles,
		}
		return
	}

	if componentDef.ModelID != nil {
		modelDef := modelv1beta1.ModelDefinition{}
		result = mh.dbHandler.Session(&gorm.Session{NewDB: true}).Model(&modelv1beta1.ModelDefinition{}).
			Where("id = ?", componentDef.ModelID).
			First(&modelDef)
		if result.Error == nil {
			componentDef.Model = &modelDef
		} else if result.Error != gorm.ErrRecordNotFound && !isMissingTableError(result.Error, modelv1beta1.ModelDefinition{}.TableName()) {
			mh.log.Error(ErrDBRead(result.Error))
		}
	}

	return
}

func isMissingTableError(err error, tableName string) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "no such table: "+tableName) ||
		strings.Contains(msg, `relation "`+tableName+`" does not exist`) ||
		strings.Contains(msg, "."+tableName+"' doesn't exist")
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
	if mh == nil {
		return
	}
	if mh.stopOnce == nil {
		mh.stopOnce = &sync.Once{}
	}
	if mh.listenerWg == nil {
		mh.listenerWg = &sync.WaitGroup{}
	}

	mh.stopOnce.Do(func() {
		if mh.stopCh != nil {
			close(mh.stopCh)
		}

		if mh.StopFunc != nil {
			mh.StopFunc()
		}

		if mh.broker != nil {
			mh.broker.CloseConnection()
		}

		mh.listenerWg.Wait()
	})
}

// PauseSubscriptions unsubscribes the handler's MeshSync subjects so a
// reconnecting broker connection stops delivering messages to it, without
// stopping the handler itself: the listener goroutines keep running, and the
// subscriptions are restored by ResumeSubscriptions. The broker supervisor
// pauses the old handler while a replacement connects, so the two handlers
// cannot receive - and persist - the same message. Idempotent, and safe to
// call on a handler whose subscriptions never started.
func (mh *MeshsyncDataHandler) PauseSubscriptions() {
	if mh == nil || mh.broker == nil {
		return
	}
	if mh.paused {
		return
	}
	mh.paused = true
	for _, subject := range []string{meshsyncEventsSubject, MeshsyncStoreUpdatesSubject} {
		if err := mh.broker.Unsubscribe(subject); err != nil {
			mh.log.Error(ErrBrokerSubscription(err))
		}
	}
}

// ResumeSubscriptions re-subscribes the MeshSync subjects paused by
// PauseSubscriptions, so the handler receives messages again. A no-op unless
// the handler is paused, and a no-op once the handler is stopped: a stopped
// handler's broker connection is closed, so re-subscribing would only fail.
func (mh *MeshsyncDataHandler) ResumeSubscriptions() {
	if mh == nil || mh.broker == nil {
		return
	}
	if !mh.paused {
		return
	}
	mh.paused = false
	select {
	case <-mh.stopSignal():
		return
	default:
	}
	if mh.eventsChan != nil {
		if err := mh.broker.SubscribeWithChannel(meshsyncEventsSubject, "", mh.eventsChan); err != nil {
			mh.log.Error(ErrBrokerSubscription(err))
		}
	}
	if mh.storeChan != nil {
		if err := mh.broker.SubscribeWithChannel(MeshsyncStoreUpdatesSubject, "", mh.storeChan); err != nil {
			mh.log.Error(ErrBrokerSubscription(err))
		}
	}
}
