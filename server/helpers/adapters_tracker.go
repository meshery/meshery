package helpers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/layer5io/meshery/server/core"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AdaptersTracker is used to hold the list of known adapters
type AdaptersTracker struct {
	adapters     map[string]models.Adapter
	adaptersLock *sync.Mutex
}

// NewAdaptersTracker returns an instance of AdaptersTracker
func NewAdaptersTracker(adapterURLs []string) *AdaptersTracker {
	initialAdapters := map[string]models.Adapter{}
	for _, u := range adapterURLs {
		initialAdapters[u] = models.Adapter{
			Location: u,
		}
	}
	a := &AdaptersTracker{
		adapters:     initialAdapters,
		adaptersLock: &sync.Mutex{},
	}

	return a
}

// AddAdapter is used to add new adapters to the collection
func (a *AdaptersTracker) AddAdapter(_ context.Context, adapter models.Adapter) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	a.adapters[adapter.Location] = adapter
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) RemoveAdapter(_ context.Context, adapter models.Adapter) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	delete(a.adapters, adapter.Location)
}

// GetAdapters returns the list of existing adapters
func (a *AdaptersTracker) GetAdapters(_ context.Context) []models.Adapter {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()

	ad := make([]models.Adapter, 0)
	for _, x := range a.adapters {
		ad = append(ad, x)
	}
	return ad
}

// AddAdapter is used to add new adapters to the collection
func (a *AdaptersTracker) DeployAdapter(ctx context.Context, adapter models.Adapter) error {
	platform, err := utils.GetPlatform()
	if err != nil {
		return ErrDockerHost(err)
	}

	// Deploy to current platform
	targetIP := "127.0.0.1"
	switch platform {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		adapterImage := "layer5/" + adapter.Name + ":stable-latest"

		// Pull the latest image
		reader, err := cli.ImagePull(ctx, adapterImage, types.ImagePullOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}
		defer reader.Close()
		_, _ = io.Copy(os.Stdout, reader)

		// Create and start the container
		portNum := adapter.Location
		port := nat.Port(portNum + "/tcp")
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: adapterImage,
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{
					{
						HostIP:   targetIP,
						HostPort: portNum,
					},
				},
			},
		}, &network.NetworkingConfig{}, nil, adapter.Name+"-"+fmt.Sprint(time.Now().Unix()))
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return ErrAdapterAdministration(err)
		}

	case "kubernetes":
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		exposedPort, err := strconv.Atoi(adapter.Location)
		if err != nil {
			return ErrAdapterAdministration(err)
		}
		err = a.applyHelmCharts(kubeClient, adapter.Name, exposedPort, true)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		svc, err := kubeClient.KubeClient.CoreV1().Services(core.MesheryNamespace).Get(context.Background(), adapter.Name, metav1.GetOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		targetIP = svc.Spec.ClusterIP

	// switch to default case if the platform specified is not supported
	default:
		return ErrAdapterAdministration(fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform))
	}

	adapter.Location = targetIP + ":" + adapter.Location
	a.AddAdapter(ctx, adapter)
	return nil
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) UndeployAdapter(ctx context.Context, adapter models.Adapter) error {
	platform, err := utils.GetPlatform()
	if err != nil {
		return ErrDockerHost(err)
	}

	// Undeploy from current platform
	switch platform {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		// Find the container ID by given exposed port
		port := strings.Split(adapter.Location, ":")[1]
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}
		var containerID string
		for _, container := range containers {
			for _, p := range container.Ports {
				if strconv.Itoa(int(p.PublicPort)) == port {
					containerID = container.ID
					break
				}
			}
		}
		if containerID == "" {
			return ErrAdapterAdministration(fmt.Errorf("no container found for port %s", port))
		}

		// Stop and remove the container
		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			return ErrAdapterAdministration(err)
		}

	case "kubernetes":
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		err = a.applyHelmCharts(kubeClient, adapter.Name, 0, false)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

	// switch to default case if the platform specified is not supported
	default:
		return ErrAdapterAdministration(fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform))
	}

	a.RemoveAdapter(ctx, adapter)
	return nil
}

// SetOverrideValues returns the value overrides to install/upgrade helm chart
func (a *AdaptersTracker) setOverrideValues(adapterName string, exposedPort int, enable bool) map[string]interface{} {
	// Initialize the valueOverrides with default values and set config for selected adapter
	valueOverrides := map[string]interface{}{}
	for _, availableAdapter := range models.ListAvailableAdapters {
		valueOverrides[availableAdapter.Name] = map[string]interface{}{
			"enabled": false,
		}
		if enable && availableAdapter.Name == adapterName {
			valueOverrides[availableAdapter.Name] = map[string]interface{}{
				"enabled": true,
				"service": setService(adapterName, exposedPort),
			}
		}
	}

	// Enable all already existing adapters
	for _, setAdapter := range a.adapters {
		setAdapterPort, _ := strconv.Atoi(strings.Split(setAdapter.Location, ":")[1])
		valueOverrides[setAdapter.Name] = map[string]interface{}{
			"enabled": true,
			"service": setService(setAdapter.Name, setAdapterPort),
		}
	}

	return valueOverrides
}

// Apply Helm charts for adapter deployment and service
func (a *AdaptersTracker) applyHelmCharts(kubeClient *meshkitkube.Client, adapterName string, exposedPort int, enable bool) error {
	// get value overrides to install the helm chart
	overrideValues := a.setOverrideValues(adapterName, exposedPort, enable)

	// install the helm charts with specified override values
	return kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
		Namespace:       core.MesheryNamespace,
		ReleaseName:     core.MesheryReleaseName,
		CreateNamespace: true,
		ChartLocation: meshkitkube.HelmChartLocation{
			Repository: utils.HelmChartURL,
			Chart:      utils.HelmChartName,
		},
		Action:         meshkitkube.INSTALL,
		DryRun:         false,
		OverrideValues: overrideValues,
	})
}

func setService(name string, port int) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: core.MesheryNamespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "gRPC",
					Protocol: corev1.ProtocolTCP,
					Port:     int32(port),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}
