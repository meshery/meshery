package models

import (
	"encoding/gob"
	"time"

	"github.com/grafana-tools/sdk"
)

// K8SNode - represents a kubernetes node
type K8SNode struct {
	InternalIP              string `json:"internal_ip,omitempty"`
	HostName                string `json:"hostname,omitempty"`
	AllocatableCPU          string `json:"allocatable_cpu,omitempty"`
	AllocatableMemory       string `json:"allocatable_memory,omitempty"`
	CapacityCPU             string `json:"capacity_cpu,omitempty"`
	CapacityMemory          string `json:"capacity_memory,omitempty"`
	OSImage                 string `json:"os_image,omitempty"`
	OperatingSystem         string `json:"operating_system,omitempty"`
	KubeletVersion          string `json:"kubelet_version,omitempty"`
	KubeProxyVersion        string `json:"kubeproxy_version,omitempty"`
	ContainerRuntimeVersion string `json:"container_runtime_version,omitempty"`
	Architecture            string `json:"architecture,omitempty"`
}

// Grafana represents the Grafana session config
type Grafana struct {
	GrafanaURL    string `json:"grafanaURL,omitempty"`
	GrafanaAPIKey string `json:"grafanaAPIKey,omitempty"`
	// GrafanaBoardSearch string          `json:"grafanaBoardSearch,omitempty"`
	GrafanaBoards []*SelectedGrafanaConfig `json:"selectedBoardsConfigs,omitempty"`
}

// SelectedGrafanaConfig represents the selected boards, panels, and template variables
type SelectedGrafanaConfig struct {
	GrafanaBoard         *GrafanaBoard `json:"board,omitempty"`
	GrafanaPanels        []*sdk.Panel  `json:"panels,omitempty"`
	SelectedTemplateVars []string      `json:"templateVars,omitempty"`
}

// Prometheus represents the prometheus session config
type Prometheus struct {
	PrometheusURL                   string                   `json:"prometheusURL,omitempty"`
	SelectedPrometheusBoardsConfigs []*SelectedGrafanaConfig `json:"selectedPrometheusBoardsConfigs,omitempty"`
}

// LoadTestPreferences represents the load test preferences
type LoadTestPreferences struct {
	ConcurrentRequests int    `json:"c,omitempty"`
	QueriesPerSecond   int    `json:"qps,omitempty"`
	Duration           string `json:"t,omitempty"`
	LoadGenerator      string `json:"gen,omitempty"`
}

// Parameters to updates Anonymous stats
type PreferenceParams struct {
	AnonymousUsageStats  bool `json:"anonymousUsageStats"`
	AnonymousPerfResults bool `json:"anonymousPerfResults"`
}

// Preference represents the data stored in session / local DB
type Preference struct {
	MeshAdapters              []*Adapter             `json:"meshAdapters,omitempty"`
	Grafana                   *Grafana               `json:"grafana,omitempty"`
	Prometheus                *Prometheus            `json:"prometheus,omitempty"`
	LoadTestPreferences       *LoadTestPreferences   `json:"loadTestPrefs,omitempty"`
	AnonymousUsageStats       bool                   `json:"anonymousUsageStats"`
	AnonymousPerfResults      bool                   `json:"anonymousPerfResults"`
	UpdatedAt                 time.Time              `json:"updated_at,omitempty"`
	UsersExtensionPreferences map[string]interface{} `json:"usersExtensionPreferences,omitempty"`
}

func init() {
	gob.Register(&Preference{})
	gob.Register(map[string]interface{}{})
}

// PreferencePersister defines methods for a session persister
type PreferencePersister interface {
	ReadFromPersister(userID string) (*Preference, error)
	WriteToPersister(userID string, data *Preference) error
	DeleteFromPersister(userID string) error

	// Lock(userID string)
	// Unlock(userID string)
	// ClosePersister()
}

// Parameters to save Grafana configuration
type GrafanaConfigParams struct {
	GrafanaURL    string `json:"grafanaURL,omitempty"`
	GrafanaAPIKey string `json:"grafanaAPIKey,omitempty"`
}
