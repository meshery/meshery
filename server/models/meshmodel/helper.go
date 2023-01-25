package meshmodel

import (
	"sync"

	"github.com/layer5io/meshery/server/helpers/utils"
)

type MeshModelSummaryChannel struct {
	channel []chan struct{}
	mx      sync.Mutex
}

func NewMeshModelSummaryHelper() *MeshModelSummaryChannel {
	return &MeshModelSummaryChannel{
		channel: make([]chan struct{}, 10),
	}
}

func (c *MeshModelSummaryChannel) Subscribe(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.channel = append(c.channel, ch)
}

func (c *MeshModelSummaryChannel) Publish() {
	for _, ch := range c.channel {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}
