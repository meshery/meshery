package models

import (
	"encoding/gob"
	"time"
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

