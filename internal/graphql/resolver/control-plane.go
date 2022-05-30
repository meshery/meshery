package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

func (r *Resolver) getControlPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.ControlPlane, error) {
	selectors := make([]model.MeshType, 0)
	if filter.Type == nil {
		for _, mesh := range model.AllMeshType {
			selectors = append(selectors, mesh)
		}
	} else {
		selectors = append(selectors, *filter.Type)
	}
	var cids []string
	if len(filter.K8sServerIDs) != 0 {
		cids = filter.K8sServerIDs
	} else { //This is a fallback
		k8sctxs, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8sctxs) == 0 {
			r.Log.Error(ErrEmptyCurrentK8sContext)
			return nil, ErrEmptyCurrentK8sContext
		}
		for _, context := range k8sctxs {
			if context.KubernetesServerID == nil {
				r.Log.Error(ErrEmptyCurrentK8sContext)
				return nil, ErrEmptyCurrentK8sContext
			}
			cids = append(cids, context.KubernetesServerID.String())
		}

	}
	controlplanelist, err := model.GetControlPlaneState(selectors, provider, cids)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return controlplanelist, nil
}

func (r *Resolver) listenToControlPlaneState(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) (<-chan []*model.ControlPlane, error) {
	if r.controlPlaneChannel == nil {
		r.controlPlaneChannel = make(chan []*model.ControlPlane)
	}

	go func() {
		r.Log.Info("Initializing ControlPlane subscription")

		for {
			select {
			case <-r.MeshSyncChannel:
				status, err := r.getControlPlanes(ctx, provider, filter)
				if err != nil {
					r.Log.Error(ErrControlPlaneSubscription(err))
					break
				}
				r.controlPlaneChannel <- status
			case <-ctx.Done():
				r.Log.Info("ControlPlane subscription stopped")
				return
			}
		}
	}()
	return r.controlPlaneChannel, nil
}
