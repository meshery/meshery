package models

import (
	"time"
)

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

// MesheryResult - represents the results from Meshery test run to be shipped
type MesheryResult struct {
	Name   string                 `json:"name,omitempty"`
	Mesh   string                 `json:"mesh,omitempty"`
	Result map[string]interface{} `json:"runner_results,omitempty"`
}
