package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// MesheryClient defines the interface required to fetch runtime resource state from Meshery
type MesheryClient interface {
	GetConnections(ctx context.Context) ([]ConnectionResource, error)
	GetProviders(ctx context.Context) ([]ProviderResource, error)
	GetAdapters(ctx context.Context) ([]AdapterResource, error)
	GetHealth(ctx context.Context) (HealthResponse, error)
	GetEnvironments(ctx context.Context) ([]EnvironmentResource, error)
}

// ListResources returns the supported Meshery runtime MCP resources metadata
func ListResources() []ResourceMetadata {
	return []ResourceMetadata{
		{
			URI:         "meshery://connections",
			Name:        "Meshery Connections",
			Description: "Current Meshery connections (clusters, registries, cloud providers)",
			MIMEType:    "application/json",
		},
		{
			URI:         "meshery://providers",
			Name:        "Meshery Providers",
			Description: "Registered Meshery providers and authentication status",
			MIMEType:    "application/json",
		},
		{
			URI:         "meshery://adapters",
			Name:        "Meshery Adapters",
			Description: "Available Meshery infrastructure adapters and running status",
			MIMEType:    "application/json",
		},
		{
			URI:         "meshery://health",
			Name:        "Meshery Runtime Health",
			Description: "Current runtime health status of Meshery components",
			MIMEType:    "application/json",
		},
		{
			URI:         "meshery://environments",
			Name:        "Meshery Environments",
			Description: "Available logical environments and associated connection references",
			MIMEType:    "application/json",
		},
	}
}

// ReadResource handles resource retrieval for meshery:// URIs
func ReadResource(ctx context.Context, uri string, client MesheryClient) ([]byte, error) {
	now := time.Now().UTC()

	var rawData interface{}

	switch uri {
	case "meshery://connections":
		conns, err := client.GetConnections(ctx)
		if err != nil {
			rawData = ConnectionsResponse{
				Timestamp:   now,
				Connections: []ConnectionResource{},
			}
		} else {
			rawData = ConnectionsResponse{
				Timestamp:   now,
				Connections: conns,
			}
		}

	case "meshery://providers":
		providers, err := client.GetProviders(ctx)
		if err != nil {
			rawData = ProvidersResponse{
				Timestamp: now,
				Providers: []ProviderResource{},
			}
		} else {
			rawData = ProvidersResponse{
				Timestamp: now,
				Providers: providers,
			}
		}

	case "meshery://adapters":
		adapters, err := client.GetAdapters(ctx)
		if err != nil {
			rawData = AdaptersResponse{
				Timestamp: now,
				Adapters:  []AdapterResource{},
			}
		} else {
			rawData = AdaptersResponse{
				Timestamp: now,
				Adapters:  adapters,
			}
		}

	case "meshery://health":
		health, err := client.GetHealth(ctx)
		if err != nil {
			rawData = HealthResponse{
				Timestamp:  now,
				Status:     "unavailable",
				Components: map[string]string{"server": "unreachable"},
			}
		} else {
			health.Timestamp = now
			rawData = health
		}

	case "meshery://environments":
		envs, err := client.GetEnvironments(ctx)
		if err != nil {
			rawData = EnvironmentsResponse{
				Timestamp:    now,
				Environments: []EnvironmentResource{},
			}
		} else {
			rawData = EnvironmentsResponse{
				Timestamp:    now,
				Environments: envs,
			}
		}

	default:
		return nil, fmt.Errorf("invalid or unsupported resource URI: %s", uri)
	}

	marshaled, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resource payload: %w", err)
	}

	return SanitizeJSON(marshaled)
}
