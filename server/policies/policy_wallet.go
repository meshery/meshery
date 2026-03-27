package policies

import "strings"

// HierarchicalWalletPolicy handles hierarchical parent wallet relationships
// (e.g., PodTemplate to Pod, EndpointSlice to Service).
type HierarchicalWalletPolicy struct{}

func (p *HierarchicalWalletPolicy) Identifier() string {
	return "hierarchical_wallet"
}

func (p *HierarchicalWalletPolicy) IsImplicatedBy(rel map[string]interface{}) bool {
	return strings.EqualFold(getMapString(rel, "kind"), "hierarchical") &&
		strings.EqualFold(getMapString(rel, "type"), "parent") &&
		strings.EqualFold(getMapString(rel, "subType"), "wallet")
}

func (p *HierarchicalWalletPolicy) IsInvalid(rel, designFile map[string]interface{}) bool {
	return fromOrToComponentsDontExist(rel, designFile)
}

func (p *HierarchicalWalletPolicy) AlreadyExists(rel, designFile map[string]interface{}) bool {
	return false
}

func (p *HierarchicalWalletPolicy) IdentifyRelationship(relDef, designFile map[string]interface{}) []map[string]interface{} {
	return identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields(relDef, designFile)
}

func (p *HierarchicalWalletPolicy) SideEffects(rel, designFile map[string]interface{}) []PolicyAction {
	status := getMapString(rel, "status")
	if status == "deleted" {
		return nil
	}
	return patchMutatorsAction(rel, designFile)
}
