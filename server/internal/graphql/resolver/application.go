package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) fetchApplications(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.ApplicationPage, error) {
	tokenString := ctx.Value(models.TokenCtxKey).(string)
	var updateAfter string
	if selector.UpdatedAfter != nil {
		updateAfter = *selector.UpdatedAfter
	}
	var order string
	if selector.Order != nil {
		order = *selector.Order
	}
	var search string
	if selector.Search != nil {
		search = *selector.Search
	}
	resp, err := provider.GetMesheryApplications(tokenString, selector.Page, selector.PageSize, search, order, updateAfter)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	applications := &model.ApplicationPage{}

	err = json.Unmarshal(resp, applications)
	if err != nil {
		obj := "result data"
		return nil, models.ErrUnmarshal(err, obj)
	}

	return applications, nil
}
