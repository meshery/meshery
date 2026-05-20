package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// registerSystemRoutes registers routes for system-level operations: GraphQL,
// version info, database management, session sync, events, and health probes.
// g and gp are the GraphQL and GraphQL playground handlers from NewRouter.
func registerSystemRoutes(gMux *mux.Router, h models.HandlerInterface, g, gp http.Handler) {
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

	gMux.Handle("/api/system/sync", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.K8sFSMMiddleware(h.SessionSyncHandler))), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/system/fileDownload", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.DownloadHandler), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/system/fileView", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ViewHandler), models.NoAuth))).Methods("GET")

	// Events
	gMux.Handle("/api/system/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ClientEventHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/events", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetAllEvents), models.ProviderAuth))).
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

	// Health probes
	gMux.HandleFunc("/healthz/live", h.K8sHealthzHandler).Methods("GET")
	gMux.HandleFunc("/healthz/ready", h.K8sHealthzHandler).Methods("GET")
}
