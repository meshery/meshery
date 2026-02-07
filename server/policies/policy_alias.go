package policies

import (
	"fmt"
	"strings"
)

// AliasPolicy handles alias relationships (hierarchical parent alias).
type AliasPolicy struct{}

func (p *AliasPolicy) Identifier() string {
	return "alias_relationships_policy"
}

func (p *AliasPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return isAliasRelationship(rel)
}

func (p *AliasPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	return !isAliasRelationshipValid(rel, designFile)
}

func (p *AliasPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return aliasRelationshipAlreadyExists(designFile, rel)
}

func (p *AliasPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	var allIdentified []map[string]interface{}
	comps := getMapSlice(designFile, "components")

	for _, c := range comps {
		comp, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		if isRelationshipFeasibleTo(comp, relDef) == nil {
			continue
		}
		identified := identifyAliasRelationships(comp, relDef, designFile)
		allIdentified = append(allIdentified, identified...)
	}
	return allIdentified
}

func (p *AliasPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	status := getMapString(rel, "status")

	if status == "identified" || status == "pending" {
		return aliasAddComponentSideEffects(rel, designFile)
	}
	if status == "deleted" {
		return aliasDeleteComponentSideEffects(rel, designFile)
	}
	return nil
}

// isAliasRelationship checks if a relationship is an alias type.
func isAliasRelationship(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "kind"), "hierarchical") &&
		strings.EqualFold(getMapString(rel, "type"), "parent") &&
		strings.EqualFold(getMapString(rel, "subType"), "alias")
}

// isAliasRelationshipValid checks if an alias relationship is still valid.
func isAliasRelationshipValid(rel, designFile map[string]interface{}) bool {
	status := getMapString(rel, "status")

	if status == "pending" {
		return true
	}

	if status != "approved" {
		return false
	}

	// Check from component exists
	fromID := fromComponentID(rel)
	fromComp := componentDeclarationByID(designFile, fromID)
	if fromComp == nil {
		return false
	}

	// Check to component exists
	toID := toComponentID(rel)
	toComp := componentDeclarationByID(designFile, toID)
	if toComp == nil {
		return false
	}

	// Check the path in the to component is still present
	ref := aliasRefFromRelationship(rel)
	if ref == nil {
		return false
	}

	config := getComponentConfiguration(toComp, designFile)
	value := objectGetNested(config, popFirst(ref), nil)
	return value != nil
}

// aliasRefFromRelationship extracts the mutatorRef path from an alias relationship.
func aliasRefFromRelationship(rel map[string]interface{}) []string {
	selectors := getMapSlice(rel, "selectors")
	for _, s := range selectors {
		sel, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(sel, "allow")
		if allow == nil {
			continue
		}
		fromArr := getMapSlice(allow, "from")
		for _, f := range fromArr {
			from, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			patch := getMapMap(from, "patch")
			if patch == nil {
				continue
			}
			mutatorRef := getMapSlice(patch, "mutatorRef")
			if len(mutatorRef) > 0 {
				return interfaceToStringSlice(mutatorRef[0])
			}
		}
	}
	return nil
}

// identifyAliasRelationships identifies alias relationships for a component.
func identifyAliasRelationships(component, relDef, designFile map[string]interface{}) []map[string]interface{} {
	var identified []map[string]interface{}

	selectors := getMapSlice(relDef, "selectors")
	for _, s := range selectors {
		sel, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(sel, "allow")
		if allow == nil {
			continue
		}
		fromArr := getMapSlice(allow, "from")
		toArr := getMapSlice(allow, "to")

		for _, f := range fromArr {
			from, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			for _, t := range toArr {
				to, ok := t.(map[string]interface{})
				if !ok {
					continue
				}

				fromPatch := getMapMap(from, "patch")
				if fromPatch == nil {
					continue
				}
				mutatorRefs := getMapSlice(fromPatch, "mutatorRef")
				if len(mutatorRefs) == 0 {
					continue
				}

				ref := interfaceToStringSlice(mutatorRefs[0])
				paths := getArrayAwareConfigPaths(ref, component, designFile)

				for _, path := range paths {
					selectorPatchDecl := map[string]interface{}{
						"patchStrategy": "replace",
						"mutatorRef":    []interface{}{stringSliceToInterface(path)},
						"mutatedRef":    []interface{}{stringSliceToInterface(path)},
					}

					compID := newUUID(map[string]interface{}{"c": component, "s": selectorPatchDecl}).String()

					newFrom := deepCopyMap(from)
					newFrom["id"] = compID
					newFrom["patch"] = selectorPatchDecl

					newTo := deepCopyMap(to)
					newTo["id"] = getMapString(component, "id")
					newTo["patch"] = selectorPatchDecl

					selectorDecl := map[string]interface{}{
						"allow": map[string]interface{}{
							"from": []interface{}{newFrom},
							"to":   []interface{}{newTo},
						},
						"deny": map[string]interface{}{},
					}

					rel := deepCopyMap(relDef)
					rel["selectors"] = []interface{}{selectorDecl}
					rel["id"] = newUUID(selectorDecl).String()
					rel["status"] = "identified"

					identified = append(identified, rel)
				}
			}
		}
	}

	return identified
}

