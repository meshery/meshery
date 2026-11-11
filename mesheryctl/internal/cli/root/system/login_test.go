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
)

// TestLoginCmdArgs pins the argument validator of `system login`:
// the command accepts configuration via flags (e.g. -p/--provider) and takes
// no positional arguments, so any supplied positional argument must be rejected
// before the login handler runs.
func TestLoginCmdArgs(t *testing.T) {
	if err := loginCmd.Args(loginCmd, []string{"unexpected-arg"}); err == nil {
		t.Error("expected an error when passing a positional argument to 'system login', got nil")
	}
	if err := loginCmd.Args(loginCmd, []string{}); err != nil {
		t.Errorf("expected no error when passing no arguments to 'system login', got %v", err)
	}
}

func TestLoginCmdExecutionRejectsPositionalArgs(t *testing.T) {
	setupContextTestEnv(t)

	SystemCmd.SetArgs([]string{"login", "unexpected-arg"})
	if err := SystemCmd.Execute(); err == nil {
		t.Fatal("expected an error when passing a positional argument to 'system login', got nil")
	}
}
