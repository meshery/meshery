package resources

import "time"

// ConnectionResource represents a Meshery connection state
type ConnectionResource struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Endpoint string `json:"endpoint"`
}

type ConnectionsResponse struct {
	Timestamp   time.Time            `json:"timestamp"`
	Connections []ConnectionResource `json:"connections"`
}

// ProviderResource represents registered Meshery providers
type ProviderResource struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type ProvidersResponse struct {
	Timestamp time.Time          `json:"timestamp"`
	Providers []ProviderResource `json:"providers"`
}

// AdapterResource represents available Meshery adapters
type AdapterResource struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

type AdaptersResponse struct {
	Timestamp time.Time         `json:"timestamp"`
	Adapters  []AdapterResource `json:"adapters"`
}

// HealthResponse represents current Meshery runtime health
type HealthResponse struct {
	Timestamp  time.Time         `json:"timestamp"`
	Status     string            `json:"status"` // "healthy", "degraded", "unhealthy", "unavailable"
	Components map[string]string `json:"components"`
}

// EnvironmentResource represents available Meshery environments
type EnvironmentResource struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	ConnectionIDs []string `json:"connection_references"`
}

type EnvironmentsResponse struct {
	Timestamp    time.Time             `json:"timestamp"`
	Environments []EnvironmentResource `json:"environments"`
}

// ResourceMetadata represents the MCP Resource metadata structure
type ResourceMetadata struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MIMEType    string `json:"mimeType"`
}
