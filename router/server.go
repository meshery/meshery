package router

import (
	"context"
	"fmt"
	"net/http"

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

	mux.Handle("/api/user", h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler)))
	mux.Handle("/api/user/stats", h.AuthMiddleware(h.SessionInjectorMiddleware(h.AnonymousStatsHandler)))
	mux.Handle("/api/config/sync", h.AuthMiddleware(h.SessionInjectorMiddleware(h.SessionSyncHandler)))

	mux.Handle("/api/k8sconfig", h.AuthMiddleware(h.SessionInjectorMiddleware(h.K8SConfigHandler)))
	mux.Handle("/api/k8sconfig/contexts", h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig)))
	mux.Handle("/api/k8sconfig/ping", h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesPingHandler)))
	mux.Handle("/api/mesh/scan", h.AuthMiddleware(h.SessionInjectorMiddleware(h.InstalledMeshesHandler)))

	mux.Handle("/api/load-test", h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler)))
	mux.Handle("/api/load-test-prefs", h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestPrefencesHandler)))
	mux.Handle("/api/results", h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler)))

	mux.Handle("/api/mesh/manage", h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshAdapterConfigHandler)))
	mux.Handle("/api/mesh/ops", h.AuthMiddleware(h.SessionInjectorMiddleware(h.MeshOpsHandler)))
	mux.Handle("/api/mesh/adapters", h.AuthMiddleware(http.HandlerFunc(h.GetAllAdaptersHandler)))
	mux.Handle("/api/mesh/adapter/ping", h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdapterPingHandler)))
	mux.Handle("/api/events", h.AuthMiddleware(h.SessionInjectorMiddleware(h.EventStreamHandler)))

	mux.Handle("/api/grafana/config", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler)))
	mux.Handle("/api/grafana/boards", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler)))
	mux.Handle("/api/grafana/query", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler)))
	mux.Handle("/api/grafana/query_range", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryRangeHandler)))
	mux.Handle("/api/grafana/ping", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler)))

	mux.Handle("/api/prometheus/config", h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler)))
	mux.Handle("/api/prometheus/board_import", h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler)))
	mux.Handle("/api/prometheus/query", h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler)))
	mux.Handle("/api/prometheus/query_range", h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryRangeHandler)))
	mux.Handle("/api/prometheus/ping", h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler)))
	mux.Handle("/api/prometheus/static_board", h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler)))
	mux.Handle("/api/prometheus/boards", h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler)))

	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		h.LoginHandler(w, r, false)
	})

	// TODO: have to change this too
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo.png")
	}))
	mux.Handle("/", h.AuthMiddleware(http.FileServer(http.Dir("../ui/out/"))))

	mux.Handle("/settings", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/settings.html")
	}))

	mux.Handle("/configure", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/configure.html")
	}))

	mux.Handle("/management", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/management.html")
	}))

	mux.Handle("/performance", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/performance.html")
	}))

	mux.Handle("/results", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/results.html")
	}))

	mux.Handle("/404", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/404.html")
	}))

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
