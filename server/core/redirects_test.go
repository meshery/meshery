package core

import (
	"net/url"
	"testing"
)

func TestEncodeRefUrl(t *testing.T) {
	testCases := []struct {
		name     string
		input    url.URL
		expected string
	}{
		{
			name:     "Empty URL",
			input:    url.URL{},
			expected: "",
		},
		{
			name:     "Root path only",
			input:    url.URL{Path: "/"},
			expected: "",
		},
		{
			name:     "Relative URL",
			input:    url.URL{Path: "/dashboard"},
			expected: "L2Rhc2hib2FyZA",
		},
		{
			name: "Absolute URL with query param",
			input: url.URL{
				Scheme:   "http",
				Host:     "localhost:9081",
				Path:     "/login",
				RawQuery: "test=true",
			},
			expected: "aHR0cDovL2xvY2FsaG9zdDo5MDgxL2xvZ2luP3Rlc3Q9dHJ1ZQ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := EncodeRefUrl(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %s failed: expected %q, but got %q", tc.name, tc.expected, actual)
			}
		})
	}
}

func TestDecodeRefURL(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{
			name:        "Decode Dashboard path",
			input:       "L2Rhc2hib2FyZA",
			expected:    "/dashboard",
			expectError: false,
		},
		{
			name:        "Decode Empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "Invalid base64 input",
			input:       "!!!not-base64!!!", // Added an actual invalid case
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := DecodeRefURL(tc.input)

			// Verify if the error state matches what we expect
			if (err != nil) != tc.expectError {
				t.Fatalf("Test %s: DecodeRefURL() error = %v, wantErr %v", tc.name, err, tc.expectError)
			}

			// Only check the result if we didn't expect (and get) an error
			if !tc.expectError && actual != tc.expected {
				t.Errorf("Test %s failed: expected %q, but got %q", tc.name, tc.expected, actual)
			}
		})
	}
}
