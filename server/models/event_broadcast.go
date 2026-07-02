package models

import (
	"sync"

	"github.com/meshery/schemas/models/core"
)

// clients holds the set of subscriber channels registered for a single id
// together with the mutex that guards that set. It is always stored in
// Broadcast.clients as a *pointer* so the mutex and the listeners slice are
// shared (never copied) across Subscribe/Publish/unsubscribe for the same id.
type clients struct {
	mu        sync.Mutex
	listeners []chan interface{}
}

type Broadcast struct {
	clients sync.Map
	Name    string
}

// Subscribe registers a new listener channel for id and returns it along with
// an unsubscribe function that removes and closes the channel. unsubscribe is
// safe to call multiple times: once the channel has been removed subsequent
// calls are no-ops (so the channel is never double-closed).
func (c *Broadcast) Subscribe(id core.Uuid) (chan interface{}, func()) {
	ch := make(chan interface{}, 1)

	// Attach ch to whichever entry is currently active in the map. If a
	// concurrent unsubscribe removed the entry we loaded (see unsubscribe
	// below) we must retry, otherwise we would register on an orphaned
	// *clients that Publish can no longer reach and the listener would never
	// receive events.
	var client *clients
	for {
		actual, _ := c.clients.LoadOrStore(id, &clients{})
		client = actual.(*clients)

		client.mu.Lock()
		if cur, ok := c.clients.Load(id); !ok || cur != actual {
			// The entry was deleted or replaced after our load; drop it and
			// try again with the fresh entry.
			client.mu.Unlock()
			continue
		}
		client.listeners = append(client.listeners, ch)
		client.mu.Unlock()
		break
	}

	unsubscribe := func() {
		client.mu.Lock()
		defer client.mu.Unlock()

		updated := make([]chan interface{}, 0, len(client.listeners))
		for _, listener := range client.listeners {
			if listener == ch {
				close(listener)
			} else {
				updated = append(updated, listener)
			}
		}
		client.listeners = updated

		// Drop the map entry once its last listener leaves so ids that are no
		// longer subscribed don't accumulate empty entries forever. The
		// CompareAndDelete runs under the same lock that guards the listener
		// set and only removes this exact *clients, so it can never delete an
		// entry a concurrent Subscribe has already attached to: that Subscribe
		// either observes a non-empty slice here (no delete) or fails its
		// post-lock validation above and retries with a fresh entry.
		if len(updated) == 0 {
			c.clients.CompareAndDelete(id, client)
		}
	}
	return ch, unsubscribe
}

// Publish delivers data to every listener currently subscribed to id.
//
// The per-id lock is held for the whole fan-out. Because unsubscribe closes a
// channel only while holding the same lock, no channel can be closed
// mid-send, so Publish can never panic with "send on closed channel" even
// during a concurrent unsubscribe. Holding the lock also makes the iteration
// over listeners race-free.
//
// The send is non-blocking: if a listener's buffer is full the event is
// skipped for that listener rather than blocking the publishing goroutine (and
// every other subscriber) behind a single slow consumer, which would otherwise
// stall the broadcaster and leak goroutines under load.
func (c *Broadcast) Publish(id core.Uuid, data interface{}) {
	actual, ok := c.clients.Load(id)
	if !ok {
		return
	}
	client, ok := actual.(*clients)
	if !ok {
		return
	}

	client.mu.Lock()
	defer client.mu.Unlock()
	for _, listener := range client.listeners {
		select {
		case listener <- data:
		default:
		}
	}
}

func NewBroadcaster(name string) *Broadcast {
	return &Broadcast{
		clients: sync.Map{},
		Name:    name,
	}
}
