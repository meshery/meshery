git switch -c# https://play.openpolicyagent.org/p/wFNhyGsIej
package hierarchical_policy

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
    # print(allowed_component)
    
    allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
    apply_patch(allowed_component, service, from_selectors, to_selectors)
}

apply_patch(mutator, mutated, from_selectors, to_selectors) {
    from_selectors[i].kind == mutator.type
    mutator_path := from_selectors[i].patch.mutatorRef
    mutator_value := json.filter(mutator, [mutator_path])
    
    mutator_keys := split(mutator_path, "/")
    value_to_update := object.get(mutator_value, mutator_keys, "")
    print(value_to_update)

    to_selectors[j].kind == mutated.type
    mutated_path := to_selectors[j].patch.mutatedRef
    sp := split(mutated_path, "_")
    coun := object.get(mutated, [sp[0]], "")
    print(mutated_path, coun, sp[0], split(sp[0], "/"), object.get(mutated, split(sp[0], "/"), "kk"))
    mutated_v := json.patch(mutated, [{
        "op": "add", 
        "path": "/settings/configPatches/0/patch/value", 
        "value": {"value": "test"}
    }])

    mutated_keys := split(mutated_path, "/")
    print(mutated_v, "-----")
}
    # print(keys, value_to_patch)
    # key := keys[count(keys) - 1]
    # print(object.get(mutator, keys))
    
    # print(value_to_patch["config"])
    # print()
    # value_to_update := mutator[path] 

# sm := filtered_services[k]
# print(sm, k, "\n")