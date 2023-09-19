package models

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
)

type clients struct {
	listeners []chan interface{}
	mu        *sync.Mutex
}

type Broadcast struct {
	clients *sync.Map
}

func (c *Broadcast) Subscribe(id uuid.UUID) (chan interface{}, func()) {
	clientMap, _ := c.clients.LoadOrStore(id, &clients{mu: new(sync.Mutex)})
	ch := make(chan interface{}, 1)
	fmt.Println("line 22:=========", id)
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
				connectedClient.listeners = append(connectedClient.listeners[:index], connectedClient.listeners[index+1:]...)
			}
		}
	}
	return ch, unsubscribe
}

func (c *Broadcast) Publish(id uuid.UUID, data interface{}) {
	fmt.Println("test------------", id)
	clientMap, ok := c.clients.Load(id)
	if !ok {
		return
	}
	
	clientToPublish, _ := clientMap.(*clients)
	for _, client := range clientToPublish.listeners {
		fmt.Println("test------------mmmmmmmmm", id)
		client <- data
	}
}

func NewBroadcaster() *Broadcast {
	return &Broadcast{
		clients: new(sync.Map),
	}
}
