package models

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
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
	clientMap, err := c.clients.LoadOrStore(id, clients{mu: &sync.Mutex{}})
	if err {
		logrus.Infof("Client for id %s does not exist in clients map", id)
	}
	ch := make(chan interface{}, 1)
	connectedClient, err := clientMap.(clients)

	if err {
		logrus.Infof("Client for id %s is not connected. Attempting to connect to client %s.", id, id)
	}

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

	// Lock to get a consistent view of listeners and synchronize with unsubscribe
	clientToPublish.mu.Lock()
	listeners := make([]chan interface{}, len(clientToPublish.listeners))
	copy(listeners, clientToPublish.listeners)
	clientToPublish.mu.Unlock()

	for _, ch := range listeners {
		c.safeSend(ch, data)
	}
}

// safeSend attempts to send data to the channel, recovering gracefully if the channel
// is closed. This handles the race condition where a channel may be closed by unsubscribe
// between the time we copy the listeners slice and when we send.
func (c *Broadcast) safeSend(ch chan interface{}, data interface{}) {
	defer func() {
		if r := recover(); r != nil {
			// Channel was closed between copy and send, ignore
		}
	}()
	ch <- data
}

func NewBroadcaster(name string) *Broadcast {
	return &Broadcast{
		clients: sync.Map{},
		Name:    name,
	}
}
