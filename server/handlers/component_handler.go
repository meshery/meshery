package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

// TODO: Swagger API docs
func (h *Handler) ComponentTypesHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := core.ComponentTypesSingleton.Get()
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// TODO: Swagger API docs
func (h *Handler) ComponentVersionsHandler(rw http.ResponseWriter, r *http.Request) {
	t := mux.Vars(r)["type"]
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := core.ComponentTypesSingleton.FilterWorkloadVersionsByType(t)
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// TODO: Swagger API docs
func (h *Handler) ComponentsByNameHandler(rw http.ResponseWriter, r *http.Request) {
	t := mux.Vars(r)["type"]
	v := r.URL.Query().Get("version")
	if v == "" {
		v = "latest"
	}
	n := mux.Vars(r)["name"]
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := core.ComponentTypesSingleton.FilterWorkloadByVersionAndTypeAndName(t, v, n)
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// TODO: Swagger API docs
func (h *Handler) ComponentsForTypeHandler(rw http.ResponseWriter, r *http.Request) {
	t := mux.Vars(r)["type"]
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "latest"
	}
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := core.ComponentTypesSingleton.FilterWorkloadByVersionAndType(t, version)
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// TODO: Swagger API docs
func (h *Handler) GetAllComponents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := core.GetWorkloads()
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

/**Meshmodel endpoints **/
const DefaultPageSizeForMeshModelComponents = 25

// swagger:route GET /api/meshmodel/model GetMeshmodelModels idGetMeshmodelModels
// Handle GET request for getting all meshmodel models. The component type/model name should be lowercase like "kubernetes", "istio"
// ?version={version} If version is unspecified then all models are returned
// ?order={field} orders on the passed field
// ?category={category of model} If category in unspecified then all models are returned
// ?sort={[asc/desc]} Default behavior is asc
// ?search={[true/false]} If search is true then a greedy search is performed
// ?page={page-number} Default page number is 1
// ?pagesize={pagesize} Default pagesize is 25. To return all results: pagesize=all
// responses:
// 200: []Model
func (h *Handler) GetMeshmodelModels(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := r.URL.Query().Get("name")
	v := r.URL.Query().Get("version")
	cat := r.URL.Query().Get("category")
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
	res := h.registryManager.GetModels(&v1alpha1.ModelFilter{
		Name:     name,
		Version:  v,
		Category: cat,
		Limit:    limit,
		Offset:   offset,
		OrderOn:  r.URL.Query().Get("order"),
		Sort:     r.URL.Query().Get("sort"),
	})
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodel/model/{model}/component/{name} MeshmodelGetByName idMeshmodelGetByName
// Handle GET request for getting meshmodel components of a specific type by name.
// Example: /api/meshmodel/model/kubernetes/component/Namespace
// Components can be further filtered through query parameter
// ?version={version} If version is unspecified then all models are returned
// ?order={field} orders on the passed field
// ?sort={[asc/desc]} Default behavior is asc
// ?search={[true/false]} If search is true then a greedy search is performed
// ?page={page-number} Default page number is 1
// ?pagesize={pagesize} Default pagesize is 25. To return all results: pagesize=all
// responses:
// 200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	typ := mux.Vars(r)["model"]
	v := r.URL.Query().Get("version")
	var search bool
	if r.URL.Query().Get("search") == "true" {
		search = true
	}
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
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:      name,
		ModelName: typ,
		Version:   v,
		Greedy:    search,
		Offset:    offset,
		Limit:     limit,
		OrderOn:   r.URL.Query().Get("order"),
		Sort:      r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		m := make(map[string]interface{})
		comp, _ := r.(v1alpha1.ComponentDefinition)
		_ = json.Unmarshal([]byte(comp.Schema), &m)
		m = k8s.Format.Prettify(m, false)
		b, _ := json.Marshal(m)
		comp.Schema = string(b)
		comps = append(comps, comp)
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger: response typesResponseWithModelname
type typesResponseWithModelname struct {
	DisplayName string   `json:"display-name"`
	Versions    []string `json:"versions"`
}

// swagger:route GET /api/meshmodel/model/{model}/component MeshmodelGetByType idMeshmodelGetByType
// Handle GET request for getting meshmodel components of a specific type. The component type/model name should be lowercase like "kubernetes", "istio"
// Example: /api/meshmodel/model/kubernetes/component
// Components can be further filtered through query parameter
// ?version={version}
// ?order={field} orders on the passed field
// ?sort={[asc/desc]} Default behavior is asc
// ?page={page-number} Default page number is 1
// ?pagesize={pagesize} Default pagesize is 25. To return all results: pagesize=all
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
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		ModelName: typ,
		Version:   v,
		Limit:     limit,
		Offset:    offset,
		OrderOn:   r.URL.Query().Get("order"),
		Sort:      r.URL.Query().Get("sort"),
	})
	var comps []v1alpha1.ComponentDefinition
	for _, r := range res {
		m := make(map[string]interface{})
		comp, _ := r.(v1alpha1.ComponentDefinition)
		_ = json.Unmarshal([]byte(comp.Schema), &m)
		m = k8s.Format.Prettify(m, false)
		b, _ := json.Marshal(m)
		comp.Schema = string(b)
		comps = append(comps, comp)
	}
	if err := enc.Encode(comps); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
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
	r.URL.Query()
	var cc meshmodel.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch cc.EntityType {
	case types.ComponentDefinition:
		var c v1alpha1.ComponentDefinition
		err = json.Unmarshal(cc.Entity, &c)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.registryManager.RegisterEntity(cc.Host, c)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

// swagger:response ModelResponse
type ModelResponse struct {
	v1alpha1.Model
	Components    []v1alpha1.ComponentDefinition    `json:"components"`
	Relationships []v1alpha1.RelationshipDefinition `json:"relationships"`
}

// swagger:route GET /api/meshmodel/model/{model} GetMeshmodelEntititiesByModel idGetMeshmodelEntititiesByModel
// Handle GET request for getting all meshmodel entities of a specific model.
// Example: /api/meshmodel/model/kubernetes
// Models can be further filtered through query parameter
// ?version={version}
// ?order={field} orders on the passed field
// ?sort={[asc/desc]} Default behavior is asc
// ?search={[true/false]} If search is true then a greedy search is performed
// ?page={page-number} Default page number is 1
// ?pagesize={pagesize} Default pagesize is 25. To return all results: pagesize=all
// responses:
//
//	200: ModelResponse
func (h *Handler) GetMeshmodelEntititiesByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	v := r.URL.Query().Get("version")
	var greedy bool
	if r.URL.Query().Get("search") == "true" {
		greedy = true
	}
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
	var mres ModelResponse
	mod := h.registryManager.GetModels(&v1alpha1.ModelFilter{
		Name:    typ,
		Version: v,
		Greedy:  greedy,
	})
	if len(mod) != 0 {
		mres.Model = mod[0]
	}
	if mres.Name != "" {
		res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
			ModelName: mres.Name,
			Version:   v,
			Limit:     limit,
			Offset:    offset,
		})
		var comps []v1alpha1.ComponentDefinition
		for _, r := range res {
			m := make(map[string]interface{})
			comp, _ := r.(v1alpha1.ComponentDefinition)
			_ = json.Unmarshal([]byte(comp.Schema), &m)
			m = k8s.Format.Prettify(m, false)
			b, _ := json.Marshal(m)
			comp.Schema = string(b)
			comps = append(comps, comp)
		}
		res2 := h.registryManager.GetEntities(&v1alpha1.RelationshipFilter{
			ModelName: mres.Name,
			Limit:     limit,
			Offset:    offset,
		})
		var relationships []v1alpha1.RelationshipDefinition
		for _, r := range res2 {
			rel, _ := r.(v1alpha1.RelationshipDefinition)

			relationships = append(relationships, rel)
		}
		mres.Relationships = relationships
		mres.Components = comps
	}
	if err := enc.Encode(mres); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
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
