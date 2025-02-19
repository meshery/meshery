package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"

	"github.com/layer5io/meshkit/models/events"

	"github.com/layer5io/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	mutils "github.com/layer5io/meshkit/utils"
)

const (
	//TODO: Might needs to be made dyanamic to support dynamically loading policies based on eval request.
	RelationshipPolicyPackageName = "data.relationship_evaluation_policy"
	suffix                        = "_relationship"
)

const RELATIONSHIP_SUBTYPE_ALIAS = "alias"

// Aliasses Are not resolved
func parseRelationshipToAlias(relationshipDeclaration relationship.RelationshipDefinition) (pattern.NonResolvedAlias, bool) {

	alias := pattern.NonResolvedAlias{}

	if relationshipDeclaration.SubType != RELATIONSHIP_SUBTYPE_ALIAS {
		return alias, false
	}

	selectors := *relationshipDeclaration.Selectors

	if len(selectors) == 0 {
		return alias, false
	}

	selector := selectors[0]
	fromSet := selector.Allow.From
	toSet := selector.Allow.To

	if len(fromSet) == 0 || len(toSet) == 0 {
		return alias, false
	}

	from := fromSet[0]
	to := toSet[0]
	if from.Patch == nil || from.Patch.MutatedRef == nil {
		return alias, false
	}
	mutatedRefs := *from.Patch.MutatedRef

	if len(mutatedRefs) == 0 {
		return alias, false
	}

	alias.ImmediateParentId = *to.Id
	alias.AliasComponentId = *from.Id
	alias.RelationshipId = relationshipDeclaration.Id
	alias.ImmediateRefFieldPath = mutatedRefs[0]

	return alias, true

}

func ParseComponentToAlias(component component.ComponentDefinition, relationships []*relationship.RelationshipDefinition) (pattern.NonResolvedAlias, bool) {

	for _, relationship := range relationships {
		alias, ok := parseRelationshipToAlias(*relationship)
		if !ok {
			continue
		}

		if alias.AliasComponentId == component.Id {
			return alias, true
		}
	}

	return pattern.NonResolvedAlias{}, false
}

// getComponentById retrieves a component from the design by its ID
func getComponentById(design pattern.PatternFile, id uuid.UUID) *component.ComponentDefinition {
	for _, comp := range design.Components {
		if comp.Id == id {
			return comp
		}
	}
	return nil
}

func ResolveAlias(nonResolvedAlias pattern.NonResolvedAlias, currentNonResolved pattern.NonResolvedAlias, path []string, design pattern.PatternFile) pattern.ResolvedAlias {
	parentComponent := getComponentById(design, currentNonResolved.ImmediateParentId)
	if parentComponent == nil {
		return pattern.ResolvedAlias{
			NonResolvedAlias:     nonResolvedAlias,
			ResolvedParentId:     currentNonResolved.ImmediateParentId,
			ResolvedRefFieldPath: path,
		}
	}

	parentAlias, ok := ParseComponentToAlias(*parentComponent, design.Relationships)

	if !ok {

		return pattern.ResolvedAlias{
			NonResolvedAlias:     nonResolvedAlias,
			ResolvedParentId:     currentNonResolved.ImmediateParentId,
			ResolvedRefFieldPath: path,
		}
	}

	// slicing from 1 to remove "configuration" prefix when building the resolved ref
	// so we dont get something like configuration,spec,configuration , containers , _
	// appending to aprentAlias.ImmediateReffiled , than path , because this a recursive function it will otherwise build the path in reverse
	return ResolveAlias(nonResolvedAlias, parentAlias, append(parentAlias.ImmediateRefFieldPath, path[1:]...), design)
}

func ResolveAliasesInDesign(design pattern.PatternFile) map[string]pattern.ResolvedAlias {

	resolvedAliases := make(map[string]pattern.ResolvedAlias)

	for _, relationship := range design.Relationships {
		nonResolvedalias, ok := parseRelationshipToAlias(*relationship)
		if ok {
			resolvedAlias := ResolveAlias(nonResolvedalias, nonResolvedalias, nonResolvedalias.ImmediateRefFieldPath, design)
			resolvedAliases[resolvedAlias.AliasComponentId.String()] = resolvedAlias
		}
	}

	return resolvedAliases

}

