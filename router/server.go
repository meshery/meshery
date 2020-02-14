package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/models"
)

// Router represents Meshery router
type Router struct {
	s    *http.ServeMux
	port int
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/provider", h.ProviderHandler)
	mux.HandleFunc("/api/providers", h.ProvidersHandler)
	mux.HandleFunc("/provider/", h.ProviderUIHandler)

	mux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler))))
	mux.Handle("/api/user/stats", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AnonymousStatsHandler))))
	mux.Handle("/api/config/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SessionSyncHandler))))

	mux.Handle("/api/k8sconfig", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler))))
	mux.Handle("/api/k8sconfig/contexts", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig))))
	mux.Handle("/api/k8sconfig/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler))))
	mux.Handle("/api/mesh/scan", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.InstalledMeshesHandler))))

	mux.Handle("/api/load-test", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler))))
	mux.Handle("/api/load-test-smps", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestUsingSMPSHandler))))
	mux.Handle("/api/load-test-prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestPrefencesHandler))))
	mux.Handle("/api/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler))))
	mux.Handle("/api/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler))))

	mux.Handle("/api/mesh/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshAdapterConfigHandler))))
	mux.Handle("/api/mesh/ops", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshOpsHandler))))
	mux.Handle("/api/mesh/adapters", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.GetAllAdaptersHandler(w, req, provider)
	})))
	mux.Handle("/api/mesh/adapter/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdapterPingHandler))))
	mux.Handle("/api/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EventStreamHandler))))

	mux.Handle("/api/grafana/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler))))
	mux.Handle("/api/grafana/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler))))
	mux.Handle("/api/grafana/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler))))
	mux.Handle("/api/grafana/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryRangeHandler))))
	mux.Handle("/api/grafana/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler))))

	mux.Handle("/api/prometheus/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler))))
	mux.Handle("/api/prometheus/board_import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler))))
	mux.Handle("/api/prometheus/query", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler))))
	mux.Handle("/api/prometheus/query_range", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryRangeHandler))))
	mux.Handle("/api/prometheus/ping", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler))))
	mux.Handle("/api/prometheus/static_board", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler))))
	mux.Handle("/api/prometheus/boards", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler))))

	mux.Handle("/logout", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LogoutHandler(w, req, provider)
	})))
	mux.Handle("/login", h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LoginHandler(w, req, provider, false)
	})))

	// TODO: have to change this too
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo.png")
	}))

	mux.Handle("/", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeUI(w, r, "", "../ui/out/")
	}))))

	return &Router{
		s:    mux,
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
