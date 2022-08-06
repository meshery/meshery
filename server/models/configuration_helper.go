package models

import (
	"sync"

	"github.com/layer5io/meshery/server/helpers/utils"
)

type ConfigurationChannel struct {
	ApplicationsChannel []chan struct{}
	PatternsChannel     []chan struct{}
	FiltersChannel      []chan struct{}
	mx                  sync.Mutex
}

func NewConfigurationHelper() *ConfigurationChannel {
	return &ConfigurationChannel{
		ApplicationsChannel: make([]chan struct{}, 10),
		PatternsChannel:     make([]chan struct{}, 10),
		FiltersChannel:      make([]chan struct{}, 10),
	}
}

func (c *ConfigurationChannel) SubscribeApplications(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.ApplicationsChannel = append(c.ApplicationsChannel, ch)
}

func (c *ConfigurationChannel) PublishApplications() {
	for _, ch := range c.ApplicationsChannel {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}

func (c *ConfigurationChannel) SubscribePatterns(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.PatternsChannel = append(c.PatternsChannel, ch)
}

func (c *ConfigurationChannel) PublishPatterns() {
	for _, ch := range c.PatternsChannel {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}

func (c *ConfigurationChannel) SubscribeFilters(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.FiltersChannel = append(c.FiltersChannel, ch)
}

func (c *ConfigurationChannel) PublishFilters() {
	for _, ch := range c.FiltersChannel {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}
