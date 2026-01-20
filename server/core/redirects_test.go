package core

import (
	"net/url"
	"testing"
)

func TestIsValidRedirectURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		// Valid redirect URLs
		{
			name:     "empty string",
			url:      "",
			expected: true,
		},
		{
			name:     "root path",
			url:      "/",
			expected: true,
		},
		{
			name:     "relative path",
			url:      "/dashboard",
			expected: true,
		},
		{
			name:     "relative path with query",
			url:      "/dashboard?foo=bar",
			expected: true,
		},
		{
			name:     "relative path with fragment",
			url:      "/dashboard#section",
			expected: true,
		},
		{
			name:     "extension path",
			url:      "/extension/navigator",
			expected: true,
		},

		// Invalid redirect URLs - open redirect attempts
		{
			name:     "absolute URL http",
			url:      "http://evil.com",
			expected: false,
		},
		{
			name:     "absolute URL https",
			url:      "https://evil.com",
			expected: false,
		},
		{
			name:     "protocol-relative URL",
			url:      "//evil.com",
			expected: false,
		},
		{
			name:     "protocol-relative with path",
			url:      "//evil.com/path",
			expected: false,
		},
		{
			name:     "javascript protocol",
			url:      "javascript:alert('xss')",
			expected: false,
		},
		{
			name:     "data protocol",
			url:      "data:text/html,<script>alert('xss')</script>",
			expected: false,
		},

		// Invalid redirect URLs - CRLF injection attempts
		{
			name:     "CRLF with carriage return",
			url:      "/path\rLocation: http://evil.com",
			expected: false,
		},
		{
			name:     "CRLF with newline",
			url:      "/path\nLocation: http://evil.com",
			expected: false,
		},
		{
			name:     "CRLF with both",
			url:      "/path\r\nLocation: http://evil.com",
			expected: false,
		},

		// Edge cases
		{
			name:     "triple slash",
			url:      "///example.com",
			expected: false,
		},
		{
			name:     "backslash instead of slash",
			url:      "\\evil.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidRedirectURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidRedirectURL(%q) = %v, expected %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestEncodeDecodeRefURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "simple path",
			input:   "/dashboard",
			wantErr: false,
		},
		{
			name:    "path with query",
			input:   "/dashboard?foo=bar",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "complex URL",
			input:   "/api/user?id=123&name=test#section",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For empty or root paths, EncodeRefUrl returns empty string
			if tt.input == "" || tt.input == "/" {
				u, _ := url.Parse(tt.input)
				encoded := EncodeRefUrl(*u)
				if encoded != "" {
					t.Errorf("EncodeRefUrl(%q) should return empty string for root/empty paths", tt.input)
				}
				return
			}

			// Test encoding and decoding round-trip
			u, _ := url.Parse(tt.input)
			encoded := EncodeRefUrl(*u)
			if encoded == "" {
				t.Errorf("EncodeRefUrl(%q) returned empty string", tt.input)
				return
			}

			decoded, err := DecodeRefURL(encoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeRefURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && decoded != tt.input {
				t.Errorf("Round-trip failed: input = %q, decoded = %q", tt.input, decoded)
			}
		})
	}
}
