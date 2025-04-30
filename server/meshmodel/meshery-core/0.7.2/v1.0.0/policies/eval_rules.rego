package eval_rules

import rego.v1

import data.core_utils.truncate_set
import  data.core_utils
import data.actions


approve_relationships_action(relationships,status,max_limit) := update_actions if {
   update_actions :=  truncate_set({action |
		some pending_rel in relationships
		pending_rel.status in status

		action := {
			"op": actions.update_relationship_op,
			"value": {
			   "id" : pending_rel.id,
			   "path": "/status",
			   "value": "approved",
			}
		}
	},max_limit)
}

approve_pending_relationships_action(relationships,max_limit) :=  approve_relationships_action(relationships,{"pending"},max_limit)
approve_identified_relationships_action(relationships,max_limit) :=  approve_relationships_action(relationships,{"identified"},max_limit)

cleanup_deleted_relationships_actions(relationships) := delete_actions if {
   delete_actions := {action |
		some rel in relationships
		rel.status == "deleted"
		action := {
			"op": actions.delete_relationship_op,
			"value": {"id": rel.id},
		}
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

      from_component := core_utils.component_declaration_by_id(design_file, from.id)
      to_component := core_utils.component_declaration_by_id(design_file, to.id)
    
      mutatorValue :=  core_utils.configuration_for_component_at_path(mutatorRef, from_component, design_file)
      old_value := core_utils.configuration_for_component_at_path(mutatedRef, to_component, design_file)

      old_value != mutatorValue

      action := {
        "op": actions.get_component_update_op(mutatedRef),
        "value": {
            "id": to.id,
            "path": mutatedRef,
            "value": mutatorValue,
        }
      }
}

same_relationship_identitfier(rel_a,rel_b) := true if {
    rel_a.kind == rel_b.kind
    rel_a.type == rel_b.type
    rel_a.subType == rel_b.subType
}

#selector clause is from or to
same_relationship_selector_clause(clause_a,clause_b) := true if {
   clause_a.kind == clause_b.kind
   clause_a.id == clause_b.id
   clause_a.patch == clause_b.patch
}

relationships_are_same(rel_a,rel_b) := true if {

    same_relationship_identitfier(rel_a,rel_b)

    some selector_a in rel_a.selectors
    some selector_b in rel_b.selectors

    same_relationship_selector_clause(selector_a.allow.from[0],selector_b.allow.from[0])
    same_relationship_selector_clause(selector_a.allow.to[0],selector_b.allow.to[0])
}

relationship_already_exists(design_file, relationship) := true if {
    some existing_rel in design_file.relationships
    existing_rel.status != "deleted" # check if the relationship is not deleted

    relationships_are_same(existing_rel, relationship)
}else := false



## Matching
#  Match supports multiple strategies
#  equal : complete equality of numbers,strigns,arrays,objects
#  equal_as_strings : coerses values to strings before equality
#  to_contains_from : checks if from is a superset of to ( for objects and arrays)
#  not_null : checks if both values are not nulll
#
#  These strategies can be grouped into a list to check if the values at mutator/mutated path match
#  And the the arrays can be grouped into a match_strategy to form a ordered list of strategies for each path

match_values(from_value,to_value,strategy) := true if {
    strategy == "equal"
    from_value == to_value
}


match_values(from_value,to_value,strategy) := true if {
    strategy == "equal_as_strings"
    s1 := sprintf("%v",[from_value])
    s2 :=  sprintf("%v",[to_value])

    s1 == s2
}


match_values(from_value,to_value,strategy) := true if {
    strategy == "to_contains_from"
    object.subset(to_value,from_value)
}

match_values(from_value,to_value,strategy) := true if {
    strategy == "not_null"
    from_value != null
    to_value != null
}

match_values_with_strategies(from_value,to_value,strategies) := true if {
    every strategy in strategies {
        match_values(from_value,to_value,strategy)
    }
}



get_match_strategy_for_selector(from_clause) := strategy if {
    strategy := from_clause.match_strategy_matrix
    strategy != null
}else := []

# if no strategies are specified for a path fallback to ["equal"]
get_strategy_for_value_at(strategies,index) := strategy if{
  strategy := strategies[index]
  strategy != null
}else := ["equal"]

matching_mutators(component_from , component_to , from_clause,to_clause,design_file) := matching_selectors if  {
     match_strategy_matrix := get_match_strategy_for_selector(from_clause)
     # print("match_strategy",match_strategy_matrix)
     mutatorCount := count(from_clause.patch.mutatorRef)

     every i in numbers.range(0, mutatorCount - 1) {
         mutatorPath := from_clause.patch.mutatorRef[i]
         mutatedPath := to_clause.patch.mutatedRef[i]

         mutatorValue := core_utils.configuration_for_component_at_path(mutatorPath, component_from, design_file)
         mutatedValue := core_utils.configuration_for_component_at_path(mutatedPath, component_to, design_file)

         strategies := get_strategy_for_value_at(match_strategy_matrix,i)
   
         matching := match_values_with_strategies(mutatorValue,mutatedValue,strategies)

    }

    matching_selectors := {
        "from": [from_clause],
        "to": [to_clause],
    }

}



identify_relationships_based_on_matching_mutator_and_mutated_fields(relationship,design_file) := { declaration |

    some selector_set in relationship.selectors
    some from_selector in selector_set.allow.from
    some to_selector in selector_set.allow.to

    some component_from in design_file.components
    some component_to in design_file.components

    not component_from.id == component_to.id

    # print("from_selector", from_selector.kind , component_from.component.kind,"to_selector", to_selector.kind , component_to.component.kind)

    component_from.component.kind == from_selector.kind
    component_to.component.kind == to_selector.kind

    matching_selectors := matching_mutators(component_from, component_to, from_selector, to_selector, design_file)

    # print("identfied rel", component_from.id, component_to.id,matching_selectors)

    selector_patch_declaration:=  from_selector.patch

    # print("selector patch declaration", matching_selectors)
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
            "value": core_utils.new_uuid(selector_declaration),
        },
        {
            "op": "replace",
            "path": "/status",
            "value": "identified",
        },
    ])
}



from_and_to_components_exist(relationship,design_file) := true if {

	# check if the from component is still present
	from_component := core_utils.component_declaration_by_id(design_file, core_utils.from_component_id(relationship))
	from_component != null

	# print("Is valid -> from_component", from_component)

	# check if the to component is still present
	to_component := core_utils.component_declaration_by_id(design_file, core_utils.to_component_id(relationship))
	to_component != null
}

from_or_to_components_dont_exist(relationship,design_file) := true if {
    not from_and_to_components_exist(relationship,design_file)
}


from_or_to_components_dont_exist(relationship,design_file) := false if {
    from_and_to_components_exist(relationship,design_file)
}


