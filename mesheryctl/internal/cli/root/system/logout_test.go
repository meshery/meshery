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

// TestLogoutCmdArgs pins the argument validator of `system logout`:
// the command takes no positional arguments, so any supplied argument
// must be rejected before the logout handler runs.
func TestLogoutCmdArgs(t *testing.T) {
	if err := logoutCmd.Args(logoutCmd, []string{"unexpected-arg"}); err == nil {
		t.Error("expected an error when passing a positional argument to 'system logout', got nil")
	}
	if err := logoutCmd.Args(logoutCmd, []string{}); err != nil {
		t.Errorf("expected no error when passing no arguments to 'system logout', got %v", err)
	}
}
