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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLogoutCmdArgs(t *testing.T) {
	expectedUsage := "Usage: mesheryctl system logout [flags]\nRun 'mesheryctl system logout --help' to see detailed help message"

	tests := []struct {
		Name          string
		Args          []string
		ExpectError   bool
		ExpectedError error
	}{
		{
			Name:        "valid with no positional arguments",
			Args:        []string{},
			ExpectError: false,
		},
		{
			Name:          "invalid with unexpected positional argument",
			Args:          []string{"unexpected-arg"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("unexpected argument(s): %q\n\n%s", []string{"unexpected-arg"}, expectedUsage)),
		},
		{
			Name:          "invalid with multiple unexpected positional arguments",
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("unexpected argument(s): %q\n\n%s", []string{"arg1", "arg2"}, expectedUsage)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := logoutCmd.Args(logoutCmd, tt.Args)
			if tt.ExpectError {
				assert.Error(t, err)
				utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLogoutCmdExecution runs `system logout` end to end against an isolated
// token file: an unexpected positional argument must fail before the handler
// runs, leaving the token untouched, while a bare invocation must succeed and
// clear the token.
func TestLogoutCmdExecution(t *testing.T) {
	tokenContent := []byte(`{"auths":{"meshery":{"token":"secret"}}}`)

	// setupLogoutEnv loads a private copy of the shared meshconfig fixture and
	// points the current context's token at a private file holding content.
	setupLogoutEnv := func(t *testing.T, content []byte) string {
		t.Helper()
		setupContextTestEnv(t)
		utils.SetupMeshkitLoggerTesting(t, false)

		tokenPath := filepath.Join(t.TempDir(), "token.json")
		if err := os.WriteFile(tokenPath, content, 0o600); err != nil {
			t.Fatalf("unable to write token fixture: %v", err)
		}
		viper.Set("tokens", []map[string]string{{"name": "default", "location": tokenPath}})
		return tokenPath
	}

	t.Run("unexpected argument fails and leaves the token file unchanged", func(t *testing.T) {
		tokenPath := setupLogoutEnv(t, tokenContent)

		SystemCmd.SetArgs([]string{"logout", "unexpected-arg"})
		if err := SystemCmd.Execute(); err == nil {
			t.Fatal("expected an error when passing a positional argument to 'system logout', got nil")
		}

		got, err := os.ReadFile(tokenPath)
		if err != nil {
			t.Fatalf("unable to read token file: %v", err)
		}
		if !bytes.Equal(got, tokenContent) {
			t.Errorf("token file was modified by the rejected invocation: got %q, want %q", got, tokenContent)
		}
	})

	t.Run("bare invocation succeeds and clears the token file", func(t *testing.T) {
		tokenPath := setupLogoutEnv(t, tokenContent)

		SystemCmd.SetArgs([]string{"logout"})
		if err := SystemCmd.Execute(); err != nil {
			t.Fatalf("expected no error when passing no arguments to 'system logout', got %v", err)
		}

		got, err := os.ReadFile(tokenPath)
		if err != nil {
			t.Fatalf("unable to read token file: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected token file to be emptied, got %q", got)
		}
	})
}
