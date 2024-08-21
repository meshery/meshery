package relationship_evaluation_policy

import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "sibling"

	some selector_set in relationship.selectors

	from_selectors := {kind: selectors |
		some selectors in selector_set.allow.from
		kind := selectors.kind
	}

	to_selectors := {kind: selectors |
		some selectors in selector_set.allow.to
		kind := selectors.kind
	}

	from := extract_components(design_file.components, from_selectors)
	to := extract_components(design_file.components, to_selectors)

	evaluation_results := evaluate_siblings with data.relationship as relationship
		with data.from as from
		with data.to as to
		with data.from_selectors as from_selectors
		with data.to_selectors as to_selectors
}

evaluate_siblings contains result if {
	some from_selector in data.from_selectors
	some to_selector in data.to_selectors

	filtered_from_decls := extract_components_by_type(data.from, from_selector)
	filtered_to_decls := extract_components_by_type(data.to, to_selector)

	some from_decl in filtered_from_decls
	some to_decl in filtered_to_decls

	from_decl.id != to_decl.id

	is_valid_siblings(from_decl, to_decl, from_selector, to_selector)

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

is_valid_siblings(from_declaration, to_declaration, from_selector, to_selector) if {
	match_from := from_selector.match.refs
	match_to := to_selector.match.refs

	match_results := [result |
		some i in numbers.range(0, count(match_from))
		result := is_feasible(match_from[i], match_to[i], from_declaration, to_declaration)
	]

	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]

	count(match_results) == count(valid_results)
}
