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

func TestStartCmd(t *testing.T) {
	// Setup
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tests := []struct {
		name     string
		args     []string
		expecter func(t *testing.T)
	}{
		{
			name: "start command with docker platform",
			args: []string{"start", "--platform", "docker"},
			expecter: func(t *testing.T) {
				// Test that start command is properly configured
				assert.NotNil(t, startCmd)
				assert.Equal(t, "start", startCmd.Use)
				assert.Equal(t, "Start Meshery Nighthawk performance testing adapter", startCmd.Short)
			},
		},
		{
			name: "start command with kubernetes platform",
			args: []string{"start", "--platform", "kubernetes"},
			expecter: func(t *testing.T) {
				// Test that start command is properly configured
				assert.NotNil(t, startCmd)
				assert.Equal(t, "start", startCmd.Use)
				assert.Equal(t, "Start Meshery Nighthawk performance testing adapter", startCmd.Short)
			},
		},
		{
			name: "start command with skip-update-config flag",
			args: []string{"start", "--skip-update-config"},
			expecter: func(t *testing.T) {
				// Test that the flag is properly configured
				flag := startCmd.Flags().Lookup("skip-update-config")
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

func TestStartCmdFlags(t *testing.T) {
	// Test that all required flags are present
	tests := []struct {
		flagName string
		flagType string
	}{
		{"skip-update-config", "bool"},
		{"platform", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.flagName, func(t *testing.T) {
			flag := startCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.flagType, flag.Value.Type(), "Flag %s should be of type %s", tt.flagName, tt.flagType)
		})
	}
}

func TestStartCmdExamples(t *testing.T) {
	// Test that the command has proper examples
	assert.NotEmpty(t, startCmd.Example)
	assert.Contains(t, startCmd.Example, "mesheryctl perf start")
	assert.Contains(t, startCmd.Example, "--platform docker")
	assert.Contains(t, startCmd.Example, "--platform kubernetes")
	assert.Contains(t, startCmd.Example, "--skip-update-config")
}

func TestStartCmdValidation(t *testing.T) {
	// Test command structure
	assert.Equal(t, "start", startCmd.Use)
	assert.NotEmpty(t, startCmd.Short)
	assert.NotEmpty(t, startCmd.Long)
	assert.NotNil(t, startCmd.RunE)
}
