package policies

import "strings"

// identifyAdditions finds components that should be added to the design based on
// inventory relationships. When a component references another component type that
// doesn't exist in the design, this suggests adding the missing component.
func identifyAdditions(relDef, designFile map[string]interface{}) []PolicyAction {
	if !strings.EqualFold(getMapString(relDef, "kind"), "hierarchical") ||
		!strings.EqualFold(getMapString(relDef, "type"), "parent") ||
		!strings.EqualFold(getMapString(relDef, "subType"), "inventory") {
		return nil
	}

	comps := getMapSlice(designFile, "components")
	selectors := getMapSlice(relDef, "selectors")

	var actions []PolicyAction
	seen := make(map[string]bool)

	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}

		mutatedSels := selectorsWithRef(allow, "mutatedRef")
		mutatorSels := selectorsWithRef(allow, "mutatorRef")

		for _, mutatedSel := range mutatedSels {
			mutatedPatch := getMapMap(mutatedSel, "patch")
			if mutatedPatch == nil {
				continue
			}
			mutatedRefs := getMapSlice(mutatedPatch, "mutatedRef")
			mutatedKind := getMapString(mutatedSel, "kind")
			mutatedComps := filterComponentsByKind(comps, mutatedKind)

			for _, mutatedComp := range mutatedComps {
				mutatedValues := extractValues(mutatedComp, mutatedRefs)
				if len(mutatedValues) != len(mutatedRefs) {
					continue
				}

				for _, mutatorSel := range mutatorSels {
					comp := processCompToAdd(designFile, relDef, mutatedComp, mutatedValues, mutatorSel)
					if comp == nil {
						continue
					}

					id := getMapString(comp, "id")
					if seen[id] {
						continue
					}
					seen[id] = true

					actions = append(actions, PolicyAction{
						Op: AddComponentOp,
						Value: map[string]interface{}{
							"item": comp,
						},
					})
				}
			}
		}
	}

	return actions
}

// processCompToAdd checks if a mutator component already exists for the given
// mutated values. If not, creates a new component suggestion.
func processCompToAdd(
	designFile, relDef, mutatedComp map[string]interface{},
	mutatedValues []interface{},
	mutatorSel map[string]interface{},
) map[string]interface{} {
	comps := getMapSlice(designFile, "components")
	mutatorKind := getMapString(mutatorSel, "kind")
	mutatorComps := filterComponentsByKind(comps, mutatorKind)

	mutatorPatch := getMapMap(mutatorSel, "patch")
	if mutatorPatch == nil {
		return nil
	}
	mutatorRefs := getMapSlice(mutatorPatch, "mutatorRef")

	// Check if any existing mutator component matches the mutated values.
	for _, mutatorComp := range mutatorComps {
		mutatorValues := extractValues(mutatorComp, mutatorRefs)
		if matchObjectValues(mutatedValues, mutatorValues) {
			return nil
		}
	}

	// No matching mutator exists. Build a new component suggestion.
	component := map[string]interface{}{
		"component": map[string]interface{}{"kind": mutatorKind},
		"model":     mutatorSel["model"],
	}

	// Check feasibility between mutated component and the new component.
	if feasibleRelationshipSelectorBetween(component, mutatedComp, relDef) == nil {
		return nil
	}

	// Set values at mutator paths.
	for i := 0; i < len(mutatorRefs) && i < len(mutatedValues); i++ {
		path := interfaceToStringSlice(mutatorRefs[i])
		setNestedValue(component, path, mutatedValues[i])
	}

	component["id"] = newUUID(component).String()
	return component
}

// selectorsWithRef returns selectors that contain the given ref type in their patch.
func selectorsWithRef(allow map[string]interface{}, refKey string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, key := range []string{"from", "to"} {
		for _, s := range getMapSlice(allow, key) {
			sel, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			patch := getMapMap(sel, "patch")
			if patch == nil {
				continue
			}
			if getMapSlice(patch, refKey) != nil {
				result = append(result, sel)
			}
		}
	}
	return result
}
