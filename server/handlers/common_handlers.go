// Package handlers : collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/meshery/meshery/server/core"
	"github.com/meshery/meshery/server/models"
)

// swagger:route GET /user/login UserAPI idGetUserLogin
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
func (h *Handler) LogoutHandler(w http.ResponseWriter, req *http.Request, user *models.User, p models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Clear all Meshery cookies to ensure complete logout
	for _, cookieName := range []string{
		h.config.ProviderCookieName,
		models.TokenCookieName,
		models.ProviderSessionCookieName,
	} {
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		})
	}
	_ = p.DeleteCapabilitiesForUser(user.ID.String())
	err := p.Logout(w, req)
	if err != nil {
		h.log.Error(models.ErrLogout(err))
		p.HandleUnAuthenticated(w, req)
		return
	}
	h.log.Info(fmt.Sprintf("logged out from %v provider", p.Name()))
	http.Redirect(w, req, "/provider", http.StatusFound)
}

// swagger:route GET /api/user/token UserAPI idGetTokenProvider
// Handle GET request for tokens
//
// Returns token from the actual provider in a file
// response:
// 	200:

// swagger:route POST /api/user/token UserAPI idPostTokenProvider
// Handle POST request for tokens
//
// Receives token from the actual provider
// response:
// 	200:

// TokenHandler Receives token from the actual provider
func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, p models.Provider, fromMiddleWare bool) {
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	p.TokenHandler(w, r, fromMiddleWare)
}

// swagger:route GET /api/system/viewFile system viewFile idFileViewer
// Handles GET request to view a file.
//
// Retrieves and displays the content of the specified file as plain text.
//
// responses:
//   200:
//   500:

// ViewHandler handles viewing the file content.
func (h *Handler) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))

	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the content type to plain text
	responseWriter.Header().Set("Content-Type", "text/plain")

	// Copy the file content to the response writer
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/downloadFile system downloadFile idDownloadFile
// Handles GET request to download a file.
//
// Retrieves and initiates a download for the specified file.
//
// responses:
//   200:
//   500:

// DownloadHandler handles downloading the file.
func (h *Handler) DownloadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	responseWriter.Header().Set("Content-Type", "text/plain")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	_, err = io.Copy(responseWriter, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Deep-link and redirect support to land user on their originally requested page post authentication instead of dropping user on the root (home) page.
func GetRefURL(req *http.Request) string {
	return core.EncodeRefUrl(*req.URL)
}

func (h *Handler) HandleErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(http.StatusInternalServerError)

	// Define the error response structure
	type ErrorResponse struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	// Create an error response instance
	errorResponse := ErrorResponse{
		Status:  http.StatusInternalServerError,
		Message: "We encountered an error while processing your request. Please try again later.",
	}

	// Encode and send the error response as JSON
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		h.log.Error(models.ErrMarshal(err, "error response"))
	}
}
