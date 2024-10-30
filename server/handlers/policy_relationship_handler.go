package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"

	"github.com/layer5io/meshkit/models/events"

	"github.com/layer5io/meshkit/models/meshmodel/registry"
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
	unknownComponents := processEvaluationResponse(h.registryManager, relationshipPolicyEvalPayload, &evaluationResponse)
	if len(unknownComponents) > 0 {
		event := events.NewEvent().FromUser(userUUID).FromSystem(*h.SystemID).WithCategory("relationship").WithAction("evaluation").WithSeverity(events.Informational).ActedUpon(patternUUID).WithDescription(fmt.Sprintf("Relationship evaluation for \"%s\" at version \"%s\" resulted in the addition of new components but they are not registered inside the registry.", evaluationResponse.Design.Name, evaluationResponse.Design.Version)).WithMetadata(map[string]interface{}{
			"ComponentsToBeAdded": unknownComponents,
		}).Build()

		_ = provider.PersistEvent(event)
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

func processEvaluationResponse(registry *registry.RegistryManager, evalPayload pattern.EvaluationRequest, evalResponse *pattern.EvaluationResponse) []*component.ComponentDefinition {
	compsUpdated := []component.ComponentDefinition{}
	compsAdded := []component.ComponentDefinition{}

	// components which were added by the evaluator based on the relationship definition, but doesn't exist in the registry.
	unknownComponents := []*component.ComponentDefinition{}

	for _, cmp := range evalResponse.Trace.ComponentsAdded {
		_c := cmp

		compFilter := &regv1beta1.ComponentFilter{
			Name:       _c.Component.Kind,
			APIVersion: _c.Component.Version,
			Version:    _c.Model.Model.Version,
			ModelName:  _c.Model.Name,
			Limit:      1,
			Trim:       true,
		}
		if _c.Model.Model.Version == "*" {
			compFilter.Version = ""
		}
		entities, _, _, _ := registry.GetEntities(compFilter)
		if len(entities) == 0 {
			unknownComponents = append(unknownComponents, &_c)
			continue
		}
		_component, _ := entities[0].(*component.ComponentDefinition)

		_c.Configuration = core.Format.Prettify(_c.Configuration, false)
		_component.Id = _c.Id
		if _c.DisplayName != "" {
			_component.DisplayName = _c.DisplayName
		} else {
			_component.DisplayName = fmt.Sprintf("%s-%s", strings.ToLower(_component.DisplayName), utils.GetRandomAlphabetsOfDigit(3))
		}
		_component.Configuration = _c.Configuration
		compsAdded = append(compsAdded, *_component)
	}

	evalResponse.Trace.ComponentsAdded = compsAdded

	for _, component := range evalResponse.Trace.ComponentsUpdated {
		_c := component

		_c.Configuration = core.Format.Prettify(_c.Configuration, false)
		compsUpdated = append(compsUpdated, _c)
	}

	evalResponse.Trace.ComponentsUpdated = compsUpdated

	cmps := append(compsAdded, compsUpdated...)

	if evalPayload.Options != nil && evalPayload.Options.ReturnDiffOnly != nil && *evalPayload.Options.ReturnDiffOnly {
		evalResponse.Design.Relationships = []*relationship.RelationshipDefinition{}
		evalResponse.Design.Components = []*component.ComponentDefinition{}

		for _, relationship := range evalResponse.Trace.RelationshipsAdded {
			_r := relationship
			evalResponse.Design.Relationships = append(evalResponse.Design.Relationships, &_r)
		}
		for _, relationship := range evalResponse.Trace.RelationshipsRemoved {
			_r := relationship
			evalResponse.Design.Relationships = append(evalResponse.Design.Relationships, &_r)
		}

		for _, relationship := range evalResponse.Trace.RelationshipsUpdated {
			_r := relationship
			evalResponse.Design.Relationships = append(evalResponse.Design.Relationships, &_r)
		}

		for _, cmp := range cmps {
			evalResponse.Design.Components = append(evalResponse.Design.Components, &cmp)
		}
		return unknownComponents
	}

	designComponents := evalResponse.Design.Components
	evalResponse.Design.Components = []*component.ComponentDefinition{}

	for _, cmp := range designComponents {
		_c := cmp

		for _, c := range cmps {
			if c.Id == _c.Id {
				_c = &c
				break
			}
		}

		_c.Configuration = core.Format.Prettify(_c.Configuration, false)
		evalResponse.Design.Components = append(evalResponse.Design.Components, _c)

	}

	return unknownComponents
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
