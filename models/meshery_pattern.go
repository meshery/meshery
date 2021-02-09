package models

import "github.com/gofrs/uuid"

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	ID uuid.UUID `json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	PatternFile string `json:"pattern_file"`
}
