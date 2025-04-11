package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils

import data.eval_rules
import data.actions




relationship_is_implicated_by_policy(rel,policy_identifier) if {
	not true
}

implicable_relationships(relationships,policy_identifier) := { rel |
  some rel in relationships
  relationship_is_implicated_by_policy(rel,policy_identifier)
}

relationship_side_effects(relationship,design,policy_identifier) := actions if {
   not true
   actions := {}
}

identify_relationship(relationship,design_file,policy) := rels if { 
	not true
	rels := {}
}

relationship_already_exists(rel,design_file,policy_identifier) if {
	not true
}

relationship_is_invalid(rel,design_file,policy_identifier) if {
	not true
}


identify_relationships_in_design(design_file, relationships_in_scope, policy_identifier) := eval_results if {

	implicable_rels := implicable_relationships(relationships_in_scope,policy_identifier)
    
	eval_results := union({ new_relationships |
	   some relationship in implicable_rels
	   identified_relationships := identify_relationship(relationship, design_file,policy_identifier)

	   new_relationships := {add_action |
            some rel in identified_relationships
            not eval_rules.relationship_already_exists(design_file, rel)
			not relationship_already_exists(rel,design_file,policy_identifier)

            add_action := {
               "op": actions.add_relationship_op,
               "value":{"item":rel}
            }
       }
	})
}

## Validate
validate_relationships_in_design(design_file, policy_identifier) := result if {

	implicable_rels := implicable_relationships(design_file.relationships,policy_identifier)
    result := {action |
       some rel in implicable_rels
       eval_rules.from_or_to_components_dont_exist(rel,design_file)
	   relationship_is_invalid(rel,design_file,policy_identifier)

	   action := {
		"op": actions.update_relationship_op,
		"value": {
			"id": rel.id,
			"path": "/status",
			"value": "deleted"
		 }
	   }
    }

}

#
### Action Phase
generate_actions_to_apply_on_design(design_file, policy_identifier) := result if {
	implicable_rels := implicable_relationships(design_file.relationships,policy_identifier)
    
	relationships_to_add := eval_rules.approve_identified_relationships_action(implicable_rels, 100)
	rels_to_delete := eval_rules.cleanup_deleted_relationships_actions(implicable_rels)
	
    policy_specific_actions :=  union({ actions |
	  some rel in implicable_rels
	  actions := relationship_side_effects(rel,design_file,policy_identifier)
	})

	print("policy side effects",policy_identifier,count(policy_specific_actions))

	result := relationships_to_add | rels_to_delete | policy_specific_actions
}
