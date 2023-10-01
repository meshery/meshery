package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"gopkg.in/yaml.v2"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/patterns"

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

	var input patterns.Pattern
	err = yaml.Unmarshal((body), &input)

	if err != nil {
		http.Error(rw, ErrDecoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}

	for _, svc := range input.Services {
		svc.Settings = patterns.Format.DePrettify(svc.Settings, false)
	}

	data, err := yaml.Marshal(input)
	if err != nil {
		http.Error(rw, models.ErrEncoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}
	// evaluate all the rego policies in the policies directory
	networkPolicy, err := h.Rego.RegoPolicyHandler("data.meshmodel_policy", data)
	if err != nil {
		h.log.Error(ErrResolvingRegoRelationship(err))
		http.Error(rw, ErrResolvingRegoRelationship(err).Error(), http.StatusInternalServerError)
		return
	}

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(networkPolicy)
	if err != nil {
		h.log.Error(models.ErrEncoding(err, "networkPolicy response"))
		http.Error(rw, models.ErrEncoding(err, "networkPolicy response").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/meshmodels/models/{model}/policies/{name} GetMeshmodelPoliciesByName idGetMeshmodelPoliciesByName
// Handle GET request for getting meshmodel policies of a specific model by name.
//
// Example: ```/api/meshmodels/models/kubernetes/policies/{name}```
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
//	200: []meshmodelPoliciesResponseWrapper
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
		if limit == 0 { //If limit is unspecified then it defaults to 25
			limit = DefaultPageSizeForMeshModelComponents
		}
	}
	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	entities, _, _ := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind:      name,
		ModelName: typ,
		Greedy:    greedy,
		Offset:    offset,
		OrderOn:   r.URL.Query().Get("order"),
		Sort:      r.URL.Query().Get("sort"),
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
		pgSize = 0
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    0,
		Policies: policies,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/meshmodels/models/{model}/policies/ GetMeshmodelPolicies idGetMeshmodelPolicies
// Handle GET request for getting meshmodel policies of a specific model by name.
//
// Example: ```/api/meshmodels/models/kubernetes/policies```
//
// // ```?order={field}``` orders on the passed field
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
//	200: []meshmodelPoliciesResponseWrapper
func (h *Handler) GetAllMeshmodelPolicies(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	typ := mux.Vars(r)["model"]

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
	offset := (page - 1) * limit
	if page <= 0 {
		page = 1
	}
	entities, _, _ := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		ModelName: typ,
		Greedy:    greedy,
		Offset:    offset,
		OrderOn:   r.URL.Query().Get("order"),
		Sort:      r.URL.Query().Get("sort"),
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
		pgSize = 0
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    0,
		Policies: policies,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}
