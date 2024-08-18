package relationship_evaluation_policy

import rego.v1

default rels_in_design_file := []

rels_in_design_file := input.relationships if {
	count(input.relationships) > 0
}

evaluate := updated_design_file if {
	# iterate relationships in the design file and resolve the patches.
	resultant_patches := {patched_declaration |
		some rel in rels_in_design_file

		# do not evaluate relationships which have status as "deleted".
		lower(rel.status) != "deleted"

		patched_declaration := perform_eval(input, rel)
	}

	# merge the patches made to the same declaration as part of relationships.
	# separate out same declarations by id.
	intermediate_result := {x |
		some val in resultant_patches

		some nval in val
		x := nval
	}

	# assign id for new identified rels
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

	updated_relationships := union({result |
		# relationships from registry
		some relationship in data.relationships
		result := identify_relationship(input, relationship)
	})
	
	relationships_added := evaluate_relationships_added(input, updated_relationships)
	
	combined_relationships := array.concat(input.relationships, relationships_added)
	
	updated_design_file := json.patch(input, [
		{
			"op": "replace",
			"path": "/components",
			"value": design_file_with_updated_declarations,
		},
		{
			"op": "replace",
			"path": "/relationships",
			"value": combined_relationships,
		}
	])
}

filter_updated_declaration(declaration, updated_declarations) := obj.declaration if {
	some obj in updated_declarations
	obj.declaration_id == declaration.id
} else := declaration

filter_relationship(rel, relationships) := relationship if {
	some relationship in relationships
	relationship.id == rel.id
} else := rel
