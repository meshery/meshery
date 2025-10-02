package connections

import (
	"context"

	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
)

// swagger:response ConnectionStatus
type ConnectionStatus = schemasConnection.ConnectionStatus

type InitFunc func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error)

// TODO
// Caps lock values are left for compatibility for now,
// update later on to Pascal case everywhere
const (
	DISCOVERED   ConnectionStatus = schemasConnection.Discovered
	REGISTERED   ConnectionStatus = schemasConnection.Registered
	CONNECTED    ConnectionStatus = schemasConnection.Connected
	IGNORED      ConnectionStatus = schemasConnection.Ignored
	MAINTENANCE  ConnectionStatus = schemasConnection.Maintenance
	DISCONNECTED ConnectionStatus = schemasConnection.Disconnected
	DELETED      ConnectionStatus = schemasConnection.Deleted
	NOTFOUND     ConnectionStatus = schemasConnection.NotFound
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
type Connection = schemasConnection.Connection

var validConnectionStatusToManage = []ConnectionStatus{
	DISCOVERED, REGISTERED, CONNECTED,
	// If the connection has status as NotFound we try to discover it again as the NotFound status indicates, connection was available previously.
	NOTFOUND,
}

// Check whether the Connection should be managed.
// Connections with status as Discovered, Registered, Connected should only be managed.
// Eg: If the status is set as Maintenance or Ignore do not try to mange it, not even during greedy import of K8sConnection from KubeConfig.
func ShouldConnectionBeManaged(c Connection) bool {
	for _, validStatus := range validConnectionStatusToManage {
		if validStatus == c.Status {
			return true
		}
	}
	return false
}

// swagger:response ConnectionPage
type ConnectionPage = schemasConnection.ConnectionPage

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
