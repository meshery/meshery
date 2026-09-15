package models

import (
	"sync"
)

type ConfigurationChannel struct {
	ApplicationsChannel []chan struct{}
	PatternsChannel     []chan struct{}
	FiltersChannel      []chan struct{}
	mx                  sync.Mutex
}

func NewConfigurationHelper() *ConfigurationChannel {
	// Start empty; the previous length-10 make seeded each slice with 10 nil
	// channels that Publish then had to skip on every call.
	return &ConfigurationChannel{
		ApplicationsChannel: make([]chan struct{}, 0),
		PatternsChannel:     make([]chan struct{}, 0),
		FiltersChannel:      make([]chan struct{}, 0),
	}
}

func (c *ConfigurationChannel) SubscribeApplications(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.ApplicationsChannel = append(c.ApplicationsChannel, ch)
}

func (c *ConfigurationChannel) PublishApplications() {
	c.mx.Lock()
	subscribers := make([]chan struct{}, len(c.ApplicationsChannel))
	copy(subscribers, c.ApplicationsChannel)
	c.mx.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (c *ConfigurationChannel) SubscribePatterns(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.PatternsChannel = append(c.PatternsChannel, ch)
}

func (c *ConfigurationChannel) PublishPatterns() {
	c.mx.Lock()
	subscribers := make([]chan struct{}, len(c.PatternsChannel))
	copy(subscribers, c.PatternsChannel)
	c.mx.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

func (c *ConfigurationChannel) SubscribeFilters(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.FiltersChannel = append(c.FiltersChannel, ch)
}

func (c *ConfigurationChannel) PublishFilters() {
	c.mx.Lock()
	subscribers := make([]chan struct{}, len(c.FiltersChannel))
	copy(subscribers, c.FiltersChannel)
	c.mx.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
