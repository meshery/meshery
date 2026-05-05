package policies

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// identifyInventoryAdditions emits add_component actions for inventory
// parents the relationship implies but the design lacks. Direct port of
// identify_additions.rego (the rule that auto-creates a Namespace when the
// design contains a Pod referencing one by metadata.namespace but no such
// Namespace component exists).
//
// Why this exists: the OPA/Rego policy engine ships this behavior, the Go
// policy engine did not. Same input design produced different results
// depending on USE_GO_POLICY_ENGINE — a real parity break that surfaced as
// "Kanvas designs are missing components when the Go engine is used."
//
// Fires only on hierarchical/parent/inventory relationships; everything
// else returns no actions, matching the rego rule's guard clauses.
func identifyInventoryAdditions(
	design *pattern.PatternFile,
	relsInScope []*relationship.RelationshipDefinition,
) []PolicyAction {
	if design == nil {
		return nil
	}
	var actions []PolicyAction
	for _, rel := range relsInScope {
		if !isHierarchicalParentInventory(rel) {
			continue
		}
		actions = append(actions, identifyInventoryAdditionsForRel(rel, design)...)
	}
	return actions
}

func isHierarchicalParentInventory(rel *relationship.RelationshipDefinition) bool {
	if rel == nil {
		return false
	}
	return strings.EqualFold(string(rel.Kind), "hierarchical") &&
		strings.EqualFold(rel.RelationshipType, "parent") &&
		strings.EqualFold(rel.SubType, "inventory")
}

// identifyInventoryAdditionsForRel runs the per-relationship branch of
// identify_additions.rego: split the from/to selectors into mutated/mutator
// pools, find existing components matching the mutated selectors, and for
// each (mutated_component, mutator_selector) pair emit a new component
// unless one already exists with matching mutatorRef values.
func identifyInventoryAdditionsForRel(
	rel *relationship.RelationshipDefinition,
	design *pattern.PatternFile,
) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	var out []PolicyAction
	emittedIDs := map[string]bool{}

	for _, ss := range *rel.Selectors {
		mutatedSels, mutatorSels := partitionByPatchRole(ss)

		for _, mutatedSel := range mutatedSels {
			mutatedRefs := *mutatedSel.RelationshipDefinitionSelectorsPatch.MutatedRef
			for _, mutatedComp := range design.Components {
				if mutatedComp == nil {
					continue
				}
				if !selectorAndComponentKindMatches(mutatedSel, mutatedComp) {
					continue
				}
				values, ok := extractValuesAtPaths(mutatedComp, mutatedRefs, design)
				if !ok {
					continue
				}

				for _, mutatorSel := range mutatorSels {
					mutatorRefs := *mutatorSel.RelationshipDefinitionSelectorsPatch.MutatorRef
					if existsMatchingMutatorComponent(design, mutatorSel, mutatorRefs, values) {
						continue
					}
					candidate := buildInventoryParentCandidate(mutatorSel, mutatorRefs, values)
					if candidate == nil {
						continue
					}
					// The rego runs feasibility against the relationship-as-a-whole;
					// re-use the existing helper so cross-model relationships
					// (azure EventSubscription→azure-storage StorageAccount) are
					// treated identically here.
					if feasibleRelationshipSelectorBetween(mutatedComp, candidate, rel) == nil {
						continue
					}
					if emittedIDs[candidate.ID.String()] {
						continue
					}
					emittedIDs[candidate.ID.String()] = true
					out = append(out, newAddComponentAction(candidate))
				}
			}
		}
	}
	return out
}

// partitionByPatchRole walks both allow.from and allow.to and returns
// selectors with mutatedRef separately from those with mutatorRef. Mirrors
// the rego helpers `mutated_selectors` / `mutator_selectors` over the union
// of from+to.
func partitionByPatchRole(ss relationship.SelectorSetItem) (mutated, mutator []relationship.SelectorItem) {
	classify := func(s relationship.SelectorItem) {
		if s.RelationshipDefinitionSelectorsPatch == nil {
			return
		}
		if s.RelationshipDefinitionSelectorsPatch.MutatedRef != nil &&
			len(*s.RelationshipDefinitionSelectorsPatch.MutatedRef) > 0 {
			mutated = append(mutated, s)
		}
		if s.RelationshipDefinitionSelectorsPatch.MutatorRef != nil &&
			len(*s.RelationshipDefinitionSelectorsPatch.MutatorRef) > 0 {
			mutator = append(mutator, s)
		}
	}
	for _, s := range ss.Allow.From {
		classify(s)
	}
	for _, s := range ss.Allow.To {
		classify(s)
	}
	return mutated, mutator
}

