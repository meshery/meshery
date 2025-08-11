package feasibility_evaluation_utils

import rego.v1

# wildcard_match ( if match is '*', then it matches any name )
match_name(name, match) if {
	match == "*"
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
	# print("is_relationship_feasible", comp.component.kind, comp.model.name, "-->", selector.kind, selector.model.name)
	match_name(comp.component.kind, selector.kind)
	match_name(comp.model.name, selector.model.name)
}

is_relationship_feasible_to(component, relationship) := to if {
	some selector in relationship.selectors
	some to in selector.allow.to

	# print("is_relationship_feasible_to", is_relationship_feasible(to, component))
	is_relationship_feasible(to, component)
}

is_relationship_feasible_from(component, relationship) := from if {
	some selector in relationship.selectors
	some from in selector.allow.from

	# print("is_relationship_feasible_from", is_relationship_feasible(from, component))
	is_relationship_feasible(from, component)
}
