package relationship_evaluation_policy

import rego.v1

apply_patch(design_file, mutator_object, to_mutate_object) := result if {
	# both mutator_object.paths and to_mutate_object.paths are paths to the config which gets manipulated.
	# Both should be of same length.
	# In case it is not, it is a bug in relationship_definition.

	range := numbers.range(0, min([count(mutator_object.paths), count(to_mutate_object.paths)]) - 1)

	# declaration that gets mutated
	declaration_to_mutate := declaration_with_id(design_file, to_mutate_object.declaration_id)

	# declaration that mutates the other declaration with its config
	mutating_declaration := declaration_with_id(design_file, mutator_object.declaration_id)

	# iterate all the paths
	patches := [patch_object |
		some i in range
		resolved_mutated_path := resolve_path(to_mutate_object.paths[i], declaration_to_mutate)

		resolved_mutator_path := resolve_path(mutator_object.paths[i], mutating_declaration)

		update_value := object.get(mutating_declaration, resolved_mutator_path, null)

		update_value != null
		patch_object := {
			"op": "add",
			"path": resolved_mutated_path,
			"value": update_value,
		}
	]

	resultant_patches_to_apply := ensureParentPathsExist(patches, declaration_to_mutate)

	mutated_declaration := json.patch(declaration_to_mutate, resultant_patches_to_apply)

	result := {
		"declaration_id": declaration_to_mutate.id,
		"mutated_declaration": mutated_declaration,
	}
}
