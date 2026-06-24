package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meshery/meshery/server/models/httputil"
)

// Response helpers
// ----------------
//
// These four helpers are the canonical way to write an HTTP response from
// server/handlers. Never use http.Error — it emits Content-Type: text/plain
// which crashes RTK Query's default baseQuery on the UI (see
// docs/content/en/project/contributing/error-contract.md).
//
// The real implementations live in server/models/httputil so both
// server/handlers and server/models (and any future sibling) can call them
// without an import cycle. These wrappers preserve the original unexported
// identifiers so none of the ~150 existing call sites in this package had to
// change during the migration.
//
// Reach for:
//   - writeMeshkitError     — ANY error path. If err wraps a *meshkiterrors.Error
//                             or *ErrorV2, the code/severity/cause/remediation
//                             survive onto the wire. If it doesn't, the .Error()
//                             string is still emitted as JSON.
//   - writeJSONError        — error paths where the message is a bare string with
//                             no MeshKit wrapper. Prefer promoting the string to
//                             a MeshKit error and using writeMeshkitError instead.
//   - writeJSONMessage      — success paths that return a small status or result
//                             payload (e.g. {"message": "deleted"}).
//   - writeJSONEmptyObject  — success paths that need to return an empty JSON
//                             object ({}) with the Content-Type header set.

func writeJSONError(w http.ResponseWriter, message string, status int) {
	httputil.WriteJSONError(w, message, status)
}

func writeMeshkitError(w http.ResponseWriter, err error, status int) {
	httputil.WriteMeshkitError(w, err, status)
}

func writeJSONMessage(w http.ResponseWriter, payload any, status int) {
	httputil.WriteJSONMessage(w, payload, status)
}

func writeJSONEmptyObject(w http.ResponseWriter, status int) {
	httputil.WriteJSONEmptyObject(w, status)
}

const (
	defaultPageSize = 25
	queryParamTrue  = "true"
)

func getPaginationParams(req *http.Request) (page, offset, limit int, search, order, sortOnCol, status string) {

	urlValues := req.URL.Query()
	page, _ = strconv.Atoi(urlValues.Get("page"))
	limitstr := urlValues.Get("pagesize")
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = defaultPageSize
		}
	}

	search = urlValues.Get("search")
	order = urlValues.Get("order")
	sortOnCol = urlValues.Get("sort")
	status = urlValues.Get("status")

	if page < 0 {
		page = 0
	}
	offset = page * limit

	if sortOnCol == "" {
		sortOnCol = "updated_at"
	}
	return
}

// Extracts specified boolean query parameters from the request and returns a map of params and their value.
func extractBoolQueryParams(r *http.Request, params ...string) (map[string]bool, error) {
	result := make(map[string]bool)
	for _, param := range params {
		val, err := strconv.ParseBool(r.URL.Query().Get(param))
		if err != nil {
			val = false
		}
		result[param] = val
	}
	return result, nil
}

// SafeFilePath validates a file path to prevent path traversal attacks.
// It ensures the path is within allowed directories (e.g., ~/.meshery or os.TempDir()).
func SafeFilePath(providedPath string) (string, error) {
	if providedPath == "" {
		return "", fmt.Errorf("file path cannot be empty")
	}

	if strings.Contains(providedPath, "..") {
		return "", fmt.Errorf("path traversal attempts are not allowed")
	}

	absPath, err := filepath.Abs(providedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %v", err)
	}

	var resolvedPath string
	if _, err := os.Stat(absPath); err == nil {
		resolvedPath, err = filepath.EvalSymlinks(absPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve symlinks: %v", err)
		}
	} else {
		dir := filepath.Dir(absPath)
		resolvedDir, err := filepath.EvalSymlinks(dir)
		if err != nil {
			return "", fmt.Errorf("failed to resolve symlinks for directory: %v", err)
		}
		resolvedPath = filepath.Join(resolvedDir, filepath.Base(absPath))
	}

	// Allowed base directories
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine user home directory")
	}

	allowedBaseDirs := []string{
		filepath.Join(home, ".meshery"),
		os.TempDir(),
	}

	isAllowed := false
	for _, baseDir := range allowedBaseDirs {
		resolvedBase, err := filepath.EvalSymlinks(baseDir)
		if err != nil {
			resolvedBase, _ = filepath.Abs(baseDir)
		}

		rel, err := filepath.Rel(resolvedBase, resolvedPath)
		if err == nil && !strings.HasPrefix(rel, "..") {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return "", fmt.Errorf("file path is outside allowed directories")
	}

	return resolvedPath, nil
}
