package feasibility_evaluation_utils

import rego.v1

is_relationship_feasible(selector, comp_type) if {
	selector.kind == "*"
}

is_relationship_feasible(selector, comp_type) if {
	selector.kind == comp_type
}

is_relationship_feasible_to(component, relationship) := to if {
	some selector in relationship.selectors
	some to in selector.allow.to
	is_relationship_feasible(to, component.component.kind)
}

is_relationship_feasible_from(component, relationship) := from if {
	some selector in relationship.selectors
	some from in selector.allow.from
	is_relationship_feasible(from, component.component.kind)
}
