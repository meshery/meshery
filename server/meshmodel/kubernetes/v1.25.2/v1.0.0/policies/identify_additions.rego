package relationship_evaluation_policy

import rego.v1

identify_additions(
	design_file,
	relationship,
) := comps_to_add if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "parent"

	comps_to_add := [{result.id: result} |
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

        print("\n MUTATED SELECTOR: ", mutated_selector)

		# find all the
		# components that gets mutated
		mutated_components := extract_components_by_type(design_file.components, mutated_selector)

		some mutated_component in mutated_components

        print("\nMUTATED COMPONENT: ", mutated_component)

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
	]
}

process_comps_to_add(design_file, mutated_values, mutator_selector) := declaration_with_id if {
	mutator_components := extract_components_by_type(design_file.components, mutator_selector)

	every mutator_component in mutator_components {
		# print("\nMutator processing")
		print("type: ", mutator_component.component.kind)

		print("Mutator Component ", mutator_component, "mutatorRef: ", mutator_selector.patch.mutatorRef)

		mutator_values := extract_values(mutator_component, mutator_selector.patch.mutatorRef)

		print("\nProcess comps to add: ", mutated_values, mutator_values)
		print("match_object(mutated_values, mutator_values)", match_object(mutated_values, mutator_values))
		not match_object(mutated_values, mutator_values)
	}

	# NOTE: Currently if atleast one of the mutator comp doesn't exist with the configuration as present in the mutaed comp (eg Pod and 2 configmMap resource presetn in design but for one of the configmap configuation is not desired, but because 1 configmap conf is desired the results for additions should be emoty but right now it isn't the cae, so change the logic such that if for every mutator type conp no comp has desired conf only then add to result set)

	print("\n\nREACHED 136", mutated_values)

	intermediate_patches := [patch |
		some key, val in mutated_values
        patch := {
            "op": "add",
            "path": key,
            "value": val,
        }
    ]

    print("\nIP : ", intermediate_patches)
    component := {
        "component": {
		"kind": mutator_selector.kind
        },
		"model": mutator_selector.model		
	}

	resultant_patches_to_apply := ensureParentPathsExist(intermediate_patches, component)
	
    print("\nIC : ", component)
    
	declaration := json.patch(component, resultant_patches_to_apply)
    
	print("\nDECL : ", declaration)
    # filter only those whch don't match later

    id := uuid.rfc4122(json.marshal(declaration))
    declaration_with_id := json.patch(declaration, [{
        "op": "add",
        "path": ["id"],
        "value": id,
    }])
	print("DECL WITH ID : ", declaration_with_id)
}
