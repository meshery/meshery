package helpers

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
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
	platform := utils.GetPlatform()

	// Deploy to current platform
	switch platform {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return ErrDeployingAdapterInDocker(err)
		}
		defer cli.Close()

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			return ErrDeployingAdapterInDocker(err)
		}
		var mesheryNetworkSettings *types.SummaryNetworkSettings
		for _, container := range containers {
			if strings.Contains(container.Image, "layer5/meshery") {
				mesheryNetworkSettings = container.NetworkSettings
			}
		}

		adapterImage := "layer5/" + adapter.Name + ":stable-latest"

		// Pull the latest image
		resp, err := cli.ImagePull(ctx, adapterImage, types.ImagePullOptions{})
		if err != nil {
			return ErrDeployingAdapterInDocker(err)
		}

		defer resp.Close()
		_, err = io.ReadAll(resp)
		if err != nil {
			return ErrDeployingAdapterInDocker(err)
		}

		if mesheryNetworkSettings == nil {
			return ErrDeployingAdapterInDocker(fmt.Errorf("meshery network not found"))
		}

		for netName := range mesheryNetworkSettings.Networks {
			nets, err := cli.NetworkList(ctx, types.NetworkListOptions{})
			if err != nil {
				return ErrDeployingAdapterInDocker(err)
			}
			for _, net := range nets {
				if net.Name == netName {
					// Create and start the container
					portNum := strings.Split(adapter.Location, ":")[1] // eg: for location=meshery-istio:10000, portNum=10000
					port := nat.Port(portNum + "/tcp")
					adapterContainerCreatedBody, err := cli.ContainerCreate(ctx, &container.Config{
						Image: adapterImage,
						ExposedPorts: nat.PortSet{
							port: struct{}{},
						},
						Labels: map[string]string{"com.centurylinklabs.watchtower.enable": "true"},
					}, &container.HostConfig{
						NetworkMode: container.NetworkMode(netName),
						PortBindings: nat.PortMap{
							port: []nat.PortBinding{
								{
									HostIP:   "",
									HostPort: portNum,
								},
							},
						},
					}, &network.NetworkingConfig{
						EndpointsConfig: map[string]*network.EndpointSettings{
							netName: {
								Aliases: []string{adapter.Name},
							},
						},
					}, nil, adapter.Name+"-"+fmt.Sprint(time.Now().Unix()))
					if err != nil {
						return ErrDeployingAdapterInDocker(err)
					}
					err = cli.NetworkConnect(ctx, net.ID, adapterContainerCreatedBody.ID, &network.EndpointSettings{
						Aliases: []string{adapter.Name},
					})
					if err != nil {
						return ErrDeployingAdapterInDocker(err)
					}
					if err := cli.ContainerStart(ctx, adapterContainerCreatedBody.ID, types.ContainerStartOptions{}); err != nil {
						return ErrDeployingAdapterInDocker(err)
					}
				}
			}
		}

	case "kubernetes":
		build := viper.GetString("BUILD")
		_, latestVersion, err := models.CheckLatestVersion(build)
		if err != nil {
			return ErrDeployingAdapterInK8s(err)
		}
		var k8scontext models.K8sContext
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(allContexts) == 0 {
			fmt.Println("No context found")
			return ErrDeployingAdapterInK8s(fmt.Errorf("no context found"))
		}
		for _, k8sctx := range allContexts {
			if k8sctx.Name == "in-cluster" {
				k8scontext = k8sctx
				break
			}
		}

		kubeclient, err := k8scontext.GenerateKubeHandler()
		if err != nil {
			return ErrDeployingAdapterInK8s(err)
		}

		overrideValues := models.SetOverrideValuesForMesheryDeploy(a.GetAdapters(ctx), adapter, true)
		err = kubeclient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
			Namespace:       "meshery",
			ReleaseName:     "meshery",
			CreateNamespace: true,
			ChartLocation: meshkitkube.HelmChartLocation{
				Repository: utils.HelmChartURL,
				Chart:      utils.HelmChartName,
				Version:    latestVersion,
			},
			OverrideValues: overrideValues,
			Action:         meshkitkube.INSTALL,
		})

		if err != nil {
			return ErrDeployingAdapterInK8s(err)
		}

	// switch to default case if the platform specified is not supported
	default:
		return ErrDeployingAdapterInUnknownPlatform(fmt.Errorf("the platform %s is not supported currently. The supported platforms are: Docker and Kubernetes", platform))
	}

	a.AddAdapter(ctx, adapter)
	return nil
}

// RemoveAdapter is used to remove existing adapters from the collection
func (a *AdaptersTracker) UndeployAdapter(ctx context.Context, adapter models.Adapter) error {
	platform := utils.GetPlatform()

	// Undeploy from current platform
	switch platform {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return ErrUnDeployingAdapterInDocker(err)
		}
		defer cli.Close()

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			return ErrUnDeployingAdapterInDocker(err)
		}
		var containerID string
		for _, container := range containers {
			for _, p := range container.Ports {
				if strconv.Itoa(int(p.PublicPort)) == strings.Split(adapter.Location, ":")[1] {
					containerID = container.ID
					break
				}
			}
		}
		if containerID == "" {
			return ErrUnDeployingAdapterInDocker(fmt.Errorf("no container found for port %s", strings.Split(adapter.Location, ":")[1]))
		}

		// Stop and remove the container
		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
		if err != nil {
			return ErrUnDeployingAdapterInDocker(err)
		}

	case "kubernetes":
		build := viper.GetString("BUILD")
		_, latestVersion, err := models.CheckLatestVersion(build)
		if err != nil {
			return ErrUnDeployingAdapterInK8s(err)
		}
		var k8scontext models.K8sContext
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(allContexts) == 0 {
			fmt.Println("No context found")
			return ErrUnDeployingAdapterInK8s(fmt.Errorf("no context found"))
		}
		for _, k8sctx := range allContexts {
			if k8sctx.Name == "in-cluster" {
				k8scontext = k8sctx
				break
			}

		}
		kubeclient, err := k8scontext.GenerateKubeHandler()
		if err != nil {
			return ErrUnDeployingAdapterInK8s(err)
		}

		overrideValues := models.SetOverrideValuesForMesheryDeploy(a.GetAdapters(ctx), adapter, false)
		err = kubeclient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
			Namespace:       "meshery",
			ReleaseName:     "meshery",
			CreateNamespace: true,
			ChartLocation: meshkitkube.HelmChartLocation{
				Repository: utils.HelmChartURL,
				Chart:      utils.HelmChartName,
				Version:    latestVersion,
			},
			OverrideValues: overrideValues,
			Action:         meshkitkube.UNINSTALL,
		})
		if err != nil {
			return ErrUnDeployingAdapterInK8s(err)
		}

	// switch to default case if the platform specified is not supported
	default:
		return ErrUnDeployingAdapterInUnknownPlatform(fmt.Errorf("the platform %s is not supported currently. The supported platforms are: Docker and Kubernetes", platform))
	}

	a.RemoveAdapter(ctx, adapter)
	return nil
}
