package policies

import (
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
	if rel.Selectors == nil || len(*rel.Selectors) == 0 {
		return ""
	}
	ss := (*rel.Selectors)[0]
	if len(ss.Allow.From) == 0 {
		return ""
	}
	if ss.Allow.From[0].ID == nil {
		return ""
	}
	return ss.Allow.From[0].ID.String()
}

// toComponentID returns the to component ID from a relationship definition.
func toComponentID(rel *relationship.RelationshipDefinition) string {
	if rel.Selectors == nil || len(*rel.Selectors) == 0 {
		return ""
	}
	ss := (*rel.Selectors)[0]
	if len(ss.Allow.To) == 0 {
		return ""
	}
	if ss.Allow.To[0].ID == nil {
		return ""
	}
	return ss.Allow.To[0].ID.String()
}

// fromAndToComponentsExist checks if both from and to components still exist in the design.
func fromAndToComponentsExist(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	fromID := fromComponentID(rel)
	toID := toComponentID(rel)
	if fromID == "" || toID == "" {
		return false
	}
	return componentDeclarationByID(design, fromID) != nil &&
		componentDeclarationByID(design, toID) != nil
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

// patchMutatorsAction creates patch actions for relationships that have mutator/mutated refs.
func patchMutatorsAction(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if rel.Selectors == nil {
		return nil
	}
	var actions []PolicyAction

	for _, ss := range *rel.Selectors {
		if len(ss.Allow.From) == 0 || len(ss.Allow.To) == 0 {
			continue
		}
		fromSel := ss.Allow.From[0]
		toSel := ss.Allow.To[0]
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
			if mutatorValue == nil || deepEqual(mutatorValue, oldValue) {
				continue
			}
			actions = append(actions, newComponentUpdateAction(getComponentUpdateOp(mutatedRefs[i]), mutatedID, mutatedRefs[i], mutatorValue))
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
	if len(path) == 0 {
		return nil
	}
	if path[0] == "configuration" {
		config := getComponentConfiguration(comp, design)
		return objectGetNested(config, popFirst(path), nil)
	}
	// For non-configuration paths, convert to map for nested access.
	compMap, err := toGenericMap(comp)
	if err != nil {
		return nil
	}
	return objectGetNested(compMap, path, nil)
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
	parent, refPath := resolveComponentAlias(comp.ID.String(), design)
	if parent != nil {
		parentMap, err := toGenericMap(parent)
		if err != nil {
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

// RelationshipsAreSame checks if two relationships are equivalent.
func RelationshipsAreSame(relA, relB *relationship.RelationshipDefinition) bool {
	if !sameRelationshipIdentifier(relA, relB) {
		return false
	}
	if relA.Selectors == nil || relB.Selectors == nil {
		return false
	}
	for _, ssA := range *relA.Selectors {
		for _, ssB := range *relB.Selectors {
			if len(ssA.Allow.From) == 0 || len(ssB.Allow.From) == 0 ||
				len(ssA.Allow.To) == 0 || len(ssB.Allow.To) == 0 {
				continue
			}
			if sameRelationshipSelectorClause(ssA.Allow.From[0], ssB.Allow.From[0]) &&
				sameRelationshipSelectorClause(ssA.Allow.To[0], ssB.Allow.To[0]) {
				return true
			}
		}
	}
	return false
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
	fromUUID, _ := uuid.FromString(compFromID)
	newFromSel.ID = &fromUUID

	newToSel := toSel
	toUUID, _ := uuid.FromString(compToID)
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
