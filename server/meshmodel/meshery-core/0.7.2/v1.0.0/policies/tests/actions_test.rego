package actions_test

import rego.v1

import data.actions

# Test action operation constants
test_update_component_op if {
actions.update_component_op == "update_component"
}

test_update_component_configuration_op if {
actions.update_component_configuration_op == "update_component_configuration"
}

test_delete_component_op if {
actions.delete_component_op == "delete_component"
}

test_add_component_op if {
actions.add_component_op == "add_component"
}

test_update_relationship_op if {
actions.update_relationship_op == "update_relationship"
}

test_delete_relationship_op if {
actions.delete_relationship_op == "delete_relationship"
}

test_add_relationship_op if {
actions.add_relationship_op == "add_relationship"
}

# Test get_component_update_op
test_get_component_update_op_configuration if {
result := actions.get_component_update_op(["configuration", "spec", "containers"])
result == "update_component_configuration"
}

test_get_component_update_op_other if {
result := actions.get_component_update_op(["metadata", "name"])
result == "update_component"
}

# Test filter_out_deleted
test_filter_out_deleted if {
items := {
{"id": "comp-1", "name": "Component 1"},
{"id": "comp-2", "name": "Component 2"},
}
delete_actions := [
{"op": "delete_component", "value": {"id": "comp-1"}},
]
result := actions.filter_out_deleted(items, delete_actions, "delete_component")
count(result) == 1
}

test_filter_out_deleted_none if {
items := {
{"id": "comp-1", "name": "Component 1"},
{"id": "comp-2", "name": "Component 2"},
}
delete_actions := []
result := actions.filter_out_deleted(items, delete_actions, "delete_component")
count(result) == 2
}

# Test append_added_items
test_append_added_items if {
items := {{"id": "comp-1", "name": "Component 1"}}
add_actions := [{"op": "add_component", "value": {"item": {"id": "comp-2", "name": "Component 2"}}}]
result := actions.append_added_items(items, add_actions, "add_component")
count(result) == 2
}

test_append_added_items_empty if {
items := {{"id": "comp-1", "name": "Component 1"}}
add_actions := []
result := actions.append_added_items(items, add_actions, "add_component")
count(result) == 1
}

# Test apply_updates_to_item
test_apply_updates_to_item if {
original := {"id": "comp-1", "status": "pending"}
updates := [
{
"op": "update_component",
"value": {"id": "comp-1", "path": "/status", "value": "approved"},
},
]
result := actions.apply_updates_to_item(original, updates, "update_component")
result.status == "approved"
}

test_apply_updates_to_item_no_matching_updates if {
original := {"id": "comp-1", "status": "pending"}
updates := [
{
"op": "update_component",
"value": {"id": "comp-2", "path": "/status", "value": "approved"},
},
]
result := actions.apply_updates_to_item(original, updates, "update_component")
result.status == "pending"
}

# Test apply_updates_to_item with mismatched operation type
test_apply_updates_to_item_wrong_op_type if {
original := {"id": "comp-1", "status": "pending"}
updates := [
{
"op": "update_relationship",
"value": {"id": "comp-1", "path": "/status", "value": "approved"},
},
]
# Using "update_component" but the update has "update_relationship" op
result := actions.apply_updates_to_item(original, updates, "update_component")
# Should not apply update since operation types don't match
result.status == "pending"
}
