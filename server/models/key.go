package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

type Key struct {

	// Unique Identifier
	ID uuid.UUID `json:"id,omitempty" db:"id"`

	Owner uuid.UUID `json:"owner,omitempty" db:"owner"`

	Function string `json:"function,omitempty" db:"function"`

	Category string `json:"category,omitempty" db:"category"`

	Description string `json:"description,omitempty" db:"description"`

	Subcategory string `json:"subcategory,omitempty" db:"subcategory"`

	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`

	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`

	DeletedAt sql.NullTime `json:"deleted_at,omitempty" db:"deleted_at"`
}

type KeysPage struct {
	Page int `json:"page"`

	PageSize int `json:"page_size"`

	TotalCount int `json:"total_count"`

	Keys []*Key `json:"keys"`
}
