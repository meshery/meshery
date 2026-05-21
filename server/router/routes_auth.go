package router

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// registerAuthRoutes registers routes for authentication, provider, user
// identity, and preferences.
func registerAuthRoutes(gMux *mux.Router, h models.HandlerInterface) {
	gMux.HandleFunc("/api/provider", h.ProviderHandler)
	gMux.HandleFunc("/error", h.HandleErrorHandler)
	gMux.HandleFunc("/api/providers", h.ProvidersHandler).
		Methods("GET")
	gMux.PathPrefix("/api/provider/extension").
		Handler(h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.ProviderComponentsHandler)), models.ProviderAuth))).
		Methods("GET", "POST", "OPTIONS", "PUT", "DELETE")
	gMux.Handle("/api/provider/capabilities", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.ProviderCapabilityHandler), models.ProviderAuth))).
		Methods("GET")

	// Provider UI
	gMux.HandleFunc("/provider/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
	})
	gMux.PathPrefix("/provider/_next").
		Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
		}))
	gMux.PathPrefix("/provider").
		Handler(http.HandlerFunc(h.ProviderUIHandler)).
		Methods("GET")
	gMux.HandleFunc("/auth/login", h.ProviderUIHandler).
		Methods("GET")

	// User / Identity
	gMux.Handle("/api/user", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/users/profile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/profile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserByIDHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/users/profile/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUserByIDHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/users", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUsers), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/identity/orgs/{orgID}/users/keys", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetUsersKeys), models.ProviderAuth))).
		Methods("GET")

	// Preferences
	gMux.Handle("/api/user/prefs", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserPrefsHandler), models.ProviderAuth))).
		Methods("GET", "POST")
	gMux.Handle("/api/user/prefs/perf", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UserTestPreferenceHandler), models.ProviderAuth))).
		Methods("GET", "POST", "DELETE")

	// Login / Logout / Token
	gMux.Handle("/user/logout", h.ProviderMiddleware(h.SessionInjectorMiddleware(func(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
		h.LogoutHandler(w, req, user, provider)
	})))
	gMux.Handle("/user/login", h.NoCacheMiddleware(h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.LoginHandler(w, req, provider, false)
	}))))
	gMux.Handle("/api/user/token", h.NoCacheMiddleware(h.ProviderMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		h.TokenHandler(w, req, provider, false)
	})))).Methods("POST", "GET")
	gMux.Handle("/api/token", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(
		func(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
			provider.ExtractToken(w, req)
		}), models.ProviderAuth))).Methods("GET")

	// Auth redirect
	gMux.HandleFunc("/auth/redirect", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		http.SetCookie(w, &http.Cookie{
			Name:     models.TokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
			Expires:  time.Now().Add(24 * time.Hour),
		})
		h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
	}).
		Methods("GET")
}
