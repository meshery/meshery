package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
	schemascoreworkspace "github.com/meshery/schemas/models/v1beta1/workspace"
)

// Workspace is the local DB-compatible representation of a workspace.
// It mirrors workspace.Workspace from the schemas package but omits the
// Metadata field (type MapObject / map[string]string) which GORM cannot
// handle natively.  All GORM operations use this type; the schema type is
// used for API serialisation.
type Workspace struct {
	ID             uuid.UUID    `json:"id,omitempty" db:"id"`
	Name           string       `json:"name,omitempty" db:"name"`
	OrganizationID uuid.UUID    `json:"organization_id,omitempty" db:"organization_id"`
	Description    string       `json:"description,omitempty" db:"description"`
	Owner          string       `json:"owner,omitempty" db:"owner"`
	CreatedAt      time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

// ToSchemaWorkspace converts the local DB model to the schema workspace type.
func (w *Workspace) ToSchemaWorkspace() *schemascoreworkspace.Workspace {
	orgID := w.OrganizationID
	return &schemascoreworkspace.Workspace{
		ID:             w.ID,
		Name:           w.Name,
		OrganizationID: &orgID,
		Description:    w.Description,
		Owner:          w.Owner,
		CreatedAt:      w.CreatedAt,
		UpdatedAt:      w.UpdatedAt,
	}
}

// WorkspaceFromSchema converts a schema workspace to the local DB model.
func WorkspaceFromSchema(ws *schemascoreworkspace.Workspace) *Workspace {
	if ws == nil {
		return nil
	}
	w := &Workspace{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		Owner:       ws.Owner,
		CreatedAt:   ws.CreatedAt,
		UpdatedAt:   ws.UpdatedAt,
	}
	if ws.OrganizationID != nil {
		w.OrganizationID = *ws.OrganizationID
	}
	return w
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
