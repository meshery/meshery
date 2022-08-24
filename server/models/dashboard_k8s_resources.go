package models

import (
	"sync"

	"github.com/layer5io/meshery/server/helpers/utils"
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
	for _, ch := range d.ResourcesChan {
		if !utils.IsClosed(ch) {
			ch <- struct{}{}
		}
	}
}
