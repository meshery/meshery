package eval

import rego.v1

import data.eval_rules

hierarchical_parent_child_policy_identifier := "hierarchical_parent_child"

relationship_is_implicated_by_policy(relationship, policy_identifier) if {
	policy_identifier == hierarchical_parent_child_policy_identifier
	relationship.kind == "hierarchical"
	relationship.type == "parent"
	relationship.subType == "inventory"
}

relationship_is_invalid(relationship, design_file, policy_identifier) if {
	policy_identifier == hierarchical_parent_child_policy_identifier
	eval_rules.from_or_to_components_dont_exist(relationship, design_file)
}
