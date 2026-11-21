package eval_rules_test

import rego.v1

import data.eval_rules

# Test same_relationship_identifier
test_same_relationship_identifier_true if {
	rel_a := {"kind": "edge", "type": "binding", "subType": "mount"}
	rel_b := {"kind": "edge", "type": "binding", "subType": "mount"}
	eval_rules.same_relationship_identifier(rel_a, rel_b)
}

test_same_relationship_identifier_different_kind if {
	rel_a := {"kind": "edge", "type": "binding", "subType": "mount"}
	rel_b := {"kind": "hierarchical", "type": "binding", "subType": "mount"}
	not eval_rules.same_relationship_identifier(rel_a, rel_b)
}

test_same_relationship_identifier_different_type if {
	rel_a := {"kind": "edge", "type": "binding", "subType": "mount"}
	rel_b := {"kind": "edge", "type": "non-binding", "subType": "mount"}
	not eval_rules.same_relationship_identifier(rel_a, rel_b)
}

# Test same_relationship_selector_clause
test_same_relationship_selector_clause_true if {
	clause_a := {"kind": "Pod", "id": "comp-1", "patch": {"mutatorRef": ["spec"]}}
	clause_b := {"kind": "Pod", "id": "comp-1", "patch": {"mutatorRef": ["spec"]}}
	eval_rules.same_relationship_selector_clause(clause_a, clause_b)
}

test_same_relationship_selector_clause_different_kind if {
	clause_a := {"kind": "Pod", "id": "comp-1", "patch": {"mutatorRef": ["spec"]}}
	clause_b := {"kind": "Deployment", "id": "comp-1", "patch": {"mutatorRef": ["spec"]}}
	not eval_rules.same_relationship_selector_clause(clause_a, clause_b)
}

# Test from_and_to_components_exist
test_from_and_to_components_exist_true if {
	design := {"components": [
		{"id": "comp-1"},
		{"id": "comp-2"},
	]}
	relationship := {"selectors": [{"allow": {
		"from": [{"id": "comp-1"}],
		"to": [{"id": "comp-2"}],
	}}]}
	eval_rules.from_and_to_components_exist(relationship, design)
}

test_from_and_to_components_exist_missing_from if {
	design := {"components": [{"id": "comp-2"}]}
	relationship := {"selectors": [{"allow": {
		"from": [{"id": "comp-1"}],
		"to": [{"id": "comp-2"}],
	}}]}
	not eval_rules.from_and_to_components_exist(relationship, design)
}

test_from_and_to_components_exist_missing_to if {
	design := {"components": [{"id": "comp-1"}]}
	relationship := {"selectors": [{"allow": {
		"from": [{"id": "comp-1"}],
		"to": [{"id": "comp-2"}],
	}}]}
	not eval_rules.from_and_to_components_exist(relationship, design)
}

# Test from_or_to_components_dont_exist
test_from_or_to_components_dont_exist_true if {
	design := {"components": [{"id": "comp-1"}]}
	relationship := {"selectors": [{"allow": {
		"from": [{"id": "comp-1"}],
		"to": [{"id": "comp-2"}],
	}}]}
	eval_rules.from_or_to_components_dont_exist(relationship, design)
}

test_from_or_to_components_dont_exist_false if {
	design := {"components": [
		{"id": "comp-1"},
		{"id": "comp-2"},
	]}
	relationship := {"selectors": [{"allow": {
		"from": [{"id": "comp-1"}],
		"to": [{"id": "comp-2"}],
	}}]}
	not eval_rules.from_or_to_components_dont_exist(relationship, design)
}

# Test match_values with different strategies
test_match_values_equal_true if {
	eval_rules.match_values("test", "test", "equal")
}

test_match_values_equal_false if {
	not eval_rules.match_values("test", "other", "equal")
}

test_match_values_equal_as_strings_number if {
	eval_rules.match_values(123, "123", "equal_as_strings")
}

test_match_values_not_null_both_non_null if {
	eval_rules.match_values("value1", "value2", "not_null")
}

