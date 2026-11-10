package policies

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// componentDeclarationByID finds a component in the design by its ID.
func componentDeclarationByID(design *pattern.PatternFile, id string) *component.ComponentDefinition {
	for _, comp := range design.Components {
		if comp.ID.String() == id {
			return comp
		}
	}
	return nil
}

// fromComponentID returns the from component ID from a relationship definition.
func fromComponentID(rel *relationship.RelationshipDefinition) string {
	if rel == nil || rel.Selectors == nil {
		return ""
	}
	for _, ss := range *rel.Selectors {
		for _, item := range ss.Allow.From {
			if item.ID != nil && item.ID.String() != "" {
				return item.ID.String()
			}
		}
	}
	return ""
}

// toComponentID returns the to component ID from a relationship definition.
func toComponentID(rel *relationship.RelationshipDefinition) string {
	if rel == nil || rel.Selectors == nil {
		return ""
	}
	for _, ss := range *rel.Selectors {
		for _, item := range ss.Allow.To {
			if item.ID != nil && item.ID.String() != "" {
				return item.ID.String()
			}
		}
	}
	return ""
}

// fromAndToComponentsExist checks if both from and to components still exist in the design.
// It inspects all selectors and all referenced component IDs across From and To clauses
// to guarantee that no dangling component references linger in multi-component or
// multi-selector relationships.
func fromAndToComponentsExist(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	if rel == nil || rel.Selectors == nil || len(*rel.Selectors) == 0 || design == nil {
		return false
	}

	for _, ss := range *rel.Selectors {
		if len(ss.Allow.From) == 0 || len(ss.Allow.To) == 0 {
			return false
		}
		for _, fromSel := range ss.Allow.From {
			if fromSel.ID == nil || fromSel.ID.String() == "" {
				return false
			}
			if componentDeclarationByID(design, fromSel.ID.String()) == nil {
				return false
			}
		}
		for _, toSel := range ss.Allow.To {
			if toSel.ID == nil || toSel.ID.String() == "" {
				return false
			}
			if componentDeclarationByID(design, toSel.ID.String()) == nil {
				return false
			}
		}
	}

	return true
}

// fromOrToComponentsDontExist returns true when one or both sides are missing.
func fromOrToComponentsDontExist(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return !fromAndToComponentsExist(rel, design)
}

// approveRelationshipsAction creates update actions to approve relationships with given statuses.
func approveRelationshipsAction(rels []*relationship.RelationshipDefinition, statuses map[string]bool, maxLimit int) []PolicyAction {
	var actions []PolicyAction
	for _, rel := range rels {
		if len(actions) >= maxLimit {
			break
		}
		if !statuses[getRelStatus(rel)] {
			continue
		}
		actions = append(actions, newUpdateRelationshipAction(rel.ID.String(), "/status", StatusApproved))
	}
	return actions
}

// approveIdentifiedRelationshipsAction approves all identified relationships (up to limit).
func approveIdentifiedRelationshipsAction(rels []*relationship.RelationshipDefinition, maxLimit int) []PolicyAction {
	return approveRelationshipsAction(rels, map[string]bool{StatusIdentified: true}, maxLimit)
}

// cleanupDeletedRelationshipsActions creates delete actions for relationships with status "deleted".
func cleanupDeletedRelationshipsActions(rels []*relationship.RelationshipDefinition) []PolicyAction {
	var actions []PolicyAction
	for _, rel := range rels {
		if getRelStatus(rel) == StatusDeleted {
			actions = append(actions, newDeleteRelationshipAction(rel.ID.String(), rel))
		}
	}
	return actions
}

// schemaDefaultAtPath walks a component's JSON schema to return the `default`
// declared for the field at the given configuration path. Returns (nil, false)
// when the schema is missing, unparseable, or defines no default at that path.
func schemaDefaultAtPath(comp *component.ComponentDefinition, path []string) (interface{}, bool) {
	if comp == nil || comp.Component.Schema == "" {
		return nil, false
	}
	var schema map[string]interface{}
	if err := json.Unmarshal([]byte(comp.Component.Schema), &schema); err != nil {
		return nil, false
	}
	segs := path
	if len(segs) > 0 && segs[0] == "configuration" {
		segs = segs[1:]
	}
	cursor := schema
	for _, seg := range segs {
		props, ok := cursor["properties"].(map[string]interface{})
		if !ok {
			return nil, false
		}
		next, ok := props[seg].(map[string]interface{})
		if !ok {
			return nil, false
		}
		cursor = next
	}
	def, ok := cursor["default"]
	return def, ok
}

