package relationship_evaluation_policy

import rego.v1

perform_eval(
	design_file,
	relationship,
) := patched_declarations if {
	applicable_on_rels := [
		{"kind": "hierarchical", "type": "parent"},
		{"kind": "edge", "type": "non-binding"},
	]

	{"kind": lower(relationship.kind), "type": lower(relationship.type)} in applicable_on_rels

	patched_declarations := [result |
		mutator_object := extract_mutator_config_from_patch(relationship.selectors[0].allow)

		# from_node_declaration_id :=
		to_mutate_object := extract_tomutate_config_from_patch(relationship.selectors[0].allow)

		result := apply_patch(design_file, mutator_object, to_mutate_object)
	]
	# perform merge of the same object after returning from this rule in the evaluation rule.
	# based on the relationship it can happen that the same object is updated at differnet location
	# hence merge all those updates.
}

# use the below rules, when dealing with relationships block in the design file.
# In the design file's relationships block, only one array item is present inside the selectors.

extract_mutator_config_from_patch(selector) := value if {
	has_key(selector.from[0].patch, "mutatorRef")
	value = {
		"paths": selector.from[0].patch.mutatorRef,
		"declaration_id": selector.from[0].id,
	}
}

extract_mutator_config_from_patch(selector) := value if {
	has_key(selector.to[0].patch, "mutatorRef")
	value = {
		"paths": selector.to[0].patch.mutatorRef,
		"declaration_id": selector.to[0].id,
	}
}

extract_tomutate_config_from_patch(selector) := value if {
	has_key(selector.from[0].patch, "mutatedRef")
	value = {
		"paths": selector.from[0].patch.mutatedRef,
		"declaration_id": selector.from[0].id,
	}
}

extract_tomutate_config_from_patch(selector) := value if {
	has_key(selector.to[0].patch, "mutatedRef")
	value = {
		"paths": selector.to[0].patch.mutatedRef,
		"declaration_id": selector.to[0].id,
	}
}
