// Package handlers Meshery API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
//
//     Schemes: http
//     BasePath: /
//     Version: 0.4.27
//     License: Apache-2.0 http://www.apache.org/licenses/LICENSE-2.0.txt
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - token: []
//
//     SecurityDefinitions:
//     token:
//          type: apiKey
//          name: token
//          in: cookie
//
//
// swagger:meta
package handlers

import (
	"github.com/go-openapi/strfmt"
	"github.com/layer5io/meshery/models"
	v1 "k8s.io/api/core/v1"
)

// Returns all meshery patterns
// swagger:response mesheryPatternsResponseWrapper
type mesheryPatternsResponseWrapper struct {
	// in: body
	Body models.PatternsAPIResponse
}

// Returns a single meshery pattern
// swagger:response mesheryPatternResponseWrapper
type mesheryPatternResponseWrapper struct {
	// in: body
	Body models.MesheryPattern
}

// swagger:response noContentWrapper
type noContentWrapper struct {
}

// swagger:parameters idGetMesheryPattern idDeleteMesheryPattern idGetSinglePerformanceProfile idDeletePerformanceProfile idGETProfileResults
type IDParameterWrapper struct {
	// id for a specific
	// in: path
	// required: true
	ID strfmt.UUID `json:"id"`
}

// Returns all the performance profiles
// swagger:response performanceProfilesResponseWrapper
type performanceProfilesResponseWrapper struct {
	// in: body
	Body models.PerformanceProfilesAPIResponse
}

// Returns all performance results
// swagger:response performanceResultsResponseWrapper
type performanceResultsResponseWrapper struct {
	// in: body
	Body models.PerformanceResultsAPIResponse
}

// Returns a single performance profile
// swagger:response performanceProfileResponseWrapper
type performanceProfileResponseWrapper struct {
	// in: body
	Body models.PerformanceProfile
}

// Save a performance profile
// swagger:parameters idSavePerformanceProfile
type performanceProfileParameterWrapper struct {
	// in: body
	// name of performance profile
	// required: true
	Name string `json:"name,omitempty"`
	// array of load generators
	// required: true
	LoadGenerators []string `json:"load_generators,omitempty" gorm:"type:text[]"`
	// array of urls of performance results
	// required: true
	Endpoints []string `json:"endpoints,omitempty" gorm:"type:text[]"`
	// service mesh for performance tests
	// required: true
	ServiceMesh string `json:"service_mesh,omitempty"`
	// number of concurrent requests
	// required: true
	ConcurrentRequest int `json:"concurrent_request,omitempty"`
	// qps in integer
	// required: true
	QPS int `json:"qps,omitempty"`
	// duration of tests e.g. 30s
	// required: true
	Duration string `json:"duration,omitempty"`
}

// Run a performance test with params
// swagger:parameters idRunPerformanceTest
type performanceTestParameterWrapper struct {
	// in: query

	// test-id of pre-existing test
	TestID strfmt.UUID `json:"uuid,omitempty"`
	// name of performance test
	Name string `json:"name"`
	// load generator for performance test
	// required: true
	LoadGenerator string `json:"loadGenerator"`
	// url for test
	// required: true
	URL string `json:"url"`
	// service mesh for performance test
	// required: true
	ServiceMesh string `json:"mesh"`
	// concurrent request in number
	// required: true
	ConcurrentRequest int `json:"c"`
	// qps in number
	// required: true
	QPS int `json:"qps"`
	// time in integer e.g. 30
	// required: true
	Time int `json:"t"`
	// duration e.g. s for second
	Duration string `json:"dur"`
}

// swagger:parameters idPostGrafanaConfig
type grafanaConfigParamsWrapper struct {
	// in: body
	// required: true
	Body *models.GrafanaConfigParams
}

// Returns Grafana configs
// swagger:response grafanaConfigResponseWrapper
type grafanaConfigResponseWrapper struct {
	// in: body
	Body *models.Grafana
}

// Parameter to search a matching grafana board
// swagger:parameters idGetGrafanaBoards
type grafanaBoardsParamsWrapper struct {
	// in: query
	DashboardSearch string `json:"dashboardSearch"`
}

// Returns Grafana boards and panels
// swagger:response grafanaBoardsResponseWrapper
type grafanaBoardsResponseWrapper struct {
	// in: body
	Body []*models.GrafanaBoard
}

// Returns a map for v1 services
// swagger:response v1ServicesMapResponseWrapper
type v1ServicesMapResponseWrapper struct {
	// in: body
	Body map[string][]v1.Service
}

// Parameters for persisting or deleting prometheus url
// swagger:parameters idPostPrometheusConfig
type prometheusConfigParamsWrapper struct {
	// in: body
	PrometheusURL string `json:"prometheusURL,omitempty"`
}

// Returns prometheus configuration
// swagger:response prometheusConfigResponseWrapper
type prometheusConfigResponseWrapper struct {
	// in: body
	Body *models.Prometheus
}

// Response for prometheus board import
// swagger:response prometheusBoardImportRespWrapper
type prometheusBoardImportRespWrapper struct {
	// in: body
	Body *models.GrafanaBoard
}

// Returns Prometheus static board
// swagger:response prometheusStaticBoardRespWrapper
type prometheusStaticBoardRespWrapper struct {
	// in: body
	Body map[string]*models.GrafanaBoard
}

// Save selected Prometheus boards
// swagger:parameters idPostPrometheusBoard
type prometheusBoardParamsWrapper struct {
	// in: body
	// required: true
	Body []*models.SelectedGrafanaConfig
}
