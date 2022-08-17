package models

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type DashboardK8sResourcesChan struct {
	ResourcesChan   []chan struct{}
	mx             sync.Mutex
}

func NewDashboardK8sResourcesHelper() *DashboardK8sResourcesChan {
	return &DashboardK8sResourcesChan{
		ResourcesChan:  make([]chan struct{}, 0),
	}
}

func (d *DashboardK8sResourcesChan) SubscribeDashbordK8Resources(ch chan struct{}) {
	d.mx.Lock()
	defer d.mx.Unlock()
	logrus.Debug("subscribing resources")
	d.ResourcesChan = append(d.ResourcesChan, ch)
}

func (d *DashboardK8sResourcesChan) PublishDashboardK8sResources() {
	logrus.Debug("publishing resources")
	logrus.Debug("resourcesChan: ", d.ResourcesChan)
	for _, ch := range d.ResourcesChan {
		logrus.Debug("ch(helper): ", ch)
		if ch != nil {
			logrus.Debug("pushing: ", ch)
			ch <- struct{}{}
			logrus.Debug("pushed")
		}	
	}
}
