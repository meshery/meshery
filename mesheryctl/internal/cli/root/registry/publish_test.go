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
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// resetPublishFlags restores every publish flag to its default so subtests
// do not leak flag values into each other.
func resetPublishFlags(t *testing.T) {
	t.Helper()
	publishCmd.Flags().VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			t.Fatalf("failed to reset flag %q: %v", f.Name, err)
		}
		f.Changed = false
	})
}

// TestPublishCmdPreRunE_FlagValidation verifies the flag validation of the
// publish command's PreRunE:
//   - meshery         → --imgs-output-path is optional
//   - remote-provider → --imgs-output-path is required
//   - website         → --imgs-output-path is required
//   - positional arguments are rejected
func TestPublishCmdPreRunE_FlagValidation(t *testing.T) {
	validFlags := func(sys string) map[string]string {
		return map[string]string{
			"system":             sys,
			"spreadsheet-cred":   "cred",
			"spreadsheet-id":     "sheetID",
			"models-output-path": "path/to/models",
		}
	}
	with := func(flags map[string]string, key, value string) map[string]string {
		flags[key] = value
		return flags
	}
	without := func(flags map[string]string, key string) map[string]string {
		delete(flags, key)
		return flags
	}

	tests := []struct {
		name        string
		flags       map[string]string
		args        []string
		wantErr     bool
		errContains string
	}{
		// ── meshery system ──────────────────────────────────────────────────────
		{
			name:    "meshery: accepted without --imgs-output-path",
			flags:   validFlags("meshery"),
			wantErr: false,
		},
		{
			name:    "meshery: accepted with --imgs-output-path",
			flags:   with(validFlags("meshery"), "imgs-output-path", "path/to/imgs"),
			wantErr: false,
		},

		// ── remote-provider system ──────────────────────────────────────────────
		{
			name:    "remote-provider: accepted with --imgs-output-path",
			flags:   with(validFlags("remote-provider"), "imgs-output-path", "path/to/imgs"),
			wantErr: false,
		},
		{
			name:        "remote-provider: rejected without --imgs-output-path",
			flags:       validFlags("remote-provider"),
			wantErr:     true,
			errContains: "--imgs-output-path is required for the 'remote-provider' system",
		},
		{
			name:        "remote-provider: rejected with empty --imgs-output-path",
			flags:       with(validFlags("remote-provider"), "imgs-output-path", ""),
			wantErr:     true,
			errContains: "--imgs-output-path is required for the 'remote-provider' system",
		},

		// ── website system ──────────────────────────────────────────────────────
		{
			name:    "website: accepted with --imgs-output-path",
			flags:   with(validFlags("website"), "imgs-output-path", "path/to/imgs"),
			wantErr: false,
		},
		{
			name:        "website: rejected without --imgs-output-path",
			flags:       validFlags("website"),
			wantErr:     true,
			errContains: "--imgs-output-path is required for the 'website' system",
		},
		{
			name:        "website: rejected with whitespace-only --imgs-output-path",
			flags:       with(validFlags("website"), "imgs-output-path", "   "),
			wantErr:     true,
			errContains: "--imgs-output-path is required for the 'website' system",
		},

		// ── required flags ──────────────────────────────────────────────────────
		{
			name:        "no flags rejected",
			flags:       map[string]string{},
			wantErr:     true,
			errContains: "missing required flag(s): --models-output-path, --spreadsheet-cred, --spreadsheet-id, --system",
		},
		{
			name:        "missing --system rejected",
			flags:       without(validFlags("meshery"), "system"),
			wantErr:     true,
			errContains: "missing required flag(s): --system",
		},
		{
			name:        "missing --spreadsheet-cred rejected",
			flags:       without(validFlags("meshery"), "spreadsheet-cred"),
			wantErr:     true,
			errContains: "missing required flag(s): --spreadsheet-cred",
		},
		{
			name:        "missing --spreadsheet-id rejected",
			flags:       without(validFlags("meshery"), "spreadsheet-id"),
			wantErr:     true,
			errContains: "missing required flag(s): --spreadsheet-id",
		},
		{
			name:        "missing --models-output-path rejected",
			flags:       without(validFlags("meshery"), "models-output-path"),
			wantErr:     true,
			errContains: "missing required flag(s): --models-output-path",
		},

		// ── global boundary ─────────────────────────────────────────────────────
		{
			name:        "unknown system rejected",
			flags:       with(validFlags("invalid-system"), "imgs-output-path", "path/to/imgs"),
			wantErr:     true,
			errContains: "invalid system: 'invalid-system'",
		},
		{
			name:        "legacy positional args without flags rejected",
			flags:       map[string]string{},
			args:        []string{"website", "cred", "sheetID", "path/to/models", "path/to/imgs"},
			wantErr:     true,
			errContains: "unexpected positional arguments: website cred sheetID path/to/models path/to/imgs",
		},
		{
			name:        "positional args with otherwise valid flags rejected",
			flags:       with(validFlags("website"), "imgs-output-path", "path/to/imgs"),
			args:        []string{"extra"},
			wantErr:     true,
			errContains: "unexpected positional arguments: extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetPublishFlags(t)
			t.Cleanup(func() { resetPublishFlags(t) })

			for name, value := range tt.flags {
				if err := publishCmd.Flags().Set(name, value); err != nil {
					t.Fatalf("failed to set flag %q: %v", name, err)
				}
			}

			err := publishCmd.PreRunE(publishCmd, tt.args)
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("did not expect an error, but got: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error, but got nil")
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Fatalf("expected error to contain %q, got: %v", tt.errContains, err)
			}
		})
	}
}
