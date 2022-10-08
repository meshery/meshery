package resolver

import (
	"context"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) getAvailableAddons(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.AddonList, error) {
	var cids []string
	if len(filter.K8sClusterIDs) != 0 {
		cids = filter.K8sClusterIDs
	}

	selectors := make([]model.MeshType, 0)
	if filter == nil || *filter.Type == model.MeshTypeAllMesh {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		selectors = append(selectors, *filter.Type)
	}

	addonlist, err := model.GetAddonsState(ctx, selectors, provider, cids)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return addonlist, nil
}

func (r *Resolver) listenToAddonState(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) (<-chan []*model.AddonList, error) {
	if r.addonChannel == nil {
		r.addonChannel = make(chan []*model.AddonList, 0)
	}

	go func() {
		r.Log.Info("Addons subscription started")
		for _, ctxID := range filter.K8sClusterIDs {
			go func(ctxID string) {
				for {
					select {
					case <-r.MeshSyncChannelPerK8sContext[ctxID]:
						status, err := r.getAvailableAddons(ctx, provider, filter)
						if err != nil {
							r.Log.Error(ErrAddonSubscription(err))
							break
						}
						r.addonChannel <- status
					case <-ctx.Done():
						r.Log.Info("Addons subscription stopped")
						return
					}
				}
			}(ctxID)
		}
	}()

	return r.addonChannel, nil
}
