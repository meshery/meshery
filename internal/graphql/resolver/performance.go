package resolver

import (
	"context"
	"fmt"

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
	var meshery_result *models.MesheryResult
	if *id == "" {
		return nil, handlers.ErrQueryGet("*id")
	}

	res := provider.GetGenericPersister().
		First(&meshery_result)

	r.Log.Info("1")

	if res.Error != nil {
		fmt.Printf(*id)
		fmt.Printf(res.Error.Error())
		r.Log.Info("1.2")
		return nil, res.Error
	}

	r.Log.Info("2")
	// sp, err := meshery_result.ConvertToSpec()

	// if err != nil {
	// 	r.Log.Error(ErrQuery(res.Error))
	// 	return nil, ErrQuery(res.Error)

	// }

	// start_time := int(sp.StartTime.Unix())
	// end_time := int(sp.EndTime.Unix())

	// return &graphqlModels.PerfResult{
	// 	SmpVersion: &sp.SMPVersion,
	// 	// ID:         &sp.EnvID,
	// 	// labels:     &sp.,
	// 	StartTime: &start_time,
	// 	EndTime:   &end_time,
	// 	LatenciesMs: &graphqlModels.LatenciesMs{
	// 		Min:     &sp.Latencies.Min,
	// 		Average: &sp.Latencies.Min,
	// 		P50:     &sp.Latencies.P50,
	// 		P90:     &sp.Latencies.P90,
	// 		P99:     &sp.Latencies.P99,
	// 		Max:     &sp.Latencies.Max,
	// 	},
	// 	ActualQPS:    &sp.ActualQPS,
	// 	DetailsURI:   &sp.DetailsURI,
	// 	TestID:       &sp.TestID,
	// 	MeshConfigID: &sp.MeshConfigID,
	// 	EnvID:        &sp.EnvID,
	// }, nil

	r.Log.Info("3")
	start_time := int(meshery_result.TestStartTime.Unix())
	server_board_config := fmt.Sprintf("%v", meshery_result.ServerBoardConfig)
	server_metrics := fmt.Sprintf("%v", meshery_result.ServerMetrics)
	runner_results := fmt.Sprintf("%v", meshery_result.Result)
	meshery_id := fmt.Sprintf("%v", meshery_result.ID)
	performance_profile := fmt.Sprintf("%v", meshery_result.PerformanceProfileInfo.ID)
	r.Log.Info("4")

	return &graphqlModels.MesheryResult{
		MesheryID:          &meshery_id,
		Name:               &meshery_result.Name,
		Mesh:               &meshery_result.Mesh,
		PerformanceProfile: &performance_profile,
		TestID:             &meshery_result.TestID,
		RunnerResults:      &runner_results,
		ServerMetrics:      &server_metrics,
		ServerBoardConfig:  &server_board_config,
		TestStartTime:      &start_time,
		// UserID: &meshery_result,
		// UpdatedAt: &meshery_result.Name,
		// CreatedAt: &meshery_result.Name,
	}, nil

	// return meshery_result, nil
}
