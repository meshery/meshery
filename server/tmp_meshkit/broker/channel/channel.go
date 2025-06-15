// nolint
// because this is temporally here and will be moved under meshkit
package channel

// TODO
// put this under meshkit
// as this is required for both meshsync and server

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	realBroker "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/utils"
)

type TMPChannelBrokerHandler struct {
	Options
	name string
	// this structure represents [subject] => [queue] => channel
	// so there is a channel per queue per subject
	storage map[string]map[string]chan *realBroker.Message
}

func NewTMPChannelBrokerHandler(optsSetters ...OptionsSetter) *TMPChannelBrokerHandler {
	options := DefautOptions
	for _, setOptions := range optsSetters {
		if setOptions != nil {
			setOptions(&options)
		}
	}
	return &TMPChannelBrokerHandler{
		name: fmt.Sprintf(
			"channel-broker-handler--%s",
			uuid.New().String(),
		),
		Options: options,
		storage: make(map[string]map[string]chan *realBroker.Message),
	}
}

func (h *TMPChannelBrokerHandler) ConnectedEndpoints() (endpoints []string) {
	// return subjects::queue list intead of connection endpoints
	list := make([]string, 0, len(h.storage))
	for subject, qstorage := range h.storage {
		if qstorage == nil {
			continue
		}
		for queue := range qstorage {
			list = append(
				list,
				fmt.Sprintf(
					"%s::%s",
					subject,
					queue,
				),
			)
		}

	}
	return list
}

func (h *TMPChannelBrokerHandler) Info() string {
	// return name because nats implementation returns name
	return h.name
}

func (h *TMPChannelBrokerHandler) CloseConnection() {
	for subject, qstorage := range h.storage {
		for queue, ch := range qstorage {
			if !utils.IsClosed(ch) {
				close(ch)
			}
			delete(qstorage, queue)
		}
		delete(h.storage, subject)
	}
}

// Publish - to publish messages
func (h *TMPChannelBrokerHandler) Publish(subject string, message *realBroker.Message) error {
	if len(h.storage[subject]) <= 0 {
		// nobody is listening => not publishing
		return nil
	}

	var successList []string
	var failedList []string

	for queue, ch := range h.storage[subject] {
		select {
		case ch <- message:
			successList = append(successList, queue)
		case <-time.After(h.PublishToChannelDelay):
			failedList = append(failedList, queue)
		}
	}

	if len(failedList) > 0 {
		return ErrChannelPublish(
			fmt.Errorf("failed to publish to one or more queue for subject %s", subject),
			successList,
			failedList,
		)
	}

	return nil
}

// PublishWithChannel - to publish messages with channel
func (h *TMPChannelBrokerHandler) PublishWithChannel(subject string, msgch chan *realBroker.Message) error {
	go func() {
		// as soon as this channel will be closed, for loop will end
		for msg := range msgch {
			h.Publish(subject, msg)
		}
	}()
	return nil
}

// Subscribe - for subscribing messages
func (h *TMPChannelBrokerHandler) Subscribe(subject, queue string, message []byte) error {
	// Looks like current version with nats does not seem to be correct
	// it adds callback which is executed on every message
	// and if queue is populated with messages it keeps waiting, in the end it returns the last one;
	// it does not unsubscribe and callback will keep executing and write to local variable, which probably will cause panic;

	// Not supported

	return nil
}

// SubscribeWithChannel will publish all the messages received to the given channel
func (h *TMPChannelBrokerHandler) SubscribeWithChannel(subject, queue string, msgch chan *realBroker.Message) error {
	if h.storage[subject] == nil {
		h.storage[subject] = make(map[string]chan *realBroker.Message)
	}

	if h.storage[subject][queue] == nil {
		h.storage[subject][queue] = make(chan *realBroker.Message, h.SingleChannelBufferSize)
	}

	go func() {
		// this loop will terminate when the h.storage[subject][queue] is closed
		for message := range h.storage[subject][queue] {
			// this flow is correct as if we have more than one consumer for one queue
			// only one will receive the message
			msgch <- message
		}
	}()

	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (h *TMPChannelBrokerHandler) DeepCopyInto(out realBroker.Handler) {
	// Not supported

	// TODO
	// it is used in meshery server in operator_helper, check if this is code base which is in use
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new Nats.
func (h *TMPChannelBrokerHandler) DeepCopy() realBroker.Handler {
	// Not supported
	return h
}

// DeepCopyObject is a deepcopy function, copying the receiver, creating a new realBroker.Handler.
func (h *TMPChannelBrokerHandler) DeepCopyObject() realBroker.Handler {
	// Not supported
	return h
}

// Check if the connection object is empty
func (h *TMPChannelBrokerHandler) IsEmpty() bool {
	return len(h.storage) <= 0
}
