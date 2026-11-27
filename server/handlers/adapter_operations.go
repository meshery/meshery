package handlers

import (
	"sync"
	"time"

	"github.com/meshery/meshery/server/meshes"
)

// adapterOperationOwnerTTL is how long the owner of an adapter operation is
// remembered. Adapters report an operation's events while it runs, which takes
// minutes, so a day leaves ample room while keeping the map bounded.
const adapterOperationOwnerTTL = 24 * time.Hour

// adapterOperationOwners remembers which user started each adapter operation.
// Adapters send every event to every subscriber of their event stream, so this
// is what keeps the events of one user's operation from reaching other users.
type adapterOperationOwners struct {
	mu     sync.Mutex
	owners map[string]adapterOperationOwner
	now    func() time.Time
}

type adapterOperationOwner struct {
	userID  string
	started time.Time
}

func newAdapterOperationOwners() *adapterOperationOwners {
	return &adapterOperationOwners{
		owners: map[string]adapterOperationOwner{},
		now:    time.Now,
	}
}

// record notes that userID started the operation, and drops expired entries.
func (o *adapterOperationOwners) record(operationID, userID string) {
	if o == nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	now := o.now()
	for id, owner := range o.owners {
		if now.Sub(owner.started) > adapterOperationOwnerTTL {
			delete(o.owners, id)
		}
	}
	o.owners[operationID] = adapterOperationOwner{userID: userID, started: now}
}

// forget drops an operation that never started.
func (o *adapterOperationOwners) forget(operationID string) {
	if o == nil {
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.owners, operationID)
}

// ownedBy reports whether userID started the operation. An unknown or expired
// operation belongs to nobody.
func (o *adapterOperationOwners) ownedBy(operationID, userID string) bool {
	if o == nil {
		return false
	}
	o.mu.Lock()
	defer o.mu.Unlock()

	owner, ok := o.owners[operationID]
	if !ok || o.now().Sub(owner.started) > adapterOperationOwnerTTL {
		return false
	}
	return owner.userID == userID
}

// forwardAdapterEvent reports whether an adapter event may go to userID. An
// event of an operation goes only to the user who started it. An event without
// an operation ID isn't tied to any user's operation and is forwarded as before.
func forwardAdapterEvent(owners *adapterOperationOwners, event *meshes.EventsResponse, userID string) bool {
	operationID := event.GetOperationId()
	if operationID == "" {
		return true
	}
	return owners.ownedBy(operationID, userID)
}
