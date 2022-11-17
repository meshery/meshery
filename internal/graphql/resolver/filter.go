package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

func (r *Resolver) fetchFilters(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.FilterPage, error) {
	tokenString := ctx.Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryFilters(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	filters := &model.FilterPage{}

	err = json.Unmarshal(resp, filters)
	if err != nil {
		obj := "result data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return filters, nil
}
