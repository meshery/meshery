package meshmodel_policy

import rego.v1

import data.helper.extract_components
import data.helper.extract_components_by_type
import data.helper.format_json_path
import data.helper.is_relationship_feasible
import data.helper.match_object
import data.helper.resolve_path
import data.patch_helper.identifyMutated
import data.patch_helper.identifyMutator
import data.patch_helper.mutated_selectors
import data.patch_helper.mutator_selectors
import data.path_builder.ensureParentPathsExist

hierarchical_inventory_relationship contains results if {
	services_map := {service.traits.meshmap.id: service |
		service := input.services[_]
	}

	# curernt set up leads to for every relationship patches, addtion  and such object will be present refer the screenshot in notes
	patches := identify_patches with input.services as input.services
		with data.services_map as services_map
		with data.relationships as data.relationships

	# add the notes written inside notes, about how this logic works and the example
	additions := identify_additions with input.services as input.services
		with data.services_map as services_map
		with data.relationships as data.relationships

	results := {
		"patches": {},
		"additions": additions
	}
}

identify_patches := result if {
	result := [updated_comp |
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

		allowed_parent_comps := extract_components(input.services, from_selectors)

		some i, j
		service := data.services_map[i]
		allowed_component := allowed_parent_comps[j]

		allowed_component.traits.meshmap["meshmodel-metadata"].parentId == i
		updated_comp := apply_patch(allowed_component, service, from_selectors, to_selectors)
	]
}

identify_additions := result if {
	result := [ comps_to_add |
		relationship := data.relationships[_]
		relationship.subType in {"Inventory", "Parent"}

		selector_set := relationship.selectors[_]


		mutated_selector_set := {selector |
			some set in selector_set.allow
			selector := mutated_selectors(set)
			count(selector) > 0
		}

		umutated_selector_set := union(mutated_selector_set)
		
		mutator_selector_set := { selector |
			some set in selector_set.allow
			selector := mutator_selectors(set)
			count(selector) > 0
		}

		umutator_selector_set := union(mutator_selector_set)
		
	# extract all the components which gets mutated
	
		# because in a particular selector set comps are related we xan be sure that whathever results we get the relationship allows for cuh realtionship but we alos needs to consider the deny rules for that particualr set		
		some mutated_selector in umutated_selector_set
		mutated_components := extract_components_by_type(input.services, mutated_selector)

		some mutated_component in mutated_components

		print("\nMutated processing", mutated_selector.kind)
		print("type: ", mutated_component.type)

		mutated_values := extract_values(mutated_component, mutated_selector.patch.mutatedRef)

		# For all paths specifed in the ref, ensure the values exists.
		count(mutated_values) == count(mutated_selector.patch.mutatedRef)

		# for each such components, find corresponding mutator component.
		# if such component doesn't exist, it indicates a scenario where,
		# a component is referring other component of different type but it doesn't exist in the design, hence give the suggestion to the client to create the component with requried configuration.
		# eg for namespace and other namespaced resource, and second example of configmap - deployment.
		some mutator_selector in umutator_selector_set
			# every mutator_component in mutator_components {
			
			comps_to_add := process_comps_to_add(mutated_values, mutator_selector)
		# }
		# print("COMPS TO ADD: ", comps_to_add)
	]
}

process_comps_to_add(mutated_values, mutator_selector) := result if {
	mutator_components := extract_components_by_type(input.services, mutator_selector)

	every mutator_component in mutator_components {

		print("\n\nMutator processing")
		print("type: ", mutator_component.type)
		print("Mutator Component ", mutator_component, "mutatorRef: ", mutator_selector.patch.mutatorRef)

		mutator_values := extract_values(mutator_component, mutator_selector.patch.mutatorRef)

		print("\nProcess comps to add: ", mutated_values, mutator_values)
		print("match_object(mutated_values, mutator_values)", match_object(mutated_values, mutator_values))
		not match_object(mutated_values, mutator_values)

	}
	# NOTE: Currently if atleast one of the mutator comp doesn't exist with the configuration as present in the mutaed comp (eg Pod and 2 configmMap resource presetn in design but for one of the configmap configuation is not desired, but because 1 configmap conf is desired the results for additions should be emoty but right now it isn't the cae, so change the logic such that if for every mutator type conp no comp has desired conf only then add to result set)
	
	print("\n\nREACHED 136")

	values := { val |
		some val in mutated_values
	}
	# filter only those whch don't match later
	result := {
		"type": mutator_selector.kind,
		"model": mutator_selector.model,
		"path": mutator_selector.patch.mutatorRef,
		"values": values,
	}
}

