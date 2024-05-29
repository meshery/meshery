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
	"time"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
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
	err := p.Logout(w, req)
	if err != nil {
		logrus.Errorf("Error performing logout: %v", err.Error())
		p.HandleUnAuthenticated(w, req)
		return
	}
	logrus.Infof("successfully logged out from %v provider", p.Name())
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

// swagger:route GET /api/system/fileView
// Handlers GET request for view file
//
// Redirects to a url to view file
// responses:
// 	200:

// ViewHandler redirects to view the file
func (h *Handler) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	fmt.Println(filePath)

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

	var jsonData interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonData); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(responseWriter)
	if err := encoder.Encode(jsonData); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/fileDownload
// Handlers GET request for download file
//
// Redirects to a url to download file
// responses:
// 	200:

// DownloadHandler redirects to download the file
func (h *Handler) DownloadHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileName := filepath.Base(filePath)
	responseWriter.Header().Set("Content-Type", "application/octet-stream")
	responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	_, err = io.Copy(responseWriter, file)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}
