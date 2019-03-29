package models

import (
	"time"
)

type LoadTestOptions struct {
	URL string

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
