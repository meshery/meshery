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

package perf

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestStopCmd(t *testing.T) {
	// Setup
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name     string
		args     []string
		expecter func(t *testing.T)
	}{
		{
			name: "stop command with docker platform",
			args: []string{"stop", "--platform", "docker"},
			expecter: func(t *testing.T) {
				// Test that stop command is properly configured
				assert.NotNil(t, stopCmd)
				assert.Equal(t, "stop", stopCmd.Use)
				assert.Equal(t, "Stop Meshery Nighthawk performance testing adapter", stopCmd.Short)
			},
		},
		{
			name: "stop command with kubernetes platform",
			args: []string{"stop", "--platform", "kubernetes"},
			expecter: func(t *testing.T) {
				// Test that stop command is properly configured
				assert.NotNil(t, stopCmd)
				assert.Equal(t, "stop", stopCmd.Use)
				assert.Equal(t, "Stop Meshery Nighthawk performance testing adapter", stopCmd.Short)
			},
		},
		{
			name: "stop command with force flag",
			args: []string{"stop", "--force"},
			expecter: func(t *testing.T) {
				// Test that the flag is properly configured
				flag := stopCmd.Flags().Lookup("force")
				assert.NotNil(t, flag)
				assert.Equal(t, "bool", flag.Value.Type())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expecter(t)
		})
	}
}

func TestStopCmdFlags(t *testing.T) {
	// Test that all required flags are present
	tests := []struct {
		flagName string
		flagType string
	}{
		{"force", "bool"},
		{"platform", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.flagName, func(t *testing.T) {
			flag := stopCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestStopCmdExamples(t *testing.T) {
	// Test that the command has proper examples
	assert.NotEmpty(t, stopCmd.Example)
	assert.Contains(t, stopCmd.Example, "mesheryctl perf stop")
	assert.Contains(t, stopCmd.Example, "--platform docker")
	assert.Contains(t, stopCmd.Example, "--platform kubernetes")
	assert.Contains(t, stopCmd.Example, "--force")
}

func TestStopCmdValidation(t *testing.T) {
	// Test command structure
	assert.Equal(t, "stop", stopCmd.Use)
	assert.NotEmpty(t, stopCmd.Short)
	assert.NotEmpty(t, stopCmd.Long)
	assert.NotNil(t, stopCmd.RunE)
}
