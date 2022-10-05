package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
)

// Router represents Meshery router
type Router struct {
	S    *mux.Router
	port int
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int, g http.Handler, gp http.Handler) *Router {
	gMux := mux.NewRouter()

	gMux.Handle("/api/system/graphql/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MesheryControllersMiddleware(h.GraphqlMiddleware(g))))))).Methods("GET", "POST")
	gMux.Handle("/api/system/graphql/playground", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MesheryControllersMiddleware(h.GraphqlMiddleware(gp))))))).Methods("GET", "POST")

	gMux.HandleFunc("/api/system/version", h.ServerVersionHandler).
		Methods("GET")
	gMux.Handle("/api/extension/version", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExtensionsVersionHandler)))).
		Methods("GET")

	gMux.HandleFunc("/api/provider", h.ProviderHandler)
	gMux.HandleFunc("/api/providers", h.ProvidersHandler).
		Methods("GET")
	gMux.PathPrefix("/api/provider/extension").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MesheryControllersMiddleware(h.ProviderComponentsHandler)))))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")
	gMux.Handle("/api/provider/capabilities", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProviderCapabilityHandler)))).
		Methods("GET")
	gMux.PathPrefix("/provider").
		Handler(http.HandlerFunc(h.ProviderUIHandler)).
		Methods("GET")
	gMux.HandleFunc("/auth/login", h.ProviderUIHandler).
		Methods("GET")
	// gMux.PathPrefix("/provider/").
	// 	Handler(http.StripPrefix("/provider/", http.FileServer(http.Dir("../provider-ui/out/")))).
	// 	Methods("GET")

	gMux.Handle("/api/system/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MesheryControllersMiddleware(h.SessionSyncHandler)))))).
		Methods("GET")

	gMux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserPrefsHandler)))).
		Methods("GET", "POST")

	gMux.Handle("/api/user/prefs/perf", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserTestPreferenceHandler)))).
		Methods("GET", "POST", "DELETE")

	gMux.Handle("/api/system/kubernetes", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler)))).
		Methods("POST")
	gMux.Handle("/api/system/kubernetes/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler)))).
		Methods("GET")

	gMux.Handle("/api/system/kubernetes/contexts", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig)))).
		Methods("POST")
	gMux.Handle("/api/system/kubernetes/contexts", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetAllContexts)))).
		Methods("GET")
	gMux.Handle("/api/system/kubernetes/contexts/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetContext)))).
		Methods("GET")
	gMux.Handle("/api/system/kubernetes/contexts/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteContext)))).
		Methods("DELETE")

	gMux.Handle("/api/perf/profile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/profile/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler)))).
		Methods("GET")
	gMux.Handle("/api/perf/profile/result/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler)))).
		Methods("GET")
	gMux.Handle("/api/mesh", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSMPServiceMeshes)))).
		Methods("GET")

	gMux.Handle("/api/smi/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSmiResultsHandler)))).
		Methods("GET")
	gMux.Handle("/api/smi/results/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSingleSmiResultHandler)))).
		Methods("GET")

	gMux.Handle("/api/system/adapter/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshAdapterConfigHandler)))))
	gMux.Handle("/api/system/adapter/operation", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshOpsHandler))))).
		Methods("POST")
	gMux.Handle("/api/system/adapters", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdaptersHandler))))
	gMux.Handle("/api/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EventStreamHandler)))).
		Methods("GET")

	gMux.Handle("/api/telemetry/metrics/grafana/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler)))).
		Methods("GET", "POST", "DELETE")
	gMux.Handle("/api/telemetry/metrics/grafana/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/telemetry/metrics/grafana/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler)))).
		Methods("GET")
	gMux.Handle("/api/grafana/query_range", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GrafanaQueryRangeHandler)))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/grafana/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler)))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/grafana/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ScanGrafanaHandler)))))

	gMux.Handle("/api/telemetry/metrics/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler)))).
		Methods("GET", "POST", "DELETE")
	gMux.Handle("/api/telemetry/metrics/board_import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler)))).
		Methods("POST")
	gMux.Handle("/api/telemetry/metrics/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler)))).
		Methods("GET")
	gMux.Handle("/api/prometheus/query_range", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.PrometheusQueryRangeHandler)))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler)))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/static-board", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler)))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler)))).
		Methods("POST")
	gMux.Handle("/api/system/meshsync/prometheus", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ScanPrometheusHandler)))))

	gMux.Handle("/api/system/meshsync/grafana", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ScanPromGrafanaHandler)))))

	gMux.Handle("/api/pattern/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.PatternFileHandler))))).
		Methods("POST", "DELETE")
	gMux.HandleFunc("/api/pattern/collab", h.PatternCollabHandler)
	gMux.Handle("/api/pattern", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PatternFileRequestHandler)))).
		Methods("POST", "GET")
	gMux.Handle("/api/pattern/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryPatternsHandler)))).
		Methods("GET")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternHandler)))).
		Methods("GET")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryPatternHandler)))).
		Methods("DELETE")
	gMux.Handle("/api/pattern/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryPatternHandler)))).
		Methods("POST")
	gMux.Handle("/api/patterns/delete", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMultiMesheryPatternsHandler)))).
		Methods("POST")

	gMux.HandleFunc("/api/oam/{type}", h.OAMRegisterHandler).Methods("GET", "POST")
	gMux.HandleFunc("/api/oam/{type}/{name}", h.OAMComponentDetailsHandler).Methods("GET")
	gMux.HandleFunc("/api/oam/{type}/{name}/{id}", h.OAMComponentDetailByIDHandler).Methods("GET")
	gMux.HandleFunc("/api/meshmodel/validate", h.ValidationHandler).Methods("POST")
	gMux.HandleFunc("/api/meshmodel/component/generate", h.ComponentGenerationHandler).Methods("POST")
	gMux.HandleFunc("/api/meshmodel/components", h.GetAllMeshmodelComponents).Methods("GET")

	gMux.HandleFunc("/api/components", h.GetAllComponents).Methods("GET")
	gMux.HandleFunc("/api/components/types", h.ComponentTypesHandler).Methods("GET")
	gMux.HandleFunc("/api/components/{type}", h.ComponentsForTypeHandler).Methods("GET")
	gMux.HandleFunc("/api/components/{type}/versions", h.ComponentVersionsHandler).Methods("GET")
	gMux.HandleFunc("/api/components/{type}/{name}", h.ComponentsByNameHandler).Methods("GET")

	gMux.Handle("/api/filter/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.FilterFileHandler))))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/filter", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FilterFileRequestHandler)))).
		Methods("POST", "GET")
	gMux.Handle("/api/filter/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryFiltersHandler)))).
		Methods("GET")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterHandler)))).
		Methods("GET")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryFilterHandler)))).
		Methods("DELETE")
	gMux.Handle("/api/filter/file/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterFileHandler)))).
		Methods("GET")
	gMux.Handle("/api/filter/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryFilterHandler)))).
		Methods("POST")

	gMux.Handle("/api/application/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ApplicationFileHandler))))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/application", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ApplicationFileRequestHandler)))).
		Methods("GET")
	gMux.Handle("/api/application/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ApplicationFileRequestHandler)))).
		Methods("POST", "PUT")
	gMux.Handle("/api/application/types", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryApplicationTypesHandler)))).
		Methods("GET")
	gMux.Handle("/api/application/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryApplicationHandler)))).
		Methods("GET")
	gMux.Handle("/api/application/download/{id}/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryApplicationSourceHandler)))).
		Methods("GET")
	gMux.Handle("/api/application/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryApplicationHandler)))).
		Methods("DELETE")

	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfilesHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfileHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeletePerformanceProfileHandler)))).
		Methods("DELETE")
	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SavePerformanceProfileHandler)))).
		Methods("POST")
	gMux.Handle("/api/user/performance/profiles/{id}/run", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.LoadTestHandler))))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler)))).
		Methods("GET")

	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSchedulesHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetScheduleHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteScheduleHandler)))).
		Methods("DELETE")
	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveScheduleHandler)))).
		Methods("POST")

	gMux.PathPrefix("/api/extensions").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExtensionsHandler)))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")

	//gMux.PathPrefix("/api/system/graphql").Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GraphqlSystemHandler)))).Methods("GET", "POST")

	gMux.Handle("/user/logout", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			logrus.Debug("Inside not OK")
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LogoutHandler(w, req, provider)
	})))
	gMux.Handle("/user/login", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LoginHandler(w, req, provider, false)
	})))
	gMux.Handle("/api/user/token", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.TokenHandler(w, req, provider, false)
	}))).Methods("POST", "GET")
	gMux.Handle("/api/token", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(
		func(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
			provider.ExtractToken(w, req)
		})))).Methods("GET")

	// TODO: have to change this too
	gMux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo/meshery-logo.svg")
	}))

	// Swagger Interactive Playground
	swaggerOpts := middleware.SwaggerUIOpts{SpecURL: "./swagger.yaml"}
	swaggerSh := middleware.SwaggerUI(swaggerOpts, nil)
	gMux.Handle("/swagger.yaml", http.FileServer(http.Dir("../helpers/")))
	gMux.Handle("/docs", swaggerSh)

	gMux.PathPrefix("/").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.ServeUI(w, r, "", "../../ui/out/")
		})))).
		Methods("GET")

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
