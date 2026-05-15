package models

import (
	"encoding/json"
	"time"

	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/schemas/models/core"
	"gopkg.in/yaml.v2"
)

// MesheryFilter represents the filters that needs to be saved
type MesheryFilter struct {
	ID *core.Uuid `json:"id,omitempty"`

	Name       string `json:"name,omitempty"`
	FilterFile []byte `json:"filterFile"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"userId"`

	Location       sql.Map    `json:"location"`
	Visibility     string     `json:"visibility"`
	CatalogData    sql.Map    `json:"catalogData"`
	FilterResource string     `json:"filterResource"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

type MesheryFilterPayload struct {
	ID *core.Uuid `json:"id,omitempty"`

	Name       string `json:"name,omitempty"`
	FilterFile []byte `json:"filterFile"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"userId"`

	Location       sql.Map    `json:"location"`
	Visibility     string     `json:"visibility"`
	CatalogData    sql.Map    `json:"catalogData"`
	FilterResource string     `json:"filterResource"`
	Config         string     `json:"config"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

// MesheryCatalogFilterRequestBody refers to the type of request body that PublishCatalogFilter would receive
type MesheryCatalogFilterRequestBody struct {
	ID          core.Uuid `json:"id,omitempty"`
	CatalogData sql.Map   `json:"catalogData,omitempty"`
}

// MesheryCatalogFilterRequestBody refers to the type of request body
// that CloneMesheryFilterHandler would receive
type MesheryCloneFilterRequestBody struct {
	Name string `json:"name,omitempty"`
}

// MesheryFilterRequestBody refers to the type of request body that
// SaveMesheryFilter would receive. Canonical wire form for the
// wrapper is `filterData` (camelCase) per the identifier-naming
// migration; the legacy snake_case `filter_data` spelling is still
// dual-accepted via UnmarshalJSON for the deprecation window because
// the existing UI (MesheryFilters/Filters.tsx) and the outbound
// remote-provider call still emit `filter_data`. Canonical wins when
// both are present.
//
// Note: the outbound wrapper in remote_provider.SaveMesheryFilter
// continues to emit `filter_data` because meshery-cloud's
// MesheryFilterRequestBody (server/handlers/meshery_filters.go:29)
// has not yet been migrated to accept `filterData`. Flipping the
// outbound without cloud dual-accept would silently drop the payload
// (encoding/json does not match `filterData` to a `filter_data`
// tag even case-insensitively — the underscore is significant). That
// cross-repo coordination is tracked as part of the per-resource
// Phase 3 migration for filter.
type MesheryFilterRequestBody struct {
	URL        string                `json:"url,omitempty"`
	Path       string                `json:"path,omitempty"`
	Save       bool                  `json:"save,omitempty"`
	Config     string                `json:"config,omitempty"`
	FilterData *MesheryFilterPayload `json:"filterData,omitempty"`
}

// UnmarshalJSON dual-accepts the canonical camelCase `filterData`
// and the legacy snake_case `filter_data` wrapper keys for
// FilterData. Canonical wins when both are present. Other fields
// unmarshal via stdlib default rules through the embedded-alias
// pattern; FilterData is explicitly re-zeroed before the precedence
// switch so a reused receiver does not carry a stale pointer when
// the next payload omits both spellings.
func (p *MesheryFilterRequestBody) UnmarshalJSON(data []byte) error {
	type alias MesheryFilterRequestBody
	aux := &struct {
		*alias
		FilterDataCanonical *MesheryFilterPayload `json:"filterData,omitempty"`
		FilterDataLegacy    *MesheryFilterPayload `json:"filter_data,omitempty"`
	}{alias: (*alias)(p)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.FilterData = nil
	switch {
	case aux.FilterDataCanonical != nil:
		p.FilterData = aux.FilterDataCanonical
	case aux.FilterDataLegacy != nil:
		p.FilterData = aux.FilterDataLegacy
	}
	return nil
}

// GetFilterName takes in a stringified filterfile and extracts the name from it
func GetFilterName(stringifiedFile string) (string, error) {
	out := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(stringifiedFile), &out); err != nil {
		return "", err
	}

	// Get Name from the file
	name, ok := out["name"].(string)
	if !ok {
		return "", ErrFilterFileName
	}

	return name, nil
}
