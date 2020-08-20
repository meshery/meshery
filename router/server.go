package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

// Router represents Meshery router
type Router struct {
	s    *mux.Router
	port int
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	gMux := mux.NewRouter()

	gMux.Handle("/api/server/version", http.HandlerFunc(h.ServerVersionHandler))

	gMux.Handle("/api/provider", http.HandlerFunc(h.ProviderHandler))
	gMux.Handle("/api/providers", http.HandlerFunc(h.ProvidersHandler))
	gMux.Handle("/provider/", http.HandlerFunc(h.ProviderUIHandler))

	gMux.Handle("/api/config/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SessionSyncHandler))))

	gMux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler))))
	gMux.Handle("/api/user/stats", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AnonymousStatsHandler))))

	gMux.Handle("/api/user/test-prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserTestPreferenceHandler))))

	gMux.Handle("/api/k8sconfig", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler))))
	gMux.Handle("/api/k8sconfig/contexts", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig))))
	gMux.Handle("/api/k8sconfig/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler))))
	gMux.Handle("/api/mesh/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.InstalledMeshesHandler))))

	gMux.Handle("/api/perf/load-test", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler))))
	gMux.Handle("/api/perf/load-test-smps", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestUsingSMPSHandler))))
	gMux.Handle("/api/perf/load-test-prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestPrefencesHandler))))
	gMux.Handle("/api/perf/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler))))
	gMux.Handle("/api/perf/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler))))

	gMux.Handle("/api/mesh/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshAdapterConfigHandler))))
	gMux.Handle("/api/mesh/ops", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshOpsHandler))))
	gMux.Handle("/api/mesh/adapters", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.GetAllAdaptersHandler(w, req, provider)
	})))
	gMux.Handle("/api/mesh/adapter/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdapterPingHandler))))
	gMux.Handle("/api/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EventStreamHandler))))

	gMux.Handle("/api/grafana/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler))))
	gMux.Handle("/api/grafana/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler))))
	gMux.Handle("/api/grafana/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler))))
	gMux.Handle("/api/grafana/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryRangeHandler))))
	gMux.Handle("/api/grafana/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler))))
	gMux.Handle("/api/grafana/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanGrafanaHandler))))

	gMux.Handle("/api/prometheus/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler))))
	gMux.Handle("/api/prometheus/board_import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler))))
	gMux.Handle("/api/prometheus/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler))))
	gMux.Handle("/api/prometheus/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryRangeHandler))))
	gMux.Handle("/api/prometheus/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler))))
	gMux.Handle("/api/prometheus/static_board", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler))))
	gMux.Handle("/api/prometheus/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler))))
	gMux.Handle("/api/prometheus/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanPrometheusHandler))))

	gMux.Handle("/api/promGrafana/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ScanPromGrafanaHandler))))

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

	gMux.Methods("GET").PathPrefix("/").Handler(h.ProviderMiddleware(h.AuthMiddleware(http.StripPrefix("/", http.FileServer(http.Dir("../ui/out/"))))))

	return &Router{
		s:    gMux,
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
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.s)
}
