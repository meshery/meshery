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

package workspaces

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestDeleteWorkspace(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	workspaceId := "0dd47d1a-d1c9-47dc-897c-40bf4a71d96b"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "Delete workspace without arguments",
			Args:             []string{"delete"},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              "/api/workspaces",
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ Workspace ID ] isn't specified\n\nUsage: mesheryctl workspace delete [workspaceId]\nRun 'mesheryctl workspace delete --help' to see detailed help message")),
		},
		{
			Name:             "Delete workspace successfully",
			Args:             []string{"delete", workspaceId},
			HttpMethod:       "DELETE",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/api/workspaces/%s", workspaceId),
			Fixture:          "delete.workspace.api.response.golden",
			ExpectedResponse: "delete.workspace.success.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}
