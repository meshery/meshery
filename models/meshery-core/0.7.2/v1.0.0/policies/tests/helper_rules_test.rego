package helper_rules_test

import rego.v1

import data.relationship_evaluation_policy

# Test has_key function
test_has_key_true if {
	obj := {"name": "test", "value": 123}
	relationship_evaluation_policy.has_key(obj, "name")
}

test_has_key_false if {
	obj := {"name": "test"}
	not relationship_evaluation_policy.has_key(obj, "missing")
}

# Test set_to_array function
test_set_to_array if {
	result := relationship_evaluation_policy.set_to_array({1, 2, 3})
	count(result) == 3
}

# Test array_to_set function
test_array_to_set if {
	result := relationship_evaluation_policy.array_to_set([1, 2, 3, 3])
	count(result) == 3
}

# Test declaration_with_id
test_declaration_with_id_found if {
	design := {"components": [
		{"id": "comp-1", "name": "First"},
		{"id": "comp-2", "name": "Second"},
	]}
	result := relationship_evaluation_policy.declaration_with_id(design, "comp-1")
	result.name == "First"
}

# Test arr_contains function
test_arr_contains_true if {
	relationship_evaluation_policy.arr_contains(["a", "b", "_"], "_")
}

test_arr_contains_false if {
	not relationship_evaluation_policy.arr_contains(["a", "b", "c"], "_")
}

# Test array_path_position
test_array_path_position if {
	result := relationship_evaluation_policy.array_path_position(["spec", "containers", "_", "env"])
	result == 2
}

# Test match_object
test_match_object_true if {
	obj1 := {"a": 1, "b": 2}
	obj2 := {"b": 2, "a": 1}
	relationship_evaluation_policy.match_object(obj1, obj2)
}

test_match_object_false if {
	obj1 := {"a": 1, "b": 2}
	obj2 := {"a": 1, "b": 3}
	not relationship_evaluation_policy.match_object(obj1, obj2)
}

# Test format_path - string to number conversion
test_format_path_numeric if {
	result := relationship_evaluation_policy.format_path("123")
	result == 123
}

test_format_path_non_numeric if {
	result := relationship_evaluation_policy.format_path("abc")
	result == "abc"
}

# Test format_json_path
test_format_json_path if {
	result := relationship_evaluation_policy.format_json_path(["spec", "0", "name"])
	result == ["spec", 0, "name"]
}

# Test is_relationship_feasible with wildcard
test_is_relationship_feasible_wildcard if {
	selector := {"kind": "*"}
	relationship_evaluation_policy.is_relationship_feasible(selector, "Pod")
}

# Test is_relationship_feasible with exact match
test_is_relationship_feasible_exact if {
	selector := {"kind": "Pod"}
	relationship_evaluation_policy.is_relationship_feasible(selector, "Pod")
}

test_is_relationship_feasible_mismatch if {
	selector := {"kind": "Deployment"}
	not relationship_evaluation_policy.is_relationship_feasible(selector, "Pod")
}

# Test extract_values
test_extract_values if {
	component := {"configuration": {"spec": {"containers": [{"name": "web"}]}}}
	refs := [["configuration", "spec", "containers"]]
	result := relationship_evaluation_policy.extract_values(component, refs)
	count(result) == 1
}

# Test resolve_path without wildcard
test_resolve_path_no_wildcard if {
	path := ["spec", "containers"]
	mutated := {"spec": {"containers": []}}
	result := relationship_evaluation_policy.resolve_path(path, mutated)
	result == ["spec", "containers"]
}

# Test array_index_to_patch with empty array
test_array_index_to_patch_empty if {
	result := relationship_evaluation_policy.array_index_to_patch(0)
	result == "0"
}

# Test array_index_to_patch with non-empty array
test_array_index_to_patch_non_empty if {
	result := relationship_evaluation_policy.array_index_to_patch(3)
	result == "2"
}

# Test array_index_to_patch with single element array
test_array_index_to_patch_single if {
	result := relationship_evaluation_policy.array_index_to_patch(1)
	result == "0"
}

# Test array_index_to_patch with large index
test_array_index_to_patch_large if {
	result := relationship_evaluation_policy.array_index_to_patch(100)
	result == "99"
}
