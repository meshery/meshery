package feasibility_evaluation_utils

import rego.v1

# wildcard_match
match_name(name, match) if {
	name == "*"
}

# exact_match
match_name(name, match) if {
	name == match
}

# regex_match
match_name(name, match) if {
	regex.match(match, name)
}

is_relationship_feasible(selector, comp) if {
	match_name(comp.component.kind, selector.kind)
	match_name(comp.model.name, selector.model.name)
}

is_relationship_feasible_to(component, relationship) := to if {
	some selector in relationship.selectors
	some to in selector.allow.to
	is_relationship_feasible(to, component)
}

is_relationship_feasible_from(component, relationship) := from if {
	some selector in relationship.selectors
	some from in selector.allow.from
	is_relationship_feasible(from, component)
}
