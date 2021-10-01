package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils/broadcast"
)

func (r *Resolver) getControlPlanes(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) ([]*model.ControlPlane, error) {
	selectors := make([]model.MeshType, 0)
	if filter.Type == nil {
		for _, mesh := range model.AllMeshType {
			selectors = append(selectors, mesh)
		}
	} else {
		selectors = append(selectors, *filter.Type)
	}

	controlplanelist, err := model.GetControlPlaneState(selectors, provider)
	if err != nil {
		r.Log.Error(err)
		return nil, err
	}

	return controlplanelist, nil
}

func (r *Resolver) listenToControlPlaneState(ctx context.Context, provider models.Provider, filter *model.ServiceMeshFilter) (<-chan []*model.ControlPlane, error) {
	controlPlaneChannel := make(chan []*model.ControlPlane)

	broadcastChannel := make(chan broadcast.BroadcastMessage)
	r.Broadcast.Register(broadcastChannel)

	go func() {
		r.Log.Info("Initializing ControlPlane subscription")
		err := r.connectToBroker(context.TODO(), provider)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		for {
			select {
			case message := <-broadcastChannel:
				if message.Type == "meshsync" {
					status, err := r.getControlPlanes(ctx, provider, filter)
					if err != nil {
						r.Log.Error(ErrControlPlaneSubscription(err))
						break
					}
					controlPlaneChannel <- status
				}
			case <-ctx.Done():
				r.Log.Info("ControlPlane subscription stopped")
				close(controlPlaneChannel)
				r.Broadcast.Unregister(broadcastChannel)
				close(broadcastChannel)
				return
			}
		}
	}()
	return controlPlaneChannel, nil
}
