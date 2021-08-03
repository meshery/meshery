package resolver

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/handlers"
	graphqlModels "github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

// func (r *Resolver) subscribePerfResults(ctx context.Context, provider models.Provider, filter *model.PageFilter) (<-chan *model.PerfPageResult, error) {
// 	if r.performanceChannel == nil {
// 		r.performanceChannel = make(chan *model.PerfPageResult)
// 		r.operatorSyncChannel = make(chan struct{})
// 	}

// 	go func() {
// 		r.Log.Info("Performance subscription started")
// 		// err := r.connectToBroker(context.TODO(), provider)
// 		// if err != nil && err != ErrNoMeshSync {
// 		// 	r.Log.Error(err)
// 		// 	return
// 		// }

// 		// Enforce enable operator
// 		status, err := r.getOperatorStatus(ctx, provider)
// 		if err != nil {
// 			r.Log.Error(ErrOperatorSubscription(err))
// 			return
// 		}
// 		if status.Status != model.StatusEnabled {
// 			_, err = r.changeOperatorStatus(ctx, provider, model.StatusEnabled)
// 			if err != nil {
// 				r.Log.Error(ErrOperatorSubscription(err))
// 				return
// 			}
// 		}

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

func (r *Resolver) getPerfResult(ctx context.Context, provider models.Provider, id *string) (*graphqlModels.MesheryResult, error) {
	if *id == "" {
		return nil, handlers.ErrQueryGet("*id")
	}

	resultID, err := uuid.FromString(*id)

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

	start_time := int(bdr.TestStartTime.Unix())
	server_board_config := fmt.Sprintf("%v", bdr.ServerBoardConfig)
	server_metrics := fmt.Sprintf("%v", bdr.ServerMetrics)
	runner_results := fmt.Sprintf("%v", bdr.Result)
	meshery_id := fmt.Sprintf("%v", bdr.ID)
	performance_profile := fmt.Sprintf("%v", bdr.PerformanceProfileInfo.ID)

	return &graphqlModels.MesheryResult{
		MesheryID:          &meshery_id,
		Name:               &bdr.Name,
		Mesh:               &bdr.Mesh,
		PerformanceProfile: &performance_profile,
		TestID:             &bdr.TestID,
		RunnerResults:      &runner_results,
		ServerMetrics:      &server_metrics,
		ServerBoardConfig:  &server_board_config,
		TestStartTime:      &start_time,
		// UserID: &meshery_result,
		// UpdatedAt: &meshery_result.Name,
		// CreatedAt: &meshery_result.Name,
	}, nil

}
