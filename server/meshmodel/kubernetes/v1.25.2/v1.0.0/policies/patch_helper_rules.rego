package patch_helper

import data.helper.has_key
import rego.v1

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj if {
	has_key(to_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutated,
		"path": to_selector.patch.mutatorRef,
	}
}

identifyMutator(from_selector, to_selector, mutator, mutated) := mutatorObj if {
	has_key(from_selector.patch, "mutatorRef")
	mutatorObj = {
		"mutator": mutator,
		"path": from_selector.patch.mutatorRef,
	}
}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj if {
	has_key(from_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutator,
		"path": from_selector.patch.mutatedRef,
	}
}

identifyMutated(from_selector, to_selector, mutator, mutated) := mutatedObj if {
	has_key(to_selector.patch, "mutatedRef")
	mutatedObj = {
		"mutated": mutated,
		"path": to_selector.patch.mutatedRef,
	}
}

mutator_selectors(selector_set) := {selector |
	some selector in selector_set
	contains_mutator_selector(selector)
}

mutated_selectors(selector_set) := {selector |
	some selector in selector_set
	contains_mutated_selector(selector)
}

contains_mutator_selector(selector) if {
	selector.patch.mutatorRef
}

contains_mutated_selector(selector) if {
	selector.patch.mutatedRef
}
