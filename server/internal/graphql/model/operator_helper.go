package model

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	operatorv1alpha1 "github.com/meshery/meshery-operator/api/v1alpha1"
	operatorClient "github.com/meshery/meshery-operator/pkg/client"
	"github.com/meshery/meshery/server/models"
	brokerpkg "github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/broker/nats"
	"github.com/meshery/meshkit/logger"

	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/utils"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	Namespace                = "meshery"
	RequestSubject           = "meshery.meshsync.request"
	MeshsyncSubject          = "meshery.meshsync.core"
	BrokerQueue              = "meshery"
	MeshSyncBrokerConnection = "meshsync"
)

type Connections struct {
	Connections []connection `json:"connections"`
}

type connection struct {
	Name string `json:"name"`
}

func GetOperator(kubeclient *mesherykube.Client) (string, string, error) {
	if kubeclient == nil || kubeclient.KubeClient == nil {
		return "", "", ErrMesheryClientNil
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

	return dep.Name, version, nil
}

func GetBrokerInfo(broker controllers.IMesheryController, log logger.Handler) OperatorControllerStatus {
	brokerStatus := broker.GetStatus().String()
	monitorEndpoint, err := broker.GetEndpointForPort("monitor")
	log.Debug("broker monitor endpoint", monitorEndpoint, err)

	if brokerStatus == controllers.Connected.String() {
		brokerEndpoint, _ := broker.GetPublicEndpoint()
		brokerStatus = fmt.Sprintf("%s %s", brokerStatus, brokerEndpoint)
	}
	// change the type of IMesheryController GetName() to to models.MesheryController
	// and use MesheryControllersStatusListItem instead of OperatorControllerStatus
	brokerControllerStatus := OperatorControllerStatus{
		Name:   broker.GetName(),
		Status: Status(brokerStatus),
	}

	brokerControllerStatus.Version, _ = broker.GetVersion()

	return brokerControllerStatus
}

func GetMeshSyncInfo(meshsync controllers.IMesheryController, broker controllers.IMesheryController, log logger.Handler) OperatorControllerStatus {
	if meshsync == nil {
		log.Warn(fmt.Errorf("operator_helper::GetMeshSyncInfo: meshsync controller is nil"))
		return OperatorControllerStatus{
			Status: StatusUnknown,
		}
	}
	meshsyncStatus := meshsync.GetStatus().String()

	// change the type of IMesheryController GetName() to to models.MesheryController
	// and use MesheryControllersStatusListItem instead of OperatorControllerStatus
	if broker == nil {
		meshsyncStatus = controllers.Unknown.String()
	} else {
		// Debug block
		monitorEndpoint, err := broker.GetEndpointForPort("monitor")
		log.Debug("broker monitor endpoint", monitorEndpoint, err)

		if meshsyncStatus == controllers.Connected.String() {
			brokerEndpoint, err := broker.GetPublicEndpoint()
			if err != nil {
				log.Warn(err) // Log the original error
			} else if brokerEndpoint == "" {
				log.Warn(fmt.Errorf("broker public endpoint is empty")) // Log a new error
			} else {
				meshsyncStatus = fmt.Sprintf("%s %s", meshsyncStatus, brokerEndpoint)
			}
		}
	}
	version, _ := meshsync.GetVersion()
	meshsyncControllerStatus := OperatorControllerStatus{
		Name:    meshsync.GetName(),
		Version: version,
		Status:  Status(meshsyncStatus),
	}

	return meshsyncControllerStatus
}

// brokerHostPort parses a "host:port" broker endpoint into a utils.HostPort.
// It returns ok=false for an empty or malformed endpoint (for example a bare
// host with no port, or an unset external endpoint). Parsing the endpoint by
// hand with strings.Split(endpoint, ":")[1] previously panicked with an index
// out of range in exactly those cases, which crashed the server because the
// caller runs SubscribeToBroker in a goroutine with no recover.
func brokerHostPort(endpoint string) (utils.HostPort, bool) {
	host, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return utils.HostPort{}, false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return utils.HostPort{}, false
	}
	return utils.HostPort{Address: host, Port: int32(port)}, true
}

func SubscribeToBroker(_ models.Provider, mesheryKubeClient *mesherykube.Client, datach chan *brokerpkg.Message, brokerConn brokerpkg.Handler, ct *K8sConnectionTracker) (string, error) {
	var broker *operatorv1alpha1.Broker
	var endpoints []string
	if ct != nil {
		endpoints = ct.ListBrokerEndpoints()
	}
	mesheryclient, err := operatorClient.New(&mesheryKubeClient.RestConfig)
	if err != nil {
		if mesheryclient == nil {
			return "", ErrMesheryClientNil
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

	// The loop above can fall through on timeout with broker still nil (the
	// broker resource never became available, or every Get errored). Guard
	// against it so a bare goroutine caller does not panic and crash the server.
	if broker == nil {
		return "", ErrMesheryClient(fmt.Errorf("timed out waiting for meshery-broker; broker resource not available"))
	}

	endpoint := broker.Status.Endpoint.Internal
	if internal, ok := brokerHostPort(broker.Status.Endpoint.Internal); ok {
		if !utils.TcpCheck(&internal, nil) {
			if external, ok := brokerHostPort(broker.Status.Endpoint.External); ok {
				endpoint = broker.Status.Endpoint.External
				if !utils.TcpCheck(&external, nil) {
					if !utils.TcpCheck(&utils.HostPort{
						Address: "host.docker.internal",
						Port:    external.Port,
					}, nil) {
						u, _ := url.Parse(mesheryKubeClient.RestConfig.Host)
						if utils.TcpCheck(&utils.HostPort{
							Address: u.Hostname(),
							Port:    external.Port,
						}, nil) {
							endpoint = fmt.Sprintf("%s:%d", u.Hostname(), external.Port)
						}
					} else {
						endpoint = fmt.Sprintf("host.docker.internal:%d", external.Port)
					}
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

	if err != nil {
		return endpoint, ErrPublishBroker(err)
	}

	return endpoint, nil
}
