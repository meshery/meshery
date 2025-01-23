package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils

# Module: Alias Relationship Evaluator
#
# Purpose: Manages relationships where one component references (aliases) a field in another component.
#
# Example: In Kubernetes, a Container component might alias the field `spec.containers[x]` within a Pod component, where 'x' represents the container's index in the Pod's array.
#
# Process Flow:
#
# 1. Validation
#    Valid relationships require:
#    ✓ Parent component still contains the referenced path
#    ✓ Child (alias) component still exists
#
#    Relationships become invalid when:
#    ✗ Child component is removed (e.g., Container deleted from Pod)
#    ✗ Referenced field is removed (e.g., `containers` array removed from Pod spec)
#
# 2. Identification
#    Alias relationships are detected by:
#    - Component fields marked as aliases in relationship block
#    - Selector configuration that specifies:
#      • 'from': The child (alias) component
#      • 'to': The parent component being referenced
#    - Single-item mutatorRef/mutatedRef array where index[0] points to
#      the child component field acting as the alias
#
#    Note: New relationships start with "pending" status
#
# 3. Actions
#    New  Relationships:
#    + Add alias component to design file (if missing)
#    + Update relationship status to "approved"
#
#    Approved Relationships:
#    • No action needed (maintain current state)
#
#    Deleted Relationships:
#    - Remove alias component from design file
#    - Clean up aliased configuration in parent component

# It is unlikely, that Meshery has a use case for supporting relationship.type == "child" aliases in the future.
is_alias_relationship(relationship) if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "parent"
	lower(relationship.subType) == "alias"
}

alias_ref_from_relationship(relationship) := ref if {
	some selector in relationship.selectors
	some from in selector.allow.from
	ref := from.patch.mutatorRef[0]
}

is_alias_policy_identifier(relationship_policy_identifier) if {
	relationship_policy_identifier == {
		"kind": "hierarchical",
		"type": "parent",
		"subtype": "alias",
	}
}

identify_relationships(design_file, relationships_in_scope, relationship_policy_identifier) := eval_results if {
	is_alias_policy_identifier(relationship_policy_identifier)

	eval_results := union({new_relationships |
		some relationship in relationships_in_scope

		is_alias_relationship(relationship)
		some component in design_file.components

		# print("is alias rel", component)
		is_relationship_feasible_to(component, relationship)

		# print("rel is feasible")

		identified_relationships := identify_alias_relationships(component, relationship)
		new_relationships := {rel |
			some rel in identified_relationships
			not alias_relationship_already_exists(design_file, rel)
		}
	})
	# print("Identify alias rels Eval results", count(eval_results))

}

identify_alias_relationships(component, relationship) := {rel |
	some selector in relationship.selectors
	some from in selector.allow.from # from is child or alias
	some to in selector.allow.to # to is parent

	# identify if alias can be created
	# identified_alias_paths := identify_alias_paths(from, to, component,input)

	array_items := get_array_aware_configuration_for_component_at_path(from.patch.mutatorRef[0], component, input)
	identified_alias_paths := array_items.paths

	# print("Identified Alias Paths", identified_alias_paths)

	count(identified_alias_paths) > 0 # if alias paths are present then alias can be created

	some path in identified_alias_paths

	selector_patch_declaration := {
		"patchStrategy": "replace",
		"mutatorRef": [path], # path to the component that needs to be mutated ( basically the ref to alias)
		"mutatedRef": [path],
	}

	# create alias relationship declaration
	selector_declaration := {
		"allow": {
			"from": [json.patch(from, [
				{
					"op": "replace",
					"path": "/id",
					"value": new_uuid({"c": component, "s": selector_patch_declaration}),
					# use both component and selector patch declaration to create a unique id
				},
				{
					"op": "replace",
					"path": "/patch",
					"value": selector_patch_declaration,
				},
			])],
			"to": [json.patch(to, [
				{
					"op": "replace",
					"path": "/id",
					"value": component.id,
				},
				{
					"op": "replace",
					"path": "/patch",
					"value": selector_patch_declaration,
				},
			])],
		},
		"deny": {},
	}

	# print("selector dec", selector)

	rel := json.patch(relationship, [
		{
			"op": "add",
			"path": "/selectors",
			"value": [selector_declaration],
		},
		{
			"op": "add",
			"path": "/id",
			"value": new_uuid(selector_declaration),
		},
		{
			"op": "replace",
			"path": "/status",
			"value": "pending",
		},
	])
}

