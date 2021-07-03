package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PatternResource represents a pattern resource that is provisioned
// by meshery
type PatternResource struct {
	ID        *uuid.UUID `json:"id,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Name      string     `json:"name,omitempty"`
	Namespace string     `json:"namespace,omitempty"`
	Type      string     `json:"type,omitempty"`
	OAMType   string     `json:"oam_type,omitempty"`
	Deleted   bool       `json:"deleted,omitempty"`
	// History   []PatternResource `json:"history,omitempty"` // Maybe reused when audit trail arrives

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
