package models

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
)

type clients struct {
	listeners []chan interface{}
	mu        *sync.Mutex
}

type Broadcast struct {
	clients sync.Map
	Name    string
}

func (c *Broadcast) Subscribe(id uuid.UUID) (chan interface{}, func()) {
	clientMap, _ := c.clients.LoadOrStore(id, clients{mu: &sync.Mutex{}})
	ch := make(chan interface{}, 1)
	connectedClient, _ := clientMap.(clients)

	connectedClient.mu.Lock()
	connectedClient.listeners = append(connectedClient.listeners, ch)
	connectedClient.mu.Unlock()

	c.clients.Store(id, connectedClient)
	unsubscribe := func() {
		cclient, ok := c.clients.Load(id)
		var client clients
		if ok {
			client, ok = cclient.(clients)
			if !ok {
				return
			}
		}
		client.mu.Lock()
		defer client.mu.Unlock()

		listeners := client.listeners
		updatedListeners := []chan interface{}{}
		for _, client := range listeners {
			if client == ch {
				close(client)
				// break
			} else {
				updatedListeners = append(updatedListeners, client)
			}
		}
		client.listeners = updatedListeners
		c.clients.Store(id, client)
	}
	return ch, unsubscribe
}

func (c *Broadcast) Publish(id uuid.UUID, data interface{}) {
	clientMap, ok := c.clients.Load(id)
	if !ok {
		return
	}

	clientToPublish, ok := clientMap.(clients)
	if !ok {
		return
	}
	for _, client := range clientToPublish.listeners {
		if !utils.IsClosed(client) {
			client <- data
		}
	}
}

func NewBroadcaster(name string) *Broadcast {
	return &Broadcast{
		clients: sync.Map{},
		Name:    name,
	}
}
