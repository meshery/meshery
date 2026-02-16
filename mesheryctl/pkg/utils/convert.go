package utils

import (
	"sort"
	"strings"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
)

// ToComposeSummaries takes takes []container.Sumamry type and returns
// a []api.ContainerSummary.
func ToComposeSummaries(containersSummary []container.Summary) ([]api.ContainerSummary, error) {
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

		out = append(out, api.ContainerSummary{
			Name:       canonicalContainerName(containerSummary),
			Image:      containerSummary.Image,
			Command:    containerSummary.Command,
			Service:    containerSummary.Labels[api.ServiceLabel],
			Created:    containerSummary.Created,
			State:      string(containerSummary.State),
			Publishers: publishers,
		})
	}

	return out, nil
}

// canonicalContainerName resolves a human-readable identifier for a container.
// It prioritizes a "canonical" name (a single-level name starting with a slash),
// then falls back to the first available name with the prefix trimmed.
// If no names are present, it returns a 12-character short ID or the full ID
// if the ID is shorter than 12 characters.
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
