package models

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/meshery/schemas/models/core"
)

// PerformanceProfileParameters structs contains parameters to save a performance profile
type PerformanceProfileParameters struct {
	// name of performance profile
	Name string `json:"name,omitempty"`
	// array of load generators
	LoadGenerators []string `json:"loadGenerators,omitempty" gorm:"type:text[]"`
	// array of urls of performance results
	Endpoints []string `json:"endpoints,omitempty" gorm:"type:text[]"`
	// infrastructure for performance tests
	ServiceMesh string `json:"serviceMesh,omitempty"`
	// number of concurrent requests
	ConcurrentRequest int `json:"concurrentRequest,omitempty"`
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
	// infrastructure for performance test
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
	PageSize   uint                 `json:"pageSize"`
	TotalCount uint                 `json:"totalCount"`
	Profiles   []PerformanceProfile `json:"profiles,omitempty"`
}

// PerformanceResultsAPIResponse response retruned by performance endpoint on meshery server
type PerformanceResultsAPIResponse struct {
	Page       uint                `json:"page"`
	PageSize   uint                `json:"pageSize"`
	TotalCount uint                `json:"totalCount"`
	Results    []PerformanceResult `json:"results,omitempty"`
}

// PerformanceResult represents the result of a performance test
type PerformanceResult struct {
	MesheryID          *core.Uuid    `json:"mesheryId,omitempty"`
	Name               string        `json:"name,omitempty"`
	Mesh               string        `json:"mesh,omitempty"`
	PerformanceProfile *core.Uuid    `json:"performanceProfile,omitempty"`
	UserID             *core.Uuid    `json:"userId"`
	RunnerResults      RunnerResults `json:"runnerResults"`
	ServerMetrics      interface{}   `json:"serverMetrics"`
	ServerBoardConfig  interface{}   `json:"serverBoardConfig,omitempty"`
	TestStartTime      *time.Time    `json:"testStartTime,omitempty"`
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
