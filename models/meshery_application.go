package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/internal/sql"
)

type ApplicationType string

type ApplicationTypeResponse struct {
	Type                ApplicationType `json:"application_type"`
	SupportedExtensions []string        `json:"supported_extensions"`
}

func GetApplicationTypes() (r []ApplicationTypeResponse) {
	r = append(r, ApplicationTypeResponse{
		Type:                HELM_CHART,
		SupportedExtensions: []string{".tgz"},
	},
		ApplicationTypeResponse{
			Type:                DOCKER_COMPOSE,
			SupportedExtensions: []string{".yaml", ".yml"},
		},
		ApplicationTypeResponse{
			Type:                HELM_CHART,
			SupportedExtensions: []string{".yaml", ".yml"},
		})
	return
}

const (
	HELM_CHART     ApplicationType = "Helm Chart"
	DOCKER_COMPOSE ApplicationType = "Docker Compose"
	K8S_MANIFEST   ApplicationType = "Kubernetes Manifest"
)

// MesheryApplication represents the applications that needs to be saved
type MesheryApplication struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name            string `json:"name,omitempty"`
	ApplicationFile string `json:"application_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID        *string         `json:"user_id" gorm:"-"`
	Location      sql.Map         `json:"location"`
	Type          ApplicationType `json:"type"`
	SourceContent []byte
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}
