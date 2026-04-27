package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJSONError writes a JSON-encoded {"error": message} body with the given
// HTTP status. Using JSON (instead of http.Error's plain text) keeps client
// response parsers — notably RTK Query's default baseQuery, which parses by
// Content-Type and treats application/json as JSON — from choking on error
// bodies that happen to start with a letter (e.g. "WorkspaceID or OrgID ...").
func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
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
