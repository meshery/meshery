package resolver

import (
	"context"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/meshery/meshery/internal/graphql/model"
	"github.com/meshery/meshery/models"
)

func (r *Resolver) getAvailableNamespaces(ctx context.Context, provider models.Provider) ([]*model.NameSpace, error) {
	resourceobjects := make([]meshsyncmodel.ResourceObjectMeta, 0)

	result := provider.GetGenericPersister().Distinct("namespace").Not("namespace = ?", "").Find(&resourceobjects)
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
