package relationship_evaluation_policy

import rego.v1

evalutate if {
	rels_in_design_file := input.relationships

	# iterate relationships in the design file
	some rel in rels_in_design_file

	resultant_patches := perform_eval(input, rel)
}
