package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/types"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

// swagger:route GET /api/meshmodel/model/{model}/policy/{name} GetMeshmodelPolicyByName idGetMeshmodelPolicyByName
// Handle GET request for getting meshmodel policies of a specific model by name.
// Example: /api/meshmodel/model/kubernetes/policy/Validation
// responses:
// 200: []PolicyDefinition
func (h *Handler) GetMeshmodelPolicyByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]

	res := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind:      name,
		ModelName: typ,
	})

	var pls []v1alpha1.PolicyDefinition
	for _, p := range res {
		pl, ok := p.(v1alpha1.PolicyDefinition)
		if ok {
			pls = append(pls, pl)
		}
	}
	if err := enc.Encode(pls); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodel/model/{model}/policy GetAllMeshmodelPolicy idGetAllMeshmodelPolicy
// Handle GET request for getting meshmodel policies of a specific model.
// Example: /api/meshmodel/model/kubernetes/policy
// responses:
// 200: []PolicyDefinition

func (h *Handler) GetAllMeshmodelPolicy(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]

	res := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		ModelName: typ,
	})
	var pls []v1alpha1.PolicyDefinition
	for _, p := range res {
		pl, ok := p.(v1alpha1.PolicyDefinition)
		if ok {
			pls = append(pls, pl)
		}
	}
	if err := enc.Encode(pls); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterMeshmodelPolicy(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc meshmodel.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch cc.EntityType {
	case types.PolicyDefinition:
		var p v1alpha1.PolicyDefinition
		err = json.Unmarshal(cc.Entity, &p)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.registryManager.RegisterEntity(cc.Host, p)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	go h.config.MeshModelSummaryChannel.Publish()
}
