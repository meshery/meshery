package resolver

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace       = "meshery"
	operatorSubject = "meshery.meshsync.core"
	brokerQueue     = "meshery"

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
		r.Log.Info("Operator operation executed")

		if !del {
			endpoint, err := r.subscribeToBroker(kubeclient, r.brokerChannel)
			r.Log.Debug("Endpoint: ", endpoint)
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
	if r.KubeClient == nil {
		return nil, ErrMesheryClient(nil)
	}

	obj, err := getOperator(r.KubeClient)
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
		r.operatorSyncChannel = make(chan struct{})
	}

	go func() {
		r.Log.Info("Operator subscription started")
		err := r.connectToBroker(context.TODO())
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		// Enforce enable operator
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

		for {
			select {
			case <-r.operatorSyncChannel:
				status, err := r.getOperatorStatus(ctx)
				if err != nil {
					r.Log.Error(ErrOperatorSubscription(err))
					return
				}
				r.operatorChannel <- status
			case <-ctx.Done():
				r.Log.Info("Operator subscription flushed")
				return
			}
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

	endpoint := broker.Status.Endpoint.Internal
	if len(strings.Split(broker.Status.Endpoint.Internal, ":")) > 1 {
		port, _ := strconv.Atoi(strings.Split(broker.Status.Endpoint.Internal, ":")[1])
		if !utils.TcpCheck(&utils.HostPort{
			Address: strings.Split(broker.Status.Endpoint.Internal, ":")[0],
			Port:    int32(port),
		}, nil) {
			endpoint = broker.Status.Endpoint.External
			port, _ = strconv.Atoi(strings.Split(broker.Status.Endpoint.External, ":")[1])
			if !utils.TcpCheck(&utils.HostPort{
				Address: strings.Split(broker.Status.Endpoint.External, ":")[0],
				Port:    int32(port),
			}, nil) {
				if !utils.TcpCheck(&utils.HostPort{
					Address: "host.docker.internal",
					Port:    int32(port),
				}, nil) {
					u, _ := url.Parse(r.KubeClient.RestConfig.Host)
					if utils.TcpCheck(&utils.HostPort{
						Address: u.Hostname(),
						Port:    int32(port),
					}, nil) {
						endpoint = fmt.Sprintf("%s:%d", u.Hostname(), int32(port))
					}
				} else {
					endpoint = fmt.Sprintf("host.docker.internal:%d", int32(port))
				}
			}
		}
	}

	// subscribing to nats
	r.BrokerConn, err = nats.New(nats.Options{
		URLS:           []string{endpoint},
		ConnectionName: "meshery",
		Username:       "",
		Password:       "",
		ReconnectWait:  2 * time.Second,
		MaxReconnect:   5,
	})
	// Hack for minikube based clusters
	if err != nil {
		return endpoint, err
	}

	err = r.BrokerConn.SubscribeWithChannel(operatorSubject, brokerQueue, datach)
	if err != nil {
		return endpoint, ErrSubscribeChannel(err)
	}

	return endpoint, nil
}

func getOperator(kubeclient *mesherykube.Client) ([]*controller, error) {
	if kubeclient == nil || kubeclient.KubeClient == nil || kubeclient.KubeClient.AppsV1() == nil {
		return nil, ErrMesheryClient(nil)
	}

	dep, err := kubeclient.KubeClient.AppsV1().Deployments("meshery").Get(context.TODO(), "meshery-operator", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return nil, ErrMesheryClient(err)
	}

	deploys := make([]*controller, 0)
	version := ""
	for _, container := range dep.Spec.Template.Spec.Containers {
		if container.Name == "manager" {
			version = strings.Split(container.Image, ":")[1]
		}
	}

	deploys = append(deploys, &controller{
		version:    version,
		name:       dep.ObjectMeta.Name,
		kind:       dep.Kind,
		apiversion: dep.APIVersion,
		namespace:  dep.ObjectMeta.Namespace,
	})

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
	if err != nil && !kubeerror.IsNotFound(err) {
		return controllers, ErrMesheryClient(err)
	}
	controllers = append(controllers, &controller{
		name:    "broker",
		version: broker.Labels["version"],
	})

	meshsync, err = mesheryclient.CoreV1Alpha1().MeshSyncs(namespace).Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
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
