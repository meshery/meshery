package relationship_evaluation_policy

import rego.v1

# Contains all the helper functions used in the policy.

set_to_array(set) := [val |
	some val in set
]

array_to_set(arr) := {val |
	some val in arr
}

# Checks if object 'x' has key 'k'.
has_key(x, k) if {
	x[k]
}

# Retrieves a declaration with a specific id from the design file.
declaration_with_id(design_file, id) := result if {
	declarations := design_file.components
	some declaration in declarations
	declaration.id == id
	result = declaration
}

# Resolves the path in an array, handling wildcards for arrays.
resolve_path(arr, mutated) := path if {
	# Check if the array contains a wildcard '_'.
	arr_contains(arr, "_")
	index := array_path_position(arr)

	# Split the path at the wildcard position.
	prefix_path := array.slice(arr, 0, index)
	suffix_path := array.slice(arr, index + 1, count(arr))

	# Get the value to patch based on the prefix path.
	value_to_patch := object.get(mutated, prefix_path, "")

	# Determine the index to patch in the array.
	index_to_patch := array_index_to_patch(count(value_to_patch))

	# Construct the intermediate path with the determined index.
	intermediate_path := array.concat(prefix_path, [index_to_patch])

	# Combine the intermediate path with the suffix path.
	path = array.concat(intermediate_path, suffix_path)
}

resolve_path(arr, mutated) := path if {
	not arr_contains(arr, "_")
	path = arr
}

# Determines the index to patch in an array field.
array_index_to_patch(no_of_elements) := index if {
	no_of_elements == 0

	# If the array is empty, start at index 0.
	index = "0"
}

array_index_to_patch(no_of_elements) := index if {
	not no_of_elements == 0

	# Use the last index in the array.
	index = format_int(no_of_elements - 1, 10)
}

# Checks if an array contains a specific key.
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

# Checks if two objects have matching values.
match_object(o1, o2) if {
	o1_values := {val |
		some val in o1
	}
	o2_values := {val |
		some val in o2
	}
	o1_values == o2_values
} else := false

# Formats a JSON path, converting numeric strings to numbers.
format_json_path(path) := [fp |
	some p in path

	# fp := is_numeric(p)

	fp := format_path(p)
]

# Converts a string to a number if it represents a numeric value.
format_path(s) := result if {
	regex.match(`^[0-9]+$`, s)
	is_string(s)
	result := to_number(s)
} else := s

# Groups objects by their declaration id.
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

# Extracts components from declarations based on selectors.
extract_components(declarations, selectors) := {declaration.id: declaration |
	selector := selectors[_]
	declaration := declarations[_]
	is_relationship_feasible(selector, declaration.component.kind)
	component := declaration
}

# Extracts components of a specific type from declarations.
extract_components_by_type(declarations, selector) := {result |
	some declaration in declarations
	is_relationship_feasible(selector, declaration.component.kind)
	result := declaration
}

# TODO: Add checks for
# 1. when operators/regex are used in the version fields
# 2. deny selctor

is_relationship_feasible_from(fromComponent, relationship) := from if {
	some selector in relationship.selectors
	some from in selector.allow.from
	is_relationship_feasible(from, fromComponent.component.kind)
}

is_relationship_feasible(selector, comp_type) if {
	selector.kind == "*"
}

is_relationship_feasible(selector, comp_type) if {
	selector.kind == comp_type
}

# Extracts values from a component based on reference paths.
extract_values(component, refs) := [component_value |
	some ref in refs
	path := resolve_path(ref, component)
	formatted_path := format_json_path(path)
	component_value := object.get(component, formatted_path, null)
	component_value != null
]
