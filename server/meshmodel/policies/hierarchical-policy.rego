
# https://play.openpolicyagent.org/p/yLQ0On30CT
package meshmodel_policy

extract_components(services, selectors) = components {
	components := [
	component |
		selector := selectors[_]
		service := services[_]
		selector.kind == service.type
		component := service
	]
}

has_key(x, k) {
	_ = x[k]
}

evaluate_hierarchical_relationship(services, selectors) = service_name {

	service_name := [service_name | 
        parent_comps := extract_components(services, selectors)
        component := parent_comps[_]

        service := services[_]
        has_key(service, lower(component.type))

        service.type != component.type
        service[lower(component.type)] == component.name

        service_name := {
            "destination_name": component.name,
            "source_name": service.name,
            "source_id": service.traits.meshmap.id,
            "destination_id": component.traits.meshmap.id,
            "namespace": component.name,
        }]
}

parent_child_relationship = service_name {
    print(data.selectors)
	allowed_selectors := data.selectors.allow.from
    
	service_name = evaluate_hierarchical_relationship(input.services, allowed_selectors)
    
}