// cleanupActionForPath returns the reverse patch for a mutated field on delete:
// the schema default when one is declared, otherwise a remove of the field.
func cleanupActionForPath(mutatedID string, mutatedPath []string, mutatedComp *component.ComponentDefinition) PolicyAction {
	if def, ok := schemaDefaultAtPath(mutatedComp, mutatedPath); ok {
		return newComponentUpdateAction(getComponentUpdateOp(mutatedPath), mutatedID, mutatedPath, def)
	}
	return newComponentRemoveAction(mutatedID, mutatedPath)
}

// patchMutatorsAction creates patch actions for relationships that have mutator/mutated refs.
// When the relationship is in StatusDeleted, it emits reverse patches (schema default when
// declared, otherwise remove) for fields that still hold the mutator's value, so deletion
// restores the pre-mutation state.
func patchMutatorsAction(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	reverse := getRelStatus(rel) == StatusDeleted
	var actions []PolicyAction

	for _, ss := range *rel.Selectors {
		fromCount := len(ss.Allow.From)
		toCount := len(ss.Allow.To)
		if fromCount == 0 || toCount == 0 {
			continue
		}
		pairCount := fromCount
		if toCount < pairCount {
			pairCount = toCount
		}
		// If one side has 1 element and the other side has N elements, fan out across all N targets.
		if fromCount == 1 || toCount == 1 {
			pairCount = max(fromCount, toCount)
		}

		for p := 0; p < pairCount; p++ {
			fromIdx := p
			if fromIdx >= fromCount {
				fromIdx = 0
			}
			toIdx := p
			if toIdx >= toCount {
				toIdx = 0
			}
			fromSel := ss.Allow.From[fromIdx]
			toSel := ss.Allow.To[toIdx]
			if fromSel.RelationshipDefinitionSelectorsPatch == nil || toSel.RelationshipDefinitionSelectorsPatch == nil {
				continue
			}

			fromID := selectorItemID(fromSel)
			toID := selectorItemID(toSel)
			fromComp := componentDeclarationByID(design, fromID)
			toComp := componentDeclarationByID(design, toID)
			if fromComp == nil || toComp == nil {
				continue
			}

			mutatorRefs, mutatedRefs, mutatorComp, mutatedComp := resolveMutatorMutatedRefs(fromSel.RelationshipDefinitionSelectorsPatch, toSel.RelationshipDefinitionSelectorsPatch, fromComp, toComp)
			if mutatorRefs == nil || mutatedRefs == nil {
				continue
			}

			mutatedID := toID
			if mutatedComp.ID == fromComp.ID {
				mutatedID = fromID
			}

			count := len(mutatorRefs)
			if len(mutatedRefs) < count {
				count = len(mutatedRefs)
			}

			for i := 0; i < count; i++ {
				mutatorValue := configurationForComponentAtPath(mutatorRefs[i], mutatorComp, design)
				oldValue := configurationForComponentAtPath(mutatedRefs[i], mutatedComp, design)
				if mutatorValue == nil {
					continue
				}
				if reverse {
					if !deepEqual(mutatorValue, oldValue) {
						continue
					}
					actions = append(actions, cleanupActionForPath(mutatedID, mutatedRefs[i], mutatedComp))
					continue
				}
				if deepEqual(mutatorValue, oldValue) {
					continue
				}
				actions = append(actions, newComponentUpdateAction(getComponentUpdateOp(mutatedRefs[i]), mutatedID, mutatedRefs[i], mutatorValue))
			}
		}
	}
	return actions
}

// selectorItemID returns the ID string from a SelectorItem, or empty if nil.
func selectorItemID(sel relationship.SelectorItem) string {
	if sel.ID == nil {
		return ""
	}
	return sel.ID.String()
}

// selectorItemKind returns the kind string from a SelectorItem, or empty if nil.
func selectorItemKind(sel relationship.SelectorItem) string {
	if sel.Kind == nil {
		return ""
	}
	return *sel.Kind
}

