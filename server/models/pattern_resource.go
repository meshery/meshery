package models

import (
	"time"

	"github.com/meshery/schemas/models/core"
)

// PatternResource represents a pattern resource that is provisioned
// by meshery
type PatternResource struct {
	ID        *core.Uuid `json:"id,omitempty"`
	UserID    *core.Uuid `json:"userId,omitempty"`
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
	OAMType   string     `json:"oamType,omitempty"`
	Deleted   bool       `json:"deleted,omitempty"`
	// History   []PatternResource `json:"history,omitempty"` // Maybe reused when audit trail arrives

	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
