package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/models"
)

// Router represents Meshery router
type Router struct {
	S    *mux.Router
	port int
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	gMux := mux.NewRouter()

	gMux.HandleFunc("/api/server/version", h.ServerVersionHandler).
		Methods("GET")

	gMux.HandleFunc("/api/provider", h.ProviderHandler)
	gMux.HandleFunc("/api/providers", h.ProvidersHandler).
		Methods("GET")
	gMux.PathPrefix("/api/provider/extension").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProviderComponentsHandler)))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")
	gMux.Handle("/api/provider/capabilities", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProviderCapabilityHandler)))).
		Methods("GET")
	gMux.PathPrefix("/provider").
		Handler(http.HandlerFunc(h.ProviderUIHandler)).
		Methods("GET")
	// gMux.PathPrefix("/provider/").
	// 	Handler(http.StripPrefix("/provider/", http.FileServer(http.Dir("../provider-ui/out/")))).
	// 	Methods("GET")

	gMux.Handle("/api/config/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SessionSyncHandler)))).
		Methods("GET")

	gMux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler)))).
		Methods("GET")
	gMux.Handle("/api/user/stats", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AnonymousStatsHandler)))).
		Methods("GET", "POST")

	gMux.Handle("/api/user/test-prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserTestPreferenceHandler)))).
		Methods("GET", "POST", "DELETE")

	gMux.Handle("/api/k8sconfig", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler)))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/k8sconfig/contexts", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig)))).
		Methods("POST")
	gMux.Handle("/api/k8sconfig/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler))))

	gMux.Handle("/api/perf/load-test", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/load-test-smps", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestUsingSMPHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/load-test-prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestPrefencesHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler)))).
		Methods("GET")
	gMux.Handle("/api/perf/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler)))).
		Methods("GET")
	gMux.Handle("/api/mesh", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSMPServiceMeshes)))).
		Methods("GET")

	gMux.Handle("/api/smi/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSmiResultsHandler)))).
		Methods("GET")

	gMux.Handle("/api/mesh/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshAdapterConfigHandler))))
	gMux.Handle("/api/mesh/ops", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshOpsHandler)))).
		Methods("POST")
	gMux.Handle("/api/mesh/adapters", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.GetAllAdaptersHandler(w, req, provider)
	})))
	gMux.Handle("/api/mesh/adapter/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdapterPingHandler)))).
		Methods("GET")
	gMux.Handle("/api/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EventStreamHandler)))).
		Methods("GET")

	gMux.Handle("/api/grafana/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler)))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/grafana/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler)))).
		Methods("GET", "POST")
	gMux.Handle("/api/grafana/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler)))).
		Methods("GET")
	gMux.Handle("/api/grafana/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryRangeHandler)))).
		Methods("GET")
	gMux.Handle("/api/grafana/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler)))).
		Methods("GET")
	gMux.Handle("/api/grafana/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanGrafanaHandler))))

	gMux.Handle("/api/prometheus/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler)))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/prometheus/board_import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler)))).
		Methods("POST")
	gMux.Handle("/api/prometheus/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler)))).
		Methods("GET")
	gMux.Handle("/api/prometheus/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryRangeHandler)))).
		Methods("GET")
	gMux.Handle("/api/prometheus/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler)))).
		Methods("GET")
	gMux.Handle("/api/prometheus/static_board", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler)))).
		Methods("GET")
	gMux.Handle("/api/prometheus/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler)))).
		Methods("POST")
	gMux.Handle("/api/prometheus/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanPrometheusHandler))))

	gMux.Handle("/api/promGrafana/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanPromGrafanaHandler))))

	gMux.Handle("/api/experimental/patternfile/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PatternFileHandler)))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/experimental/patternfile/export/cytoscapejs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExportPatternFile)))).
		Methods("POST")
	gMux.Handle("/api/experimental/patternfile/import/cytoscapejs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ImportPatternFile)))).
		Methods("POST")
	gMux.Handle("/api/experimental/patternfile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PatternFileRequestHandler)))).
		Methods("POST", "GET")
	gMux.Handle("/api/experimental/patternfile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternHandler)))).
		Methods("GET")
	gMux.Handle("/api/experimental/patternfile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryPatternHandler)))).
		Methods("DELETE")
	gMux.HandleFunc("/api/experimental/oam/{type}", h.OAMRegisterHandler).Methods("GET", "POST")

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
	gMux.Handle("/api/user/performance/profiles/{id}/run", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler)))).
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

	gMux.PathPrefix("/api/system/graphql").Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GraphqlSystemHandler)))).Methods("GET", "POST")

	gMux.Handle("/logout", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LogoutHandler(w, req, provider)
	})))
	gMux.Handle("/login", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LoginHandler(w, req, provider, false)
	})))
	gMux.Handle("/api/token", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.TokenHandler(w, req, provider, false)
	})))
	gMux.Handle("/api/gettoken", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(
		func(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
			provider.ExtractToken(w, req)
		}))))

	// TODO: have to change this too
	gMux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo/meshery-logo.svg")
	}))

	// gMux.Handle("/", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	handlers.ServeUI(w, r, "", "../ui/out/")
	// }))))

	gMux.PathPrefix("/").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handlers.ServeUI(w, r, "", "../ui/out/")
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
