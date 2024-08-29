package relationship_evaluation_policy

import rego.v1

# Evaluates the relationships which needs to be added based on the current state of design file.

# "identified_relationships": Always contains valid set of relationships with status as "approved"

evaluate_relationships_deleted(
	design_relationships,
	identified_relationships,
) := [rel |
	some existing_rel in design_relationships

	existing_rel.status != "deleted"

	# if the existing rel is not present in the identified_relationships,
	# it indicates it must be deleted.

	not does_relationship_exist_in_design(identified_relationships, existing_rel)
	not existing_rel.subType == "annotation"
	relationship := json.patch(existing_rel, [{
		"op": "replace",
		"path": "/status",
		"value": "deleted",
	}])
	rel := relationship
]

# Evaluates the relationships which needs to be added based on the current state of design file.

evaluate_relationships_added(
	design_relationships,
	identified_relationships,
) := [relationship |
	some rel in identified_relationships

	not does_relationship_exist_in_design(design_relationships, rel)
	relationship := rel
]

does_relationship_exist_in_design(relationships, relationship) if {
	some rel in relationships
	rel.status == relationship.status
	is_of_same_kind(rel, relationship)
	does_belongs_to_same_model(rel, relationship)
	is_same_selector(rel, relationship)
}

is_of_same_kind(existing_rel, new_rel) if {
	lower(new_rel.kind) == lower(existing_rel.kind)
	lower(new_rel.type) == lower(existing_rel.type)
	lower(new_rel.subType) == lower(existing_rel.subType)
}

# default does_belongs_to_same_model := undefined

# consider for wildcard in name and version?
does_belongs_to_same_model(existing_rel, new_rel) if {
	object.get(existing_rel.model, "name", "") == object.get(new_rel.model, "name", "")

	# use model.model.version and same model check o mode defined in selectrors to be added
	object.get(existing_rel.model, "version", "") == object.get(new_rel.model, "version", "")
}

is_same_selector(existing_rel, new_rel) if {
	some ex_selector in existing_rel.selectors
	some selector in new_rel.selectors

	# the rels that are inside design file, have only one element in the "from" and "to".
	ex_from_selector := ex_selector.allow.from[0]
	from_selector := selector.allow.from[0]

	ex_to_selector := ex_selector.allow.to[0]
	to_selector := selector.allow.to[0]

	# is relationship between same components or different
	ex_from_selector.id == from_selector.id
	ex_to_selector.id == to_selector.id
	# check if the relationship includes binding component.
	# If present, verify the binding component for the existing and the identified relationship are same or different.

}

is_same_binding(ex_from_selector, from_selector, ex_to_selector, to_selector) if {
	has_key(ex_from_selector, "match")
	has_key(from_selector, "match")

	is_same_binding_declaration(ex_from_selector, ex_to_selector, from_selector, to_selector)
}

is_same_binding(ex_from_selector, from_selector, ex_to_selector, to_selector) if {
	not has_key(ex_from_selector, "match")
	not has_key(from_selector, "match")
}

is_same_binding_declaration(ex_from_selector, ex_to_selector, from_selector, to_selector) if {
	existing_binding := {
		ex_from_selector.match.from[0].id,
		ex_from_selector.match.to[0].id,
		ex_to_selector.match.to[0].id,
	}

	identified_binding := {
		from_selector.match.from[0].id,
		from_selector.match.to[0].id,
		to_selector.match.to[0].id,
	}

	existing_binding == identified_binding
}
