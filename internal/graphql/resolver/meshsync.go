package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

var (
	meshsyncName = "meshsync"
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func (r *subscriptionResolver) subscribeToMeshSync(ctx context.Context) (<-chan *model.OperatorControllerStatus, error) {
	channel := make(chan *model.OperatorControllerStatus)
	status := model.StatusUnknown

	go func(ch chan *model.OperatorControllerStatus) {
		err := listernToEvents(r.DBHandler, r.meshsyncChannel)
		if err != nil {
			ch <- &model.OperatorControllerStatus{
				Name:   &meshsyncName,
				Status: &status,
				Error: &model.Error{
					Code:        errCode,
					Description: err.Error(),
				},
			}
			return
		}
		// extension to notify other channel when data comes in
	}(channel)

	return channel, nil
}

func runMeshSync(client *mesherykube.Client, delete bool) error {
	err := applyYaml(client, delete, meshsyncYaml)
	if err != nil {
		return err
	}
	return nil
}

func recordMeshSyncData(handler *database.Handler, object meshsyncmodel.Object) error {
	result := handler.Create(&object)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func listernToEvents(handler *database.Handler, datach chan *broker.Message) error {
	for {
		select {
		case msg := <-datach:
			objectJSON, _ := utils.Marshal(msg.Object)
			object := meshsyncmodel.Object{}
			err := utils.Unmarshal(string(objectJSON), &object)
			if err != nil {
				return err
			}

			// persist the object
			err = recordMeshSyncData(handler, object)
			if err != nil {
				return err
			}
		}
	}
}
