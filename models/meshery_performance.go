package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// PerformanceResult represents the result of a performance test
type PerformanceResult struct {
	MesheryID          *uuid.UUID    `json:"meshery_id,omitempty"`
	Name               string        `json:"name,omitempty"`
	Mesh               string        `json:"mesh,omitempty"`
	PerformanceProfile *uuid.UUID    `json:"performance_profile,omitempty"`
	TestID             string        `json:"test_id"`
	RunnerResults      RunnerResults `json:"runner_results"`
	ServerMatrics      interface{}   `json:"server_metrics"`
	ServerBoardConfig  interface{}   `json:"server_board_config,omitempty"`
	TestStartTime      *time.Time    `json:"test_start_time,omitempty"`
}

type RunnerResults struct {
	Duration          string     `json:"RequestedDuration"`
	Qps               float64    `json:"ActualQPS"`
	StartTime         *time.Time `json:"StartTime"`
	DurationHistogram struct {
		Percentiles []struct {
			Percentile float64 `json:"Percentile,omitempty"`
			Value      float64 `json:"Value,omitempty"`
		} `json:"Percentiles,omitempty"`
	} `json:"DurationHistogram,omitempty"`
}
