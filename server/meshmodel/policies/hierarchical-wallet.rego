package hierarchical_wallet_policy

import future.keywords.every

extract_components(services, selectors) = components {
    components := {component.traits.meshmap.id: component | 
        selector := selectors[_]
        service := services[_]
        selector.kind == service.type; 
        component := service
    }
}

parent_child_relationship = updated_design {
    from_selectors := data.selectors.allow.from
    to_selectors := data.selectors.allow.to

    # contains "selectors.from" components only, eg: WASMFilters comps only
    allowed_parent_comps := extract_components(input.services, from_selectors)
    
    services_map := { service.traits.meshmap.id: service |
        service := input.services[_]
    }
    
    # contains all except "selector.from" components, eg: all except WASMFilters comps
    filtered_services := object.remove(services_map, object.keys(allowed_parent_comps))

    updated_comps := {id: updated_comp |
        some i, j
        service := filtered_services[i]
        allowed_component := allowed_parent_comps[j]

        allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
        updated_comp := apply_patch(allowed_component, service, from_selectors, to_selectors)
        id := updated_comp.traits.meshmap.id
    }

    remaining_comps := object.remove(services_map, object.keys(updated_comps))
    
    updated_design = object.union_n([remaining_comps, updated_comps])
}

# TODO: break into common policy funcs
apply_patch(mutator, mutated, from_selectors, to_selectors) := mutated_design {
    from_selectors[i].kind == mutator.type
        mutator_path := from_selectors[i].patch.mutatorRef
        update_value := object.get(mutator, mutator_path, "")

    to_selectors[j].kind == mutated.type
        mutated_path := to_selectors[j].patch.mutatedRef
        index := get_array_pos(mutated_path)

        prefix_path := array.slice(mutated_path, 0, index)
        suffix_path := array.slice(mutated_path, index + 1, count(mutated_path))
        value_to_patch := object.get(mutated, prefix_path, "")

        intermediate_path := array.concat(prefix_path, [count(value_to_patch) - 1])
        final_path := array.concat(intermediate_path, suffix_path)

    mutated_design = json.patch(mutated, [{
        "op": "add", 
        "path":  final_path, 
        "value":  update_value
    }])
}

get_array_pos(arr_path) = index {
    arr_path[k] == "_"
    index = k
}
