package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils
import data.eval_rules
import data.actions


edge_network_policy_identifier := "edge-non-binding"


relationship_is_implicated_by_policy(relationship,policy_identifier) := true if {
    policy_identifier == edge_network_policy_identifier
	relationship.kind == "edge"
}


relationship_is_invalid(relationship,design_file,policy_identifier) := true if {
	policy_identifier == edge_network_policy_identifier
    eval_rules.from_or_to_components_dont_exist(relationship,design_file)
}


identify_relationship(rel_definition,design_file,policy_identifier) := identified_relationships if {
	policy_identifier == edge_network_policy_identifier
	identified_relationships := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(rel_definition, design_file)
}


relationship_side_effects(relationship,design_file,policy_identifier) := actions if {
   policy_identifier == edge_network_policy_identifier
   not relationship.status == "deleted"
   actions := eval_rules.patch_mutators_action(relationship, design_file)
}

