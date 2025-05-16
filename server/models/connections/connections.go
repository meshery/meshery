package connections

import (
	"context"
	"database/sql"
	"time"

	"github.com/layer5io/meshkit/models/events"
	"github.com/spf13/viper"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models/environments"
	"github.com/layer5io/meshkit/logger"
)

// swagger:response ConnectionStatus
type ConnectionStatus string

type InitFunc func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error)

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

type ConnectionRegisterPayload struct {
	EventType string
	// It is different from connection id, this is used to track the registration process for the connection.
	// Connection ID is generated after the registration process is completed.
	ID    uuid.UUID
	Model string
	// The concrete type depends on the type of connection and the corresponding connection definition.
	Connection struct {
		ConnMetadata interface{}
		CredMetadata interface{}
	}
}

type PromConn struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type PromCred struct {
	Name string `json:"name,omitempty"`
	// If Basic then it should be formatted as username:password
	APIKeyOrBasicAuth string `json:"secret,omitempty"`
}

type GrafanaConn struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

type GrafanaCred struct {
	Name string `json:"name,omitempty"`
	// If Basic then it should be formatted as username:password
	APIKeyOrBasicAuth string `json:"secret,omitempty"`
}

// swagger:response Connection
type Connection struct {
	ID           uuid.UUID                      `json:"id,omitempty" db:"id"`
	Name         string                         `json:"name,omitempty" db:"name"`
	CredentialID uuid.UUID                      `json:"credential_id,omitempty" db:"credential_id"`
	Type         string                         `json:"type,omitempty" db:"type"`
	SubType      string                         `json:"sub_type,omitempty" db:"sub_type"`
	Kind         string                         `json:"kind,omitempty" db:"kind"`
	Metadata     utils.JSONMap                  `json:"metadata,omitempty" db:"metadata" gorm:"type:JSONB"`
	Status       ConnectionStatus               `json:"status,omitempty" db:"status"`
	UserID       *uuid.UUID                     `json:"user_id,omitempty" db:"user_id"`
	CreatedAt    time.Time                      `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time                      `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt    sql.NullTime                   `json:"deleted_at,omitempty" db:"deleted_at"`
	Environments []environments.EnvironmentData `json:"environments,omitempty" db:"environments" gorm:"-"`
}

var validConnectionStatusToManage = []ConnectionStatus{
	DISCOVERED, REGISTERED, CONNECTED,
	// If the connection has status as NotFound we try to discover it again as the NotFound status indicates, connection was available previously.
	NOTFOUND,
}

// Check whether the Connection should be managed.
// Connections with status as Discovered, Registered, Connected should only be managed.
// Eg: If the status is set as Maintenance or Ignore do not try to mange it, not even during greedy import of K8sConnection from KubeConfig.
func (c *Connection) ShouldConnectionBeManaged() bool {
	for _, validStatus := range validConnectionStatusToManage {
		if validStatus == c.Status {
			return true
		}
	}
	return false
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

type ConnectionPayload struct {
	ID                         uuid.UUID              `json:"id,omitempty"`
	Kind                       string                 `json:"kind,omitempty"`
	SubType                    string                 `json:"sub_type,omitempty"`
	Type                       string                 `json:"type,omitempty"`
	MetaData                   map[string]interface{} `json:"metadata,omitempty"`
	Status                     ConnectionStatus       `json:"status,omitempty"`
	CredentialSecret           map[string]interface{} `json:"credential_secret,omitempty"`
	Name                       string                 `json:"name,omitempty"`
	CredentialID               *uuid.UUID             `json:"credential_id,omitempty"`
	Model                      string                 `json:"model,omitempty"`
	SkipCredentialVerification bool                   `json:"skip_credential_verification"`
}

func BuildMesheryConnectionPayload(serverURL string, credential map[string]interface{}) *ConnectionPayload {
	metadata := map[string]interface{}{
		"server_id":        viper.GetString("INSTANCE_ID"),
		"server_version":   viper.GetString("BUILD"),
		"server_build_sha": viper.GetString("COMMITSHA"),
		"server_location":  serverURL,
	}
	return &ConnectionPayload{
		Kind:             "meshery",
		Type:             "platform",
		SubType:          "management",
		MetaData:         metadata,
		Status:           CONNECTED,
		CredentialSecret: credential,
	}
}
