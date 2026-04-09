package models

import (
	"encoding/base64"
	"testing"
)

func TestResolvePostLoginRedirect(t *testing.T) {
	t.Parallel()

	const fallback = "/"

	tests := []struct {
		name     string
		rawRef   string
		expected string
	}{
		{
			name:     "empty ref falls back",
			rawRef:   "",
			expected: fallback,
		},
		{
			name:     "encoded in app path is decoded",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap")),
			expected: "/extension/meshmap",
		},
		{
			name:     "plain in app path is preserved",
			rawRef:   "/extension/meshmap",
			expected: "/extension/meshmap",
		},
		{
			name:     "encoded absolute url falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/phish")),
			expected: fallback,
		},
		{
			name:     "plain absolute url falls back",
			rawRef:   "https://evil.example/phish",
			expected: fallback,
		},
		{
			name:     "invalid base64 falls back",
			rawRef:   "not-base64",
			expected: fallback,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := resolvePostLoginRedirect(tc.rawRef, fallback)
			if actual != tc.expected {
				t.Fatalf("expected redirect %q, got %q", tc.expected, actual)
			}
		})
	}
}
