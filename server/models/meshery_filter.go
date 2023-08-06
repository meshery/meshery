package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/sql"
	"gopkg.in/yaml.v2"
)

// MesheryFilter represents the filters that needs to be saved
type MesheryFilter struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name       string `json:"name,omitempty"`
	FilterFile []byte `json:"filter_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id"`

	Location       sql.Map    `json:"location"`
	Visibility     string     `json:"visibility"`
	CatalogData    sql.Map    `json:"catalogData"`
	FilterResource string     `json:"filter_resource"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

type MesheryFilterPayload struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name       string `json:"name,omitempty"`
	FilterFile []byte `json:"filter_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id"`

	Location       sql.Map    `json:"location"`
	Visibility     string     `json:"visibility"`
	CatalogData    sql.Map    `json:"catalogData"`
	FilterResource string     `json:"filter_resource"`
	Config         string     `json:"config"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

// MesheryCatalogFilterRequestBody refers to the type of request body that PublishCatalogFilter would receive
type MesheryCatalogFilterRequestBody struct {
	ID          uuid.UUID `json:"id,omitempty"`
	CatalogData sql.Map   `json:"catalog_data,omitempty"`
}

// MesheryCatalogFilterRequestBody refers to the type of request body
// that CloneMesheryFilterHandler would receive
type MesheryCloneFilterRequestBody struct {
	Name string `json:"name,omitempty"`
}

// MesheryFilterRequestBody refers to the type of request body that
// SaveMesheryFilter would receive
type MesheryFilterRequestBody struct {
	URL        string                `json:"url,omitempty"`
	Path       string                `json:"path,omitempty"`
	Save       bool                  `json:"save,omitempty"`
	Config     string                `json:"config,omitempty"`
	FilterData *MesheryFilterPayload `json:"filter_data,omitempty"`
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
