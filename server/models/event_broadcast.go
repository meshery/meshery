package models

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// listenerBufferSize is the per-subscriber channel buffer. It absorbs short
// bursts of events published faster than a subscriber drains them; once it is
// full, Publish drops the newest event for that subscriber rather than blocking
// the publisher (see Publish).
const listenerBufferSize = 256

// replayBufferSize bounds how many recent ID-bearing events are retained per
// user for SSE Last-Event-ID replay across reconnects.
const replayBufferSize = 256

type replayEntry struct {
	id   string
	data interface{}
}

type replayBuffer struct {
	mu      sync.Mutex
	entries [replayBufferSize]replayEntry
	head    int
	count   int
}

func (rb *replayBuffer) append(id string, data interface{}) {
	rb.entries[rb.head] = replayEntry{id: id, data: data}
	rb.head = (rb.head + 1) % replayBufferSize
	if rb.count < replayBufferSize {
		rb.count++
	}
}

// after returns a chronological snapshot of all entries recorded strictly
// after lastEventID. If lastEventID is empty or not present in the buffer, it
// returns nil. Caller must hold rb.mu.
func (rb *replayBuffer) after(lastEventID string) []interface{} {
	if lastEventID == "" || rb.count == 0 {
		return nil
	}
	start := (rb.head - rb.count + replayBufferSize) % replayBufferSize
	for i := rb.count - 1; i >= 0; i-- {
		idx := (start + i) % replayBufferSize
		if rb.entries[idx].id == lastEventID {
			remaining := rb.count - 1 - i
			if remaining == 0 {
				return nil
			}
			out := make([]interface{}, 0, remaining)
			for j := i + 1; j < rb.count; j++ {
				out = append(out, rb.entries[(start+j)%replayBufferSize].data)
			}
			return out
		}
	}
	return nil
}

type clients struct {
	listeners []chan interface{}
	mu        *sync.Mutex
	evicted   bool
}

type Broadcast struct {
	clients sync.Map
	history sync.Map // core.Uuid -> *replayBuffer
	Name    string
}

// Subscribe registers a new listener for id and returns its channel plus an
// idempotent unsubscribe func. The clients entry is stored behind a pointer so
// Subscribe, Publish and unsubscribe all share one mutex and one listeners
// slice - no stale copy can outlive a concurrent mutation. The channel is
// buffered (listenerBufferSize); the caller must drain it or unsubscribe,
// otherwise a full buffer makes Publish drop the newest events for this listener.
// When all listeners for id unsubscribe, the entry is pruned from the clients
// map to prevent memory leaks across disconnected sessions.
func (c *Broadcast) Subscribe(id core.Uuid) (chan interface{}, func()) {
	ch := make(chan interface{}, listenerBufferSize)

	for {
		actual, ok := c.clients.Load(id)
		if !ok {
			actual, _ = c.clients.LoadOrStore(id, &clients{mu: &sync.Mutex{}})
		}
		cl := actual.(*clients)

		cl.mu.Lock()
		if cl.evicted {
			cl.mu.Unlock()
			continue
		}
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
				if len(cl.listeners) == 0 {
					cl.evicted = true
					c.clients.CompareAndDelete(id, actual)
				}
			})
		}
		return ch, unsubscribe
	}
}

// SubscribeFrom atomically snapshots any buffered events strictly after
// lastEventID and registers a live listener channel for id. Holding the
// per-user replay lock across both steps guarantees no ID-bearing event is
// missed or duplicated between the replay slice and the live channel.
func (c *Broadcast) SubscribeFrom(id core.Uuid, lastEventID string) ([]interface{}, chan interface{}, func()) {
	if lastEventID == "" {
		ch, unsubscribe := c.Subscribe(id)
		return nil, ch, unsubscribe
	}
	rb := c.getOrCreateReplayBuffer(id)
	rb.mu.Lock()
	replay := rb.after(lastEventID)
	ch, unsubscribe := c.Subscribe(id)
	rb.mu.Unlock()
	return replay, ch, unsubscribe
}

func (c *Broadcast) getOrCreateReplayBuffer(id core.Uuid) *replayBuffer {
	if actual, ok := c.history.Load(id); ok {
		return actual.(*replayBuffer)
	}
	actual, _ := c.history.LoadOrStore(id, &replayBuffer{})
	return actual.(*replayBuffer)
}

func extractEventID(data interface{}) string {
	switch e := data.(type) {
	case *events.Event:
		if e != nil && e.ID != uuid.Nil {
			return e.ID.String()
		}
	case events.Event:
		if e.ID != uuid.Nil {
			return e.ID.String()
		}
	}
	return ""
}

// Publish records ID-bearing events in the per-user replay buffer and delivers
// data to every current listener for id. Sends are non-blocking and run under
// the same mutex that guards unsubscribe, which makes this safe on two fronts:
// a channel present in listeners under the lock is guaranteed still open (so
// the send can never panic on a closed channel), and a slow consumer whose
// buffer is full simply has the newest event dropped rather than blocking the
// publisher.
func (c *Broadcast) Publish(id core.Uuid, data interface{}) {
	if eventID := extractEventID(data); eventID != "" {
		rb := c.getOrCreateReplayBuffer(id)
		rb.mu.Lock()
		rb.append(eventID, data)
		c.publishToClients(id, data)
		rb.mu.Unlock()
		return
	}
	c.publishToClients(id, data)
}

func (c *Broadcast) publishToClients(id core.Uuid, data interface{}) {
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
	if cl.evicted {
		return
	}
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
