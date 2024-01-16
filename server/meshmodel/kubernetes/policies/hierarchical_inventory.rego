package meshmodel_policy

import future.keywords.every
import future.keywords.in
import data.common.extract_components
import data.common.contains
import data.common.get_array_pos
import data.common.get_path

heirarchical_inventory_relationship [updated_comps] {
    selector_set := data["inventory"].selectors[_]

    from_selectors := {kind: selectors |
		selectors := selector_set.allow.from[_]
		kind := selectors.kind
	}

    to_selectors := {kind: selectors |
		selectors := selector_set.allow.to[_]
		kind := selectors.kind
	}

    # contains "selectors.from" components only, eg: WASMFilters comps only
    allowed_parent_comps := extract_components(input.services, from_selectors)
    
    services_map := { service.traits.meshmap.id: service |
        service := input.services[_]
    }
    
    filtered_services := object.remove(services_map, object.keys(allowed_parent_comps))

    updated_comps = {id: updated_comp |
        some i, j
        service := filtered_services[i]
        allowed_component := allowed_parent_comps[j]

        allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
        updated_comp := apply_patch(allowed_component, service, from_selectors, to_selectors)
        id := updated_comp.traits.meshmap.id
    }
}

apply_patch(mutator, mutated, from_selectors, to_selectors) := mutated_design {
    from_selectors[i].kind == mutator.type
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