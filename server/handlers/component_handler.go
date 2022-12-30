package handlers

import (
	"encoding/json"
	"net/http"

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

// swagger:route GET /api/meshmodel/components MeshmodelGet idMeshmodelGet
// Handle GET request for getting all meshmodel components.
// Components can be further filtered through query parameter ?version=
// responses:
//
//	200: []ComponentDefinition
func (h *Handler) GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{})
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

// swagger:route GET /api/meshmodel/components/{type}/{name} MeshmodelGetByName idMeshmodelGetByName
// Handle GET request for getting meshmodel components of a specific type by name.
// Example: /api/meshmodel/components/kubernetes/Namespace
// Components can be further filtered through query parameter ?version=
// responses:
//
//	200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	typ := mux.Vars(r)["type"]
	v := r.URL.Query().Get("version")
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		Name:      name,
		ModelName: typ,
		Version:   v,
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

// swagger: response responseTypesWithModelName
type responseTypesWithModelName map[string]*typesResponseWithModelname

// swagger:route GET /api/meshmodel/components/types MeshmodelComponentsForTypeHandler idMeshmodelComponentsForTypeHandler
// Handle GET request for getting meshmodel types or model names.
// response:
//
//	200: responseTypesWithModelName
func (h *Handler) MeshmodelComponentsForTypeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{})
	var response = make(responseTypesWithModelName)
	for _, r := range res {
		def, _ := r.(v1alpha1.ComponentDefinition)
		if response[def.Model.Name] == nil {
			response[def.Model.Name] = &typesResponseWithModelname{
				DisplayName: def.Model.DisplayName,
				Versions:    []string{def.Model.Version},
			}
		} else {
			response[def.Model.Name].Versions = append(response[def.Model.Name].Versions, def.Model.Version)
		}
	}
	for _, x := range response {
		x.Versions = filterUniqueElementsArray(x.Versions)
	}
	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodel/components/{type} MeshmodelGetByType idMeshmodelGetByType
// Handle GET request for getting meshmodel components of a specific type. The component type/model name should be lowercase like "kubernetes", "istio"
// Example: /api/meshmodel/components/kubernetes
// Components can be further filtered through query parameter ?version=
// responses:
//
//	200: []ComponentDefinition
func (h *Handler) GetMeshmodelComponentsByType(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["type"]
	v := r.URL.Query().Get("version")
	res := h.registryManager.GetEntities(&v1alpha1.ComponentFilter{
		ModelName: typ,
		Version:   v,
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
