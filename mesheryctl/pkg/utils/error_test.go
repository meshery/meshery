// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"strings"
	"testing"

	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

// containsColonPrefix asserts that the error message contains the prefix followed by a colon and space.
func containsColonPrefix(t *testing.T, err error, prefix string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	msg := err.Error()
	expected := prefix + ": "
	if !strings.Contains(msg, expected) {
		t.Fatalf("error message %q does not contain expected prefix %q", msg, expected)
	}
}

// TestErrorDescriptionsSeparateCause validates description strings for constructors
// that append the underlying error with a colon separator.
func TestErrorDescriptionsSeparateCause(t *testing.T) {
	cause := fmt.Errorf("connection refused")

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrLoadConfig", ErrLoadConfig(cause), "Error processing config: connection refused"},
		{"ErrParseGithubFile", ErrParseGithubFile(cause, "https://github.com/meshery/meshery"), "Failed to parse github file: connection refused"},
		{"ErrRequestResponse", ErrRequestResponse(cause), "Unable to create a response from request: connection refused"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, meshkitErrors.GetLDescription(tt.err))
		})
	}
}

// TestErrorConstructorsFormatting checks that error constructors embed a colon+space prefix.
func TestErrorConstructorsFormatting(t *testing.T) {
	cause := fmt.Errorf("dummy error")

	// ErrLoadConfig (utils/error.go)
	containsColonPrefix(t, ErrLoadConfig(cause), "Error processing config")

	// ErrParseGithubFile (utils/error.go)
	containsColonPrefix(t, ErrParseGithubFile(cause, "http://example.com"), "Failed to parse github file")

	// ErrRequestResponse (utils/error.go)
	containsColonPrefix(t, ErrRequestResponse(cause), "Unable to create a response from request")
}
