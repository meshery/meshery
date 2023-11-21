package connections

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"
)

// swagger:response ConnectionStatus
type ConnectionStatus string

const (
	DISCOVERED   ConnectionStatus = "discovered"
	REGISTERED   ConnectionStatus = "registered"
	CONNECTED    ConnectionStatus = "connected"
	IGNORED      ConnectionStatus = "ignored"
	MAINTENANCE  ConnectionStatus = "maintenance"
	DISCONNECTED ConnectionStatus = "disconnected"
	DELETED      ConnectionStatus = "deleted"
	NOTFOUND     ConnectionStatus = "not found"
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
	Status       ConnectionStatus       `json:"status,omitempty" db:"status"`
	UserID       *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	CreatedAt    time.Time              `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt    sql.NullTime           `json:"deleted_at,omitempty" db:"deleted_at"`
}

// swagger:response ConnectionPage
type ConnectionPage struct {
	Connections []*Connection `json:"connections"`
	TotalCount  int           `json:"total_count"`
	Page        int           `json:"page"`
	PageSize    int           `json:"page_size"`
}

type ConnectionStatusInfo struct {
	Status string `json:"status" db:"status"`
	Count  int    `json:"count" db:"count"`
}

// swagger:response ConnectionsStatusPage
type ConnectionsStatusPage struct {
	ConnectionsStatus []*ConnectionStatusInfo `json:"connections_status"`
}