test_match_values_not_null_false_first if {
	not eval_rules.match_values(null, "value", "not_null")
}

test_match_values_not_null_false_second if {
	not eval_rules.match_values("value", null, "not_null")
}

# Test cleanup_deleted_relationships_actions
test_cleanup_deleted_relationships_actions if {
	relationships := [
		{"id": "rel-1", "status": "deleted"},
		{"id": "rel-2", "status": "approved"},
	]
	result := eval_rules.cleanup_deleted_relationships_actions(relationships)
	count(result) == 1
}

test_cleanup_deleted_relationships_actions_none if {
	relationships := [
		{"id": "rel-1", "status": "approved"},
		{"id": "rel-2", "status": "pending"},
	]
	result := eval_rules.cleanup_deleted_relationships_actions(relationships)
	count(result) == 0
}

# Test identify_relationships_based_on_matching_mutator_and_mutated_fields
# Regression tests for meshery/meshery#24006: the identify rule paired components
# to selectors by kind only, ignoring the selector model, so two components of
# the same kind from different models could be wrongly identified as related.

edge_non_binding_mutator_relationship := {
	"id": "test-rel-id",
	"kind": "edge",
	"type": "non-binding",
	"subType": "network",
	"status": "pending",
	"selectors": [{
		"allow": {
			"from": [{
				"id": "sel-from-1",
				"kind": "Pod",
				"model": {"name": "kubernetes"},
				"patch": {"mutatorRef": [["configuration", "metadata", "namespace"]]},
			}],
			"to": [{
				"id": "sel-to-1",
				"kind": "Pod",
				"model": {"name": "kubernetes"},
				"patch": {"mutatedRef": [["configuration", "metadata", "namespace"]]},
			}],
		},
	}],
}

test_identify_mutator_no_match_across_models if {
	design_file := {"components": [
		{
			"id": "pod-1",
			"component": {"kind": "Pod"},
			"model": {"name": "kubernetes"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
		{
			"id": "pod-2",
			"component": {"kind": "Pod"},
			"model": {"name": "istio"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
	]}

	result := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(edge_non_binding_mutator_relationship, design_file)
	count(result) == 0
}

test_identify_mutator_matches_same_model if {
	design_file := {"components": [
		{
			"id": "pod-1",
			"component": {"kind": "Pod"},
			"model": {"name": "kubernetes"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
		{
			"id": "pod-2",
			"component": {"kind": "Pod"},
			"model": {"name": "kubernetes"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
	]}

	result := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(edge_non_binding_mutator_relationship, design_file)
	count(result) == 2
}

test_identify_mutator_matches_wildcard_selector_model if {
	relationship := json.patch(edge_non_binding_mutator_relationship, [
		{"op": "replace", "path": "/selectors/0/allow/from/0/model/name", "value": "*"},
		{"op": "replace", "path": "/selectors/0/allow/to/0/model/name", "value": "*"},
	])
	design_file := {"components": [
		{
			"id": "pod-1",
			"component": {"kind": "Pod"},
			"model": {"name": "kubernetes"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
		{
			"id": "pod-2",
			"component": {"kind": "Pod"},
			"model": {"name": "istio"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
	]}

	result := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(relationship, design_file)
	count(result) == 2
}

test_identify_mutator_matches_selector_without_model if {
	relationship := json.patch(edge_non_binding_mutator_relationship, [
		{"op": "remove", "path": "/selectors/0/allow/from/0/model"},
		{"op": "remove", "path": "/selectors/0/allow/to/0/model"},
	])
	design_file := {"components": [
		{
			"id": "pod-1",
			"component": {"kind": "Pod"},
			"model": {"name": "kubernetes"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
		{
			"id": "pod-2",
			"component": {"kind": "Pod"},
			"model": {"name": "istio"},
			"configuration": {"metadata": {"namespace": "team-a"}},
		},
	]}

	result := eval_rules.identify_relationships_based_on_matching_mutator_and_mutated_fields(relationship, design_file)
	count(result) == 2
}
