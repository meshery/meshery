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
	"github.com/stretchr/testify/assert"
)

func TestDeleteCmdArgs(t *testing.T) {
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
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("delete does not accept positional arguments\n\n%s", deleteUsageMsg)),
		},
		{
			Name:          "invalid with multiple unexpected positional arguments",
			Args:          []string{"arg1", "arg2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("delete does not accept positional arguments\n\n%s", deleteUsageMsg)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := deleteCmd.Args(deleteCmd, tt.Args)
			if tt.ExpectError {
				assert.Error(t, err)
				utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
