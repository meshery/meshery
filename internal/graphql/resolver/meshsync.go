package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	brokerpkg "github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	"gorm.io/gorm"
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
		r.Log.Info("MeshSync subscription started")
		go listernToEvents(r.Log, provider.GetGenericPersister(), r.brokerChannel, r.MeshSyncChannel, r.operatorSyncChannel, r.meshsyncLivenessChannel)

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
		err := r.BrokerConn.Publish(requestSubject, &brokerpkg.Message{
			Request: &brokerpkg.RequestObject{
				Entity: brokerpkg.ReSyncDiscoveryEntity,
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
		endpoint, err := r.subscribeToBroker(provider, r.Config.KubeClient, r.brokerChannel)
		if err != nil {
			r.Log.Error(ErrAddonSubscription(err))
			r.operatorChannel <- &model.OperatorStatus{
				Status: model.StatusDisabled,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return err
		}
		r.Log.Info("Connected to broker at:", endpoint)
		return nil
	}
	return ErrNoMeshSync
}

func runMeshSync(client *mesherykube.Client, delete bool) error {
	err := applyYaml(client, delete, meshsyncYaml)
	if err != nil {
		return err
	}
	return nil
}

func recordMeshSyncData(eventtype broker.EventType, handler *database.Handler, object *meshsyncmodel.Object) error {
	if handler == nil {
		return ErrEmptyHandler
	}

	handler.Lock()
	switch eventtype {
	case broker.Add, broker.Update:
		result := handler.Create(object)
		if result.Error != nil {
			result = handler.Session(&gorm.Session{FullSaveAssociations: true}).Updates(object)
			if result.Error != nil {
				return ErrCreateData(result.Error)
			}
		}
	case broker.Delete:
		result := handler.Delete(object)
		if result.Error != nil {
			return ErrDeleteData(result.Error)
		}
	case broker.ErrorEvent:
		return nil
	}
	handler.Unlock()
	return nil
}
