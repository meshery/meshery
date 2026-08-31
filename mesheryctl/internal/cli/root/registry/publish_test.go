// # Copyright Meshery Authors
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

package registry

import (
	"testing"
)

// TestPublishCmdPreRunE_ArgValidation verifies system-specific argument
// validation for the publish command's PreRunE:
//   - meshery        → 4 or 5 args accepted (imgs-output-path is optional)
//   - remote-provider → exactly 5 args required (imgs-output-path is mandatory)
//   - website         → exactly 5 args required (imgs-output-path is mandatory)
func TestPublishCmdPreRunE_ArgValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		// ── meshery system ──────────────────────────────────────────────────────
		{
			name:    "meshery: 4 args accepted (imgs-output-path optional)",
			args:    []string{"meshery", "cred", "sheetID", "path/to/models"},
			wantErr: false,
		},
		{
			name:    "meshery: 5 args accepted",
			args:    []string{"meshery", "cred", "sheetID", "path/to/models", "path/to/imgs"},
			wantErr: false,
		},
		{
			name:    "meshery: 3 args rejected (too few)",
			args:    []string{"meshery", "cred", "sheetID"},
			wantErr: true,
		},
		{
			name:    "meshery: 6 args rejected (too many)",
			args:    []string{"meshery", "cred", "sheetID", "path/to/models", "path/to/imgs", "extra"},
			wantErr: true,
		},

		// ── remote-provider system ──────────────────────────────────────────────
		{
			name:    "remote-provider: 5 args accepted",
			args:    []string{"remote-provider", "cred", "sheetID", "path/to/models", "path/to/imgs"},
			wantErr: false,
		},
		{
			name:    "remote-provider: 5 args rejected (empty imgs-output-path)",
			args:    []string{"remote-provider", "cred", "sheetID", "path/to/models", ""},
			wantErr: true,
		},
		{
			name:    "remote-provider: 4 args rejected (imgs-output-path mandatory)",
			args:    []string{"remote-provider", "cred", "sheetID", "path/to/models"},
			wantErr: true,
		},
		{
			name:    "remote-provider: 3 args rejected",
			args:    []string{"remote-provider", "cred", "sheetID"},
			wantErr: true,
		},

		// ── website system ──────────────────────────────────────────────────────
		{
			name:    "website: 5 args accepted",
			args:    []string{"website", "cred", "sheetID", "path/to/models", "path/to/imgs"},
			wantErr: false,
		},
		{
			name:    "website: 5 args rejected (empty imgs-output-path)",
			args:    []string{"website", "cred", "sheetID", "path/to/models", ""},
			wantErr: true,
		},
		{
			name:    "website: 4 args rejected (imgs-output-path mandatory)",
			args:    []string{"website", "cred", "sheetID", "path/to/models"},
			wantErr: true,
		},
		{
			name:    "website: 3 args rejected",
			args:    []string{"website", "cred", "sheetID"},
			wantErr: true,
		},

		// ── global boundary ─────────────────────────────────────────────────────
		{
			name:    "unknown system rejected",
			args:    []string{"invalid-system", "cred", "sheetID", "path"},
			wantErr: true,
		},
		{
			name:    "0 args rejected",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := publishCmd.PreRunE(publishCmd, tt.args)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error, but got: %v", err)
				}
			}
		})
	}
}
