package models

import (
	"sync"
)

type K8scontextChan struct {
	mx          sync.Mutex
	contextchan []chan struct{}
}

func NewContextHelper() *K8scontextChan {
	return &K8scontextChan{
		contextchan: make([]chan struct{}, 0),
	}
}

func (k *K8scontextChan) SubscribeContext(ch chan struct{}) {
	k.mx.Lock()
	defer k.mx.Unlock()

	k.contextchan = append(k.contextchan, ch)
}

func (k *K8scontextChan) UnsubscribeContext(ch chan struct{}) {
	k.mx.Lock()
	defer k.mx.Unlock()

	for i := 0; i < len(k.contextchan); i++ {
		if k.contextchan[i] == ch {
			k.contextchan = append(k.contextchan[:i], k.contextchan[i+1:]...)
			return
		}
	}
}

func (k *K8scontextChan) PublishContext() {
	k.mx.Lock()
	subscribers := make([]chan struct{}, len(k.contextchan))
	copy(subscribers, k.contextchan)
	k.mx.Unlock()

	for _, ch := range subscribers {
		// Non-blocking coalescing send. A subscriber that already has a pending
		// signal buffered doesn't need another; a departed subscriber is removed
		// via UnsubscribeContext and its channel is never closed, so this can
		// neither block the publisher nor panic on a closed channel. The former
		// utils.IsClosed guard was a destructive receive that silently consumed a
		// buffered signal, dropping the refetch it was meant to trigger.
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
