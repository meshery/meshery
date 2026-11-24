package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUUID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"Valid UUID with hyphens", "123e4567-e89b-12d3-a456-426614174000", true},
		{"Valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", true},
		{"Valid Nil UUID", "00000000-0000-0000-0000-000000000000", true},
		{"Valid UUID no hyphens", "123e4567e89b12d3a456426614174000", true},
		{"Empty string", "", false},
		{"Invalid random string", "not-a-uuid", false},
		{"Invalid short string", "123", false},
		{"Invalid non-hex character", "g123e4567-e89b-12d3-a456-426614174000", false},
		{"Invalid length too short", "123e4567-e89b-12d3-a456-42661417400", false},
		{"Invalid length too short no-hyphen", "123e4567e89b12d3a45642661417400", false},
		{"Invalid hyphen placement", "0000-0000-0000-0000-000000000000", false},
	}

	for _, tc := range cases {
		tc := tc 
		t.Run(tc.name, func(t *testing.T) {
			got := IsUUID(tc.in)
			assert.Equal(t, tc.want, got, "Input: %s", tc.in)
		})
	}
}