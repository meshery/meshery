package core

import (
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestEncodeRefUrl(t *testing.T) {
	testCases := []struct {
		name     string
		input    url.URL
		expected string
	}{
		{
			name:     "given empty url when EncodeRefUrl then return empty string",
			input:    url.URL{},
			expected: "",
		},
		{
			name:     "given root path url when EncodeRefUrl then return empty string",
			input:    url.URL{Path: "/"},
			expected: "",
		},
		{
			name:     "given relative url when EncodeRefUrl then return encoded string",
			input:    url.URL{Path: "/dashboard"},
			expected: "L2Rhc2hib2FyZA",
		},
		{
			name: "given absolute url with query params when EncodeRefUrl then return encoded string",
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
			assert.Equal(t, tc.expected, actual, tc.name)
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
			name:        "given valid base64 encoded string when DecodeRefURL then return decoded url",
			input:       "L2Rhc2hib2FyZA",
			expected:    "/dashboard",
			expectError: false,
		},
		{
			name:        "given empty string when DecodeRefURL then return empty string",
			input:       "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "given invalid base64 string when DecodeRefURL then return error",
			input:       "!!!not-base64!!!", 
			expected:    "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := DecodeRefURL(tc.input)

		if tc.expectError {
				assert.Error(t, err, tc.name)
			} else {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, tc.expected, actual, tc.name)
			}
		})
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	testCases := []struct {
		name  string
		input url.URL
	}{
		{
			name:  "given relative url when EncodeRefUrl and DecodeRefURL then return original url",
			input: url.URL{Path: "/settings/system"},
		},
		{
			name: "given absolute url when EncodeRefUrl and DecodeRefURL then return original url",
			input: url.URL{
				Scheme:   "https",
				Host:     "meshery.io",
				Path:     "/catalog",
				RawQuery: "type=wasm",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encoded := EncodeRefUrl(tc.input)
			decoded, err := DecodeRefURL(encoded)
			
			assert.NoError(t, err, tc.name)
			assert.Equal(t, tc.input.String(), decoded, tc.name)
		})
	}
}