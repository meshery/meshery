package resolver

import (
	"context"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

func (r *Resolver) getControlPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.ControlPlane, error) {
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
	controlplanelist, err := model.GetControlPlaneState(ctx, selectors, provider, cids)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return controlplanelist, nil
}