extract_values(component, refs) := values if {
	values := { formatted_path: component_value |
		some ref in refs
		path := resolve_path(ref, component)
		formatted_path := format_json_path(path)

		component_value := object.get(component, formatted_path, null)
		component_value != null
	}

	print("Returning from extract_values", values)
}

apply_patch(mutator, mutated, from_selectors, to_selectors) := result if {
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
	result := {mutatedObj.mutated.traits.meshmap.id: mutated_design}
}

# Example used:
# {
#     "name": "test",
#     "services": {
#         "cm": {
#             "name": "sc-1",
#             "type": "ConfigMap",
#             "apiVersion": "v1",
#             "namespace": "stg",
#             "version": "v1.25.2",
#             "model": "kubernetes",
#             "traits": {
#                 "meshmap": {
#                     "id": "1b85f558-1798-4930-a6b1-bc08443b69asds",
#                     "meshmodel-metadata": {
#                         "parentId": "1b85f558-1798-4930-a6b1-bc08443b69ac"
#                     }
#                 }
#             }
#         },
#         "cm-2": {
#             "name": "cm-2",
#             "type": "ConfigMap",
#             "apiVersion": "v1",
#             "namespace": "prod",
#             "version": "v1.25.2",
#             "model": "kubernetes",
#             "traits": {
#                 "meshmap": {
#                     "id": "1b85f558-1798-4930-a6b1-bc08443b69asd3",
#                     "meshmodel-metadata": {
#                         "parentId": "1b85f558-1798-4930-a6b1-bc08443b69ab"
                        
#                     }
#                 },
#                 "label": "pod-di"             
#             }
#         },
#         "pod-di-2": {
#             "name": "pod-3",
#             "type": "Pod",
#             "apiVersion": "v1",
#             "namespace": "prod",
#             "version": "v1.25.2",
#             "model": "kubernetes",
#             "traits": {
#                 "meshmap": {
#                     "id": "1b85f558-1798-4930-a6b1-bc08443b69ac",
#                     "label": "pod-di"
#                 }
#             },
#             "settings": {
#                 "spec": {
#                     "containers": [
#                         {
#                             "envFrom": [
#                                 {
#                                     "configMapRef": {
#                                         "name": "scm-1"
#                                     }
#                                 }
#                             ]
#                         }
#                     ]
#                 }
#             }
#         },
#         "pod-di-3": {
#             "name": "pod-3",
#             "type": "Pod",
#             "apiVersion": "v1",
#             "namespace": "stg",
#             "version": "v1.25.2",
#             "model": "kubernetes",
#             "traits": {
#                 "meshmap": {
#                     "id": "1b85f558-1798-4930-a6b1-bc08443b69ab",
#                     "label": "pod-di"
#                 }
#             },
#             "settings": {
#                 "spec": {
#                     "containers": [
#                         {
#                             "envFrom": [
#                                 {
#                                     "configMapRef": {
#                                         "name": "scm-1"
#                                     }
#                                 }
#                             ]
#                         }
#                     ]
#                 }
#             }
#         }
#     }
# }

