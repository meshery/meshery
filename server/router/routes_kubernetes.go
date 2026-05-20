package router

import (
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func registerKubernetesRoutes(gMux *mux.Router, h models.HandlerInterface) {
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
}
