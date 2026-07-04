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

// UnsubscribeDashboardK8sResources removes a listener registered with
// SubscribeDashbordK8Resources. Subscribers must call it when their context is
// done, otherwise ResourcesChan grows without bound as subscriptions come and go.
func (d *DashboardK8sResourcesChan) UnsubscribeDashboardK8sResources(ch chan struct{}) {
	d.mx.Lock()
	defer d.mx.Unlock()

	for i := 0; i < len(d.ResourcesChan); i++ {
		if d.ResourcesChan[i] == ch {
			d.ResourcesChan = append(d.ResourcesChan[:i], d.ResourcesChan[i+1:]...)
			return
		}
	}
}

func (d *DashboardK8sResourcesChan) PublishDashboardK8sResources() {
	d.mx.Lock()
	subscribers := make([]chan struct{}, len(d.ResourcesChan))
	copy(subscribers, d.ResourcesChan)
	d.mx.Unlock()

	for _, ch := range subscribers {
		// Non-blocking coalescing send - see the note in K8scontextChan.
		// PublishContext. The former utils.IsClosed guard destructively consumed
		// buffered signals and could block the publisher on a slow subscriber.
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}
