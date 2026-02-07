package policies

import (
	"fmt"
	"strings"
)

// componentDeclarationByID finds a component in the design by its ID.
func componentDeclarationByID(design map[string]interface{}, id string) map[string]interface{} {
	comps := getMapSlice(design, "components")
	for _, c := range comps {
		comp, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		compID := getMapString(comp, "id")
		if compID == id {
			return comp
		}
	}
	return nil
}

// fromComponentID returns the from component ID from a relationship declaration.
func fromComponentID(rel map[string]interface{}) string {
	selectors := getMapSlice(rel, "selectors")
	if len(selectors) == 0 {
		return ""
	}
	sel, ok := selectors[0].(map[string]interface{})
	if !ok {
		return ""
	}
	allow := getMapMap(sel, "allow")
	if allow == nil {
		return ""
	}
	from := getMapSlice(allow, "from")
	if len(from) == 0 {
		return ""
	}
	fromSel, ok := from[0].(map[string]interface{})
	if !ok {
		return ""
	}
	return getMapString(fromSel, "id")
}

// toComponentID returns the to component ID from a relationship declaration.
func toComponentID(rel map[string]interface{}) string {
	selectors := getMapSlice(rel, "selectors")
	if len(selectors) == 0 {
		return ""
	}
	sel, ok := selectors[0].(map[string]interface{})
	if !ok {
		return ""
	}
	allow := getMapMap(sel, "allow")
	if allow == nil {
		return ""
	}
	to := getMapSlice(allow, "to")
	if len(to) == 0 {
		return ""
	}
	toSel, ok := to[0].(map[string]interface{})
	if !ok {
		return ""
	}
	return getMapString(toSel, "id")
}

// fromAndToComponentsExist checks if both from and to components still exist in the design.
func fromAndToComponentsExist(rel, design map[string]interface{}) bool {
	fromID := fromComponentID(rel)
	toID := toComponentID(rel)

	if fromID == "" || toID == "" {
		return false
	}

	fromComp := componentDeclarationByID(design, fromID)
	toComp := componentDeclarationByID(design, toID)

	return fromComp != nil && toComp != nil
}

// fromOrToComponentsDontExist returns true when one or both sides are missing.
func fromOrToComponentsDontExist(rel, design map[string]interface{}) bool {
	return !fromAndToComponentsExist(rel, design)
}

// approveRelationshipsAction creates update actions to approve relationships with given statuses.
func approveRelationshipsAction(relationships []map[string]interface{}, statuses map[string]bool, maxLimit int) []PolicyAction {
	var actions []PolicyAction
	for _, rel := range relationships {
		if len(actions) >= maxLimit {
			break
		}
		status := getMapString(rel, "status")
		if !statuses[status] {
			continue
		}
		actions = append(actions, PolicyAction{
			Op: UpdateRelationshipOp,
			Value: map[string]interface{}{
				"id":    getMapString(rel, "id"),
				"path":  "/status",
				"value": "approved",
			},
		})
	}
	return actions
}

// approveIdentifiedRelationshipsAction approves all identified relationships (up to limit).
func approveIdentifiedRelationshipsAction(relationships []map[string]interface{}, maxLimit int) []PolicyAction {
	return approveRelationshipsAction(relationships, map[string]bool{"identified": true}, maxLimit)
}

// cleanupDeletedRelationshipsActions creates delete actions for relationships with status "deleted".
func cleanupDeletedRelationshipsActions(relationships []map[string]interface{}) []PolicyAction {
	var actions []PolicyAction
	for _, rel := range relationships {
		if getMapString(rel, "status") == "deleted" {
			actions = append(actions, PolicyAction{
				Op: DeleteRelationshipOp,
				Value: map[string]interface{}{
					"id":           getMapString(rel, "id"),
					"relationship": rel,
				},
			})
		}
	}
	return actions
}

// patchMutatorsAction creates patch actions for relationships that have mutator/mutated refs.
func patchMutatorsAction(rel, designFile map[string]interface{}) []PolicyAction {
	var actions []PolicyAction
	selectors := getMapSlice(rel, "selectors")

	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}

		fromArr := getMapSlice(allow, "from")
		toArr := getMapSlice(allow, "to")
		if len(fromArr) == 0 || len(toArr) == 0 {
			continue
		}

		from, _ := fromArr[0].(map[string]interface{})
		to, _ := toArr[0].(map[string]interface{})
		if from == nil || to == nil {
			continue
		}

		fromPatch := getMapMap(from, "patch")
		toPatch := getMapMap(to, "patch")
		if fromPatch == nil || toPatch == nil {
			continue
		}

		mutatorRefRaw := getMapSlice(fromPatch, "mutatorRef")
		mutatedRefRaw := getMapSlice(toPatch, "mutatedRef")
		if mutatorRefRaw == nil || mutatedRefRaw == nil {
			continue
		}

		fromComp := componentDeclarationByID(designFile, getMapString(from, "id"))
		toComp := componentDeclarationByID(designFile, getMapString(to, "id"))
		if fromComp == nil || toComp == nil {
			continue
		}

		count := len(mutatorRefRaw)
		if len(mutatedRefRaw) < count {
			count = len(mutatedRefRaw)
		}

		for i := 0; i < count; i++ {
			mutatorRef := interfaceToStringSlice(mutatorRefRaw[i])
			mutatedRef := interfaceToStringSlice(mutatedRefRaw[i])

			mutatorValue := configurationForComponentAtPath(mutatorRef, fromComp, designFile)
			oldValue := configurationForComponentAtPath(mutatedRef, toComp, designFile)

			if deepEqual(mutatorValue, oldValue) {
				continue
			}

			actions = append(actions, PolicyAction{
				Op: getComponentUpdateOp(mutatedRef),
				Value: map[string]interface{}{
					"id":    getMapString(to, "id"),
					"path":  mutatedRef,
					"value": mutatorValue,
				},
			})
		}
	}

	return actions
}

