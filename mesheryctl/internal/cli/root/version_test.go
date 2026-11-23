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

package root

import (
	"strings"
	"testing"
)

func TestVersionCmdArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
	}{
		{
			name:      "no arguments passed",
			args:      []string{},
			expectErr: false,
		},
		{
			name:      "extra positional argument passed",
			args:      []string{"extra"},
			expectErr: true,
		},
		{
			name:      "multiple positional arguments passed",
			args:      []string{"foo", "bar"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if versionCmd.Args == nil {
				t.Fatalf("expected versionCmd.Args to be configured")
			}
			err := versionCmd.Args(versionCmd, tt.args)
			if (err != nil) != tt.expectErr {
				t.Errorf("versionCmd.Args() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestVersionCmdProperties(t *testing.T) {
	if versionCmd.Use != "version" {
		t.Errorf("expected Use to be 'version', got '%s'", versionCmd.Use)
	}

	if !strings.Contains(versionCmd.Example, "mesheryctl version") {
		t.Errorf("expected Example to contain 'mesheryctl version', got: %s", versionCmd.Example)
	}

	if versionCmd.Short == "" {
		t.Errorf("expected Short description to be non-empty")
	}
}
