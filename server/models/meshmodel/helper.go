package meshmodel

import (
	"sync"
)

type SummaryChannel struct {
	channel []chan struct{}
	mx      sync.Mutex
}

func NewSummaryHelper() *SummaryChannel {
	// Start empty. The previous length-10 make seeded the slice with 10 nil
	// channels that Publish then had to skip on every call.
	return &SummaryChannel{
		channel: make([]chan struct{}, 0),
	}
}

func (c *SummaryChannel) Subscribe(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.channel = append(c.channel, ch)
}

// Unsubscribe removes a listener registered with Subscribe. Subscribers must
// call it when their context is done: it detaches the channel so Publish stops
// sending to it (and so the slice doesn't grow without bound). Removing rather
// than closing keeps Publish's send panic-free.
func (c *SummaryChannel) Unsubscribe(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	for i := 0; i < len(c.channel); i++ {
		if c.channel[i] == ch {
			c.channel = append(c.channel[:i], c.channel[i+1:]...)
			return
		}
	}
}

func (c *SummaryChannel) Publish() {
	c.mx.Lock()
	subscribers := make([]chan struct{}, len(c.channel))
	copy(subscribers, c.channel)
	c.mx.Unlock()

	for _, ch := range subscribers {
		// Non-blocking coalescing send. The former utils.IsClosed guard did a
		// destructive receive that consumed a buffered refetch signal (dropping
		// the update) and did not actually make the following blocking send safe.
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
