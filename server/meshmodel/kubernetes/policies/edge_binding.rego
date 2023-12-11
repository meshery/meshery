package meshmodel_policy

import data.common.contains
import data.common.extract_components
import data.common.get_array_pos
import future.keywords.in

binding_types := ["mount", "permission"]

edge_binding_relationship[results] {
	binding_type := binding_types[_]

	from_selectors := {kind: selectors |
		s := data[binding_type].selectors.allow.from[_]
		kind := s.kind
		selectors := s
	}

	to_selectors := {kind: selectors |
		t := data[binding_type].selectors.allow.to[_]

		# kind := t.kind
		selectors := t

		# match := t.match
		val := t.match[key]

		key != "self"

		# kind := t.kind
		kind := concat("#", {t.kind, key})
	}
	# contains "selectors.from" components only, eg: Role/ClusterRole(s) comps only
	from := extract_components(input.services, from_selectors)
	to := extract_components(input.services, to_selectors)
	binding_comps := {type |
		selector := data[binding_type].selectors.allow.from[_].match
		selector[key]
		key != "self"
		type = key
	}

	# This is a set of set as it contains results for a particular binding_type and each binding_type can be binded by different type of nodes.
	evaluation_results := {result |
		some comp in binding_comps
		binding_resources := extract_components(input.services, [{"kind": comp}])

		count(binding_resources) > 0

		result := evaluate with data.binding_resources as binding_resources with data.from as from with data.to as to with data.from_selectors as from_selectors with data.to_selectors as to_selectors
	}

	union_of_results := union(evaluation_results)
	edges_set := {e |
		some comp in from

		some edge in comp.traits.meshmap.edges
		contains(binding_types, lower(edge.data.subType))

		e := {
			"from": {"id": edge.data.source},
			"to": {"id": edge.data.target},
			"binded_by": {"id": edge.data.metadata.binded_by},
		}
	}

	# print("edges_set: ", edges_set)
	# print("edges_set - evaluation_results: ", edges_set - union_of_results)
	# print("union_of_results: ", union_of_results)

	results = {binding_type: {
		"edges_to_add": union_of_results,
		"edges_to_remove": edges_set - union_of_results,
	}}
}

evaluate[results] {
	services_map := {service.traits.meshmap.id: service |
		service := input.services[_]
	}
	some i, j, k
	resource := data.from[i]
	binding_resource := data.binding_resources[j]

	r := is_related(resource, binding_resource, data.from_selectors[resource.type])

	r == true

	to_resource := data.to[k]
	q := is_related(to_resource, binding_resource, data.to_selectors[concat("#",{to_resource.type, binding_resource.type})])
	q == true

	results = {
		"from": {"id": resource.traits.meshmap.id},
		"to": {"id": to_resource.traits.meshmap.id},
		"binded_by": {"id": binding_resource.traits.meshmap.id},
	}
}

is_related(resource1, resource2, from_selectors) {
print(resource1, resource2, from_selectors)
	match_results := [result |
		some i
		match_from := from_selectors.match.self[i]
		match_to := from_selectors.match[resource2.type][i]
		ans := is_feasible(match_from, match_to, resource1, resource2)
		ans == true
		result := true
	]

	# ensure all the atribute present in the match field are equal
	count(match_results) == count(from_selectors.match.self)

}

# If none of the match paths ("from" and "to") doesn't contain array field in between, then it is a normal lookup.
is_feasible(from, to, resource1, resource2) {
	not contains(to, "_")
	object.get(resource1, from, "") == object.get(resource2, to, "")
}

# If any of the match paths contains array field in between then the path needs to be resolved before checking for their equality.
# "from" or "to" any of them can contain array in their path hence "is_feasible" is overrided. 
is_feasible(from, to, resource1, resource2) {
	match(from, to, resource1, resource2)
}

is_feasible(from, to, resource1, resource2) {
	match(to, from, resource2, resource1)
}

match(from, to, resource1, resource2) {
	contains(to, "_")
	index := get_array_pos(to)
	prefix_path := array.slice(to, 0, index)
	suffix_path := array.slice(to, index + 1, count(to))

	resource_val := object.get(resource2, prefix_path, "")
	value := [val |
		v := resource_val[_]
		val := object.get(v, suffix_path, "")
	]

	some i in value
	i == object.get(resource1, from, "")
}
