package models

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/models/events"
)

type clients struct {
	listeners []chan interface{}
	mu *sync.Mutex
}

type EventBroadcast struct {
	clients *sync.Map
}

func (c *EventBroadcast) Subscribe(id uuid.UUID) (chan interface{}, func()) {
	clientMap, _ := c.clients.LoadOrStore(id, &clients{mu: new(sync.Mutex)})
	ch := make(chan interface{})
	connectedClient := clientMap.(*clients)
	
	connectedClient.mu.Lock()
	connectedClient.listeners = append(connectedClient.listeners, ch)
	connectedClient.mu.Unlock()

	unsubscribe := func() {
		connectedClient.mu.Lock()
		defer connectedClient.mu.Unlock()
		for index, client := range connectedClient.listeners {
			if client == ch {
				fmt.Println("inside listeneres: ", id)
				close(client)
				connectedClient.listeners = append(connectedClient.listeners[:index], connectedClient.listeners[index + 1:]...)
			}
		}
	}
	return ch, unsubscribe
}

func (c *EventBroadcast) Publish(id uuid.UUID, event *events.Event) {
	clientMap, ok := c.clients.Load(id)
	if !ok {
		return
	}

	clientToPublish, _ := clientMap.(*clients)
	for _, client := range clientToPublish.listeners {
		client <- event
	}
}

func NewEventBroadcaster() *EventBroadcast {
	return &EventBroadcast{
		clients: new(sync.Map),
	}
}
