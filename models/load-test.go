package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// LoadGenerator - represents the load generator type
type LoadGenerator string

const (
	// FortioLG - represents the Fortio load generator
	FortioLG LoadGenerator = "fortio"

	// Wrk2LG - represents the wrk2 load generator
	Wrk2LG LoadGenerator = "wrk2"
)

// Name - retrieves a string value for the generator
func (l LoadGenerator) Name() string {
	return string(l)
}

// LoadTestOptions represents the load test options
type LoadTestOptions struct {
	Name string
	URL  string

	HTTPQPS float64

	HTTPNumThreads int

	IsInsecure bool
	Duration   time.Duration

	LoadGenerator LoadGenerator

	Cert, Key, CACert string

	AllowInitialErrors bool

	IsGRPC           bool
	GRPCStreamsCount int
	GRPCDoHealth     bool
	GRPCHealthSvc    string
	GRPCDoPing       bool
	GRPCPingDelay    time.Duration
}

// LoadTestStatus - used for representing load test status
type LoadTestStatus string

const (
	// LoadTestError - respresents an error status
	LoadTestError LoadTestStatus = "error"

	// LoadTestInfo - represents a info status
	LoadTestInfo LoadTestStatus = "info"

	// LoadTestSuccess - represents a success status
	LoadTestSuccess LoadTestStatus = "success"
)

// LoadTestResponse - used to bundle the response with status to the client
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
