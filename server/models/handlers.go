package models

import (
	"net/http"

	"time"

	"github.com/layer5io/meshery/server/models/meshmodel"
	"github.com/layer5io/meshkit/utils/events"
	"github.com/vmihailenco/taskq/v3"
)

// HandlerInterface defines the methods a Handler should define
type HandlerInterface interface {
	ServerVersionHandler(w http.ResponseWriter, r *http.Request)

	ProviderMiddleware(http.Handler) http.Handler

	//Set the AuthenticationMechanism as NoAuth to skip provider authentication for certain endpoints. If the provider is enforced, then this flag will not be respected.
	//Make sure all the endpoints are behind this middleware thereby protecting them. The reason for not just skipping this middleware is:
	//1. So that we can enfore provider through this middleware whenever want for use cases where no unauthenticated endpoints should be there, at buildtime.
	//2. For adapter and other components of Meshery, they register/use endpoints without any provider authentication. Although we can have a different type of authentication built for
	//such external systems trying to communicate without provider authentication. So for different endpoints, different authentication mechanisms other than provider can be used.
	AuthMiddleware(http.Handler, AuthenticationMechanism) http.Handler
	KubernetesMiddleware(func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)) func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)
	MesheryControllersMiddleware(func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)) func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)
	SessionInjectorMiddleware(func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)) http.Handler
	GraphqlMiddleware(http.Handler) func(http.ResponseWriter, *http.Request, *Preference, *User, Provider)

	ProviderHandler(w http.ResponseWriter, r *http.Request)
	ProvidersHandler(w http.ResponseWriter, r *http.Request)
	ProviderUIHandler(w http.ResponseWriter, r *http.Request)
	ProviderCapabilityHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	ProviderComponentsHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	TokenHandler(w http.ResponseWriter, r *http.Request, provider Provider, fromMiddleWare bool)
	LoginHandler(w http.ResponseWriter, r *http.Request, provider Provider, fromMiddleWare bool)
	LogoutHandler(w http.ResponseWriter, req *http.Request, provider Provider)
	UserHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetUserByIDHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetUsers(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	K8SConfigHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request)
	KubernetesPingHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	GetAllContexts(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetContext(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteContext(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	// GetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	// SetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	LoadTestHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	LoadTestUsingSMPHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	CollectStaticMetrics(config *SubmitMetricsConfig) error
	FetchResultsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	FetchAllResultsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetResultHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetSMPServiceMeshes(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	FetchSmiResultsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	FetchSingleSmiResultHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	MeshOpsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	AdaptersHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	AvailableAdaptersHandler(w http.ResponseWriter, req *http.Request)
	EventStreamHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	AdapterPingHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request)
	GrafanaPingHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	// ScanPromGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	// ScanPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	// ScanGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	PrometheusConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request)
	PrometheusPingHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	UserPrefsHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	UserTestPreferenceHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	UserTestPreferenceStore(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	UserTestPreferenceGet(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	UserTestPreferenceDelete(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	SavePerformanceProfileHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetPerformanceProfilesHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetPerformanceProfileHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeletePerformanceProfileHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	SessionSyncHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	PatternFileHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMeshmodelCategories(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelCategoriesByName(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelModelsByName(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelModelsByCategories(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelModelsByCategoriesByModel(rw http.ResponseWriter, r *http.Request)
	ValidationHandler(rw http.ResponseWriter, r *http.Request)
	MeshModelGenerationHandler(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelModels(rw http.ResponseWriter, r *http.Request)
	RegisterMeshmodelComponents(rw http.ResponseWriter, r *http.Request)

	GetMeshmodelComponentByModel(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelComponentByModelByCategory(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelComponentByCategory(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelComponentsByNameByModelByCategory(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelComponentsByNameByCategory(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelComponentsByNameByModel(rw http.ResponseWriter, r *http.Request)
	GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request)
	GetAllMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request)
	GetAllMeshmodelRelationships(rw http.ResponseWriter, r *http.Request)
	GetMeshmodelRelationshipByName(rw http.ResponseWriter, r *http.Request)
	RegisterMeshmodelRelationships(rw http.ResponseWriter, r *http.Request)

	PatternFileRequestHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	CloneMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DownloadMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteMultiMesheryPatternsHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetCatalogMesheryPatternsHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	PublishCatalogPatternHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryPatternHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	FilterFileHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryFilterFileHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	FilterFileRequestHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetCatalogMesheryFiltersHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	PublishCatalogFilterHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	CloneMesheryFilterHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	ApplicationFileHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	ApplicationFileRequestHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryApplicationTypesHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryApplicationHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryApplicationSourceHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetMesheryApplicationFile(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteMesheryApplicationHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	ShareDesignHandler(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	ExtensionsEndpointHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	LoadExtensionFromPackage(w http.ResponseWriter, req *http.Request, provider Provider) error
	ExtensionsVersionHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)

	SaveScheduleHandler(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetSchedulesHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	GetScheduleHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteScheduleHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	ExtensionsHandler(w http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)

	SaveUserCredential(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetUserCredentials(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	UpdateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	DeleteUserCredential(w http.ResponseWriter, req *http.Request, prefObj *Preference, user *User, provider Provider)
	GetRegoPolicyForDesignFile(rw http.ResponseWriter, r *http.Request, prefObj *Preference, user *User, provider Provider)
}

// HandlerConfig holds all the config pieces needed by handler methods
type HandlerConfig struct {
	// SessionName   string
	// RefCookieName string

	// SessionStore sessions.Store

	AdapterTracker AdaptersTrackerInterface
	QueryTracker   QueryTrackerInterface

	Queue taskq.Queue

	KubeConfigFolder string

	GrafanaClient         *GrafanaClient
	GrafanaClientForQuery *GrafanaClient

	PrometheusClient         *PrometheusClient
	PrometheusClientForQuery *PrometheusClient

	// GraphQLHandler           http.Handler
	// GraphQLPlaygroundHandler http.Handler
	PlaygroundBuild        bool
	Providers              map[string]Provider
	ProviderCookieName     string
	ProviderCookieDuration time.Duration

	// to be removed
	BrokerEndpointURL *string

	PerformanceChannel       chan struct{}
	PerformanceResultChannel chan struct{}

	ConfigurationChannel *ConfigurationChannel

	DashboardK8sResourcesChan *DashboardK8sResourcesChan
	MeshModelSummaryChannel   *meshmodel.SummaryChannel

	K8scontextChannel *K8scontextChan
	EventsBuffer      *events.EventStreamer
	OperatorTracker   *OperatorTracker
}

// SubmitMetricsConfig is used to store config used for submitting metrics
type SubmitMetricsConfig struct {
	TestUUID, ResultID, PromURL string
	StartTime, EndTime          time.Time
	// TokenKey,
	TokenVal string
	Provider Provider
}

type AuthenticationMechanism int

func (a AuthenticationMechanism) String() string {
	switch a {
	case 0:
		return "no_auth"
	case 1:
		return "provider_auth"
	}
	return ""
}

const (
	NoAuth AuthenticationMechanism = iota
	ProviderAuth
)
