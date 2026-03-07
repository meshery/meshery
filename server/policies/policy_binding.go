package policies

import "strings"

// EdgeBindingPolicy handles edge binding relationships (3-party: from, binding, to).
type EdgeBindingPolicy struct{}

func (p *EdgeBindingPolicy) Identifier() string {
	return "edge_binding"
}

func (p *EdgeBindingPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "kind"), "edge") &&
		strings.EqualFold(getMapString(rel, "type"), "binding")
}

func (p *EdgeBindingPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	return fromOrToComponentsDontExist(rel, designFile)
}

func (p *EdgeBindingPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return false
}

func (p *EdgeBindingPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	return identifyBindingRelationships(relDef, designFile)
}

func (p *EdgeBindingPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	return nil
}

// identifyBindingRelationships identifies 3-party binding relationships.
// For each (from, binding, to) triple, validates that from-binding and binding-to match
// via the selector's match field, then creates an identified relationship.
func identifyBindingRelationships(relDef, designFile map[string]interface{}) []map[string]interface{} {
	comps := getMapSlice(designFile, "components")
	selectors := getMapSlice(relDef, "selectors")
	var identified []map[string]interface{}

	for _, s := range selectors {
		selectorSet, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(selectorSet, "allow")
		if allow == nil {
			continue
		}
		deny := getMapMap(selectorSet, "deny")

		fromSelectorArr := getMapSlice(allow, "from")
		toSelectorArr := getMapSlice(allow, "to")

		// Index from selectors by component kind.
		fromSelByKind := make(map[string]map[string]interface{})
		for _, fs := range fromSelectorArr {
			fromSel, ok := fs.(map[string]interface{})
			if !ok {
				continue
			}
			fromSelByKind[getMapString(fromSel, "kind")] = fromSel
		}

		// Collect binding component kinds from match selectors.
		bindingKinds := collectBindingKinds(fromSelectorArr)

		// Index to selectors by "toKind#bindingKind".
		toSelByPair := indexToSelectorsByKindPair(toSelectorArr)

		// Pre-filter from and to components.
		fromComps := filterComponentsByKindSet(comps, fromSelByKind)
		toComps := filterComponentsByToSelectors(comps, toSelectorArr)

		for bindingKind := range bindingKinds {
			bindingDecls := filterComponentsByKind(comps, bindingKind)
			if len(bindingDecls) == 0 {
				continue
			}

			for _, fromDecl := range fromComps {
				fromKind := getMapString(getMapMap(fromDecl, "component"), "kind")
				fromSel, exists := fromSelByKind[fromKind]
				if !exists {
					continue
				}

				for _, bindingDecl := range bindingDecls {
					if getMapString(fromDecl, "id") == getMapString(bindingDecl, "id") {
						continue
					}
					if !isValidBinding(fromDecl, bindingDecl, fromSel) {
						continue
					}

					for _, toDecl := range toComps {
						if getMapString(toDecl, "id") == getMapString(bindingDecl, "id") {
							continue
						}

						toKind := getMapString(getMapMap(toDecl, "component"), "kind")
						toSel, exists := toSelByPair[toKind+"#"+bindingKind]
						if !exists {
							continue
						}

						if deny != nil && isRelationshipDenied(fromDecl, toDecl, deny) {
							continue
						}

						if !isValidBinding(bindingDecl, toDecl, toSel) {
							continue
						}

						fromDeclID := getMapString(fromDecl, "id")
						bindingDeclID := getMapString(bindingDecl, "id")
						toDeclID := getMapString(toDecl, "id")

						newFromSel := deepCopyMap(fromSel)
						patchBindingSelectorIDs(newFromSel, fromDeclID, fromDeclID, bindingDeclID)

						newToSel := deepCopyMap(toSel)
						patchBindingSelectorIDs(newToSel, toDeclID, bindingDeclID, toDeclID)

						selectorDecl := map[string]interface{}{
							"allow": map[string]interface{}{
								"from": []interface{}{newFromSel},
								"to":   []interface{}{newToSel},
							},
						}

						seed := map[string]interface{}{
							"from":    fromDeclID,
							"binding": bindingDeclID,
							"to":      toDeclID,
							"relId":   getMapString(relDef, "id"),
						}

						rel := deepCopyMap(relDef)
						rel["id"] = newUUID(seed).String()
						rel["selectors"] = []interface{}{selectorDecl}
						rel["status"] = "approved"

						identified = append(identified, rel)
					}
				}
			}
		}
	}

	return identified
}

// collectBindingKinds extracts binding component kinds from from selectors' match fields.
func collectBindingKinds(fromSelectorArr []interface{}) map[string]bool {
	kinds := make(map[string]bool)
	for _, fs := range fromSelectorArr {
		fromSel, ok := fs.(map[string]interface{})
		if !ok {
			continue
		}
		match := getMapMap(fromSel, "match")
		if match == nil {
			continue
		}
		for _, key := range []string{"from", "to"} {
			for _, mv := range getMapSlice(match, key) {
				mSel, ok := mv.(map[string]interface{})
				if !ok {
					continue
				}
				kind := getMapString(mSel, "kind")
				if kind != "" && !strings.EqualFold(kind, "self") {
					kinds[kind] = true
				}
			}
		}
	}
	return kinds
}

