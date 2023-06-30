<<<<<<< HEAD
<<<<<<< HEAD
package hierarchical_wallet_policy

import future.keywords.every
import future.keywords.in

=======
git switch -c# https://play.openpolicyagent.org/p/wFNhyGsIej
package hierarchical_policy
>>>>>>> 672c69639 (policy)
=======
package hierarchical_wallet_policy
>>>>>>> 2c7185097 (update policy and path)

extract_components(services, selectors) = components {
    components := {component.traits.meshmap.id: component | 
        selector := selectors[_]
        service := services[_]
        selector.kind == service.type; 
        component := service
    }
}

parent_child_relationship {
    from_selectors := data.selectors.allow.from
    to_selectors := data.selectors.allow.to

    # from components, contains WASMFilters comps only
    allowed_parent_comps := extract_components(input.services, from_selectors)
    
    # to components, contains all except WASMFilters comps
    services_map := { service.traits.meshmap.id: service |
        service := input.services[_]
    }
    
    filtered_services := object.remove(services_map, object.keys(allowed_parent_comps))
    
    some i, j
    service := filtered_services[i]
    allowed_component := allowed_parent_comps[j]
    
    allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
    apply_patch(allowed_component, service, from_selectors, to_selectors)
}

apply_patch(mutator, mutated, from_selectors, to_selectors) {
    from_selectors[i].kind == mutator.type
<<<<<<< HEAD
<<<<<<< HEAD
        mutator_paths := from_selectors[i].patch.mutatorRef
        

    to_selectors[j].kind == mutated.type
        mutated_paths := to_selectors[j].patch.mutatedRef

    patches := [ patch |  
        some i
            mutator_path := get_path(mutator_paths[i], mutator)
            update_value := object.get(mutator, mutator_path, "")
            update_value != null
            mutated_path := get_path(mutated_paths[i], mutated)
            patch := {
                "op": "add",
                "path": mutated_path,
                "value": update_value
            }
    ]
    mutated_design = json.patch(mutated, patches)
}

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
=======
    mutator_path := from_selectors[i].patch.mutatorRef
    mutator_value := json.filter(mutator, [mutator_path])
    
    mutator_keys := split(mutator_path, "/")
    value_to_update := object.get(mutator_value, mutator_keys, "")
    print(value_to_update)
=======
        mutator_path := from_selectors[i].patch.mutatorRef
        update_value := object.get(mutator, mutator_path, "")
        print(update_value)
>>>>>>> 2c7185097 (update policy and path)

    to_selectors[j].kind == mutated.type
        mutated_path := to_selectors[j].patch.mutatedRef
        index := get_array_pos(mutated_path)

        prefix_path := array.slice(mutated_path, 0, index)
        suffix_path := array.slice(mutated_path, index + 1, count(mutated_path))
        value_to_patch := object.get(mutated, prefix_path, "")

        intermediate_path := array.concat(prefix_path, [count(value_to_patch) - 1])
        final_path := array.concat(intermediate_path, suffix_path)

    mutated_v := json.patch(mutated, [{
        "op": "add", 
        "path":  final_path, 
        "value":  update_value
    }])
<<<<<<< HEAD
>>>>>>> 672c69639 (policy)

    mutated_keys := split(mutated_path, "/")
    print(mutated_v, "-----")
=======
>>>>>>> 2c7185097 (update policy and path)
}

get_array_pos(arr_path) = index {
    arr_path[k] == "_"
    index = k
}