package models

import (
	"sync"
)

type DashboardK8sResourcesChan struct {
	ResourcesChan []chan struct{}
	mx            sync.Mutex
}

func NewDashboardK8sResourcesHelper() *DashboardK8sResourcesChan {
	return &DashboardK8sResourcesChan{
		ResourcesChan: make([]chan struct{}, 0),
	}
}

func (d *DashboardK8sResourcesChan) SubscribeDashbordK8Resources(ch chan struct{}) {
	d.mx.Lock()
	defer d.mx.Unlock()

	d.ResourcesChan = append(d.ResourcesChan, ch)
}

func (d *DashboardK8sResourcesChan) PublishDashboardK8sResources() {
	d.mx.Lock()
	subscribers := make([]chan struct{}, len(d.ResourcesChan))
	copy(subscribers, d.ResourcesChan)
	d.mx.Unlock()

	for _, ch := range subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
