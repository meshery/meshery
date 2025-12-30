package core_utils_test

import rego.v1

import data.core_utils

# Test set_to_array function
test_set_to_array if {
result := core_utils.set_to_array({1, 2, 3})
count(result) == 3
}

test_set_to_array_empty if {
result := core_utils.set_to_array(set())
count(result) == 0
}

# Test array_to_set function
test_array_to_set if {
result := core_utils.array_to_set([1, 2, 3, 2])
count(result) == 3
}

test_array_to_set_empty if {
result := core_utils.array_to_set([])
count(result) == 0
}

# Test pop_last function
test_pop_last if {
result := core_utils.pop_last([1, 2, 3])
result == [1, 2]
}

test_pop_last_single_element if {
result := core_utils.pop_last([1])
result == []
}

# Test pop_first function
test_pop_first if {
result := core_utils.pop_first([1, 2, 3])
result == [2, 3]
}

test_pop_first_single_element if {
result := core_utils.pop_first([1])
result == []
}

# Test array_endswith function
test_array_endswith_true if {
core_utils.array_endswith([1, 2, 3], 3)
}

test_array_endswith_false if {
not core_utils.array_endswith([1, 2, 3], 2)
}

# Test coalesce function
test_coalesce_non_null if {
result := core_utils.coalesce("value", "default")
result == "value"
}

test_coalesce_null if {
result := core_utils.coalesce(null, "default")
result == "default"
}

# Test truncate_set smaller than max
test_truncate_set_smaller if {
result := core_utils.truncate_set({1, 2}, 5)
count(result) == 2
}

# Test truncate_set larger than max
test_truncate_set_larger if {
result := core_utils.truncate_set({1, 2, 3, 4, 5}, 2)
count(result) == 2
}

# Test normalize_path with string
test_normalize_path_string if {
result := core_utils.normalize_path("/spec/containers")
result == "/spec/containers"
}

# Test normalize_path with array
test_normalize_path_array if {
result := core_utils.normalize_path(["spec", "containers"])
result == "/spec/containers"
}

# Test object_get_nested
test_object_get_nested_found if {
obj := {"spec": {"containers": [{"name": "test"}]}}
result := core_utils.object_get_nested(obj, ["spec", "containers"], null)
result == [{"name": "test"}]
}

test_object_get_nested_not_found if {
obj := {"spec": {"containers": []}}
result := core_utils.object_get_nested(obj, ["spec", "missing"], "default")
result == "default"
}

# Test component_declaration_by_id
test_component_declaration_by_id_found if {
design := {"components": [
{"id": "comp-1", "name": "Component 1"},
{"id": "comp-2", "name": "Component 2"},
]}
result := core_utils.component_declaration_by_id(design, "comp-1")
result.name == "Component 1"
}

test_component_declaration_by_id_not_found if {
	design := {"components": [{"id": "comp-1", "name": "Component 1"}]}
	# The function returns undefined (not null) when not found
	not core_utils.component_declaration_by_id(design, "comp-999")
}

# Test component_declaration_by_id with empty components array
test_component_declaration_by_id_empty_components if {
	design := {"components": []}
	not core_utils.component_declaration_by_id(design, "comp-1")
}

# Test from_component_id
test_from_component_id if {
relationship := {
"selectors": [{
"allow": {
"from": [{"id": "from-comp-id"}],
"to": [{"id": "to-comp-id"}],
},
}],
}
result := core_utils.from_component_id(relationship)
result == "from-comp-id"
}

# Test to_component_id
test_to_component_id if {
relationship := {
"selectors": [{
"allow": {
"from": [{"id": "from-comp-id"}],
"to": [{"id": "to-comp-id"}],
},
}],
}
result := core_utils.to_component_id(relationship)
result == "to-comp-id"
}

# Test is_direct_reference
test_is_direct_reference_true if {
core_utils.is_direct_reference(["spec", "containers", "0"])
}

test_is_direct_reference_false if {
not core_utils.is_direct_reference(["spec", "containers", "_"])
}
