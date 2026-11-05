package policies

import (
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// RelationshipPolicy defines the interface each relationship type policy must implement.
type RelationshipPolicy interface {
	Identifier() string
	IsImplicatedBy(rel *relationship.RelationshipDefinition) bool
	IsInvalid(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool
	AlreadyExists(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) bool
	IdentifyRelationship(relDef *relationship.RelationshipDefinition, design *pattern.PatternFile) []*relationship.RelationshipDefinition
	SideEffects(rel *relationship.RelationshipDefinition, design *pattern.PatternFile) []PolicyAction
}

// implicableRelationships filters relationships that belong to a specific policy.
func implicableRelationships(rels []*relationship.RelationshipDefinition, policy RelationshipPolicy) []*relationship.RelationshipDefinition {
	var result []*relationship.RelationshipDefinition
	for _, rel := range rels {
		if policy.IsImplicatedBy(rel) {
			result = append(result, rel)
		}
	}
	return result
}

// validateRelationshipsInDesign validates relationships in the design for a policy.
func validateRelationshipsInDesign(design *pattern.PatternFile, policy RelationshipPolicy) []PolicyAction {
	implicated := implicableRelationships(design.Relationships, policy)

	var actions []PolicyAction
	for _, rel := range implicated {
		if fromOrToComponentsDontExist(rel, design) && policy.IsInvalid(rel, design) {
			actions = append(actions, newUpdateRelationshipAction(rel.ID.String(), "/status", StatusDeleted))
		}
	}
	return actions
}

// identifyRelationshipsInDesign identifies new relationships in the design for a policy.
func identifyRelationshipsInDesign(
	design *pattern.PatternFile,
	relationshipsInScope []*relationship.RelationshipDefinition,
	policy RelationshipPolicy,
) []PolicyAction {
	implicated := implicableRelationships(relationshipsInScope, policy)

	var actions []PolicyAction
	seen := make(map[string]bool)

	for _, relDef := range implicated {
		identified := policy.IdentifyRelationship(relDef, design)

		for _, rel := range identified {
			if policy.AlreadyExists(rel, design) {
				continue
			}

			id := rel.ID.String()
			if seen[id] {
				continue
			}
			seen[id] = true

			actions = append(actions, newAddRelationshipAction(rel))
		}
	}
	return actions
}

// generateActionsToApplyOnDesign generates all actions for a policy.
func generateActionsToApplyOnDesign(design *pattern.PatternFile, policy RelationshipPolicy) []PolicyAction {
	implicated := implicableRelationships(design.Relationships, policy)

	var allActions []PolicyAction
	allActions = append(allActions, approveIdentifiedRelationshipsAction(implicated, 100)...)
	allActions = append(allActions, cleanupDeletedRelationshipsActions(implicated)...)

	for _, rel := range implicated {
		allActions = append(allActions, policy.SideEffects(rel, design)...)
	}

	return allActions
}
