package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils/broadcast"
)

func (r *Resolver) getDataPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.DataPlane, error) {
	selectors := make([]model.MeshType, 0)
	if filter.Type == nil {
		for _, mesh := range model.AllMeshType {
			selectors = append(selectors, mesh)
		}
	} else {
		selectors = append(selectors, *filter.Type)
	}

	dataPlaneList, err := model.GetDataPlaneState(selectors, provider)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return dataPlaneList, nil
}

func (r *Resolver) listenToDataPlaneState(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) (<-chan []*model.DataPlane, error) {
	dataPlaneChannel := make(chan []*model.DataPlane)

	broadcastChannel := make(chan broadcast.BroadcastMessage)
	r.Broadcast.Register(broadcastChannel)

	go func() {
		r.Log.Info("Initializing DataPlane subscription")
		err := r.connectToBroker(context.TODO(), provider)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		for {
			select {
			case message := <-broadcastChannel:
				if message.Type == "meshsync" {
					r.Log.Info("Dataplane sync channel called")
					containers, err := r.getDataPlanes(ctx, provider, filter)
					if err != nil {
						r.Log.Error(ErrDataPlaneSubscription(err))
						break
					}

					dataPlaneChannel <- containers
				}
			case <-ctx.Done():
				r.Log.Info("DataPlane subscription stopped")
				close(dataPlaneChannel)
				r.Broadcast.Unregister((broadcastChannel))
				close(broadcastChannel)
				return
			}
		}
	}()
	return dataPlaneChannel, nil
}
