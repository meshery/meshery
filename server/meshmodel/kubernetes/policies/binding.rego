package binding_policy

import data.common
import future.keywords.if

binding_relationship = results {
    # pattern := json.marshal(input.services)    
    from_selectors := { kind: selectors |
        s := data.selectors.allow.from[_]
        kind := s.kind
        selectors := s
    }

    to_selectors :=  { kind: selectors |
        t := data.selectors.allow.to[_]
        kind := t.kind
        selectors := t
    }

    # print(input.services)
    # contains "selectors.from" components only, eg: Role/ClusterRole(s) comps only
    from := common.extract_components(input.services, from_selectors)
    to := common.extract_components(input.services, to_selectors)

    # binding_selectors := { data.bindingResource : [] }
    binding_resources := common.extract_components(input.services, [data.bindingResource])

    # print("from:", from, "\n\n", "to:", to, "\n\n", "binding_resources\n\n",  binding_resources)
    # print("from-selectors: ", from_selectors, "to-selectors: ", to_selectors, "binding-selectors: ", data.bindingResource)
    services_map := { service.traits.meshmap.id: service |
        service := input.services[_]
    }

    results = [ result |
        some i, j, k
            resource := from[i]
            # some j
                binding_resource := binding_resources[j]

            # print(resource.type, binding_resource.type, from_selectors[resource.type])
            r := is_match(resource, binding_resource, from_selectors[resource.type])
            r == true  
            print("[[[]]]")
            # some k
                to_resource := to[k]
                q := is_match(to_resource, binding_resource, to_selectors[to_resource.type])
                q == true

        result := {
            "from": {
                "id": resource.traits.meshmap.id
            },
            "to": {
                "id": to.traits.meshmap.id
            },
            "binded_by": {
                "id": binding_resource.traits.meshmap.id
            }
        }
    ]
}

is_match(resource1, resource2, from_selectors) = match {
    # print("from: ", from_selectors)
    # match = true
    value1 := [ val |
        some i
        match_from := from_selectors.match.self[i]
        print(match_from)
        val := object.get(resource1, match_from, "") if {
            not common.contains(match_from, "_")
        } else := ["test"]
    ]
    
    value2 := [ val |
        some i
        match_to := from_selectors.match["bindingResource"][i]
        print(match_to)
        val := object.get(resource2, match_to, "")
    ]
    
    print("r1:", resource1.type, "v1:", value1, "r2:", resource2.type, "v2:", value2)

    match := value1 == value2
}