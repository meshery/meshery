package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	gopolicies "github.com/meshery/meshery/server/policies"
	"github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/capability"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"

	"github.com/meshery/meshkit/models/events"

	"github.com/meshery/meshkit/models/meshmodel/registry"
	regv1alpha3 "github.com/meshery/meshkit/models/meshmodel/registry/v1alpha3"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	patternHelpers "github.com/meshery/meshkit/models/patterns"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/spf13/viper"
)

const (
	//TODO: Might needs to be made dyanamic to support dynamically loading policies based on eval request.
	RelationshipPolicyPackageName = "data.relationship_evaluation_policy"
	suffix                        = "_relationship"
)

const RELATIONSHIP_SUBTYPE_ALIAS = "alias"

// Aliasses Are not resolved
func parseRelationshipToAlias(relationshipDeclaration relationship.RelationshipDefinition) (core.NonResolvedAlias, bool) {

	alias := core.NonResolvedAlias{}

	if relationshipDeclaration.SubType != RELATIONSHIP_SUBTYPE_ALIAS {
		return alias, false
	}

	if relationshipDeclaration.Selectors == nil {
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
	if from.RelationshipDefinitionSelectorsPatch == nil || from.RelationshipDefinitionSelectorsPatch.MutatedRef == nil {
		return alias, false
	}
	mutatedRefs := *from.RelationshipDefinitionSelectorsPatch.MutatedRef

	if len(mutatedRefs) == 0 {
		return alias, false
	}

	if to.ID == nil || from.ID == nil {
		return alias, false
	}

	alias.ImmediateParentId = *to.ID
	alias.AliasComponentId = *from.ID
	alias.RelationshipId = relationshipDeclaration.ID
	alias.ImmediateRefFieldPath = mutatedRefs[0]

	return alias, true

}

func ParseComponentToAlias(component component.ComponentDefinition, relationships []*relationship.RelationshipDefinition) (core.NonResolvedAlias, bool) {

	for _, relationship := range relationships {
		alias, ok := parseRelationshipToAlias(*relationship)
		if !ok {
			continue
		}

		if alias.AliasComponentId == component.ID {
			return alias, true
		}
	}

	return core.NonResolvedAlias{}, false
}

// getComponentById retrieves a component from the design by its ID
func getComponentById(design pattern.PatternFile, id core.Uuid) *component.ComponentDefinition {
	for _, comp := range design.Components {
		if comp.ID == id {
			return comp
		}
	}
	return nil
}

func ResolveAlias(nonResolvedAlias core.NonResolvedAlias, currentNonResolved core.NonResolvedAlias, path []string, design pattern.PatternFile) core.ResolvedAlias {
	parentComponent := getComponentById(design, currentNonResolved.ImmediateParentId)
	if parentComponent == nil {
		return core.ResolvedAliasFromNonResolved(nonResolvedAlias, currentNonResolved.ImmediateParentId, path)
	}

	parentAlias, ok := ParseComponentToAlias(*parentComponent, design.Relationships)

	if !ok {

		return core.ResolvedAliasFromNonResolved(nonResolvedAlias, currentNonResolved.ImmediateParentId, path)

	}

	// slicing from 1 to remove "configuration" prefix when building the resolved ref
	// so we dont get something like configuration,spec,configuration , containers , _
	// appending to aprentAlias.ImmediateReffiled , than path , because this a recursive function it will otherwise build the path in reverse
	return ResolveAlias(nonResolvedAlias, parentAlias, append(parentAlias.ImmediateRefFieldPath, path[1:]...), design)
}

func ResolveAliasesInDesign(design pattern.PatternFile) map[string]core.ResolvedAlias {

	resolvedAliases := make(map[string]core.ResolvedAlias)

	for _, relationship := range design.Relationships {
		nonResolvedalias, ok := parseRelationshipToAlias(*relationship)
		if ok {
			resolvedAlias := ResolveAlias(nonResolvedalias, nonResolvedalias, nonResolvedalias.ImmediateRefFieldPath, design)
			resolvedAliases[resolvedAlias.AliasComponentId.String()] = resolvedAlias
		}
	}

	return resolvedAliases

}

// mergeTraceUnique appends trace entries from src into dst, skipping duplicates by ID.
func mergeTraceUnique(dst, src *pattern.Trace) {
	compSeen := make(map[core.Uuid]bool)
	for _, c := range dst.ComponentsAdded {
		compSeen[c.ID] = true
	}
	for _, c := range dst.ComponentsUpdated {
		compSeen[c.ID] = true
	}
	for _, c := range dst.ComponentsRemoved {
		compSeen[c.ID] = true
	}
	relSeen := make(map[core.Uuid]bool)
	for _, r := range dst.RelationshipsAdded {
		relSeen[r.ID] = true
	}
	for _, r := range dst.RelationshipsUpdated {
		relSeen[r.ID] = true
	}
	for _, r := range dst.RelationshipsRemoved {
		relSeen[r.ID] = true
	}

	for _, c := range src.ComponentsAdded {
		if !compSeen[c.ID] {
			compSeen[c.ID] = true
			dst.ComponentsAdded = append(dst.ComponentsAdded, c)
		}
	}
	for _, c := range src.ComponentsRemoved {
		if !compSeen[c.ID] {
			compSeen[c.ID] = true
			dst.ComponentsRemoved = append(dst.ComponentsRemoved, c)
		}
	}
	for _, c := range src.ComponentsUpdated {
		if !compSeen[c.ID] {
			compSeen[c.ID] = true
			dst.ComponentsUpdated = append(dst.ComponentsUpdated, c)
		}
	}
	for _, r := range src.RelationshipsAdded {
		if !relSeen[r.ID] {
			relSeen[r.ID] = true
			dst.RelationshipsAdded = append(dst.RelationshipsAdded, r)
		}
	}
	for _, r := range src.RelationshipsRemoved {
		if !relSeen[r.ID] {
			relSeen[r.ID] = true
			dst.RelationshipsRemoved = append(dst.RelationshipsRemoved, r)
		}
	}
	for _, r := range src.RelationshipsUpdated {
		if !relSeen[r.ID] {
			relSeen[r.ID] = true
			dst.RelationshipsUpdated = append(dst.RelationshipsUpdated, r)
		}
	}
}

// deduplicateActions removes duplicate actions by (op, id) key.
func deduplicateActions(actions []pattern.Action) []pattern.Action {
	type key struct {
		op string
		id string
	}
	seen := make(map[key]bool)
	var result []pattern.Action
	for _, a := range actions {
		id := ""
		if a.Value != nil {
			if v, ok := a.Value["id"]; ok {
				if s, ok := v.(string); ok {
					id = s
				}
			}
			// For add_component/add_relationship, the id is inside "item"
			if id == "" {
				if item, ok := a.Value["item"].(map[string]interface{}); ok {
					if v, ok := item["id"]; ok {
						if s, ok := v.(string); ok {
							id = s
						}
					}
				}
			}
		}
		k := key{op: a.Op, id: id}
		if !seen[k] {
			seen[k] = true
			result = append(result, a)
		}
	}
	return result
}

func doesntNeedReeval(response pattern.EvaluationResponse) bool {

	for _, action := range response.Actions {
		if action.Op == "delete_component" || action.Op == "add_component" || action.Op == "update_component_configuration" {
			return false
		}
	}

	return true
}

// max number of time to keep revaluating the design till there are no reval triggering actions in the response
const MAX_RE_EVALUATION_DEPTH = 5

// defaultPolicyEvalTimeout bounds a single evaluation. Override via POLICY_EVAL_TIMEOUT.
const defaultPolicyEvalTimeout = 3 * time.Minute

var errEvalTimeout = errors.New("relationship policy evaluation timed out")

func policyEvalTimeout() time.Duration {
	if d := viper.GetDuration("POLICY_EVAL_TIMEOUT"); d > 0 {
		return d
	}
	return defaultPolicyEvalTimeout
}

// Helper method to make design evaluation based on the relationship policies.
// evalIterations is num of passes the evaluator needs to go through to do complete evaluation
func (h *Handler) EvaluateDesign(
	relationshipPolicyEvalPayload pattern.EvaluationRequest,
	evalIterations int,
) (pattern.EvaluationResponse, error) {

	// hydrate the design file components from the registry if needed
	if err := patternHelpers.HydratePattern(&relationshipPolicyEvalPayload.Design, h.registryManager); err != nil {
		h.log.Warnf("failed to hydrate pattern for evaluation: %v", err)
	}

	defer mutils.TrackTime(h.log, time.Now(), "EvaluateDesign")

	var lastEvaluationResponse pattern.EvaluationResponse
	lastEvaluationResponse.Design = relationshipPolicyEvalPayload.Design

	useGoEngine := viper.GetBool("USE_GO_POLICY_ENGINE")

	// Pre-fetch and convert registered relationships once, outside the re-evaluation loop.
	var convertedRels []*relationship.RelationshipDefinition
	if useGoEngine && h.GoEngine != nil {
		registeredRels, _, _, relErr := h.registryManager.GetEntities(&regv1alpha3.RelationshipFilter{})
		if relErr != nil {
			return lastEvaluationResponse, fmt.Errorf("failed to get relationships for Go engine: %w", relErr)
		}
		relInterfaces := make([]interface{}, len(registeredRels))
		for idx, r := range registeredRels {
			relInterfaces[idx] = r
		}
		convertedRels = gopolicies.ConvertRelationships(relInterfaces)
	}

	for i := range MAX_RE_EVALUATION_DEPTH {

		var evaluationResponse pattern.EvaluationResponse
		var err error

		if useGoEngine && h.GoEngine != nil {
			evaluationResponse, err = h.GoEngine.EvaluateDesign(lastEvaluationResponse.Design, convertedRels)
		} else {
			// Use OPA/Rego policy engine (default)
			evaluationResponse, err = h.Rego.RegoPolicyHandler(lastEvaluationResponse.Design,
				RelationshipPolicyPackageName,
			)
		}

		if err != nil {
			h.log.Debug(err)
			// log an event
			return pattern.EvaluationResponse{}, err
		}

		// Save trace and clear it before processEvaluationResponse to prevent
		// it from replacing components with pre-patch versions from trace.
		// ComponentsAdded is kept for hydration (styles/icons).
		savedTrace := evaluationResponse.Trace
		evaluationResponse.Trace = pattern.Trace{
			ComponentsAdded: savedTrace.ComponentsAdded,
		}
		_ = processEvaluationResponse(h.registryManager, relationshipPolicyEvalPayload, &evaluationResponse)
		evaluationResponse.Trace = savedTrace

		evaluatedAliases := ResolveAliasesInDesign(evaluationResponse.Design)
		if evaluationResponse.Design.Metadata == nil {
			evaluationResponse.Design.Metadata = &pattern.PatternFile_Metadata{}
		}
		evaluationResponse.Design.Metadata.ResolvedAliases = &evaluatedAliases

		lastEvaluationResponse.Design = evaluationResponse.Design
		lastEvaluationResponse.Actions = append(lastEvaluationResponse.Actions, evaluationResponse.Actions...)
		mergeTraceUnique(&lastEvaluationResponse.Trace, &evaluationResponse.Trace)

		if evalIterations == i+1 || doesntNeedReeval(evaluationResponse) {
			h.log.Info("Evaluation completed in iteration ", i+1)
			break
		}
		if i == (MAX_RE_EVALUATION_DEPTH - 1) {
			h.log.Warnf("Evaluation depth limit of %d reached; returning partial result", MAX_RE_EVALUATION_DEPTH)
			break
		}

	}

	currentTime := time.Now()
	lastEvaluationResponse.Timestamp = &currentTime

	// Deduplicate actions by (op, id) to avoid duplicates from re-evaluation iterations.
	lastEvaluationResponse.Actions = deduplicateActions(lastEvaluationResponse.Actions)

	// dehydrate the design file components to remove unnecessary details
	patternHelpers.DehydratePattern(&lastEvaluationResponse.Design)

	return lastEvaluationResponse, nil
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
			Limit:      1,
			Trim:       true,
		}

		if _c.Model != nil {
			compFilter.Version = _c.Model.Model.Version
			compFilter.ModelName = _c.Model.Name
			if _c.Model.Model.Version == "*" {
				compFilter.Version = ""
			}
		} else {
			compFilter.ModelName = _c.ModelReference.Name
		}

		entities, _, _, _ := registryManager.GetEntitiesMemoized(compFilter, registryCache)
		if len(entities) == 0 {
			unknownComponents = append(unknownComponents, &_c)
			continue
		}
		_component, _ := entities[0].(*component.ComponentDefinition)

		_component.ID = _c.ID
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
			if c.ID == _c.ID {
				_c = &c
				break
			}
		}

		evalResponse.Design.Components = append(evalResponse.Design.Components, _c)

	}

	return unknownComponents
}

