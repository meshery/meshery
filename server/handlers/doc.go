package handlers

import (
	"bytes"

	"github.com/go-openapi/strfmt"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshery/server/models/environments"
	"github.com/meshery/meshkit/models/events"
	workspace "github.com/meshery/schemas/models/v1beta1/workspace"
	v1 "k8s.io/api/core/v1"
)

type userInfo struct {
	// in: body
	Body models.User
}

type users struct {
	// in: body
	Body models.AllUsers
}

type usersParameterWrapper struct {
	// in: query
	Page int64 `json:"page"`
	// in: query
	PageSize int64 `json:"pageSize"`
	// in: query
	Search string `json:"search"`
	// in: order
	Order string `json:"order"`
	// in: filter
	Filter string `json:"filter"`
}

type usersKeys struct {
	Body models.UserKeys
}

type mesheryPatternsResponseWrapper struct {
	// in: body
	Body models.PatternsAPIResponse
}

type mesheryPatternResponseWrapper struct {
	// in: body
	Body models.MesheryPattern
}

type noContentWrapper struct {
}

type IDParameterWrapper struct {
	// id for a specific
	// in: path
	// required: true
	ID strfmt.UUID `json:"id"`
}

type performanceProfilesResponseWrapper struct {
	// in: body
	Body models.PerformanceProfilesAPIResponse
}

type performanceResultsResponseWrapper struct {
	// in: body
	Body models.PerformanceResultsAPIResponse
}

type performanceProfileResponseWrapper struct {
	// in: body
	Body models.PerformanceProfile
}

type performanceProfileParameterWrapper struct {
	// in: body
	Body *models.PerformanceProfileParameters
}

type performanceTestParameterWrapper struct {
	// in: query
	Body *models.PerformanceTestParameters
}

type grafanaConfigParamsWrapper struct {
	// in: body
	// required: true
	Body *models.GrafanaConfigParams
}

type grafanaConfigResponseWrapper struct {
	// in: body
	Body *models.Grafana
}

type grafanaBoardsParamsWrapper struct {
	// in: query
	DashboardSearch string `json:"dashboardSearch"`
}

type grafanaBoardsResponseWrapper struct {
	// in: body
	Body []*models.GrafanaBoard
}

type v1ServicesMapResponseWrapper struct {
	// in: body
	Body map[string][]v1.Service
}

type prometheusConfigParamsWrapper struct {
	// in: body
	PrometheusURL string `json:"prometheusURL,omitempty"`
}

type prometheusConfigResponseWrapper struct {
	// in: body
	Body *models.Prometheus
}

type prometheusBoardImportRespWrapper struct {
	// in: body
	Body *models.GrafanaBoard
}

type prometheusStaticBoardRespWrapper struct {
	// in: body
	Body map[string]*models.GrafanaBoard
}

type prometheusBoardParamsWrapper struct {
	// in: body
	// required: true
	Body []*models.SelectedGrafanaConfig
}

type userLoadTestPrefsRespWrapper struct {
	// in: body
	Body *models.Preference
}

type anonymousStatsParamsWrapper struct {
	// in: body
	Body                *models.PreferenceParams
	LoadTestPreferences *models.LoadTestPreferences
}

type loadTestPreferencesRespWrapper struct {
	// in: body
	Body SMP.PerformanceTestConfig
}

type loadTestPreferencesParamsWrapper struct {
	// in: body
	Body SMP.PerformanceTestConfig
}

type UUIDParamsWrapper struct {
	// in: query
	UUID strfmt.UUID `json:"uuid"`
}

type perfTestParamsWrapper struct {
	// in: query
	Query *models.PerformanceTestParameters
	// in: body
	Body *SMP.PerformanceTestConfig
}

type perfSingleResultRespWrapper struct {
	// in: body
	Body *models.PerformanceSpec
}

type perfTestPrefsRespWrapper struct {
	// in: body
	Body *models.Preference
}

type schedulesResponseWrapper struct {
	// in: body
	Body models.SchedulesAPIResponse
}

type singleScheduleResponseWrapper struct {
	// in: body
	Body models.Schedule
}

type systemAdaptersRespWrapper struct {
	// in: body
	Body []models.Adapter
}

type adapterParamsWrapper struct {
	// in: query
	Adapter string `json:"adapter"`
}

type k8sConfigRespWrapper struct {
	// in: body
	// Body *models.K8SConfig
}

