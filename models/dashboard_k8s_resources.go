package models

import (
	"sync"
	"context"

	"github.com/layer5io/meshery/helpers/utils"
	"github.com/sirupsen/logrus"
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
		if resChannels, ok := d.ResourcesChan[ctxID]; ok {
			_ch := resChannels
			_ch = append(_ch, ch)
			d.ResourcesChan[ctxID] = _ch
		} else {
			_ch := make([]chan struct{}, 0)
			_ch = append(_ch, ch)
			d.ResourcesChan[ctxID] = _ch
		}
	}
	logrus.Debug("chans: ", d.ResourcesChan)
}

func (d *DashboardK8sResourcesChan) PublishDashboardK8sResources(ctx context.Context, k8scontextID string) {
	var cid string
	k8sCtxs, ok := ctx.Value(AllKubeClusterKey).([]K8sContext)
	if !ok || len(k8sCtxs) == 0 {
		return
	}
	for _, k8sContext := range k8sCtxs {
		if k8sContext.KubernetesServerID != nil && k8sContext.ID == k8scontextID {
			cid = k8sContext.KubernetesServerID.String()
		}
	}
	for ctxID, ch := range d.ResourcesChan {
		logrus.Debug("given ctxId: ", cid)
		logrus.Debug("id: ", ctxID)
		if ctxID == cid {
			logrus.Debug("match found")
			for _, c := range ch {
				if !utils.IsClosed(c){
					logrus.Debug("publishing to ch: ", ch)
					c <- struct{}{}
				}	
			}
		}
	}
}
