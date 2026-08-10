package models

import (
	"sync"

	"github.com/meshery/schemas/models/core"
)

// listenerBufferSize is the per-subscriber channel buffer. It absorbs short
// bursts of events published faster than a subscriber drains them; once it is
// full, Publish drops the newest event for that subscriber rather than blocking
// the publisher (see Publish).
const listenerBufferSize = 256

type clients struct {
	listeners []chan interface{}
	mu        *sync.Mutex
}

type Broadcast struct {
	clients sync.Map
	Name    string
}

// Subscribe registers a new listener for id and returns its channel plus an
// idempotent unsubscribe func. The clients entry is stored behind a pointer so
// Subscribe, Publish and unsubscribe all share one mutex and one listeners
// slice - no stale copy can outlive a concurrent mutation. The channel is
// buffered (listenerBufferSize); the caller must drain it or unsubscribe,
// otherwise a full buffer makes Publish drop the newest events for this listener.
func (c *Broadcast) Subscribe(id core.Uuid) (chan interface{}, func()) {
	actual, _ := c.clients.LoadOrStore(id, &clients{mu: &sync.Mutex{}})
	cl := actual.(*clients)

	ch := make(chan interface{}, listenerBufferSize)

	cl.mu.Lock()
	cl.listeners = append(cl.listeners, ch)
	cl.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			cl.mu.Lock()
			defer cl.mu.Unlock()
			for i, listener := range cl.listeners {
				if listener == ch {
					cl.listeners = append(cl.listeners[:i], cl.listeners[i+1:]...)
					close(ch)
					break
				}
			}
		})
	}
	return ch, unsubscribe
}

// Publish delivers data to every current listener for id. Sends are
// non-blocking and run under the same mutex that guards unsubscribe, which makes
// this safe on two fronts: a channel present in listeners under the lock is
// guaranteed still open (so the send can never panic on a closed channel), and a
// slow consumer whose buffer is full simply has the newest event dropped rather
// than blocking the publisher.
func (c *Broadcast) Publish(id core.Uuid, data interface{}) {
	actual, ok := c.clients.Load(id)
	if !ok {
		return
	}
	cl, ok := actual.(*clients)
	if !ok {
		return
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()
	for _, listener := range cl.listeners {
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
