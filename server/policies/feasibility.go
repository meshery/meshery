package policies

// isRelationshipFeasible checks if a selector matches a component (kind + model name).
func isRelationshipFeasible(selector, comp map[string]interface{}) bool {
	selectorKind := getMapString(selector, "kind")
	compComponent := getMapMap(comp, "component")
	compKind := getMapString(compComponent, "kind")

	if !matchName(compKind, selectorKind) {
		return false
	}

	selectorModel := getMapMap(selector, "model")
	compModel := getMapMap(comp, "model")

	if selectorModel != nil && compModel != nil {
		sModelName := getMapString(selectorModel, "name")
		cModelName := getMapString(compModel, "name")
		if sModelName != "" && !matchName(cModelName, sModelName) {
			return false
		}
	}

	return true
}

// isRelationshipFeasibleTo checks if a component can be the "to" side of a relationship.
// Returns the matching "to" selector, or nil.
func isRelationshipFeasibleTo(component map[string]interface{}, relationship map[string]interface{}) map[string]interface{} {
	selectors := getMapSlice(relationship, "selectors")
	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}
		toSelectors := getMapSlice(allow, "to")
		for _, ts := range toSelectors {
			toSel, ok := ts.(map[string]interface{})
			if !ok {
				continue
			}
			if isRelationshipFeasible(toSel, component) {
				return toSel
			}
		}
	}
	return nil
}

// isRelationshipFeasibleFrom checks if a component can be the "from" side of a relationship.
// Returns the matching "from" selector, or nil.
func isRelationshipFeasibleFrom(component map[string]interface{}, relationship map[string]interface{}) map[string]interface{} {
	selectors := getMapSlice(relationship, "selectors")
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
		for _, fs := range fromSelectors {
			fromSel, ok := fs.(map[string]interface{})
			if !ok {
				continue
			}
			if isRelationshipFeasible(fromSel, component) {
				return fromSel
			}
		}
	}
	return nil
}

// feasibleRelationshipSelectorBetween finds a feasible selector pair between two components.
// Returns the from and to selector pair, or nil if not feasible.
func feasibleRelationshipSelectorBetween(fromComp, toComp, relationship map[string]interface{}) map[string]interface{} {
	selectors := getMapSlice(relationship, "selectors")
	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}

		// Check deny selectors first
		deny := getMapMap(selectorSet, "deny")
		if deny != nil {
			denyTo := getMapSlice(deny, "to")
			denyFrom := getMapSlice(deny, "from")
			if isSelectorSetFeasibleBetween(toComp, fromComp, denyTo, denyFrom) != nil {
				continue // denied
			}
		}

		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}
		allowTo := getMapSlice(allow, "to")
		allowFrom := getMapSlice(allow, "from")

		result := isSelectorSetFeasibleBetween(toComp, fromComp, allowTo, allowFrom)
		if result != nil {
			return result
		}
	}
	return nil
}

func isSelectorSetFeasibleBetween(toComp, fromComp map[string]interface{}, toSelectors, fromSelectors []interface{}) map[string]interface{} {
	for _, ts := range toSelectors {
		toSel, ok := ts.(map[string]interface{})
		if !ok {
			continue
		}
		if !isRelationshipFeasible(toSel, toComp) {
			continue
		}

		for _, fs := range fromSelectors {
			fromSel, ok := fs.(map[string]interface{})
			if !ok {
				continue
			}
			if isRelationshipFeasible(fromSel, fromComp) {
				return map[string]interface{}{
					"from": fromSel,
					"to":   toSel,
				}
			}
		}
	}
	return nil
}

// isRelationshipDenied checks if a relationship between two declarations is denied.
func isRelationshipDenied(fromDecl, toDecl map[string]interface{}, denySelectors map[string]interface{}) bool {
	if denySelectors == nil {
		return false
	}

	denyFrom := getMapSlice(denySelectors, "from")
	denyTo := getMapSlice(denySelectors, "to")

	if len(denyFrom) == 0 || len(denyTo) == 0 {
		return false
	}

	for _, df := range denyFrom {
		dfSel, ok := df.(map[string]interface{})
		if !ok {
			continue
		}
		if !anySelectorMatches(fromDecl, dfSel) {
			continue
		}

		for _, dt := range denyTo {
			dtSel, ok := dt.(map[string]interface{})
			if !ok {
				continue
			}
			if anySelectorMatches(toDecl, dtSel) {
				return true
			}
		}
	}
	return false
}

// anySelectorMatches checks if a selector matches a declaration's kind and model.
func anySelectorMatches(declaration, selector map[string]interface{}) bool {
	return selectorAndDeclarationKindMatches(selector, declaration) &&
		selectorAndDeclarationModelMatches(selector, declaration)
}

func selectorAndDeclarationKindMatches(selector, declaration map[string]interface{}) bool {
	sKind := getMapString(selector, "kind")
	if sKind == "*" {
		return true
	}
	dComp := getMapMap(declaration, "component")
	dKind := getMapString(dComp, "kind")
	return sKind == dKind
}

func selectorAndDeclarationModelMatches(selector, declaration map[string]interface{}) bool {
	sModel := getMapMap(selector, "model")
	dModel := getMapMap(declaration, "model")
	if sModel == nil {
		return true
	}

	sName := getMapString(sModel, "name")
	dName := getMapString(dModel, "name")
	if sName != "*" && sName != dName {
		if !matchName(dName, sName) {
			return false
		}
	}

	// Check registrant
	sReg := sModel["registrant"]
	dReg := dModel["registrant"]

	switch sr := sReg.(type) {
	case string:
		if sr == "*" {
			return true
		}
		dr, ok := dReg.(string)
		if !ok {
			return false
		}
		return sr == dr
	case map[string]interface{}:
		dr, ok := dReg.(map[string]interface{})
		if !ok {
			return false
		}
		srKind := getMapString(sr, "kind")
		drKind := getMapString(dr, "kind")
		return srKind == drKind
	default:
		return true // no registrant constraint
	}
}
