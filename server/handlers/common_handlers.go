// Package handlers : collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/server/core"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/viper"
)

var (
	errMissingFilePath     = errors.New("missing file query parameter")
	errFileAccessDenied    = errors.New("requested file is not available")
	errUnsupportedFileType = errors.New("requested path must point to a file")
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
	h.serveAllowedFile(responseWriter, request, false)
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
	h.serveAllowedFile(responseWriter, request, true)
}

func (h *Handler) serveAllowedFile(responseWriter http.ResponseWriter, request *http.Request, download bool) {
	filePath, err := validateMesheryLogFilePath(request.URL.Query().Get("file"))
	if err != nil {
		switch {
		case errors.Is(err, errMissingFilePath), errors.Is(err, errUnsupportedFileType):
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
		case errors.Is(err, errFileAccessDenied):
			http.Error(responseWriter, err.Error(), http.StatusForbidden)
		case errors.Is(err, os.ErrNotExist):
			http.Error(responseWriter, "requested file was not found", http.StatusNotFound)
		default:
			h.log.Error(err)
			http.Error(responseWriter, "failed to process requested file", http.StatusInternalServerError)
		}
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(responseWriter, "requested file was not found", http.StatusNotFound)
			return
		}

		h.log.Error(err)
		http.Error(responseWriter, "failed to open requested file", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	responseWriter.Header().Set("Content-Type", "text/plain")
	if download {
		responseWriter.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(filePath)))
	}

	if _, err = io.Copy(responseWriter, file); err != nil {
		h.log.Error(err)
		http.Error(responseWriter, "failed to stream requested file", http.StatusInternalServerError)
	}
}

func validateMesheryLogFilePath(filePath string) (string, error) {
	if filePath == "" {
		return "", errMissingFilePath
	}

	if !filepath.IsAbs(filePath) {
		return "", errFileAccessDenied
	}

	allowedRoots, err := mesheryLogRoots()
	if err != nil {
		return "", err
	}

	cleanFilePath := filepath.Clean(filePath)
	isAllowed, err := pathWithinAllowedRoots(cleanFilePath, allowedRoots)
	if err != nil {
		return "", err
	}
	if !isAllowed {
		return "", errFileAccessDenied
	}

	resolvedFilePath, err := filepath.EvalSymlinks(cleanFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		return "", fmt.Errorf("resolve requested file path: %w", err)
	}

	isAllowed, err = pathWithinAllowedRoots(resolvedFilePath, allowedRoots)
	if err != nil {
		return "", err
	}
	if !isAllowed {
		return "", errFileAccessDenied
	}

	info, err := os.Stat(resolvedFilePath)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", errUnsupportedFileType
	}

	return resolvedFilePath, nil
}

func mesheryLogRoots() ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("determine home directory: %w", err)
	}

	rootCandidates := []string{
		filepath.Join(homeDir, ".meshery", "logs"),
	}

	if registryLogFile := viper.GetString("REGISTRY_LOG_FILE"); registryLogFile != "" {
		rootCandidates = append(rootCandidates, filepath.Dir(registryLogFile))
	}

	roots := make([]string, 0, len(rootCandidates))
	seen := make(map[string]struct{}, len(rootCandidates))
	for _, rootCandidate := range rootCandidates {
		resolvedRoot, err := normalizeAllowedRoot(rootCandidate)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[resolvedRoot]; ok {
			continue
		}
		seen[resolvedRoot] = struct{}{}
		roots = append(roots, resolvedRoot)
	}

	return roots, nil
}

func normalizeAllowedRoot(root string) (string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve allowed file root %q: %w", root, err)
	}

	resolvedRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return absRoot, nil
		}
		return "", fmt.Errorf("resolve allowed file root %q: %w", root, err)
	}

	return resolvedRoot, nil
}

func pathWithinAllowedRoots(path string, allowedRoots []string) (bool, error) {
	for _, allowedRoot := range allowedRoots {
		relPath, err := filepath.Rel(allowedRoot, path)
		if err != nil {
			return false, fmt.Errorf("resolve relative path from %q to %q: %w", allowedRoot, path, err)
		}
		if relPath == "." || (relPath != ".." && !strings.HasPrefix(relPath, ".."+string(os.PathSeparator))) {
			return true, nil
		}
	}

	return false, nil
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
