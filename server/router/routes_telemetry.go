package router

import (
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func registerTelemetryRoutes(gMux *mux.Router, h models.HandlerInterface) {
	gMux.Handle("/api/telemetry/metrics/grafana/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaConfigHandler), models.ProviderAuth))).
		Methods("GET", "POST", "DELETE")
	gMux.Handle("/api/telemetry/metrics/grafana/boards/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardsHandler), models.ProviderAuth))).
		Methods("GET", "POST")
	gMux.Handle("/api/telemetry/metrics/grafana/query/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/grafana/query_range/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaQueryRangeHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/grafana/ping/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaPingHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/telemetry/metrics/config", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusConfigHandler), models.ProviderAuth))).
		Methods("GET", "POST", "DELETE")
	gMux.Handle("/api/telemetry/metrics/board_import/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GrafanaBoardImportForPrometheusHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/telemetry/metrics/query/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/prometheus/query_range/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusQueryRangeHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/ping/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusPingHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/static-board/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PrometheusStaticBoardHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/telemetry/metrics/boards/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveSelectedPrometheusBoardsHandler), models.ProviderAuth))).
		Methods("POST")
}
