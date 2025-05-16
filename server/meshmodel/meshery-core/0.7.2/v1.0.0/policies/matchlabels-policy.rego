package eval

import data.core_utils
import data.feasibility_evaluation_utils
import rego.v1

import data.core_utils.component_alias
import data.core_utils.component_declaration_by_id
import data.core_utils.configuration_for_component_at_path
import data.core_utils.from_component_id
import data.core_utils.get_component_configuration
import data.core_utils.new_uuid
import data.core_utils.object_get_nested
import data.core_utils.pop_first
import data.core_utils.static_uuid
import data.core_utils.to_component_id
import data.core_utils.truncate_set
import data.feasibility_evaluation_utils.is_relationship_feasible_from
import data.feasibility_evaluation_utils.is_relationship_feasible_to

import data.actions
import data.eval_rules

# Notes :
# these policies take advantage of set comprehensions to weedout duplicates at any stage
# so replacing them with array comprehensions will cause duplicacy issue
# subtype: matchlabels
# in each eval phase check the identifier matches
# whenerever going through all relationships in desing always check if it matches the kind,type,subtype for the policy
# because the relationships are never scoped down by default
# use partial rules for all of the exported eval stages so they dont conflict with others

MAX_MATCHLABELS := 20

match_labels_policy_identifier := "sibling_match_labels_policy"


relationship_is_implicated_by_policy(relationship,policy_identifier) := true if {
    policy_identifier == match_labels_policy_identifier
	relationship.type == "sibling"
}


relationship_already_exists(rel,design_file,policy_identifier) := false if {
  policy_identifier == match_labels_policy_identifier
}

# these rels are stateless right now so just invalidate all then re identify
# to bypass any expensive validation at this stage as we are eitherway going to do that at identication stage
relationship_is_invalid(relationship,design_file,policy_identifier) := true if {
	policy_identifier == match_labels_policy_identifier
	relationship_is_implicated_by_policy(relationship,policy_identifier)
}


identify_relationship(rel_definition,design_file,policy_identifier) := identified_relationships if {
	policy_identifier == match_labels_policy_identifier
	identified_relationships := identify_matchlabel_relationships(rel_definition,design_file)
}

# Matchlabels specific logic

# returns set {
# 	{
# 		field,
# 		value,
# 		components
# 	}
# }
identify_matchlabels(design_file, relationship) := all_match_labels if {
	field_pairs := {pair |

		# identify against pair of components
		some component in design_file.components
		some other_component in design_file.components
		component.id != other_component.id

		# print("is alias rel", component)
		from := is_relationship_feasible_from(component, relationship)
		is_relationship_feasible_to(other_component, relationship)

		path := from.match.refs[0]

		# print("from mutator",path,configuration_for_component_at_path(path, component,design_file))
		some field, value in configuration_for_component_at_path(path, component, design_file)
		some field2, value2 in configuration_for_component_at_path(path, other_component, design_file)

		field == field2
		value == value2

		# print("matching fields",field,value)
		components := {component, other_component}

		pair := {
			"field": field,
			"value": value,
			"components": components,
		}
	}

	# merges the pairs into single groups on n components for each matching field ,value
	#  set comprehenshion to weed out duplicate matchlabels
	all_match_labels := {
	{
		"field": pair.field,
		"value": pair.value,
		"components": get_components_merged(pair.field, pair.value, field_pairs),
	} |
		some pair in field_pairs
	}

}

# merges the pairs into single groups on n components for each matching field ,value
# set comprehension to weed out duplicates in components
get_components_merged(field, value, field_pairs) := {component |
	some pair in field_pairs
	pair.field == field
	pair.value == value
	some component in pair.components
}


identify_matchlabel_relationships(relationship,design_file) := identified_rels if {
	

	identified_rels := {new_relationship |
		
		# limit matchlabel relationships
		identified_matchlabels := truncate_set(identify_matchlabels(design_file, relationship), MAX_MATCHLABELS)
		some match_label in identified_matchlabels

		# these need to be set comprehensions to automatically weed out duplicated declarations
		from_selectors := {from_selector |
			some component in match_label.components
			from_selector := {
				"id": component.id,
				"kind": component.component.kind,
				"model": component.model,
				"patch": {
					"patchStrategy": "replace",
					"mutatorRef": [[match_label.field]],
				},
			}
		}

		selector_declaration := {
			"allow": {
				# from and to both need to be defined for the visualization
				"from": from_selectors,
				"to": from_selectors,
			},
			"deny": {
				"from": [],
				"to": [],
			},
		}

		new_relationship := json.patch(relationship, [
			{
				"op": "replace",
				"path": "/selectors",
				"value": [selector_declaration],
			},
			{
				"op": "replace",
				"path": "/id",
				# use relationship_policy_identifier instead of relationship as seed to prevent duplicate
				# relationships being created (due to duplicate registered relationships)
				"value": static_uuid({selector_declaration, match_labels_policy_identifier}),
			},
			{
				"op": "replace",
				"path": "/status",
				"value": "identified",
			},
		])
	}
}
