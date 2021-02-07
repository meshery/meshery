package models

import "github.com/gofrs/uuid"

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	ID uuid.UUID `json:"pattern_id,omitempty" db:"id"`

	Name        string    `json:"name,omitempty" db:"name"`
	UserID      uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	PatternFile string    `json:"pattern_file" db:"pattern_file"`
}
