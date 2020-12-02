package nats

import (
	"github.com/sirupsen/logrus"

	nats "github.com/nats-io/nats.go"
)

// Message (temporary) - will import this from meshsync
type Message struct {
	Type   string
	Object interface{}
}

// Nats - this will implement the nats subsriber for the meshery server
type Nats struct {
	ec *nats.EncodedConn
}

// New - returns a new nats client
func New(url string) (*Nats, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		logrus.Printf("Nats connect Error: %v", err)
		return nil, err
	}
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logrus.Printf("Nats couldnt create encoded connection Error: %v", err)
		return nil, err
	}

	return &Nats{
		ec: ec,
	}, nil
}

// Subscribe - this method will be used to subsribe to subjects.
// subject - string, the subject to subscribe to
func (n *Nats) Subscribe(subject string) error {
	_, err := n.ec.Subscribe(subject, func(msg Message) {
		// read the type and do unmarshalling of object
		logrus.Println("Received a message from Meshsync")
		// logrus.Println(msg.Object)
	})
	if err != nil {
		logrus.Printf("Nats: couldnt subscribe %v", err)
		return err
	}
	return nil
}
