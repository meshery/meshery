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
	"os"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCmdArgs(t *testing.T) {
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
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", updateUsageMsg)),
		},
		{
			Name:          "invalid with multiple unexpected positional arguments",
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", updateUsageMsg)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := updateCmd.Args(updateCmd, tt.Args)
			if tt.ExpectError {
				assert.Error(t, err)
				utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSupersedePinnedVersion(t *testing.T) {
	tests := []struct {
		name        string
		skipReset   bool
		answer      string
		wantProceed bool
		wantVersion string
	}{
		{
			// The Helm upgrade is skipped, so the cluster still runs the pinned version.
			name:        "kubernetes --skip-reset keeps the pinned version",
			skipReset:   true,
			answer:      "y\n",
			wantProceed: true,
			wantVersion: "v0.8.214",
		},
		{
			name:        "kubernetes update moves the pinned version to latest after confirmation",
			answer:      "y\n",
			wantProceed: true,
			wantVersion: "latest",
		},
		{
			name:        "kubernetes update keeps the pinned version when confirmation is declined",
			answer:      "n\n",
			wantProceed: false,
			wantVersion: "v0.8.214",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.SetupMeshkitLoggerTesting(t, false)
			prevSkipReset, prevSilent := utils.SkipResetFlag, utils.SilentFlag
			utils.SkipResetFlag, utils.SilentFlag = tt.skipReset, false
			t.Cleanup(func() { utils.SkipResetFlag, utils.SilentFlag = prevSkipReset, prevSilent })

			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			_, _ = w.WriteString(tt.answer)
			_ = w.Close()
			stdin := os.Stdin
			os.Stdin = r
			t.Cleanup(func() {
				os.Stdin = stdin
				_ = r.Close()
			})

			currCtx := &config.Context{Platform: platformKubernetes, Channel: "stable", Version: "v0.8.214"}
			proceed, err := supersedePinnedVersion(currCtx)

			assert.NoError(t, err)
			assert.Equal(t, tt.wantProceed, proceed)
			assert.Equal(t, tt.wantVersion, currCtx.GetVersion())
		})
	}
}
