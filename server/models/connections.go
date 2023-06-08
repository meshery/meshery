package models

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

// swagger:response ConnectionStatus
type ConnectionStatus string

const (
	DISCOVERED 		ConnectionStatus = "discovered"
	REGISTERED 		ConnectionStatus = "registered"
	CONNECTED 		ConnectionStatus = "connected"
	IGNORED 		ConnectionStatus = "ignored"
	MAINTENANCE 	ConnectionStatus = "maintenance"
	DISCONNECTED 	ConnectionStatus = "disconnected"
	DELETED 		ConnectionStatus = "deleted"
	NOT_FOUND 		ConnectionStatus = "not_found"
)

// swagger:response Connection
type Connection struct {
	ID           uuid.UUID              `json:"id,omitempty" db:"id"`
	Name         string                 `json:"name,omitempty" db:"name"`
	CredentialID uuid.UUID              `json:"credential_id,omitempty" db:"credential_id"`
	Type         string                 `json:"type,omitempty" db:"type"`
	SubType      string                 `json:"sub_type,omitempty" db:"sub_type"`
	Kind         string                 `json:"kind,omitempty" db:"kind"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	Status     	 ConnectionStatus 		`json:"status,omitempty" db:"status"`
	UserID       *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	CreatedAt    time.Time              `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt    sql.NullTime           `json:"deleted_at,omitempty" db:"deleted_at"`
}

// swagger:response ConnectionPage
type ConnectionPage struct {
	Connection []Connection `json:"connection"`
	TotalCount int          `json:"total_count"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
}
