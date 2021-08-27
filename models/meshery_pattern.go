package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/internal/sql"
	"gopkg.in/yaml.v2"
)

// MesheryPattern represents the patterns that needs to be saved
type MesheryPattern struct {
	// example: [0]
	ID *uuid.UUID `json:"id,omitempty"`
	// example: string
	Name string `json:"name,omitempty"`
	// example: string
	PatternFile string `json:"pattern_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	// example: string
	UserID *string `json:"user_id" gorm:"-"`
	// example: { "host": "github.com", "path": "service-mesh-patterns/service-mesh-patterns", "type": "remote" }
	Location sql.Map `json:"location"`
	// example: 2021-08-02T18:30:22.995Z
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	// example: 2021-08-02T18:30:22.995Z
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
