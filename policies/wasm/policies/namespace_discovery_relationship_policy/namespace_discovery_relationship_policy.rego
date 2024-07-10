# https://play.openpolicyagent.org/p/wFNhyGsIej
package hierarchical_policy

available_namespaces[namespace_key] {
  service := input.services[_];
  namespace_key = service.namespace
}

available_namespaces[namespace_key] {
  service := input.services[_];
  service.type == "Namespace";
  namespace_key = service.name
}

parent_child_mapping = { namespace_map: ns |
	some namespace_map
	available_namespaces[namespace_map]
    service := input.services[key]
    service.namespace
    service.namespace == namespace_map
    ns := {ns_serv: comp_id |
        some ns_serv
         x := input.services[ns_serv]
         x.namespace == namespace_map
         comp_id := x.traits.meshmap.id
    }
}

namespaces_to_create[namespaces2] {
	some namespace
    available_namespaces[namespace]
    print(namespace)
    ns_creation_status = is_present(check_namespace_present_status(namespace))
    print(ns_creation_status)
    ns_creation_status != true
    namespaces2 = namespace
}

# incase of present: {"present": true}, in absent, returns empty set
check_namespace_present_status(ns) = is_present {
  is_present := {"present": is |
    some svc
    s := input.services[svc]
    s.type == "Namespace"
    s.name == ns
    is = true
  }
}

# Is present function wraps the result of check_namespace_present_status to true or false
is_present(obj) {
	obj.present == true
}

is_present(obj) = pre {
 not obj.present
 pre = false
}




