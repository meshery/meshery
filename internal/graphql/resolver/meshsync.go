package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	"gorm.io/gorm"
)

var (
	meshsyncName = "meshsync"
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func (r *Resolver) listenToMeshSyncEvents(ctx context.Context) (<-chan *model.OperatorControllerStatus, error) {
	channel := make(chan *model.OperatorControllerStatus)
	if r.brokerChannel == nil {
		r.brokerChannel = make(chan *broker.Message)
	}

	go func(ch chan *model.OperatorControllerStatus) {
		r.Log.Info("MeshSync subscription started")
		go listernToEvents(r.Log, r.DBHandler, r.brokerChannel, r.MeshSyncChannel)

		// signal to install operator when initialized
		r.MeshSyncChannel <- struct{}{}
		// extension to notify other channel when data comes in
	}(channel)

	return channel, nil
}

func (r *Resolver) connectToBroker(ctx context.Context) error {
	status, err := r.getOperatorStatus(ctx)
	if err != nil {
		return err
	}
	if r.brokerConn == nil && status != nil && status.Status == model.StatusEnabled {
		endpoint, err := r.subscribeToBroker(r.KubeClient, r.brokerChannel)
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
	return nil
}
