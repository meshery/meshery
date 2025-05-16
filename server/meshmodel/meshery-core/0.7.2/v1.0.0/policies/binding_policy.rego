package relationship_evaluation_policy

import rego.v1

# The `perform_eval` function evaluates relationships of kind "edge" and type "binding".
# It mutates declarations in the `design_file` according to the selectors defined in the relationship.
# Returns a list of patched declarations resulting from the applied mutations.
perform_eval(
	design_file,
	relationship,
) := patched_declarations if {
	# Ensure the relationship is of kind "edge" and type "binding"
	lower(relationship.kind) == "edge"
	lower(relationship.type) == "binding"

	# Initialize `patched_declarations` as a list comprehension
	patched_declarations := [result |
		# Iterate over each allowed selector in the first selector set of the relationship
		# Note: There should be only one selector set and one component in `from` and `to` if the relationship exists
		some allowed_selector in relationship.selectors[0].allow

		# Extract the mutator object from the match in the allowed selector
		mutator_object := extract_mutator_from_match(allowed_selector[0].match)

		# Extract the object to be mutated from the match in the allowed selector
		to_mutate_object := extract_tomutate_from_match(allowed_selector[0].match)

		# Both `mutator_object["paths"]` and `to_mutate_object["paths"]` should be arrays of JSON paths
		# from which configuration manipulation will happen. They should be of the same length.
		# If not, it indicates a bug in the relationship definition.

		# Apply the patch to the design file using the mutator and to-mutate objects
		result := apply_patch(design_file, mutator_object, to_mutate_object)
	]
	# Note: After returning from this rule, perform a merge of the same object in the evaluation rule.
	# Based on the relationship, the same object might be updated at different locations,
	# hence merge all those updates to maintain consistency.
}

# Function to extract the mutator object from the match.
# It checks if "mutatorRef" exists in `match.from[0]` or `match.to[0]`,
# and returns an object containing "paths" and "declaration_id".
extract_mutator_from_match(match) := value if {
	# Check if "mutatorRef" is present in `match.from[0]`
	has_key(match.from[0], "mutatorRef")

	# Construct the mutator object with "paths" and "declaration_id"
	value := {
		"paths": match.from[0].mutatorRef,
		"declaration_id": match.from[0].id,
	}
}

extract_mutator_from_match(match) := value if {
	# Check if "mutatorRef" is present in `match.to[0]`
	has_key(match.to[0], "mutatorRef")

	# Construct the mutator object with "paths" and "declaration_id"
	value := {
		"paths": match.to[0].mutatorRef,
		"declaration_id": match.to[0].id,
	}
}

# Function to extract the object to be mutated from the match.
# It checks if "mutatedRef" exists in `match.from[0]` or `match.to[0]`,
# and returns an object containing "paths" and "declaration_id".
extract_tomutate_from_match(match) := value if {
	# Check if "mutatedRef" is present in `match.from[0]`
	has_key(match.from[0], "mutatedRef")

	# Construct the to-mutate object with "paths" and "declaration_id"
	value := {
		"paths": match.from[0].mutatedRef,
		"declaration_id": match.from[0].id,
	}
}

extract_tomutate_from_match(match) := value if {
	# Check if "mutatedRef" is present in `match.to[0]`
	has_key(match.to[0], "mutatedRef")

	# Construct the to-mutate object with "paths" and "declaration_id"
	value := {
		"paths": match.to[0].mutatedRef,
		"declaration_id": match.to[0].id,
	}
}
