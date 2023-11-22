package meshmodel_policy
# https://play.openpolicyagent.org/p/EApfRBcqPq

group_objects_with_matching_labels[key]{
	key := {result |
		some value
		some k
		svc := input.services[_]
		labels := svc.labels
		value = labels[k]
		comp := {name: id |
			some name
			x := input.services[name]
			x.labels == labels
			id := x.traits.meshmap.id
		}
		result := {k: {value: comp}}
	}
}
