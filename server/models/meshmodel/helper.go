package meshmodel

import (
	"sync"
)

type SummaryChannel struct {
	channel []chan struct{}
	mx      sync.Mutex
}

func NewSummaryHelper() *SummaryChannel {
	return &SummaryChannel{
		channel: make([]chan struct{}, 10),
	}
}

func (c *SummaryChannel) Subscribe(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.channel = append(c.channel, ch)
}

func (c *SummaryChannel) Publish() {
	c.mx.Lock()
	subscribers := make([]chan struct{}, len(c.channel))
	copy(subscribers, c.channel)
	c.mx.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
