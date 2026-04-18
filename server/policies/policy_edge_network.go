package policies

import (
	"strings"

	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// EdgeNonBindingPolicy handles edge non-binding (network) relationships.
type EdgeNonBindingPolicy struct{}

func (p *EdgeNonBindingPolicy) Identifier() string {
	return "edge-non-binding"
}

func (p *EdgeNonBindingPolicy) IsImplicatedBy(rel *relationship.RelationshipDefinition) bool {
	return strings.EqualFold(string(rel.Kind), "edge") &&
		strings.EqualFold(rel.RelationshipType, "non-binding")
}

func (p *EdgeNonBindingPolicy) IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return fromOrToComponentsDontExist(rel, design)
}

func (p *EdgeNonBindingPolicy) AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool {
	return relationshipAlreadyExists(design, rel)
}

func (p *EdgeNonBindingPolicy) IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition {
	return identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, design)
}

func (p *EdgeNonBindingPolicy) SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction {
	if getRelStatus(rel) == StatusDeleted {
		return nil
	}
	return patchMutatorsAction(rel, design)
}
