package relationship_evaluation_policy

import rego.v1

evalutate if {
	
	rels_in_design_file := input.relationships

	# iterate relationships in the design file
	some rel in rels_in_design_file
	
	# relationship := rel.selectors[0].allow
	

	resultant_patches := perform_eval(input, rel)
	# mutated_design_file = json.patch(, resultantPatchesToApply)
}

# opa eval -i ${design} --data ${policy} data.relationship_evaluation_policy.evaluate