# Relationship used: Add the namespace as one of the patch refs for this ckinf of relationhsip as configmmap are namespaed resourceds and hence udpating(when user drop a configmap or ismisalr reosuce on the pod)/evaluating considering the namespace is ideal
# {"relationships":[
#   {
#   "schemaVersion": "relationships.meshery.io/v1alpha2",
#   "version": "v1.0.0",
#   "kind": "Hierarchical",
#   "metadata": {
#     "description": "A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component."
#   },
#   "model": {
#     "schemaVersion": "models.meshery.io/v1beta1",
#     "version": "v1.0.0",
#     "name": "kubernetes",
#     "model": {
#       "version": "v1.25.2"
#     },
#     "displayName": "Kubernetes",
#     "category": {
#       "name": "Orchestration \u0026 Management",
#       "metadata": null
#     },
#     "metadata": {}
#   },
#   "subType": "Inventory",
#   "evaluationQuery": "hierarchical_inventory_relationship",
#   "selectors": [
#     {
#       "allow": {
#         "from": [
#           {
#             "kind": "ConfigMap",
#             "model": "kubernetes",
#             "patch": {
#               "patchStrategy": "replace",
#               "mutatorRef": [
#                 [
#                   "name"
#                 ],
#                 [
#                   "namespace"
#                 ]
#               ],
#               "description": "In Kubernetes, ConfigMaps are a versatile resource that can be referenced by various other resources to provide configuration data to applications or other Kubnernetes resources.\n\nBy referencing ConfigMaps in these various contexts, you can centralize and manage configuration data more efficiently, allowing for easier updates, versioning, and maintenance of configurations in a Kubernetes environment."
#             }
#           }
#         ],
#         "to": [
#           {
#             "kind": "Deployment",
#             "model": "kubernetes",
#             "patch": {
#               "patchStrategy": "replace",
#               "mutatedRef": [
#                 [
#                   "settings",
#                   "spec",
#                   "template",
#                   "spec",
#                   "containers",
#                   "_",
#                   "envFrom",
#                   0,
#                   "configMapRef",
#                   "name"
#                 ],
#                 [
#                   "namespace"
#                 ]
#               ],
#               "description": "Deployments can reference ConfigMaps to inject configuration data into the Pods they manage. This is useful for maintaining consistent configuration across replica sets.\n\nThe keys from the ConfigMap will be exposed as environment variables to the containers within the pods managed by the Deployment."
#             }
#           },
#           {
#             "kind": "Pod",
#             "model": "kubernetes",
#             "patch": {
#               "patchStrategy": "replace",
#               "mutatedRef": [
#                 [
#                   "settings",
#                   "spec",
#                   "containers",
#                   "_",
#                   "envFrom",
#                   "0",
#                   "configMapRef",
#                   "name"
#                 ],
#                 [
#                   "namespace"
#                 ]
#               ],
#               "description": "ConfigMaps can be referenced in the Pod specification to inject configuration data into the Pod's environment.\n\nThe keys from the ConfigMap will be exposed as environment variables to the container within the Pod."
#             }
#           }
#         ]
#       },
#       "deny": {
#         "from": [],
#         "to": []
#       }
#     }
#   ]
# },
# {
#   "schemaVersion": "relationships.meshery.io/v1alpha2",
#   "version": "v1.0.0",
#   "kind": "Hierarchical",
#   "metadata": {
#     "description": "A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component."
#   },
#   "model": {
#     "schemaVersion": "models.meshery.io/v1beta1",
#     "version": "v1.0.0",
#     "name": "kubernetes",
#     "model": {
#       "version": "v1.25.2"
#     },
#     "displayName": "Kubernetes",
#     "category": {
#       "name": "Orchestration \u0026 Management",
#       "metadata": null
#     },
#     "metadata": {}
#   },
#   "subType": "Inventory",
#   "evaluationQuery": "hierarchical_inventory_relationship",
#   "selectors": [
#     {
#       "allow": {
#         "from": [
#           {
#             "kind": "ConfigMap",
#             "model": "kubernetes",
#             "patch": {
#               "patchStrategy": "replace",
#               "mutatorRef": [
#                 [
#                   "name"
#                 ]
#               ],
#               "description": "In Kubernetes, ConfigMaps are a versatile resource that can be referenced by various other resources to provide configuration data to applications or other Kubnernetes resources.\n\nBy referencing ConfigMaps in these various contexts, you can centralize and manage configuration data more efficiently, allowing for easier updates, versioning, and maintenance of configurations in a Kubernetes environment."
#             }
#           }
#         ],
#         "to": [
#           {
#             "kind": "Job",
#             "model": "kubernetes",
#             "patch": {
#               "patchStrategy": "replace",
#               "mutatedRef": [
#                 [
#                   "settings",
#                   "spec",
#                   "template",
#                   "spec",
#                   "containers",
#                   "_",
#                   "envFrom",
#                   "0",
#                   "configMapRef",
#                   "name"
#                 ]
#               ],
#               "description": "ConfigMaps can be referenced in the Pod template specification within the Job definition to inject configuration data into the pods.\n\nThe keys from the ConfigMap will be exposed as environment variables to the container within the pod created by the Job."
#             }
#           }
#         ]
#       },
#       "deny": {
#         "from": [],
#         "to": []
#       }
#     }
#   ]
# }]
# }