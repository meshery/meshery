package common

get_path(obj, mutated) = path {
    path = is_array(obj, mutated)        
}

is_array(arr, mutated) = path {
    contains(arr, "_")
    index := get_array_pos(arr)
    prefix_path := array.slice(arr, 0, index)
    suffix_path := array.slice(arr, index + 1, count(arr))
    value_to_patch := object.get(mutated, prefix_path, "")
    intermediate_path := array.concat(prefix_path, [count(value_to_patch) - 1])
    path = array.concat(intermediate_path, suffix_path)
} 

is_array(arr, mutated) = path {
	not contains(arr, "_")
    path = arr
}

contains(arr, elem) {
  arr[_] = elem
}

get_array_pos(arr_path) = index {
    arr_path[k] == "_"
    index = k
}

extract_components(services, selectors) = components {
    components := {component.traits.meshmap.id: component | 
        selector := selectors[_]
        service := services[_]
        selector.kind == service.type; 
        component := service
    }
}