alias_relationship_already_exists(design_file, relationship) := existing_rel if {
	some existing_rel in design_file.relationships
	existing_rel.status != "deleted" # check if the relationship is not deleted
	is_alias_relationship(existing_rel)

	some selector in existing_rel.selectors
	some to in selector.allow.to

	some new_selector in relationship.selectors
	some new_to in new_selector.allow.to

	to.kind == new_to.kind
	to.patch == new_to.patch
	to.id == new_to.id
}

## Validate

# validate relationship and return the updated relationship
is_alias_relationship_valid(relationship, design_file) if {
	relationship.status == "approved"

	# check if the from component is still present
	from_component := component_declaration_by_id(design_file, from_component_id(relationship))
	from_component != null

	# print("Is valid -> from_component", from_component)

	# check if the to component is still present
	to_component := component_declaration_by_id(design_file, to_component_id(relationship))
	to_component != null

	# print("Is valid -> to_component", to_component)

	# check if the path in the to component is still present

	ref := alias_ref_from_relationship(relationship)

	# print("Is valid -> ref", ref,relationship.id)
	value := object_get_nested(get_component_configuration(to_component, design_file), pop_first(ref), null)

	# print("Is valid -> value", value)
	value != null
}

is_alias_relationship_valid(relationship, design_file) if {
	relationship.status == "pending"
}

validate_relationship(relationship, design_file) := relationship if {
	is_alias_relationship_valid(relationship, design_file)
}

validate_relationship(relationship, design_file) := updated_relationship if {
	not is_alias_relationship_valid(relationship, design_file)
	updated_relationship := json.patch(relationship, [{
		"op": "replace",
		"path": "/status",
		"value": "deleted",
	}])
}

# validate all relationships in the design file ( use partial rule so it doesnt conflict with other policies)
validate_relationships_phase(design_file, relationship_policy_identifier) := result if {
	is_alias_policy_identifier(relationship_policy_identifier)

	result := {validated |
		some rel in design_file.relationships
		is_alias_relationship(rel)
		validated := validate_relationship(rel, design_file)
	}
}

## Action Phase

alias_components_to_add(design_file, alias_relationships) := {action |
	some relationship in alias_relationships
	relationship.status == "pending"
	some selector in relationship.selectors
	some from in selector.allow.from

	component := {
		"id": from.id,
		"component": {"kind": from.kind},
		"model": from.model,
	}

	action := {
		"op": "add_component",
		"value": component,
	}
}

alias_components_to_remove(design_file, alias_relationships) := {action |
	some relationship in alias_relationships
	relationship.status == "deleted"
	some selector in relationship.selectors
	some from in selector.allow.from

	component := component_declaration_by_id(design_file, from.id)

	action := {
		"op": "delete_component",
		"value": component,
	}
}

# action response {
#   components_added :      list of components added
#   components_deleted :    list of components deleted
#   components_updated :    list of components updated
#   relationships_added :   list of relationships added
#   relationships_deleted : list of relationships deleted
#   relationships_updated : list of relationships updated
# }
action_phase(design_file, relationship_policy_identifier) := result if {
	is_alias_policy_identifier(relationship_policy_identifier)

	alias_relationships := {rel |
		some rel in design_file.relationships
		is_alias_relationship(rel)
	}

	components_to_add := alias_components_to_add(design_file, alias_relationships)
	relationships_to_add := {action |
		some alias_rel in alias_relationships
		alias_rel.status == "pending"
		rel := json.patch(alias_rel, [{
			"op": "replace",
			"path": "/status",
			"value": "approved",
		}])
		action := {
			"op": "add_relationship",
			"value": rel,
		}
	}

	# Relationships that are deleted already at the validation phase
	relationships_to_delete := {action |
		some alias_rel in alias_relationships
		alias_rel.status == "deleted"

		action := {
			"op": "delete_relationship",
			"value": alias_rel,
		}
	}
	components_to_delete := alias_components_to_remove(design_file, alias_relationships)

	# print("Relationships Deleted by alias policy", count(relationships_deleted))

	# print("Components added", count(components_added))

	result := ((components_to_add | components_to_delete) | relationships_to_add) | relationships_to_delete
}
