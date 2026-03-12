package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUUID(t *testing.T) {
	// Table-driven test cases
	cases := []struct {
		in   string
		want bool
	}{
		// Valid UUIDs
		{"123e4567-e89b-12d3-a456-426614174000", true},
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"550E8400-E29B-41D4-A716-446655440000", true},
		{"00000000-0000-0000-0000-000000000000", true}, // nil UUID
		{"123e4567e89b12d3a456426614174000", true},     // no-hyphen form

		// Invalid UUIDs
		{"", false},
		{"not-a-uuid", false},
		{"123", false},
		{"g123e4567-e89b-12d3-a456-426614174000", false}, // non-hex char
		{"123e4567-e89b-12d3-a456-42661417400", false},   // too short
		{"123e4567e89b12d3a45642661417400", false},       // too short no-hyphen
		{"0000-0000-0000-0000-000000000000", false},      // wrong hyphen placement
	}

	for _, tc := range cases {
		tc := tc // Capture range variable
		
		// Handle display logic for the test name
		displayInput := tc.in
		if displayInput == "" {
			displayInput = "empty_string"
		}

		testName := fmt.Sprintf("Given_input_%s_When_IsUUID_is_called_Then_it_should_return_%v", displayInput, tc.want)

		t.Run(testName, func(t *testing.T) {
			got := IsUUID(tc.in)
			assert.Equal(t, tc.want, got, "Input: %s", tc.in)
		})
	}
}