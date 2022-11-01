package resolver

import (
	"context"
	"encoding/json"

	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) fetchPatterns(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.PatternPageResult, error) {
	tokenString := ctx.Value(models.TokenCtxKey).(string)
	user := ctx.Value(models.UserCtxKey).(*models.User)
	prefObj := ctx.Value(models.PerfObjCtxKey).(*models.Preference)

	resp, err := provider.GetMesheryPatterns(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order)

	if err != nil {
		r.Log.Error(ErrFetchingPatterns(err))
		return nil, err
	}

	mc := handlers.NewContentModifier(tokenString, provider, prefObj, user.UserID)
	err = mc.AddMetadataForPatterns(ctx, &resp)
	if err != nil {
		r.Log.Error(ErrFetchingPatterns(err))
	}

	patterns := &model.PatternPageResult{}

	if err := json.Unmarshal(resp, patterns); err != nil {
		obj := "result data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return patterns, nil
}
