package models

import (
	"sync"
	"sync/atomic"

	"github.com/meshery/meshkit/logger"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

var (
	once              sync.Once
	registrationQueue *MeshSyncRegistrationQueue
)

// DefaultQueueCapacity is the default buffer size for the registration channel.
const DefaultQueueCapacity = 100

// MeshSyncRegistrationQueue holds a channel for queuing Kubernetes resources that are
// eligible to be registered as Meshery connections.
//
// TODO: MeshSync additionally publishes these resources to a broker topic/subject
// (meshsync.registerconnection.queue?); Meshery Server subscribes and performs the
// necessary action.
type MeshSyncRegistrationQueue struct {
	RegChan chan MeshSyncRegistrationData
	log     logger.Handler
	dropped int64
}

type MeshSyncRegistrationData struct {
	MeshsyncDataHandler MeshsyncDataHandler
	Obj                 meshsyncmodel.KubernetesResource
}

// InitMeshSyncRegistrationQueue initializes the registration queue singleton.
func InitMeshSyncRegistrationQueue() {
	initQueue()
}

// GetMeshSyncRegistrationQueue returns the registration queue singleton,
// initializing it if necessary.
func GetMeshSyncRegistrationQueue() *MeshSyncRegistrationQueue {
	initQueue()
	return registrationQueue
}

// NewMeshSyncRegistrationQueue creates a new MeshSyncRegistrationQueue with the specified channel capacity.
func NewMeshSyncRegistrationQueue(capacity int) *MeshSyncRegistrationQueue {
	if capacity <= 0 {
		capacity = DefaultQueueCapacity
	}
	return &MeshSyncRegistrationQueue{
		RegChan: make(chan MeshSyncRegistrationData, capacity),
	}
}

// SetLogger assigns an optional logger to the registration queue for observability.
func (mrq *MeshSyncRegistrationQueue) SetLogger(log logger.Handler) {
	if mrq != nil {
		mrq.log = log
	}
}

// Dropped returns the total number of registration items dropped due to a saturated queue buffer.
func (mrq *MeshSyncRegistrationQueue) Dropped() int64 {
	if mrq == nil {
		return 0
	}
	return atomic.LoadInt64(&mrq.dropped)
}

// initQueue initializes the registration queue exactly once using sync.Once.
func initQueue() {
	once.Do(func() {
		registrationQueue = NewMeshSyncRegistrationQueue(DefaultQueueCapacity)
	})
}

// Send enqueues data into RegChan in a non-blocking manner.
// If mrq is nil, RegChan is nil, or the channel buffer is full, Send drops
// the item immediately and returns false to prevent unbounded goroutine leaks.
// When dropping an item due to a full buffer, a warning is logged and the dropped counter is incremented.
func (mrq *MeshSyncRegistrationQueue) Send(data MeshSyncRegistrationData) bool {
	if mrq == nil || mrq.RegChan == nil {
		return false
	}
	select {
	case mrq.RegChan <- data:
		return true
	default:
		atomic.AddInt64(&mrq.dropped, 1)
		var resName, resNamespace string
		if data.Obj.KubernetesResourceMeta != nil {
			resName = data.Obj.KubernetesResourceMeta.Name
			resNamespace = data.Obj.KubernetesResourceMeta.Namespace
		}
		resKind := data.Obj.Kind
		if resKind == "" {
			resKind = "resource"
		}
		if data.MeshsyncDataHandler.log != nil {
			data.MeshsyncDataHandler.log.Warnf("MeshSync registration queue is full (capacity %d); dropping registration event for %s %s/%s", cap(mrq.RegChan), resKind, resNamespace, resName)
		} else if mrq.log != nil {
			mrq.log.Warnf("MeshSync registration queue is full (capacity %d); dropping registration event for %s %s/%s", cap(mrq.RegChan), resKind, resNamespace, resName)
		}
		return false
	}
}
