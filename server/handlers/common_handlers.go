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

// LoginHandler redirects user for auth or issues session
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request, p models.Provider, fromMiddleWare bool) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	p.InitiateLogin(w, r, fromMiddleWare)
}

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

// TokenHandler Receives token from the actual provider
func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request, p models.Provider, fromMiddleWare bool) {
	// if r.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	p.TokenHandler(w, r, fromMiddleWare)
}

// ViewHandler handles viewing the file content.
func (h *Handler) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))

	if err != nil {
		writeMeshkitError(responseWriter, ErrInvalidFileRequest(err), http.StatusBadRequest)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		writeMeshkitError(responseWriter, ErrReadFileContent(err, filePath), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	// Set the content type to plain text
	responseWriter.Header().Set("Content-Type", "text/plain")

	// Copy the file content to the response writer. If io.Copy fails mid-stream
	// the response status and headers are already committed, so we log the
	// error for diagnostics and return rather than attempting a second write.
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		h.log.Error(models.ErrCopy(err, filePath))
		return
	}
}

// DownloadHandler handles downloading the file.
func (h *Handler) DownloadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	if err != nil {
		writeMeshkitError(responseWriter, ErrInvalidFileRequest(err), http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		writeMeshkitError(responseWriter, ErrReadFileContent(err, filePath), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	fileName := filepath.Base(filePath)
	responseWriter.Header().Set("Content-Type", "text/plain")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// See ViewHandler: response has already started streaming, so a second
	// error write would corrupt the body. Log and return.
	_, err = io.Copy(responseWriter, file)
	if err != nil {
		h.log.Error(models.ErrCopy(err, filePath))
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
