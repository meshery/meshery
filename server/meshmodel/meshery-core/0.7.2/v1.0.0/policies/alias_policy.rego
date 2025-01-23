package eval

import rego.v1

import data.core_utils.get_array_aware_configuration_for_component_at_path
import data.core_utils.new_uuid
import data.core_utils.from_component_id
import data.core_utils.component_declaration_by_id
import data.core_utils.to_component_id
import data.core_utils.get_component_configuration
import data.core_utils.pop_first
import data.core_utils.object_get_nested
import data.core_utils.component_alias
import data.feasibility_evaluation_utils.is_relationship_feasible_to
import data.feasibility_evaluation_utils.is_relationship_feasible_from


# ----------- Module Alias Policy ------------------------------------------------------------------------
# This module is responsible for evaluating alias relationships.
# Alias relationships occur when a component serves as an alias
# to a field in another component.
# Example: A Container acts as an alias to the field `spec.containers[x]`
# in a Pod (where `x` is the index of the container in the Pod).
#
# 1. Validation of Alias Relationships:
#    - A previously approved alias relationship remains valid if:
#        - The path it aliases still exists in the parent component.
#        - The child alias component is still present.
#    - The alias relationship becomes invalid if:
#        - The child component is missing (e.g., the user deleted the container component inside the Pod).
#        - The aliased configuration in the parent is missing
#          (e.g., the user deleted the `containers` field in the Pod).
#
# 2. Identification of Alias Relationships:
#    - An alias relationship exists if fields in the component are defined as aliases in the relationship block.
#    - The alias relationship is identified by the presence of a selector where:
#        - `from` specifies the child component.
#        - `to` specifies the parent component.
#    - The `mutatorRef`/`mutatedRef` must be a single-item array,
#      where the item at index `0` points to the field in the child component that acts as an alias to the parent.
#    - Newly identified relationships have a status of **pending**.
#
# 3. Action Phase:
#    - For each pending relationship:
#        - Add the alias component to the design file (if not already present).
#        - Update the status of the relationship to **approved**.
#    - For approved relationships:
#        - No action is required.
#    - For deleted relationships:
#        - Remove the alias component from the design file.
#        - Remove the aliased configuration from the parent component.


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
	"kind":"hierarchical",
	"type":"parent",
	"subtype":"alias"
  }
}

identify_relationships(design_file, relationships_in_scope,relationship_policy_identifier) := eval_results if {

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
	#print("Identify alias rels Eval results", count(eval_results))

}

 
identify_alias_relationships(component, relationship) := {rel |
	some selector in relationship.selectors
	some from in selector.allow.from # from is child or alias
	some to in selector.allow.to # to is parent

	# identify if alias can be created
	# identified_alias_paths := identify_alias_paths(from, to, component,input)

    array_items := get_array_aware_configuration_for_component_at_path(from.patch.mutatorRef[0], component,input)
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

	#print("Is valid -> from_component", from_component)

	# check if the to component is still present
	to_component := component_declaration_by_id(design_file, to_component_id(relationship))
	to_component != null

	#print("Is valid -> to_component", to_component)

	# check if the path in the to component is still present

	ref := alias_ref_from_relationship(relationship)

	#print("Is valid -> ref", ref,relationship.id)
	value := object_get_nested(get_component_configuration(to_component,design_file), pop_first(ref), null)
	#print("Is valid -> value", value)
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
validate_relationships_phase(design_file,relationship_policy_identifier) := result if  {

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
		"op":"add_component",
		"value":component
	}
}

alias_components_to_remove(design_file, alias_relationships) := {action |
	some relationship in alias_relationships
	relationship.status == "deleted"
	some selector in relationship.selectors
	some from in selector.allow.from

	component := component_declaration_by_id(design_file, from.id)

	action := {
		"op":"delete_component",
		"value":component
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
action_phase(design_file,relationship_policy_identifier) := result if {

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
			"op":"add_relationship",
			"value":rel
		}
	}

	# Relationships that are deleted already at the validation phase
	relationships_to_delete := {action |
		some alias_rel in alias_relationships
		alias_rel.status == "deleted"

		action:= {
			"op":"delete_relationship",
			"value":alias_rel
		}
	}
	components_to_delete := alias_components_to_remove(design_file,alias_relationships)
	
	# print("Relationships Deleted by alias policy", count(relationships_deleted))


	# print("Components added", count(components_added))

	result := components_to_add | components_to_delete | relationships_to_add | relationships_to_delete 

} 
