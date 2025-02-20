package core_utils

import rego.v1

#--- General datastructures and algorithm  utils

new_uuid(seed) := id if {
	now := format_int(time.now_ns(), 10)
	id := uuid.rfc4122(sprintf("%s%s", [seed, now]))
}

# is deterministic in eval
static_uuid(seed) := id if {
	id := uuid.rfc4122(sprintf("%s%s", [seed]))
}

object_get_nested(obj, path, default_value) := current_value if {
	stringfied_path := [sprintf("%v", [v]) | some v in path]
	[current_path, current_value] := walk(obj)
	stringfied_current_path := [sprintf("%v", [v]) | some v in current_path]
	stringfied_current_path == stringfied_path
} else := default_value

pop_last(arr) := array.slice(arr, 0, count(arr) - 1)

pop_first(arr) := array.slice(arr, 1, count(arr))

array_endswith(arr, item) if {
	arr[count(arr) - 1] == item
}

# truncate_set restricts a set to a maximum number of elements
#
# Note: It's legal to define the same function twice in Rego,
# but with a crucial caveat: the function definitions must have different
# arguments or conditions. This is how Rego implements function overloading.
# Rego will evaluate the conditions and use the appropriate definition based
# on the input.
#
# Args:
#   s: The input set to be limited
#   max_length: Maximum number of elements to keep
#
# Returns:
#   A set containing up to max_length elements from the original set
#
# Behavior:
#   - If the input set is smaller than or equal to max_length, returns the original set
#   - If the input set is larger than max_length, returns a set with only the first max_length elements
truncate_set(s, max_length) := result if {
	arr := [x | x := s[_]]

	count(arr) <= max_length
	result := s
}

truncate_set(s, max_length) := result if {
	arr := [x | x := s[_]]
	count(arr) > max_length
	result := {arr[i] | i < max_length}
}

#-----------

#-------- Get Component Configuration -----------

component_alias(component_id) := alias if {
	alias := input.metadata.resolvedAliases[component_id]
}

get_component_configuration(component, design) := configuration if {
	alias := component_alias(component.id)

	# print("configuration from Alias is ==>",alias)

	parent := component_declaration_by_id(design, alias.resolved_parent_id)

	configuration := object_get_nested(parent, alias.resolved_ref_field_path, null)
	# print("Configuration got from Alias " ,configuration)
}

get_component_configuration(component, design) := configuration if {
	not component_alias(component.id)
	configuration := component.configuration
	# print("Configuration direct",configuration)
}

# ------------------------------------------------

# Get the from component id ( assumes only one selector and from in declaration)
from_component_id(relationship) := component if {
	component := relationship.selectors[0].allow.from[0].id
}

# Get the to component id ( assumes only one selector and from in declaration)
to_component_id(relationship) := component if {
	component := relationship.selectors[0].allow.to[0].id
}

# Get the component declaration by id
component_declaration_by_id(design_file, id) := component if {
	some component in design_file.components
	component.id == id
}

#----------------- COMPONENT CONFIGURATION UTILITIES ------------------

# is_direct_reference checks if the reference is a direct reference or an array reference.
# if the reference is a direct reference then it should not end with _
is_direct_reference(ref) if {
	not array_endswith(ref, "_")
}

configuration_for_component_at_path(path, component, design) := result if {
	result := object_get_nested(get_component_configuration(component, design), pop_first(path), null)
}

# get_array_aware_configuration_for_component_at_path returns the configuration for a component at a given path. If the path is an array reference, it returns the configuration for each element in the array. Otherwise, it returns the configuration for the path.
get_array_aware_configuration_for_component_at_path(ref, component, design) := result if {
	# print("ref",ref)
	not is_direct_reference(ref)
	direct_ref := pop_last(ref)

	# remove nullish values
	items := [item |
		some item in object_get_nested(get_component_configuration(component, design), pop_first(direct_ref), [])
		item != null
	]

	count(items) > 0

	paths := [path |
		some index in numbers.range(0, count(items) - 1)
		path := array.concat(direct_ref, [sprintf("%d", [index])])
	]

	result := {
		"items": items,
		"paths": paths,
	}
	# print("Paths", paths)
}

# get_array_aware_configuration_for_component_at_path returns the configuration for a component at a given path. If the path is an array reference, it returns the configuration for each element in the array. Otherwise, it returns the configuration for the path.
get_array_aware_configuration_for_component_at_path(ref, component, design) := result if {
	is_direct_reference(ref)
	value := object_get_nested(get_component_configuration(component, design), pop_first(ref), null)
	value != null

	result := {
		"items": [value],
		"paths": [ref],
	}
}
