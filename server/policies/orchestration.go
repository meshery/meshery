package policies

// Shared, registry-INDEPENDENT relationship-evaluation orchestration.
//
// This file is the single source of truth for the parts of relationship
// evaluation that do not depend on the component registry:
//
//   - the MAX_RE_EVALUATION_DEPTH re-evaluation loop
//   - alias resolution (ResolveAliasesInDesign + ResolvedAliasesV1beta2ToV1beta1)
//   - per-iteration trace merge (mergeTraceUnique)
//   - design version bump (semver IncPatch)
//   - cross-iteration action de-duplication (deduplicateActions)
//   - the completed-at Timestamp
//
// Both the HTTP handler (server/handlers/policy_relationship_handler.go,
// which additionally runs the registry-bound HydratePattern /
// processEvaluationResponse / DehydratePattern steps) and the WASM
// entrypoint (cmd/wasm/relationship-engine) call into these shared
// functions so the two evaluators cannot drift. Do NOT copy-paste this
// orchestration anywhere; extend it here.
//
// The registry-bound steps (HydratePattern, processEvaluationResponse,
// reg.GetEntitiesMemoized) deliberately live ONLY in the handler — there
// is no registry in the browser, and the in-browser client compensates
// for component hydration on its side.

import (
	"time"

	"github.com/Masterminds/semver/v3"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	legacycoremodel "github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	coremodelv1beta2 "github.com/meshery/schemas/models/v1beta2/core"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// MAX_RE_EVALUATION_DEPTH is the maximum number of times the design is
// re-evaluated until no re-eval-triggering actions remain in the
// response. This constant is the parity pin: the handler and the WASM
// build MUST share this exact value (asserted by
// TestMaxReEvaluationDepthParity).
const MAX_RE_EVALUATION_DEPTH = 5

const relationshipSubtypeAlias = "alias"

// EvaluationOptions mirrors the registry-independent subset of
// pattern.Options that the orchestration honors.
type EvaluationOptions struct {
	// EnableTrace is accepted for contract symmetry. The trace is always
	// computed (the handler relies on it for hydration); the flag does
	// not gate trace construction in either evaluator today, so it is
	// retained only so callers can pass the parsed options through
	// unchanged.
	EnableTrace bool
	// ReturnDiffOnly, when true, replaces the returned design's
	// Components/Relationships with only the trace deltas — mirroring the
	// registry-independent branch of processEvaluationResponse.
	ReturnDiffOnly bool
}

// EvaluateDesignOrchestrated runs the full registry-INDEPENDENT
// relationship-evaluation pipeline: the MAX_RE_EVALUATION_DEPTH re-eval
// loop around (*GoEngine).EvaluateDesign, per-iteration version bump,
// alias resolution, unique trace merge, cross-iteration action
// de-duplication, optional return-diff-only projection, and the
// completed-at timestamp.
//
// It deliberately does NOT hydrate components from a registry. The HTTP
// handler still performs HydratePattern / processEvaluationResponse /
// DehydratePattern around its own loop; this function is the shared core
// that the WASM build uses verbatim and that the handler's helpers
// (mergeTraceUnique / deduplicateActions / doesntNeedReeval /
// bumpDesignVersion / ResolveAliasesInDesign) are now sourced from, so
// the two evaluators stay in lockstep.
//
// iterations bounds the loop the same way the handler's evalIterations
// argument does; callers pass MAX_RE_EVALUATION_DEPTH for full
// evaluation.
func (e *GoEngine) EvaluateDesignOrchestrated(
	design pattern.PatternFile,
	registeredRelationships []*relationship.RelationshipDefinition,
	iterations int,
	opts EvaluationOptions,
) (pattern.EvaluationResponse, error) {
	var lastEvaluationResponse pattern.EvaluationResponse
	lastEvaluationResponse.Design = design

	if iterations <= 0 || iterations > MAX_RE_EVALUATION_DEPTH {
		iterations = MAX_RE_EVALUATION_DEPTH
	}

	for i := range MAX_RE_EVALUATION_DEPTH {
		evaluationResponse, err := e.EvaluateDesign(lastEvaluationResponse.Design, registeredRelationships)
		if err != nil {
			return pattern.EvaluationResponse{}, err
		}

		// Registry-independent equivalent of the version bump that the
		// handler performs inside processEvaluationResponse. Bumped once
		// per iteration, exactly as the handler does.
		BumpDesignVersion(&evaluationResponse.Design)

		// ReturnDiffOnly is the registry-independent branch of
		// processEvaluationResponse: project the design down to just the
		// trace deltas. Component hydration is intentionally skipped here
		// (no registry in the browser); the in-browser client compensates.
		if opts.ReturnDiffOnly {
			applyReturnDiffOnly(&evaluationResponse)
		}

		evaluatedAliases := ResolveAliasesInDesign(evaluationResponse.Design)
		if evaluationResponse.Design.Metadata == nil {
			evaluationResponse.Design.Metadata = &pattern.PatternFile_Metadata{}
		}
		evaluationResponse.Design.Metadata.ResolvedAliases = patternutils.ResolvedAliasesV1beta2ToV1beta1(&evaluatedAliases)

		lastEvaluationResponse.Design = evaluationResponse.Design
		lastEvaluationResponse.Actions = append(lastEvaluationResponse.Actions, evaluationResponse.Actions...)
		MergeTraceUnique(&lastEvaluationResponse.Trace, &evaluationResponse.Trace)

		if iterations == i+1 || DoesntNeedReeval(evaluationResponse) {
			break
		}
		if i == (MAX_RE_EVALUATION_DEPTH - 1) {
			break
		}
	}

	currentTime := time.Now()
	lastEvaluationResponse.Timestamp = &currentTime

	// Deduplicate actions by (op, id) to avoid duplicates from
	// re-evaluation iterations.
	lastEvaluationResponse.Actions = DeduplicateActions(lastEvaluationResponse.Actions)

	return lastEvaluationResponse, nil
}

// BumpDesignVersion increments the patch component of the design's
// semantic version, falling back to 0.0.1 when the existing version does
// not parse. This is the registry-independent version bump that the
// handler performs inside processEvaluationResponse; it is shared here so
// the handler and the WASM build can never diverge on versioning.
func BumpDesignVersion(design *pattern.PatternFile) {
	oldVersion, versionParseErr := semver.NewVersion(design.Version)
	if versionParseErr != nil {
		oldVersion = semver.MustParse("0.0.0")
	}
	newVersion := oldVersion.IncPatch()
	design.Version = newVersion.String()
}

// applyReturnDiffOnly mirrors the ReturnDiffOnly branch of
// processEvaluationResponse, minus the registry hydration of added
// components (which the browser cannot do and the in-browser client
// reconstructs from the trace). The design is reduced to exactly the
// components and relationships present in the trace deltas.
func applyReturnDiffOnly(resp *pattern.EvaluationResponse) {
	resp.Design.Relationships = []*relationship.RelationshipDefinition{}
	resp.Design.Components = []*component.ComponentDefinition{}

	for _, rel := range resp.Trace.RelationshipsAdded {
		_r := rel
		resp.Design.Relationships = append(resp.Design.Relationships, &_r)
	}
	for _, rel := range resp.Trace.RelationshipsRemoved {
		_r := rel
		resp.Design.Relationships = append(resp.Design.Relationships, &_r)
	}
	for _, rel := range resp.Trace.RelationshipsUpdated {
		_r := rel
		resp.Design.Relationships = append(resp.Design.Relationships, &_r)
	}

	cmps := append(append([]component.ComponentDefinition{}, resp.Trace.ComponentsAdded...), resp.Trace.ComponentsUpdated...)
	for _, cmp := range cmps {
		_c := cmp
		resp.Design.Components = append(resp.Design.Components, &_c)
	}
}

// ----------------------------------------------------------------------
// Registry-independent helpers, relocated from package handlers so that
// both the handler and the WASM build consume one implementation. The
// handler now delegates to these (see policy_relationship_handler.go).
// ----------------------------------------------------------------------

// MergeTraceUnique appends trace entries from src into dst, skipping
// duplicates by ID.
func MergeTraceUnique(dst, src *pattern.Trace) {
	compSeen := make(map[legacycoremodel.Uuid]bool)
	for _, c := range dst.ComponentsAdded {
		compSeen[c.ID] = true
	}
	for _, c := range dst.ComponentsUpdated {
		compSeen[c.ID] = true
	}
	for _, c := range dst.ComponentsRemoved {
		compSeen[c.ID] = true
	}
	relSeen := make(map[legacycoremodel.Uuid]bool)
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

// DeduplicateActions removes duplicate actions by (op, id) key.
func DeduplicateActions(actions []pattern.Action) []pattern.Action {
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

// DoesntNeedReeval reports whether the evaluation response is stable —
// i.e. it produced no component-mutating actions that would require
// another evaluation pass.
func DoesntNeedReeval(response pattern.EvaluationResponse) bool {
	for _, action := range response.Actions {
		if action.Op == "delete_component" || action.Op == "add_component" || action.Op == "update_component_configuration" {
			return false
		}
	}
	return true
}

// ----------------------------------------------------------------------
// Alias resolution, relocated from package handlers. Pure design-graph
// traversal — no registry involved.
// ----------------------------------------------------------------------

// ParseRelationshipToAlias extracts a NonResolvedAlias from an alias
// relationship declaration. Aliases are not yet resolved at this point.
func ParseRelationshipToAlias(relationshipDeclaration relationship.RelationshipDefinition) (coremodelv1beta2.NonResolvedAlias, bool) {
	alias := coremodelv1beta2.NonResolvedAlias{}

	if relationshipDeclaration.SubType != relationshipSubtypeAlias {
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

// ParseComponentToAlias finds the alias (if any) whose alias component is
// the given component.
func ParseComponentToAlias(comp component.ComponentDefinition, relationships []*relationship.RelationshipDefinition) (coremodelv1beta2.NonResolvedAlias, bool) {
	for _, rel := range relationships {
		alias, ok := ParseRelationshipToAlias(*rel)
		if !ok {
			continue
		}
		if alias.AliasComponentId == comp.ID {
			return alias, true
		}
	}
	return coremodelv1beta2.NonResolvedAlias{}, false
}

// getComponentById retrieves a component from the design by its ID.
func getComponentById(design pattern.PatternFile, id legacycoremodel.Uuid) *component.ComponentDefinition {
	for _, comp := range design.Components {
		if comp.ID == id {
			return comp
		}
	}
	return nil
}

// ResolveAlias recursively resolves a single alias up the parent chain to
// produce a fully-qualified ResolvedAlias.
func ResolveAlias(nonResolvedAlias coremodelv1beta2.NonResolvedAlias, currentNonResolved coremodelv1beta2.NonResolvedAlias, path []string, design pattern.PatternFile) coremodelv1beta2.ResolvedAlias {
	parentComponent := getComponentById(design, currentNonResolved.ImmediateParentId)
	if parentComponent == nil {
		return coremodelv1beta2.ResolvedAliasFromNonResolved(nonResolvedAlias, currentNonResolved.ImmediateParentId, path)
	}

	parentAlias, ok := ParseComponentToAlias(*parentComponent, design.Relationships)
	if !ok {
		return coremodelv1beta2.ResolvedAliasFromNonResolved(nonResolvedAlias, currentNonResolved.ImmediateParentId, path)
	}

	// slicing from 1 to remove "configuration" prefix when building the
	// resolved ref so we don't get something like
	// configuration,spec,configuration,containers,_; appending to
	// parentAlias.ImmediateRefFieldPath, then path, because this is a
	// recursive function it will otherwise build the path in reverse.
	return ResolveAlias(nonResolvedAlias, parentAlias, append(parentAlias.ImmediateRefFieldPath, path[1:]...), design)
}

// ResolveAliasesInDesign resolves every alias relationship in the design
// into the keyed ResolvedAlias map stored on the design metadata.
func ResolveAliasesInDesign(design pattern.PatternFile) map[string]coremodelv1beta2.ResolvedAlias {
	resolvedAliases := make(map[string]coremodelv1beta2.ResolvedAlias)

	for _, rel := range design.Relationships {
		nonResolvedalias, ok := ParseRelationshipToAlias(*rel)
		if ok {
			resolvedAlias := ResolveAlias(nonResolvedalias, nonResolvedalias, nonResolvedalias.ImmediateRefFieldPath, design)
			resolvedAliases[resolvedAlias.AliasComponentId.String()] = resolvedAlias
		}
	}

	return resolvedAliases
}
