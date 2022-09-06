// Package handlers Meshery API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
//	Schemes: http
//	BasePath: /
//	Version: 0.4.27
//	License: Apache-2.0 http://www.apache.org/licenses/LICENSE-2.0.txt
//
//	Consumes:
//	- application/json
//	- multipart/form-data
//
//	Produces:
//	- application/json
//
//	Security:
//	- token: []
//
//	SecurityDefinitions:
//	token:
//	     type: apiKey
//	     name: token
//	     in: cookie
//
// swagger:meta
package handlers

import (
	"bytes"

	"github.com/go-openapi/strfmt"
	"github.com/layer5io/meshery/server/models"
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

// swagger:parameters idGetMesheryPattern idDeleteMesheryPattern idGetSinglePerformanceProfile idDeletePerformanceProfile idGETProfileResults idDeleteSchedules idGetSingleSchedule idDeleteMesheryApplicationFile idGetMesheryApplication idDeleteMesheryFilter idGetMesheryFilter
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

// Return all the adapters
// swagger:response systemAdaptersRespWrapper
type systemAdaptersRespWrapper struct {
	// in: body
	Body []models.Adapter
}

// swagger:parameters idDeleteAdapterConfig idGetSystemAdapters
type adapterParamsWrapper struct {
	// in: query
	Adapter string `json:"adapter"`
}

// Returns saved kubernetes config
// swagger:response k8sConfigRespWrapper
type k8sConfigRespWrapper struct {
	// in: body
	// Body *models.K8SConfig
}

// Returns kubernetes context list
// swagger:response k8sContextsRespWrapper
type k8sContextsRespWrapper struct {
	// in: body
	// Body []*models.K8SContext
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

// Returns the response of the application files
// swagger:response applicationFilesResponseWrapper
type applicationFilesResponseWrapper struct {
	// in: body
	Body *models.MesheryApplication
}

// Parameters for uploading a yaml file
// swagger:parameters idPostDeployApplicationFile idPostDeployPattern typeGetApplication typePutApplication
type applicationFileParamsWrapper struct {
	// in: formData
	//
	// swagger:file
	FormFile *bytes.Buffer `json:"Upload Yaml/Yml File"`
}

// Fetches a single Meshery Application
// swagger:response mesheryApplicationResponseWrapper
type mesheryApplicationResponseWrapper struct {
	// in: body
	Body models.MesheryApplication
}

// Returns all meshery applications
// swagger:response mesheryApplicationsResponseWrapper
type mesheryApplicationsResponseWrapper struct {
	// in: body
	Body models.ApplicationsAPIResponse
}

// Returns all the meshery adapters
// swagger:response mesheryAdaptersRespWrapper
type mesheryAdaptersRespWrapper struct {
	// in: body
	Body []*models.Adapter
}

// Parameter for meshery adapter location-url
// swagger:parameters idPostAdapterConfig
type mesheryAdapterParamsWrapper struct {
	// in: body
	MeshLocationURL string `json:"meshLocationURL"`
}

// Parameters for meshery operations
// swagger:parameters idPostAdapterOperation
type adapterOpsParamsWrapper struct {
	Adapter    string `json:"adapter"`
	Query      string `json:"query"`
	CustomBody string `json:"customBody"`
	Namespace  string `json:"namespace"`
	Delete     string `json:"deleteOp"`
}

// Returns a single meshery filter
// swagger:response mesheryFilterResponseWrapper
type mesheryFilterResponseWrapper struct {
	// in: body
	Body models.MesheryFilter
}

// Returns all meshery filters
// swagger:response mesheryFiltersResponseWrapper
type mesheryFiltersResponseWrapper struct {
	// in: body
	Body models.FiltersAPIResponse
}

// Returns the response of the Filter files
// swagger:response FilterFilesResponseWrapper
type filterFilesResponseWrapper struct {
	// in: body
	Body *models.MesheryFilter
}