type k8sContextsRespWrapper struct {
	// in: body
	// Body []*models.K8SContext
}

type mesheryProviderParamsWrapper struct {
	// in: query
	Provider string `json:"provider"`
}

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

type mesheryVersionRespWrapper struct {
	// in: body
	Body Version
}

type applicationFilesResponseWrapper struct {
	// in: body
	Body *models.MesheryApplication
}

type applicationFileParamsWrapper struct {
	FormFile *bytes.Buffer `json:"Upload Yaml/Yml File"`
}

type mesheryApplicationResponseWrapper struct {
	// in: body
	Body models.MesheryApplication
}

type mesheryApplicationsResponseWrapper struct {
	// in: body
	Body models.ApplicationsAPIResponse
}

type mesheryAdaptersRespWrapper struct {
	// in: body
	Body []*models.Adapter
}

type mesheryAdapterParamsWrapper struct {
	// in: body
	MeshLocationURL string `json:"meshLocationURL"`
}

type adapterOpsParamsWrapper struct {
	Adapter    string `json:"adapter"`
	Query      string `json:"query"`
	CustomBody string `json:"customBody"`
	Namespace  string `json:"namespace"`
	Delete     string `json:"deleteOp"`
}

type mesheryFilterResponseWrapper struct {
	// in: body
	Body models.MesheryFilter
}

type mesheryFiltersResponseWrapper struct {
	// in: body
	Body models.FiltersAPIResponse
}

type filterFilesResponseWrapper struct {
	// in: body
	Body *models.MesheryFilter
}

type meshmodelModelsResponseWrapper struct {
	// in: body
	Body *models.MeshmodelsAPIResponse
}

type meshmodelModelsDuplicateResponseWrapper struct {
	// in: body
	Body *models.MeshmodelsDuplicateAPIResponse
}

type meshmodelComponentsResponseWrapper struct {
	// in: body
	Body *models.MeshmodelComponentsAPIResponse
}

type meshmodelComponentsDuplicateResponseWrapper struct {
	// in: body
	Body *models.MeshmodelComponentsDuplicateAPIResponse
}

type meshsyncResourcesResponseWrapper struct {
	// in: body
	Body *models.MeshSyncResourcesAPIResponse
}

type meshmodelRelationshipsResponseWrapper struct {
	// in: body
	Body *models.MeshmodelRelationshipsAPIResponse
}

type meshmodelPoliciesResponseWrapper struct {
	// in: body
	Body *models.MeshmodelPoliciesAPIResponse
}

type meshmodelCategoriesResponseWrapper struct {
	// in: body
	Body *models.MeshmodelCategoriesAPIResponse
}

type meshmodelRegistrantsResponseWrapper struct {
	// in: body
	Body *models.MeshmodelRegistrantsAPIResponse
}

type systemDatabaseResponseWrapper struct {
	// in: body
	Body *models.DatabaseSummary
}

type systemK8sContextsResponseWrapper struct {
	// in: body
	Body *models.MesheryK8sContextPage
}

type smiResultsResponseWrapper struct {
	// in: body
	Body *models.SmiResultPage
}

type mesheryApplicationTypesResponseWrapper struct {
	// in: body
	Body []models.ApplicationTypeResponse
}

type mesheryConnectionResponseWrapper struct {
	// in: body
	Body connections.Connection
}

type mesheryConnectionsResponseWrapper struct {
	// in: body
	Body *connections.ConnectionPage
}

type mesheryConnectionsStatusPage struct {
	// in: body
	Body *connections.ConnectionsStatusPage
}

type environmentResponseWrapper struct {
	// in: body
	Body *environments.EnvironmentData
}

type environmentsResponseWrapper struct {
	// in: body
	Body *environments.EnvironmentPage
}

type workspacesResponseWrapper struct {
	// in: body
	Body *workspace.WorkspacePage
}

type workspaceResponseWrapper struct {
	// in: body
	Body *workspace.Workspace
}

type workspaceDesignsMappingResponseWrapper struct {
	// in: body
	Body *workspace.WorkspacesDesignsMapping
}

type workspaceEnvironmentsMappingResponseWrapper struct {
	// in: body
	Body *workspace.WorkspacesEnvironmentsMapping
}

type eventResponseWrapper struct {
	// in: body
	Body *events.Event
}

type eventsResponseWrapper struct {
	// in: body
	Body *models.EventsResponse
}

