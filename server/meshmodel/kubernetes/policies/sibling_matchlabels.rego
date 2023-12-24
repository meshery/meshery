package meshmodel_policy

get_tag_keys(tag_type) := tag_keys {
	tag_keys := {tag_key |
		some tag_key
		svc := input.services[_]
		tags := svc[tag_type]
		tags[tag_key]
	}
}

group_nodes(tag_type, tag_key) := key {
	key := {value: group |
		svc := input.services[_]
		tags := svc[tag_type]
		some value
		some key
		value = tags[key]
		key == tag_key
		group = {g |
			some name
			x := input.services[name]
			x[tag_type][key]
			x[tag_type][key] == value
			id := x.traits.meshmap.id
			g := {"id": id, "name": name}
		}
	}
}

sibling_matchlabels_relationship = result {
	labels_result := {tag_labels: key |
		tag_labels_array := get_tag_keys("labels")
		tag_labels_array[tag_labels]
		key := group_nodes("labels", tag_labels)
	}

	annotations_result := {tag_annotations: key |
		tag_annotations_array := get_tag_keys("annotations")
		tag_annotations_array[tag_annotations]
		key := group_nodes("annotations", tag_annotations)
	}

	result := {"annotation": annotations_result, "labels": labels_result}
}
