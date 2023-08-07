package helpers

import (
	"context"
	"fmt"
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
		fmt.Printf("cli: %#v\n", cli)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}
		var mesheryNetworkSettings *types.SummaryNetworkSettings
		for _, container := range containers {
			if strings.Contains(container.Image, "thebeginner86/meshery-dev") {
				mesheryNetworkSettings = container.NetworkSettings
			}
		}

		adapterImage := "layer5/" + adapter.Name + ":stable-latest"

		// Pull the latest image
		_, err = cli.ImagePull(ctx, adapterImage, types.ImagePullOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		for netName := range mesheryNetworkSettings.Networks {
			fmt.Printf("network: %#v\n", netName)
			nets, err := cli.NetworkList(ctx, types.NetworkListOptions{})
			fmt.Printf("networks: %#v\n", nets)
			if err != nil {
				// fmt.Errorf("error listing networks: %v", err)
				return ErrAdapterAdministration(err)
			}
			for _, net := range nets {
				if net.Name == netName {
					fmt.Printf("found net\n: %#v", net)
					// Create and start the container
					portNum := adapter.Location
					adapter.Location = "localhost:" + adapter.Location
					port := nat.Port(portNum + "/tcp")
					adapterContainerCreatedBody, err := cli.ContainerCreate(ctx, &container.Config{
						Image: adapterImage,
						ExposedPorts: nat.PortSet{
							port: struct{}{},
						},
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
					}, &network.NetworkingConfig{}, nil, adapter.Name+"-"+fmt.Sprint(time.Now().Unix()))
					if err != nil {
						return ErrAdapterAdministration(err)
					}
					fmt.Printf("adapterContainerCreatedBody: %#v\n", adapterContainerCreatedBody)
					err = cli.NetworkConnect(ctx, net.ID, adapterContainerCreatedBody.ID, &network.EndpointSettings{})
					if err != nil {
						return ErrAdapterAdministration(err)
					}

					fmt.Printf("adapterContainerCreatedBody: %#v\n", adapterContainerCreatedBody)
							if err := cli.ContainerStart(ctx, adapterContainerCreatedBody.ID, types.ContainerStartOptions{}); err != nil {
			return ErrAdapterAdministration(err)
		}
				}
			}
		}

	// switch to default case if the platform specified is not supported
	default:
		return ErrAdapterAdministration(fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes", platform))
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
		fmt.Printf("cli: %#v", cli)
		if err != nil {
			return ErrAdapterAdministration(err)
		}

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			return ErrAdapterAdministration(err)
		}
		var containerID string
		for _, container := range containers {
			for _, p := range container.Ports {
				if strconv.Itoa(int(p.PublicPort)) == adapter.Location {
					containerID = container.ID
					break
				}
			}
		}
		if containerID == "" {
			return ErrAdapterAdministration(fmt.Errorf("no container found for port %s", adapter.Location))
		}

		// Stop and remove the container
		err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force:         true,
			RemoveVolumes: true,
		})
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
