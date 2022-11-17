package models

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
)

// PerformanceProfileParameters structs contains parameters to save a performance profile
type PerformanceProfileParameters struct {
	// name of performance profile
	Name string `json:"name,omitempty"`
	// array of load generators
	LoadGenerators []string `json:"load_generators,omitempty" gorm:"type:text[]"`
	// array of urls of performance results
	Endpoints []string `json:"endpoints,omitempty" gorm:"type:text[]"`
	// service mesh for performance tests
	ServiceMesh string `json:"service_mesh,omitempty"`
	// number of concurrent requests
	ConcurrentRequest int `json:"concurrent_request,omitempty"`
	// qps in integer
	QPS int `json:"qps,omitempty"`
	// duration of tests e.g. 30s
	Duration string `json:"duration,omitempty"`
}

// PerformanceTestParameters contains parameters to run a performance test
type PerformanceTestParameters struct {
	// test-id of pre-existing test
	TestID strfmt.UUID `json:"uuid,omitempty"`
	// name of performance test
	Name string `json:"name"`
	// load generator for performance test
	LoadGenerator string `json:"loadGenerator"`
	// url for test
	URL string `json:"url"`
	// service mesh for performance test
	ServiceMesh string `json:"mesh"`
	// concurrent request in number
	ConcurrentRequest int `json:"c"`
	// qps in number
	QPS int `json:"qps"`
	// time in integer e.g. 30
	Time int `json:"t"`
	// duration e.g. s for second
	Duration string `json:"dur"`
}

// PerformanceProfilesAPIResponse response retruned by performance endpoint on meshery server
type PerformanceProfilesAPIResponse struct {
	Page       uint                 `json:"page"`
	PageSize   uint                 `json:"page_size"`
	TotalCount uint                 `json:"total_count"`
	Profiles   []PerformanceProfile `json:"profiles,omitempty"`
}

// PerformanceResultsAPIResponse response retruned by performance endpoint on meshery server
type PerformanceResultsAPIResponse struct {
	Page       uint                `json:"page"`
	PageSize   uint                `json:"page_size"`
	TotalCount uint                `json:"total_count"`
	Results    []PerformanceResult `json:"results,omitempty"`
}

// PerformanceResult represents the result of a performance test
type PerformanceResult struct {
	MesheryID          *uuid.UUID    `json:"meshery_id,omitempty"`
	Name               string        `json:"name,omitempty"`
	Mesh               string        `json:"mesh,omitempty"`
	PerformanceProfile *uuid.UUID    `json:"performance_profile,omitempty"`
	UserID             *uuid.UUID    `json:"user_id"`
	RunnerResults      RunnerResults `json:"runner_results"`
	ServerMatrics      interface{}   `json:"server_metrics"`
	ServerBoardConfig  interface{}   `json:"server_board_config,omitempty"`
	TestStartTime      *time.Time    `json:"test_start_time,omitempty"`
}

type RunnerResults struct {
	URL               string     `json:"URL"`
	LoadGenerator     string     `json:"load-generator"`
	ActualDuration    uint64     `json:"ActualDuration"`
	RequestedDuration string     `json:"RequestedDuration"`
	QPS               float64    `json:"ActualQPS"`
	StartTime         *time.Time `json:"StartTime"`
	DurationHistogram struct {
		Average     float64 `json:"Avg,omitempty"`
		Max         float64 `json:"Max,omitempty"`
		Min         float64 `json:"Min,omitempty"`
		Percentiles []struct {
			Percentile float64 `json:"Percentile,omitempty"`
			Value      float64 `json:"Value,omitempty"`
		} `json:"Percentiles,omitempty"`
	} `json:"DurationHistogram,omitempty"`
}
