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
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	brokerpkg "github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace       = "meshery"
	requestSubject  = "meshery.meshsync.request"
	meshsyncSubject = "meshery.meshsync.core"
	brokerQueue     = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
)

func (r *Resolver) changeOperatorStatus(ctx context.Context, provider models.Provider, status model.Status) (model.Status, error) {
	delete := true
	if status == model.StatusEnabled {
		r.Log.Info("Installing Operator")
		delete = false
	}

	if r.Config.KubeClient.KubeClient == nil {
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
			status, err := r.resyncCluster(context.TODO(), provider, &model.ReSyncActions{
				ReSync:  "false",
				ClearDb: "true",
			})
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

			endpoint, err := r.subscribeToBroker(provider, kubeclient, r.brokerChannel)
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
	}(delete, r.Config.KubeClient)

	r.operatorChannel <- &model.OperatorStatus{
		Status: model.StatusProcessing,
	}

	return model.StatusProcessing, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context, provider models.Provider) (*model.OperatorStatus, error) {
	status := model.StatusUnknown
	version := string(model.StatusUnknown)
	if r.Config.KubeClient == nil {
		return nil, ErrMesheryClient(nil)
	}

	name, version, err := getOperator(r.Config.KubeClient)
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
	if name == "" {
		status = model.StatusDisabled
	} else {
		status = model.StatusEnabled
	}

	controllers, err := getControllersInfo(r.Config.KubeClient, r.BrokerConn, r.meshsyncLivenessChannel)
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

	return &model.OperatorStatus{
		Status:      status,
		Version:     version,
		Controllers: controllers,
	}, nil
}

func (r *Resolver) listenToOperatorState(ctx context.Context, provider models.Provider) (<-chan *model.OperatorStatus, error) {
	if r.operatorChannel == nil {
		r.operatorChannel = make(chan *model.OperatorStatus)
	}
	if r.operatorSyncChannel == nil {
		r.operatorSyncChannel = make(chan struct{})
	}
	if r.meshsyncLivenessChannel == nil {
		r.meshsyncLivenessChannel = make(chan struct{})
	}

	go func() {
		r.Log.Info("Operator subscription started")
		err := r.connectToBroker(context.TODO(), provider)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		// Enforce enable operator
		status, err := r.getOperatorStatus(ctx, provider)
		if err != nil {
			r.Log.Error(ErrOperatorSubscription(err))
			return
		}
		if status.Status != model.StatusEnabled {
			_, err = r.changeOperatorStatus(ctx, provider, model.StatusEnabled)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				return
			}
		}

		for {
			select {
			case <-r.operatorSyncChannel:
				status, err := r.getOperatorStatus(ctx, provider)
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

func (r *Resolver) subscribeToBroker(provider models.Provider, mesheryKubeClient *mesherykube.Client, datach chan *brokerpkg.Message) (string, error) {
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
					u, _ := url.Parse(r.Config.KubeClient.RestConfig.Host)
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
	conn, err := nats.New(nats.Options{
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
	conn.DeepCopyInto(r.BrokerConn)

	err = r.BrokerConn.SubscribeWithChannel(meshsyncSubject, brokerQueue, datach)
	if err != nil {
		return endpoint, ErrSubscribeChannel(err)
	}

	err = r.BrokerConn.Publish(requestSubject, &brokerpkg.Message{
		Request: &brokerpkg.RequestObject{
			Entity: brokerpkg.ReSyncDiscoveryEntity,
		},
	})
	if err != nil {
		return endpoint, ErrPublishBroker(err)
	}

	return endpoint, nil
}

func getOperator(kubeclient *mesherykube.Client) (string, string, error) {
	if kubeclient == nil || kubeclient.KubeClient == nil || kubeclient.KubeClient.AppsV1() == nil {
		return "", "", ErrMesheryClient(nil)
	}

	dep, err := kubeclient.KubeClient.AppsV1().Deployments("meshery").Get(context.TODO(), "meshery-operator", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return "", "", ErrMesheryClient(err)
	}

	version := ""
	if err == nil {
		for _, container := range dep.Spec.Template.Spec.Containers {
			if container.Name == "manager" {
				version = strings.Split(container.Image, ":")[1]
			}
		}
	}

	return dep.ObjectMeta.Name, version, nil
}

func getControllersInfo(mesheryKubeClient *mesherykube.Client, brokerConn broker.Handler, ch chan struct{}) ([]*model.OperatorControllerStatus, error) {
	controllers := make([]*model.OperatorControllerStatus, 0)
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
	if err == nil {
		status := model.StatusEnabled
		if brokerConn.Info() == brokerpkg.NotConnected {
			status = model.StatusDisabled
		}
		controllers = append(controllers, &model.OperatorControllerStatus{
			Name:    "broker",
			Version: broker.Labels["version"],
			Status:  status,
		})
	}

	meshsync, err = mesheryclient.CoreV1Alpha1().MeshSyncs(namespace).Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return controllers, ErrMesheryClient(err)
	}
	if err == nil {
		status := model.StatusDisabled
		flag := false
		for start := time.Now(); time.Since(start) < 5*time.Second; {
			select {
			case <-ch:
				flag = true
				break
			default:
				continue
			}
		}
		if flag {
			status = model.StatusEnabled
		}
		controllers = append(controllers, &model.OperatorControllerStatus{
			Name:    "meshsync",
			Version: meshsync.Labels["version"],
			Status:  status,
		})
	}
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
