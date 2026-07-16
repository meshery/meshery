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
	"strings"

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

// mesheryFileServingDir returns the canonical directory that ViewHandler and
// DownloadHandler are permitted to serve files from. Every producer of the
// "ViewLink"/"DownloadLink" event metadata that drives these endpoints writes
// under ~/.meshery/logs - registration logs (REGISTRY_LOG_FILE) and model
// generation logs. Confining reads to that directory is the trust boundary: it
// keeps ~/.meshery/config (the mesherydb.sql dump and provider auth tokens) and
// everything outside the Meshery home off-limits.
//
// NOTE: the server creates ~/.meshery/logs on startup (see cmd/main.go). Until
// it exists, resolveFileWithinDir fails closed and every request is rejected
// with 403 - correct for security, but worth knowing when debugging a fresh
// install that has not yet produced any logs.
func mesheryFileServingDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".meshery", "logs"), nil
}

// resolveFileWithinDir resolves requested against baseDir and reports whether
// the result is a real file confined to baseDir, returning the safe absolute
// path to open. It is the single validation boundary for the file view/download
// handlers and defends against:
//   - "../" traversal and absolute paths outside the base
//   - symlink traversal, e.g. a link inside the base pointing to /etc/shadow
//   - sibling-prefix bypass such as "<base>-backup" (via the trailing separator)
//   - NUL-byte injection
//
// It fails closed: ok is false for every rejection - including a target that
// cannot be resolved (a missing file) - so callers never receive a path outside
// the permitted directory. reason is a short diagnostic for server-side logging;
// it is deliberately not surfaced to the client.
func resolveFileWithinDir(requested, baseDir string) (resolved string, ok bool, reason string) {
	if requested == "" {
		return "", false, "empty file path"
	}
	if strings.ContainsRune(requested, '\x00') {
		return "", false, "file path contains a NUL byte"
	}

	// Canonicalize the base so the containment comparison runs against a fully
	// symlink-resolved absolute path (macOS /var, /home symlinks, t.TempDir).
	base, err := filepath.EvalSymlinks(baseDir)
	if err != nil {
		return "", false, "serving directory is unavailable: " + err.Error()
	}

	target := requested
	if !filepath.IsAbs(target) {
		target = filepath.Join(base, target)
	}

	// Resolve symlinks on the target so links that escape the base are caught,
	// then require the resolved path to sit strictly inside base. The trailing
	// separator stops a sibling like "<base>-backup" from satisfying the prefix.
	target, err = filepath.EvalSymlinks(filepath.Clean(target))
	if err != nil {
		return "", false, "file path could not be resolved: " + err.Error()
	}
	if !strings.HasPrefix(target, base+string(os.PathSeparator)) {
		return "", false, "resolved path escapes permitted directory " + base
	}

	// Only serve regular files. This rejects directories (which would fail with
	// EISDIR on read) and, more importantly, refuses to open a FIFO, socket, or
	// device node - opening a named pipe would block the handler goroutine
	// indefinitely waiting for a writer.
	info, err := os.Stat(target)
	if err != nil {
		return "", false, "resolved path could not be inspected: " + err.Error()
	}
	if !info.Mode().IsRegular() {
		return "", false, "resolved path is not a regular file"
	}
	return target, true, ""
}

// validateRequestedFile resolves the user-supplied file path against the
// permitted serving directory, returning the safe absolute path to open.
func (h *Handler) validateRequestedFile(filePath string) (resolved string, ok bool, reason string) {
	baseDir, err := mesheryFileServingDir()
	if err != nil {
		return "", false, "could not determine serving directory: " + err.Error()
	}
	return resolveFileWithinDir(filePath, baseDir)
}

// ViewHandler handles viewing the file content.
func (h *Handler) ViewHandler(responseWriter http.ResponseWriter, request *http.Request) {
	filePath, err := url.QueryUnescape(request.URL.Query().Get("file"))

	if err != nil {
		writeMeshkitError(responseWriter, ErrInvalidFileRequest(err), http.StatusBadRequest)
		return
	}

	resolvedPath, ok, reason := h.validateRequestedFile(filePath)
	if !ok {
		// Log the offending path server-side only; it is deliberately kept out
		// of the response so probing input is not reflected back to the caller.
		h.log.Warnf("rejected fileView request for %q: %s", filePath, reason)
		writeMeshkitError(responseWriter, ErrInvalidFilePath(), http.StatusForbidden)
		return
	}

	file, err := os.Open(resolvedPath)
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

	resolvedPath, ok, reason := h.validateRequestedFile(filePath)
	if !ok {
		// Log the offending path server-side only; it is deliberately kept out
		// of the response so probing input is not reflected back to the caller.
		h.log.Warnf("rejected fileDownload request for %q: %s", filePath, reason)
		writeMeshkitError(responseWriter, ErrInvalidFilePath(), http.StatusForbidden)
		return
	}

	file, err := os.Open(resolvedPath)
	if err != nil {
		writeMeshkitError(responseWriter, ErrReadFileContent(err, filePath), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Error(err)
		}
	}()

	fileName := filepath.Base(resolvedPath)
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
