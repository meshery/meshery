package relationship_evaluation_policy

import rego.v1

rels_in_design_file := input.relationships if {
	count(input.relationships) > 0
}

evaluate := updated_design_file if {
	# iterate relationships in the design file and resolve the patches.
	resultant_patches := {patched_declaration |
		some rel in rels_in_design_file
		patched_declaration := perform_eval(input, rel)
	}

	# merge the patches made to the same declaration as part of relationships.
	# separate out same declarations by id.
	intermediate_result := {x |
		some val in resultant_patches

		some nval in val
		x := nval
	}

	ans := group_by_id(intermediate_result)

	result := {mutated |
		some val in ans
		merged := object.union_n(val)
		mutated := {
			"declaration_id": merged.id,
			"declaration": merged,
		}
	}

	design_file_with_updated_declarations := [declaration |
		some val in input.components

		declaration := filter_updated_declaration(val, result)
	]

	updated_relationships := {result |
		# relationships from registry
		some relationship in data.relationships
		result := identify_relationship(input, relationship)
	}

	updated_design_file := json.patch(input, [
		{
			"op": "replace",
			"path": "/components",
			"value": design_file_with_updated_declarations,
		},
		{
			"op": "add", # include those relationships, which do not exist or should be removed.
			"path": "/relationships",
			"value": union(updated_relationships),
		},
	])
}

filter_updated_declaration(declaration, updated_declarations) := obj.declaration if {
	some obj in updated_declarations
	obj.declaration_id == declaration.id
} else := declaration
