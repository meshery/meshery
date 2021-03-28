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
	namespace       = "meshery"
	operatorSubject = "meshery.>"
	operatorQueue   = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
)

type controller struct {
	id         string
	name       string
	version    string
	kind       string
	apiversion string
	namespace  string
}

func (r *Resolver) changeOperatorStatus(ctx context.Context, status model.Status) (model.Status, error) {
	delete := true
	if status == model.StatusEnabled {
		r.Log.Info("Installing Operator")
		delete = false
	}

	if r.KubeClient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		return model.StatusUnknown, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		err := r.cleanEntries(del)
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

		err = initialize(kubeclient, del)
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
		r.Log.Info("Operator operation executed")

		if !del {
			endpoint, err := r.subscribeToBroker(kubeclient, r.brokerChannel)
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
			r.Log.Info("Connected to broker at:", endpoint)
		}

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
		r.Log.Info("Meshsync operation executed")

		r.operatorChannel <- &model.OperatorStatus{
			Status: status,
		}
	}(delete, r.KubeClient)

	r.operatorChannel <- &model.OperatorStatus{
		Status: model.StatusProcessing,
	}

	return model.StatusProcessing, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context) (*model.OperatorStatus, error) {
	status := model.StatusUnknown
	version := string(model.StatusUnknown)

	obj, err := getOperator(r.DBHandler)
	if err != nil {
		r.Log.Error(err)
		return &model.OperatorStatus{
			Status: status,
			Error: &model.Error{
				Code:        "",
				Description: err.Error(),
			},
		}, nil
	}

	if len(obj) > 0 {
		status = model.StatusEnabled
		version = obj[0].version
	} else {
		status = model.StatusDisabled
	}

	objs, err := getControllersInfo(r.KubeClient)
	if err != nil {
		r.Log.Error(err)
		return &model.OperatorStatus{
			Status: status,
			Error: &model.Error{
				Code:        "",
				Description: err.Error(),
			},
		}, nil
	}

	controllers := make([]*model.OperatorControllerStatus, 0)
	for _, obj := range objs {
		controllers = append(controllers, &model.OperatorControllerStatus{
			Name:    obj.name,
			Version: obj.version,
			Status:  model.StatusEnabled,
		})
	}

	return &model.OperatorStatus{
		Status:      status,
		Version:     version,
		Controllers: controllers,
	}, nil
}

func (r *Resolver) listenToOperatorState(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	if r.operatorChannel == nil {
		r.operatorChannel = make(chan *model.OperatorStatus)
	}

	go func() {
		r.Log.Info("Operator subscription started")
		err := r.connectToBroker(context.TODO())
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		select {
		case <-r.MeshSyncChannel:
			status, err := r.getOperatorStatus(ctx)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				return
			}

			if status.Status != model.StatusEnabled {
				_, err = r.changeOperatorStatus(ctx, model.StatusEnabled)
				if err != nil {
					r.Log.Error(ErrOperatorSubscription(err))
					return
				}
			}
			r.operatorChannel <- status
		}
	}()

	return r.operatorChannel, nil
}

func (r *Resolver) subscribeToBroker(mesheryKubeClient *mesherykube.Client, datach chan *broker.Message) (string, error) {
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
	r.brokerConn, err = nats.New(nats.Options{
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
			r.brokerConn, er = nats.New(nats.Options{
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

	err = r.brokerConn.SubscribeWithChannel(operatorSubject, operatorQueue, datach)
	if err != nil {
		return endpoint, ErrSubscribeChannel(err)
	}

	return endpoint, nil
}

func getOperator(handler *database.Handler) ([]*controller, error) {
	objects := make([]meshsyncmodel.Object, 0)

	subquery1 := handler.Select("id").Where("kind = ? AND key = ? AND value = ?", meshsyncmodel.KindLabel, "app", "meshery").Table("key_values")
	subquery2 := handler.Select("id").Where("id IN (?) AND kind = ? AND key = ? AND value = ?", subquery1, meshsyncmodel.KindLabel, "component", "operator").Table("key_values")
	result := handler.
		Preload("ObjectMeta").
		Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
		Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
		Preload("Spec").
		Preload("Status").
		Find(&objects, "id IN (?) AND kind = ?", subquery2, "Deployment")
	if result.Error != nil {
		return nil, ErrQuery(result.Error)
	}

	deploys := make([]*controller, 0)
	for _, obj := range objects {
		if meshsyncmodel.IsObject(obj) {
			version := "latest"
			for _, label := range obj.ObjectMeta.Labels {
				if label.Key == "version" {
					version = label.Value
				}
			}

			deploys = append(deploys, &controller{
				id:         obj.ID,
				version:    version,
				name:       obj.ObjectMeta.Name,
				kind:       obj.Kind,
				apiversion: obj.APIVersion,
				namespace:  obj.ObjectMeta.Namespace,
			})
		}
	}

	return deploys, nil
}

func getControllersInfo(mesheryKubeClient *mesherykube.Client) ([]*controller, error) {
	controllers := make([]*controller, 0)
	var broker *operatorv1alpha1.Broker
	var meshsync *operatorv1alpha1.MeshSync
	mesheryclient, err := client.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return controllers, ErrMesheryClient(nil)
		}
		return controllers, ErrMesheryClient(err)
	}

	broker, err = mesheryclient.CoreV1Alpha1().Brokers(namespace).Get(context.TODO(), "meshery-broker", metav1.GetOptions{})
	if err != nil {
		return controllers, ErrMesheryClient(err)
	}
	controllers = append(controllers, &controller{
		name:    "broker",
		version: broker.Labels["version"],
	})

	meshsync, err = mesheryclient.CoreV1Alpha1().MeshSyncs(namespace).Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	if err != nil {
		return controllers, ErrMesheryClient(err)
	}
	controllers = append(controllers, &controller{
		name:    "meshsync",
		version: meshsync.Labels["version"],
	})
	return controllers, nil
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

func (r *Resolver) cleanEntries(del bool) error {
	if del {
		r.Log.Info("Removing Operator")
		objs, err := getOperator(r.DBHandler)
		if err != nil {
			return err
		}

		for _, obj := range objs {
			err := recordMeshSyncData(broker.Delete, r.DBHandler, &meshsyncmodel.Object{
				ID:         obj.id,
				Kind:       obj.kind,
				APIVersion: obj.apiversion,
				ObjectMeta: &meshsyncmodel.ResourceObjectMeta{
					Name:      obj.name,
					Namespace: obj.namespace,
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
