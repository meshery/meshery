package eval

import rego.v1

import data.core_utils
import data.feasibility_evaluation_utils

import data.core_utils.component_alias
import data.core_utils.component_declaration_by_id
import data.core_utils.from_component_id
import data.core_utils.get_array_aware_configuration_for_component_at_path
import data.core_utils.get_component_configuration
import data.core_utils.new_uuid
import data.core_utils.object_get_nested
import data.core_utils.pop_first
import data.core_utils.to_component_id
import data.core_utils.truncate_set
import data.feasibility_evaluation_utils.is_relationship_feasible_from
import data.feasibility_evaluation_utils.is_relationship_feasible_to

import data.actions
import data.eval_rules

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

MAX_ALIASES := 20

# It is unlikely, that Meshery has a use case for supporting relationship.type == "child" aliases in the future.
is_alias_relationship(relationship) if {
	lower(relationship.subType) == "child"
}

is_alias_relationship(relationship) if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "parent"
	lower(relationship.subType) == "alias"
}

alias_policy_identifier := "alias_relationships_policy"

relationship_is_implicated_by_policy(relationship, policy_identifier) if {
	policy_identifier == alias_policy_identifier
	is_alias_relationship(relationship)
}

relationship_is_invalid(relationship, design_file, policy_identifier) if {
	policy_identifier == alias_policy_identifier
	not is_alias_relationship_valid(relationship, design_file)
}

relationship_already_exists(rel, design_file, policy_identifier) if {
	policy_identifier == alias_policy_identifier
	alias_relationship_already_exists(design_file, rel)
}

identify_relationship(rel_definition, design_file, policy_identifier) := identified_relationships if {
	policy_identifier == alias_policy_identifier
	identified_relationships := union({rels |
		some component in design_file.components

		# print("checking feasibility for component", rel_definition.metadata.description, component.id, rel_definition.kind, rel_definition.type, rel_definition.subType)
		is_relationship_feasible_to(component, rel_definition)
		rels := identify_alias_relationships(component, rel_definition)
	})
}

default_display_name_for_added_component(component, path) := display_name if {
	# print("default_display_name_for_added_component", component)
	length := count(path)
	length > 1
	display_name := sprintf("%s.%s", [path[length - 2], path[length - 1]])
}

default_display_name_for_added_component(component, path) := display_name if {
	length := count(path)
	length == 1
	display_name := sprintf("%s", [path[0]])
}

default_display_name_for_added_component(component, path) := display_name if {
	length := count(path)
	length == 0
	display_name := "default"
}

new_added_component_declaration(component, path, selector, relationship, design_file) := new_component if {
	# print("new_added_component_declaration", component, path, relationship, design_file)
	display_name := default_display_name_for_added_component(component, path)
	new_component := {
		"id": component.id,
		"component": {"kind": component.kind},
		"model": component.model,
		"displayName": display_name,
		"metadata": {"isAnnotation": relationship.subType == "alias"},
	}
}

# set configuration for added component ( not for alias)
relationship_side_effects(relationship, design_file, policy_identifier) := side_effects if {
	policy_identifier == alias_policy_identifier
	relationship.status == "added-component"
	relationship.subType != "alias"

	set_config_side_effects := {action |
		to := relationship.selectors[0].allow.to[0]
		from := relationship.selectors[0].allow.from[0]
		some i in numbers.range(0, count(from.patch.mutatorRef) - 1)
		mutator_ref := from.patch.mutatorRef[i]
		mutated_ref := to.patch.mutatedRef[i]
		to_component := core_utils.component_declaration_by_id(design_file, to.id)
		mutatedValue := core_utils.configuration_for_component_at_path(mutated_ref, to_component, design_file)

		action := {
			"op": actions.get_component_update_op(mutator_ref),
			"value": {
				"id": from.id,
				"path": mutator_ref,
				"value": mutatedValue,
			},
		}
	}

	update_status := {
		"op": actions.update_relationship_op,
		"value": {
			"id": relationship.id,
			"path": "/status",
			"value": "identified",
		},
	}
	side_effects := set_config_side_effects | {update_status}
}

# add components for relationships that have only been identified (with a status of "added-relationship")
relationship_side_effects(relationship, design_file, policy_identifier) := side_effects if {
	policy_identifier == alias_policy_identifier

	relationship.status == "added-relationship"

	add_side_effects := {action |
		some selector in relationship.selectors
		some from in selector.allow.from

		path := from.patch.mutatorRef[0]
		display_name := default_display_name_for_added_component(from, path)
		component := new_added_component_declaration(from, path, selector, relationship, design_file)
		action := {
			"op": actions.add_component_op,
			"value": {"item": component},
		}
	}

	update_status := {
		"op": actions.update_relationship_op,
		"value": {
			"id": relationship.id,
			"path": "/status",
			"value": "added-component",
		},
	}

	side_effects := add_side_effects | {update_status}
}

# remove alias components for deleted/invalid relationships
relationship_side_effects(relationship, design_file, policy_identifier) := side_effects if {
	policy_identifier == alias_policy_identifier

	relationship.status == "deleted"

	side_effects := {action |
		some selector in relationship.selectors
		some from in selector.allow.from

		action := {
			"op": actions.delete_component_op,
			"value": {"id": from.id, "component": component_declaration_by_id(design_file, from.id)},
		}
	}
}

## policy specific logic

alias_ref_from_relationship(relationship) := ref if {
	some selector in relationship.selectors
	some from in selector.allow.from
	ref := from.patch.mutatorRef[0]
}

selector_for_identified_relationship(to, from, component, relationship, design_file) := selector_declaration if {
	relationship.subType == "alias"

	array_items := get_array_aware_configuration_for_component_at_path(to.patch.mutatedRef[0], component, design_file)
	identified_alias_paths := array_items.paths

	count(identified_alias_paths) > 0 # if alias paths are present then alias can be created

	some path in identified_alias_paths

	selector_patch_declaration := {
		"patchStrategy": "replace",
		"mutatorRef": [path], # path to the component that needs to be mutated ( basically the ref to alias)
		"mutatedRef": [path],
	}

	selector_declaration := {
		"to": selector_patch_declaration,
		"from": selector_patch_declaration,
	}
}

# need to modify mutators schemas ( use a map rather than ordered arrays) to support raay wildcards in path
selector_for_identified_relationship(to, from, component, relationship, design_file) := selector_declaration if {
	relationship.subType != "alias"

	# if value is defined at the path then we can create an relationship
	some mutated_ref in to.patch.mutatedRef
	value := core_utils.configuration_for_component_at_path(mutated_ref, component, design_file)
	value != null
	value != ""

	selector_declaration := {
		"from": from.patch,
		"to": to.patch,
	}
}

identify_alias_relationships(component, relationship) := {rel |
	some selector in relationship.selectors
	some from in selector.allow.from # from is child or alias
	some to in selector.allow.to # to is parent

	# identify if alias can be created
	# identified_alias_paths := identify_alias_paths(from, to, component,input)

	selector_patch_declaration := selector_for_identified_relationship(to, from, component, relationship, input)

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
					"value": selector_patch_declaration.from,
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
					"value": selector_patch_declaration.to,
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
			"value": "added-relationship",
		},
	])
}

# this check is different than the comman duplicacy check in eval_rules.rego
# as the to selector is newly created with new id so we need to check
# if there is any existing relationship in the design of kind "alias" and the to ( alias component) is same
# only one alias relationship is allowed for a component at a particular path is allowed so we dont need to check for the from selector
alias_relationship_already_exists(design_file, relationship) if {
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
