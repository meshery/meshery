package resolver

import (
	"context"
	"time"

	operatorClient "github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/utils/broadcast"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	MeshSyncSubscriptionError = model.Error{
		Description: "Failed to get MeshSync data",
		Code:        ErrMeshsyncSubscriptionCode,
	}
	MeshSyncMesheryClientMissingError = model.Error{
		Code:        ErrMeshsyncSubscriptionCode,
		Description: "Cannot find Meshery Client",
	}
)

func (r *Resolver) listenToMeshSyncEvents(ctx context.Context, provider models.Provider) (<-chan *model.OperatorControllerStatus, error) {
	channel := make(chan *model.OperatorControllerStatus)
	if r.brokerChannel == nil {
		r.brokerChannel = make(chan *broker.Message)
	}
	prevStatus := r.getMeshSyncStatus()

	go func(ch chan *model.OperatorControllerStatus) {
		r.Log.Info("Initializing MeshSync subscription")

		go model.ListernToEvents(r.Log, provider.GetGenericPersister(), r.brokerChannel, r.MeshSyncChannel, r.operatorSyncChannel, r.controlPlaneSyncChannel, r.meshsyncLivenessChannel, r.Broadcast)

		// signal to install operator when initialized
		r.MeshSyncChannel <- struct{}{}
		// extension to notify other channel when data comes in
		for {
			status := r.getMeshSyncStatus()

			ch <- &status
			if status != prevStatus {
				prevStatus = status
			}

			time.Sleep(10 * time.Second)
		}
	}(channel)

	return channel, nil
}

func (r *Resolver) getMeshSyncStatus() model.OperatorControllerStatus {
	var status model.OperatorControllerStatus
	mesheryclient, err := operatorClient.New(&r.Config.KubeClient.RestConfig)

	if err != nil {
		return model.OperatorControllerStatus{
			Name:    "",
			Version: "",
			Status:  model.StatusDisabled,
			Error:   &MeshSyncMesheryClientMissingError,
		}
	}

	status, err = model.GetMeshSyncInfo(mesheryclient, r.meshsyncLivenessChannel)

	if err != nil {
		return model.OperatorControllerStatus{
			Name:    "",
			Version: "",
			Status:  model.StatusDisabled,
			Error:   &MeshSyncSubscriptionError,
		}
	}

	return status
}

func (r *Resolver) resyncCluster(ctx context.Context, provider models.Provider, actions *model.ReSyncActions) (model.Status, error) {
	if actions.ClearDb == "true" {
		// Clear existing data
		err := provider.GetGenericPersister().Migrator().DropTable(
			&meshsyncmodel.KeyValue{},
			&meshsyncmodel.Object{},
			&meshsyncmodel.ResourceSpec{},
			&meshsyncmodel.ResourceStatus{},
			&meshsyncmodel.ResourceObjectMeta{},
		)
		if err != nil {
			if provider.GetGenericPersister() == nil {
				return "", ErrEmptyHandler
			}
			r.Log.Warn(ErrDeleteData(err))
		}
		err = provider.GetGenericPersister().Migrator().CreateTable(
			&meshsyncmodel.KeyValue{},
			&meshsyncmodel.Object{},
			&meshsyncmodel.ResourceSpec{},
			&meshsyncmodel.ResourceStatus{},
			&meshsyncmodel.ResourceObjectMeta{},
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
	status, err := r.getOperatorStatus(ctx, provider)
	if err != nil {
		return err
	}

	if r.BrokerConn.IsEmpty() && status != nil && status.Status == model.StatusEnabled {
		endpoint, err := model.SubscribeToBroker(provider, r.Config.KubeClient, r.brokerChannel, r.BrokerConn)
		if err != nil {
			r.Log.Error(ErrAddonSubscription(err))

			r.Broadcast.Submit(broadcast.BroadcastMessage{
				Source: broadcast.OperatorSyncChannel,
				Type:   "error",
				Data:   err,
			})

			return err
		}
		r.Log.Info("Connected to broker at:", endpoint)
		r.Config.BrokerEndpointURL = &endpoint
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data:   false,
			Type:   "health",
		})
		return nil
	}

	if r.BrokerConn.Info() == broker.NotConnected {
		return ErrBrokerNotConnected
	}

	return nil
}

func (r *Resolver) deployMeshsync(ctx context.Context, provider models.Provider) (model.Status, error) {
	//err := model.RunMeshSync(r.Config.KubeClient, false)
	r.Log.Info("Installing Meshsync")
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data:   true,
		Type:   "health",
	})

	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data:   false,
		Type:   "health",
	})
	return model.StatusProcessing, nil
}

func (r *Resolver) connectToNats(ctx context.Context, provider models.Provider) (model.Status, error) {
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data:   true,
		Type:   "health",
	})
	err := r.connectToBroker(ctx, provider)
	if err != nil {
		r.Log.Error(err)
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data:   err,
			Type:   "error",
		})
		return model.StatusDisabled, err
	}

	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data:   false,
		Type:   "health",
	})
	return model.StatusConnected, nil
}
