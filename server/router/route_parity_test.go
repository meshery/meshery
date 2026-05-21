package router

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

type mockHandler struct{}

var _ models.HandlerInterface = (*mockHandler)(nil)

func (m *mockHandler) ProviderMiddleware(next http.Handler) http.Handler {
	return next
}

func (m *mockHandler) NoCacheMiddleware(next http.Handler) http.Handler {
	return next
}

func (m *mockHandler) AuthMiddleware(next http.Handler, _ models.AuthenticationMechanism) http.Handler {
	return next
}

func (m *mockHandler) KubernetesMiddleware(f func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return f
}

func (m *mockHandler) K8sFSMMiddleware(f func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return f
}

func (m *mockHandler) SessionInjectorMiddleware(f func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f(w, r, nil, nil, nil)
	})
}

func (m *mockHandler) GraphqlMiddleware(next http.Handler) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
		next.ServeHTTP(w, r)
	}
}

func (m *mockHandler) ServerVersionHandler(w http.ResponseWriter, r *http.Request) {}
func (m *mockHandler) ProviderHandler(w http.ResponseWriter, r *http.Request)       {}
func (m *mockHandler) HandleErrorHandler(w http.ResponseWriter, r *http.Request)    {}
func (m *mockHandler) ProvidersHandler(w http.ResponseWriter, r *http.Request)      {}
func (m *mockHandler) ProviderUIHandler(w http.ResponseWriter, r *http.Request)     {}
func (m *mockHandler) ProviderCapabilityHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ProviderComponentsHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) TokenHandler(w http.ResponseWriter, r *http.Request, _ models.Provider, _ bool) {
}
func (m *mockHandler) LoginHandler(w http.ResponseWriter, r *http.Request, _ models.Provider, _ bool) {
}
func (m *mockHandler) LogoutHandler(w http.ResponseWriter, _ *http.Request, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetUserByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetUsers(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetUsersKeys(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) K8SConfigHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetContextsFromK8SConfig(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) KubernetesPingHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) K8sRegistrationHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetAllContexts(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetContext(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteContext(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) LoadTestHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) LoadTestUsingSMPHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) CollectStaticMetrics(_ *models.SubmitMetricsConfig) error {
	return nil
}
func (m *mockHandler) FetchResultsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) FetchAllResultsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetResultHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetSMPServiceMeshes(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetSystemDatabase(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ResetSystemDatabase(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ServerEventConfigurationHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ServerEventConfigurationGet(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ServerEventConfigurationSet(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) FetchSmiResultsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) FetchSingleSmiResultHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) MeshAdapterConfigHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) MeshOpsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AdaptersHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AvailableAdaptersHandler(w http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) EventStreamHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ClientEventHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AdapterPingHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DownloadHandler(w http.ResponseWriter, r *http.Request) {}
func (m *mockHandler) ViewHandler(w http.ResponseWriter, r *http.Request)     {}
func (m *mockHandler) GetAllEvents(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetEventTypes(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UpdateEventStatus(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) BulkUpdateEventStatus(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteEvent(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) BulkDeleteEvent(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaConfigHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaBoardsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaQueryHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaQueryRangeHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaPingHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PrometheusConfigHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PrometheusQueryHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PrometheusQueryRangeHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PrometheusPingHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PrometheusStaticBoardHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserPrefsHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserTestPreferenceHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserTestPreferenceStore(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserTestPreferenceGet(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UserTestPreferenceDelete(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SavePerformanceProfileHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetPerformanceProfilesHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetPerformanceProfileHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeletePerformanceProfileHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SessionSyncHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PatternFileHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMeshmodelCategories(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelCategoriesByName(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelModelsByName(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelModelsByCategories(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelModelsByCategoriesByModel(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) ValidationHandler(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) MeshModelGenerationHandler(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelModels(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) RegisterMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) UpdateEntityStatus(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMeshmodelRegistrants(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) RegisterMeshmodels(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) HandleResourceSchemas(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentByModel(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentByModelByCategory(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentByCategory(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentsByNameByModelByCategory(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentsByNameByCategory(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelComponentsByNameByModel(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetAllMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetAllMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetMeshmodelRelationshipByName(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetAllMeshmodelPolicies(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetAllMeshmodelPoliciesByName(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) RegisterMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) EvaluateRelationshipPolicy(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PatternFileRequestHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) CloneMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DownloadMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteMultiMesheryPatternsHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetCatalogMesheryPatternsHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PublishCatalogPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UnPublishCatalogPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DesignFileRequestHandlerWithSourceType(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DesignFileImportHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMesheryDesignTypesHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMesheryPatternSourceHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) FilterFileHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMesheryFilterFileHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) FilterFileRequestHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetCatalogMesheryFiltersHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) PublishCatalogFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UnPublishCatalogFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) CloneMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ShareDesignHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ShareFilterHandler(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ExtensionsEndpointHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) LoadExtensionFromPackage(w http.ResponseWriter, r *http.Request, _ models.Provider) error {
	return nil
}
func (m *mockHandler) ExtensionsVersionHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveScheduleHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetSchedulesHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetScheduleHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteScheduleHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ExtensionsHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveUserCredential(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetUserCredentials(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetUserCredentialByID(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UpdateUserCredential(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteUserCredential(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveConnection(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetConnections(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetConnectionsByKind(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetConnectionByID(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UpdateConnectionById(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteConnection(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ProcessConnectionRegistration(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) ExportModel(w http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) GetEnvironments(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetEnvironmentByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveEnvironment(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UpdateEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AddConnectionToEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) RemoveConnectionFromEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetConnectionsOfEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMeshSyncResources(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMeshSyncResourceByID(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetMeshSyncResourcesSummary(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteMeshSyncResource(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteModel(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetOrganizations(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetWorkspacesHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetWorkspaceByIdHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) SaveWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) DeleteWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) UpdateWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetEnvironmentsOfWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AddEnvironmentToWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) RemoveEnvironmentFromWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetDesignsOfWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AddDesignToWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) RemoveDesignFromWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetViewsOfWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AddViewToWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) RemoveViewFromWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) GetTeamsOfWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) AddTeamToWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) RemoveTeamFromWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
}
func (m *mockHandler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {
}
func (m *mockHandler) ServeUI(w http.ResponseWriter, r *http.Request, _, _ string) {
}

func TestRouteParity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping route parity test in short mode")
	}

	h := &mockHandler{}
	r := NewRouter(context.Background(), h, 0, http.DefaultServeMux, http.DefaultServeMux)
	if r == nil {
		t.Fatal("NewRouter returned nil")
	}

	var paths []string
	_ = r.S.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, _ := route.GetPathTemplate()
		methods, _ := route.GetMethods()
		if path != "" {
			paths = append(paths, path)
			t.Logf("Route: %s %v", path, methods)
		}
		return nil
	})

	if len(paths) == 0 {
		t.Fatal("Walk returned zero routes")
	}
	t.Logf("Total unique route paths: %d", len(paths))
}
