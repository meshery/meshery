package resolver

import (
	"context"
	"fmt"
	"strings"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	natspackage "github.com/nats-io/nats.go"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
)

func (r *Resolver) changeOperatorStatus(ctx context.Context, status *model.Status) (*model.Status, error) {
	delete := true
	st := model.StatusProcessing
	if *status == model.StatusEnabled {
		delete = false
	}

	if r.KubeClient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		return nil, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		err := initialize(kubeclient, del)
		if err != nil {
			r.Log.Error(err)
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Initialized")

		endpoint, err := subscribeToBroker(kubeclient, r.meshsyncChannel)
		if err != nil {
			r.Log.Error(err)
			r.Log.Info(endpoint)
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Connected to broker at:", endpoint)

		// installMeshsync
		err = runMeshSync(kubeclient, del)
		if err != nil {
			r.Log.Error(err)
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Meshsync started")

		r.operatorChannel <- &model.OperatorStatus{
			Status: status,
		}
	}(delete, r.KubeClient)

	r.operatorChannel <- &model.OperatorStatus{
		Status: &st,
	}

	return &st, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context) (*model.OperatorStatus, error) {
	status := model.StatusUnknown

	obj, err := getOperator(r.DBHandler)
	if err != nil {
		r.Log.Error(err)
		return &model.OperatorStatus{
			Status: &status,
			Error: &model.Error{
				Code:        "",
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

func (r *Resolver) listenToOperatorState(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	r.operatorChannel = make(chan *model.OperatorStatus)

	go func() {
		r.Log.Info("Operator subscription started")
		endpoint, err := subscribeToBroker(r.KubeClient, r.meshsyncChannel)
		if err != nil {
			r.Log.Error(ErrOperatorSubscription(err))
			disabledStatus := model.StatusDisabled
			r.operatorChannel <- &model.OperatorStatus{
				Status: &disabledStatus,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Connected to broker at:", endpoint)

		select {
		case <-r.meshsyncChannel:
			status, err := r.getOperatorStatus(ctx)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				return
			}
			r.Log.Info("Operator status updated")
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
		Preload("ObjectMeta").
		Preload("ObjectMeta.Labels").
		Preload("ObjectMeta.Annotations", "id IN (?)", subquery2).
		Preload("Spec").
		Preload("Status").
		Find(&objects, "kind = ?", "Deployment")
	if result.Error != nil {
		return nil, ErrQuery(result.Error)
	}

	deploys := []string{}
	for _, obj := range objects {
		if meshsyncmodel.IsObject(obj) {
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

func subscribeToBroker(mesheryKubeClient *mesherykube.Client, datach chan *broker.Message) (string, error) {
	var broker *operatorv1alpha1.Broker
	mesheryclient, err := client.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return "", ErrMesheryClient(nil)
		}
		return "", ErrMesheryClient(err)
	}

	timeout := 60
	for timeout > 0 {
		broker, err = mesheryclient.CoreV1Alpha1().Brokers(namespace).Get(context.Background(), "meshery-broker", metav1.GetOptions{})
		if err == nil && broker.Status.Endpoint.External != "" {
			break
		}
		timeout--
		time.Sleep(1 * time.Second)
	}

	// subscribing to nats
	endpoint := broker.Status.Endpoint.External
	natsClient, err := nats.New(nats.Options{
		URLS:           []string{endpoint},
		ConnectionName: "meshery",
		Username:       "",
		Password:       "",
		ReconnectWait:  2 * time.Second,
		MaxReconnect:   5,
	})
	// Hack for minikube based clusters
	if err != nil {
		if err.Error() == nats.ErrConnect(natspackage.ErrNoServers).Error() {
			var er error
			var port, address string
			if len(strings.Split(broker.Status.Endpoint.External, ":")) > 1 {
				port = strings.Split(broker.Status.Endpoint.External, ":")[1]
			}
			if len(strings.SplitAfter(mesheryKubeClient.RestConfig.Host, "://")) > 1 {
				address = strings.SplitAfter(strings.SplitAfter(mesheryKubeClient.RestConfig.Host, "://")[1], ":")[0]
				if len(address) > 0 {
					address = address[:len(address)-1]
				}
			}
			endpoint = fmt.Sprintf("%s:%s", address, port)
			natsClient, er = nats.New(nats.Options{
				URLS:           []string{endpoint},
				ConnectionName: "meshery",
				Username:       "",
				Password:       "",
				ReconnectWait:  2 * time.Second,
				MaxReconnect:   5,
			})
			if er != nil {
				return endpoint, er
			}
		} else {
			return endpoint, err
		}
	}

	err = natsClient.SubscribeWithChannel(meshsyncSubject, meshsyncQueue, datach)
	if err != nil {
		return endpoint, err
	}

	return endpoint, nil
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