// indexToSelectorsByKindPair indexes to selectors by "toKind#matchKind".
func indexToSelectorsByKindPair(toSelectorArr []interface{}) map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	for _, ts := range toSelectorArr {
		toSel, ok := ts.(map[string]interface{})
		if !ok {
			continue
		}
		toKind := getMapString(toSel, "kind")
		match := getMapMap(toSel, "match")
		if match == nil {
			continue
		}
		for _, key := range []string{"from", "to"} {
			for _, mv := range getMapSlice(match, key) {
				mSel, ok := mv.(map[string]interface{})
				if !ok {
					continue
				}
				matchKind := getMapString(mSel, "kind")
				if matchKind != "" && !strings.EqualFold(matchKind, "self") {
					result[toKind+"#"+matchKind] = toSel
				}
			}
		}
	}
	return result
}

// patchBindingSelectorIDs sets the IDs in a binding selector and its match field.
func patchBindingSelectorIDs(sel map[string]interface{}, selectorID, matchFromID, matchToID string) {
	sel["id"] = selectorID
	match := getMapMap(sel, "match")
	if match == nil {
		return
	}
	fromArr := getMapSlice(match, "from")
	if len(fromArr) > 0 {
		if f, ok := fromArr[0].(map[string]interface{}); ok {
			f["id"] = matchFromID
		}
	}
	toArr := getMapSlice(match, "to")
	if len(toArr) > 0 {
		if t, ok := toArr[0].(map[string]interface{}); ok {
			t["id"] = matchToID
		}
	}
}

// isValidBinding checks if two components match via a binding selector's match field.
func isValidBinding(comp1, comp2, selector map[string]interface{}) bool {
	match := getMapMap(selector, "match")
	if match == nil {
		return false
	}

	mutatorDecl, mutatorPaths := extractMutatorFromMatch(match, comp1, comp2)
	mutatedDecl, mutatedPaths := extractMutatedFromMatch(match, comp1, comp2)
	if mutatorDecl == nil || mutatedDecl == nil {
		return false
	}
	if len(mutatorPaths) == 0 || len(mutatedPaths) == 0 {
		return false
	}

	count := len(mutatorPaths)
	if len(mutatedPaths) < count {
		count = len(mutatedPaths)
	}

	for i := 0; i < count; i++ {
		fromPath := interfaceToStringSlice(mutatorPaths[i])
		toPath := interfaceToStringSlice(mutatedPaths[i])

		resolvedFrom := resolvePath(fromPath, mutatorDecl)
		resolvedTo := resolvePath(toPath, mutatedDecl)

		val1 := objectGetNested(mutatorDecl, resolvedFrom, nil)
		val2 := objectGetNested(mutatedDecl, resolvedTo, nil)

		if !deepEqual(val1, val2) {
			return false
		}
	}

	return true
}

// extractMutatorFromMatch finds the mutator declaration and ref paths from a match field.
func extractMutatorFromMatch(match, comp1, comp2 map[string]interface{}) (map[string]interface{}, []interface{}) {
	fromArr := getMapSlice(match, "from")
	if len(fromArr) > 0 {
		if from, ok := fromArr[0].(map[string]interface{}); ok {
			if refs := getMapSlice(from, "mutatorRef"); refs != nil {
				return comp1, refs
			}
		}
	}
	toArr := getMapSlice(match, "to")
	if len(toArr) > 0 {
		if to, ok := toArr[0].(map[string]interface{}); ok {
			if refs := getMapSlice(to, "mutatorRef"); refs != nil {
				return comp2, refs
			}
		}
	}
	return nil, nil
}

// extractMutatedFromMatch finds the mutated declaration and ref paths from a match field.
func extractMutatedFromMatch(match, comp1, comp2 map[string]interface{}) (map[string]interface{}, []interface{}) {
	fromArr := getMapSlice(match, "from")
	if len(fromArr) > 0 {
		if from, ok := fromArr[0].(map[string]interface{}); ok {
			if refs := getMapSlice(from, "mutatedRef"); refs != nil {
				return comp1, refs
			}
		}
	}
	toArr := getMapSlice(match, "to")
	if len(toArr) > 0 {
		if to, ok := toArr[0].(map[string]interface{}); ok {
			if refs := getMapSlice(to, "mutatedRef"); refs != nil {
				return comp2, refs
			}
		}
	}
	return nil, nil
}

// filterComponentsByKind returns components matching the given kind.
func filterComponentsByKind(comps []interface{}, kind string) []map[string]interface{} {
	var result []map[string]interface{}
	for _, ci := range comps {
		comp, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}
		compKind := getMapString(getMapMap(comp, "component"), "kind")
		if matchName(compKind, kind) {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByKindSet returns components whose kind is in the given set.
func filterComponentsByKindSet(comps []interface{}, kindSet map[string]map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, ci := range comps {
		comp, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}
		compKind := getMapString(getMapMap(comp, "component"), "kind")
		if _, exists := kindSet[compKind]; exists {
			result = append(result, comp)
		}
	}
	return result
}

// filterComponentsByToSelectors returns components matching any to selector kind.
func filterComponentsByToSelectors(comps, toSelectorArr []interface{}) []map[string]interface{} {
	toKinds := make(map[string]bool)
	for _, ts := range toSelectorArr {
		toSel, ok := ts.(map[string]interface{})
		if !ok {
			continue
		}
		toKinds[getMapString(toSel, "kind")] = true
	}

	var result []map[string]interface{}
	for _, ci := range comps {
		comp, ok := ci.(map[string]interface{})
		if !ok {
			continue
		}
		compKind := getMapString(getMapMap(comp, "component"), "kind")
		if toKinds[compKind] || toKinds["*"] {
			result = append(result, comp)
		}
	}
	return result
}
