package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/internal/graphql/model"
	graphqlModels "github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

// func (r *Resolver) subscribePerfResults(ctx context.Context, provider models.Provider, filter *graphqlModels.PageFilter) (<-chan *graphqlModels.PerfPageResult, error) {
// 	if r.performanceChannel == nil {
// 		r.performanceChannel = make(chan *graphqlModels.PerfPageResult)
// 		r.operatorSyncChannel = make(chan struct{})
// 	}

// 	go func() {
// 		r.Log.Info("Performance subscription started")

// 		tokenString := ctx.Value("token").(string)

// 		provider.FetchAllResults()

// 		for {
// 			select {
// 			case <-r.operatorSyncChannel:
// 				status, err := r.getOperatorStatus(ctx, provider)
// 				if err != nil {
// 					r.Log.Error(ErrOperatorSubscription(err))
// 					return
// 				}
// 				r.performanceChannel <- status
// 			case <-ctx.Done():
// 				r.Log.Info("Operator subscription flushed")
// 				return
// 			}
// 		}
// 	}()

// 	return r.performanceChannel, nil
// }

func (r *Resolver) getPerfResult(ctx context.Context, provider models.Provider, id string) (*graphqlModels.MesheryResult, error) {
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

	return &graphqlModels.MesheryResult{
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

func (r *Resolver) fetchResults(ctx context.Context, provider models.Provider, selector graphqlModels.PageFilter, profileID string) (*graphqlModels.PerfPageResult, error) {
	if profileID == "" {
		return nil, handlers.ErrQueryGet("*profileID")
	}

	tokenString := ctx.Value("token").(string)

	bdr, err := provider.FetchResults(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order, profileID)

	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	result := &graphqlModels.PerfPageResult{}

	if err := json.Unmarshal(bdr, result); err != nil {
		obj := "result data"
		return nil, handlers.ErrUnmarshal(err, obj)
	}

	return result, nil
}

func (r *Resolver) listenToPerformanceResult(ctx context.Context, provider models.Provider, profileID string) (<-chan *model.MesheryResult, error) {
	// fmt.Println("SUBSCRIPTION STARTED!")
	// fmt.Println(r.Config.PerformanceChannels)
	if r.Config.PerformanceChannels == nil {
		r.Config.PerformanceChannels = make(map[string](chan *model.MesheryResult))
	}
	if r.Config.PerformanceChannels[profileID] == nil {
		r.Config.PerformanceChannels[profileID] = make(chan *model.MesheryResult)
	}
	// fmt.Println(r.Config.PerformanceChannels)
	go func() {
		r.Log.Info("Performance subscription started")

		for {
			select {
			case <-ctx.Done():
				r.Log.Info("Performance subscription stopped")
				// delete channel from channels
				delete(r.Config.PerformanceChannels, profileID)
				return
			}
		}
	}()
	return r.Config.PerformanceChannels[profileID], nil
}