// configurationForComponentAtPath gets the value at a path in a component,
// handling the "configuration" prefix and alias resolution.
func configurationForComponentAtPath(path []string, comp *component.ComponentDefinition, design *pattern.PatternFile) interface{} {
	return configurationForComponentAtPathWithCache(path, comp, design, nil)
}

// configurationForComponentAtPathWithCache is configurationForComponentAtPath
// with an optional pass-level component map cache.
func configurationForComponentAtPathWithCache(path []string, comp *component.ComponentDefinition, design *pattern.PatternFile, cache componentMapCache) interface{} {
	if len(path) == 0 {
		return nil
	}
	if path[0] == "configuration" {
		config := getComponentConfigurationWithCache(comp, design, cache)
		return objectGetNested(config, popFirst(path), nil)
	}
	return objectGetNested(cache.genericMap(comp), path, nil)
}

// resolveComponentAlias checks if a component is an alias and returns the
// resolved parent component and reference path.
func resolveComponentAlias(compID string, design *pattern.PatternFile) (*component.ComponentDefinition, []string) {
	if design.Metadata == nil || design.Metadata.ResolvedAliases == nil {
		return nil, nil
	}
	alias, exists := (*design.Metadata.ResolvedAliases)[compID]
	if !exists {
		return nil, nil
	}
	parent := componentDeclarationByID(design, alias.ResolvedParentId.String())
	if parent == nil {
		return nil, nil
	}
	return parent, alias.ResolvedRefFieldPath
}

// getComponentConfiguration returns the configuration for a component,
// resolving aliases if present.
func getComponentConfiguration(comp *component.ComponentDefinition, design *pattern.PatternFile) interface{} {
	return getComponentConfigurationWithCache(comp, design, nil)
}

// getComponentConfigurationWithCache is getComponentConfiguration with an
// optional pass-level component map cache for resolved alias parents.
func getComponentConfigurationWithCache(comp *component.ComponentDefinition, design *pattern.PatternFile, cache componentMapCache) interface{} {
	parent, refPath := resolveComponentAlias(comp.ID.String(), design)
	if parent != nil {
		parentMap := cache.genericMap(parent)
		if parentMap == nil {
			return nil
		}
		return objectGetNested(parentMap, refPath, nil)
	}
	return comp.Configuration
}

// sameRelationshipIdentifier checks if two relationships have the same kind/type/subType.
func sameRelationshipIdentifier(relA, relB *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(relA.Kind), string(relB.Kind)) &&
		strings.EqualFold(relA.RelationshipType, relB.RelationshipType) &&
		strings.EqualFold(relA.SubType, relB.SubType)
}

// sameRelationshipSelectorClause checks if two selector items reference the same component.
func sameRelationshipSelectorClause(a, b relationship.SelectorItem) bool {
	return selectorItemKind(a) == selectorItemKind(b) &&
		selectorItemID(a) == selectorItemID(b) &&
		deepEqual(a.RelationshipDefinitionSelectorsPatch, b.RelationshipDefinitionSelectorsPatch)
}

