package policies

// RelationshipPolicy defines the interface each relationship type policy must implement.
type RelationshipPolicy interface {
	// Identifier returns the unique policy identifier string.
	Identifier() string

	// IsImplicatedBy returns true if the relationship belongs to this policy.
	IsImplicatedBy(rel map[string]interface{}) bool

	// IsInvalid checks whether an existing relationship is now invalid in the design.
	IsInvalid(rel, designFile map[string]interface{}) bool

	// AlreadyExists provides policy-specific duplicate detection (beyond the generic check).
	AlreadyExists(rel, designFile map[string]interface{}) bool

	// IdentifyRelationship identifies new relationships from a relationship definition against a design.
	IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{}

	// SideEffects returns additional actions for a relationship (e.g., adding/removing components).
	SideEffects(rel, designFile map[string]interface{}) []PolicyAction
}

// implicableRelationships filters relationships that belong to a specific policy.
func implicableRelationships(relationships []map[string]interface{}, policy RelationshipPolicy) []map[string]interface{} {
	var result []map[string]interface{}
	for _, rel := range relationships {
		if policy.IsImplicatedBy(rel) {
			result = append(result, rel)
		}
	}
	return result
}

// validateRelationshipsInDesign validates relationships in the design for a policy.
// Returns delete/update actions for relationships that are now invalid.
func validateRelationshipsInDesign(designFile map[string]interface{}, policy RelationshipPolicy) []PolicyAction {
	rels := extractMapSlice(designFile, "relationships")
	implicated := implicableRelationships(rels, policy)

	var actions []PolicyAction
	for _, rel := range implicated {
		if fromOrToComponentsDontExist(rel, designFile) && policy.IsInvalid(rel, designFile) {
			actions = append(actions, PolicyAction{
				Op: UpdateRelationshipOp,
				Value: map[string]interface{}{
					"id":    getMapString(rel, "id"),
					"path":  "/status",
					"value": "deleted",
				},
			})
		}
	}
	return actions
}

// identifyRelationshipsInDesign identifies new relationships in the design for a policy.
func identifyRelationshipsInDesign(
	designFile map[string]interface{},
	relationshipsInScope []map[string]interface{},
	policy RelationshipPolicy,
) []PolicyAction {
	implicated := implicableRelationships(relationshipsInScope, policy)

	var actions []PolicyAction
	seen := make(map[string]bool)

	for _, relDef := range implicated {
		identified := policy.IdentifyRelationship(relDef, designFile)

		for _, rel := range identified {
			// Check generic duplicate
			if relationshipAlreadyExists(designFile, rel) {
				continue
			}
			// Check policy-specific duplicate
			if policy.AlreadyExists(rel, designFile) {
				continue
			}

			id := getMapString(rel, "id")
			if seen[id] {
				continue
			}
			seen[id] = true

			actions = append(actions, PolicyAction{
				Op: AddRelationshipOp,
				Value: map[string]interface{}{
					"item": rel,
				},
			})
		}
	}
	return actions
}

// generateActionsToApplyOnDesign generates all actions for a policy.
func generateActionsToApplyOnDesign(designFile map[string]interface{}, policy RelationshipPolicy) []PolicyAction {
	rels := extractMapSlice(designFile, "relationships")
	implicated := implicableRelationships(rels, policy)

	var allActions []PolicyAction

	// Approve identified relationships
	allActions = append(allActions, approveIdentifiedRelationshipsAction(implicated, 100)...)

	// Cleanup deleted relationships
	allActions = append(allActions, cleanupDeletedRelationshipsActions(implicated)...)

	// Policy-specific side effects
	for _, rel := range implicated {
		allActions = append(allActions, policy.SideEffects(rel, designFile)...)
	}

	return allActions
}
