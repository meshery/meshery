package meshmodel_policy
#https://play.openpolicyagent.org/p/YcqqNCLuJx

label_obj[label_key] {
	some key
	svc := input.services[_]
	labels := svc.labels
	labels[key]
	label_key = key
}

annotation_obj[annotation_key] {
	some key
	svc := input.services[_]
	annotation := svc.annotations
	annotation[key]
	annotation_key = key
}


group_objects_with_matching_field = result {
	labels_result := {labels_map: key |
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
            result = {r |
                some name
                x := input.services[name]
                x.labels == labels
                id := x.traits.meshmap.id
                r := {"id": id, "name": name}
            }
        }
    } 
    
	annotation_result := {annotations_map: key |
        some annotations_map
        annotation_obj[annotations_map]
        svc := input.services[_]
        svc.annotations
        some keys
        svc.annotations[keys]
        keys == annotations_map

        key := {value: result |
            svc := input.services[_]
            annotations := svc.annotations
            some value
            some k
            value = annotations[k]
            k == annotations_map
            result = {r |
                some name
                x := input.services[name]
                x.annotations == annotations
                id := x.traits.meshmap.id
                r := {"id": id, "name": name}
            }
        }
    } 
    
    result := {"annotations": annotation_result, "labels": labels_result}
}
