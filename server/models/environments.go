package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type EnvironmentData struct {
	ID             uuid.UUID    `json:"id,omitempty" db:"id"`
	Name           string       `json:"name,omitempty" db:"name"`
	Owner          string       `json:"owner,omitempty" db:"owner"`
	OrganizationID uuid.UUID    `json:"org_id,omitempty" db:"org_id"`
	Description    string       `json:"description,omitempty" db:"description"`
	CreatedAt      time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt      sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type EnvironmentPage struct {
	Environments []Environment `json:"environments"`
	TotalCount   int           `json:"total_count"`
	Page         int           `json:"page"`
	PageSize     int           `json:"page_size"`
}
