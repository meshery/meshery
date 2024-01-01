package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Organization struct {
	// Unique identifier
	ID          *uuid.UUID   `json:"id,omitempty" db:"id"`
	Name        string       `json:"name" db:"name"`
	Country     string       `json:"country,omitempty" db:"country"`
	Region      string       `json:"region,omitempty" db:"region"`
	Description string       `json:"description,omitempty" db:"description"`
	Owner       uuid.UUID    `json:"owner,omitempty" db:"owner"`
	CreatedAt   time.Time    `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type OrganizationsPage struct {
	Organizations []*Organization `json:"organizations"`
	TotalCount    int             `json:"total_count"`
	Page          uint64          `json:"page"`
	PageSize      uint64          `json:"page_size"`
}
