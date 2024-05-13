package meshmodel_policy


# gets unique tag keys for a given tag type
# example:
#   get_tag_keys("labels") := {"app", "env" "version"} for input  {"labels": {"app": "redis","env": "stg", "version": "edge"}}
#   get_tag_keys("annotations") := {"meshery.io/namespace", "meshery.io/service"} for input {"annotations": {"meshery.io/namespace": "default", "meshery.io/service": "redis"}}
get_tag_keys(tag_type) := tag_keys {
	tag_keys := {tag_key |
		some tag_key
		svc := input.services[_]
		tags := svc[tag_type]
		tags[tag_key]
	}
}



# groups nodes based on a given tag type and tag key
# example:
#   group_nodes("labels", "app") := {
#       "redis": [{"id": "1", "name": "prod-ns"}}, {"id": "2", "name": "prod-pod"}], 
#       "sql": [{"id": "3", "name": "stg-ns"}}]
#   }
#
#  for input:
# {
#     "name": "rel-policy-demo",
#     "services": {
#         "prod-ns": {
#             "apiVersion": "v1",
#             "labels": {
#                 "app": "redis",
#                 "env": "stg",
#             },
#             "name": "prod-ns",
#             "traits": {
#                 "meshmap": {
#                     "id": "1"
#                 }
#             }
#         },
#         "prod-pod": {
#             "labels": {
#                 "app": "redis"
#             },
#             "name": "prod-pod",
#             "traits": {
#                 "meshmap": {
#                     "id": "2"
#                 }
#             }
#         },
#         "stg-ns": {
#             "labels": {
#                 "app": "sql"
#             },
#             "name": "stg-ns",
#             "traits": {
#                 "meshmap": {
#                     "id": "3"
#                 }
#             }
#         },
#     }
# }

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


# Rego query
sibling_matchlabels_relationship = result {
    # group nodes based on labels
	labels_result := {tag_labels: key |
		tag_labels_array := get_tag_keys("labels")
		tag_labels_array[tag_labels]
		key := group_nodes("labels", tag_labels)
	}

    # group nodes based on annotations
	annotations_result := {tag_annotations: key |
		tag_annotations_array := get_tag_keys("annotations")
		tag_annotations_array[tag_annotations]
		key := group_nodes("annotations", tag_annotations)
	}

    # final evaluation result
	result := {"annotations": annotations_result, "labels": labels_result}
}
