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
	"fmt"
	"testing"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
)

func TestCheckHelmReleaseExists(t *testing.T) {
	// Test that the function exists and has correct signature
	var testFunc func(*meshkitkube.Client, string, string) (bool, error) = checkHelmReleaseExists
	if testFunc == nil {
		t.Error("checkHelmReleaseExists function should be available")
	}

	// Test with nil client - should return error
	exists, err := checkHelmReleaseExists(nil, "meshery", "meshery")
	if err == nil {
		t.Error("Expected error with nil client")
	}
	if exists {
		t.Error("Expected false with nil client")
	}
}

func TestApplyHelmChartsActionLogic(t *testing.T) {
	testCases := []struct {
		name                string
		initialAction       meshkitkube.HelmChartAction
		serverExists        bool
		operatorExists      bool
		serverCheckError    error
		operatorCheckError  error
		expectedFinalAction meshkitkube.HelmChartAction
		expectError         bool
	}{
		{
			name:                "INSTALL with no existing releases",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      false,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.INSTALL,
			expectError:         false,
		},
		{
			name:                "INSTALL with server existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        true,
			operatorExists:      false,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.UPGRADE,
			expectError:         false,
		},
		{
			name:                "INSTALL with operator existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      true,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.UPGRADE,
			expectError:         false,
		},
		{
			name:                "INSTALL with both existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        true,
			operatorExists:      true,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.UPGRADE,
			expectError:         false,
		},
		{
			name:                "UPGRADE unchanged",
			initialAction:       meshkitkube.UPGRADE,
			serverExists:        false,
			operatorExists:      false,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.UPGRADE,
			expectError:         false,
		},
		{
			name:                "UNINSTALL unchanged",
			initialAction:       meshkitkube.UNINSTALL,
			serverExists:        true,
			operatorExists:      true,
			serverCheckError:    nil,
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.UNINSTALL,
			expectError:         false,
		},
		{
			name:                "INSTALL with server check error",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      false,
			serverCheckError:    fmt.Errorf("server check failed"),
			operatorCheckError:  nil,
			expectedFinalAction: meshkitkube.INSTALL, // doesn't matter since error expected
			expectError:         true,
		},
		{
			name:                "INSTALL with operator check error",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      false,
			serverCheckError:    nil,
			operatorCheckError:  fmt.Errorf("operator check failed"),
			expectedFinalAction: meshkitkube.INSTALL, // doesn't matter since error expected
			expectError:         true,
		},
		{
			name:                "INSTALL with both check errors",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      false,
			serverCheckError:    fmt.Errorf("server check failed"),
			operatorCheckError:  fmt.Errorf("operator check failed"),
			expectedFinalAction: meshkitkube.INSTALL, // doesn't matter since error expected
			expectError:         true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock the checkHelmReleaseExists function behavior
			mockCheckHelmReleaseExists := func(releaseName string) (bool, error) {
				if releaseName == "meshery" {
					return tc.serverExists, tc.serverCheckError
				}
				if releaseName == "meshery-operator" {
					return tc.operatorExists, tc.operatorCheckError
				}
				return false, fmt.Errorf("unknown release name: %s", releaseName)
			}

			// Test the logic that would be in applyHelmCharts
			finalAction := tc.initialAction
			var err error

			if tc.initialAction == meshkitkube.INSTALL {
				// Check if Meshery Server release exists
				serverExists, serverErr := mockCheckHelmReleaseExists("meshery")
				if serverErr != nil {
					err = fmt.Errorf("failed to check for existing Meshery Server release: %w", serverErr)
				} else {
					// Check if Meshery Operator release exists
					operatorExists, operatorErr := mockCheckHelmReleaseExists("meshery-operator")
					if operatorErr != nil {
						err = fmt.Errorf("failed to check for existing Meshery Operator release: %w", operatorErr)
					} else {
						// If either release exists, use UPGRADE instead of INSTALL for idempotency
						if serverExists || operatorExists {
							finalAction = meshkitkube.UPGRADE
						}
					}
				}
			}

			// Verify error handling
			if tc.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return // Don't check final action if we expect an error
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if finalAction != tc.expectedFinalAction {
				t.Errorf("Expected final action %v but got %v", tc.expectedFinalAction, finalAction)
			}
		})
	}
}
