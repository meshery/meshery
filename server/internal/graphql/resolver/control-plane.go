package resolver

import (
	"context"
	"strings"

	"github.com/meshery/meshery/server/internal/graphql/dataloader"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

func (r *Resolver) getControlPlanes(ctx context.Context, provider models.GenericPersisterGetter, filter *model.ServiceMeshFilter) ([]*model.ControlPlane, error) {
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
	if loader != nil && loader.ControlPlaneLoader != nil && len(cids) > 0 {
		clusterResults, err := loader.ControlPlaneLoader.LoadMany(cids, selectors)
		if err != nil {
			r.Log.Error(err)
			return nil, err
		}

		controlplanelist := make([]*model.ControlPlane, 0)
		for _, selector := range selectors {
			selectorName := strings.ToLower(selector.String())
			var members []*model.ControlPlaneMember
			for _, cid := range cids {
				for _, cp := range clusterResults[cid] {
					if cp.Name == selectorName {
						members = append(members, cp.Members...)
					}
				}
			}
			controlplanelist = append(controlplanelist, &model.ControlPlane{
				Name:    selectorName,
				Members: members,
			})
		}
		return controlplanelist, nil
	}

	controlplanelist, err := model.GetControlPlaneState(ctx, selectors, provider, cids)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return controlplanelist, nil
}
