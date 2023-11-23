package meshmodel_policy
# https://play.openpolicyagent.org/p/EApfRBcqPq


label_obj[label_key] {
	some key
	svc := input.services[_]
	labels := svc.labels
	labels[key]
	label_key = key
}

group_objects_with_matching_labels = {labels_map: key |
	some labels_map
	label_obj[labels_map]
	svc := input.services[_]
	svc.labels
	some keys
	svc.labels[keys]
	keys == labels_map

	key := {value: result |
		svc := input.services[_]
		labels := svc.labels
		some value
		some k
		value = labels[k]
		k == labels_map
		result = {name: id |
			some name
			x := input.services[name]
			x.labels == labels
			id := x.traits.meshmap.id
		}
	}
}
