package models

import "github.com/layer5io/meshery/server/helpers/utils"

type K8scontextChan struct {
	contextchan []chan struct{}
}

func NewContextHelper() *K8scontextChan {
	return &K8scontextChan{
		contextchan: make([]chan struct{}, 0),
	}
}

func (k *K8scontextChan) SubscribeContext(ch chan struct{}) {
	k.contextchan = append(k.contextchan, ch)
}

func (k *K8scontextChan) PublishContext() {
	for _, ch := range k.contextchan {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}
