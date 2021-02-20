package resolver

import (
	"context"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
	"github.com/layer5io/meshsync/pkg/broker/nats"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
)

func (r *mutationResolver) changeOperatorStatus(ctx context.Context, status *model.Status) (*model.Status, error) {
	delete := true
	if *status == model.StatusEnabled {
		delete = false
	}

	go func(del bool) {
		er := initialize(r.KubeClient, del)
		if er != nil {
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        errCode,
					Description: er.Error(),
				},
			}
			return
		}

		r.meshsyncChannel = make(chan *broker.Message)
		err := subscribeToBroker(r.KubeClient, r.meshsyncChannel)
		if err != nil {
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        errCode,
					Description: err.Error(),
				},
			}
			return
		}

		// installMeshsync
		err = runMeshSync(r.KubeClient, del)
		if err != nil {
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        errCode,
					Description: err.Error(),
				},
			}
			return
		}
	}(delete)

	r.operatorChannel <- &model.OperatorStatus{
		Status: status,
	}

	return status, nil
}

func (r *queryResolver) getOperatorStatus(ctx context.Context) (*model.OperatorStatus, error) {
	return nil, nil
}

func (r *subscriptionResolver) listenToOperatorEvents(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	r.operatorChannel = make(chan *model.OperatorStatus)
	return r.operatorChannel, nil
}

func initialize(client *mesherykube.Client, delete bool) error {
	// installOperator
	err := applyYaml(client, delete, operatorYaml)
	if err != nil {
		return err
	}

	// installBroker
	err = applyYaml(client, delete, brokerYaml)
	if err != nil {
		return err
	}

	return nil
}

func subscribeToBroker(mesheryKubeClient *mesherykube.Client, datach chan *broker.Message) error {
	var broker *operatorv1alpha1.Broker
	mesheryclient, err := client.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		return err
	}

	for {
		broker, err = mesheryclient.CoreV1Alpha1().Brokers(namespace).Get(context.Background(), "meshery-broker", metav1.GetOptions{})
		if err == nil && broker.Status.Endpoint.External != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// subscribing to nats
	natsClient, err := nats.New(broker.Status.Endpoint.External)
	if err != nil {
		return err
	}

	err = natsClient.SubscribeWithChannel(meshsyncName, namespace, datach)
	if err != nil {
		return err
	}

	return nil
}

func applyYaml(client *mesherykube.Client, delete bool, file string) error {
	contents, err := utils.ReadRemoteFile(file)
	if err != nil {
		return err
	}

	err = client.ApplyManifest([]byte(contents), mesherykube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    delete,
	})
	if err != nil {
		return err
	}

	return nil
}