// configurationForComponentAtPath gets the value at a path in a component,
// handling the "configuration" prefix and alias resolution.
func configurationForComponentAtPath(path []string, component, design map[string]interface{}) interface{} {
	if len(path) == 0 {
		return nil
	}

	if path[0] == "configuration" {
		config := getComponentConfiguration(component, design)
		return objectGetNested(config, popFirst(path), nil)
	}

	return objectGetNested(component, path, nil)
}

// getComponentConfiguration returns the configuration for a component,
// resolving aliases if present.
func getComponentConfiguration(component, design map[string]interface{}) interface{} {
	compID := getMapString(component, "id")

	// Check for alias
	metadata := getMapMap(design, "metadata")
	if metadata != nil {
		aliases := getMapMap(metadata, "resolvedAliases")
		if aliases != nil {
			alias := getMapMap(aliases, compID)
			if alias != nil {
				parentID := getMapString(alias, "resolved_parent_id")
				parent := componentDeclarationByID(design, parentID)
				if parent != nil {
					refPath := interfaceToStringSlice(alias["resolved_ref_field_path"])
					return objectGetNested(parent, refPath, nil)
				}
			}
		}
	}

	return component["configuration"]
}

// sameRelationshipIdentifier checks if two relationships have the same kind/type/subType.
func sameRelationshipIdentifier(relA, relB map[string]interface{}) bool {
	return strings.EqualFold(getMapString(relA, "kind"), getMapString(relB, "kind")) &&
		strings.EqualFold(getMapString(relA, "type"), getMapString(relB, "type")) &&
		strings.EqualFold(getMapString(relA, "subType"), getMapString(relB, "subType"))
}

// sameRelationshipSelectorClause checks if two selector clauses reference the same component.
func sameRelationshipSelectorClause(clauseA, clauseB map[string]interface{}) bool {
	return getMapString(clauseA, "kind") == getMapString(clauseB, "kind") &&
		getMapString(clauseA, "id") == getMapString(clauseB, "id") &&
		deepEqual(clauseA["patch"], clauseB["patch"])
}

// relationshipsAreSame checks if two relationships are equivalent.
func relationshipsAreSame(relA, relB map[string]interface{}) bool {
	if !sameRelationshipIdentifier(relA, relB) {
		return false
	}

	selectorsA := getMapSlice(relA, "selectors")
	selectorsB := getMapSlice(relB, "selectors")

	for _, sa := range selectorsA {
		selA, ok := sa.(map[string]interface{})
		if !ok {
			continue
		}
		for _, sb := range selectorsB {
			selB, ok := sb.(map[string]interface{})
			if !ok {
				continue
			}
			allowA := getMapMap(selA, "allow")
			allowB := getMapMap(selB, "allow")
			if allowA == nil || allowB == nil {
				continue
			}

			fromA := getMapSlice(allowA, "from")
			fromB := getMapSlice(allowB, "from")
			toA := getMapSlice(allowA, "to")
			toB := getMapSlice(allowB, "to")

			if len(fromA) == 0 || len(fromB) == 0 || len(toA) == 0 || len(toB) == 0 {
				continue
			}

			fa, _ := fromA[0].(map[string]interface{})
			fb, _ := fromB[0].(map[string]interface{})
			ta, _ := toA[0].(map[string]interface{})
			tb, _ := toB[0].(map[string]interface{})

			if fa != nil && fb != nil && ta != nil && tb != nil {
				if sameRelationshipSelectorClause(fa, fb) && sameRelationshipSelectorClause(ta, tb) {
					return true
				}
			}
		}
	}
	return false
}

// relationshipAlreadyExists checks if a relationship already exists in the design.
func relationshipAlreadyExists(design, relationship map[string]interface{}) bool {
	rels := getMapSlice(design, "relationships")
	for _, r := range rels {
		existing, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if getMapString(existing, "status") == "deleted" {
			continue
		}
		if relationshipsAreSame(existing, relationship) {
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
		s1 := fmt.Sprintf("%v", fromValue)
		s2 := fmt.Sprintf("%v", toValue)
		return s1 == s2
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
			if !exists {
				return false
			}
			if !objectSubset(tv, fv) {
				return false
			}
		}
		return true
	}

	return deepEqual(fromValue, toValue)
}

// identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields identifies relationships
// where components have matching values at mutator/mutated paths.
func identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(
	relDef map[string]interface{},
	designFile map[string]interface{},
) []map[string]interface{} {
	var identified []map[string]interface{}

	selectors := getMapSlice(relDef, "selectors")
	comps := getMapSlice(designFile, "components")

	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}

		fromSelectors := getMapSlice(allow, "from")
		toSelectors := getMapSlice(allow, "to")

		for _, fs := range fromSelectors {
			fromSel, ok := fs.(map[string]interface{})
			if !ok {
				continue
			}
			for _, ts := range toSelectors {
				toSel, ok := ts.(map[string]interface{})
				if !ok {
					continue
				}

				fromKind := getMapString(fromSel, "kind")
				toKind := getMapString(toSel, "kind")

				for _, cf := range comps {
					compFrom, ok := cf.(map[string]interface{})
					if !ok {
						continue
					}
					compFromComponent := getMapMap(compFrom, "component")
					if getMapString(compFromComponent, "kind") != fromKind && fromKind != "*" {
						continue
					}

					for _, ct := range comps {
						compTo, ok := ct.(map[string]interface{})
						if !ok {
							continue
						}
						if getMapString(compFrom, "id") == getMapString(compTo, "id") {
							continue
						}

						compToComponent := getMapMap(compTo, "component")
						if getMapString(compToComponent, "kind") != toKind && toKind != "*" {
							continue
						}

						if !matchingMutators(compFrom, compTo, fromSel, toSel, designFile) {
							continue
						}

						// Build the identified relationship declaration
						newFromSel := deepCopyMap(fromSel)
						newFromSel["id"] = getMapString(compFrom, "id")
						newToSel := deepCopyMap(toSel)
						newToSel["id"] = getMapString(compTo, "id")

						selectorDecl := map[string]interface{}{
							"allow": map[string]interface{}{
								"from": []interface{}{newFromSel},
								"to":   []interface{}{newToSel},
							},
							"deny": map[string]interface{}{},
						}

						decl := deepCopyMap(relDef)
						decl["selectors"] = []interface{}{selectorDecl}
						decl["id"] = newUUID(selectorDecl).String()
						decl["status"] = "identified"

						identified = append(identified, decl)
					}
				}
			}
		}
	}

	return identified
}

// matchingMutators checks if the mutator/mutated refs of two components match.
func matchingMutators(compFrom, compTo, fromClause, toClause, designFile map[string]interface{}) bool {
	fromPatch := getMapMap(fromClause, "patch")
	toPatch := getMapMap(toClause, "patch")
	if fromPatch == nil || toPatch == nil {
		return false
	}

	mutatorRefs := getMapSlice(fromPatch, "mutatorRef")
	mutatedRefs := getMapSlice(toPatch, "mutatedRef")
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

	// Get match strategy matrix
	strategies := getMatchStrategyForSelector(fromClause)

	for i := 0; i < count; i++ {
		mutatorPath := interfaceToStringSlice(mutatorRefs[i])
		mutatedPath := interfaceToStringSlice(mutatedRefs[i])

		mutatorValue := configurationForComponentAtPath(mutatorPath, compFrom, designFile)
		mutatedValue := configurationForComponentAtPath(mutatedPath, compTo, designFile)

		strategy := getStrategyForValueAt(strategies, i)

		if !matchValuesWithStrategies(mutatorValue, mutatedValue, strategy) {
			return false
		}
	}

	return true
}

// getMatchStrategyForSelector returns the match strategy matrix from a selector.
func getMatchStrategyForSelector(fromClause map[string]interface{}) []interface{} {
	raw, ok := fromClause["match_strategy_matrix"]
	if !ok || raw == nil {
		return nil
	}
	arr, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	return arr
}

// getStrategyForValueAt returns the match strategy for a specific index.
func getStrategyForValueAt(strategies []interface{}, index int) []string {
	if index < len(strategies) && strategies[index] != nil {
		if arr, ok := strategies[index].([]interface{}); ok {
			return toStringSlice(arr)
		}
	}
	return []string{"equal"}
}

// extractValues extracts non-nil values from a component at the given reference paths.
func extractValues(component map[string]interface{}, refs []interface{}) []interface{} {
	var values []interface{}
	for _, ref := range refs {
		path := interfaceToStringSlice(ref)
		resolved := resolvePath(path, component)
		value := objectGetNested(component, resolved, nil)
		if value != nil {
			values = append(values, value)
		}
	}
	return values
}

// matchObjectValues checks if two value slices contain the same set of values.
func matchObjectValues(values1, values2 []interface{}) bool {
	if len(values1) != len(values2) {
		return false
	}
	set := make(map[string]bool, len(values1))
	for _, v := range values1 {
		set[fmt.Sprintf("%v", v)] = true
	}
	for _, v := range values2 {
		if !set[fmt.Sprintf("%v", v)] {
			return false
		}
	}
	return true
}
