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

import "testing"

// TestLoginCmdArgs pins the argument validator of `system login`:
// the command takes no positional arguments, so any supplied argument
// must be rejected before the login handler runs.
func TestLoginCmdArgs(t *testing.T) {
	if err := loginCmd.Args(loginCmd, []string{"unexpected-arg"}); err == nil {
		t.Error("expected an error when passing a positional argument to 'system login', got nil")
	}
	if err := loginCmd.Args(loginCmd, []string{"unexpected-arg-1", "unexpected-arg-2"}); err == nil {
		t.Error("expected an error when passing multiple positional arguments to 'system login', got nil")
	}
	if err := loginCmd.Args(loginCmd, []string{}); err != nil {
		t.Errorf("expected no error when passing no arguments to 'system login', got %v", err)
	}
}

// TestLoginCmdFlags verifies that the provider flag remains registered
// and supported by loginCmd and is not affected by argument validation changes.
func TestLoginCmdFlags(t *testing.T) {
	flag := loginCmd.PersistentFlags().Lookup("provider")
	if flag == nil {
		t.Fatal("expected 'provider' flag to be registered on loginCmd, got nil")
	}
	if flag.Shorthand != "p" {
		t.Errorf("expected shorthand for 'provider' flag to be 'p', got %q", flag.Shorthand)
	}
}
