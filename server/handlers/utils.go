package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	// DefaultMaxBodySize is the default maximum request body size (10 MB)
	DefaultMaxBodySize int64 = 10 * 1024 * 1024
)

// readBodyWithLimit reads the request body up to maxBytes.
// Returns the body bytes or an error if the body exceeds the limit.
func readBodyWithLimit(r *http.Request, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBodySize
	}

	limitedBody := http.MaxBytesReader(nil, r.Body, maxBytes)
	defer func() {
		_ = limitedBody.Close()
	}()

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		if err.Error() == "http: request body too large" {
			return nil, fmt.Errorf("request body exceeds maximum allowed size of %d bytes", maxBytes)
		}
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	return body, nil
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

// TODO: Remone completely after confirm is no more needed
// func getLatestKubeVersionFromRegistry(reg *registry.RegistryManager) string {
// 	entities, _, _, _ := reg.GetEntities(&v1beta1.ModelFilter{
// 		Name: "kubernetes",
// 	})

// 	versions := []string{}

// 	for _, entity := range entities {
// 		modelDef, err := utils.Cast[*model.ModelDefinition](entity)
// 		if err != nil {
// 			continue
// 		}
// 		versions = append(versions, modelDef.Model.Version)
// 	}
// 	if len(versions) == 0 {
// 		return ""
// 	}

// 	versions = utils.SortDottedStringsByDigits(versions)

// 	return versions[0]
// }
