package relationship_evaluation_policy

import rego.v1

is_alias_relationship(relationship) if {
	lower(relationship.kind) == "hierarchical"
	lower(relationship.type) == "parent"
	lower(relationship.subType) == "alias"
}

identify_relationships(design_file, relationships_in_scope) := eval_results if {
	eval_results := union({new_relationships |
		some relationship in relationships_in_scope

		is_alias_relationship(relationship)
		some component in design_file.components

		print("is alias rel", component)
		is_relationship_feasible_to(component, relationship)
		print("rel is feasible")

		identified_relationships := identify_alias_relationships(component, relationship)
		new_relationships := { rel |
		    some rel in identified_relationships
	     	not alias_relationship_already_exists(design_file, rel)
	    }
	})

	print("Identify alias rels Eval results", eval_results)
}

new_uuid(seed) := id if {
	now := format_int(time.now_ns(), 10)
	id := uuid.rfc4122(sprintf("%s%s", [seed, now]))
}

identify_alias_relationships(component, relationship) := {rel |
	some selector in relationship.selectors
	some from in selector.allow.from
	some to in selector.allow.to


	selector_declaration := {
	  "allow" : {
		"from": [json.patch(from, [{
			"op": "replace",
			"path": "/id",
			"value": new_uuid(component),
		}])],
		"to": [json.patch(to, [{
			"op": "replace",
			"path": "/id",
			"value": component.id,
		}])]
	  },

	  "deny" :  {}
	}

	print("selector dec", selector)

	rel := json.patch(relationship, [
		{
			"op": "add",
			"path": "/selectors",
			"value": [selector_declaration],
		},
		{
			"op": "add",
			"path": "/id",
			"value": new_uuid(component.id),
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
    is_alias_relationship(existing_rel)

    some selector in existing_rel.selectors
    some to in selector.allow.to

    some new_selector in relationship.selectors
    some new_to in new_selector.allow.to

    to.kind == new_to.kind
    to.patch == new_to.patch
    to.id == new_to.id

}



is_relationship_feasible_to(component, relationship) := to if {
	some selector in relationship.selectors
	some to in selector.allow.to
	is_relationship_feasible(to, component.component.kind)
}

## Validate

validate_relationships_phase(design_file) := validated_rels if {
	# print("Validating relationships", design_file)
	validated_rels := design_file.relationships
}

## Action Phase

add_components_action(design_file,alias_relationships) := components_added if {

    components_added := {component |

        some relationship in alias_relationships
        relationship.status == "pending"
        some selector in relationship.selectors
        some from in selector.allow.from

        print("To Add", from)
        component := {
            "id": from.id,
            "component": {"kind": from.kind},
            "model": from.model,
        }
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
action_phase(design_file) := result if {
    alias_relationships := {rel |
        some rel in design_file.relationships
        is_alias_relationship(rel)
    }


    components_added := add_components_action(design_file, alias_relationships)
    relationships_added := { rel |
        some alias_rel in alias_relationships
        alias_rel.status == "pending"
        rel := json.patch(alias_rel, [{
            "op": "replace",
            "path": "/status",
            "value": "approved",
        }])
    }
    relationships_deleted := { alias_rel |
        some alias_rel in alias_relationships
        alias_rel.status == "deleteted"
    }

    print("Components added", count(components_added))

    result := {
        "components_added": components_added,
        "components_deleted": {},
        "components_updated": {},
        "relationships_updated": {},
        "relationships_added": relationships_added,
        "relationships_deleted": relationships_deleted,

    }

} else := {
    "components_added": {},
    "components_deleted": {},
    "components_updated": {},
    "relationships_updated": {},
    "relationships_added": {},
}
