package eval_rules

import rego.v1

import data.core_utils.truncate_set


approve_pending_relationships_action(relationships,max_limit) := relationships_to_add if {
   relationships_to_add :=  truncate_set({action |
		some pending_rel in relationships

		pending_rel.status == "pending"

		rel := json.patch(pending_rel, [{
			"op": "replace",
			"path": "/status",
			"value": "approved",
		}])

		action := {
			"op": "add_relationship",
			"value": rel,
		}
	},max_limit)

}

same_relationship_identitfier(rel_a,rel_b) := true if {
    rel_a.kind == rel_b.kind
    rel_a.type == rel_b.type
    rel_a.subType == rel_b.subType
}

#selector clause is from or to
same_relationship_selector_clause(clause_a,clause_b) := true if {
   clause_a.kind == clause_b.kind
   clause_a.id == clause_b.id
   clause_a.patch == clause_b.patch
}

relationships_are_same(rel_a,rel_b) := true if {

   same_relationship_identitfier(rel_a,rel_b)

    some selector_a in rel_a.selectors
    some selector_b in rel_b.selectors

    same_relationship_selector_clause(selector_a.allow.from[0],selector_b.allow.from[0])
    same_relationship_selector_clause(selector_a.allow.to[0],selector_b.allow.to[0])
}

relationship_already_exists(design_file, relationship) := true if {
    some existing_rel in design_file.relationships
    existing_rel.status != "deleted" # check if the relationship is not deleted

    relationships_are_same(existing_rel, relationship)
}

#any_comman_selector(relationship_a, relationship_b) := true if {
#    some selector_a in relationship_a.selectors
#    some selector_b in relationship_b.selectors
#
#    some from_a in selector_a.allow.from
#    some from_b in selector_b.allow.from
#
#    some to_a in selector_a.allow.to
#    some to_b in selector_b.allow.to
#
#    same_relationship_selector_clause(from_a,from_b)
#    same_relationship_selector_clause(to_a,to_b)
#}
#
#same_from_clauses(relationship_a, relationship_b) := true if {
#     # checks if all the s
#
##    some selector_a in relationship_a.selectors
##    some selector_b in relationship_b.selectors
##
##    some from_a in selector_a.allow.from
##    some from_b in selector_b.allow.from
##
##
##    same_relationship_selector_clause(from_a,from_b)
#}
#
#any_comman_to_clause(relationship_a, relationship_b) := true if {
#    some selector_a in relationship_a.selectors
#    some selector_b in relationship_b.selectors
#
#    some to_a in selector_a.allow.to
#    some to_b in selector_b.allow.to
#
#    same_relationship_selector_clause(to_a,to_b)
#}

# checks if a valid relationship already exists in the design file
# make it or
#relationship_already_exists(design_file, relationship) := existing_rel if {
#    print("checking for existing rel",relationship.kind,relationship.type,relationship.subType,relationship.status,count(design_file.relationships))
#
#	some existing_rel in design_file.relationships
#	existing_rel.status != "deleted" # check if the relationship is not deleted
#
#	relationship.kind == existing_rel.kind
#	relationship.type == existing_rel.type
#	relationship.subType == existing_rel.subType
#
#
#	some selector in existing_rel.selectors
#	some to in selector.allow.to
#	some from in selector.allow.from
#
#	some new_selector in relationship.selectors
#	some new_to in new_selector.allow.to
#	some new_from in new_selector.allow.from
#
#    print("check","old",to.id , to.kind,"-->" ,from.id , from.kind,"new" , new_to.id , new_to.kind,"-->" ,new_from.id , new_from.kind)
#
##    print("to",to.id , new_to.id, to.kind, new_to.kind)
#	to.patch == new_to.patch
#	to.id == new_to.id
#
##	print("from",from.id , new_from.id, from.kind, new_from.kind)
#	from.patch == new_from.patch
#	from.id == new_from.id
#
#	print("matched")
#
#}



