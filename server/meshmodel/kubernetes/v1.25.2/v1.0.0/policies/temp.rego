package meshmodel_policy

import rego.v1

import data.helper.extract_components_by_type
import data.helper.extract_components
import data.helper.array_path_position
import data.helper.resolve_path
import data.helper.has_key
import data.helper.is_relationship_feasible
import data.path_builder.ensureParentPathsExist

hierarchical_inventory_relationship contains results if {
	relationship := data.relationships[_]
	relationship.subType in {"Inventory", "Parent"}

	selector_set := relationship.selectors[_]

	from_selectors := {kind: selector |
		selector := selector_set.allow.from[_]
		kind := selector.kind
	}


	to_selectors := {kind: selector |
		selector := selector_set.allow.to[_]
		kind := selector.kind
	}

	# contains "selectors.from" components only, eg: WASMFilters comps only
	allowed_parent_comps := extract_components(input.services, from_selectors)

	services_map := {service.traits.meshmap.id: service |
		service := input.services[_]
	}

	# break this 3 actions into separate rules?

	# curernt set up leads to for every relationship patches, addtion  and such object will be present refer the screenshot in notes
	patches := {id: updated_comp |
		some i, j
		service := data.services_map[i]
		allowed_component := allowed_parent_comps[j]

		allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
		updated_comp := apply_patch(allowed_component, service, from_selectors, to_selectors)
		# The id is the combination of the <id of the node from which the config was patched>_<id of the node from which the config was patched>
		# eg: <configmap node id>_<deployment_node_id>
		# The client interceptor should rely on the id present in the traits only and hsould not bother to parse the above id, the above id is included to provided more context.
		id := sprintf("%s_%s", [j, i])
	} 

	## does below rule run in contenxt of loops aboe?
	

	# add the notes written inside notes, about how this logic works and the example
	addition := { result |
		print("DASDASDASD")
		mutated_selector_set := { selector |
			some set in selector_set.allow
			selector := mutated_selectors(set)
			count(selector) > 0
		}

	 	umutated_selector_set := union(mutated_selector_set)
		# print("\nUNION OF mutated_selector_set : ", umutated_selector_set)

		mutator_selector_set := { selector |
			some set in selector_set.allow
			selector := mutator_selectors(set)
			count(selector) > 0
		}

		umutator_selector_set := union(mutator_selector_set)
		# print("\nUNION OF mutator_selector_set : ", umutator_selector_set)
		
		
		# extract all the components which gets mutated 
		some mutated_selector in umutated_selector_set 
		# print("\nMutated Selector: ", mutated_selector)
		
		mutated_components := extract_components_by_type(input.services, mutated_selector)
		
		some mutated_component in mutated_components
			print("\n Mutated_component: ", mutated_component)

			mutated_values := { formatted_path: mutated_component_value |
				some ref in mutated_selector.patch.mutatedRef
				mutated_path := resolve_path(ref, mutated_component)
				formatted_path := format_json_path(mutated_path)
				print("Mutated path : ", formatted_path)
				mutated_component_value := object.get(mutated_component, formatted_path, null)
				mutated_component_value != null
				print("\nmutated_component_value : ", mutated_component_value)
			}
			

		# for each such components, find corresponding mutator component.
		# if such component doesn't exist, it indicates a scenario where,
		# a component is referring other component of different type but it doesn't exist in the design, hence give the suggestion to the client to create the component with requried configuration.
		# eg for namespace and other namespaced resource, and second example of configmap - deployment.
		some mutator_selector in umutator_selector_set
		
		result := process_comps_to_add(mutated_values, mutator_selector)

		print("\n Mutator_values : ", result)
		print("\n Mutated_values : ", mutated_values)
	}

	results := {
		"patches": patches,
		"addition": addition
	}
}

extract_mutator(mutator_components, mutator_selector) := mutator_values if {
	some mutator_component in mutator_components

	mutator_values := { formatted_path: mutator_component_value |

		some ref in mutator_selector.patch.mutatorRef

		mutator_path := resolve_path(ref, mutator_component)
		formatted_path := format_json_path(mutator_path)
		
		print("Mutator path : ", formatted_path)
		
		mutator_component_value := object.get(mutator_component, formatted_path, null)
		mutator_component_value != null
		
		print("\nmutator_component_value : ", mutator_component_value)
	}	
}

process_comps_to_add(mutated_values, mutator_selector) := result if {
	
	mutator_components := extract_components_by_type(input.services, mutator_selector)
	print("\nmutator_component: ", mutator_components)
	
	count(mutator_components) > 0
	mutator_values := extract_mutator(mutator_components, mutator_selector)
		
	count(mutator_values) != 0

	print("mutated_values: ", mutated_values, "mutator_values: ", mutator_values)
	mutated_values != mutator_values
	# filter only those whch don't match later
	result := {
		"type": mutator_selector.kind,
		"model": mutator_selector.model,
		# "path": mutator_path,
		"values": mutated_values	
	}
	print("RRRRR--------", result)
}

format_json_path(path) := [fp |
	some p in path 
	# fp := is_numeric(p)

	fp := format_path(p)
]

format_path(s) := result if {
	regex.match("^[0-9]+$", s)
	is_string(s)
	result := to_number(s)
} else := s

mutator_selectors(selector_set) := { selector |
	some selector in selector_set
	contains_mutator_selector(selector)
}

mutated_selectors(selector_set) := { selector |
	some selector in selector_set
	contains_mutated_selector(selector)
}

contains_mutator_selector(selector) if {
	selector.patch.mutatorRef 	
}

contains_mutated_selector(selector) if {
	selector.patch.mutatedRef 	
}

apply_patch(mutator, mutated, from_selectors, to_selectors) := design_diff if {
	some i, j

	is_relationship_feasible(from_selectors[i], mutator.type)

	is_relationship_feasible(to_selectors[j], mutated.type)

	mutatorObj := identifyMutator(from_selectors[i], to_selectors[j], mutator, mutated)

	mutatedObj := identifyMutated(from_selectors[i], to_selectors[j], mutator, mutated)

	patches := [patch |
		some i
		mutator_path := resolve_path(mutatorObj.path[i], mutatorObj.mutator)
		update_value := object.get(mutatorObj.mutator, mutator_path, "")
		update_value != null
		mutated_path := resolve_path(mutatedObj.path[i], mutatedObj.mutated)
		patch := {
			"op": "add",
			"path": mutated_path,
			"value": update_value,
		}
	]
	resultantPatchesToApply := ensureParentPathsExist(patches, mutatedObj.mutated)
	mutated_design = json.patch(mutatedObj.mutated, resultantPatchesToApply)
	
	# stores the index at which the final unresolved patch is stored.
	index := count(resultantPatchesToApply) - 1
	# return only the diff as result of evaluation which needs to be applied by the client
	design_diff := json.filter(mutated_design, resultantPatchesToApply[ind].path)
}

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj if {
	has_key(to_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutated,
		"path": to_selector.patch.mutatorRef,
	}
}

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj if {
	has_key(from_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutator,
		"path": from_selector.patch.mutatorRef,
	}
}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj if {
	has_key(from_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutator,
		"path": from_selector.patch.mutatedRef,
	}
}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj if {
	has_key(to_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutated,
		"path": to_selector.patch.mutatedRef,
	}
}
