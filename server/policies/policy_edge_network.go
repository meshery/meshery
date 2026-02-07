package policies

import "strings"

// EdgeNonBindingPolicy handles edge non-binding (network) relationships.
type EdgeNonBindingPolicy struct{}

func (p *EdgeNonBindingPolicy) Identifier() string {
	return "edge-non-binding"
}

func (p *EdgeNonBindingPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "kind"), "edge") &&
		strings.EqualFold(getMapString(rel, "type"), "non-binding")
}

func (p *EdgeNonBindingPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	return fromOrToComponentsDontExist(rel, designFile)
}

func (p *EdgeNonBindingPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return false // uses default check
}

func (p *EdgeNonBindingPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	return identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, designFile)
}

func (p *EdgeNonBindingPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	status := getMapString(rel, "status")
	if status == "deleted" {
		return nil
	}
	return patchMutatorsAction(rel, designFile)
}
