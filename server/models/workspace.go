package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// MapObject is a GORM-compatible map[string]string type that serializes as JSON text.
type MapObject map[string]string

func (m *MapObject) Scan(src interface{}) error {
	if src == nil {
		*m = MapObject{}
		return nil
	}
	var b []byte
	switch t := src.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		return fmt.Errorf("scan source was not []byte nor string but %T", src)
	}
	return json.Unmarshal(b, m)
}

func (m MapObject) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

type Workspace struct {
	ID             uuid.UUID    `json:"id,omitempty" db:"id"`
	Name           string       `json:"name,omitempty" db:"name"`
	OrganizationID *uuid.UUID   `json:"organization_id,omitempty" db:"organization_id"`
	Description    string       `json:"description,omitempty" db:"description"`
	Owner          string       `json:"owner,omitempty" db:"owner"`
	Metadata       MapObject    `json:"metadata,omitempty" db:"metadata" gorm:"type:text"`
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
