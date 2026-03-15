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
		name     string
		input    string
		expected string
	}{
		{
			name:     "Decode Dashboard path",
			input:    "L2Rhc2hib2FyZA",
			expected: "/dashboard",
		},
		{
			name:     "Decode Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We add '_ ' here to ignore the error return value for now
			actual, _ := DecodeRefURL(tc.input)
			if actual != tc.expected {
				t.Errorf("Test %s failed: expected %q, but got %q", tc.name, tc.expected, actual)
			}
		})
	}
}
