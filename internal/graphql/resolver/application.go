package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

func (r *Resolver) fetchApplications(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.ApplicationPage, error) {
	tokenString := ctx.Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryApplications(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	applications := &model.ApplicationPage{}

	err = json.Unmarshal(resp, applications)
	if err != nil {
		obj := "result data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return applications, nil
}
