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

	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestRestartCmdArgs(t *testing.T) {
	// restart does not take any arguments
	err := restartCmd.Args(restartCmd, []string{})
	assert.NoError(t, err)

	err = restartCmd.Args(restartCmd, []string{"extra-arg"})
	assert.Error(t, err)
}

func TestRestartCmdFlags(t *testing.T) {
	skipUpdateFlag := restartCmd.Flags().Lookup("skip-update")
	assert.NotNil(t, skipUpdateFlag)
	assert.Equal(t, "false", skipUpdateFlag.DefValue)

	providerFlag := restartCmd.Flags().Lookup("provider")
	assert.NotNil(t, providerFlag)
	assert.Equal(t, "", providerFlag.DefValue)
}

func TestErrRestartMeshery(t *testing.T) {
	cause := fmt.Errorf("underlying failure")
	err := ErrRestartMeshery(cause)
	assert.Error(t, err)
	assert.Equal(t, ErrRestartMesheryCode, errors.GetCode(err))
	assert.ErrorContains(t, err, cause.Error())
	assert.Contains(t, err.Error(), "Error restarting Meshery")
}
