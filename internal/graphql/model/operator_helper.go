package model

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
	operatorClient "github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/models"
	brokerpkg "github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace       = "meshery"
	RequestSubject  = "meshery.meshsync.request"
	MeshsyncSubject = "meshery.meshsync.core"
	BrokerQueue     = "meshery"
)

func Initialize(client *mesherykube.Client, delete bool, adapterTracker models.AdaptersTrackerInterface) error {
	// installOperator
	err := installUsingHelm(client, delete, adapterTracker)
	if err != nil {
		return err
	}

	return nil
}

func GetOperator(kubeclient *mesherykube.Client) (string, string, error) {
	if kubeclient == nil || kubeclient.KubeClient == nil {
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

func GetControllersInfo(mesheryKubeClient *mesherykube.Client, brokerConn brokerpkg.Handler, ch chan struct{}) ([]*OperatorControllerStatus, error) {
	controllers := make([]*OperatorControllerStatus, 0)

	mesheryclient, err := operatorClient.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return controllers, ErrMesheryClient(nil)
		}
		return controllers, ErrMesheryClient(err)
	}

	broker, err := GetBrokerInfo(mesheryclient, brokerConn)
	if err != nil {
		return controllers, err
	}
	controllers = append(controllers, &broker)

	meshsync, err := GetMeshSyncInfo(mesheryclient, ch)
	if err != nil {
		return controllers, err
	}
	controllers = append(controllers, &meshsync)

	return controllers, nil
}

func GetBrokerInfo(mesheryclient operatorClient.Interface, brokerConn brokerpkg.Handler) (OperatorControllerStatus, error) {
	var brokerStatus OperatorControllerStatus

	broker, err := mesheryclient.CoreV1Alpha1().Brokers(Namespace).Get(context.TODO(), "meshery-broker", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return brokerStatus, ErrMesheryClient(err)
	}
	if err == nil {
		brokerStatus.Status = StatusConnected
		if brokerConn.Info() == brokerpkg.NotConnected {
			brokerStatus.Status = StatusEnabled
		}
		brokerStatus.Name = "broker"
		brokerStatus.Version = broker.Labels["version"]
	}

	return brokerStatus, nil
}

func GetMeshSyncInfo(mesheryclient operatorClient.Interface, ch chan struct{}) (OperatorControllerStatus, error) {
	var meshsyncStatus OperatorControllerStatus

	meshsync, err := mesheryclient.CoreV1Alpha1().MeshSyncs(Namespace).Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return meshsyncStatus, ErrMesheryClient(err)
	}

	// Synthetic Check for MeshSync data is too time consuming. Commented for now.

	// if err == nil {
	// 	status := StatusDisabled
	// 	flag := false
	// 	for start := time.Now(); time.Since(start) < 5*time.Second; {
	// 		select {
	// 		case <-ch:
	// 			flag = true
	// 			break
	// 		default:
	// 			continue
	// 		}
	// 	}
	// 	if flag {
	// 		status = StatusEnabled
	// 	}
	// 	meshsyncStatus = OperatorControllerStatus{
	// 		Name:    "meshsync",
	// 		Version: meshsync.Labels["version"],
	// 		Status:  status,
	// 	})
	// }

	meshsyncStatus.Status = StatusEnabled
	meshsyncStatus.Name = "meshsync"
	meshsyncStatus.Version = meshsync.Labels["version"]

	return meshsyncStatus, nil
}

func SubscribeToBroker(provider models.Provider, mesheryKubeClient *mesherykube.Client, datach chan *brokerpkg.Message, brokerConn brokerpkg.Handler) (string, error) {
	var broker *operatorv1alpha1.Broker

	mesheryclient, err := operatorClient.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return "", ErrMesheryClient(nil)
		}
		return "", ErrMesheryClient(err)
	}

	timeout := 60
	for timeout > 0 {
		broker, err = mesheryclient.CoreV1Alpha1().Brokers(Namespace).Get(context.Background(), "meshery-broker", metav1.GetOptions{})
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
					u, _ := url.Parse(mesheryKubeClient.RestConfig.Host)
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
	conn.DeepCopyInto(brokerConn)

	err = brokerConn.SubscribeWithChannel(MeshsyncSubject, BrokerQueue, datach)
	if err != nil {
		return endpoint, ErrSubscribeChannel(err)
	}

	err = brokerConn.Publish(RequestSubject, &brokerpkg.Message{
		Request: &brokerpkg.RequestObject{
			Entity: brokerpkg.ReSyncDiscoveryEntity,
		},
	})
	if err != nil {
		return endpoint, ErrPublishBroker(err)
	}

	return endpoint, nil
}
