package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *queryResolver) fetchCatalogPattern(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogPattern, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetCatalogMesheryPatterns(token, selector.Search, selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}
	var catalog []*model.CatalogPattern

	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(handlers.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog, nil
}

func (r *queryResolver) fetchCatalogFilter(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogFilter, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetCatalogMesheryFilters(token, selector.Search, selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	var catalog []*model.CatalogFilter

	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(handlers.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog, nil
}
