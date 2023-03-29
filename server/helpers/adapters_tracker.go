package helpers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/layer5io/meshery/core"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
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
func (a *AdaptersTracker) AddAdapter(ctx context.Context, adapter models.Adapter) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	a.adapters[adapter.Location] = adapter
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) RemoveAdapter(ctx context.Context, adapter models.Adapter) {
	a.adaptersLock.Lock()
	defer a.adaptersLock.Unlock()
	delete(a.adapters, adapter.Location)
}

// GetAdapters returns the list of existing adapters
func (a *AdaptersTracker) GetAdapters(ctx context.Context) []models.Adapter {
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
	platform := utils.GetPlatform()

	// Deploy to current platform
	switch platform {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}

		adapterImage := "layer5/" + adapter.Name

		// Pull the latest image
		reader, err := cli.ImagePull(ctx, adapterImage, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer reader.Close()
		_, _ = io.Copy(os.Stdout, reader)

		// Create and start the container
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: adapterImage,
			ExposedPorts: nat.PortSet{
				"10000/tcp": struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				"10000/tcp": []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: "10000",
					},
				},
			},
		}, &network.NetworkingConfig{}, nil, adapter.Name+"-"+fmt.Sprint(time.Now().Unix()))
		if err != nil {
			return err
		}

		adapter.Name = adapter.Name + ":" + resp.ID
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}

	case "kubernetes":
		if err := adapterDeployKubernetes(adapter, meshkitkube.UNINSTALL); err != nil {
			return err
		}

	// switch to default case if the platform specified is not supported
	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform)
	}

	a.AddAdapter(ctx, adapter)
	return nil
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) UndeployAdapter(ctx context.Context, adapter models.Adapter) error {
	platform := utils.GetPlatform()

	// Undeploy from current platform
	switch utils.GetPlatform() {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}

		// Stop and remove the container
		containerID := strings.Split(adapter.Name, ":")[1]
		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			return err
		}

	case "kubernetes":
		if err := adapterDeployKubernetes(adapter, meshkitkube.UNINSTALL); err != nil {
			return err
		}

	// switch to default case if the platform specified is not supported
	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform)
	}

	a.RemoveAdapter(ctx, adapter)
	return nil
}

func adapterDeployKubernetes(adapter models.Adapter, act meshkitkube.HelmChartAction) error {
	kubeClient, err := meshkitkube.New([]byte(""))
	if err != nil {
		return err
	}

	if err := core.CreateManifestsFolder(); err != nil {
		return err
	}

	// Applying Meshery Helm charts for installing selected adapter
	mesheryReleaseVersion := viper.GetString("BUILD")
	err = kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
		Namespace:       "meshery",
		ReleaseName:     adapter.Name,
		CreateNamespace: true,
		ChartLocation: meshkitkube.HelmChartLocation{
			Repository: models.ChartRepo,
			Chart:      adapter.Name,
			Version:    mesheryReleaseVersion,
		},
		OverrideValues: nil,
		Action:         act,
		// the helm chart will be downloaded to ~/.meshery/manifests if it doesn't exist
		DownloadLocation: "~/.meshery/manifests",
		DryRun:           false,
	})
	return err
}
