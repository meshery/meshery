package models

import (
	"sync"

	"github.com/gofrs/uuid"
)

type ConfigurationChannel struct {
	ApplicationsChannel map[uuid.UUID]chan struct{}
	PatternsChannel     map[uuid.UUID]chan struct{}
	FiltersChannel      map[uuid.UUID]chan struct{}
	mx                  sync.Mutex
}

func NewConfigurationHelper() *ConfigurationChannel {
	return &ConfigurationChannel{
		ApplicationsChannel: make(map[uuid.UUID]chan struct{}, 10),
		PatternsChannel:     make(map[uuid.UUID]chan struct{}, 10),
		FiltersChannel:      make(map[uuid.UUID]chan struct{}, 10),
	}
}

func (c *ConfigurationChannel) PublishApplications() {
	for _, ch := range c.ApplicationsChannel {
		ch <- struct{}{}
	}
}

func (c *ConfigurationChannel) SubscribeApplications(ch chan struct{}, key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.ApplicationsChannel[key] = ch
}

func (c *ConfigurationChannel) UnSubscribeApplications(key uuid.UUID) {
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

func (c *ConfigurationChannel) SubscribePatterns(ch chan struct{}, key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.PatternsChannel[key] = ch
}

func (c *ConfigurationChannel) UnSubscribePatterns(key uuid.UUID) {
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

func (c *ConfigurationChannel) SubscribeFilters(ch chan struct{}, key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.FiltersChannel[key] = ch
}

func (c *ConfigurationChannel) UnSubscribeFilters(key uuid.UUID) {
	c.mx.Lock()
	defer c.mx.Unlock()
	close(c.FiltersChannel[key])
	delete(c.FiltersChannel, key)
}

