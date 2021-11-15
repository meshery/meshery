package helpers

import (
	"context"
	"path/filepath"
	"time"

	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewDynamicClientGenerator takes in the kubeconfig *directory* path and returns a
// function which can generate dynamic client
func NewDynamicClientGenerator(path string) func() (dynamic.Interface, error) {
	configPath := filepath.Join(path, "config")

	config, err := utils.ReadLocalFile(configPath)
	return func() (dynamic.Interface, error) {
		if err != nil {
			return nil, ErrNewDynamicClientGenerator(err)
		}

		return NewDynamicClient([]byte(config))
	}
}

// NewKubeClientGenerator takes in the kubeconfig *directory* path and returns a
// function which can generate dynamic client
func NewKubeClientGenerator(path string) func() (*mesherykube.Client, error) {
	configPath := filepath.Join(path, "config")

	config, err := utils.ReadLocalFile(configPath)
	return func() (*mesherykube.Client, error) {
		if err != nil {
			return nil, ErrNewKubeClientGenerator(err)
		}

		return NewKubeClient([]byte(config))
	}
}

// NewDynamicClient generates new dynamic go client
func NewDynamicClient(kubeconfig []byte) (dynamic.Interface, error) {
	var (
		restConfig *rest.Config
		err        error
	)

	if len(kubeconfig) > 0 {
		restConfig, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig)
		if err != nil {
			return nil, ErrRestConfigFromKubeConfig(err)
		}
	} else {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, ErrInClusterConfig(err)
		}
	}

	return dynamic.NewForConfig(restConfig)
}

func NewKubeClient(kubeconfig []byte) (*mesherykube.Client, error) {
	client, err := mesherykube.New(kubeconfig)
	if err != nil {
		return nil, ErrNewKubeClient(err)
	}

	return client, nil
}

func getK8SClientSet(kubeconfig []byte, contextName string) (*kubernetes.Clientset, error) {
	var clientConfig *rest.Config
	var err error
	if len(kubeconfig) == 0 {
		clientConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, ErrInClusterConfig(err)
		}
	} else {
		config, err := clientcmd.Load(kubeconfig)
		if err != nil {
			return nil, ErrInvalidK8SConfig(err)
		}
		if contextName != "" {
			config.CurrentContext = contextName
		}
		clientConfig, err = clientcmd.NewDefaultClientConfig(*config, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return nil, ErrClientConfig(err)
		}
	}
	clientConfig.Timeout = 2 * time.Minute
	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, ErrClientSet(err)
	}
	return clientset, nil
}

// FetchKubernetesNodes - function used to fetch nodes metadata
func FetchKubernetesNodes(kubeconfig []byte, contextName string) ([]*models.K8SNode, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return nil, ErrFetchKubernetesNodes(err)
	}
	var nodes []*models.K8SNode

	// nodes
	nodesClient := clientset.CoreV1().Nodes()
	logrus.Debugf("Listing nodes")
	nodelist, err := nodesClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, ErrFetchNodes(err)
	}
	for _, n := range nodelist.Items {
		// logrus.Debugf(" * %s (%d replicas)", n.Name, *d.Spec.Replicas)
		node := &models.K8SNode{}
		addresses := n.Status.Addresses
		for _, address := range addresses {
			logrus.Debugf("Type: %s, Address: %s", address.Type, address.Address)
			if address.Type == "InternalIP" {
				node.InternalIP = address.Address
			} else if address.Type == "Hostname" {
				node.HostName = address.Address
			}
		}

		logrus.Debugf("Allocatable CPU: %s", n.Status.Allocatable.Cpu())
		node.AllocatableCPU = n.Status.Allocatable.Cpu().String()
		logrus.Debugf("Allocatable CPU: %s", n.Status.Allocatable.Memory())
		node.AllocatableMemory = n.Status.Allocatable.Memory().String()
		logrus.Debugf("Capacity CPU: %s", n.Status.Capacity.Cpu())
		node.CapacityCPU = n.Status.Capacity.Cpu().String()
		logrus.Debugf("Capacity CPU: %s", n.Status.Capacity.Memory())
		node.CapacityMemory = n.Status.Capacity.Memory().String()

		nodeInfo := n.Status.NodeInfo
		logrus.Debugf("OS Image: %s", nodeInfo.OSImage)
		node.OSImage = nodeInfo.OSImage
		logrus.Debugf("Operating system: %s", nodeInfo.OperatingSystem)
		node.OperatingSystem = nodeInfo.OperatingSystem
		logrus.Debugf("Kubelet version: %s", nodeInfo.KubeletVersion)
		node.KubeletVersion = nodeInfo.KubeletVersion
		logrus.Debugf("Kubeproxy version: %s", nodeInfo.KubeProxyVersion)
		node.KubeProxyVersion = nodeInfo.KubeProxyVersion
		logrus.Debugf("Container runtime version: %s", nodeInfo.ContainerRuntimeVersion)
		node.ContainerRuntimeVersion = nodeInfo.ContainerRuntimeVersion
		logrus.Debugf("Architecture: %s", nodeInfo.Architecture)
		node.Architecture = nodeInfo.Architecture

		nodes = append(nodes, node)
	}
	return nodes, nil
}

// FetchKubernetesVersion - function used to fetch kubernetes server version
func FetchKubernetesVersion(kubeconfig []byte, contextName string) (string, error) {
	clientset, err := getK8SClientSet(kubeconfig, contextName)
	if err != nil {
		return "", ErrFetchKubernetesVersion(err)
	}

	serverVersion, err := clientset.ServerVersion()
	if err != nil {
		return "", ErrFetchKubernetesVersion(err)
	}
	logrus.Debugf("Kubernetes API Server version: %s", serverVersion.String())
	return serverVersion.String(), nil
}
