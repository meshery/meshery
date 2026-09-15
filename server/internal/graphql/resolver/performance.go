package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

func (r *Resolver) getPerfResult(ctx context.Context, provider models.Provider, id string) (*model.MesheryResult, error) {
	if id == "" {
		return nil, handlers.ErrQueryGet("*id")
	}

	resultID, err := uuid.FromString(id)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	tokenString, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok || tokenString == "" {
		return nil, ErrInvalidRequest
	}

	bdr, err := provider.GetResult(tokenString, resultID)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	startTime := fmt.Sprintf("%v", bdr.TestStartTime)
	serverBoardConfig := fmt.Sprintf("%v", bdr.ServerBoardConfig)
	serverMetrics := fmt.Sprintf("%v", bdr.ServerMetrics)
	mesheryID := fmt.Sprintf("%v", bdr.ID)
	performanceProfile := fmt.Sprintf("%v", bdr.PerformanceProfileInfo.ID)

	return &model.MesheryResult{
		MesheryID:          &mesheryID,
		Name:               &bdr.Name,
		Mesh:               &bdr.Mesh,
		PerformanceProfile: &performanceProfile,
		TestID:             &bdr.TestID,
		RunnerResults:      bdr.Result,
		ServerMetrics:      &serverMetrics,
		ServerBoardConfig:  &serverBoardConfig,
		TestStartTime:      &startTime,
		Owner:              &bdr.Owner,
		UpdatedAt:          &bdr.UpdatedAt,
		CreatedAt:          &bdr.CreatedAt,
	}, nil
}

func (r *Resolver) fetchResults(ctx context.Context, provider models.Provider, selector model.PageFilter, profileID string) (*model.PerfPageResult, error) {
	if profileID == "" {
		return nil, handlers.ErrQueryGet("*profileID")
	}

	tokenString, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok || tokenString == "" {
		return nil, ErrInvalidRequest
	}
	search := ""
	if selector.Search != nil {
		search = *selector.Search
	}
	order := ""
	if selector.Order != nil {
		order = *selector.Order
	}

	bdr, err := provider.FetchResults(tokenString, selector.Page, selector.PageSize, search, order, profileID)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	result := &model.PerfPageResult{}

	if err := json.Unmarshal(bdr, result); err != nil {
		obj := "result data"
		return nil, models.ErrUnmarshal(err, obj)
	}

	return result, nil
}

func (r *Resolver) getPerformanceProfiles(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.PerfPageProfiles, error) {
	tokenString, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok || tokenString == "" {
		return nil, ErrInvalidRequest
	}
	search := ""
	if selector.Search != nil {
		search = *selector.Search
	}
	order := ""
	if selector.Order != nil {
		order = *selector.Order
	}

	bdr, err := provider.GetPerformanceProfiles(tokenString, selector.Page, selector.PageSize, search, order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	profiles := &model.PerfPageProfiles{}

	if err := json.Unmarshal(bdr, profiles); err != nil {
		obj := "performance profiles data"
		return nil, models.ErrUnmarshal(err, obj)
	}

	return profiles, nil
}

func (r *Resolver) fetchAllResults(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.PerfPageResult, error) {
	tokenString, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok || tokenString == "" {
		return nil, ErrInvalidRequest
	}
	search := ""
	if selector.Search != nil {
		search = *selector.Search
	}
	order := ""
	if selector.Order != nil {
		order = *selector.Order
	}
	from := ""
	if selector.From != nil {
		from = *selector.From
	}
	to := ""
	if selector.To != nil {
		to = *selector.To
	}

	bdr, err := provider.FetchAllResults(tokenString, selector.Page, selector.PageSize, search, order, from, to)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	performanceResults := &model.PerfPageResult{}

	if err := json.Unmarshal(bdr, performanceResults); err != nil {
		obj := "performance results data"
		return nil, models.ErrUnmarshal(err, obj)
	}

	return performanceResults, nil
}
