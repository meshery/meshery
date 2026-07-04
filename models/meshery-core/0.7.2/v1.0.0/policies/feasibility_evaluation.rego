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

feasible_relationship_selector_between(from_component, to_component, relationship) := feasible_selector if {
	some selector in relationship.selectors
	not is_selector_set_feasible_between(to_component, from_component, selector.deny.to, selector.deny.from)
	feasible_selector := is_selector_set_feasible_between(to_component, from_component, selector.allow.to, selector.allow.from)
}

is_selector_set_feasible_between(to_component, from_component, to_selectors, from_selectors) := feasible_selector if {
	some to_selector in to_selectors
	is_relationship_feasible(to_selector, to_component)

	some from_selector in from_selectors
	is_relationship_feasible(from_selector, from_component)

	feasible_selector := {
		"from": from_selector,
		"to": to_selector,
	}
}
