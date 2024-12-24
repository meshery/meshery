package relationship_evaluation_policy

import rego.v1

# This rule checks if a relationship between two declarations is denied based on deny selectors.
is_relationship_denied(from_declaration, to_declaration, deny_selectors) if {
	# Iterate over all denied selectors in 'from' and 'to' deny selectors.
	some denied_from_selector in deny_selectors.from
	some denied_to_selector in deny_selectors.to

	# Check if any selector matches the 'from' and 'to' declarations.
	any_selector_matches(from_declaration, denied_from_selector)
	any_selector_matches(to_declaration, denied_to_selector)
}

# Determines if any selector matches the given declaration.
any_selector_matches(declaration, selector) if {
	# Check if the kind matches.
	is_selector_and_declaration_kind_matches(selector, declaration)

	# Check if the model matches.
	is_selector_and_declaration_model_matches(selector, declaration)
}

# If properties are not present in the relationship selector, then the rule evaluates to undefined.
# and the expresssion "not is_relationship_denied" in other rules evaluates to true. And the relationship gets created.
is_selector_and_declaration_kind_matches(selector, declaration) if {
	selector.kind == "*"
}

is_selector_and_declaration_kind_matches(selector, declaration) if {
	selector.kind == declaration.component.kind
}

# Checks if both the model name and registrant match.
is_selector_and_declaration_model_matches(selector, declaration) if {
	is_selector_and_declaration_model_name_matches(selector, declaration)
	is_selector_and_declaration_model_registrant_matches(selector, declaration)
}

# Matches any model name if selector model name is '*', otherwise matches specific name.
is_selector_and_declaration_model_name_matches(selector, declaration) if {
	selector.model.name == "*"
}

is_selector_and_declaration_model_name_matches(selector, declaration) if {
	selector.model.name == declaration.model.name
}

# Matches any registrant if selector model registrant is '*', otherwise matches specific registrant.
is_selector_and_declaration_model_registrant_matches(selector, declaration) if {
	selector.model.registrant == "*"
}

is_selector_and_declaration_model_registrant_matches(selector, declaration) if {
	selector.model.registrant == declaration.model.registrant
}
