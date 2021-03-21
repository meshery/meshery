package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func (r *Resolver) getAvailableNamespaces(ctx context.Context) ([]*model.NameSpace, error) {
	resourceobjects := make([]meshsyncmodel.ResourceObjectMeta, 0)

	result := r.DBHandler.Distinct("namespace").Not("namespace = ?", "").Find(&resourceobjects)
	if result.Error != nil {
		r.Log.Error(result.Error)
		return nil, result.Error
	}
	namespaces := make([]*model.NameSpace, 0)
	for _, obj := range resourceobjects {
		namespaces = append(namespaces, &model.NameSpace{
			Namespace: obj.Namespace,
		})
	}

	return namespaces, nil
}
