package handlers

import (
	"net/http"
	"strconv"

	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/model"
)

const (
	defaultPageSize = 25
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

func getLatestKubeVersionFromRegistry(reg *registry.RegistryManager) string {
	entities, _, _, _ := reg.GetEntities(&v1beta1.ModelFilter{
		Name: "kubernetes",
	})

	versions := []string{}

	for _, entity := range entities {
		modelDef, err := utils.Cast[*model.ModelDefinition](entity)
		if err != nil {
			continue
		}
		versions = append(versions, modelDef.Model.Version)
	}
	if len(versions) == 0 {
		return ""
	}

	versions = utils.SortDottedStringsByDigits(versions)

	return versions[0]
}
