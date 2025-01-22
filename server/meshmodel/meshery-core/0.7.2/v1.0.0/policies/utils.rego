package relationship_evaluation_policy
import rego.v1

#--- General datastructures and algorithm  utils


new_uuid(seed) := id if {
	now := format_int(time.now_ns(), 10)
	id := uuid.rfc4122(sprintf("%s%s", [seed, now]))
}

object_get_nested(obj, path, default_value) := current_value if {
	stringfied_path := [sprintf("%v", [v]) | some v in path]
	[current_path, current_value] := walk(obj)
	stringfied_current_path := [sprintf("%v", [v]) | some v in current_path]
	stringfied_current_path == stringfied_path
} else := default_value

pop_last(arr) := array.slice(arr, 0, count(arr) - 1)
pop_first(arr) := array.slice(arr,1,count(arr))

array_endswith(arr, item) if {
	arr[count(arr) - 1] == item
}

#----------- 


#-------- Get Component Configuration -----------

get_component_configuration(component,design) := configuration if {
	alias := component_alias(component.id)
	# print("configuration from Alias is ==>",alias)

	parent := component_declaration_by_id(design,alias.resolved_parent_id)
	
	configuration := object_get_nested(parent,alias.resolved_ref_field_path,null)
	# print("Configuration got from Alias " ,configuration)
}

get_component_configuration(component,design) := configuration if {
   not component_alias(component.id)
   configuration := component.configuration

   # print("Configuration direct",configuration)
}

# ------------------------------------------------

# Get the from component id
from_component_id(relationship) := component if {
	some selector in relationship.selectors
	some from in selector.allow.from
	component := from.id
}

# Get the to component id
to_component_id(relationship) := component if {
	some selector in relationship.selectors
	some to in selector.allow.to
	component := to.id
}

# Get the component declaration by id
component_declaration_by_id(design_file, id) := component if {
	some component in design_file.components
	component.id == id
}


#----------------- Component configuration utils------------------

# check if the reference is a direct reference or an array reference
# if the reference is a direct reference then it should not end with _
is_direct_reference(ref) if {
	not array_endswith(ref, "_")
} 
get_array_aware_configuration_for_component_at_path(ref, component,design) := result if {
    print("ref",ref)
	not is_direct_reference(ref)
	direct_ref := pop_last(ref)

	# remove nullish values
	items := [item |
		some item in object_get_nested(get_component_configuration(component,design), pop_first(direct_ref), [])
		item != null
	]

	count(items) > 0

	paths := [path |
		some index in numbers.range(0, count(items) - 1)
		path := array.concat(direct_ref, [sprintf("%d", [index])])
	]

	result := {
		"items": items ,
		"paths": paths
	}
	#print("Paths", paths)
}


get_array_aware_configuration_for_component_at_path(ref, component,design) := result if {
	is_direct_reference(ref)
	value := object_get_nested(get_component_configuration(component,design), pop_first(ref), null)
	value != null

	result := {
		"items": [value],
		"paths": [ref]
	}
}

