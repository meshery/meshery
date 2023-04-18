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
	appsv1 "k8s.io/api/apps/v1"
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
						HostIP:   "127.0.0.1",
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

		port, _ := strconv.Atoi(adapter.Location)

		// Create a deployment
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      adapter.Name,
				Namespace: core.MesheryNamespace,
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": adapter.Name,
					},
				},
				Replicas: int32Ptr(1),
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": adapter.Name,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  adapter.Name,
								Image: "layer5/" + adapter.Name + ":stable-latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: int32(port),
									},
								},
							},
						},
					},
				},
			},
		}

		_, err = kubeClient.KubeClient.AppsV1().Deployments(core.MesheryNamespace).Create(context.Background(), deployment, metav1.CreateOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}

	// switch to default case if the platform specified is not supported
	default:
		return ErrAdapterAdministration(fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform))
	}

	adapter.Location = "localhost:" + adapter.Location
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

		if adapter.Name == "" {
			adapter = a.adapters[adapter.Location]
		}

		// Delete the Deployment
		err = kubeClient.KubeClient.AppsV1().Deployments(core.MesheryNamespace).Delete(context.Background(), adapter.Name, metav1.DeleteOptions{})
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

func int32Ptr(n int32) *int32 {
	return &n
}
