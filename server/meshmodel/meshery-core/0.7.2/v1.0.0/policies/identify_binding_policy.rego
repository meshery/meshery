package relationship_evaluation_policy

import rego.v1

# This file needs to be loaded as it contains helper functions like is_feasible

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

		# Prevent self-referencing conflicts
		val[0].kind != "self"

		# Combine kinds to form a unique selector key
		kind := concat("#", {selectors.kind, val[0].kind})
	}

	# Extract source and target components
	from := extract_components(design_file.components, from_selectors)
	to := extract_components(design_file.components, to_selectors)

	# Determine binding component kinds
	binding_comps := {type |
		some match_selector in selector_set.allow.from[_].match
		match_selector[0].kind != "self"
		type = match_selector[0].kind
	}

	# Evaluate bindings for all valid combinations
	evaluation_results := union({result |
		some comp in binding_comps
		binding_declarations := extract_components(
			design_file.components,
			[{"kind": comp}],
		)

		count(binding_declarations) > 0

		result := evaluate_bindings(
			binding_declarations,
			relationship,
			from,
			to,
			from_selectors,
			to_selectors,
		)
	})
}

# Evaluate binding feasibility
evaluate_bindings(
	binding_declarations,
	relationship,
	from,
	to,
	from_selectors,
	to_selectors,
) := {result |
	some i, from_declaration in from
	some j, binding_declaration in binding_declarations

	from_declaration.id != binding_declaration.id

	selector := from_selectors[from_declaration.component.kind]

	is_valid_binding(from_declaration, binding_declaration, selector)

	some k, to_declaration in to
	to_declaration.id != binding_declaration.id

	to_selector := to_selectors[
		concat("#", {to_declaration.component.kind, binding_declaration.component.kind})
	]

	is_valid_binding(binding_declaration, to_declaration, to_selector)

	match_selector_for_from := json.patch(selector, [
		{"op": "add", "path": "id", "value": from_declaration.id},
		{"op": "add", "path": "/match/from/0/id", "value": from_declaration.id},
		{"op": "add", "path": "/match/to/0/id", "value": binding_declaration.id},
	])

	match_selector_for_to := json.patch(to_selector, [
		{"op": "add", "path": "id", "value": to_declaration.id},
		{"op": "add", "path": "/match/from/0/id", "value": binding_declaration.id},
		{"op": "add", "path": "/match/to/0/id", "value": to_declaration.id},
	])

	now := format_int(time.now_ns(), 10)

	id := uuid.rfc4122(sprintf(
		"%s%s%s%s%s",
		[from_declaration.id, binding_declaration.id, to_declaration.id, relationship.id, now],
	))

	cloned_selectors := {
		"id": id,
		"selectors": [{
			"allow": {
				"from": [match_selector_for_from],
				"to": [match_selector_for_to],
			},
		}],
	}

	result := object.union_n([
		relationship,
		cloned_selectors,
		{"status": "approved"},
	])
}

is_valid_binding(resource1, resource2, selectors) if {
	match_from := extract_mutator_path(selectors.match, resource1, resource2)
	match_to := extract_tomutate_path(selectors.match, resource1, resource2)

	match_results := [
		result |
		some i in numbers.range(0, count(match_from.paths))
		result := is_feasible(
			match_from.paths[i],
			match_to.paths[i],
			match_from.declaration,
			match_to.declaration,
			"",
			"",
		)
	]

	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]

	count(match_results) == count(valid_results)
}

# If none of the paths contain arrays, perform direct value comparison
is_feasible(from, to, resource1, resource2, default_value1, default_value2) if {
	from_path := resolve_path(from, resource1)
	formatted_from_path := format_json_path(from_path)

	to_path := resolve_path(to, resource2)
	formatted_to_path := format_json_path(to_path)

	val1 := object.get(resource1, formatted_from_path, default_value1)
	val2 := object.get(resource2, formatted_to_path, default_value2)

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
