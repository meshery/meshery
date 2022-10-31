package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	meshmodelcore "github.com/layer5io/meshery/server/models/meshmodel/core"
	"github.com/layer5io/meshery/server/models/pattern/core"
	mcore "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
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

func (h *Handler) GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	res := meshmodelcore.GetComponents(h.DbHandler)
	if err := enc.Encode(res); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}
func (h *Handler) GetMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["name"]
	res := meshmodelcore.GetComponentsByName(h.DbHandler, name)
	if err := enc.Encode(res); err != nil {
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
	var cc mcore.ComponentCapability
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = meshmodelcore.RegisterComponentCapability(h.DbHandler, mcore.ComponentCapabilityDBFromCC(cc))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
