package evaluation_pipeline_test

import rego.v1

import data.actions

# =============================================================================
# Evaluation Pipeline State Transition Tests
# =============================================================================
# These tests verify that each step in the evaluation pipeline (evaluation.rego)
# correctly chains its input and actions:
#
#   Step 1: design_file_with_validated_rels = apply(design_file_to_evaluate, validation_actions)
#   Step 2: design_file_with_identified_rels = apply(design_file_with_validated_rels, new_identified_rels_actions)
#
# A previous copy-paste bug caused step 2 to apply validation_actions to the
# original design instead of applying identified relationship actions to the
# validated design. These tests catch that class of error by asserting the
# intermediate state after each step.
# =============================================================================

# -- Fixtures --

base_design := {
	"components": {
		{"id": "comp-1", "name": "Service A", "status": "active"},
		{"id": "comp-2", "name": "Service B", "status": "active"},
	},
	"relationships": {
		{
			"id": "rel-existing",
			"kind": "hierarchical",
			"type": "parent",
			"subType": "inventory",
			"status": "approved",
		},
		{
			"id": "rel-invalid",
			"kind": "hierarchical",
			"type": "parent",
			"subType": "inventory",
			"status": "approved",
		},
	},
}

# Step 1 actions: validation marks an invalid relationship as deleted
validation_actions_fixture := {
	{
		"op": actions.update_relationship_op,
		"value": {
			"id": "rel-invalid",
			"path": "/status",
			"value": "deleted",
		},
	},
}

# Step 2 actions: identification adds a newly discovered relationship
identified_rels_actions_fixture := {
	{
		"op": actions.add_relationship_op,
		"value": {"item": {
			"id": "rel-new-identified",
			"kind": "edge",
			"type": "network",
			"subType": "network-policy",
			"status": "identified",
		}},
	},
}

# -- Step 1 Tests: Validation --

test_step1_validation_marks_invalid_relationship_deleted if {
	result := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	# The invalid relationship should now have status "deleted"
	some rel in result.relationships
	rel.id == "rel-invalid"
	rel.status == "deleted"
}

test_step1_validation_preserves_valid_relationship if {
	result := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	some rel in result.relationships
	rel.id == "rel-existing"
	rel.status == "approved"
}

test_step1_validation_preserves_components if {
	result := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)
	count(result.components) == count(base_design.components)
}

test_step1_validation_preserves_relationship_count if {
	result := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)
	count(result.relationships) == count(base_design.relationships)
}

# -- Step 2 Tests: Identification uses validated design, not original --

test_step2_identification_adds_new_relationship if {
	# Step 1
	validated := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	# Step 2: apply identified actions to validated design (the fix)
	result := actions.apply_all_actions_to_design(validated, identified_rels_actions_fixture)

	# New relationship must be present
	rel_ids := {rel.id | some rel in result.relationships}
	"rel-new-identified" in rel_ids
}

test_step2_identification_retains_validation_changes if {
	# Step 1
	validated := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	# Step 2
	result := actions.apply_all_actions_to_design(validated, identified_rels_actions_fixture)

	# The validation change (status -> deleted) must still be present
	some rel in result.relationships
	rel.id == "rel-invalid"
	rel.status == "deleted"
}

test_step2_identification_relationship_count if {
	# Step 1
	validated := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	# Step 2
	result := actions.apply_all_actions_to_design(validated, identified_rels_actions_fixture)

	# Original 2 + 1 newly identified = 3
	count(result.relationships) == 3
}

# -- Regression: Applying the wrong actions to the wrong base --

test_regression_wrong_base_loses_identified_rels if {
	# Simulates the bug: applying identified actions to the ORIGINAL design
	# instead of the validated design. The new relationship should still appear
	# (it does, but validation changes are lost).
	result := actions.apply_all_actions_to_design(base_design, identified_rels_actions_fixture)

	# New rel is added (this part works even with the bug)
	rel_ids := {rel.id | some rel in result.relationships}
	"rel-new-identified" in rel_ids

	# But validation change is NOT applied (rel-invalid still has "approved")
	some rel in result.relationships
	rel.id == "rel-invalid"
	rel.status == "approved"
}

test_regression_wrong_actions_skips_identification if {
	# Simulates the other half of the bug: applying validation_actions
	# to the validated design a second time instead of identified actions.
	validated := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)
	result := actions.apply_all_actions_to_design(validated, validation_actions_fixture)

	# No new relationship is added
	rel_ids := {rel.id | some rel in result.relationships}
	not "rel-new-identified" in rel_ids
}

# -- Full pipeline correctness --

test_full_pipeline_correct_chaining if {
	# Step 1: validate
	validated := actions.apply_all_actions_to_design(base_design, validation_actions_fixture)

	# Step 2: identify (uses validated output + identification actions)
	identified := actions.apply_all_actions_to_design(validated, identified_rels_actions_fixture)

	# All original components preserved
	count(identified.components) == 2

	# 3 relationships: existing (approved), invalid (deleted), new (identified)
	count(identified.relationships) == 3

	# Verify each relationship's state
	rel_map := {rel.id: rel | some rel in identified.relationships}
	rel_map["rel-existing"].status == "approved"
	rel_map["rel-invalid"].status == "deleted"
	rel_map["rel-new-identified"].status == "identified"
}
