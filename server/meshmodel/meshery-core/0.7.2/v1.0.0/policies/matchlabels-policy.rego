package eval

import data.core_utils
import data.feasibility_evaluation_utils
import rego.v1

# Notes :
# these policies take advantage of set comprehensions to weedout duplicates at any stage
# so replacing them with array comprehensions will cause duplicacy issue
# subtype: matchlabels
# in each eval phase check the identifier matches
# whenerever going through all relationships in desing always check if it matches the kind,type,subtype for the policy
# because the relationships are never scoped down by default
# use partial rules for all of the exported eval stages so they dont conflict with others

MAX_MATCHLABELS := 20

is_matchlabel_relationship(relationship) if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "sibling"
}

# returns set {
# 	{
# 		field,
# 		value,
# 		components
# 	}
# }
identify_matchlabels(design_file, relationship) := all_match_labels if {
	field_pairs := {pair |
		is_matchlabel_relationship(relationship)

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

is_match_labels_policy_identifier(relationship_policy_identifier) if {
	relationship_policy_identifier == {
		"kind": "hierarchical",
		"type": "sibling",
		"subtype": "matchlabels",
	}
}

identify_relationships(design_file, relationships_in_scope, relationship_policy_identifier) := eval_results if {
	is_match_labels_policy_identifier(relationship_policy_identifier)

	eval_results := {new_relationship |
		some relationship in relationships_in_scope

		is_matchlabel_relationship(relationship)

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
				"value": static_uuid({selector_declaration, relationship_policy_identifier}),
			},
			{
				"op": "replace",
				"path": "/status",
				"value": "pending",
			},
		])
	}

	print("Identify match rels Eval results", count(eval_results))
}

# validate all relationships in the design file
# as matchlabels are stateless and have no sideeffects on the design , we just delete
# all the matchlabel relationships in validate phase and reidentify at next stage . this helps
# to bypass any expensive validation at this stage as we are eitherway going to do that at identication stage
validate_relationships_phase(design_file, relationship_policy_identifier) := result if {
	is_match_labels_policy_identifier(relationship_policy_identifier)

	result := {validated |
		some rel in design_file.relationships
		is_matchlabel_relationship(rel)
		validated := json.patch(rel, [{
			"op": "replace",
			"path": "/status",
			"value": "deleted",
		}])
	}
}

action_phase(design_file, relationship_policy_identifier) := result if {
	is_match_labels_policy_identifier(relationship_policy_identifier)

	relationships_to_add := {action |
		some rel in design_file.relationships
		is_matchlabel_relationship(rel)
		rel.status == "pending"
		rel_to_add := json.patch(rel, [{
			"op": "replace",
			"path": "/status",
			"value": "approved",
		}])

		action := {
			"op": "add_relationship",
			"value": rel_to_add,
		}
	}

	relationships_to_delete := {action |
		some rel in design_file.relationships
		is_matchlabel_relationship(rel)
		rel.status == "deleted"
		action := {
			"op": "delete_relationship",
			"value": rel,
		}
	}

	result := relationships_to_add | relationships_to_delete
}
