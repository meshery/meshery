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

# Test case insensitivity for alias relationship
test_is_alias_relationship_case_insensitive if {
relationship := {
"kind": "Hierarchical",
"type": "Parent",
"subType": "Alias",
}
eval.is_alias_relationship(relationship)
}
