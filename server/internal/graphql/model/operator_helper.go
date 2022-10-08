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
	"github.com/layer5io/meshery/server/models"
	brokerpkg "github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/broker/nats"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace       = "meshery"
	RequestSubject  = "meshery.meshsync.request"
	MeshsyncSubject = "meshery.meshsync.core"
	BrokerQueue     = "meshery"
)

var (
	meshsyncVersion string
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

func GetControllersInfo(mesheryKubeClient *mesherykube.Client, brokerConn brokerpkg.Handler) ([]*OperatorControllerStatus, error) {
	controllers := make([]*OperatorControllerStatus, 0)

	mesheryclient, err := operatorClient.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return controllers, ErrMesheryClient(nil)
		}
		return controllers, ErrMesheryClient(err)
	}

	broker, err := GetBrokerInfo(mesheryclient, mesheryKubeClient, brokerConn)
	if err != nil {
		return controllers, err
	}

	controllers = append(controllers, &broker)

	meshsync, err := GetMeshSyncInfo(mesheryclient, mesheryKubeClient)
	if err != nil {
		return controllers, err
	}

	controllers = append(controllers, &meshsync)

	return controllers, nil
}

func GetBrokerInfo(mesheryclient operatorClient.Interface, mesheryKubeClient *mesherykube.Client, brokerConn brokerpkg.Handler) (OperatorControllerStatus, error) {
	var brokerStatus OperatorControllerStatus

	broker, err := mesheryclient.CoreV1Alpha1().Brokers(Namespace).Get(context.TODO(), "meshery-broker", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return brokerStatus, ErrMesheryClient(err)
	}
	statefulSet, err := mesheryKubeClient.KubeClient.AppsV1().StatefulSets("meshery").Get(context.TODO(), "meshery-broker", metav1.GetOptions{})
	brokerVersion := ""
	if err == nil {
		brokerVersion = imageVersionExtractUtil(statefulSet.Spec.Template, "nats")
	}
	if err == nil {
		status := fmt.Sprintf("%s %s", StatusConnected, broker.Status.Endpoint.External)
		if brokerConn.Info() != brokerpkg.NotConnected {
			brokerStatus.Status = Status(status)
		}
		brokerStatus.Name = "broker"
		brokerStatus.Version = brokerVersion
	}

	return brokerStatus, nil
}

func GetMeshSyncInfo(mesheryclient operatorClient.Interface, mesheryKubeClient *mesherykube.Client) (OperatorControllerStatus, error) {
	var meshsyncStatus OperatorControllerStatus
	meshsync, err := mesheryclient.CoreV1Alpha1().MeshSyncs(Namespace).Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	if err != nil && !kubeerror.IsNotFound(err) {
		return meshsyncStatus, ErrMesheryClient(err)
	}

	meshsyncDeployment, err := mesheryKubeClient.KubeClient.AppsV1().Deployments("meshery").Get(context.TODO(), "meshery-meshsync", metav1.GetOptions{})
	meshsyncVersion := ""
	var status string
	if err == nil {
		status = fmt.Sprintf("%s %s", StatusEnabled, meshsync.Status.PublishingTo)
		meshsyncVersion = imageVersionExtractUtil(meshsyncDeployment.Spec.Template, "meshsync")
	}
	if err != nil && kubeerror.IsNotFound(err) {
		status = StatusDisabled.String()
	}

	meshsyncStatus.Status = Status(status)
	meshsyncStatus.Name = "meshsync"
	meshsyncStatus.Version = meshsyncVersion
	return meshsyncStatus, nil
}

func SubscribeToBroker(provider models.Provider, mesheryKubeClient *mesherykube.Client, datach chan *brokerpkg.Message, brokerConn brokerpkg.Handler, ct *K8sConnectionTracker) (string, error) {
	var broker *operatorv1alpha1.Broker
	var endpoints []string
	if ct != nil {
		endpoints = ct.ListBrokerEndpoints()
	}
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
	endpoints = append(endpoints, endpoint)
	// subscribing to nats
	conn, err := nats.New(nats.Options{
		URLS:           endpoints,
		ConnectionName: "meshery",
		Username:       "",
		Password:       "",
		ReconnectWait:  2 * time.Second,
		MaxReconnect:   5,
	})
	// Hack for minikube based clusters
	if err != nil && conn == nil {
		return endpoint, err
	}
	defer func() {
		if conn == nil {
			return
		}
		available := make(map[string]bool)
		for _, server := range conn.ConnectedEndpoints() {
			available[server] = true
		}
		ct.ResetEndpoints(available)
	}()
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

	go getVersion(brokerConn)

	if err != nil {
		return endpoint, ErrPublishBroker(err)
	}

	return endpoint, nil
}

func imageVersionExtractUtil(container v1.PodTemplateSpec, containerName string) string {
	version := ""
	for _, container := range container.Spec.Containers {
		if strings.Compare(container.Name, containerName) == 0 {
			version = strings.Split(container.Image, ":")[1]
		}
	}
	return version
}

func getVersion(brokerConn brokerpkg.Handler) {
	versionch := make(chan *brokerpkg.Message)

	err := brokerConn.SubscribeWithChannel("meshsync-meta", "meshery", versionch) // what is this queue used for now just using "meshery"

	if err != nil {
		logrus.Error(err.Error())
		return
	}
	err = brokerConn.Publish(RequestSubject, &brokerpkg.Message{
		Request: &brokerpkg.RequestObject{
			Entity: "meshsync-meta",
		},
	})
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	ch := <-versionch
	meshsyncVersion = "stable-" + ch.Object.(string)
	fmt.Println("VERSION", meshsyncVersion)
}
