package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Credential struct {
	ID     uuid.UUID  `json:"id,omitempty" db:"id"`
	Name   string     `json:"name,omitempty" db:"name"`
	UserID *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Type   string     `json:"type,omitempty" db:"type"`
	// change to just interface
	Secret    map[string]interface{} `json:"secret,omitempty" db:"secret"`
	CreatedAt time.Time              `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time              `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt sql.NullTime           `json:"deleted_at,omitempty" db:"deleted_at"`
}

type CredentialsPage struct {
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
	TotalCount  int           `json:"total_count"`
	Credentials []*Credential `json:"credentials"`
}
