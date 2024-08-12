package relationship_evaluation_policy

import rego.v1

evalutate := result if {
	rels_in_design_file := input.relationships

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

	updated_relationships := [result |
		# relationships from registry
		some relationship in data.relationships
		result := identify_relationship(input, relationship)
	]

	print("updated relationships", updated_relationships)
}
