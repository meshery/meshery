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

