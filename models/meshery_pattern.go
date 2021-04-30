package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	PatternFile string `json:"pattern_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id" gorm:"-"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}
