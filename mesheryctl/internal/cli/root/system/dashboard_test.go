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

package system

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldUseEphemeralPortFallback(t *testing.T) {
	tests := []struct {
		name              string
		portExplicitlySet bool
		expected          bool
	}{
		{
			name:              "given port not explicitly set when fallback is evaluated then fallback is enabled",
			portExplicitlySet: false,
			expected:          true,
		},
		{
			name:              "given port explicitly set when fallback is evaluated then fallback is disabled",
			portExplicitlySet: true,
			expected:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := shouldUseEphemeralPortFallback(tt.portExplicitlySet)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestIsPortAlreadyInUseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "given address already in use error then returns true",
			err:      errors.New("listen tcp 127.0.0.1:9081: bind: address already in use"),
			expected: true,
		},
		{
			name:     "given bind error then returns true",
			err:      errors.New("bind failed for namespace/pod"),
			expected: true,
		},
		{
			name:     "given spdy negotiation error then returns false",
			err:      errors.New("error upgrading connection: malformed response"),
			expected: false,
		},
		{
			name:     "given nil error then returns false",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isPortAlreadyInUseError(tt.err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
