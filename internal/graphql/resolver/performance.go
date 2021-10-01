package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
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

	tokenString := ctx.Value("token").(string)

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
		UserID:             &bdr.UserID,
		UpdatedAt:          &bdr.UpdatedAt,
		CreatedAt:          &bdr.CreatedAt,
	}, nil
}

func (r *Resolver) fetchResults(ctx context.Context, provider models.Provider, selector model.PageFilter, profileID string) (*model.PerfPageResult, error) {
	if profileID == "" {
		return nil, handlers.ErrQueryGet("*profileID")
	}

	tokenString := ctx.Value(models.TokenCtxKey).(string)

	bdr, err := provider.FetchResults(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order, profileID)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	result := &model.PerfPageResult{}

	if err := json.Unmarshal(bdr, result); err != nil {
		obj := "result data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return result, nil
}

// func (r *Resolver) listenToPerformanceResult(ctx context.Context, provider models.Provider, profileID string) (<-chan *model.MesheryResult, error) {
// 	// fmt.Println("SUBSCRIPTION STARTED!")
// 	// fmt.Println(r.Config.PerformanceChannels)
// 	if r.Config.PerformanceChannels == nil {
// 		r.Config.PerformanceChannels = make(map[string](chan *model.MesheryResult))
// 	}
// 	if r.Config.PerformanceChannels[profileID] == nil {
// 		r.Config.PerformanceChannels[profileID] = make(chan *model.MesheryResult)
// 	}
// 	// fmt.Println(r.Config.PerformanceChannels)
// 	go func() {
// 		r.Log.Info("Performance subscription started")

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				r.Log.Info("Performance subscription stopped")
// 				// delete channel from channels
// 				delete(r.Config.PerformanceChannels, profileID)
// 				return
// 			}
// 		}
// 	}()
// 	return r.Config.PerformanceChannels[profileID], nil
// }

func (r *Resolver) getPerformanceProfiles(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.PerfPageProfiles, error) {
	tokenString := ctx.Value("token").(string)

	bdr, err := provider.GetPerformanceProfiles(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	profiles := &model.PerfPageProfiles{}

	if err := json.Unmarshal(bdr, profiles); err != nil {
		obj := "performance profiles data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return profiles, nil
}
