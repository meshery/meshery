package relationship_evaluation_policy

import rego.v1

has_key(x, k) if {
	x[k]
}

declaration_with_id(design_file, id) := result if {
	declarations := design_file.components
	some declaration in declarations

	declaration.id == id
	result = declaration
}

resolve_path(arr, mutated) := path if {
	arr_contains(arr, "_")
	index := array_path_position(arr)

	prefix_path := array.slice(arr, 0, index)
	suffix_path := array.slice(arr, index + 1, count(arr))

	value_to_patch := object.get(mutated, prefix_path, "")

	index_to_patch := array_index_to_patch(count(value_to_patch))

	intermediate_path := array.concat(prefix_path, [index_to_patch])

	path = array.concat(intermediate_path, suffix_path)
}

resolve_path(arr, mutated) := path if {
	not arr_contains(arr, "_")
	path = arr
}

array_index_to_patch(no_of_elements) := index if {
	no_of_elements == 0

	# 0 based array indexing is followed
	index = "0"
}

array_index_to_patch(no_of_elements) := index if {
	not no_of_elements == 0

	# 0 based array indexing is followed
	index = format_int(no_of_elements - 1, 10)
}

arr_contains(arr, key) if {
	some element in arr
	key == element
} else := false

# Returns the index at which the mutated or mutator paths indicates the presence of an object of type: "array".
# eg: [spec, containers, _, env]:
#  It indicates that the type of "containers" field inside the configuration is of type "array".
# Hence try to append the configuration at the last available index.
array_path_position(arr_path) := index if {
	some index
	arr_path[index] == "_"
}

match_object(o1, o2) if {
	o1_values := {val |
		some val in o1
	}

	o2_values := {val |
		some val in o2
	}

	o1_values == o2_values
} else := false

format_json_path(path) := [fp |
	some p in path

	# fp := is_numeric(p)

	fp := format_path(p)
]

format_path(s) := result if {
	regex.match(`^[0-9]+$`, s)
	is_string(s)
	result := to_number(s)
} else := s

group_by_id(objects) := {obj |
	some val in objects
	grouped_objects := [p |
		some o in objects
		o.declaration_id == val.declaration_id
		some p in o.patches
	]

	obj := {
		"declaration_id": val.declaration_id,
		"declaration": val.declaration,
		"patches": grouped_objects,
	}
}

extract_components(declarations, selectors) := {declaration.id: declaration |
	selector := selectors[_]
	declaration := declarations[_]
	is_relationship_feasible(selector, declaration.component.kind)
	component := declaration
}

extract_components_by_type(declarations, selector) := {result |
	some declaration in declarations

	is_relationship_feasible(selector, declaration.component.kind)
	result := declaration
}

# TODO: Add checks for
# 1. when operators/regex are used in the version fields
# 2. deny selctor

is_relationship_feasible(selector, comp_type) if {
	selector.kind == "*"
}

is_relationship_feasible(selector, comp_type) if {
	selector.kind == comp_type
}
