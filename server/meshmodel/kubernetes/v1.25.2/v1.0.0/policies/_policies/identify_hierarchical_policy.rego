package relationship_evaluation_policy

import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	lower(relationship.kind) == "hierarchical"

	selector_set := relationship.selectors[_]
	from_selectors := {kind: selectors |
		some selectors in selector_set.allow.from
		kind := selectors.kind
	}

	to_selectors := {kind: selectors |
		some selectors in selector_set.allow.to
		kind := selectors.kind
	}

	# contains "selectors.from" components only, eg: Role/ClusterRole(s) comps only
	from := extract_components(input.components, from_selectors)
	to := extract_components(input.components, to_selectors)

	results := {}
	evaluation_results := evaluate_hierarchy with data.relationship as relationship
		with data.from as from
		with data.to as to
		with data.from_selectors as from_selectors
		with data.to_selectors as to_selectors
}

evaluate_hierarchy contains result if {
	some i, j

	from_declaration := data.from[i]
	from_selector := data.from_selectors[from_declaration.component.kind]

	to_declaration := data.to[j]
	to_selector := data.to_selectors[to_declaration.component.kind]

	is_valid_hierarchy(from_declaration, to_declaration, from_selector, to_selector)

	# Th criteria for relationship is met hence add the relationship if not exists already.

	match_selector_for_from := json.patch(from_selector, [{
		"op": "add",
		"path": "/id",
		"value": from_declaration.id,
	}])

	match_selector_for_to := json.patch(to_selector, [{
		"op": "add",
		"path": "/id",
		"value": to_declaration.id,
	}])

	cloned_selectors := {"selectors": [{"allow": {
		"from": [match_selector_for_from],
		"to": [match_selector_for_to],
	}}]}

	result := object.union(data.relationship, cloned_selectors)
}

is_valid_hierarchy(from_declaration, to_declaration, from_selector, to_selector) if {
	mutator_selector := identify_mutator(from_selector, to_selector, from_declaration, to_declaration)

	mutated_selector := identify_mutated(from_selector, to_selector, from_declaration, to_declaration)

	match_results := [result |
		some i
		is_feasible(mutator_selector.paths[i], mutated_selector.paths[i], mutator_selector.declaration, mutated_selector.declaration)

		result := true
	]

	count(match_results) == count(mutator_selector.paths)
	count(match_results) == count(mutated_selector.paths)
}
