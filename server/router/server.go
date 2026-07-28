package router

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// Router represents Meshery router
type Router struct {
	S    *mux.Router
	port int
}

// registerRegistryRoute registers a capabilities-registry route at its
// canonical /api/registry path and at its deprecated /api/meshmodels alias,
// binding both to the same handler and methods so existing clients keep
// working during the migration window (#17904, #17965).
func registerRegistryRoute(gMux *mux.Router, subpath string, handler http.Handler, methods ...string) {
	gMux.Handle("/api/registry"+subpath, handler).Methods(methods...)
	gMux.Handle("/api/meshmodels"+subpath, handler).Methods(methods...) // Deprecated: /api/meshmodels alias
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(_ context.Context, h models.HandlerInterface, port int, g http.Handler, gp http.Handler) *Router {
	gMux := mux.NewRouter()
	gMux.Use(otelmux.Middleware("meshery-server"))

	gMux.Handle("/api/system/graphql/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GraphqlMiddleware(g)), models.ProviderAuth))).Methods("GET", "POST")
	gMux.Handle("/api/system/graphql/playground", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GraphqlMiddleware(gp)), models.ProviderAuth))).Methods("GET", "POST")

	gMux.HandleFunc("/api/system/version", h.ServerVersionHandler).
		Methods("GET")
	gMux.Handle("/api/extension/version", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExtensionsVersionHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/database", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSystemDatabase), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/database/reset", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ResetSystemDatabase), models.ProviderAuth))).
		Methods("DELETE")

	gMux.HandleFunc("/api/provider", h.ProviderHandler)
	gMux.HandleFunc("/error", h.HandleErrorHandler)
	gMux.HandleFunc("/api/providers", h.ProvidersHandler).
		Methods("GET")
	// Server-Sent Events stream of per-provider availability transitions.
	// The provider-ui chooser subscribes here so it can render the local
	// provider immediately and update each remote entry independently as
	// its capability probe settles, instead of blocking on the slowest
	// remote like the old polling-based /api/providers did.
	gMux.HandleFunc("/api/providers/stream", h.ProvidersStreamHandler).
		Methods("GET")
	gMux.Handle("/api/provider/extension/install", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.InstallExtensionHandler), models.ProviderAuth))).Methods("POST")
	gMux.Handle("/api/provider/extension/remove", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveExtensionHandler), models.ProviderAuth))).Methods("POST")

	gMux.PathPrefix("/api/provider/extension").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ProviderComponentsHandler)), models.ProviderAuth))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")
	gMux.Handle("/api/provider/capabilities", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProviderCapabilityHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.HandleFunc("/provider/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
	})
	gMux.PathPrefix("/provider/_next").
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
		}))

	gMux.PathPrefix("/provider").
		Handler(http.HandlerFunc(h.ProviderUIHandler)).
		Methods("GET")
	gMux.HandleFunc("/auth/login", h.ProviderUIHandler).
		Methods("GET")
	// gMux.PathPrefix("/provider/").
	// 	Handler(http.StripPrefix("/provider/", http.FileServer(http.Dir("../provider-ui/out/")))).
	// 	Methods("GET")

	gMux.Handle("/api/system/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.K8sFSMMiddleware(h.SessionSyncHandler))), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/users/profile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/profile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserByIDHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/users/profile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserByIDHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/identity/users", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUsers), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/orgs/{orgID}/users/keys", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUsersKeys), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/user/prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserPrefsHandler), models.ProviderAuth))).
		Methods("GET", "POST")

	gMux.Handle("/api/user/prefs/perf", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserTestPreferenceHandler), models.ProviderAuth))).
		Methods("GET", "POST", "DELETE")

	gMux.Handle("/api/system/kubernetes", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/kubernetes/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/system/kubernetes/contexts", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetContextsFromK8SConfig), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/kubernetes/register", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8sRegistrationHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/kubernetes/contexts", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetAllContexts), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/kubernetes/contexts/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetContext), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/kubernetes/contexts/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteContext), models.ProviderAuth))).
		Methods("DELETE")

	gMux.Handle("/api/perf/profile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler), models.ProviderAuth))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/profile/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/perf/profile/result/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/mesh", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSMPServiceMeshes), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/smi/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSmiResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/smi/results/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSingleSmiResultHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/fileDownload", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.DownloadHandler), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/system/fileView", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ViewHandler), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/system/adapter/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshAdapterConfigHandler)), models.ProviderAuth)))
	gMux.Handle("/api/system/adapter/operation", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshOpsHandler)), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/adapters", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdaptersHandler), models.ProviderAuth)))
	gMux.Handle("/api/system/availableAdapters", http.HandlerFunc(h.AvailableAdaptersHandler)).
		Methods("GET")
	gMux.Handle("/api/system/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ClientEventHandler), models.ProviderAuth))).
		Methods("POST")
		// This will be changed to /api/events once the UI is compeltely updated to use new events and SSE is tunred off, otherwise existing events will break.
	gMux.Handle("/api/system/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetAllEvents), models.ProviderAuth))).
		Methods("GET")
	// SSE stream of live events; replaces the subscribeEvents GraphQL subscription.
	gMux.Handle("/api/system/events/subscribe", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SubscribeEventsHandler), models.ProviderAuth))).
		Methods("GET")

	// Controller (operator/MeshSync/broker) status: SSE stream replaces the
	// subscribeMesheryControllersStatus GraphQL subscription; the three one-shot
	// GETs replace getOperatorStatus/getMeshsyncStatus/getNatsStatus.
	gMux.Handle("/api/system/controllers/status/subscribe", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SubscribeMesheryControllersStatusHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/controllers/operator/status", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.OperatorStatusHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/controllers/meshsync/status", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshsyncStatusHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/controllers/broker/status", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.BrokerStatusHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/controllers/diagnostics", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ControllerDiagnosticsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/events/types", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEventTypes), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/events/status/bulk", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.BulkUpdateEventStatus), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/system/events/status/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateEventStatus), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/system/events/bulk", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.BulkDeleteEvent), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/system/events/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteEvent), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/system/events/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ServerEventConfigurationHandler), models.ProviderAuth))).
		Methods("GET", "PUT")

	// Telemetry (v2) — clean, connection-driven Grafana dashboard browsing & rendering.
	gMux.Handle("/api/telemetry/grafana/{connectionID}/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryPingHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryBoardsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/boards/{uid}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryBoardHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/datasources", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryDatasourcesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryQueryRangeHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/query_range_batch", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryQueryRangeBatchHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/telemetry/grafana/{connectionID}/pinned", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaTelemetryPinnedBoardsHandler), models.ProviderAuth))).
		Methods("GET", "POST")

	// Telemetry (v2) — clean, connection-driven Prometheus-native metric browsing & rendering.
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryPingHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/metrics", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryMetricsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/labels", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryLabelsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/label_values", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryLabelValuesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/metadata", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryMetadataHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryQueryHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryQueryRangeHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/query_range_batch", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryQueryRangeBatchHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/telemetry/prometheus/{connectionID}/panels", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusTelemetryPanelsHandler), models.ProviderAuth))).
		Methods("GET", "POST")

	gMux.Handle("/api/pattern/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.PatternFileHandler)), models.ProviderAuth))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/pattern", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PatternFileRequestHandler), models.ProviderAuth))).
		Methods("POST", "GET")

	gMux.Handle("/api/pattern/import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DesignFileImportHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DesignFileRequestHandlerWithSourceType), models.ProviderAuth))).
		Methods("POST", "PUT")
	gMux.Handle("/api/pattern/types", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryDesignTypesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryPatternsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/catalog/unpublish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UnPublishCatalogPatternHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/pattern/catalog/publish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PublishCatalogPatternHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryPatternHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/pattern/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryPatternHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/download/{id}/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternSourceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/download/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DownloadMesheryPatternHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/patterns/delete", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMultiMesheryPatternsHandler), models.ProviderAuth))).
		Methods("POST")

	gMux.Handle("/api/schema/resource/{resourceName}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.HandleResourceSchemas), models.NoAuth))).
		Methods("GET")

	// Capabilities-registry routes (#17904, #17965). The canonical paths
	// live under /api/registry; each route is also registered at its
	// /api/meshmodels form, which is deprecated but retained as an alias of
	// the same handler for the migration window so existing clients keep
	// working. Registration order is significant for gorilla/mux matching
	// and mirrors the historical /api/meshmodels order.
	handleRegistry := func(subpath string, handler http.Handler, methods ...string) {
		registerRegistryRoute(gMux, subpath, handler, methods...)
	}

	handleRegistry("/validate", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ValidationHandler), models.NoAuth)), "POST")
	handleRegistry("/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.RegisterMeshmodelComponents), models.NoAuth)), "POST") //This should also be left with NoAuth
	gMux.Handle("/api/meshmodel/components/register", h.ProviderMiddleware((http.HandlerFunc(h.RegisterMeshmodelComponents)))).Methods("POST") //For backwards compatibility with previous registrants
	handleRegistry("/{entityType}/status", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateEntityStatus), models.NoAuth)), "POST")
	handleRegistry("/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelComponents), models.NoAuth)), "GET")

	// Connection definitions are first-class registry entities, authored per-model in a connections/ folder and registered alongside components and relationships.
	handleRegistry("/connections", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetConnectionDefinitions), models.NoAuth)), "GET")
	handleRegistry("/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RegisterConnectionDefinition), models.ProviderAuth)), "POST")
	handleRegistry("/connections/{connectionDefinitionId}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetConnectionDefinitionByID), models.NoAuth)), "GET")
	handleRegistry("/connections/{connectionDefinitionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateConnectionDefinition), models.ProviderAuth)), "PUT")
	handleRegistry("/connections/{connectionDefinitionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteConnectionDefinition), models.ProviderAuth)), "DELETE")

	handleRegistry("/categories", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelCategories), models.NoAuth)), "GET")
	handleRegistry("/models", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModels), models.NoAuth)), "GET")
	handleRegistry("/models/{model}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByName), models.NoAuth)), "GET")
	handleRegistry("/models/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteModel), models.ProviderAuth)), "DELETE")

	handleRegistry("/registrants", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelRegistrants), models.NoAuth)), "GET")

	handleRegistry("/register", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RegisterMeshmodels), models.ProviderAuth)), "POST")
	handleRegistry("/categories/{category}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelCategoriesByName), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/models", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByCategories), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/models/{model}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByCategoriesByModel), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/models/{model}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByModelByCategory), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/models/{model}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByModelByCategory), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByCategory), models.NoAuth)), "GET")
	handleRegistry("/categories/{category}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByCategory), models.NoAuth)), "GET")

	handleRegistry("/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelComponentsByName), models.NoAuth)), "GET")
	handleRegistry("/generate", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.MeshModelGenerationHandler), models.NoAuth)), "POST")
	handleRegistry("/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelRelationships), models.NoAuth)), "GET")

	handleRegistry("/models/{model}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByModel), models.NoAuth)), "GET")
	handleRegistry("/models/{model}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByModel), models.NoAuth)), "GET")

	handleRegistry("/models/{model}/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelRelationships), models.NoAuth)), "GET")
	handleRegistry("/models/{model}/relationships/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelRelationshipByName), models.NoAuth)), "GET")
	handleRegistry("/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.RegisterMeshmodelRelationships), models.NoAuth)), "POST") //This should also be left with NoAuth

	handleRegistry("/relationships/evaluate", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EvaluateRelationshipPolicy), models.ProviderAuth)), "POST")

	handleRegistry("/models/{model}/policies", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelPolicies), models.NoAuth)), "GET")
	handleRegistry("/models/{model}/policies/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelPoliciesByName), models.NoAuth)), "GET")
	// Historical /api/meshmodels form of the by-name policies route was
	// registered without the "/" separator before {name}; retained verbatim
	// so any caller of the old pattern keeps resolving. The /api/registry
	// form above uses the corrected /policies/{name} shape.
	gMux.Handle("/api/meshmodels/models/{model}/policies{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelPoliciesByName), models.NoAuth))).Methods("GET")

	handleRegistry("/export", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ExportModel), models.NoAuth)), "GET")

	gMux.Handle("/api/filter/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.FilterFileHandler)), models.ProviderAuth))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/filter", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FilterFileRequestHandler), models.ProviderAuth))).
		Methods("POST", "GET")
	gMux.Handle("/api/filter/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryFiltersHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/catalog/publish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PublishCatalogFilterHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/filter/catalog/unpublish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UnPublishCatalogFilterHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryFilterHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/filter/download/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterFileHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryFilterHandler), models.ProviderAuth))).
		Methods("POST")

	gMux.Handle("/api/content/design/share", h.ProviderMiddleware((h.AuthMiddleware(h.SessionInjectorMiddleware(h.ShareDesignHandler), models.ProviderAuth)))).
		Methods("POST")
	gMux.Handle("/api/content/filter/share", h.ProviderMiddleware((h.AuthMiddleware(h.SessionInjectorMiddleware(h.ShareFilterHandler), models.ProviderAuth)))).
		Methods("POST")

	gMux.Handle("/api/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfilesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfilesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/performance/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/performance/profiles/{performanceProfileId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfileHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{performanceProfileId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfileHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/performance/profiles/{performanceProfileId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeletePerformanceProfileHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/user/performance/profiles/{performanceProfileId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeletePerformanceProfileHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SavePerformanceProfileHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SavePerformanceProfileHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/user/performance/profiles/{id}/run", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.LoadTestHandler)), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/performance/profiles/{performanceProfileId}/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSchedulesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetScheduleHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteScheduleHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveScheduleHandler), models.ProviderAuth))).
		Methods("POST")

	gMux.Handle("/api/system/meshsync/resources", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResources), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/system/meshsync/resources/summary", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResourcesSummary), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/meshsync/resources/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMeshSyncResource), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/system/meshsync/resources/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResourceByID), models.ProviderAuth))).
		Methods("GET")

	// Handlers for server-wide Meshery Operator / MeshSync / Broker
	// configuration defaults
	gMux.Handle("/api/system/controllers/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetControllersDefaultConfig), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/controllers/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateControllersDefaultConfig), models.ProviderAuth))).
		Methods("PUT")
	// Handlers for User Credentials

	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserCredentials), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/credentials/{credentialID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserCredentialByID), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateUserCredential), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteUserCredential), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveUserCredential), models.ProviderAuth))).
		Methods("POST")

	gMux.PathPrefix("/api/extensions").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExtensionsHandler), models.ProviderAuth))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")

	gMux.Handle("/api/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironments), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironmentByIDHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteEnvironmentHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveEnvironment), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateEnvironmentHandler), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/environments/{environmentID}/connections/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddConnectionToEnvironmentHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/environments/{environmentID}/connections/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveConnectionFromEnvironmentHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/environments/{environmentID}/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnectionsOfEnvironmentHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/workspaces", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetWorkspacesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateWorkspaceHandler), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetWorkspaceByIdHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironmentsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/environments/{environmentID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddEnvironmentToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/environments/{environmentID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveEnvironmentFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}/designs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetDesignsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/designs/{designID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddDesignToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/designs/{designID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveDesignFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")

	gMux.Handle("/api/workspaces/{id}/views", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetViewsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/views/{viewID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddViewToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/views/{viewID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveViewFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")

	gMux.Handle("/api/workspaces/{id}/teams", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetTeamsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/teams/{teamID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddTeamToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/teams/{teamID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveTeamFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")

	gMux.Handle("/api/identity/orgs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetOrganizations), models.ProviderAuth))).
		Methods("GET")

	//gMux.PathPrefix("/api/system/graphql").Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GraphqlSystemHandler)))).Methods("GET", "POST")

	gMux.Handle("/user/logout", h.ProviderMiddleware(h.SessionInjectorMiddleware(func(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
		h.LogoutHandler(w, req, user, provider)
	})))
	gMux.Handle("/user/login", h.NoCacheMiddleware(h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LoginHandler(w, req, provider, false)
	}))))
	gMux.Handle("/api/user/token", h.NoCacheMiddleware(h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.TokenHandler(w, req, provider, false)
	})))).Methods("POST", "GET")
	gMux.Handle("/api/token", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(
		func(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
			provider.ExtractToken(w, req)
		}), models.ProviderAuth))).Methods("GET")

	// TODO: have to change this too
	gMux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo/meshery-logo.svg")
	}))

	// Handlers for User Credentials
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveUserCredential), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserCredentials), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateUserCredential), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteUserCredential), models.ProviderAuth))).
		Methods("DELETE")

	// Handlers for User Connections
	gMux.Handle("/api/integrations/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveConnection), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/integrations/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnections), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnectionByID), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateConnectionById), models.ProviderAuth))).
		Methods("PUT")
	// Side-effecting connection operations (e.g. switch MeshSync deployment
	// mode). Separate from PUT so resource updates don't carry cluster effects.
	gMux.Handle("/api/integrations/connections/{connectionId}/actions", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PerformConnectionAction), models.ProviderAuth))).
		Methods("POST")
	// POST drives the registration state machine; the DELETE method on the same
	// path is the deprecated body-based cancel (see ProcessConnectionRegistration).
	gMux.Handle("/api/integrations/connections/register", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProcessConnectionRegistration), models.ProviderAuth))).
		Methods("POST", "DELETE")
	// Schemas-defined cancel: the registration tracker id travels in the path
	// (operationId cancelConnectionRegister).
	gMux.Handle("/api/integrations/connections/register/{registrationId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CancelConnectionRegister), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteConnection), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/integrations/connections/{connectionId}/controllers/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnectionControllersConfig), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/connections/{connectionId}/controllers/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateConnectionControllersConfig), models.ProviderAuth))).
		Methods("PUT")

	gMux.HandleFunc("/auth/redirect", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		http.SetCookie(w, &http.Cookie{
			Name:     models.TokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
	}).
		Methods("GET")

	// Kubernetes Health Probes
	gMux.HandleFunc("/healthz/live", h.K8sHealthzHandler).Methods("GET")
	gMux.HandleFunc("/healthz/ready", h.K8sHealthzHandler).Methods("GET")

	fs := http.FileServer(http.Dir("../../ui"))
	gMux.PathPrefix("/ui/public/static/img/meshmodels").Handler(http.StripPrefix("/ui/", fs)).Methods("GET", "HEAD")

	// Serve Next.js build artifacts (/_next/data/<buildId>/<route>.json,
	// /_next/static/*) without the auth middleware. Next.js client-side
	// route transitions fetch the <route>.json blob to hydrate page props;
	// when that fetch is gated by AuthMiddleware and auth fails, the 302
	// redirect HTML cannot be parsed as JSON and the SPA navigation
	// silently breaks (e.g. clicking through to /extension/meshmap).
	// This mirrors the existing /provider/_next passthrough above.
	gMux.PathPrefix("/_next").
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "", "../../ui/out/")
		})).
		Methods("GET", "HEAD")

	gMux.PathPrefix("/").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "", "../../ui/out/")
		}), models.ProviderAuth))).
		Methods("GET", "HEAD")

	return &Router{
		S:    gMux,
		port: port,
	}
}

// Run starts the http server
func (r *Router) Run() error {
	// s := &http.Server{
	// 	Addr:           fmt.Sprintf(":%d", r.port),
	// 	Handler:        r.s,
	// 	ReadTimeout:    5 * time.Second,
	// 	WriteTimeout:   2 * time.Minute,
	// 	MaxHeaderBytes: 1 << 20,
	// 	IdleTimeout:    0, //time.Second,
	// }
	// return s.ListenAndServe()
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.S)
}
