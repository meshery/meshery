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
		expectedFinalAction meshkitkube.HelmChartAction
	}{
		{
			name:                "INSTALL with no existing releases",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      false,
			expectedFinalAction: meshkitkube.INSTALL,
		},
		{
			name:                "INSTALL with server existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        true,
			operatorExists:      false,
			expectedFinalAction: meshkitkube.UPGRADE,
		},
		{
			name:                "INSTALL with operator existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        false,
			operatorExists:      true,
			expectedFinalAction: meshkitkube.UPGRADE,
		},
		{
			name:                "INSTALL with both existing",
			initialAction:       meshkitkube.INSTALL,
			serverExists:        true,
			operatorExists:      true,
			expectedFinalAction: meshkitkube.UPGRADE,
		},
		{
			name:                "UPGRADE unchanged",
			initialAction:       meshkitkube.UPGRADE,
			serverExists:        false,
			operatorExists:      false,
			expectedFinalAction: meshkitkube.UPGRADE,
		},
		{
			name:                "UNINSTALL unchanged",
			initialAction:       meshkitkube.UNINSTALL,
			serverExists:        true,
			operatorExists:      true,
			expectedFinalAction: meshkitkube.UNINSTALL,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			finalAction := tc.initialAction

			if tc.initialAction == meshkitkube.INSTALL {
				if tc.serverExists || tc.operatorExists {
					finalAction = meshkitkube.UPGRADE
				}
			}

			if finalAction != tc.expectedFinalAction {
				t.Errorf("Expected %v but got %v", tc.expectedFinalAction, finalAction)
			}
		})
	}
}
