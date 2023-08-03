package hierarchichal
# https://play.openpolicyagent.org/p/I25U0udl0I

namespace_Keys[namespace_key] {
  service := input.services[_];
  namespace_key = service.namespace
}

namespaces = { namespace_map: ns |
	some namespace_map
	namespace_Keys[namespace_map]
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

