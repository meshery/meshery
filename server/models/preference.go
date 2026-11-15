package models

import (
	"encoding/gob"
	"encoding/json"
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

// PreferenceParams holds the parameters used to update anonymous usage stats.
type PreferenceParams struct {
	AnonymousUsageStats  bool `json:"anonymousUsageStats"`
	AnonymousPerfResults bool `json:"anonymousPerfResults"`
}

// Preference represents the data stored in session / local DB.
//
// The wire form follows the schemas v1beta1 user.Preference construct, which is
// also what meshery-cloud reads and writes.
type Preference struct {
	MeshAdapters                      []*Adapter             `json:"meshAdapters,omitempty"`
	LoadTestPreferences               *LoadTestPreferences   `json:"loadTestPrefs,omitempty"`
	AnonymousUsageStats               bool                   `json:"anonymousUsageStats"`
	AnonymousPerfResults              bool                   `json:"anonymousPerfResults"`
	UpdatedAt                         time.Time              `json:"updatedAt,omitempty"`
	DashboardPreferences              map[string]interface{} `json:"dashboardPreferences,omitempty"`
	SelectedOrganizationID            string                 `json:"selectedOrganizationId,omitempty"`
	SelectedWorkspaceForOrganizations map[string]string      `json:"selectedWorkspaceForOrganizations,omitempty"` // map[orgID]workspaceID
	UsersExtensionPreferences         map[string]interface{} `json:"usersExtensionPreferences,omitempty"`
	RemoteProviderPreferences         map[string]interface{} `json:"remoteProviderPreferences,omitempty"`
}

// UnmarshalJSON accepts the legacy all-caps `selectedOrganizationID` spelling in
// addition to the canonical `selectedOrganizationId`.
//
// The canonical key is what schemas declares and what meshery-cloud reads, so
// the all-caps local spelling meant the selected organization never round
// tripped through a remote provider: executePrefSync PUT `selectedOrganizationID`,
// which meshery-cloud ignores, and the provider's reply carried
// `selectedOrganizationId`, which this struct ignored. Reading both keeps
// preferences already persisted under the legacy key - in the local provider's
// database, and in any remote provider not yet on the canonical spelling -
// readable after the rename.
//
// Which key wins is decided on payload key *presence*, not on whether the
// destination is already populated: UserPrefsHandler decodes the request body
// onto the Preference the session middleware already read from the persister,
// so the destination is normally non-empty before this runs. Both keys are
// therefore decoded through pointers - the shallower fields shadow the
// embedded struct's own `selectedOrganizationId` - and the destination is left
// untouched when the payload carries neither.
func (p *Preference) UnmarshalJSON(data []byte) error {
	type preferenceAlias Preference
	aliased := struct {
		*preferenceAlias
		SelectedOrganizationID       *string `json:"selectedOrganizationId,omitempty"`
		LegacySelectedOrganizationID *string `json:"selectedOrganizationID,omitempty"`
	}{preferenceAlias: (*preferenceAlias)(p)}

	if err := json.Unmarshal(data, &aliased); err != nil {
		return err
	}
	switch {
	case aliased.SelectedOrganizationID != nil:
		p.SelectedOrganizationID = *aliased.SelectedOrganizationID
	case aliased.LegacySelectedOrganizationID != nil:
		p.SelectedOrganizationID = *aliased.LegacySelectedOrganizationID
	}
	return nil
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
