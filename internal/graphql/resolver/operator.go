package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	namespace = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func (r *mutationResolver) changeOperatorStatus(ctx context.Context, status *model.Status) (*model.Status, error) {
	delete := true
	if *status == model.StatusEnabled {
		delete = false
	}

	go func(del bool) {
		er := initialize(r.KubeClient, del)
		if er != nil {
			return
		}
	}(delete)

	return status, nil
}

func (r *queryResolver) getOperatorStatus(ctx context.Context) (*model.Status, error) {
	return nil, nil
}

func (r *subscriptionResolver) listenToOperatorEvents(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	return nil, nil
}

func (r *subscriptionResolver) subscribeToMeshSync(ctx context.Context) (<-chan *model.OperatorControllerStatus, error) {
	channel := make(chan *model.OperatorControllerStatus)
	return channel, nil
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

// func subscribeToBroker(client *mesherykube.Client, datach chan *broker.Message) error {
// 	var broker *operatorv1alpha1.Broker
// 	mesheryclient, err := client.New(&client.RestConfig)
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		broker, err = mesheryclient.CoreV1Alpha1().Brokers(namespace).Get(context.Background(), "meshery-broker", metav1.GetOptions{})
// 		if err == nil && broker.Status.Endpoint.External != "" {
// 			break
// 		}
// 		time.Sleep(1 * time.Second)
// 	}

// 	// subscribing to nats
// 	natsClient, err := nats.New(broker.Status.Endpoint.External)
// 	if err != nil {
// 		return err
// 	}

// 	err = natsClient.SubscribeWithChannel("meshsync", namespace, datach)
// 	if err != nil {
// 		return err
// 	}

// 	status.SubscriptionStarted = "true"
// 	return nil
// }

// func runMeshSync(delete bool) error {
// 	err := applyYaml(delete, meshsyncYaml)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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
