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

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestDashboardPortFallbackConfiguration(t *testing.T) {
	tests := []struct {
		name                      string
		setupFlags                func(*cobra.Command)
		expectDefaultPort         int
		expectedPortExplicitlySet bool
	}{
		{
			name: "Default port configuration",
			setupFlags: func(cmd *cobra.Command) {
				// Don't set any flags, use defaults
			},
			expectDefaultPort:         9081,
			expectedPortExplicitlySet: false,
		},
		{
			name: "Explicitly set port with -p flag",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("port", "8080")
			},
			expectDefaultPort:         8080,
			expectedPortExplicitlySet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test command
			testCmd := &cobra.Command{
				Use: "test",
			}

			testCmd.Flags().IntVarP(&localPort, "port", "p", defaultPort, "test port flag")

			// Setup flags according to test case
			tt.setupFlags(testCmd)

			// Verify the flag state
			portExplicitlySet := testCmd.Flags().Changed("port")
			assert.Equal(t, tt.expectedPortExplicitlySet, portExplicitlySet, "port explicitly set flag should match")
		})
	}
}

func TestDefaultPortConstant(t *testing.T) {
	// Verify the default port is set correctly
	assert.Equal(t, 9081, defaultPort, "default port should be 9081")
}

func TestDashboardPortFlagValidation(t *testing.T) {
	tests := []struct {
		name        string
		portValue   string
		shouldError bool
	}{
		{
			name:        "Valid port number",
			portValue:   "8080",
			shouldError: false,
		},
		{
			name:        "Valid ephemeral port indicator",
			portValue:   "0",
			shouldError: false,
		},
		{
			name:        "Default port",
			portValue:   "9081",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCmd := &cobra.Command{
				Use: "test",
			}

			testCmd.Flags().IntVarP(&localPort, "port", "p", defaultPort, "test port flag")
			err := testCmd.Flags().Set("port", tt.portValue)

			if tt.shouldError {
				assert.Error(t, err, "setting port should error")
			} else {
				assert.NoError(t, err, "setting port should not error")
			}
		})
	}
}

func TestEphemeralPortFallbackLogic(t *testing.T) {
	tests := []struct {
		name                           string
		portExplicitlySet              bool
		currentPort                    int
		defaultPort                    int
		shouldAttemptEphemeralFallback bool
	}{
		{
			name:                           "Attempt fallback when default port is in use and not explicitly set",
			portExplicitlySet:              false,
			currentPort:                    9081,
			defaultPort:                    9081,
			shouldAttemptEphemeralFallback: true,
		},
		{
			name:                           "Do not fallback when port is explicitly set",
			portExplicitlySet:              true,
			currentPort:                    8080,
			defaultPort:                    9081,
			shouldAttemptEphemeralFallback: false,
		},
		{
			name:                           "Do not fallback when user set non-default port",
			portExplicitlySet:              true,
			currentPort:                    8080,
			defaultPort:                    9081,
			shouldAttemptEphemeralFallback: false,
		},
		{
			name:                           "Attempt fallback when default port is in use and not changed",
			portExplicitlySet:              false,
			currentPort:                    9081,
			defaultPort:                    9081,
			shouldAttemptEphemeralFallback: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the logic from dashboard.go
			shouldFallback := !tt.portExplicitlySet && tt.currentPort == tt.defaultPort
			assert.Equal(t, tt.shouldAttemptEphemeralFallback, shouldFallback,
				"fallback decision should match expected behavior")
		})
	}
}

func TestDashboardCommandStructure(t *testing.T) {
	// Verify the dashboard command is properly configured
	assert.NotNil(t, dashboardCmd, "dashboard command should be defined")
	assert.Equal(t, "dashboard", dashboardCmd.Use, "command name should be 'dashboard'")
	assert.NotEmpty(t, dashboardCmd.Short, "command should have a short description")
	assert.NotEmpty(t, dashboardCmd.Example, "command should have examples")
	assert.NoError(t, dashboardCmd.ValidateArgs(nil), "command should accept no arguments")
}

func TestPortForwardFlags(t *testing.T) {
	// Verify all port-forward related flags are properly defined
	t.Run("port-forward flag exists", func(t *testing.T) {
		flag := dashboardCmd.Flags().Lookup("port-forward")
		assert.NotNil(t, flag, "port-forward flag should exist")
		assert.Equal(t, "", flag.Shorthand, "port-forward should not have shorthand")
	})

	t.Run("port flag exists", func(t *testing.T) {
		flag := dashboardCmd.Flags().Lookup("port")
		assert.NotNil(t, flag, "port flag should exist")
		assert.Equal(t, "p", flag.Shorthand, "port should have shorthand 'p'")
	})

	t.Run("skip-browser flag exists", func(t *testing.T) {
		flag := dashboardCmd.Flags().Lookup("skip-browser")
		assert.NotNil(t, flag, "skip-browser flag should exist")
		assert.Equal(t, "", flag.Shorthand, "skip-browser should not have shorthand")
	})
}

func TestFlagHelp(t *testing.T) {
	// Verify flags have updated help text that mentions the fallback behavior
	t.Run("port flag help mentions ephemeral port", func(t *testing.T) {
		flag := dashboardCmd.Flags().Lookup("port")
		assert.NotNil(t, flag)
		// Help text should mention ephemeral port fallback
		assert.Contains(t, flag.Usage, "ephemeral",
			"port flag help should mention ephemeral port fallback")
	})
}
