package relationship_evaluation_policy

import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "sibling"

	# Process selectors from the relationship.
	some selector_set in relationship.selectors
	from_selectors := {kind: selectors |
		some selectors in selector_set.allow.from
		kind := selectors.kind
	}
	to_selectors := {kind: selectors |
		some selectors in selector_set.allow.to
		kind := selectors.kind
	}
	deny_selectors := object.get(selector_set, "deny", [])

	# Extract components based on selectors.
	from := extract_components(design_file.components, from_selectors)
	to := extract_components(design_file.components, to_selectors)

	# Evaluate siblings.
	evaluation_results := evaluate_siblings(relationship, from, to, from_selectors, to_selectors, deny_selectors)
}

# Evaluates sibling relationships between components.
evaluate_siblings(relationship, from, to, from_selectors, to_selectors, deny_selectors) := {result |
	some from_selector in from_selectors
	some to_selector in to_selectors
	filtered_from_decls := extract_components_by_type(from, from_selector)
	filtered_to_decls := extract_components_by_type(to, to_selector)

	some from_decl in filtered_from_decls
	some to_decl in filtered_to_decls
	from_decl.id != to_decl.id

	# Ensure the relationship is not denied.
	not is_relationship_denied(from_decl, to_decl, deny_selectors)

	# Check if the siblings are valid.
	is_valid_siblings(from_decl, to_decl, from_selector, to_selector)

	# Prepare match selectors for 'from' and 'to'.
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

	# Generate a unique id for the relationship.
	now := format_int(time.now_ns(), 10)
	id := uuid.rfc4122(sprintf("%s%s%s%s", [from_decl.id, to_decl.id, relationship.id, now]))

	# Clone selectors and prepare the result.
	cloned_selectors := {
		"id": id,
		"selectors": [{"allow": {
			"from": [match_selector_for_from],
			"to": [match_selector_for_to],
		}}],
	}

	result := object.union_n([relationship, cloned_selectors, {"status": "approved"}])
}

# Checks if two declarations are valid siblings.
is_valid_siblings(from_declaration, to_declaration, from_selector, to_selector) if {
	match_from := from_selector.match.refs
	match_to := to_selector.match.refs

	# Check feasibility for each reference.
	match_results := [result |
		some i in numbers.range(0, count(match_from))

		# if the value at the specified refs does not exist/has empty value
		# and relationship should not be allowed to be created pass different default values.
		# For some rels eg: Network, when rel is created it is possible that the values are empty at refs
		# but because user has explicitly created it, we should allow that.
		# Hence in the identify_parent_and_network rule, both default values are same.
		result := is_feasible(match_from[i], match_to[i], from_declaration, to_declaration, "", null)
	]

	# Ensure all match results are true.
	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]
	count(match_results) == count(valid_results)
}