// extractValuesAtPaths returns one value per path; a path that resolves to
// nil disqualifies the whole list (matching the rego's
// `count(values) == count(refs)` guard via `value != null` filter).
func extractValuesAtPaths(
	comp *component.ComponentDefinition,
	paths [][]string,
	design *pattern.PatternFile,
) ([]interface{}, bool) {
	if len(paths) == 0 {
		return nil, false
	}
	values := make([]interface{}, 0, len(paths))
	for _, p := range paths {
		v := configurationForComponentAtPath(p, comp, design)
		if v == nil {
			return nil, false
		}
		values = append(values, v)
	}
	return values, true
}

// existsMatchingMutatorComponent answers "is there already a parent in the
// design with values at mutatorRef matching the mutated child's values?" If
// yes, the rego skips emitting the add_component action — duplicates aren't
// wanted.
func existsMatchingMutatorComponent(
	design *pattern.PatternFile,
	mutatorSel relationship.SelectorItem,
	mutatorRefs [][]string,
	mutatedValues []interface{},
) bool {
	for _, comp := range design.Components {
		if comp == nil {
			continue
		}
		if !selectorAndComponentKindMatches(mutatorSel, comp) {
			continue
		}
		existing, ok := extractValuesAtPaths(comp, mutatorRefs, design)
		if !ok {
			continue
		}
		if equalValueLists(mutatedValues, existing) {
			return true
		}
	}
	return false
}

// equalValueLists compares two value lists element-wise via JSON
// equality. The rego match_object uses set-equality but in practice
// mutated/mutator pairs are 1:1 in current relationship JSONs; element-wise
// is sufficient and avoids a quadratic permutation match.
func equalValueLists(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !jsonValuesEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// jsonValuesEqual compares two configuration values using their JSON
// serialization, matching how the rego compares values across types
// (strings, numbers, nested maps).
func jsonValuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aBytes) == string(bBytes)
}

// buildInventoryParentCandidate constructs the new component that would be
// added to the design — the equivalent of the rego's `declaration_with_id`
// at the end of process_comps_to_add.
//
// The candidate's model comes from the mutator selector (the to-side of an
// inventory relationship), which is always concrete in shipped relationship
// JSONs. The candidate's displayName / configuration is populated by
// applying mutator-path patches: typically mutatorRef=[["displayName"]]
// pairs with mutatedValues=["default"] to produce DisplayName="default".
func buildInventoryParentCandidate(
	mutatorSel relationship.SelectorItem,
	mutatorRefs [][]string,
	mutatedValues []interface{},
) *component.ComponentDefinition {
	kind := selectorItemKind(mutatorSel)
	if kind == "" {
		return nil
	}
	id, err := uuid.NewV4()
	if err != nil {
		return nil
	}

	candidate := &component.ComponentDefinition{
		ID:        id,
		Component: component.Component{Kind: kind},
	}

	// Populate both Model and ModelReference. The rego sets only the
	// `model` JSON field (which unmarshals into Go's .Model). Callers
	// downstream of this (registry hydration in
	// policy_relationship_handler.processEvaluationResponse) read .Model
	// first and fall back to .ModelReference; populating both keeps the
	// candidate consistent with what rego emits AND with engine.go's
	// pre-evaluation normalization invariant.
	if mutatorSel.Model != nil {
		candidate.ModelReference = *mutatorSel.Model
		candidate.Model = &modelv1beta1.ModelDefinition{
			Name:        mutatorSel.Model.Name,
			Version:     mutatorSel.Model.Version,
			DisplayName: mutatorSel.Model.DisplayName,
			Model:       mutatorSel.Model.Model,
		}
	}

	// Apply mutator-path patches (e.g. mutatorRef[i]=["displayName"] with
	// mutatedValues[i]="default" → DisplayName="default"). Inventory
	// relationship JSONs in tree today use only top-level paths and
	// configuration.* paths; cover both shapes.
	count := len(mutatorRefs)
	if len(mutatedValues) < count {
		count = len(mutatedValues)
	}
	for i := 0; i < count; i++ {
		applyPathToCandidate(candidate, mutatorRefs[i], mutatedValues[i])
	}

	return candidate
}

// applyPathToCandidate writes value at path on the candidate. Supports the
// two path shapes inventory relationship JSONs actually emit today:
// top-level [displayName] and configuration.<...>. Anything else is
// dropped silently — same fail-mode as the rego when json.patch is given a
// path that doesn't match the target object.
func applyPathToCandidate(comp *component.ComponentDefinition, path []string, value interface{}) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 && path[0] == "displayName" {
		if s, ok := value.(string); ok {
			comp.DisplayName = s
		}
		return
	}
	if path[0] == "configuration" {
		if comp.Configuration == nil {
			comp.Configuration = map[string]interface{}{}
		}
		setNestedConfigValue(comp.Configuration, path[1:], value)
	}
}

func setNestedConfigValue(m map[string]interface{}, path []string, value interface{}) {
	if len(path) == 0 {
		return
	}
	if len(path) == 1 {
		m[path[0]] = value
		return
	}
	next, ok := m[path[0]].(map[string]interface{})
	if !ok {
		next = map[string]interface{}{}
		m[path[0]] = next
	}
	setNestedConfigValue(next, path[1:], value)
}
