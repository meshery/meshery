package relationship_evaluation_policy

import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	applicable_on_rels := [{"hierarchical": "parent"}, {"edge": "network"}]

	{lower(relationship.kind): lower(relationship.type)} in applicable_on_rels

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
	from := extract_components(design_file.components, from_selectors)
	to := extract_components(design_file.components, to_selectors)

	evaluation_results := evaluate_hierarchy with data.relationship as relationship
		with data.from as from
		with data.to as to
		with data.from_selectors as from_selectors
		with data.to_selectors as to_selectors
}

evaluate_hierarchy contains result if {
	some from_selector in data.from_selectors
	some to_selector in data.to_selectors

	filtered_from_decls := extract_components_by_type(data.from, from_selector)

	filtered_to_decls := extract_components_by_type(data.to, to_selector)

	some from_decl in filtered_from_decls
	some to_decl in filtered_to_decls

	from_decl.id != to_decl.id
	is_valid_hierarchy(from_decl, to_decl, from_selector, to_selector)

	# The criteria for relationship is met hence add the relationship.

	match_selector_for_from := json.patch(from_selector, [{
		"op": "add",
		"path": "/id",
		"value": from_decl.id,
	}])

	match_selector_for_to := json.patch(to_selector, [{
		"op": "add",
		"path": "/id",
		"value": to_decl.id,
	}])

	cloned_selectors := {"selectors": [{"allow": {
		"from": [match_selector_for_from],
		"to": [match_selector_for_to],
	}}]}

	result := object.union_n([data.relationship, cloned_selectors, {"status": "approved"}])
}

is_valid_hierarchy(from_declaration, to_declaration, from_selector, to_selector) if {
	mutator_selector := identify_mutator(from_selector, to_selector, from_declaration, to_declaration)

	mutated_selector := identify_mutated(from_selector, to_selector, from_declaration, to_declaration)

	match_results := [result |
		some i
		result := is_feasible(
			mutator_selector.paths[i],
			mutated_selector.paths[i],
			mutator_selector.declaration,
			mutated_selector.declaration,
		)
	]
	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]

	count(valid_results) == count(match_results)
}
