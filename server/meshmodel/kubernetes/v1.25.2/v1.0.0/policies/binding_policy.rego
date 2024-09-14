package relationship_evaluation_policy

import rego.v1

# performs the evaluation to mutate the declaration according to the
# selectors defined inside the relationship definition.
perform_eval(
	design_file,
	relationship,
) := patched_declarations if {
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

		# both mutator_value and mutated_value should be array of paths from which config manipulation will happen.
		# Both should be of same length.
		# In case it is not it is a bug in relationship_definition.
		result := apply_patch(design_file, mutator_object, to_mutate_object)
	]
	# perform merge of the same object after returning from this rule in the evaluation rule.
	# based on the relationship it can happen that the same object is updated at differnet location
	# hence merge all those updates.
}

extract_mutator_from_match(match) := value if {
	has_key(match.from[0], "mutatorRef")
	value := {
		"paths": match.from[0].mutatorRef,
		"declaration_id": match.from[0].id,
	}
}

extract_mutator_from_match(match) := value if {
	has_key(match.to[0], "mutatorRef")
	value := {
		"paths": match.to[0].mutatorRef,
		"declaration_id": match.to[0].id,
	}
}

extract_tomutate_from_match(match) := value if {
	has_key(match.from[0], "mutatedRef")
	value := {
		"paths": match.from[0].mutatedRef,
		"declaration_id": match.from[0].id,
	}
}

extract_tomutate_from_match(match) := value if {
	has_key(match.to[0], "mutatedRef")
	value := {
		"paths": match.to[0].mutatedRef,
		"declaration_id": match.to[0].id,
	}
}
