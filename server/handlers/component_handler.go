package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

/**Meshmodel endpoints **/
const DefaultPageSizeForMeshModelComponents = 25

// swagger:route GET /api/meshmodels/categories/{category}/models GetMeshmodelModelsByCategories idGetMeshmodelModelsByCategories
// Handle GET request for getting all meshmodel models for a given category. The component type/model name should be lowercase like "kubernetes", "istio"
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={modelname}``` If search is non empty then a greedy search is performed
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []Model
func (h *Handler) GetMeshmodelModelsByCategories(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	cat := mux.Vars(r)["category"]
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ModelFilter{
		Category: cat,
		Limit:    limit,
		Offset:   offset,
		OrderOn:  r.URL.Query().Get("order"),
		Sort:     r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.Name = r.URL.Query().Get("search")
	}
	res, _ := h.registryManager.GetModels(h.dbHandler, filter)

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category}/models/{model} GetMeshmodelModelsByCategoriesByModel idGetMeshmodelModelsByCategoriesByModel
// Handle GET request for getting all meshmodel models for a given category. The component type/model name should be lowercase like "kubernetes", "istio"
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []Model
func (h *Handler) GetMeshmodelModelsByCategoriesByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	cat := mux.Vars(r)["category"]
	model := mux.Vars(r)["model"]
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetModels(h.dbHandler, &v1alpha1.ModelFilter{
		Category: cat,
		Name:     model,
		Version:  r.URL.Query().Get("version"),
		Limit:    limit,
		Offset:   offset,
		OrderOn:  r.URL.Query().Get("order"),
		Sort:     r.URL.Query().Get("sort"),
	})

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET ```/api/meshmodels/models``` GetMeshmodelModels idGetMeshmodelModels
// Handle GET request for getting all meshmodel models and their total count.
//
// # Returns a list of registered models across all categories
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={modelname}``` If search is non empty then a greedy search is performed
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
//
//	200: meshmodelModelsResponseWrapper
func (h *Handler) GetMeshmodelModels(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ModelFilter{
		Version: v,
		Limit:   limit,
		Offset:  offset,
		OrderOn: r.URL.Query().Get("order"),
		Sort:    r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.DisplayName = r.URL.Query().Get("search")
		filter.Greedy = true
	}
	models, count := h.registryManager.GetModels(h.dbHandler, filter)

	res := struct {
		Count  int64            `json:"total_count"`
		Models []v1alpha1.Model `json:"models"`
	}{
		Count:  count,
		Models: models,
	}

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET ```/api/meshmodels/models/{model}``` GetMeshmodelModelsByName idGetMeshmodelModelsByName
// Handle GET request for getting all meshmodel models. The component type/model name should be lowercase like "kubernetes", "istio"
// Returns a list of registered models across all categories
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []Model
func (h *Handler) GetMeshmodelModelsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["model"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetModels(h.dbHandler, &v1alpha1.ModelFilter{
		Name:    name,
		Version: v,
		Limit:   limit,
		Offset:  offset,
		OrderOn: r.URL.Query().Get("order"),
		Sort:    r.URL.Query().Get("sort"),
	})

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories GetMeshmodelCategories idGetMeshmodelCategories
// Handle GET request for getting all meshmodel categories
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?search={categoryName}``` If search is non empty then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []Category
func (h *Handler) GetMeshmodelCategories(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.CategoryFilter{
		Limit:   limit,
		Offset:  offset,
		OrderOn: r.URL.Query().Get("order"),
		Sort:    r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.Name = r.URL.Query().Get("search")
	}
	res := h.registryManager.GetCategories(h.dbHandler, filter)

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category} GetMeshmodelCategoriesByName idGetMeshmodelCategoriesByName
// Handle GET request for getting all meshmodel categories of a given name
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
// responses:
// 200: []Category
func (h *Handler) GetMeshmodelCategoriesByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["category"]
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res := h.registryManager.GetCategories(h.dbHandler, &v1alpha1.CategoryFilter{
		Name:    name,
		Limit:   limit,
		Offset:  offset,
		OrderOn: r.URL.Query().Get("order"),
		Sort:    r.URL.Query().Get("sort"),
	})

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category}/models/{model}/components/{name} GetMeshmodelComponentsByNameByModelByCategory idGetMeshmodelComponentsByNameByModelByCategory
// Handle GET request for getting meshmodel components of a specific type by model and category.
//
// Example: ```/api/meshmodels/categories/Orchestration``` and Management/models/kubernetes/components/Namespace
// Components can be further filtered through query parameter
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByNameByModelByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	typ := mux.Vars(r)["model"]
	cat := mux.Vars(r)["category"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:         name,
		CategoryName: cat,
		ModelName:    typ,
		APIVersion:   r.URL.Query().Get("apiVersion"),
		Version:      v,
		Offset:       offset,
		Limit:        limit,
		OrderOn:      r.URL.Query().Get("order"),
		Sort:         r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category}/components/{name} GetMeshmodelComponentsByNameByCategory idGetMeshmodelComponentsByNameByCategory
// Handle GET request for getting meshmodel components of a specific type category.
//
// Example: ```/api/meshmodels/categories/Orchestration``` and Management/components/Namespace
// Components can be further filtered through query parameter
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByNameByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	cat := mux.Vars(r)["category"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:         name,
		CategoryName: cat,
		APIVersion:   r.URL.Query().Get("apiVersion"),
		Version:      v,
		Offset:       offset,
		Limit:        limit,
		OrderOn:      r.URL.Query().Get("order"),
		Sort:         r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/models/{model}/components/{name} GetMeshmodelComponentsByNameByModel idGetMeshmodelComponentsByNameByModel
// Handle GET request for getting meshmodel components of a specific  model.
//
// Example: ```/api/meshmodels/models/kubernetes/components/Namespace```
// Components can be further filtered through query parameter
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
//
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByNameByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	typ := mux.Vars(r)["model"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:       name,
		ModelName:  typ,
		APIVersion: r.URL.Query().Get("apiVersion"),
		Version:    v,
		Offset:     offset,
		Limit:      limit,
		OrderOn:    r.URL.Query().Get("order"),
		Sort:       r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/components/{name} GetAllMeshmodelComponentsByName idGetAllMeshmodelComponentsByName
// Handle GET request for getting meshmodel components of a specific type by name across all models and categories
//
// Example: ```/api/meshmodels/components/Namespace```
// Components can be further filtered through query parameter
//
// ```?version={version}``` If version is unspecified then all models are returned
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?trim={[true]}``` When trim is set to true, the underlying schemas are not returned for entities
//
// ```?search={[true/false]}``` If search is true then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetAllMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	res, _ := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:       name,
		Trim:       r.URL.Query().Get("trim") == "true",
		APIVersion: r.URL.Query().Get("apiVersion"),
		Version:    v,
		Offset:     offset,
		Limit:      limit,
		OrderOn:    r.URL.Query().Get("order"),
		Sort:       r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/models/{model}/components GetMeshmodelComponentByModel idGetMeshmodelComponentByModel
// Handle GET request for getting meshmodel components of a specific model. The component type/model name should be lowercase like "kubernetes", "istio"
//
// Example: ```/api/meshmodels/models/kubernetes/components```
// Components can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?trim={[true]}``` When trim is set to true, the underlying schemas are not returned for entities
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?search={componentname}``` If search is non empty then a greedy search is performed
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ComponentFilter{
		ModelName:  typ,
		Version:    v,
		Trim:       r.URL.Query().Get("trim") == "true",
		APIVersion: r.URL.Query().Get("apiVersion"),
		Limit:      limit,
		Offset:     offset,
		OrderOn:    r.URL.Query().Get("order"),
		Sort:       r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.Name = r.URL.Query().Get("search")
	}
	res, _ := h.registryManager.GetEntities(filter)
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category}/models/{model}/components GetMeshmodelComponentByModelByCategory idGetMeshmodelComponentByModelByCategory
// Handle GET request for getting meshmodel components of a specific model and category. The component type/model name should be lowercase like "kubernetes", "istio"
//
// Example: ```/api/meshmodels/categories/Orchestration``` and Management/models/kubernetes/components
// Components can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?trim={[true]}``` When trim is set to true, the underlying schemas are not returned for entities
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={componentname}``` If search is non empty then a greedy search is performed
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentByModelByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	cat := mux.Vars(r)["category"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ComponentFilter{
		CategoryName: cat,
		ModelName:    typ,
		Version:      v,
		Trim:         r.URL.Query().Get("trim") == "true",
		APIVersion:   r.URL.Query().Get("apiVersion"),
		Limit:        limit,
		Offset:       offset,
		OrderOn:      r.URL.Query().Get("order"),
		Sort:         r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.Name = r.URL.Query().Get("search")
	}
	res, _ := h.registryManager.GetEntities(filter)
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/categories/{category}/components GetMeshmodelComponentByCategory idGetMeshmodelComponentByCategory
// Handle GET request for getting meshmodel components of a specific model and category.
//
// Example: ```/api/meshmodels/categories/Orchestration``` and Management/models/kubernetes/components
// Components can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?trim={[true]}``` When trim is set to true, the underlying schemas are not returned for entities
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={componentname}``` If search is non empty then a greedy search is performed
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	cat := mux.Vars(r)["category"]
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ComponentFilter{
		CategoryName: cat,
		Version:      v,
		Trim:         r.URL.Query().Get("trim") == "true",
		APIVersion:   r.URL.Query().Get("apiVersion"),
		Limit:        limit,
		Offset:       offset,
		OrderOn:      r.URL.Query().Get("order"),
		Sort:         r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.Name = r.URL.Query().Get("search")
	}
	res, _ := h.registryManager.GetEntities(filter)
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/components GetAllMeshmodelComponents idGetAllMeshmodelComponents
// Handle GET request for getting meshmodel components across all models and categories and their total count.
//
// # Components can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?order={field}``` orders on the passed field
//
// ```?search={componentname}``` If search is non empty then a greedy search is performed
//
// ```?trim={[true]}``` When trim is set to true, the underlying schemas are not returned for entities
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
//  200: allMeshmodelComponentsResponseWrapper

func (h *Handler) GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	v := r.URL.Query().Get("version")
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}

	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit
	filter := &v1alpha1.ComponentFilter{
		Version:    v,
		Trim:       r.URL.Query().Get("trim") == "true",
		APIVersion: r.URL.Query().Get("apiVersion"),
		Limit:      limit,
		Offset:     offset,
		OrderOn:    r.URL.Query().Get("order"),
		Sort:       r.URL.Query().Get("sort"),
	}
	if r.URL.Query().Get("search") != "" {
		filter.Greedy = true
		filter.DisplayName = r.URL.Query().Get("search")
	}
	entities, count := h.registryManager.GetEntities(filter)
	var comps []v1alpha1.ComponentDefinition
	for _, r := range entities {
		comp, ok := r.(v1alpha1.ComponentDefinition)
		if ok {
			m := make(map[string]interface{})
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = core.Format.Prettify(m, true)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
	}

	res := struct {
		Count      int64                          `json:"total_count"`
		Components []v1alpha1.ComponentDefinition `json:"components"`
	}{
		Count:      *count,
		Components: comps,
	}

	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrGetMeshModels(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route POST /api/meshmodel/components/register MeshmodelValidate idPostMeshModelValidate
// Handle POST request for registering meshmodel components.
//
// Validate the given value with the given schema
// responses:
// 	200:

// request body should be json
// request body should be of ComponentCapability format
func (h *Handler) RegisterMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc meshmodel.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var c v1alpha1.ComponentDefinition
	switch cc.EntityType {
	case types.ComponentDefinition:
		err = json.Unmarshal(cc.Entity, &c)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		utils.WriteSVGsOnFileSystem(&c)
		err = h.registryManager.RegisterEntity(cc.Host, c)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	go h.config.MeshModelSummaryChannel.Publish()
}

func filterUniqueElementsArray(s []string) []string {
	m := make(map[string]bool)
	for _, ele := range s {
		m[ele] = true
	}
	ans := make([]string, 0)
	for a := range m {
		ans = append(ans, a)
	}
	return ans
}
