package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// LoadTestOptions represents the load test options
type LoadTestOptions struct {
	Name string
	URL  string

	HTTPQPS float64

	HTTPNumThreads int

	IsInsecure bool
	Duration   time.Duration

	Cert, Key, CACert string

	AllowInitialErrors bool

	IsGRPC           bool
	GRPCStreamsCount int
	GRPCDoHealth     bool
	GRPCHealthSvc    string
	GRPCDoPing       bool
	GRPCPingDelay    time.Duration
}

type LoadTestStatus string

const (
	LoadTestError   LoadTestStatus = "error"
	LoadTestInfo    LoadTestStatus = "info"
	LoadTestSuccess LoadTestStatus = "success"
)

type LoadTestResponse struct {
	Status  LoadTestStatus `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Result  *MesheryResult `json:"result,omitempty"`
}

// MesheryResult - represents the results from Meshery test run to be shipped
type MesheryResult struct {
	ID     uuid.UUID              `json:"meshery_id,omitempty"`
	Name   string                 `json:"name,omitempty"`
	Mesh   string                 `json:"mesh,omitempty"`
	Result map[string]interface{} `json:"runner_results,omitempty"`

	ServerMetrics     interface{} `json:"server_metrics,omitempty"`
	ServerBoardConfig interface{} `json:"server_board_config,omitempty"`
}
