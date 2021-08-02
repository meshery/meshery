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
	SMP "github.com/layer5io/service-mesh-performance/spec"
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

// swagger:parameters idGetMesheryPattern idDeleteMesheryPattern idGetSinglePerformanceProfile idDeletePerformanceProfile idGETProfileResults idDeleteSchedules idGetSingleSchedule
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
	Body *models.PerformanceProfileParameters
}

// Run a performance test with params
// swagger:parameters idRunPerformanceTest
type performanceTestParameterWrapper struct {
	// in: query
	Body *models.PerformanceTestParameters
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

// Returns User Load Test Preferencee
// swagger:response userLoadTestPrefsRespWrapper
type userLoadTestPrefsRespWrapper struct {
	// in: body
	Body *models.Preference
}

// Updates Anonymous stats
// swagger:parameters idPostAnonymousStats
type anonymousStatsParamsWrapper struct {
	// in: body
	Body                *models.PreferenceParams
	LoadTestPreferences *models.LoadTestPreferences
}

// Returns load test preferences
// swagger:response loadTestPreferencesWrapper
type loadTestPreferencesRespWrapper struct {
	// in: body
	Body SMP.PerformanceTestConfig
}

// Parameters Persists load test preferences
// swagger:parameters idPostLoadPreferences
type loadTestPreferencesParamsWrapper struct {
	// in: body
	Body SMP.PerformanceTestConfig
}

// Parameter
// swagger:parameters idDeleteLoadPreferences idGetLoadPreferences
type UUIDParamsWrapper struct {
	// in: query
	UUID strfmt.UUID `json:"uuid"`
}

// Parameters to run performance tests
// swagger:parameters idRunPerfTest
type perfTestParamsWrapper struct {
	// in: query
	Query *models.PerformanceTestParameters
	// in: body
	Body *SMP.PerformanceTestConfig
}

// Returns Single test result
// swagger:response perfSingleResultRespWrapper
type perfSingleResultRespWrapper struct {
	// in: body
	Body *models.PerformanceSpec
}

// Returns Perf test preference
// swagger:response perfTestPrefsRespWrapper
type perfTestPrefsRespWrapper struct {
	// in: body
	Body *models.Preference
}

// Returns List of saved schedules
// swagger:response schedulesResponseWrapper
type schedulesResponseWrapper struct {
	// in: body
	Body models.SchedulesAPIResponse
}

// Returns a single schedules
// swagger:response singleScheduleResponseWrapper
type singleScheduleResponseWrapper struct {
	// in: body
	Body models.Schedule
}

// Parameters for updating provider choice
// swagger:parameters idChoiceProvider
type mesheryProviderParamsWrapper struct {
	// in: query
	Provider string `json:"provider"`
}

// Returns a list of available providers
// swagger:response listProvidersRespWrapper
type listProvidersRespWrapper struct {
	// in: body
	Body map[string]models.ProviderProperties
}

// Returns provider capabilities
// swaggere:response providerPropertiesRespWrapper
type providerPropertiesRespWrapper struct {
	// in: body
	Body models.ProviderProperties
}

// Returns Meshery version
// swagger:response mesheryVersionRespWrapper
type mesheryVersionRespWrapper struct {
	// in: body
	Body Version
}
