package relationship_diff_policy_test

import rego.v1

import data.relationship_evaluation_policy

# Test is_same_binding with matching binding declarations
test_is_same_binding_with_matching_bindings if {
	ex_from := {
		"id": "from-1",
		"kind": "Role",
		"match": {
			"from": [{"id": "role-1"}],
			"to": [{"id": "binding-1"}],
		},
	}
	from := {
		"id": "from-1",
		"kind": "Role",
		"match": {
			"from": [{"id": "role-1"}],
			"to": [{"id": "binding-1"}],
		},
	}
	ex_to := {
		"id": "to-1",
		"kind": "ServiceAccount",
		"match": {"to": [{"id": "sa-1"}]},
	}
	to := {
		"id": "to-1",
		"kind": "ServiceAccount",
		"match": {"to": [{"id": "sa-1"}]},
	}

	relationship_evaluation_policy.is_same_binding(ex_from, from, ex_to, to)
}

# Test is_same_binding with different binding declarations
test_is_same_binding_with_different_bindings if {
	ex_from := {
		"id": "from-1",
		"kind": "Role",
		"match": {
			"from": [{"id": "role-1"}],
			"to": [{"id": "binding-1"}],
		},
	}
	from := {
		"id": "from-1",
		"kind": "Role",
		"match": {
			"from": [{"id": "role-1"}],
			"to": [{"id": "binding-2"}],
		},
	}
	ex_to := {
		"id": "to-1",
		"kind": "ServiceAccount",
		"match": {"to": [{"id": "sa-1"}]},
	}
	to := {
		"id": "to-1",
		"kind": "ServiceAccount",
		"match": {"to": [{"id": "sa-1"}]},
	}

	not relationship_evaluation_policy.is_same_binding(ex_from, from, ex_to, to)
}

# Test is_same_binding without match fields (non-binding relationships)
test_is_same_binding_without_match_fields if {
	ex_from := {"id": "from-1", "kind": "Service"}
	from := {"id": "from-1", "kind": "Service"}
	ex_to := {"id": "to-1", "kind": "Deployment"}
	to := {"id": "to-1", "kind": "Deployment"}

	relationship_evaluation_policy.is_same_binding(ex_from, from, ex_to, to)
}

# Test does_relationship_exist_in_design for non-binding relationships
test_does_relationship_exist_non_binding if {
	relationships := [{
		"kind": "edge",
		"type": "non-binding",
		"subType": "network",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{"id": "svc-1", "kind": "Service"}],
			"to": [{"id": "deploy-1", "kind": "Deployment"}],
		}}],
	}]

	rel := {
		"kind": "edge",
		"type": "non-binding",
		"subType": "network",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{"id": "svc-1", "kind": "Service"}],
			"to": [{"id": "deploy-1", "kind": "Deployment"}],
		}}],
	}

	relationship_evaluation_policy.does_relationship_exist_in_design(relationships, rel)
}

# Test does_relationship_exist_in_design for binding relationships with same binding
test_does_relationship_exist_binding_same if {
	relationships := [{
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-1"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}]

	rel := {
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-1"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}

	relationship_evaluation_policy.does_relationship_exist_in_design(relationships, rel)
}

# Test does_relationship_exist_in_design for binding relationships with different binding component
test_does_relationship_not_exist_binding_different if {
	relationships := [{
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-1"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}]

	# Same from/to components but different binding component (rb-2 instead of rb-1)
	rel := {
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-2"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}

	# Should NOT exist because the binding component differs
	not relationship_evaluation_policy.does_relationship_exist_in_design(relationships, rel)
}

# Test evaluate_relationships_added includes binding rels with different binding components
test_evaluate_relationships_added_different_binding if {
	design_relationships := [{
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-1"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}]

	# A new binding relationship with a different binding component
	identified_relationships := [{
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"status": "approved",
		"model": {"name": "kubernetes", "version": ""},
		"selectors": [{"allow": {
			"from": [{
				"id": "role-1",
				"kind": "Role",
				"match": {
					"from": [{"id": "role-1"}],
					"to": [{"id": "rb-2"}],
				},
			}],
			"to": [{
				"id": "sa-1",
				"kind": "ServiceAccount",
				"match": {"to": [{"id": "sa-1"}]},
			}],
		}}],
	}]

	result := relationship_evaluation_policy.evaluate_relationships_added(design_relationships, identified_relationships)
	count(result) == 1
}

# Test is_of_same_kind
test_is_of_same_kind_true if {
	rel_a := {"kind": "edge", "type": "binding", "subType": "mount"}
	rel_b := {"kind": "Edge", "type": "Binding", "subType": "Mount"}
	relationship_evaluation_policy.is_of_same_kind(rel_a, rel_b)
}

test_is_of_same_kind_false if {
	rel_a := {"kind": "edge", "type": "binding", "subType": "mount"}
	rel_b := {"kind": "hierarchical", "type": "parent", "subType": "inventory"}
	not relationship_evaluation_policy.is_of_same_kind(rel_a, rel_b)
}

# Test does_belongs_to_same_model
test_does_belongs_to_same_model_true if {
	rel_a := {"model": {"name": "kubernetes", "version": "v1.25.0"}}
	rel_b := {"model": {"name": "kubernetes", "version": "v1.25.0"}}
	relationship_evaluation_policy.does_belongs_to_same_model(rel_a, rel_b)
}

test_does_belongs_to_same_model_different_name if {
	rel_a := {"model": {"name": "kubernetes", "version": "v1.25.0"}}
	rel_b := {"model": {"name": "istio", "version": "v1.25.0"}}
	not relationship_evaluation_policy.does_belongs_to_same_model(rel_a, rel_b)
}
