package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/internal/sql"
)

type ApplicationType string

const (
	HELM_CHART     ApplicationType = "HELM_CHART"
	DOCKER_COMPOSE ApplicationType = "DOCKER_COMPOSE"
)

// MesheryApplication represents the applications that needs to be saved
type MesheryApplication struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name            string `json:"name,omitempty"`
	ApplicationFile string `json:"application_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID    *string                  `json:"user_id" gorm:"-"`
	Location  sql.Map                  `json:"location"`
	Source    MesheryApplicationSource `gorm:"foreignKey:SourceID" json:"source"`
	SourceID  *uuid.UUID
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

type MesheryApplicationSource struct {
	ID   *uuid.UUID      `json:"id,omitempty"`
	Type ApplicationType `json:"type,omitempty"`
	Data []byte
}
