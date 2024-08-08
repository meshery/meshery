package relationship_evaluation_policy

import rego.v1

evalutate if {
	
	rels_in_design_file := input.relationships

	# iterate relationships in the design file
	some rel in rels_in_design_file
	# In the design file relationships block,
	# if relationship exist there should be only 1 selector set and one one component in from and to.
	relationship_from := rel.selectors[0].allow.from[0]
	relationship_to := rel.selectors[0].allow.to[0]

	resultant_patches := build_patch_object(input, relationship_from, relationship_to)
	# mutated_design_file = json.patch(, resultantPatchesToApply)
}

# opa eval -i ${design} --data ${policy} data.relationship_evaluation_policy.evaluate
