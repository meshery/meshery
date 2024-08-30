package relationship_evaluation_policy

import rego.v1

identify_additions(
	design_file,
	relationship,
) := unique_comps if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "parent"

	comps_to_add := {result |
		some selector_set in relationship.selectors

		mutated_selector_set := union({selector |
			some set in selector_set.allow
			selector := mutated_selectors(set)
			count(selector) > 0
		})

		mutator_selector_set := union({selector |
			some set in selector_set.allow
			selector := mutator_selectors(set)
			count(selector) > 0
		})

		some mutated_selector in mutated_selector_set

		# find all the
		# components that gets mutated
		mutated_components := extract_components_by_type(design_file.components, mutated_selector)

		some mutated_component in mutated_components

		mutated_values := extract_values(mutated_component, mutated_selector.patch.mutatedRef)

		# For all paths specifed in the ref, ensure the values exists.
		count(mutated_values) == count(mutated_selector.patch.mutatedRef)

		# for each such components, find corresponding mutator component.
		# if such component doesn't exist, it indicates a scenario where,
		# a component is referring other component of different type but it doesn't exist in the design,
		# hence give the suggestion to the client to create the component with requried configuration.
		# eg for namespace and other namespaced resource, and second example of configmap - deployment.
		some mutator_selector in mutator_selector_set

		result := process_comps_to_add(design_file, mutated_values, mutator_selector)
	}

	unique_comp_ids := {result.id |
		some key, result in comps_to_add
	}

	unique_comps := {result |
		some id in unique_comp_ids
		component := find_in_comp_array(comps_to_add, id)
		result := component
	}
}

find_in_comp_array(comps, id) := result if {
	some key, comp in comps
	comp.id == id
	result := comp
}

process_comps_to_add(
	design_file,
	mutated_values, mutator_selector,
) := declaration_with_id if {
	mutator_components := extract_components_by_type(design_file.components, mutator_selector)
	every mutator_component in mutator_components {
		mutator_values := extract_values(mutator_component, mutator_selector.patch.mutatorRef)
		not match_object(mutated_values, mutator_values)
	}

	# location/refs to the configuration for the newly created component.
	mutator_paths := mutator_selector.patch.mutatorRef

	# values that needs to be assigned at the location/refs for the newly created component.
	values := [val |
		some val in mutated_values
	]

	# set the values that needs to be assigned in the patch object
	intermediate_patches := [patch |
		some i in numbers.range(0, count(mutator_paths) - 1)
		patch := {
			"op": "add",
			"path": mutator_paths[i],
			"value": values[i],
		}
	]

	component := {
		"component": {"kind": mutator_selector.kind},
		"model": mutator_selector.model,
	}

	resultant_patches_to_apply := ensure_parent_paths_exist(intermediate_patches, component)

	declaration := json.patch(component, resultant_patches_to_apply)

	id := uuid.rfc4122(json.marshal(declaration))
	declaration_with_id := json.patch(declaration, [{
		"op": "add",
		"path": ["id"],
		"value": id,
	}])
}
