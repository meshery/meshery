package relationship_evaluation_policy

import data.eval
import data.actions

import rego.v1

# METADATA
# entrypoint: true
# description: "Evaluates relationships in the design file and updates the design file with the results.
# Flow that the policy follows from evaluation to final trace:
# 1. Evaluate relationships in the design file
# 2. Identify relationships in the design file
# 3. Perform actions based on the identified relationships
# 4. Prepare the final design to return"
default rels_in_design_file := []

# Loads relationships from the design file if any exist.
rels_in_design_file := input.relationships if {
	count(input.relationships) > 0
}

# Filters out relationships that are pending.
filter_pending_relationships(rel, relationships) := rel if {
	every relationship in relationships {
		relationship.id == rel.id
	}
}

# scope for relationships to evaluate against
# NEEDS IMPROVEMENT: make this dynamic based on the models referenced in the design file

relationship_preference_key(rel) := sprintf("%s-%s-%s",[lower(rel.kind),lower(rel.type),lower(rel.subType)])

models_in_design :=  {component.model | some component in input.components }

is_rel_enabled(rel) := true if {
	not input.preferences.layers.relationships
}

is_rel_enabled(rel) := true if {
	rel_key := relationship_preference_key(rel)
    not input.preferences.layers.relationships[rel_key] == false
}

is_rel_disabled(rel) := true  if {
	not is_rel_enabled(rel)
}



relationships_to_evaluate_against := { rel |
    some rel in data.relationships
	some model in models_in_design
	rel_key := relationship_preference_key(rel)
	# print("rel_key",rel_key)
	model.name == rel.model.name
	is_rel_enabled(rel) == true
	# print("model implicated and rel is enabled",model.name,rel_key)
	# print("is_rel_enabled",rel_key,is_rel_enabled(rel))
}

# Main evaluation function that processes relationships and updates the design.
evaluate := eval_results if {

    print("model in design ",count(models_in_design))
	print("all registered rels count",count(data.relationships))
	print("rels to evaluate count",count(relationships_to_evaluate_against))

	

    # list of policies that will be used during evaluation
	relationship_policy_identifiers := [
        eval.match_labels_policy_identifier,
		eval.alias_policy_identifier,
		eval.edge_network_policy_identifier
	]


	# Iterate over relationships in the design file and resolve patches.
	resultant_patches := {patched_object |
		some rel in rels_in_design_file

        rel.subType in  {"inventory"}
		# Skip relationships with status "deleted".
		status := lower(rel.status)
		allowed_status := {"pending", "approved"}
		allowed_status[status]

		# Perform evaluation specific to the kind and subtype of the relationship.
		resultant_patch := perform_eval(input, rel)

		# Change status for pending relationships to "approved".
		r = json.patch(rel, [{
			"op": "replace",
			"path": "/status",
			"value": "approved",
		}])

		patched_object := {
			"patches": resultant_patch,
			"relationship": r,
		}
	}

	# Collect intermediate relationships after patching.
	intermediate_rels := [relationship |
		some val in resultant_patches
		relationship := val.relationship
	]

	# Update pending relationships with the intermediate results.
	updated_pending_rels := [rel |
		some relationship in rels_in_design_file
		rel := filter_relationship(relationship, intermediate_rels)
	]

	# Separate out declarations by id.
	intermediate_result := {x |
		some val in resultant_patches
		some nval in val.patches
		x := nval
	}

	# Group patches by declaration id.
	ans := group_by_id(intermediate_result)

	# Apply patches to mutated declarations.
	result := {mutated |
		some val in ans
		mutated_declaration := json.patch(val.declaration, val.patches)
		mutated := {
			"declaration_id": mutated_declaration.id,
			"declaration": mutated_declaration,
		}
	}

	# Update the design file with mutated declarations.
	design_file_with_updated_declarations := [declaration |
		some val in input.components
		declaration := filter_updated_declaration(val, result)
	]

	# Collect updated declarations.
	updated_declarations := [decl |
		some obj in result
		decl := obj.declaration
	]

	# Components configurations have been updated.
	updated_design_file := json.patch(input, [{
		"op": "replace",
		"path": "/components",
		"value": design_file_with_updated_declarations,
	}])

	# Identify additional components that need to be added.
	components_added := [result |
		some relationship in relationships_to_evaluate_against
		new_comps := identify_additions(updated_design_file, relationship)
		some new_comp in new_comps
		result := new_comp
	]

	# Concatenate new components to the design file components.
	final_set_of_comps := array.concat(updated_design_file.components, components_added)

	# Update the design file with new components.
	updated_design_file_with_new_comps := json.patch(updated_design_file, [{
		"op": "replace",
		"path": "/components",
		"value": final_set_of_comps,
	}])

	# Identify all valid relationships after updates.
	all_valid_relationships := union({result |
		some relationship in relationships_to_evaluate_against
		result := identify_relationship(updated_design_file_with_new_comps, relationship)
	})

	# The evaluate_relationships_added rule can work on the original design or on the updated design.
	# It is concerned only about relationships which have not changed until this point.

	relationships_added := evaluate_relationships_added(updated_pending_rels, all_valid_relationships)
	relationships_deleted := evaluate_relationships_deleted(updated_pending_rels, all_valid_relationships)


	# Concatenate updated relationships.
	final_rels_added := array.concat(updated_pending_rels, relationships_added)

	# Filter out relationships that have been deleted.
	final_rels_with_deletions := [rel |
		some relationship in final_rels_added
		rel := filter_relationship(relationship, relationships_deleted)
	]

	# Update the design file with the final set of relationships.
	final_design_file_old = json.patch(updated_design_file_with_new_comps, [{
		"op": "add",
		"path": "/relationships",
		"value": final_rels_with_deletions,
	}])

	# New Evaluation Flow

	design_file_to_evaluate := json.patch(final_design_file_old, [
		{
			"op": "replace",
			"path": "/relationships",
			"value": {rel | some rel in final_design_file_old.relationships},
		},
		{
			"op": "replace",
			"path": "/components",
			"value": {comp | some comp in final_design_file_old.components},
		},
	])


	#1. Validate Relationships
	validation_actions := union({rels |
		some identifier in relationship_policy_identifiers
		rels := eval.validate_relationships_in_design(design_file_to_evaluate, identifier)
		print("Validated by ",identifier,rels)
	})


	design_file_with_validated_rels := actions.apply_all_actions_to_design(
        design_file_to_evaluate,
        validation_actions,
    )
    
	
	print("Validation actions",count(design_file_to_evaluate.relationships),count(design_file_with_validated_rels.relationships))

	# 2. Identify relationships in the design file.
	new_identified_rels_actions := union({rels |
		some identifier in relationship_policy_identifiers
		rels := eval.identify_relationships_in_design(design_file_with_validated_rels, relationships_to_evaluate_against, identifier)
	})



    design_file_with_identified_rels := actions.apply_all_actions_to_design(
        design_file_to_evaluate,
        validation_actions,
    )
    


	print("New identified rels", count(new_identified_rels_actions))

	#3. Actions

	design_file_to_apply_actions := actions.apply_all_actions_to_design(
        design_file_with_identified_rels,
        new_identified_rels_actions,
    )

	print("All relationships", count(design_file_to_apply_actions.relationships))

	actions_to_apply := union({actions |
		some identifier in relationship_policy_identifiers
		actions := eval.generate_actions_to_apply_on_design(design_file_to_apply_actions, identifier)
		print("action for",identifier,count(actions))
	})

	print("All Actions",count(actions_to_apply))


	actions_response := trace_from_actions(actions_to_apply)
    design_to_return := actions.apply_all_actions_to_design(design_file_to_apply_actions, actions_to_apply)

	# Prepare the evaluation results with updated design and trace information.
	eval_results := {
		"design": design_to_return,
		"actions": actions_to_apply,
		"trace": {
			"componentsUpdated": [],
			"componentsAdded": array_to_set(components_added) | actions_response.components_to_add,
			"componentsRemoved":[],
			"relationshipsRemoved":[],
			"relationshipsAdded":[],
			"relationshipsUpdated":[],
#			"componentsRemoved": actions_response.components_to_delete,
#			"relationshipsAdded": array_to_set(relationships_added) | actions_response.relationships_to_add,
#			"relationshipsRemoved": array_to_set(relationships_deleted) | actions_response.relationships_to_delete,
#			"relationshipsUpdated": intermediate_rels,
		},
	}

	print("Evaluation complete")
}

