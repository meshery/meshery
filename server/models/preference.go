package models

import (
	"encoding/gob"
	"encoding/json"
	"time"

	"github.com/grafana-tools/sdk"
)

// K8SNode - represents a kubernetes node
type K8SNode struct {
	InternalIP              string `json:"internalIp,omitempty"`
	HostName                string `json:"hostname,omitempty"`
	AllocatableCPU          string `json:"allocatableCpu,omitempty"`
	AllocatableMemory       string `json:"allocatableMemory,omitempty"`
	CapacityCPU             string `json:"capacityCpu,omitempty"`
	CapacityMemory          string `json:"capacityMemory,omitempty"`
	OSImage                 string `json:"osImage,omitempty"`
	OperatingSystem         string `json:"operatingSystem,omitempty"`
	KubeletVersion          string `json:"kubeletVersion,omitempty"`
	KubeProxyVersion        string `json:"kubeproxyVersion,omitempty"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion,omitempty"`
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
	MeshAdapters                      []*Adapter             `json:"meshAdapters,omitempty"`
	Grafana                           *Grafana               `json:"grafana,omitempty"`
	Prometheus                        *Prometheus            `json:"prometheus,omitempty"`
	LoadTestPreferences               *LoadTestPreferences   `json:"loadTestPrefs,omitempty"`
	AnonymousUsageStats               bool                   `json:"anonymousUsageStats"`
	AnonymousPerfResults              bool                   `json:"anonymousPerfResults"`
	UpdatedAt                         time.Time              `json:"updatedAt,omitempty"`
	DashboardPreferences              map[string]interface{} `json:"dashboardPreferences,omitempty"`
	SelectedOrganizationID            string                 `json:"selectedOrganizationID,omitempty"`
	SelectedWorkspaceForOrganizations map[string]string      `json:"selectedWorkspaceForOrganizations,omitempty"` // map[orgID]workspaceID
	UsersExtensionPreferences         map[string]interface{} `json:"usersExtensionPreferences,omitempty"`
	RemoteProviderPreferences         map[string]interface{} `json:"remoteProviderPreferences,omitempty"`
}

func (p Preference) MarshalJSON() ([]byte, error) {
	type preferenceAlias Preference
	data, err := json.Marshal(preferenceAlias(p))
	if err != nil {
		return nil, err
	}
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	if p.SelectedOrganizationID == "" {
		return marshalPreferenceCanonicalFields(fields, p)
	}
	fields["selectedOrganizationId"] = p.SelectedOrganizationID
	return marshalPreferenceCanonicalFields(fields, p)
}

func marshalPreferenceCanonicalFields(fields map[string]interface{}, p Preference) ([]byte, error) {
	if p.Grafana != nil {
		grafana, _ := fields["grafana"].(map[string]interface{})
		if grafana == nil {
			grafana = map[string]interface{}{}
		}
		if p.Grafana.GrafanaURL != "" {
			grafana["grafanaUrl"] = p.Grafana.GrafanaURL
		}
		if p.Grafana.GrafanaAPIKey != "" {
			grafana["grafanaApiKey"] = p.Grafana.GrafanaAPIKey
		}
		fields["grafana"] = grafana
	}
	if p.Prometheus != nil {
		prometheus, _ := fields["prometheus"].(map[string]interface{})
		if prometheus == nil {
			prometheus = map[string]interface{}{}
		}
		if p.Prometheus.PrometheusURL != "" {
			prometheus["prometheusUrl"] = p.Prometheus.PrometheusURL
		}
		fields["prometheus"] = prometheus
	}
	return json.Marshal(fields)
}

// NewDefaultPreference returns a preference initialized with Meshery's default opt-in values.
func NewDefaultPreference() *Preference {
	return &Preference{
		AnonymousUsageStats:  true,
		AnonymousPerfResults: true,
	}
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

// CapabilitiesPersister defines methods for a capability persister
type CapabilitiesPersister interface {
	ReadCapabilitiesForUser(userID string) (*ProviderProperties, error)
	WriteCapabilitiesForUser(userID string, data *ProviderProperties) error
	DeleteCapabilitiesForUser(userID string) error
}

// Parameters to save Grafana configuration
type GrafanaConfigParams struct {
	GrafanaURL    string `json:"grafanaURL,omitempty"`
	GrafanaAPIKey string `json:"grafanaAPIKey,omitempty"`
}
