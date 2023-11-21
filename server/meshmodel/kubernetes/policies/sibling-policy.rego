package meshmodel_policy

group_objects_with_matching_labels[unique_id] {
    input.services[namespace] != null
    unique_id = input.services[_].traits.meshmap.id
    labels = input.services[namespace].labels
    some key
    key = labels[key]
    input.services[other_namespace] != null
    other_namespace != namespace
    other_labels = input.services[other_namespace].labels
    key = other_labels[key]
}
