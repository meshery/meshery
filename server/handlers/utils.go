package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meshery/meshery/server/models/httputil"
	meshkitutils "github.com/meshery/meshkit/utils"
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

// allowedDirError signals that a requested path resolved outside every directory
// the file endpoints may serve. A distinct sentinel type lets callers map it to
// HTTP 400 while treating any other failure (a missing or unreadable file) as
// HTTP 500, so a legitimate not-found is never reported as an attack. It is an
// internal control-flow signal; the handler translates it into a structured
// MeshKit error at the response boundary.
type allowedDirError string

func (e allowedDirError) Error() string { return string(e) }

// errFileOutsideAllowedDirs is the sentinel returned by SafeOpenFile for a path
// outside the allowed directories.
const errFileOutsideAllowedDirs = allowedDirError("file path is outside allowed directories")

// allowedFileDirs returns the directories the fileView/fileDownload endpoints
// are permitted to serve from. Only the Meshery registry log directory is
// exposed: these endpoints exist solely to surface those logs, and broadening
// the set to the whole home directory or os.TempDir() would let the
// unauthenticated endpoints read the Meshery database, provider credentials, or
// other processes' temporary files.
//
// Both home-directory derivations are included because the server writes these
// logs through two different helpers whose results can diverge: registry-logs.log
// via os.UserHomeDir() ($HOME) in cmd/main.go, and model-generation.log via
// meshkit's GetHome() (the passwd entry). $HOME and the passwd home differ under
// systemd units, sudo, or containers started with an explicit HOME.
func allowedFileDirs() []string {
	dirs := make([]string, 0, 2)
	seen := make(map[string]bool, 2)
	add := func(home string) {
		if home == "" {
			return
		}
		dir := filepath.Join(home, ".meshery", "logs")
		if !seen[dir] {
			seen[dir] = true
			dirs = append(dirs, dir)
		}
	}
	// A failed home lookup only drops that one entry; it must not disable the
	// remaining allowed directories.
	if home, err := os.UserHomeDir(); err == nil {
		add(home)
	}
	add(meshkitutils.GetHome())
	return dirs
}

// SafeOpenFile opens providedPath for reading, but only when it resolves to a
// file inside one of the allowedFileDirs. The file is opened through os.Root so
// that ".." components and symlinks cannot escape the allowed directory, and so
// that validation and open are a single operation with no time-of-check/
// time-of-use gap (validating a path string that the caller later re-opens lets
// a component be swapped for a symlink in between). The caller owns the returned
// *os.File and must close it.
func SafeOpenFile(providedPath string) (*os.File, error) {
	if providedPath == "" {
		return nil, errFileOutsideAllowedDirs
	}

	// filepath.Abs applies filepath.Clean, lexically resolving any ".." elements
	// up front, so no separate substring check is needed.
	absPath, err := filepath.Abs(providedPath)
	if err != nil {
		return nil, err
	}

	for _, dir := range allowedFileDirs() {
		rel, err := filepath.Rel(dir, absPath)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			continue // absPath is not inside this directory
		}

		root, err := os.OpenRoot(dir)
		if err != nil {
			continue // directory is missing or not a directory
		}
		// os.Root confines the open to dir: a rel that escapes via ".." or a
		// symlink pointing outside dir fails here instead of being followed. The
		// returned file remains valid after the root is closed.
		file, err := root.Open(rel)
		_ = root.Close()
		if err != nil {
			// Inside an allowed dir but not openable (missing, unreadable, or a
			// symlink escaping the root). Never serve it; surface a read error
			// rather than an out-of-bounds rejection.
			return nil, err
		}
		return file, nil
	}

	return nil, errFileOutsideAllowedDirs
}
