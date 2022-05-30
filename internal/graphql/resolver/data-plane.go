package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

func (r *Resolver) getDataPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.DataPlane, error) {
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
	dataPlaneList, err := model.GetDataPlaneState(selectors, provider, cids)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return dataPlaneList, nil
}

func (r *Resolver) listenToDataPlaneState(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) (<-chan []*model.DataPlane, error) {
	if r.dataPlaneChannel == nil {
		r.dataPlaneChannel = make(chan []*model.DataPlane)
	}

	go func() {
		r.Log.Info("Initializing DataPlane subscription")

		for {
			select {
			case <-r.MeshSyncChannel:
				containers, err := r.getDataPlanes(ctx, provider, filter)
				if err != nil {
					r.Log.Error(ErrDataPlaneSubscription(err))
					break
				}
				r.dataPlaneChannel <- containers
			case <-ctx.Done():
				r.Log.Info("DataPlane subscription stopped")
				return
			}
		}
	}()
	return r.dataPlaneChannel, nil
}
