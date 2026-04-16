package models

import (
	"time"

	"github.com/meshery/schemas/models/core"
)

// PatternResource represents a pattern resource that is provisioned
// by meshery
type PatternResource struct {
	ID        *core.Uuid `json:"id,omitempty"`
	UserID    *core.Uuid `json:"user_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
	OAMType   string     `json:"oam_type,omitempty"`
	Deleted   bool       `json:"deleted,omitempty"`
	// History   []PatternResource `json:"history,omitempty"` // Maybe reused when audit trail arrives

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
