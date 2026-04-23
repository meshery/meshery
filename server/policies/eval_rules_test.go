package policies

import (
	"testing"
)

// TestShouldSkipMutatorUpdate locks in the mutator/mutated skip predicate
// that guards patchMutatorsAction from emitting redundant ComponentUpdate
// actions and from the infinite-loop case when the mutator resolves to
// nil.
func TestShouldSkipMutatorUpdate(t *testing.T) {
	cases := []struct {
		name       string
		mutator    any
		old        any
		wantSkip   bool
		reasonHint string
	}{
		{
			name:       "identical non-nil values: skip (redundant update)",
			mutator:    "abc",
			old:        "abc",
			wantSkip:   true,
			reasonHint: "deep-equal short-circuit",
		},
		{
			name:       "different non-nil values: emit update",
			mutator:    "new",
			old:        "old",
			wantSkip:   false,
			reasonHint: "legitimate propagation",
		},
		{
			name:       "nil mutator vs. non-nil old: skip (infinite-loop guard)",
			mutator:    nil,
			old:        "existing",
			wantSkip:   true,
			reasonHint: "prevents re-emit loop until patch pipeline clears atomically",
		},
		{
			name:       "both nil: skip (no-op)",
			mutator:    nil,
			old:        nil,
			wantSkip:   true,
			reasonHint: "deep-equal short-circuit first, nil check redundant here",
		},
		{
			name:       "nil mutator, zero-string old: skip (still an infinite-loop guard)",
			mutator:    nil,
			old:        "",
			wantSkip:   true,
			reasonHint: "nil vs. \"\" are deep-unequal, but nil mutator trips the guard regardless",
		},
		{
			name:       "nested map equal: skip",
			mutator:    map[string]any{"a": 1, "b": "x"},
			old:        map[string]any{"a": 1, "b": "x"},
			wantSkip:   true,
			reasonHint: "reflect.DeepEqual through deepEqual",
		},
		{
			name:       "nested map different: emit",
			mutator:    map[string]any{"a": 2},
			old:        map[string]any{"a": 1},
			wantSkip:   false,
			reasonHint: "values differ",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSkipMutatorUpdate(tc.mutator, tc.old)
			if got != tc.wantSkip {
				t.Errorf("shouldSkipMutatorUpdate(%v, %v) = %v, want %v (%s)", tc.mutator, tc.old, got, tc.wantSkip, tc.reasonHint)
			}
		})
	}
}
