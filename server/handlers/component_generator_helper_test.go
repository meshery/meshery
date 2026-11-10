package handlers

import "testing"

func TestDeterminePluralWord(t *testing.T) {
	tests := []struct {
		name  string
		count int
		word  string
		want  string
	}{
		// A count of zero takes the plural form in English ("0 components").
		// This is the case the previous `count > 1` guard got wrong.
		{name: "zero is plural", count: 0, word: "component", want: "components"},
		{name: "zero is plural (relationship)", count: 0, word: "relationship", want: "relationships"},
		{name: "zero is plural (y-suffix)", count: 0, word: "entity", want: "entities"},
		{name: "zero is plural (s-suffix)", count: 0, word: "class", want: "classes"},

		// Exactly one is singular.
		{name: "one is singular", count: 1, word: "component", want: "component"},
		{name: "one is singular (y-suffix)", count: 1, word: "entity", want: "entity"},
		{name: "one is singular (s-suffix)", count: 1, word: "class", want: "class"},

		// Two or more is plural across each suffix branch.
		{name: "many default suffix", count: 2, word: "component", want: "components"},
		{name: "many y-suffix to ies", count: 3, word: "entity", want: "entities"},
		{name: "many s-suffix to es", count: 2, word: "class", want: "classes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determinePluralWord(tt.count, tt.word); got != tt.want {
				t.Errorf("determinePluralWord(%d, %q) = %q, want %q", tt.count, tt.word, got, tt.want)
			}
		})
	}
}
