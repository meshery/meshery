package relationship_evaluation_policy

import rego.v1

# based on the relationship kind, type and subtype accordingly, an object that will be used to patch is created.
build_patch_object(design_file, relationship_from, relationship_to) if {
	has_key(relationship_from, "match")
	has_key(relationship_to, "match") # the 2nd check might be redundant. As if one of the selector has "match"
	# then according to current relationship defs, the other selector should also have "match" attribute

	# {
	# "declaration_id": <some id>,
	# "paths": [array of json paths]
	# }
	mutator_object := extract_mutator_from_match(relationship_from.match, relationship_to.match)

	to_mutate_object := extract_to_mutate_from_match(relationship_from.match, relationship_to.match)

	print("MUTATOR OBJECT ", mutator_object, count(mutator_object.paths), count(to_mutate_object.paths))

	print("TO MUTATE OBJECT ", to_mutate_object)
	# both mutator_value and mutated_value should be array of paths from which config manipulation will happen.
	# Both should be of same length.
	# In case it is not it is a bug in relationship_definition.
	range := numbers.range(0, min([count(mutator_object.paths), count(to_mutate_object.paths)]) - 1)

	# declaration that gets mutated
	declaration_to_mutate := declaration_with_id(design_file, to_mutate_object.declaration_id)

	# declaration that mutates the other declaration with its own config
	mutating_declaration := declaration_with_id(design_file, mutator_object.declaration_id)
	# iterate all the paths
	# patches := [patch_object |
	# 	some i in range

	# 	resolved_mutated_path := resolve_path(to_mutate_object[i], declaration_to_mutate)

	# 	resolved_mutator_path := resolve_path(mutator_object[i], mutating_declaration)

	# 	update_value := object.get(mutating_declaration, resolved_mutator_path, null)

	# 	update_value != null
	# 	patch_object := {
	# 		"op": "copy",
	# 		"from": resolved_mutator_path,
	# 		"value": update_value,
	# 	}
	# ]

	# resultant_patches_to_apply := ensureParentPathsExist(patches, declaration_to_mutate)
	# mutated_declaration = json.patch(declaration_to_mutate, resultant_patches_to_apply)

	# result := {
	# 	"declaration_id": declaration_to_mutate.id,
	# 	"mutated_declaration": mutated_declaration,
	# }
}

extract_mutator_from_match(match_from, match_to) := value if {
	print("line 58 : ", match_from.from[0])
	has_key(match_from.from[0], "mutatorRef")
	value := {
		"paths": match_from.from[0].mutatorRef,
		"declaration_id": match_from.from[0].id,
	}
}

extract_mutator_from_match(match_from, match_to) := value if {
	print("line 67 : ", match_from.to[0])
	has_key(match_from.to[0], "mutatorRef")
	value := {
		"paths": match_from.to[0].mutatorRef,
		"declaration_id": match_from.to[0].id,
	}
}

extract_to_mutate_from_match(match_from, match_to) := value if {
	print("line 76 : ", match_to.from[0])
	has_key(match_to.from[0], "mutatorRef")
	value := {
		"paths": match_to.from[0].mutatorRef,
		"declaration_id": match_to.from[0].id,
	}
}

extract_to_mutate_from_match(match_from, match_to) := value if {
	print("line 85 : ", match_to.to[0])
	has_key(match_to.to[0], "mutatorRef")
	value := {
		"paths": match_to.to[0].mutatorRef,
		"declaration_id": match_to.to[0].id,
	}
}
