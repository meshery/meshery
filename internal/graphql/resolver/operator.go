package resolver

import (
	"context"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
	"github.com/layer5io/meshsync/pkg/broker/nats"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
)

func (r *Resolver) changeOperatorStatus(ctx context.Context, status *model.Status) (*model.Status, error) {
	delete := true
	if *status == model.StatusEnabled {
		delete = false
	}

	go func(del bool) {
		kubeClient, err := r.GetKubeClient()
		if err != nil {
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        errCode,
					Description: "failed to create meshkube client",
				},
			}
			return
		}

		er := initialize(kubeClient, del)
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
		err = subscribeToBroker(kubeClient, r.meshsyncChannel)
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
		err = runMeshSync(kubeClient, del)
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

func (r *Resolver) getOperatorStatus(ctx context.Context) (*model.OperatorStatus, error) {
	status := model.StatusUnknown

	obj, err := getOperator(r.DBHandler)
	if err != nil {
		return &model.OperatorStatus{
			Status: &status,
			Error: &model.Error{
				Code:        errCode,
				Description: err.Error(),
			},
		}, nil
	}

	if len(obj) > 0 {
		status = model.StatusEnabled
	} else {
		status = model.StatusDisabled
	}

	return &model.OperatorStatus{
		Status: &status,
	}, nil
}

func (r *Resolver) listenToOperatorEvents(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	r.operatorChannel = make(chan *model.OperatorStatus)

	go func() {
		select {
		case <-r.meshsyncChannel:
			status, err := r.getOperatorStatus(ctx)
			if err != nil {
				return
			}
			r.operatorChannel <- status
		}
	}()

	return r.operatorChannel, nil
}

func getOperator(handler *database.Handler) ([]string, error) {
	objects := make([]meshsyncmodel.Object, 0)

	subquery1 := handler.Select("id").Where("key = ? AND value = ?", "app", "meshery").Table("key_values")
	subquery2 := handler.Select("id").Where("id IN (?) AND key = ? AND value = ?", subquery1, "component", "operator").Table("key_values")
	result := handler.
		Preload("TypeMeta", "kind = ?", "Deployment").
		Preload("ObjectMeta").
		Preload("ObjectMeta.Labels").
		Preload("ObjectMeta.Annotations").
		Preload("Spec").
		Preload("Status").
		Find(&objects, "id IN (?)", subquery2)
	if result.Error != nil {
		return nil, result.Error
	}

	deploys := []string{}
	for _, obj := range objects {
		if obj.TypeMeta != nil && obj.ObjectMeta != nil && obj.Spec != nil && obj.Status != nil {
			deploys = append(deploys, obj.ObjectMeta.Name)
		}
	}

	return deploys, nil
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
