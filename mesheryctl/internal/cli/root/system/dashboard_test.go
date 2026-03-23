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
