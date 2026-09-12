package models

import (
	"time"

	"github.com/meshery/schemas/models/core"
)

// PatternResource represents a pattern resource that is provisioned
// by meshery.
//
// The wire form is the schemas v1beta3 pattern_resource.MesheryPatternResource
// contract, which spells the owner `userId` - it is what meshery-cloud both
// accepts on SaveMesheryPatternResource and emits on read. The struct itself
// stays local because it doubles as the GORM model for the built-in provider's
// table, and the schemas models carry `db:` tags that GORM does not read; the
// explicit column tag below is what keeps that table's `owner` column stable
// while the wire key follows the schema.
type PatternResource struct {
	ID        *core.Uuid `json:"id,omitempty"`
	UserID    *core.Uuid `json:"userId,omitempty" gorm:"column:owner"`
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
	OAMType   string     `json:"oamType,omitempty"`
	Deleted   bool       `json:"deleted,omitempty"`
	// History   []PatternResource `json:"history,omitempty"` // Maybe reused when audit trail arrives

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
