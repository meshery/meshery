package models

import (
	"sync"
)

type ConfigurationChannel struct {
	ApplicationsChannel map[int]chan struct{}
	PatternsChannel     map[int]chan struct{}
	FiltersChannel      map[int]chan struct{}
	mx                  sync.Mutex
}

func NewConfigurationHelper() *ConfigurationChannel {
	return &ConfigurationChannel{
		ApplicationsChannel: make(map[int]chan struct{}, 10),
		PatternsChannel:     make(map[int]chan struct{}, 10),
		FiltersChannel:      make(map[int]chan struct{}, 10),
	}
}

func (c *ConfigurationChannel) PublishApplications() {
	for _, ch := range c.ApplicationsChannel {
		ch <- struct{}{}
	}
}

func (c *ConfigurationChannel) SubscribeApplications(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.ApplicationsChannel[len(c.ApplicationsChannel)] = ch
}

func (c *ConfigurationChannel) UnSubscribeApplications(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	close(c.ApplicationsChannel[key])
	delete(c.ApplicationsChannel, key)
}

func (c *ConfigurationChannel) PublishPatterns() {
	for _, ch := range c.PatternsChannel {
		ch <- struct{}{}		
	}
}

func (c *ConfigurationChannel) SubscribePatterns(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.PatternsChannel[len(c.PatternsChannel)] = ch
}

func (c *ConfigurationChannel) UnSubscribePatterns(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	close(c.PatternsChannel[key])
	delete(c.PatternsChannel, key)
}

func (c *ConfigurationChannel) PublishFilters() {
	for _, ch := range c.FiltersChannel {
		ch <- struct{}{}
	}
}

func (c *ConfigurationChannel) SubscribeFilters(ch chan struct{}) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.FiltersChannel[len(c.FiltersChannel)] = ch
}

func (c *ConfigurationChannel) UnSubscribeFilters(key int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	close(c.FiltersChannel[key])
	delete(c.FiltersChannel, key)
}

