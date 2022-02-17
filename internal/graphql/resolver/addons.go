package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (r *Resolver) changeAddonStatus(ctx context.Context, provider models.Provider) (model.Status, error) {
	return model.StatusProcessing, nil
}

func (r *Resolver) getAvailableAddons(ctx context.Context, provider models.Provider, selector *model.MeshType) ([]*model.AddonList, error) {
	kubeclient, ok := ctx.Value(models.KubeHanderKey).(*mesherykube.Client)
	if !ok {
		r.Log.Error(ErrNilClient)
		return nil, ErrNilClient
	}

	selectors := make([]model.MeshType, 0)
	if selector == nil || *selector == model.MeshTypeAllMesh {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		selectors = append(selectors, *selector)
	}

	addonlist, err := model.GetAddonsState(selectors, kubeclient, provider)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return addonlist, nil
}

func (r *Resolver) listenToAddonState(ctx context.Context, provider models.Provider, selector *model.MeshType) (<-chan []*model.AddonList, error) {
	if r.addonChannel == nil {
		r.addonChannel = make(chan []*model.AddonList, 0)
	}

	go func() {
		r.Log.Info("Addons subscription started")
		err := r.connectToBroker(ctx, provider)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		for {
			select {
			case <-r.MeshSyncChannel:
				status, err := r.getAvailableAddons(ctx, provider, selector)
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
	}()

	return r.addonChannel, nil
}
