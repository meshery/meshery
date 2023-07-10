package resolver

import (
	"context"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) getDataPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.DataPlane, error) {
	selectors := make([]model.MeshType, 0)
	if filter.Type == nil {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		selectors = append(selectors, *filter.Type)
	}
	var cids []string
	if len(filter.K8sClusterIDs) != 0 {
		cids = filter.K8sClusterIDs
	}
	dataPlaneList, err := model.GetDataPlaneState(ctx, selectors, provider, cids)
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
		for _, ctxID := range filter.K8sClusterIDs {
			go func(ctxID string) {
				for {
					select {
					case <-r.MeshSyncChannelPerK8sContext[ctxID]:
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
			}(ctxID)
		}
	}()
	return r.dataPlaneChannel, nil
}
