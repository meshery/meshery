package models

import (
	"github.com/layer5io/meshkit/broker"
)

type TmpMeshsyncBrokerHandler struct {
	transport <-chan *broker.Message
}

func NewTmpMeshsyncBrokerHandler(transport <-chan *broker.Message) *TmpMeshsyncBrokerHandler {
	return &TmpMeshsyncBrokerHandler{
		transport: transport,
	}
}
func (h *TmpMeshsyncBrokerHandler) ConnectedEndpoints() (endpoints []string) {
	return []string{"meshery broker stub"}
}

func (h *TmpMeshsyncBrokerHandler) Info() string {
	return "meshery broker stub"
}

func (h *TmpMeshsyncBrokerHandler) CloseConnection() {
}

// Publish - to publish messages
func (h *TmpMeshsyncBrokerHandler) Publish(subject string, message *broker.Message) error {
	return nil
}

// PublishWithChannel - to publish messages with channel
func (h *TmpMeshsyncBrokerHandler) PublishWithChannel(subject string, msgch chan *broker.Message) error {
	return nil
}

// Subscribe - for subscribing messages
// TODO Ques: Do we want to unsubscribe
// TODO will the method-user just subsribe, how will it handle the received messages?
func (h *TmpMeshsyncBrokerHandler) Subscribe(subject, queue string, message []byte) error {
	return nil
}

// SubscribeWithChannel will publish all the messages received to the given channel
func (h *TmpMeshsyncBrokerHandler) SubscribeWithChannel(subject, queue string, msgch chan *broker.Message) error {
	go func() {
		for val := range h.transport {
			msgch <- val
		}
	}()

	return nil
}

// DeepCopyInto is a deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (h *TmpMeshsyncBrokerHandler) DeepCopyInto(out broker.Handler) {
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new Nats.
func (h *TmpMeshsyncBrokerHandler) DeepCopy() *TmpMeshsyncBrokerHandler {
	return h
}

// DeepCopyObject is a deepcopy function, copying the receiver, creating a new broker.Handler.
func (h *TmpMeshsyncBrokerHandler) DeepCopyObject() broker.Handler {
	return h
}

// Check if the connection object is empty
func (h *TmpMeshsyncBrokerHandler) IsEmpty() bool {
	return false
}
