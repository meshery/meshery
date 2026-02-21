package policies

import "strings"

// HierarchicalParentChildPolicy handles hierarchical parent-child (inventory) relationships.
type HierarchicalParentChildPolicy struct{}

func (p *HierarchicalParentChildPolicy) Identifier() string {
	return "hierarchical_parent_child"
}

func (p *HierarchicalParentChildPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "kind"), "hierarchical") &&
		strings.EqualFold(getMapString(rel, "type"), "parent") &&
		strings.EqualFold(getMapString(rel, "subType"), "inventory")
}

func (p *HierarchicalParentChildPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	return fromOrToComponentsDontExist(rel, designFile)
}

func (p *HierarchicalParentChildPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return false
}

func (p *HierarchicalParentChildPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	candidates := identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, designFile)

	// Filter by feasibility and deny selectors.
	var result []map[string]interface{}
	for _, rel := range candidates {
		fromID := fromComponentID(rel)
		toID := toComponentID(rel)
		fromComp := componentDeclarationByID(designFile, fromID)
		toComp := componentDeclarationByID(designFile, toID)
		if fromComp == nil || toComp == nil {
			continue
		}

		if feasibleRelationshipSelectorBetween(fromComp, toComp, relDef) == nil {
			continue
		}

		denied := false
		for _, s := range getMapSlice(relDef, "selectors") {
			selectorSet, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			deny := getMapMap(selectorSet, "deny")
			if deny != nil && isRelationshipDenied(fromComp, toComp, deny) {
				denied = true
				break
			}
		}
		if denied {
			continue
		}

		rel["status"] = "approved"
		result = append(result, rel)
	}
	return result
}

func (p *HierarchicalParentChildPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	status := getMapString(rel, "status")
	if status == "deleted" {
		return nil
	}
	return patchMutatorsAction(rel, designFile)
}