// sameSelectorSlice checks if two slices of SelectorItem match in multiset comparison.
func sameSelectorSlice(sliceA, sliceB []relationship.SelectorItem) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}
	used := make([]bool, len(sliceB))
	for _, itemA := range sliceA {
		matched := false
		for j, itemB := range sliceB {
			if !used[j] && sameRelationshipSelectorClause(itemA, itemB) {
				used[j] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// sameSelectorSetItem checks if two SelectorSetItem values match across both Allow and Deny clauses.
func sameSelectorSetItem(ssA, ssB relationship.SelectorSetItem) bool {
	if !sameSelectorSlice(ssA.Allow.From, ssB.Allow.From) {
		return false
	}
	if !sameSelectorSlice(ssA.Allow.To, ssB.Allow.To) {
		return false
	}
	if (ssA.Deny == nil) != (ssB.Deny == nil) {
		return false
	}
	if ssA.Deny != nil {
		if !sameSelectorSlice(ssA.Deny.From, ssB.Deny.From) {
			return false
		}
		if !sameSelectorSlice(ssA.Deny.To, ssB.Deny.To) {
			return false
		}
	}
	return true
}

// RelationshipsAreSame checks if two relationships are equivalent across their full selector sets.
func RelationshipsAreSame(relA, relB *relationship.RelationshipDefinition) bool {
	if !sameRelationshipIdentifier(relA, relB) {
		return false
	}
	if relA.Selectors == nil || relB.Selectors == nil {
		return false
	}
	if len(*relA.Selectors) != len(*relB.Selectors) {
		return false
	}

	used := make([]bool, len(*relB.Selectors))
	for _, ssA := range *relA.Selectors {
		matched := false
		for j, ssB := range *relB.Selectors {
			if !used[j] && sameSelectorSetItem(ssA, ssB) {
				used[j] = true
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// relationshipAlreadyExists checks if a relationship already exists in the design.
func relationshipAlreadyExists(design *pattern.PatternFile, rel *relationship.RelationshipDefinition) bool {
	for _, existing := range design.Relationships {
		if getRelStatus(existing) == StatusDeleted {
			continue
		}
		if RelationshipsAreSame(existing, rel) {
			return true
		}
	}
	return false
}

// matchValues checks if two values match according to a strategy.
func matchValues(fromValue, toValue interface{}, strategy string) bool {
	switch strategy {
	case "equal":
		return deepEqual(fromValue, toValue)
	case "equal_as_strings":
		return fmt.Sprintf("%v", fromValue) == fmt.Sprintf("%v", toValue)
	case "to_contains_from":
		return objectSubset(toValue, fromValue)
	case "not_null":
		return fromValue != nil && toValue != nil
	default:
		return deepEqual(fromValue, toValue)
	}
}

// matchValuesWithStrategies checks if values match all strategies.
func matchValuesWithStrategies(fromValue, toValue interface{}, strategies []string) bool {
	for _, strategy := range strategies {
		if !matchValues(fromValue, toValue, strategy) {
			return false
		}
	}
	return true
}

// objectSubset checks if fromValue is a subset of toValue (for maps and slices).
func objectSubset(toValue, fromValue interface{}) bool {
	if fromValue == nil {
		return true
	}
	fromMap, fromOk := fromValue.(map[string]interface{})
	toMap, toOk := toValue.(map[string]interface{})
	if fromOk && toOk {
		for k, fv := range fromMap {
			tv, exists := toMap[k]
			if !exists || !objectSubset(tv, fv) {
				return false
			}
		}
		return true
	}
	return deepEqual(fromValue, toValue)
}

// componentMatchesKind checks if a component matches the given kind selector.
func componentMatchesKind(comp *component.ComponentDefinition, kind string) bool {
	return matchName(comp.Component.Kind, kind)
}

// buildIdentifiedRelationship builds an identified relationship declaration
// from matched component IDs, selectors, and the relationship definition.
func buildIdentifiedRelationship(
	fromSel, toSel relationship.SelectorItem,
	compFromID, compToID string,
	relDef *relationship.RelationshipDefinition,
) *relationship.RelationshipDefinition {
	newFromSel := fromSel
	fromUUID := uuid.FromStringOrNil(compFromID)
	newFromSel.ID = &fromUUID

	newToSel := toSel
	toUUID := uuid.FromStringOrNil(compToID)
	newToSel.ID = &toUUID

	selectorSet := relationship.SelectorSetItem{
		Allow: relationship.Selector{
			From: []relationship.SelectorItem{newFromSel},
			To:   []relationship.SelectorItem{newToSel},
		},
	}

	decl := deepCopyRelDef(relDef)
	selectors := relationship.SelectorSet{selectorSet}
	decl.Selectors = &selectors
	decl.ID = staticUUID(selectorSet)
	setRelStatus(decl, StatusIdentified)
	return decl
}

// matchComponentPairs finds all matching component pairs for a given from/to selector pair.
func matchComponentPairs(
	fromSel, toSel relationship.SelectorItem,
	comps []*component.ComponentDefinition,
	relDef *relationship.RelationshipDefinition,
	design *pattern.PatternFile,
) []*relationship.RelationshipDefinition {
	var results []*relationship.RelationshipDefinition
	fromKind := selectorItemKind(fromSel)
	toKind := selectorItemKind(toSel)

	for _, compFrom := range comps {
		if !componentMatchesKind(compFrom, fromKind) {
			continue
		}
		for _, compTo := range comps {
			if compFrom.ID == compTo.ID {
				continue
			}
			if !componentMatchesKind(compTo, toKind) {
				continue
			}
			if !matchingMutators(compFrom, compTo, fromSel, toSel, design) {
				continue
			}
			results = append(results, buildIdentifiedRelationship(fromSel, toSel, compFrom.ID.String(), compTo.ID.String(), relDef))
		}
	}
	return results
}

// identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields identifies relationships
// where components have matching values at mutator/mutated paths.
func identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(
	relDef *relationship.RelationshipDefinition,
	design *pattern.PatternFile,
) []*relationship.RelationshipDefinition {
	if relDef.Selectors == nil {
		return nil
	}
	var identified []*relationship.RelationshipDefinition

	for _, ss := range *relDef.Selectors {
		for _, fromSel := range ss.Allow.From {
			for _, toSel := range ss.Allow.To {
				matches := matchComponentPairs(fromSel, toSel, design.Components, relDef, design)
				identified = append(identified, matches...)
			}
		}
	}
	return identified
}

// resolveMutatorMutatedRefs extracts mutatorRef/mutatedRef paths and the corresponding
// components from the from/to patches.
func resolveMutatorMutatedRefs(
	fromPatch, toPatch *relationship.RelationshipDefinitionSelectorsPatch,
	compFrom, compTo *component.ComponentDefinition,
) (mutatorRefs, mutatedRefs [][]string, mutatorComp, mutatedComp *component.ComponentDefinition) {
	if fromPatch.MutatorRef != nil && toPatch.MutatedRef != nil {
		return *fromPatch.MutatorRef, *toPatch.MutatedRef, compFrom, compTo
	}
	if toPatch.MutatorRef != nil && fromPatch.MutatedRef != nil {
		return *toPatch.MutatorRef, *fromPatch.MutatedRef, compTo, compFrom
	}
	return nil, nil, nil, nil
}

// matchingMutators checks if the mutator/mutated refs of two components match.
func matchingMutators(
	compFrom, compTo *component.ComponentDefinition,
	fromClause, toClause relationship.SelectorItem,
	design *pattern.PatternFile,
) bool {
	if fromClause.RelationshipDefinitionSelectorsPatch == nil || toClause.RelationshipDefinitionSelectorsPatch == nil {
		return false
	}

	mutatorRefs, mutatedRefs, mutatorComp, mutatedComp := resolveMutatorMutatedRefs(fromClause.RelationshipDefinitionSelectorsPatch, toClause.RelationshipDefinitionSelectorsPatch, compFrom, compTo)
	if mutatorRefs == nil || mutatedRefs == nil {
		return false
	}

	count := len(mutatorRefs)
	if len(mutatedRefs) < count {
		count = len(mutatedRefs)
	}
	if count == 0 {
		return false
	}

	strategies := getMatchStrategyForSelector(fromClause)

	for i := 0; i < count; i++ {
		mutatorValue := configurationForComponentAtPath(mutatorRefs[i], mutatorComp, design)
		mutatedValue := configurationForComponentAtPath(mutatedRefs[i], mutatedComp, design)
		strategy := getStrategyForValueAt(strategies, i)
		if !matchValuesWithStrategies(mutatorValue, mutatedValue, strategy) {
			return false
		}
	}
	return true
}

// getMatchStrategyForSelector returns the match strategy matrix from a selector item.
func getMatchStrategyForSelector(sel relationship.SelectorItem) [][]string {
	if sel.MatchStrategyMatrix == nil {
		return nil
	}
	return *sel.MatchStrategyMatrix
}

// getStrategyForValueAt returns the match strategy for a specific index.
func getStrategyForValueAt(strategies [][]string, index int) []string {
	if index < len(strategies) && strategies[index] != nil {
		return strategies[index]
	}
	return []string{"equal"}
}

// findRelationshipByID finds a relationship in the design by its ID.
func findRelationshipByID(design *pattern.PatternFile, id string) *relationship.RelationshipDefinition {
	for _, rel := range design.Relationships {
		if rel.ID.String() == id {
			return rel
		}
	}
	return nil
}
