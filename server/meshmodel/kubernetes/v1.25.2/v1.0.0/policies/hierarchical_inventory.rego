package meshmodel_policy

import data.common.contains
import data.common.is_relationship_feasible
import data.common.has_key
import data.common.extract_components
import data.common.get_array_pos
import data.common.get_path
import data.path_builder.ensureParentPathsExist
import future.keywords.every
import future.keywords.in

hierarchical_inventory_relationship[updated_comps] {
	relationship := data.relationships[_]
	relationship.subType in {"Inventory", "Parent"}
		
	selector_set := relationship.selectors[_]
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

	services_map := {service.traits.meshmap.id: service |
		service := input.services[_]
	}

	updated_comps = {id: updated_comp |
		some i, j
		service := services_map[i]
		allowed_component := allowed_parent_comps[j]

		allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
		updated_comp := apply_patch(allowed_component, service, from_selectors, to_selectors)
		# The id is the combination of the <id of the node from which the config was patched>_<id of the node from which the config was patched>
		# eg: <configmap node id>_<deployment_node_id>
		# The client interceptor should rely on the id present in the traits only and hsould not bother to parse the above id, the above id is included to provided more context.
		id := sprintf("%s_%s", [j, i])
	}
}

apply_patch(mutator, mutated, from_selectors, to_selectors) := mutated_design {
	some i, j

	is_relationship_feasible(from_selectors[i], mutator.type)
	
	is_relationship_feasible(to_selectors[j], mutated.type)

	mutatorObj := identifyMutator(from_selectors[i], to_selectors[j], mutator, mutated)
	
	mutatedObj := identifyMutated(from_selectors[i], to_selectors[j], mutator, mutated)


	patches := [patch |
	some i
		mutator_path := get_path(mutatorObj.path[i], mutatorObj.mutator)
		update_value := object.get(mutatorObj.mutator, mutatorObj.path[i], "")
		update_value != null
		mutated_path := get_path(mutatedObj.path[i], mutatedObj.mutated)

		patch := {
			"op": "add",
			"path": mutated_path,
			"value": update_value,
		}
	]
	resultantPatchesToApply := ensureParentPathsExist(patches, mutatedObj.mutated)
	mutated_design = json.patch(mutatedObj.mutated, resultantPatchesToApply)
}

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj {
	has_key(to_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutated,
		"path": to_selector.patch.mutatorRef
	}

}

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj {
	has_key(from_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutator,
		"path": from_selector.patch.mutatorRef
	}

}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj {
	has_key(from_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutator,
		"path": from_selector.patch.mutatedRef
	}

}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj {
	has_key(to_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutated,
		"path": to_selector.patch.mutatedRef
	}

}