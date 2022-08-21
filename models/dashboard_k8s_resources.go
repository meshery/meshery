package models

import (
	"sync"

	"github.com/layer5io/meshery/helpers/utils"
)

type DashboardK8sResourcesChan struct {
	ResourcesChan   map[string][]chan struct{}
	mx             sync.Mutex
}

func NewDashboardK8sResourcesHelper() *DashboardK8sResourcesChan {
	return &DashboardK8sResourcesChan{
		ResourcesChan:  make(map[string][]chan struct{}, 0),
	}
}

func (d *DashboardK8sResourcesChan) SubscribeDashbordK8Resources(ch chan struct{}, ctxIDs []string) {
	d.mx.Lock()
	defer d.mx.Unlock()

	for _, ctxID := range ctxIDs {
		if _, ok := d.ResourcesChan[ctxID]; !ok {
			d.ResourcesChan[ctxID] = append(d.ResourcesChan[ctxID], ch)
		}
	}
}

func (d *DashboardK8sResourcesChan) PublishDashboardK8sResources(k8scontextID string) {
	for ctxID, ch := range d.ResourcesChan {
		if ctxID == k8scontextID {
			for _, c := range ch {
				if !utils.IsClosed(c){
					c <- struct{}{}
				}	
			}
		}
	}
}
