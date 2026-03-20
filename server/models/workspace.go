package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Workspace struct {
	ID             uuid.UUID    `json:"id,omitempty" db:"id"`
	Name           string       `json:"name,omitempty" db:"name"`
	OrganizationID uuid.UUID    `json:"organization_id,omitempty" db:"organization_id"`
	Description    string       `json:"description,omitempty" db:"description"`
	CreatedAt      time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type WorkspacePayload struct {
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

type WorkspacePage struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalCount int         `json:"total_count"`
	Workspaces []Workspace `json:"workspaces"`
}

// MesheryDesignPage - represents a page of meshery patterns
type MesheryDesignPage struct {
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalCount int               `json:"total_count"`
	Designs    []*MesheryPattern `json:"designs"`
}

type WorkspacesEnvironmentsMapping struct {
	ID            uuid.UUID    `json:"ID,omitempty" db:"id"`
	WorkspaceId   uuid.UUID    `json:"workspace_id,omitempty" db:"workspace_id"`
	EnvironmentId uuid.UUID    `json:"environment_id,omitempty" db:"environment_id"`
	CreatedAt     time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt     sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type WorkspacesDesignsMapping struct {
	ID          uuid.UUID    `json:"ID,omitempty" db:"id"`
	WorkspaceId uuid.UUID    `json:"workspace_id,omitempty" db:"workspace_id"`
	DesignId    uuid.UUID    `json:"design_id,omitempty" db:"design_id"`
	CreatedAt   time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}
