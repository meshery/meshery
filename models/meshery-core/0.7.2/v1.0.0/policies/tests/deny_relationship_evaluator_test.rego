package relationship_evaluation_policy_test

import rego.v1

import data.relationship_evaluation_policy

# Test is_relationship_denied with matching registrant.kind
test_is_relationship_denied_with_matching_registrant_kind if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-1",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-2",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
	}

	relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test is_relationship_denied with non-matching registrant.kind
test_is_relationship_not_denied_with_different_registrant_kind if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-1",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-2",
				"kind": "artifacthub",
				"status": "discovered",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
	}

	not relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test is_relationship_denied with different component kinds
test_is_relationship_not_denied_with_different_component_kinds if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Deployment"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-1",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-2",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": {"kind": "github"},
			},
		}],
	}

	not relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test is_relationship_denied with wildcard registrant
test_is_relationship_denied_with_wildcard_registrant if {
	from_declaration := {
		"id": "from-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-1",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	to_declaration := {
		"id": "to-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-2",
				"kind": "artifacthub",
				"status": "discovered",
			},
		},
	}

	deny_selectors := {
		"from": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": "*",
			},
		}],
		"to": [{
			"kind": "Namespace",
			"model": {
				"name": "kubernetes",
				"registrant": "*",
			},
		}],
	}

	relationship_evaluation_policy.is_relationship_denied(from_declaration, to_declaration, deny_selectors)
}

# Test any_selector_matches with registrant.kind
test_any_selector_matches_with_registrant_kind if {
	declaration := {
		"id": "test-id",
		"component": {"kind": "Namespace"},
		"model": {
			"name": "kubernetes",
			"registrant": {
				"id": "registrant-1",
				"kind": "github",
				"status": "discovered",
			},
		},
	}

	selector := {
		"kind": "Namespace",
		"model": {
			"name": "kubernetes",
			"registrant": {"kind": "github"},
		},
	}

	relationship_evaluation_policy.any_selector_matches(declaration, selector)
}

# Test is_selector_and_declaration_model_registrant_matches with nested kind
test_is_selector_and_declaration_model_registrant_matches_with_kind if {
	selector := {"model": {"registrant": {"kind": "github"}}}

	declaration := {"model": {"registrant": {
		"id": "registrant-1",
		"kind": "github",
		"status": "discovered",
	}}}

	relationship_evaluation_policy.is_selector_and_declaration_model_registrant_matches(selector, declaration)
}
