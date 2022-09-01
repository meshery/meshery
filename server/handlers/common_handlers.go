// Package handlers : collection of handlers (aka "HTTP middleware")
package handlers

import (
	"net/http"
	"time"

	"github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/user/login UserAPI idGetUserLogin
// Handlers GET request for User login
//
// Redirects user for auth or issues session
// responses:
// 	200:

// LoginHandler redirects user for auth or issues session
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request, p models.Provider, fromMiddleWare bool) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	p.InitiateLogin(w, r, fromMiddleWare)
}

// swagger:route GET /api/user/logout UserAPI idGetUserLogout
// Handlers GET request for User logout
//
// Redirects user for auth or issues session
// responses:
// 	200:

// LogoutHandler destroys the session and redirects to home.
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request, p models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     h.config.ProviderCookieName,
		Value:    p.Name(),
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HttpOnly: true,
	})
	p.Logout(w, req)
}

// swagger:route GET /api/user/token UserAPI idGetTokenProvider
// Handle GET request for tokens
//
// Returns token from the actual provider in a file
// resposese:
// 	200:

// swagger:route POST /api/user/token UserAPI idPostTokenProvider
// Handle POST request for tokens
//
// Receives token from the actual provider
// resposese:
// 	200:

// TokenHandler Receives token from the actual provider
func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, p models.Provider, fromMiddleWare bool) {
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	p.TokenHandler(w, r, fromMiddleWare)
}