func (h *Handler) EvaluateRelationshipPolicy(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// Wall-clock timeout bounds the HTTP response; the eval goroutine
	// may still run to completion in the background.
	evalCtx, cancel := context.WithTimeout(r.Context(), policyEvalTimeout())
	defer cancel()

	token, _ := evalCtx.Value(models.TokenCtxKey).(string)

	userUUID := user.ID
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
	patternUUID := relationshipPolicyEvalPayload.Design.ID
	eventBuilder.ActedUpon(patternUUID)

	// Coalesce concurrent evaluations of the same design (rage-click guard).
	designKey := patternUUID.String()
	leader, waitCh := h.evalTracker.acquire(designKey)

	if !leader {
		h.log.Debug("coalescing relationship evaluation request for design ", designKey)
		select {
		case result := <-waitCh:
			h.writeEvaluationResult(rw, result)
			return
		case <-evalCtx.Done():
			h.writeEvalCtxError(rw, evalCtx)
			return
		}
	}

	evalRespChan := make(chan pattern.EvaluationResponse, 1)
	evalErrChan := make(chan error, 1)

	go func() {
		// Skip evaluation if the request was already cancelled
		if evalCtx.Err() != nil {
			h.evalTracker.publish(designKey, evalResult{err: evalCtx.Err()})
			return
		}
		evaluationResponse, err := h.EvaluateDesign(relationshipPolicyEvalPayload, MAX_RE_EVALUATION_DEPTH)

		if err != nil {
			h.evalTracker.publish(designKey, evalResult{err: err})
			evalErrChan <- err
		} else {
			h.evalTracker.publish(designKey, evalResult{resp: evaluationResponse})
			evalRespChan <- evaluationResponse
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
		description := fmt.Sprintf("Relationship evaluation complete: %d changes in '%s' at version '%s'", len(evaluationResponse.Actions), evaluationResponse.Design.Name, evaluationResponse.Design.Version)
		event := eventBuilder.WithDescription(description).
			WithMetadata(map[string]interface{}{
				"history_title":       fmt.Sprintf("%d changes made at version %s", len(evaluationResponse.Actions), evaluationResponse.Design.Version),
				"trace":               evaluationResponse.Trace,
				"evaluation_response": evaluationResponse,
				"evaluated_at":        *evaluationResponse.Timestamp,
			}).WithSeverity(events.Informational).Build()
		go func() {
			_ = provider.PersistEvent(*event, token)
		}()

		h.writeEvaluationResult(rw, evalResult{resp: evaluationResponse})
	case <-evalCtx.Done():
		// Unblock any followers waiting on this designID.
		h.evalTracker.publish(designKey, evalResult{err: errEvalTimeout})
		h.writeEvalCtxError(rw, evalCtx)
		return
	}
}

func (h *Handler) writeEvaluationResult(rw http.ResponseWriter, result evalResult) {
	if result.err != nil {
		h.log.Debug(result.err)
		http.Error(rw, result.err.Error(), http.StatusInternalServerError)
		return
	}
	ec := json.NewEncoder(rw)
	if err := ec.Encode(result.resp); err != nil {
		h.log.Error(models.ErrEncoding(err, "policy evaluation response"))
		http.Error(rw, models.ErrEncoding(err, "failed to generate policy evaluation results").Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) writeEvalCtxError(rw http.ResponseWriter, ctx context.Context) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		h.log.Warnf("relationship policy evaluation exceeded %s", policyEvalTimeout())
		http.Error(rw, errEvalTimeout.Error(), http.StatusGatewayTimeout)
		return
	}
	h.log.Info("Evaluation terminated: request context closed")
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

func (h *Handler) GetAllMeshmodelPoliciesByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]
	name := mux.Vars(r)["name"]
	var greedy bool
	if search == queryParamTrue {
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
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetAllMeshmodelPolicies(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]

	var greedy bool
	if search == queryParamTrue {
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
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}
