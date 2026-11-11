package hierarchical_policy_test

import rego.v1

import data.eval

# Test hierarchical_parent_child_policy_identifier
test_hierarchical_parent_child_policy_identifier if {
	eval.hierarchical_parent_child_policy_identifier == "hierarchical_parent_child"
}

# Test relationship_is_implicated_by_policy for hierarchical parent/child
test_relationship_is_implicated_hierarchical_parent_child if {
	relationship := {
		"kind": "hierarchical",
		"type": "parent",
		"subType": "inventory",
	}
	eval.relationship_is_implicated_by_policy(relationship, "hierarchical_parent_child")
}

test_relationship_not_implicated_wrong_kind if {
	relationship := {
		"kind": "edge",
		"type": "parent",
		"subType": "inventory",
	}
	not eval.relationship_is_implicated_by_policy(relationship, "hierarchical_parent_child")
}

test_relationship_not_implicated_wrong_type if {
	relationship := {
		"kind": "hierarchical",
		"type": "child",
		"subType": "inventory",
	}
	not eval.relationship_is_implicated_by_policy(relationship, "hierarchical_parent_child")
}

test_relationship_not_implicated_wrong_subtype if {
	relationship := {
		"kind": "hierarchical",
		"type": "parent",
		"subType": "alias",
	}
	not eval.relationship_is_implicated_by_policy(relationship, "hierarchical_parent_child")
}

# Test alias relationship identification
test_alias_policy_identifier if {
	eval.alias_policy_identifier == "alias_relationships_policy"
}

test_is_alias_relationship_true if {
	relationship := {
		"kind": "hierarchical",
		"type": "parent",
		"subType": "alias",
	}
	eval.is_alias_relationship(relationship)
}

test_is_alias_relationship_false_wrong_kind if {
	relationship := {
		"kind": "edge",
		"type": "parent",
		"subType": "alias",
	}
	not eval.is_alias_relationship(relationship)
}

test_is_alias_relationship_false_wrong_type if {
	relationship := {
		"kind": "hierarchical",
		"type": "child",
		"subType": "alias",
	}
	not eval.is_alias_relationship(relationship)
}

test_is_alias_relationship_false_wrong_subtype if {
	relationship := {
		"kind": "hierarchical",
		"type": "parent",
		"subType": "inventory",
	}
	not eval.is_alias_relationship(relationship)
}

# Test relationship_is_implicated_by_policy for alias
test_relationship_is_implicated_alias if {
	relationship := {
		"kind": "hierarchical",
		"type": "parent",
		"subType": "alias",
	}
	eval.relationship_is_implicated_by_policy(relationship, "alias_relationships_policy")
}

# Test case insensitivity for alias relationship - title case
test_is_alias_relationship_case_insensitive_title if {
	relationship := {
		"kind": "Hierarchical",
		"type": "Parent",
		"subType": "Alias",
	}
	eval.is_alias_relationship(relationship)
}

# Test case insensitivity for alias relationship - mixed case
test_is_alias_relationship_case_insensitive_mixed if {
	relationship := {
		"kind": "hIeRaRcHiCaL",
		"type": "pArEnT",
		"subType": "aLiAs",
	}
	eval.is_alias_relationship(relationship)
}

# Test case insensitivity for alias relationship - all uppercase
test_is_alias_relationship_case_insensitive_upper if {
	relationship := {
		"kind": "HIERARCHICAL",
		"type": "PARENT",
		"subType": "ALIAS",
	}
	eval.is_alias_relationship(relationship)
}

base_alias_relationship := {
	"id": "rel-alias-1",
	"kind": "hierarchical",
	"type": "parent",
	"subType": "alias",
	"status": "identified",
	"selectors": [{
		"allow": {
			"from": [{
				"id": "alias-container-1",
				"kind": "Container",
				"model": {"name": "meshery-core", "registrant": {"kind": "meshery"}},
				"patch": {"mutatorRef": [["configuration", "spec", "template", "spec", "containers", "0"]]},
			}],
			"to": [{
				"id": "deploy-1",
				"kind": "Deployment",
				"model": {"name": "kubernetes", "registrant": {"kind": "github"}},
			}],
		},
	}],
}

alias_relationship_with_metadata(metadata) := object.union(base_alias_relationship, {"metadata": metadata})

test_alias_side_effects_is_annotation_false if {
	rel := alias_relationship_with_metadata({"isAnnotation": false})
	side_effects := eval.relationship_side_effects(rel, {"components": []}, "alias_relationships_policy")
	count(side_effects) == 1
	some action in side_effects
	action.op == "add_component"
	action.value.item.component.kind == "Container"
	action.value.item.displayName == "containers.0"
	action.value.item.metadata.isAnnotation == false
}

test_alias_side_effects_is_annotation_true if {
	rel := alias_relationship_with_metadata({"isAnnotation": true})
	side_effects := eval.relationship_side_effects(rel, {"components": []}, "alias_relationships_policy")
	count(side_effects) == 1
	some action in side_effects
	action.op == "add_component"
	action.value.item.component.kind == "Container"
	action.value.item.metadata.isAnnotation == true
}

test_alias_side_effects_is_annotation_field_missing if {
	rel := alias_relationship_with_metadata({"description": "metadata without isAnnotation"})
	side_effects := eval.relationship_side_effects(rel, {"components": []}, "alias_relationships_policy")
	count(side_effects) == 1
	some action in side_effects
	action.op == "add_component"
	action.value.item.metadata.isAnnotation == false
}

test_alias_side_effects_missing_metadata if {
	side_effects := eval.relationship_side_effects(base_alias_relationship, {"components": []}, "alias_relationships_policy")
	count(side_effects) == 1
	some action in side_effects
	action.op == "add_component"
	action.value.item.metadata.isAnnotation == false
}
