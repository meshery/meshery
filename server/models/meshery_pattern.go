package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/sql"
	"gopkg.in/yaml.v2"
)

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	ID *uuid.UUID `json:"id,omitempty"`

	Name        string `json:"name,omitempty"`
	PatternFile string `json:"pattern_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID *string `json:"user_id" gorm:"-"`

	Location    sql.Map `json:"location"`
	Visibility  string  `json:"visibility"`
	CatalogData sql.Map `json:"catalogData"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// GetPatternName takes in a stringified patternfile and extracts the name from it
func GetPatternName(stringifiedFile string) (string, error) {
	out := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(stringifiedFile), &out); err != nil {
		return "", err
	}

	// Get Name from the file
	name, ok := out["name"].(string)
	if !ok {
		return "", ErrPatternFileName
	}

	return name, nil
}
