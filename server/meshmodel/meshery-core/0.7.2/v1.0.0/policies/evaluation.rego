package relationship_evaluation_policy

import rego.v1

# Should be loaded always
# This is the first file that is executed that is this is the entry point of the policy
# The flow that the policy follows from evaluation to final trace
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
# TODO: make this dynamic based on the models referenced in the design file
relationships_to_evaluate_against := data.relationships

# Main evaluation function that processes relationships and updates the design.
evaluate := eval_results if {
	print("Wooo in evbal")

	# Iterate over relationships in the design file and resolve patches.
	resultant_patches := {patched_object |
		some rel in rels_in_design_file

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
		some relationship in data.relationships
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
		some relationship in data.relationships
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
	final_design_file = json.patch(updated_design_file_with_new_comps, [{
		"op": "add",
		"path": "/relationships",
		"value": final_rels_with_deletions,
	}])

	# New Evaluation Flow

	# 1. Identify relationships in the design file.

	design_file_to_evaluate := json.patch(input, [
		{
			"op": "replace",
			"path": "/relationships",
			"value": {rel | some rel in input.relationships},
		},
		{
			"op": "replace",
			"path": "/components",
			"value": {comp | some comp in input.components},
		},
	])

	new_identified_rels := identify_relationships(design_file_to_evaluate, relationships_to_evaluate_against)

	#2. Validate Relationships
	# newly identified relationships dont need to be validated ( as they are valid or pending)
	validated_rels := validate_relationships_phase(design_file_to_evaluate)

	print("New identified rels", count(new_identified_rels))
	print("Validated rels", count(validated_rels))

	#3. Actions

	design_file_to_apply_actions := json.patch(design_file_to_evaluate, [{
		"op": "replace",
		"path": "/relationships",
		"value": new_identified_rels | validated_rels,
	}])

	print("All relationships", count(design_file_to_apply_actions.relationships))

	actions_response := action_phase(design_file_to_apply_actions)

	# print("Final Design", final_evaluated_design_file)

	design_to_return := json.patch(final_design_file, [
		{
			"op": "replace",
			"path": "/relationships",
			"value": array.concat(set_to_array(actions_response.relationships_added), final_design_file.relationships),
		},
		{
			"op": "replace",
			"path": "/components",
			"value": array.concat(set_to_array(actions_response.components_added), final_design_file.components),
		},
	])

	# Prepare the evaluation results with updated design and trace information.
	eval_results := {
		"design": design_to_return,
		"trace": {
			"componentsUpdated": updated_declarations,
			"componentsAdded": {x | some x in components_added} | actions_response.components_added,
			"componentsRemoved": actions_response.components_deleted,
			"relationshipsAdded": relationships_added,
			"relationshipsRemoved": {x | some x in relationships_deleted} | actions_response.relationships_deleted,
			"relationshipsUpdated": intermediate_rels,
		},
	}
}

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
