package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils
import data.eval_rules
import data.actions


edge_network_policy_identifier := "generic-mutator-policy"


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

### Action Phase
# action_phase(design_file, relationship_policy_identifier) := result if {
# 	relationship_policy_identifier == edge_network_policy_identifier
# 	implicable_relationships := edge_network_implicable_relationships(design_file.relationships)
    
# 	relationships_to_add := eval_rules.approve_identified_relationships_action(implicable_relationships, 100)

# 	rels_to_delete := eval_rules.cleanup_deleted_relationships_actions(implicable_relationships)
	
#     components_to_update := union({actions |
#         some relationship in implicable_relationships
#         actions := eval_rules.patch_mutators_action(relationship, design_file)
#     })

# 	result := ( relationships_to_add | components_to_update) | rels_to_delete
# }
