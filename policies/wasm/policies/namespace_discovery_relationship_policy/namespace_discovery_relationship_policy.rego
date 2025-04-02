package hierarchical_policy

import rego.v1

available_namespaces contains service.namespace if {
	some service in input.services
}

available_namespaces contains service.name if {
	some service in input.services
	service.type == "Namespace"
}

parent_child_mapping[namespace_map] := ns if {
	some namespace_map in available_namespaces
	some key, service in input.services
	service.namespace == namespace_map
	ns := {ns_serv: comp_id |
		some ns_serv, x in input.services
		x.namespace == namespace_map
		comp_id := x.traits.meshmap.id
	}
}

namespaces_to_create contains namespace if {
	some namespace in available_namespaces
	not namespace_present(namespace, input.services)
}

# incase of present: {"present": true}, in absent, returns empty set
namespace_present(ns, all_services) if {
	some s in all_services
	s.type == "Namespace"
	s.name == ns
}
