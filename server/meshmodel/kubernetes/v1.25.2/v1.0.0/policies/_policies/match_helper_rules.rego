package relationship_evaluation_policy

import rego.v1

# based on the relationship kind, type and subtype accordingly, an object that will be used to patch is created.
perform_eval(design_file, relationship) := patched_declarations if {
	lower(relationship.kind) == "edge"
	lower(relationship.type) == "binding"

	# {
	# "declaration_id": <some id>,
	# "paths": [array of json paths]
	# }

	# In the design file relationships block,
	# if relationship exist there should be only 1 selector set and one component in from and to.

	patched_declarations := [result |
		some allowed_selector in relationship.selectors[0].allow
		mutator_object := extract_mutator_from_match(allowed_selector[0].match)

		to_mutate_object := extract_tomutate_from_match(allowed_selector[0].match)

		# print("MUTATOR OBJECT ", mutator_object)

		# print("TO MUTATE OBJECT ", to_mutate_object)
		# both mutator_value and mutated_value should be array of paths from which config manipulation will happen.
		# Both should be of same length.
		# In case it is not it is a bug in relationship_definition.
		range := numbers.range(0, min([count(mutator_object.paths), count(to_mutate_object.paths)]) - 1)

		# print("RANGE ", range)
		# declaration that gets mutated
		declaration_to_mutate := declaration_with_id(design_file, to_mutate_object.declaration_id)

		# declaration that mutates the other declaration with its own config
		mutating_declaration := declaration_with_id(design_file, mutator_object.declaration_id)

		# print("declaration_to_mutate: ", declaration_to_mutate)
		# iterate all the paths

		patches := [patch_object |
			some i in range
			resolved_mutated_path := resolve_path(to_mutate_object.paths[i], declaration_to_mutate)

			print("resolved_mutated_path ", resolved_mutated_path)

			resolved_mutator_path := resolve_path(mutator_object.paths[i], mutating_declaration)

			print("resolved_mutator_path ", resolved_mutator_path)

			update_value := object.get(mutating_declaration, resolved_mutator_path, null)

			# as operation is add I think we don't need to do ensure parent paths exist
			update_value != null
			patch_object := {
				"op": "add",
				"path": resolved_mutated_path,
				"value": update_value,
			}
		]

		resultant_patches_to_apply := ensureParentPathsExist(patches, declaration_to_mutate)

		mutated_declaration := json.patch(declaration_to_mutate, resultant_patches_to_apply)

		# print("mutated_declaration: ", mutated_declaration)
		result := {
			"declaration_id": declaration_to_mutate.id,
			"mutated_declaration": mutated_declaration,
		}
	]

	print("\n\nPATCHES  ", patched_declarations)
	# perform merge of the same object after returning from this rule in the evaluation rule.
	# based on the relationship it can happen that the same object is updated at differnet location hence merge all those changes.

	# resultant_patches_to_apply := ensureParentPathsExist(patches, declaration_to_mutate)
	# print("resultant_patches_to_apply: ", resultant_patches_to_apply)

}

extract_mutator_from_match(match) := value if {
	has_key(match.from[0], "mutatorRef")

	# print("line 58 : ", match.from[0].mutatorRef)

	value := {
		"paths": match.from[0].mutatorRef,
		"declaration_id": match.from[0].id,
	}
}

extract_mutator_from_match(match) := value if {
	has_key(match.to[0], "mutatorRef")

	# print("line 67 : ", match.to[0].mutatorRef)

	value := {
		"paths": match.to[0].mutatorRef,
		"declaration_id": match.to[0].id,
	}
}

extract_tomutate_from_match(match) := value if {
	has_key(match.from[0], "mutatedRef")

	# print("line 76 : ", match.from[0].mutatedRef)

	value := {
		"paths": match.from[0].mutatedRef,
		"declaration_id": match.from[0].id,
	}
}

extract_tomutate_from_match(match) := value if {
	has_key(match.to[0], "mutatedRef")

	# print("line 85 : ", match.to[0].mutatedRef)

	value := {
		"paths": match.to[0].mutatedRef,
		"declaration_id": match.to[0].id,
	}
}