// Helper method to make design evaluation based on the relationship policies.
func (h *Handler) EvaluateDesign(
	relationshipPolicyEvalPayload pattern.EvaluationRequest,
) (pattern.EvaluationResponse, error) {

	defer mutils.TrackTime(h.log, time.Now(), "EvaluateDesign")

	// evaluate specified relationship policies
	// on successful eval the event containing details like comps evaulated, relationships indeitified should be emitted and peristed.
	evaluationResponse, err := h.Rego.RegoPolicyHandler(relationshipPolicyEvalPayload.Design,
		RelationshipPolicyPackageName,
	)

	if err != nil {
		h.log.Debug(err)
		// log an event
		return pattern.EvaluationResponse{}, err
	}

	currentTime := time.Now()
	evaluationResponse.Timestamp = &currentTime

	// Create the event but do not notify the client immediately, as the evaluations are frequent and takes up the view area.
	now := time.Now()
	_ = processEvaluationResponse(h.registryManager, relationshipPolicyEvalPayload, &evaluationResponse)

	evaluatedAliases := ResolveAliasesInDesign(evaluationResponse.Design)
	if evaluationResponse.Design.Metadata == nil {
		evaluationResponse.Design.Metadata = &pattern.PatternFileMetadata{}
	}
	evaluationResponse.Design.Metadata.ResolvedAliases = evaluatedAliases

	mutils.TrackTime(h.log, now, "PostProcessEvaluationResponse")
	return evaluationResponse, nil
}

func processEvaluationResponse(registryManager *registry.RegistryManager, evalPayload pattern.EvaluationRequest, evalResponse *pattern.EvaluationResponse) []*component.ComponentDefinition {

	registryCache := &registry.RegistryEntityCache{}

	compsUpdated := []component.ComponentDefinition{}
	compsAdded := []component.ComponentDefinition{}

	// Bump the version of design
	oldVersion, versionParseErr := semver.NewVersion(evalResponse.Design.Version)
	if versionParseErr != nil {
		oldVersion = semver.MustParse("0.0.0")
	}

	newVersion := oldVersion.IncPatch()
	evalResponse.Design.Version = newVersion.String()

	// components which were added by the evaluator based on the relationship definition, but doesn't exist in the registry.
	unknownComponents := []*component.ComponentDefinition{}

	// Hydrate (replace the partial definition with a complete declaration) the newly added components with the actual
	// component definition from the registry. and add a complete component declaration to the design
	// Refactor To make a single batch call to the registry to get all the components.
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

		// NOTE: this is not deterministic as any version of the component can be used.
		// NOTE: Confirm that the registry returns the latest version in case version is not specified.
		if _c.Model.Model.Version == "*" {
			compFilter.Version = ""
		}

		entities, _, _, _ := registryManager.GetEntitiesMemoized(compFilter, registryCache)
		if len(entities) == 0 {
			unknownComponents = append(unknownComponents, &_c)
			continue
		}
		_component, _ := entities[0].(*component.ComponentDefinition)

		_component.Id = _c.Id
		if _c.DisplayName != "" {
			_component.DisplayName = _c.DisplayName
		} else {
			_component.DisplayName = fmt.Sprintf("%s-%s", strings.ToLower(_component.DisplayName), utils.GetRandomAlphabetsOfDigit(3))
		}

		defaultCapabilities := []capability.Capability{} // only assign empty capabilities for component declarations
		_component.Metadata.IsAnnotation = _c.Metadata.IsAnnotation
		_component.Configuration = _c.Configuration
		_component.Capabilities = &defaultCapabilities
		compsAdded = append(compsAdded, *_component)
	}

	evalResponse.Trace.ComponentsAdded = compsAdded

	for _, component := range evalResponse.Trace.ComponentsUpdated {
		_c := component
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

		evalResponse.Design.Components = append(evalResponse.Design.Components, _c)

	}

	return unknownComponents
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
	evalCtx := r.Context()

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

	evalRespChan := make(chan pattern.EvaluationResponse)
	evalErrChan := make(chan error)

	go func() {
		// Evaluate specified relationship policies
		// Perform the CPU-intensive work
		evaluationResponse, err := h.EvaluateDesign(relationshipPolicyEvalPayload)

		if err != nil {
			evalErrChan <- err // Send the error
		} else {
			evalRespChan <- evaluationResponse // Send the response
		}
	}()

	select {

	case err := <-evalErrChan:
		h.log.Debug(err)
		// log an event
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return

	case evaluationResponse := <-evalRespChan:
		// include trace instead of design file in the event
		event := eventBuilder.WithDescription(fmt.Sprintf("Relationship evaluation completed for design \"%s\" at version \"%s\"", evaluationResponse.Design.Name, evaluationResponse.Design.Version)).
			WithMetadata(map[string]interface{}{
				"trace":        evaluationResponse.Trace,
				"evaluated_at": *evaluationResponse.Timestamp,
			}).WithSeverity(events.Informational).Build()
		_ = provider.PersistEvent(event)

		// write the response
		ec := json.NewEncoder(rw)
		err = ec.Encode(evaluationResponse)
		if err != nil {
			h.log.Error(models.ErrEncoding(err, "policy evaluation response"))
			http.Error(rw, models.ErrEncoding(err, "failed to generate policy evaluation results").Error(), http.StatusInternalServerError)
			return
		}
	case <-evalCtx.Done():
		h.log.Info("Evaluation terminated: request context closed")
		return
	}
}

// Needs to be reinstiated inorder to load the policies based on the evaluation queries.
// This should load policies identified by relationship declarations in the design file.

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