// getArrayAwareConfigPaths resolves paths, handling array wildcards.
func getArrayAwareConfigPaths(ref []string, component, design map[string]interface{}) [][]string {
	if isDirectReference(ref) {
		config := getComponentConfiguration(component, design)
		value := objectGetNested(config, popFirst(ref), nil)
		if value == nil {
			return nil
		}
		return [][]string{ref}
	}

	// Array reference (ends with "_")
	directRef := popLast(ref)
	config := getComponentConfiguration(component, design)
	items := objectGetNested(config, popFirst(directRef), nil)

	arr, ok := items.([]interface{})
	if !ok || len(arr) == 0 {
		return nil
	}

	var paths [][]string
	for i := range arr {
		if arr[i] == nil {
			continue
		}
		path := make([]string, len(directRef))
		copy(path, directRef)
		path = append(path, fmt.Sprintf("%d", i))
		paths = append(paths, path)
	}
	return paths
}

// aliasRelationshipAlreadyExists checks for duplicate alias relationships.
func aliasRelationshipAlreadyExists(designFile, relationship map[string]interface{}) bool {
	rels := getMapSlice(designFile, "relationships")

	for _, r := range rels {
		existing, ok := r.(map[string]interface{})
		if !ok {
			continue
		}
		if getMapString(existing, "status") == "deleted" {
			continue
		}
		if !isAliasRelationship(existing) {
			continue
		}

		existingTo := getFirstTo(existing)
		newTo := getFirstTo(relationship)

		if existingTo != nil && newTo != nil {
			if getMapString(existingTo, "kind") == getMapString(newTo, "kind") &&
				getMapString(existingTo, "id") == getMapString(newTo, "id") &&
				deepEqual(existingTo["patch"], newTo["patch"]) {
				return true
			}
		}
	}
	return false
}

// getFirstTo returns the first "to" selector from a relationship.
func getFirstTo(rel map[string]interface{}) map[string]interface{} {
	selectors := getMapSlice(rel, "selectors")
	for _, s := range selectors {
		sel, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(sel, "allow")
		if allow == nil {
			continue
		}
		toArr := getMapSlice(allow, "to")
		if len(toArr) > 0 {
			if to, ok := toArr[0].(map[string]interface{}); ok {
				return to
			}
		}
	}
	return nil
}

// aliasAddComponentSideEffects creates actions to add alias components.
func aliasAddComponentSideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	var actions []PolicyAction
	selectors := getMapSlice(rel, "selectors")

	for _, s := range selectors {
		sel, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(sel, "allow")
		if allow == nil {
			continue
		}
		fromArr := getMapSlice(allow, "from")

		for _, f := range fromArr {
			from, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			patch := getMapMap(from, "patch")
			if patch == nil {
				continue
			}
			mutatorRefs := getMapSlice(patch, "mutatorRef")
			if len(mutatorRefs) == 0 {
				continue
			}

			path := interfaceToStringSlice(mutatorRefs[0])
			length := len(path)
			displayName := ""
			if length >= 2 {
				displayName = fmt.Sprintf("%s.%s", path[length-2], path[length-1])
			}

			component := map[string]interface{}{
				"id":          getMapString(from, "id"),
				"component":   map[string]interface{}{"kind": getMapString(from, "kind")},
				"model":       from["model"],
				"displayName": displayName,
				"metadata":    map[string]interface{}{"isAnnotation": true},
			}

			actions = append(actions, PolicyAction{
				Op: AddComponentOp,
				Value: map[string]interface{}{
					"item": component,
				},
			})
		}
	}
	return actions
}

// aliasDeleteComponentSideEffects creates actions to delete alias components.
func aliasDeleteComponentSideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	var actions []PolicyAction
	selectors := getMapSlice(rel, "selectors")

	for _, s := range selectors {
		sel, ok := s.(map[string]interface{})
		if !ok {
			continue
		}
		allow := getMapMap(sel, "allow")
		if allow == nil {
			continue
		}
		fromArr := getMapSlice(allow, "from")

		for _, f := range fromArr {
			from, ok := f.(map[string]interface{})
			if !ok {
				continue
			}
			fromID := getMapString(from, "id")
			actions = append(actions, PolicyAction{
				Op: DeleteComponentOp,
				Value: map[string]interface{}{
					"id":        fromID,
					"component": componentDeclarationByID(designFile, fromID),
				},
			})
		}
	}
	return actions
}
