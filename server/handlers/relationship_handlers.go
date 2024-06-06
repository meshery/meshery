package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha2"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	regv1alpha2 "github.com/layer5io/meshkit/models/meshmodel/registry/v1alpha2"
)

// swagger:route GET /api/meshmodels/models/{model}/relationships/{name} GetMeshmodelRelationshipByName idGetMeshmodelRelationshipByName
// Handle GET request for getting meshmodel relationships of a specific model by name.
//
// Example: ```/api/meshmodels/models/kubernetes/relationships/Edge```
//
// # Relationships can be further filtered through query parameter
//
// ```?version={version}```
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
//
//	200: []meshmodelRelationshipsResponseWrapper
func (h *Handler) GetMeshmodelRelationshipByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]
	var greedy bool

	page, _, limit,
		search, order, sort, _ := getPaginationParams(r)
	if search == "true" {
		greedy = true
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	entities, count, _, _ := h.registryManager.GetEntities(&regv1alpha2.RelationshipFilter{
		Version:   r.URL.Query().Get("version"),
		Kind:      name,
		ModelName: typ,
		Greedy:    greedy,
		Limit:     limit,
		Offset:    offset,
		OrderOn:   order,
		Sort:      sort,
	})

	var pgSize int64
	limitstr := r.URL.Query().Get("pagesize")
	if limitstr == "all" {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelRelationshipsAPIResponse{
		Page:          page,
		PageSize:      int(pgSize),
		Count:         count,
		Relationships: entities,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/relationships GetAllMeshmodelRelationships idGetAllMeshmodelRelationships
// Handle GET request for getting all meshmodel relationships
//
// # Relationships can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
//	200: meshmodelRelationshipsResponseWrapper

// swagger:route GET /api/meshmodels/models/{model}/relationships GetAllMeshmodelRelationships idGetAllMeshmodelRelationshipsByModel
// Handle GET request for getting meshmodel relationships of a specific model
//
// Example: ```/api/meshmodel/model/kubernetes/relationship```
//
// # Relationships can be further filtered through query parameter
//
// ```?version={version}```
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
// responses:
//
//	200: meshmodelRelationshipsResponseWrapper
func (h *Handler) GetAllMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	page, _, limit,
		_, order, sort, _ := getPaginationParams(r)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	entities, count, _, _ := h.registryManager.GetEntities(&regv1alpha2.RelationshipFilter{
		Version:   r.URL.Query().Get("version"),
		ModelName: typ,
		Limit:     limit,
		Offset:    offset,
		OrderOn:   order,
		Sort:      sort,
	})

	var pgSize int64
	limitstr := r.URL.Query().Get("pagesize")
	if limitstr == "all" {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelRelationshipsAPIResponse{
		Page:          page,
		PageSize:      int(pgSize),
		Count:         count,
		Relationships: entities,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) RegisterMeshmodelRelationships(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc registry.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	switch cc.EntityType {
	case entity.RelationshipDefinition:
		var r v1alpha2.RelationshipDefinition
		err = json.Unmarshal(cc.Entity, &r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.registryManager.RegisterEntity(cc.Host, &r)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	go h.config.MeshModelSummaryChannel.Publish()
}