trace_from_actions(response) := {
	"components_to_add": {action.value.item |
		some action in response
		action.op == actions.add_component_op
	},
#	"components_to_delete": {action.value |
#		some action in response
#		action.op == "delete_component"
#	},
#	"relationships_to_delete": {action.value |
#		some action in response
#		action.op == "delete_relationship"
#	},
#	"relationships_to_add": {action.value |
#		some action in response
#		action.op == "add_relationship"
#	},
}

# --- Post Processing Phase ---##
delete_components(all_comps, comps_to_delete) := new_comps if {
	count(comps_to_delete) > 0
	ids_to_delete := {comp.id | some comp in comps_to_delete}
	new_comps := {comp |
		some comp in all_comps
		not comp.id in ids_to_delete
	}
} else := all_comps

delete_relationships(all_rels, rels_to_delete) := new_rels if {
	count(rels_to_delete) > 0
	ids_to_delete := {rel.id | some rel in rels_to_delete}
	new_rels := {rel |
		some rel in all_rels
		not rel.id in ids_to_delete
	}
} else := all_rels

final_design_from_actions(old_design, actions_response) := new_design if {
	# Add new components to the design

	with_new_components := array_to_set(old_design.components) | actions_response.components_to_add
	final_components := delete_components(with_new_components, actions_response.components_to_delete)

	# Add new relationships to the design
	with_new_relationships := array_to_set(old_design.relationships) | actions_response.relationships_to_add
	final_relationships := delete_relationships(with_new_relationships, actions_response.relationships_to_delete)

	new_design := json.patch(old_design, [
		{
			"op": "replace",
			"path": "/components",
			"value": final_components,
		},
		{
			"op": "replace",
			"path": "/relationships",
			"value": final_relationships,
		},
	])
} else := old_design

# ----------------------------------------------#

# Returns the updated declaration if it exists; otherwise, returns the original declaration.
filter_updated_declaration(declaration, updated_declarations) := obj.declaration if {
	some obj in updated_declarations
	obj.declaration_id == declaration.id
} else := declaration

# Returns the updated relationship if it exists; otherwise, returns the original relationship.
filter_relationship(rel, relationships) := relationship if {
	some relationship in relationships
	relationship.id == rel.id
} else := rel
