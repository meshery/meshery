package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"gopkg.in/yaml.v2"

	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"

	"github.com/sirupsen/logrus"
)

const (
	relationshipPolicyPackageName = "data.meshmodel_policy"
	suffix                        = "_relationship"
)

type relationshipPolicyEvalPayload struct {
	PatternFile       string   `json:"pattern_file"`
	EvaluationQueries []string `json:"evaluation_queries"`
}

// swagger:route POST /api/meshmodels/relationships/evaluate EvaluateRelationshipPolicy relationshipPolicyEvalPayloadWrapper
// Handle POST request for evaluating relationships in the provided design file by running a set of provided evaluation queries on the design file
//
// responses:
// 200
func (h *Handler) EvaluateRelationshipPolicy(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	userUUID := uuid.FromStringOrNil(user.ID)
	defer func() {
		_ = r.Body.Close()
	}()

	eventBuilder := events.NewEvent().FromSystem(*h.SystemID).FromUser(userUUID).WithCategory("relationship").WithAction("evaluation")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	relationshipPolicyEvalPayload := relationshipPolicyEvalPayload{}
	err = json.Unmarshal(body, &relationshipPolicyEvalPayload)

	if err != nil {
		http.Error(rw, ErrDecoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}
	var patternFile core.Pattern

	err = yaml.Unmarshal([]byte(relationshipPolicyEvalPayload.PatternFile), &patternFile)
	if err != nil {
		http.Error(rw, ErrDecoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}

	evaluationQueries := relationshipPolicyEvalPayload.EvaluationQueries

	for _, svc := range patternFile.Services {
		svc.Settings = core.Format.DePrettify(svc.Settings, false)
	}

	data, err := yaml.Marshal(patternFile)
	if err != nil {
		http.Error(rw, models.ErrEncoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}

	patternUUID := uuid.FromStringOrNil(patternFile.PatternID)
	eventBuilder.ActedUpon(patternUUID)

	var evalResults interface{}

	// evaluate specified relationship policies
	verifiedEvaluationQueries := h.verifyEvaluationQueries(evaluationQueries)
	if len(verifiedEvaluationQueries) == 0 {
		event := eventBuilder.WithDescription("Invalid or unsupported evaluation queries provided").WithSeverity(events.Error).WithMetadata(map[string]interface{}{"evaluationQueries": evaluationQueries}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userUUID, event)
		return
	}

	evalresults := make(map[string]interface{}, 0)
	for _, query := range verifiedEvaluationQueries {
		result, err := h.Rego.RegoPolicyHandler(fmt.Sprintf("%s.%s", relationshipPolicyPackageName, query), data)
		if err != nil {
			h.log.Debug(err)
			continue
		}
		evalresults[query] = result
	}
	// Before starting the eval the design is de-prettified, so that we can use the relationships def correctly.
	// The results contain the updated config.
	// Prettify the results before sending it to client.
	evalResults = core.Format.Prettify(evalresults, false)

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(evalResults)
	if err != nil {
		h.log.Error(models.ErrEncoding(err, "policy evaluation response"))
		http.Error(rw, models.ErrEncoding(err, "failed to generate policy evaluation results").Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) verifyEvaluationQueries(evaluationQueries []string) (verifiedEvaluationQueries []string) {
	registeredRelationships, _, _ := h.registryManager.GetEntities(&v1alpha1.RelationshipFilter{})

	var relationships []v1alpha1.RelationshipDefinition
	for _, entity := range registeredRelationships {
		relationship, err := utils.Cast[v1alpha1.RelationshipDefinition](entity)
		if err != nil {
			return
		}
		relationships = append(relationships, relationship)
	}

	if len(evaluationQueries) == 0 || (len(evaluationQueries) == 1 && evaluationQueries[0] == "all") {
		for _, relationship := range relationships {
			if relationship.EvaluationQuery != "" {
				verifiedEvaluationQueries = append(verifiedEvaluationQueries, relationship.EvaluationQuery)
			}
		}
	} else {
		for _, regoQuery := range evaluationQueries {
			for _, relationship := range relationships {
				if strings.TrimSuffix(regoQuery, suffix) == fmt.Sprintf("%s_%s", strings.ToLower(relationship.Kind), strings.ToLower(relationship.SubType)) {
					verifiedEvaluationQueries = append(verifiedEvaluationQueries, relationship.EvaluationQuery)
					break
				}
			}
		}
	}
	return
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
