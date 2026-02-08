package utils

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
)

type containerInspector interface {
	ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error)
}

func ToComposeSummaries(ctx context.Context, inspector containerInspector, containersSummary []container.Summary) ([]api.ContainerSummary, error) {
	fmt.Println("============== containersSummary ================")
	fmt.Printf("%+v\n", containersSummary)
	out := make([]api.ContainerSummary, 0, len(containersSummary))

	for _, containerSummary := range containersSummary {
		ports := append([]container.Port(nil), containerSummary.Ports...)
		sort.Slice(ports, func(i, j int) bool { return ports[i].PrivatePort < ports[j].PrivatePort })

		publishers := make(api.PortPublishers, 0, len(ports))
		for _, port := range ports {
			publishers = append(publishers, api.PortPublisher{
				URL:           port.IP,
				TargetPort:    int(port.PrivatePort),
				PublishedPort: int(port.PublicPort),
				Protocol:      port.Type,
			})
		}

		mounts := make([]string, 0, len(containerSummary.Mounts))
		localVolumes := 0
		for _, mount := range containerSummary.Mounts {
			name := mount.Name
			if name == "" {
				name = mount.Source
			}
			if mount.Driver == "local" {
				localVolumes++
			}
			mounts = append(mounts, name)
		}

		networks := []string{}
		if containerSummary.NetworkSettings != nil {
			for n := range containerSummary.NetworkSettings.Networks {
				networks = append(networks, n)
			}
			sort.Strings(networks)
		}

		inspect, err := inspector.ContainerInspect(ctx, containerSummary.ID)
		if err != nil {
			return nil, fmt.Errorf("inspect container %s: %w", containerSummary.ID, err)
		}

		health := ""
		exitCode := 0
		if inspect.State != nil {
			if inspect.State.Health != nil {
				health = inspect.State.Health.Status
			}
			if inspect.State.Status == container.StateExited || inspect.State.Status == container.StateDead {
				exitCode = int(inspect.State.ExitCode)
			}
		}

		out = append(out, api.ContainerSummary{
			ID:           containerSummary.ID,
			Name:         canonicalContainerName(containerSummary),
			Names:        containerSummary.Names,
			Image:        containerSummary.Image,
			Command:      containerSummary.Command,
			Project:      containerSummary.Labels[api.ProjectLabel],
			Service:      containerSummary.Labels[api.ServiceLabel],
			Created:      containerSummary.Created,
			State:        string(containerSummary.State),
			Status:       containerSummary.Status,
			Health:       health,
			ExitCode:     exitCode,
			Publishers:   publishers,
			Labels:       containerSummary.Labels,
			SizeRw:       containerSummary.SizeRw,
			SizeRootFs:   containerSummary.SizeRootFs,
			Mounts:       mounts,
			Networks:     networks,
			LocalVolumes: localVolumes,
		})
	}

	return out, nil
}

func canonicalContainerName(containerSummary container.Summary) string {
	if len(containerSummary.Names) == 0 {
		if len(containerSummary.ID) > 12 {
			return containerSummary.ID[:12]
		}
		return containerSummary.ID
	}
	for _, name := range containerSummary.Names {
		if strings.LastIndex(name, "/") == 0 {
			return name[1:]
		}
	}
	return strings.TrimPrefix(containerSummary.Names[0], "/")
}