type possibleTransitions struct {
	//in: body
	Body map[string]map[connections.ConnectionStatus][]connections.ConnectionStatus
}

type orgsResponseWrapper struct {
	// in: body
	Body *models.OrganizationsPage
}

type mesheryPatternSourceContentResponseWrapper struct {
	// in: body
	Body []byte
}

type meshsyncResourcesSummaryResponseWrapper struct {
	// in: body
	Body *models.MeshSyncResourcesSummaryAPIResponse
}

// Payload for meshery pattern file deploy handler idPostDeployPattern idDeleteDeployPattern
type mesheryPatternFileDeployPayloadWrapper struct {
	// in: body
	Body *models.MesheryPatternFileDeployPayload
}

// Blank identifier declarations keep swagger-only wrapper types reachable for linting.
var (
	_ userInfo
	_ users
	_ usersParameterWrapper
	_ usersKeys
	_ mesheryPatternsResponseWrapper
	_ mesheryPatternResponseWrapper
	_ noContentWrapper
	_ IDParameterWrapper
	_ performanceProfilesResponseWrapper
	_ performanceResultsResponseWrapper
	_ performanceProfileResponseWrapper
	_ performanceProfileParameterWrapper
	_ performanceTestParameterWrapper
	_ grafanaConfigParamsWrapper
	_ grafanaConfigResponseWrapper
	_ grafanaBoardsParamsWrapper
	_ grafanaBoardsResponseWrapper
	_ v1ServicesMapResponseWrapper
	_ prometheusConfigParamsWrapper
	_ prometheusConfigResponseWrapper
	_ prometheusBoardImportRespWrapper
	_ prometheusStaticBoardRespWrapper
	_ prometheusBoardParamsWrapper
	_ userLoadTestPrefsRespWrapper
	_ anonymousStatsParamsWrapper
	_ loadTestPreferencesRespWrapper
	_ loadTestPreferencesParamsWrapper
	_ UUIDParamsWrapper
	_ perfTestParamsWrapper
	_ perfSingleResultRespWrapper
	_ perfTestPrefsRespWrapper
	_ schedulesResponseWrapper
	_ singleScheduleResponseWrapper
	_ systemAdaptersRespWrapper
	_ adapterParamsWrapper
	_ k8sConfigRespWrapper
	_ k8sContextsRespWrapper
	_ mesheryProviderParamsWrapper
	_ listProvidersRespWrapper
	_ providerPropertiesRespWrapper
	_ mesheryVersionRespWrapper
	_ applicationFilesResponseWrapper
	_ applicationFileParamsWrapper
	_ mesheryApplicationResponseWrapper
	_ mesheryApplicationsResponseWrapper
	_ mesheryAdaptersRespWrapper
	_ mesheryAdapterParamsWrapper
	_ adapterOpsParamsWrapper
	_ mesheryFilterResponseWrapper
	_ mesheryFiltersResponseWrapper
	_ filterFilesResponseWrapper
	_ meshmodelModelsResponseWrapper
	_ meshmodelModelsDuplicateResponseWrapper
	_ meshmodelComponentsResponseWrapper
	_ meshmodelComponentsDuplicateResponseWrapper
	_ meshsyncResourcesResponseWrapper
	_ meshmodelRelationshipsResponseWrapper
	_ meshmodelPoliciesResponseWrapper
	_ meshmodelCategoriesResponseWrapper
	_ meshmodelRegistrantsResponseWrapper
	_ systemDatabaseResponseWrapper
	_ systemK8sContextsResponseWrapper
	_ smiResultsResponseWrapper
	_ mesheryApplicationTypesResponseWrapper
	_ mesheryConnectionResponseWrapper
	_ mesheryConnectionsResponseWrapper
	_ mesheryConnectionsStatusPage
	_ environmentResponseWrapper
	_ environmentsResponseWrapper
	_ workspacesResponseWrapper
	_ workspaceResponseWrapper
	_ workspaceDesignsMappingResponseWrapper
	_ workspaceEnvironmentsMappingResponseWrapper
	_ eventResponseWrapper
	_ eventsResponseWrapper
	_ possibleTransitions
	_ orgsResponseWrapper
	_ mesheryPatternSourceContentResponseWrapper
	_ meshsyncResourcesSummaryResponseWrapper
	_ mesheryPatternFileDeployPayloadWrapper
)
