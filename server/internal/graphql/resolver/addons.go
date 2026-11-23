package resolver

import (
	"context"

	"github.com/meshery/meshery/server/internal/graphql/dataloader"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

func (r *Resolver) getAvailableAddons(ctx context.Context, provider models.GenericPersisterGetter, filter *model.ServiceMeshFilter) ([]*model.AddonList, error) {
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

	seenCids := make(map[string]struct{}, len(cids))
	uniqCids := make([]string, 0, len(cids))
	for _, c := range cids {
		if _, ok := seenCids[c]; !ok {
			seenCids[c] = struct{}{}
			uniqCids = append(uniqCids, c)
		}
	}
	cids = uniqCids

	loader := dataloader.For(ctx)
	if loader != nil && loader.AddonLoader != nil && len(cids) > 0 {
		clusterResults, err := loader.AddonLoader.LoadMany(cids, selectors)
		if err != nil {
			if r.Log != nil {
				r.Log.Error(err)
			}
			return nil, err
		}

		addonlist := make([]*model.AddonList, 0)
		for _, selector := range selectors {
			selectorOwner := selector.String()
			for _, cid := range cids {
				for _, addon := range clusterResults[cid] {
					if addon.Owner == selectorOwner {
						addonlist = append(addonlist, addon)
					}
				}
			}
		}
		return addonlist, nil
	}

	addonlist, err := model.GetAddonsState(ctx, selectors, provider, cids)
	if err != nil {
		if r.Log != nil {
			r.Log.Error(err)
		}
		return nil, err
	}

	return addonlist, nil
}
