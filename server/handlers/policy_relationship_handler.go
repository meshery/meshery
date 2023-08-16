package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"

	"github.com/layer5io/meshkit/models/meshmodel/core/policies"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/sirupsen/logrus"
)

// swagger:route POST /api/policies/run_policy GetRegoPolicyForDesignFile idGetRegoPolicyForDesignFile
// Handle POST request for running the set of policies on the design file, the policies are picked from the policies directory and query is sent to find all the relationships around the services in the given design file
//
// responses:
// 200
func (h *Handler) GetRegoPolicyForDesignFile(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	_ models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)

		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	// evaluate all the rego policies in the policies directory
	networkPolicy, err := policies.RegoPolicyHandler(context.Background(), []string{"../meshmodel/kubernetes/policies"}, "data.network_policy", body)
	if err != nil {
		h.log.Error(ErrResolvingRegoRelationship(err))
		http.Error(rw, ErrResolvingRegoRelationship(err).Error(), http.StatusInternalServerError)
		return
	}

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(networkPolicy)
	if err != nil {
		h.log.Error(ErrEncoding(err, "networkPolicy response"))
		http.Error(rw, ErrEncoding(err, "networkPolicy response").Error(), http.StatusInternalServerError)
		return
	}
}


func (h *Handler) GetAllMeshmodelPoliciesByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]
	var greedy bool
	if r.URL.Query().Get("search") == "true" {
		greedy = true
	}
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	entities, count, _ := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind: name,
		SubType: typ,
		Greedy: greedy,
		Limit: limit,
		Offset: offset,
		OrderOn: r.URL.Query().Get("order"),
	})
	var policies []v1alpha1.PolicyDefinition
	for _, p := range entities {
		policy, ok := p.(v1alpha1.PolicyDefinition)
		if ok {
			policies = append(policies, policy)
		} 
	}

	var pgSize int64
	if limitstr == "all" {
		pgSize = *count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:          page,
		PageSize:      int(pgSize),
		Count:         *count,
		Policies: 	   policies,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) GetAllMeshmodelPolicies(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]
	var greedy bool
	if r.URL.Query().Get("search") == "true" {
		greedy = true
	}
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	entities, count, _ := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind: name,
		SubType: typ,
		Greedy: greedy,
		Limit: limit,
		Offset: offset,
		OrderOn: r.URL.Query().Get("order"),
	})
	var policies []v1alpha1.PolicyDefinition
	for _, p := range entities {
		policy, ok := p.(v1alpha1.PolicyDefinition)
		if ok {
			policies = append(policies, policy)
		} 
	}

	var pgSize int64
	if limitstr == "all" {
		pgSize = *count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:          page,
		PageSize:      int(pgSize),
		Count:         *count,
		Policies: 	   policies,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}