package policies

import (
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// HierarchicalParentChildPolicy handles hierarchical parent-child (inventory) relationships.
type HierarchicalParentChildPolicy struct{}

func (p *HierarchicalParentChildPolicy) Identifier() string {
	return "hierarchical_parent_child"
}

func (p *HierarchicalParentChildPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "hierarchical") &&
		strings.EqualFold(rel.RelationshipType, "parent") &&
		strings.EqualFold(rel.SubType, "inventory")
}

func (p *HierarchicalParentChildPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return fromOrToComponentsDontExist(rel, design)
}

func (p *HierarchicalParentChildPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *HierarchicalParentChildPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	candidates := identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, design)

	var result []*relationship.RelationshipDefinition
	for _, rel := range candidates {
		fID := fromComponentID(rel)
		tID := toComponentID(rel)
		fromComp := componentDeclarationByID(design, fID)
		toComp := componentDeclarationByID(design, tID)
		if fromComp == nil || toComp == nil {
			continue
		}

		if feasibleRelationshipSelectorBetween(fromComp, toComp, relDef) == nil {
			continue
		}

		denied := false
		if relDef.Selectors != nil {
			for _, ss := range *relDef.Selectors {
				if ss.Deny != nil && isRelationshipDenied(fromComp, toComp, ss.Deny) {
					denied = true
					break
				}
			}
		}
		if denied {
			continue
		}

		setRelStatus(rel, StatusApproved)
		result = append(result, rel)
	}
	return result
}

func (p *HierarchicalParentChildPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if getRelStatus(rel) == StatusDeleted {
		return nil
	}
	return patchMutatorsAction(rel, design)
}
