package models

import (
	"database/sql"
	"time"

	isql "github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/schemas/models/core"
)

type ApplicationType string

type ApplicationTypeResponse struct {
	Type                ApplicationType `json:"applicationType"`
	SupportedExtensions []string        `json:"supportedExtensions"`
}

func GetApplicationTypes() (r []DesignTypeResponse) {
	r = append(r, DesignTypeResponse{
		Type:                HelmChart,
		SupportedExtensions: []string{".tgz"},
	},
		DesignTypeResponse{
			Type:                DockerCompose,
			SupportedExtensions: []string{".yaml", ".yml"},
		},
		DesignTypeResponse{
			Type:                K8sManifest,
			SupportedExtensions: []string{".yaml", ".yml"},
		})
	return
}

// MesheryApplication represents the applications that needs to be saved
type MesheryApplication struct {
	ID *core.Uuid `json:"id,omitempty"`

	Name            string `json:"name,omitempty"`
	ApplicationFile string `json:"applicationFile"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID        *string        `json:"userId" gorm:"-"`
	Location      isql.Map       `json:"location"`
	Type          sql.NullString `json:"type"`
	SourceContent []byte         `json:"sourceContent"`
	UpdatedAt     *time.Time     `json:"updatedAt,omitempty"`
	CreatedAt     *time.Time     `json:"createdAt,omitempty"`
}
