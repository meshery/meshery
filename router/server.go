package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/models"
)

type Router struct {
	s    *http.ServeMux
	port int
}

// New returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	mux := http.NewServeMux()

	mux.Handle("/api/user", h.AuthMiddleware(http.HandlerFunc(h.UserHandler)))
	mux.Handle("/api/config/sync", h.AuthMiddleware(http.HandlerFunc(h.SessionSyncHandler)))

	mux.Handle("/api/k8sconfig", h.AuthMiddleware(http.HandlerFunc(h.K8SConfigHandler)))
	mux.Handle("/api/k8sconfig/contexts", h.AuthMiddleware(http.HandlerFunc(h.GetContextsFromK8SConfig)))
	mux.Handle("/api/load-test", h.AuthMiddleware(http.HandlerFunc(h.LoadTestHandler)))
	mux.Handle("/api/results", h.AuthMiddleware(http.HandlerFunc(h.FetchResultsHandler)))

	mux.Handle("/api/mesh/manage", h.AuthMiddleware(http.HandlerFunc(h.MeshAdapterConfigHandler)))
	mux.Handle("/api/mesh/ops", h.AuthMiddleware(http.HandlerFunc(h.MeshOpsHandler)))
	mux.Handle("/api/mesh/adapters", h.AuthMiddleware(http.HandlerFunc(h.GetAllAdaptersHandler)))
	mux.Handle("/api/events", h.AuthMiddleware(http.HandlerFunc(h.EventStreamHandler)))

	mux.Handle("/api/grafana/config", h.AuthMiddleware(http.HandlerFunc(h.GrafanaConfigHandler)))
	mux.Handle("/api/grafana/boards", h.AuthMiddleware(http.HandlerFunc(h.GrafanaBoardsHandler)))
	mux.Handle("/api/grafana/query", h.AuthMiddleware(http.HandlerFunc(h.GrafanaQueryHandler)))
	mux.Handle("/api/grafana/query_range", h.AuthMiddleware(http.HandlerFunc(h.GrafanaQueryRangeHandler)))
	// mux.Handle("/api/grafana/boards", h.AuthMiddleware(http.HandlerFunc(h.SaveSelectedGrafanaBoardsHandler)))

	mux.Handle("/api/prometheus/config", h.AuthMiddleware(http.HandlerFunc(h.PrometheusConfigHandler)))
	mux.Handle("/api/prometheus/board_import", h.AuthMiddleware(http.HandlerFunc(h.GrafanaBoardImportForPrometheusHandler)))
	mux.Handle("/api/prometheus/query", h.AuthMiddleware(http.HandlerFunc(h.PrometheusQueryHandler)))
	mux.Handle("/api/prometheus/query_range", h.AuthMiddleware(http.HandlerFunc(h.PrometheusQueryRangeHandler)))
	mux.Handle("/api/prometheus/static_board", h.AuthMiddleware(http.HandlerFunc(h.PrometheusStaticBoardHandler)))
	mux.Handle("/api/prometheus/boards", h.AuthMiddleware(http.HandlerFunc(h.SaveSelectedPrometheusBoardsHandler)))

	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/login", h.LoginHandler)

	// TODO: have to change this too
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo.png")
	}))
	mux.Handle("/", h.AuthMiddleware(http.FileServer(http.Dir("../ui/out/"))))

	return &Router{
		s:    mux,
		port: port,
	}
}

func (r *Router) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.s)
}
