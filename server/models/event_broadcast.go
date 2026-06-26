package models

import (
	"sync"

	"github.com/meshery/schemas/models/core"
)

type clients struct {
	listeners []chan interface{}
	mu        sync.Mutex
}

type Broadcast struct {
	clients sync.Map
	Name    string
}

func (c *Broadcast) Subscribe(id core.Uuid) (chan interface{}, func()) {
	actual, _ := c.clients.LoadOrStore(id, &clients{})
	cc := actual.(*clients)

	ch := make(chan interface{}, 1)

	cc.mu.Lock()
	cc.listeners = append(cc.listeners, ch)
	cc.mu.Unlock()

	unsubscribe := func() {
		v, ok := c.clients.Load(id)
		if !ok {
			return
		}
		cc, ok := v.(*clients)
		if !ok {
			return
		}

		cc.mu.Lock()
		found := false
		updated := make([]chan interface{}, 0, len(cc.listeners))
		for _, l := range cc.listeners {
			if l == ch {
				found = true
				continue
			}
			updated = append(updated, l)
		}
		cc.listeners = updated
		cc.mu.Unlock()

		if found {
			close(ch)
		}
	}
	return ch, unsubscribe
}

// Publish delivers data to every subscriber for id. Delivery is
// best-effort: each listener channel has a buffer of 1, and a non-blocking
// send is used so a slow subscriber cannot block other subscribers or the
// publisher. If the listener buffer is full, the event is dropped for that
// subscriber.
func (c *Broadcast) Publish(id core.Uuid, data interface{}) {
	v, ok := c.clients.Load(id)
	if !ok {
		return
	}
	cc, ok := v.(*clients)
	if !ok {
		return
	}

	cc.mu.Lock()
	defer cc.mu.Unlock()
	for _, client := range cc.listeners {
		select {
		case client <- data:
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
