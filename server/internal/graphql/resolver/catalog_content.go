package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

// CatalogPatternPage - represents a page of meshery patterns
type catalogPatternPage struct {
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalCount int                     `json:"total_count"`
	Patterns   []*model.CatalogPattern `json:"patterns"`
}

// CatalogFilterPage - represents a page of meshery filters
type catalogFilterPage struct {
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	TotalCount int                    `json:"total_count"`
	Filters    []*model.CatalogFilter `json:"filters"`
}

func (r *queryResolver) fetchCatalogPattern(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogPattern, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	metrics := "false"

	// Convert []*string to []string for class
	class := []string{}
	for _, classPtr := range selector.Class {
		if classPtr != nil {
			class = append(class, *classPtr)
		}
	}

	// Convert []*string to []string for technology
	technology := []string{}
	for _, techPtr := range selector.Technology {
		if techPtr != nil {
			technology = append(technology, *techPtr)
		}
	}

	// Convert []*string to []string for patternType
	patternType := []string{}
	for _, typePtr := range selector.PatternType {
		if typePtr != nil {
			patternType = append(patternType, *typePtr)
		}
	}

	orgID := []string{}
	for _, orgIDPtr := range selector.OrgID {
		if orgIDPtr != nil {
			orgID = append(orgID, *orgIDPtr)
		}
	}

	workspaceID := []string{}
	for _, workspaceIDPtr := range selector.WorkspaceID {
		if workspaceIDPtr != nil {
			workspaceID = append(workspaceID, *workspaceIDPtr)
		}
	}

	userid := []string{}
	for _, useridPtr := range selector.Userid {
		if useridPtr != nil {
			userid = append(userid, *useridPtr)
		}
	}

	populate := []string{}
	for _, populatePtr := range selector.Populate {
		if populatePtr != nil {
			populate = append(populate, *populatePtr)
		}
	}

	resp, err := provider.GetCatalogMesheryPatterns(token, selector.Page, selector.Pagesize, selector.Search, selector.Order, metrics, populate, class, technology, patternType, orgID, workspaceID, userid)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}
	var catalog catalogPatternPage
	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(models.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog.Patterns, nil
}

func (r *queryResolver) fetchCatalogFilter(ctx context.Context, provider models.Provider, selector *model.CatalogSelector) ([]*model.CatalogFilter, error) {
	token := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetCatalogMesheryFilters(token, selector.Page, selector.Pagesize, selector.Search, selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	var catalog catalogFilterPage

	err = json.Unmarshal(resp, &catalog)
	if err != nil {
		r.Log.Error(models.ErrUnmarshal(err, "catalog data"))
		return nil, err
	}
	return catalog.Filters, nil
}
