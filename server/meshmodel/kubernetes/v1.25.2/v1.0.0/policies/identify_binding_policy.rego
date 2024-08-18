package relationship_evaluation_policy

import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	lower(relationship.kind) == "edge"
	lower(relationship.type) == "binding"

	some selector_set in relationship.selectors
	from_selectors := {kind: selectors |
		some selectors in selector_set.allow.from
		kind := selectors.kind
	}

	to_selectors := {kind: selectors |
		some selectors in selector_set.allow.to

		some val in selectors.match

		# why this logic, i guess to ensure conflict doesn't occur.
		# Then also add it for from, as for now "from" might not be having conflicts but can happen,
		# also add in comments about this.

		# The relationship schema, can be used to allow more than one bindings b/w the edge,
		# but currently the implementation considers only one binding component between the edge.
		val[0].kind != "self"

		kind := concat("#", {selectors.kind, val[0].kind})
	}

	# contains "selectors.from" components only, eg: Role/ClusterRole(s) comps only
	from := extract_components(input.components, from_selectors)
	to := extract_components(input.components, to_selectors)

	# should consider model and versions (regex/operator/normal string)
	binding_comps := {type |
		some match_selector in selector_set.allow.from[_].match

		match_selector[0].kind != "self"
		type = match_selector[0].kind
	}

	# This is a set of set,
	# It contains results for a particular binding_type and each binding_type can be binded by different type of nodes.
	evaluation_results := union({result |
		some comp in binding_comps
		binding_declarations := extract_components(input.components, [{"kind": comp}])

		count(binding_declarations) > 0

		result := evaluate_bindings with data.binding_declarations as binding_declarations
			with data.relationship as relationship
			with data.from as from
			with data.to as to
			with data.from_selectors as from_selectors
			with data.to_selectors as to_selectors
	})
}

evaluate_bindings contains result if {
	some i, j, k

	from_declaration := data.from[i]
	binding_declaration := data.binding_declarations[j]

	from_declaration.id != binding_declaration.id

	selector := data.from_selectors[from_declaration.component.kind]

	is_valid_binding(from_declaration, binding_declaration, selector)

	to_declaration := data.to[k]

	to_declaration.id != binding_declaration.id

	to_selector := data.to_selectors[concat("#", {to_declaration.component.kind, binding_declaration.component.kind})]

	is_valid_binding(binding_declaration, to_declaration, to_selector)

	match_selector_for_from := json.patch(selector, [
		{
			"op": "add",
			"path": "id",
			"value": from_declaration.id,
		},
		{
			"op": "add",
			"path": "/match/from/0/id",
			"value": from_declaration.id,
		},
		{
			"op": "add",
			"path": "/match/to/0/id",
			"value": binding_declaration.id,
		},
	])

	match_selector_for_to := json.patch(to_selector, [
		{
			"op": "add",
			"path": "id",
			"value": to_declaration.id,
		},
		{
			"op": "add",
			"path": "/match/from/0/id",
			"value": binding_declaration.id,
		},
		{
			"op": "add",
			"path": "/match/to/0/id",
			"value": to_declaration.id,
		},
	])

	cloned_selectors := {"selectors": [{"allow": {
		"from": [match_selector_for_from],
		"to": [match_selector_for_to],
	}}]}

	result := object.union(data.relationship, cloned_selectors)
}

is_valid_binding(resource1, resource2, selectors) if {
	match_from := extract_mutator_path(selectors.match, resource1, resource2)
	match_to := extract_tomutate_path(selectors.match, resource1, resource2)

	match_results := [result |
		some i in numbers.range(0, count(match_from))

		result := is_feasible(match_from.paths[i], match_to.paths[i], match_from.declaration, match_to.declaration)
	]
	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]

	count(match_results) == count(valid_results)
}

# If none of the match paths ("from" and "to") doesn't contain array field in between, then it is a normal lookup.
is_feasible(from, to, resource1, resource2) if {
	from_path := resolve_path(from, resource1)
	formatted_from_path = format_json_path(from_path)

	to_path := resolve_path(to, resource2)
	formatted_to_path = format_json_path(to_path)

	val1 := object.get(resource1, formatted_from_path, "")
	val2 := object.get(resource2, formatted_to_path, null)

	val1 == val2
} else := false

extract_mutator_path(match, from_declaration, to_declaration) := value if {
	has_key(match.from[0], "mutatorRef")
	value := {
		"paths": match.from[0].mutatorRef,
		"declaration": from_declaration,
	}
}

extract_mutator_path(match, from_declaration, to_declaration) := value if {
	has_key(match.to[0], "mutatorRef")
	value := {
		"paths": match.to[0].mutatorRef,
		"declaration": to_declaration,
	}
}

extract_tomutate_path(match, from_declaration, to_declaration) := value if {
	has_key(match.from[0], "mutatedRef")
	value := {
		"paths": match.from[0].mutatedRef,
		"declaration": from_declaration,
	}
}

extract_tomutate_path(match, from_declaration, to_declaration) := value if {
	has_key(match.to[0], "mutatedRef")
	value := {
		"paths": match.to[0].mutatedRef,
		"declaration": to_declaration,
	}
}
