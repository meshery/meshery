package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"

	"github.com/layer5io/meshkit/models/events"

	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
)

const (
	relationshipPolicyPackageName = "data.relationship_evaluation_policy"
	suffix                        = "_relationship"
)

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
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	relationshipPolicyEvalPayload := pattern.EvaluationRequest{}
	err = json.Unmarshal(body, &relationshipPolicyEvalPayload)

	if err != nil {
		http.Error(rw, ErrDecoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}
	// decode the pattern file

	patternUUID := relationshipPolicyEvalPayload.Design.Id
	eventBuilder.ActedUpon(patternUUID)

	for _, component := range relationshipPolicyEvalPayload.Design.Components {
		component.Configuration = core.Format.DePrettify(component.Configuration, false)
	}

	// evaluate specified relationship policies
	// on successful eval the event containing details like comps evaulated, relationships indeitified should be emitted and peristed.
	evaluationResponse, err := h.Rego.RegoPolicyHandler(relationshipPolicyEvalPayload.Design,
		relationshipPolicyPackageName,
	)
	if err != nil {
		h.log.Debug(err)
		// log an event
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	currentTime := time.Now()
	evaluationResponse.Timestamp = &currentTime
	// include trace instead of design file
	event := eventBuilder.WithDescription(fmt.Sprintf("Relationship evaluation completed for \"%s\" at version \"%s\"", evaluationResponse.Design.Name, evaluationResponse.Design.Version)).
		WithMetadata(map[string]interface{}{
			"trace":        evaluationResponse.Trace,
			"evaluated_at": *evaluationResponse.Timestamp,
		}).WithSeverity(events.Informational).Build()
	_ = provider.PersistEvent(event)

	// Create the event but do not notify the client immediately, as the evaluations are frequent and takes up the view area.
	// go h.config.EventBroadcaster.Publish(userUUID, event)

	if relationshipPolicyEvalPayload.Options != nil && relationshipPolicyEvalPayload.Options.ReturnDiffOnly != nil &&
		*relationshipPolicyEvalPayload.Options.ReturnDiffOnly {
		compsUpdated := []component.ComponentDefinition{}
		evaluationResponse.Design.Components = []*component.ComponentDefinition{}
		evaluationResponse.Design.Relationships = []*relationship.RelationshipDefinition{}
		for _, component := range evaluationResponse.Trace.ComponentsUpdated {
			_c := component

			_c.Configuration = core.Format.Prettify(_c.Configuration, false)
			evaluationResponse.Design.Components = append(evaluationResponse.Design.Components, &_c)
			compsUpdated = append(compsUpdated, _c)
		}

		for _, relationship := range evaluationResponse.Trace.RelationshipsAdded {
			_r := relationship
			evaluationResponse.Design.Relationships = append(evaluationResponse.Design.Relationships, &_r)
		}
		for _, relationship := range evaluationResponse.Trace.RelationshipsRemoved {
			_r := relationship
			evaluationResponse.Design.Relationships = append(evaluationResponse.Design.Relationships, &_r)
		}

		for _, relationship := range evaluationResponse.Trace.RelationshipsUpdated {
			_r := relationship
			evaluationResponse.Design.Relationships = append(evaluationResponse.Design.Relationships, &_r)
		}

		evaluationResponse.Trace.ComponentsUpdated = compsUpdated

		ec := json.NewEncoder(rw)
		err = ec.Encode(evaluationResponse)
		if err != nil {
			h.log.Error(models.ErrEncoding(err, "policy evaluation response"))
			http.Error(rw, models.ErrEncoding(err, "failed to generate policy evaluation results").Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Before starting the eval the design is de-prettified, so that we can use the relationships def correctly.
	// The results contain the updated config.
	// Prettify the design before sending it to client.

	for _, component := range evaluationResponse.Design.Components {
		component.Configuration = core.Format.Prettify(component.Configuration, false)
	}

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(evaluationResponse)
	if err != nil {
		h.log.Error(models.ErrEncoding(err, "policy evaluation response"))
		http.Error(rw, models.ErrEncoding(err, "failed to generate policy evaluation results").Error(), http.StatusInternalServerError)
		return
	}
}

// unused currently

// func (h *Handler) verifyEvaluationQueries(evaluationQueries []string) (verifiedEvaluationQueries []string) {
// 	registeredRelationships, _, _, _ := h.registryManager.GetEntities(&regv1alpha3.RelationshipFilter{})

// 	var relationships []relationship.RelationshipDefinition
// 	for _, entity := range registeredRelationships {
// 		relationship, err := mutils.Cast[*relationship.RelationshipDefinition](entity)

// 		if err != nil {
// 			return
// 		}
// 		relationships = append(relationships, *relationship)
// 	}

// 	if len(evaluationQueries) == 0 || (len(evaluationQueries) == 1 && evaluationQueries[0] == "all") {
// 		for _, relationship := range relationships {
// 			if relationship.EvaluationQuery != nil {
// 				verifiedEvaluationQueries = append(verifiedEvaluationQueries, *relationship.EvaluationQuery)
// 			} else {
// 				verifiedEvaluationQueries = append(verifiedEvaluationQueries, relationship.GetDefaultEvaluationQuery())
// 			}
// 		}
// 	} else {
// 		for _, regoQuery := range evaluationQueries {
// 			for _, relationship := range relationships {
// 				if (relationship.EvaluationQuery != nil && regoQuery == *relationship.EvaluationQuery) || regoQuery == relationship.GetDefaultEvaluationQuery() {
// 					verifiedEvaluationQueries = append(verifiedEvaluationQueries, *relationship.EvaluationQuery)
// 					break
// 				}
// 			}
// 		}
// 	}
// 	return
// }

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
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]
	var greedy bool
	if search == "true" {
		greedy = true
	}

	entities, _, _, _ := h.registryManager.GetEntities(&regv1beta1.PolicyFilter{
		Kind:      name,
		ModelName: typ,
		Greedy:    greedy,
		Offset:    offset,
		OrderOn:   order,
		Sort:      sort,
	})

	var pgSize int64
	if limit == 0 {
		pgSize = 0
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    0,
		Policies: entities,
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
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]

	var greedy bool
	if search == "true" {
		greedy = true
	}

	entities, _, _, _ := h.registryManager.GetEntities(&regv1beta1.PolicyFilter{
		ModelName: typ,
		Greedy:    greedy,
		Offset:    offset,
		OrderOn:   order,
		Sort:      sort,
	})

	var pgSize int64

	if limit == 0 {
		pgSize = 0
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelPoliciesAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    0,
		Policies: entities,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrWorkloadDefinition(err)) //TODO: Add appropriate meshkit error
		http.Error(rw, ErrWorkloadDefinition(err).Error(), http.StatusInternalServerError)
	}
}
