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

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestStatusCmdArgs validates that system status rejects unexpected positional arguments
// with a structured utils.ErrInvalidArgument and usage guidance.
func TestStatusCmdArgs(t *testing.T) {
	expectedUsage := "Usage: mesheryctl system status [flags]\nRun 'mesheryctl system status --help' to see detailed help message"
	expectedErr := utils.ErrInvalidArgument(fmt.Errorf("unexpected argument(s)\n\n%s", expectedUsage))

	tests := []struct {
		name      string
		args      []string
		wantError error
	}{
		{
			name:      "valid with no positional arguments",
			args:      []string{},
			wantError: nil,
		},
		{
			name:      "invalid with unexpected positional argument",
			args:      []string{"unexpected-arg"},
			wantError: expectedErr,
		},
		{
			name:      "invalid with multiple unexpected positional arguments",
			args:      []string{"arg1", "arg2"},
			wantError: expectedErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := statusCmd.Args(statusCmd, tt.args)
			if tt.wantError == nil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			utils.AssertMeshkitErrorsEqual(t, tt.wantError, err)
		})
	}
}

// TestStatusCmdFlags verifies that the --verbose flag is registered correctly
// on statusCmd and is unaffected by argument validation.
func TestStatusCmdFlags(t *testing.T) {
	flag := statusCmd.Flags().Lookup("verbose")
	if flag == nil {
		t.Fatal("expected 'verbose' flag to be registered on statusCmd, got nil")
	}
	if flag.Value.Type() != "bool" {
		t.Errorf("expected 'verbose' flag to be boolean type, got %s", flag.Value.Type())
	}
}
