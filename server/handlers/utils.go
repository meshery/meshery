package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/logger"
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
	maxPageSize     = 100
	queryParamTrue  = "true"
)

func getPaginationParams(req *http.Request) (page, offset, limit int, search, order, sortOnCol, status string, err error) {
	urlValues := req.URL.Query()
	page, err = parsePaginationQueryValue(urlValues.Get("page"), "page")
	if err != nil {
		return
	}

	limitstr := urlValues.Get("pagesize")
	switch limitstr {
	case "":
		limit = defaultPageSize
	case "all":
		limit = maxPageSize
	default:
		limit, err = parsePageSize(limitstr)
		if err != nil {
			return
		}
	}

	search = urlValues.Get("search")
	order = urlValues.Get("order")
	sortOnCol = urlValues.Get("sort")
	status = urlValues.Get("status")

	offset = page * limit

	if sortOnCol == "" {
		sortOnCol = "updated_at"
	}
	return
}

func parsePaginationQueryValue(value, queryParam string) (int, error) {
	if value == "" {
		return 0, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, ErrPaginationQuery(queryParam, value, err.Error())
	}

	if parsedValue < 0 {
		return 0, ErrPaginationQuery(queryParam, value, "must be greater than or equal to 0")
	}

	return parsedValue, nil
}

func parsePageSize(value string) (int, error) {
	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, ErrPaginationQuery("pagesize", value, err.Error())
	}

	if parsedValue <= 0 || parsedValue > maxPageSize {
		return 0, ErrPaginationQuery("pagesize", value, fmt.Sprintf("must be between 1 and %d or 'all'", maxPageSize))
	}

	return parsedValue, nil
}

func writePaginationError(log logger.Handler, w http.ResponseWriter, err error) {
	if log != nil {
		log.Error(err)
	}
	writeMeshkitError(w, err, http.StatusBadRequest)
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
