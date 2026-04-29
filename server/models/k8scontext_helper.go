package models

import (
	"sync"

	"github.com/meshery/meshery/server/helpers/utils"
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
		if !utils.IsClosed(ch) {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
	}
}
