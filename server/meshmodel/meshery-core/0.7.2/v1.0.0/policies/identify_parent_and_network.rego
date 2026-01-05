package relationship_evaluation_policy

import data.feasibility_evaluation_utils
import rego.v1

identify_relationship(
	design_file,
	relationship,
) := evaluation_results if {
	applicable_on_rels := [{"kind": "hierarchical", "type": "parent"}]

	{"kind": lower(relationship.kind), "type": lower(relationship.type)} in applicable_on_rels

	# Annotation edges have wildcard selectors and may create duplicate relationships.
	# Hence, do not identify annotation relationships here.
	relationship.subType != "annotation"
	relationship.metadata.isAnnotation != true

	some selector_set in relationship.selectors

	from_selectors := {kind: selectors |
		some selectors in selector_set.allow.from
		kind := selectors.kind
	}

	to_selectors := {kind: selectors |
		some selectors in selector_set.allow.to
		kind := selectors.kind
	}

	# Extract source and target components
	from := extract_components(design_file.components, from_selectors)
	to := extract_components(design_file.components, to_selectors)

	evaluation_results := evaluate_hierarchy(
		relationship,
		from,
		to,
		from_selectors,
		to_selectors,
	)
}

evaluate_hierarchy(
	relationship,
	from,
	to,
	from_selectors,
	to_selectors,
) := {result |
	some from_selector in from_selectors
	some to_selector in to_selectors

	filtered_from_decls := extract_components_by_type(from, from_selector)
	filtered_to_decls := extract_components_by_type(to, to_selector)

	some from_decl in filtered_from_decls
	some to_decl in filtered_to_decls

	from_decl.id != to_decl.id


	# Ensure the relationship is feasible
	s := feasibility_evaluation_utils.feasible_relationship_selector_between(
		from_decl,
		to_decl,
		relationship,
	)

	

	is_valid_hierarchy(
		from_decl,
		to_decl,
		from_selector,
		to_selector,
	)

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

	now := format_int(time.now_ns(), 10)

	id := uuid.rfc4122(sprintf(
		"%s%s%s%s",
		[from_decl.id, to_decl.id, relationship.id, now],
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

is_valid_hierarchy(
	from_declaration,
	to_declaration,
	from_selector,
	to_selector,
) if {
	mutator_selector := identify_mutator(
		from_selector,
		to_selector,
		from_declaration,
		to_declaration,
	)

	mutated_selector := identify_mutated(
		from_selector,
		to_selector,
		from_declaration,
		to_declaration,
	)

	match_results := [
		result |
		range := numbers.range(
			0,
			min([
				count(mutator_selector.paths),
				count(mutated_selector.paths),
			]) - 1,
		)
		some i in range
		result := is_feasible(
			mutator_selector.paths[i],
			mutated_selector.paths[i],
			mutator_selector.declaration,
			mutated_selector.declaration,
			null,
			"",
		)
	]

	valid_results := [i |
		some result in match_results
		result == true
		i := result
	]

	count(valid_results) == count(match_results)
}
