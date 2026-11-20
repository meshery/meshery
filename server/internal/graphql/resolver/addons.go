package resolver

import (
	"context"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

func (r *Resolver) getAvailableAddons(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.AddonList, error) {
	var cids []string
	if filter != nil && len(filter.K8sClusterIDs) != 0 {
		cids = filter.K8sClusterIDs
	}

	selectors := make([]model.MeshType, 0)
	if filter == nil || filter.Type == nil || *filter.Type == model.MeshTypeAllMesh {
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
