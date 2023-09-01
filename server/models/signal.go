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

type Signal struct {
	clients *sync.Map
}

func (c *Signal) Subscribe(id uuid.UUID) (chan interface{}, func()) {
	clientMap, loaded := c.clients.LoadOrStore(id, &clients{mu: new(sync.Mutex)})
	ch := make(chan interface{})
	connectedClient := clientMap.(*clients)
	
	connectedClient.mu.Lock()
	connectedClient.listeners = append(connectedClient.listeners, ch)
	connectedClient.mu.Unlock()
	fmt.Println("line 22: ", clientMap, loaded, id, connectedClient)

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
	fmt.Println("exiting Subscribe: ")
	return ch, unsubscribe
}

func (c *Signal) Publish(id uuid.UUID, event *events.Event) {
	clientMap, ok := c.clients.Load(id)
	fmt.Println("entered Publish", ok, clientMap, id)
	if !ok {
		return
	}

	clientToPublish, _ := clientMap.(*clients)
	fmt.Println("before signalling")
	for _, client := range clientToPublish.listeners {
		client <- event
	}
	fmt.Println("exiting Publish")
}

func NewSignal() *Signal {
	return &Signal{
		clients: new(sync.Map),
	}
}
