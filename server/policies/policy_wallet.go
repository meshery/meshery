package policies

import (
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// HierarchicalWalletPolicy handles hierarchical parent wallet relationships.
type HierarchicalWalletPolicy struct{}

func (p *HierarchicalWalletPolicy) Identifier() string {
	return "hierarchical_wallet"
}

func (p *HierarchicalWalletPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "hierarchical") &&
		strings.EqualFold(rel.RelationshipType, "parent") &&
		strings.EqualFold(rel.SubType, "wallet")
}

func (p *HierarchicalWalletPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return fromOrToComponentsDontExist(rel, design)
}

func (p *HierarchicalWalletPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *HierarchicalWalletPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, design)
}

func (p *HierarchicalWalletPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if getRelStatus(rel) == StatusDeleted {
		return nil
	}
	return patchMutatorsAction(rel, design)
}
