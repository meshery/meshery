package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

// CatalogPatternPage - represents a page of meshery patterns
type catalogPatternPage struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int               `json:"total_count"`
	Patterns   []*model.CatalogPattern `json:"patterns"`
}

// CatalogFilterPage - represents a page of meshery filters
type catalogFilterPage struct {
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalCount int              `json:"total_count"`
	Filters    []*model.CatalogFilter `json:"filters"`
}

func (r *queryResolver) fetchCatalogPattern(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogPattern, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetCatalogMesheryPatterns(token, selector.Search, selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}
	var catalog catalogPatternPage
	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(handlers.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog.Patterns, nil
}

func (r *queryResolver) fetchCatalogFilter(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogFilter, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetCatalogMesheryFilters(token, selector.Search, selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	var catalog catalogFilterPage

	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(handlers.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog.Filters, nil
}
