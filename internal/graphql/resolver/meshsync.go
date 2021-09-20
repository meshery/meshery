package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func (r *Resolver) listenToMeshSyncEvents(ctx context.Context, provider models.Provider) (<-chan *model.OperatorControllerStatus, error) {
	channel := make(chan *model.OperatorControllerStatus)
	if r.brokerChannel == nil {
		r.brokerChannel = make(chan *broker.Message)
	}

	go func(ch chan *model.OperatorControllerStatus) {
		r.Log.Info("Initializing MeshSync subscription")
		go model.ListernToEvents(r.Log, provider.GetGenericPersister(), r.brokerChannel, r.MeshSyncChannel, r.operatorSyncChannel, r.controlPlaneSyncChannel, r.meshsyncLivenessChannel, r.Broadcast)

		// signal to install operator when initialized
		r.MeshSyncChannel <- struct{}{}
		// extension to notify other channel when data comes in
	}(channel)

	return channel, nil
}

func (r *Resolver) resyncCluster(ctx context.Context, provider models.Provider, actions *model.ReSyncActions) (model.Status, error) {
	if actions.ClearDb == "true" {
		// Clear existing data
		err := provider.GetGenericPersister().Migrator().DropTable(
			meshsyncmodel.KeyValue{},
			meshsyncmodel.Object{},
		)
		if err != nil {
			if provider.GetGenericPersister() == nil {
				return "", ErrEmptyHandler
			}
			r.Log.Warn(ErrDeleteData(err))
		}
	}
	if actions.ReSync == "true" {
		err := r.BrokerConn.Publish(model.RequestSubject, &broker.Message{
			Request: &broker.RequestObject{
				Entity: broker.ReSyncDiscoveryEntity,
			},
		})
		if err != nil {
			return "", ErrPublishBroker(err)
		}
	}
	return model.StatusProcessing, nil
}

func (r *Resolver) connectToBroker(ctx context.Context, provider models.Provider) error {
	kubeclient, ok := ctx.Value(models.KubeHanderKey).(*mesherykube.Client)
	if !ok {
		r.Log.Error(ErrNilClient)
		return ErrNilClient
	}

	status, err := r.getOperatorStatus(ctx, provider)
	if err != nil {
		return err
	}

	if r.BrokerConn.IsEmpty() && status != nil && status.Status == model.StatusEnabled {
		endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn)
		if err != nil {
			r.Log.Error(ErrAddonSubscription(err))

			r.Broadcast.Submit(broadcast.BroadcastMessage{
				Type:    broadcast.OperatorSyncChannel,
				Message: err,
			})
			return err
		}
		r.Log.Info("Connected to broker at:", endpoint)

		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Type:    broadcast.OperatorSyncChannel,
			Message: false,
		})
		return nil
	}

	if r.BrokerConn.Info() == broker.NotConnected {
		return ErrBrokerNotConnected
	}

	return nil
}

func (r *Resolver) deployMeshsync(ctx context.Context, provider models.Provider) (model.Status, error) {
	kubeclient, ok := ctx.Value(models.KubeHanderKey).(*mesherykube.Client)
	if !ok {
		return model.StatusProcessing, ErrMesheryClient(nil)
	}

	err := model.RunMeshSync(kubeclient, false)
	r.Log.Info("Installing Meshsync")
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Type:    broadcast.OperatorSyncChannel,
		Message: true,
	})

	if err != nil {
		r.Log.Error(err)
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Type:    broadcast.OperatorSyncChannel,
			Message: err,
		})
		return model.StatusDisabled, err
	}

	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Type:    broadcast.OperatorSyncChannel,
		Message: false,
	})

	return model.StatusProcessing, nil
}

func (r *Resolver) connectToNats(ctx context.Context, provider models.Provider) (model.Status, error) {
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Type:    broadcast.OperatorSyncChannel,
		Message: true,
	})
	err := r.connectToBroker(ctx, provider)
	if err != nil {
		r.Log.Error(err)
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Type:    broadcast.OperatorSyncChannel,
			Message: err,
		})
		return model.StatusDisabled, err
	}

	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Type:    broadcast.OperatorSyncChannel,
		Message: false,
	})
	return model.StatusConnected, nil
}
