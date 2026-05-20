package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func registerMeshRoutes(gMux *mux.Router, h models.HandlerInterface) {
	gMux.Handle("/api/system/adapter/manage", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshAdapterConfigHandler)), models.ProviderAuth)))
	gMux.Handle("/api/system/adapter/operation", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.MeshOpsHandler)), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/system/adapters", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AdaptersHandler), models.ProviderAuth)))
	gMux.Handle("/api/system/availableAdapters", http.HandlerFunc(h.AvailableAdaptersHandler)).
		Methods("GET")

	gMux.PathPrefix("/api/extensions").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ExtensionsHandler), models.ProviderAuth))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")
}
