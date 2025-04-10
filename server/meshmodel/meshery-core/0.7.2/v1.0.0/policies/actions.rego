package actions

import rego.v1
import data.core_utils

# { "op": "update_component", "value": { "id": "xxx-xxx", "path": ["name"], "value": 2 }  op: "replace"| "merge" }
update_component_op := "update_component"


# This is specially handled as configurations are complex and need to be specialy patched ( lot of array index edgecases)
# right now this is handled from golang
# { "op": "update_component_configuration", "value": { "id": "xxx-xxx", "path": ["configuration","spec"], "value": 2 }  op: "replace"| "merge" }
update_component_configuration_op := "update_component_configuration"


# { "op": "delete_component", "value": { "id": "xxx-xxx" } }
delete_component_op := "delete_component"
# { "op": "add_component", "value": { item : <component_declaration> }
add_component_op := "add_component"

# { "op": "update_relationship", "value": { "id": "xxx-xxx", "path": ["status"], "value": "approve" } }
update_relationship_op := "update_relationship"
# { "op": "delete_relationship", "value": { "id": "xxx-xxx" } }
delete_relationship_op := "delete_relationship"
# { "op": "add_relationship", "value": { item : <relationship_declaration> } }
add_relationship_op := "add_relationship"

get_component_update_op(path) := op if {
    path[0] == "configuration"
    op := update_component_configuration_op
} else :=  update_component_op



# applies all of the updates for a single item like component or relationship
apply_updates_to_item(original,updates,op) := updated if {

    valid_updates := {update |
        some update in updates
        update.op == op
        update.value.id == original.id
    }

    count(valid_updates) > 0

    json_ops := [patch |
         some update in valid_updates
         path := core_utils.normalize_path(update.value.path)
         patch := {
             "op":   "replace",
             "path": path,
             "value": update.value.value,
         }
    ]

    updated := json.patch(original, json_ops)


} else  := original

filter_out_deleted(items,actions,delete_op) := filtered_items if {
    deleted_ids := {action.value.id |
       some action in actions
       action.op == delete_op
    }

    filtered_items := {item |
        some item in items
        not item.id in deleted_ids
    }
} else  := items

append_added_items(items,actions,add_op) := added if {
    added_items := {action.value.item |
        some action in actions
        action.op == add_op
    }
    added := items | added_items
} else  := items

apply_component_update_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/components",
            "value": {updated_component |
                some component in design_file.components
                updated_component := apply_updates_to_item(component, actions, update_component_op)
            },
        },
])

apply_component_add_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/components",
            "value": append_added_items(design_file.components, actions, add_component_op)
        },
])

apply_component_delete_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/components",
            "value": filter_out_deleted(design_file.components, actions, delete_component_op)
        },
])

apply_relationship_update_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/relationships",
            "value": {updated_rel |
                some rel in design_file.relationships
                updated_rel := apply_updates_to_item(rel, actions, update_relationship_op)
            },
        },
])

apply_relationship_add_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/relationships",
            "value": append_added_items(design_file.relationships, actions, add_relationship_op)
        },
])

apply_relationship_delete_actions(design_file,actions) := json.patch(design_file, [
        {
            "op": "replace",
            "path": "/relationships",
            "value": filter_out_deleted(design_file.relationships, actions, delete_relationship_op)
        },
])


apply_all_actions_to_design(design_file,actions) := final_design if {
   with_deleted_components := apply_component_delete_actions(design_file,actions)
   with_added_components := apply_component_add_actions(with_deleted_components,actions)
   with_updated_components := apply_component_update_actions(with_added_components,actions) # move to golang to handle complex paths

   with_deleted_relationships := apply_relationship_delete_actions(with_updated_components,actions)
   with_added_relationships := apply_relationship_add_actions(with_deleted_relationships,actions)
   with_updated_relationships := apply_relationship_update_actions(with_added_relationships,actions)

   final_design := with_updated_relationships
}

