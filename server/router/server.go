package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// Router represents Meshery router
type Router struct {
	S    *mux.Router
	port int
}

// NewRouter returns a new ServeMux with app routes.
func NewRouter(_ context.Context, h models.HandlerInterface, port int, g http.Handler, gp http.Handler) *Router {
	gMux := mux.NewRouter()
	gMux.Use(otelmux.Middleware("meshery-server"))

	// Domain-route registration (see per-file register*Routes functions)
	registerSystemRoutes(gMux, h, g, gp)
	registerAuthRoutes(gMux, h)
	registerPatternRoutes(gMux, h)
	registerPerformanceRoutes(gMux, h)
	registerKubernetesRoutes(gMux, h)
	registerTelemetryRoutes(gMux, h)
	registerMeshRoutes(gMux, h)
	registerRemainingRoutes(gMux, h)

	return &Router{
		S:    gMux,
		port: port,
	}
}

// registerRemainingRoutes covers routes that don't warrant a dedicated file:
// credentials, legacy environments/workspaces/connections, identity/orgs,
// static file serving, and the SPA catch-all.
func registerRemainingRoutes(gMux *mux.Router, h models.HandlerInterface) {
	// User Credentials
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserCredentials), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/credentials/{credentialID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserCredentialByID), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateUserCredential), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteUserCredential), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/integrations/credentials", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveUserCredential), models.ProviderAuth))).
		Methods("POST")

	// Legacy environments (under /api/environments/*, distinct from /api/integrations/environments/*)
	gMux.Handle("/api/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironments), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironmentByIDHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteEnvironmentHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveEnvironment), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/environments/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateEnvironmentHandler), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/environments/{environmentID}/connections/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddConnectionToEnvironmentHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/environments/{environmentID}/connections/{connectionID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveConnectionFromEnvironmentHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/environments/{environmentID}/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnectionsOfEnvironmentHandler), models.ProviderAuth))).
		Methods("GET")

	// Legacy workspaces (under /api/workspaces/*)
	gMux.Handle("/api/workspaces", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetWorkspacesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateWorkspaceHandler), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetWorkspaceByIdHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/environments", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetEnvironmentsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/environments/{environmentID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddEnvironmentToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/environments/{environmentID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveEnvironmentFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}/designs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetDesignsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/designs/{designID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddDesignToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/designs/{designID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveDesignFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}/views", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetViewsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/views/{viewID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddViewToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/views/{viewID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveViewFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/workspaces/{id}/teams", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetTeamsOfWorkspaceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/workspaces/{id}/teams/{teamID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.AddTeamToWorkspaceHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/workspaces/{id}/teams/{teamID}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RemoveTeamFromWorkspaceHandler), models.ProviderAuth))).
		Methods("DELETE")

	// Identity / Organizations
	gMux.Handle("/api/identity/orgs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetOrganizations), models.ProviderAuth))).
		Methods("GET")

	// Legacy connections (under /api/integrations/connections/{connectionId} with old handlers)
	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnectionByID), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateConnectionById), models.ProviderAuth))).
		Methods("PUT")
	gMux.Handle("/api/integrations/connections/register", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProcessConnectionRegistration), models.ProviderAuth))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/integrations/connections/{connectionId}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteConnection), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/integrations/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveConnection), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/integrations/connections", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetConnections), models.ProviderAuth))).
		Methods("GET")

	// Favicon
	gMux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		http.ServeFile(w, r, "../ui/out/static/img/meshery-logo/meshery-logo.svg")
	}))

	// Static files (mesh model images)
	fs := http.FileServer(http.Dir("../../ui"))
	gMux.PathPrefix("/ui/public/static/img/meshmodels").Handler(http.StripPrefix("/ui/", fs)).Methods("GET", "HEAD")

	// Next.js build artifacts (no auth)
	gMux.PathPrefix("/_next").
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "", "../../ui/out/")
		})).
		Methods("GET", "HEAD")

	// SPA catch-all (auth required)
	gMux.PathPrefix("/").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "", "../../ui/out/")
		}), models.ProviderAuth))).
		Methods("GET", "HEAD")
}

// Run starts the http server
func (r *Router) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.S)
}
