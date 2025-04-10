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

# Module:
# Purpose:
# Example:
# Process Flow:
#
# 1. Validation
#    Valid relationships require:

#
# 2. Identification
#
#
# 3. Actions
#    New  Relationships:
#
#    Approved Relationships:
#    Deleted Relationships:


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

matching_mutators(component_from , component_to , from_clause,to_clause,design_file) := matching_selectors if  {

     mutatorCount := count(from_clause.patch.mutatorRef)

#     print("mutator count", mutatorCount)

     # any matching not null mutator makes the relationship valid ( which is kind of wierd as we dont have proper
     # control to specific which ones can be null or not)
     some i in numbers.range(0, mutatorCount - 1)
         mutatorPath := from_clause.patch.mutatorRef[i]
         mutatedPath := to_clause.patch.mutatedRef[i]

         mutatorValue := core_utils.configuration_for_component_at_path(mutatorPath, component_from, design_file)
         mutatedValue := core_utils.configuration_for_component_at_path(mutatedPath, component_to, design_file)

#         print("mutator value", mutatorValue, "mutated value", mutatedValue)
         mutatorValue != null
         mutatedValue != null


         mutatorValue == mutatedValue

         print("same value", mutatorValue, "mutated value", mutatedValue)


    matching_selectors := {
        "from": [from_clause],
        "to": [to_clause],
    }

}


patch_mutators_action(relationship,design_file) := { action |
#      print("patch mutators action", relationship.kind)
      some selector in relationship.selectors
      from := selector.allow.from[0]
      to := selector.allow.to[0]
      some i in numbers.range(0, count(from.patch.mutatorRef) - 1)
      mutatorRef  := from.patch.mutatorRef[i]
      mutatedRef :=  to.patch.mutatedRef[i]
#      print("mutatorRef", mutatorRef, "mutatedRef", mutatedRef)
      from_component := component_declaration_by_id(design_file, from.id)
      mutatorValue :=  core_utils.configuration_for_component_at_path(mutatorRef, from_component, design_file)
#      print("mutator value", mutatorValue)

      action := {
        "op": actions.update_component_op,
        "value": {
            "id": to.id,
            "path": mutatedRef,
            "value": mutatorValue,
        }
      }
}

#patch_mutators_action(relationship,design_file) := actions

identify_edge_network_relationships(relationship,design_file) := { declaration |

    some selector_set in relationship.selectors
    some from_selector in selector_set.allow.from
    some to_selector in selector_set.allow.to

    some component_from in design_file.components
    some component_to in design_file.components

    not component_from.id == component_to.id

#    print("from_selector", from_selector.kind , component_from.component.kind,"to_selector", to_selector.kind , component_to.component.kind)

    component_from.component.kind == from_selector.kind
    component_to.component.kind == to_selector.kind

    matching_selectors := matching_mutators(component_from, component_to, from_selector, to_selector, design_file)

#    print("identfied rel", component_from.id, component_to.id,matching_selectors)

    selector_patch_declaration:=  from_selector.patch

#    print("selector patch declaration", matching_selectors)
    # create alias relationship declaration
    selector_declaration := {
        "allow": {
            "from": [json.patch(from_selector, [
                {
                    "op": "replace",
                    "path": "/id",
                    "value" : component_from.id
                },
            ])],
            "to": [json.patch(to_selector, [
                {
                    "op": "replace",
                    "path": "/id",
                    "value": component_to.id,
                },
            ])],
        },
        "deny": {},
    }


    declaration := json.patch(relationship, [
        {
            "op": "add",
            "path": "/selectors",
            "value": [selector_declaration],
        },
        {
            "op": "add",
            "path": "/id",
            "value": new_uuid(selector_declaration),
        },
        {
            "op": "replace",
            "path": "/status",
            "value": "identified",
        },
    ])
}

identify_relationships(design_file, relationships_in_scope, relationship_policy_identifier) := eval_results if {
	relationship_policy_identifier == edge_network_policy_identifier
	implicable_relationships := {
        relationship |
        some relationship in relationships_in_scope
        is_non_binding_edge_relationship(relationship)
    }

	eval_results := union({ new_relationships |
	   some relationship in implicable_relationships
	   identified_relationships := identify_edge_network_relationships(relationship, design_file)

	   print("identified rels edge",count(identified_relationships))

       new_relationships := {add_action |
            some rel in identified_relationships
            not eval_rules.relationship_already_exists(design_file, rel)

            add_action := {
               "op": actions.add_relationship_op,
               "value":{"item":rel}
            }
       }


	   print("identified rels edge dedup",count(new_relationships))

	})

	print("eval results", eval_results)

	print("Identifying relationships for policy identifier: ", relationship_policy_identifier,count(eval_results),count(implicable_relationships))

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
	implicable_relationships := {
        relationship |
        some relationship in design_file.relationships
        is_non_binding_edge_relationship(relationship)
    }

	relationships_to_add := eval_rules.approve_identified_relationships_action(implicable_relationships, 100)

    components_to_update := union({actions |
        some relationship in implicable_relationships
        actions := patch_mutators_action(relationship, design_file)
    })

#    print("Actions to be taken", components_to_update)

	result := ( relationships_to_add | components_to_update)

}
