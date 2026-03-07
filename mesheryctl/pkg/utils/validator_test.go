package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsUUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string // Explicit name for Given, When, Then convention
		in   string
		want bool
	}{
		// Valid cases
		{
			name: "Given a valid standard UUID, then return true",
			in:   "123e4567-e89b-12d3-a456-426614174000",
			want: true,
		},
		{
			name: "Given a valid nil UUID, then return true",
			in:   "00000000-0000-0000-0000-000000000000",
			want: true,
		},
		{
			name: "Given a valid braced UUID, then return true",
			in:   "{123e4567-e89b-12d3-a456-426614174000}",
			want: true,
		},
		{
			name: "Given a valid URN UUID, then return true",
			in:   "urn:uuid:123e4567-e89b-12d3-a456-426614174000",
			want: true,
		},

		// Invalid cases
		{
			name: "Given an empty string, then return false",
			in:   "",
			want: false,
		},
		{
			name: "Given a UUID with leading whitespace, then return false",
			in:   " 123e4567-e89b-12d3-a456-426614174000",
			want: false,
		},
		{
			name: "Given a UUID with trailing whitespace, then return false",
			in:   "123e4567-e89b-12d3-a456-426614174000 ",
			want: false,
		},
		{
			name: "Given a string exceeding UUID length, then return false",
			in:   "123e4567-e89b-12d3-a456-426614174000-extra",
			want: false,
		},
		{
			name: "Given an invalid random string, then return false",
			in:   "not-a-uuid",
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := IsUUID(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}