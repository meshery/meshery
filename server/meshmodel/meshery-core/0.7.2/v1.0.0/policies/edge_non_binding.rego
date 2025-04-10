package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils

import data.core_utils.component_alias
import data.core_utils.component_declaration_by_id
import data.core_utils.from_component_id
import data.core_utils.get_array_aware_configuration_for_component_at_path
import data.core_utils.get_component_configuration
import data.core_utils.new_uuid
import data.core_utils.object_get_nested
import data.core_utils.pop_first
import data.core_utils.to_component_id
import data.feasibility_evaluation_utils.is_relationship_feasible_from
import data.feasibility_evaluation_utils.is_relationship_feasible_to
import data.core_utils.truncate_set
import data.core_utils.array_to_set

import data.eval_rules
import data.actions


edge_network_policy_identifier := {
    "kind": "edge",
    "type": "non-binding",
    "subtype": "network",
}

is_non_binding_edge_relationship(relationship) if {
    relationship.kind == edge_network_policy_identifier.kind
    relationship.type == edge_network_policy_identifier.type
    relationship.subType == edge_network_policy_identifier.subtype
}

edge_network_implicable_relationships(relationships) :=  { relationship |
    some relationship in relationships
    is_non_binding_edge_relationship(relationship)
}

identify_relationships(design_file, relationships_in_scope, relationship_policy_identifier) := eval_results if {
	relationship_policy_identifier == edge_network_policy_identifier
    print("coutn scope",count(relationships_in_scope))

	implicable_relationships := edge_network_implicable_relationships(relationships_in_scope)
    
	eval_results := union({ new_relationships |
	   some relationship in implicable_relationships
	   identified_relationships := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(relationship, design_file)

       new_relationships := {add_action |
            some rel in identified_relationships
            not eval_rules.relationship_already_exists(design_file, rel)

            add_action := {
               "op": actions.add_relationship_op,
               "value":{"item":rel}
            }
       }
	})
}

## Validate
## validate all relationships in the design file ( use partial rule so it doesnt conflict with other policies)
#validate_relationships_phase(design_file, relationship_policy_identifier) := result if {
#	relationship_policy_identifier == edge_network_policy_identifier
#    result := {rel |
#       some rel in design_file.relationships
#       is_non_binding_edge_relationship(rel)
#    }
#}

#
### Action Phase
action_phase(design_file, relationship_policy_identifier) := result if {
	relationship_policy_identifier == edge_network_policy_identifier
	implicable_relationships := edge_network_implicable_relationships(design_file.relationships)
    
	relationships_to_add := eval_rules.approve_identified_relationships_action(implicable_relationships, 100)

    components_to_update := union({actions |
        some relationship in implicable_relationships
        actions := eval_rules.patch_mutators_action(relationship, design_file)
    })

	result := ( relationships_to_add | components_to_update)
}
