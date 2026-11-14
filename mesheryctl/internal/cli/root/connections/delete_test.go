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

package connections

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

const connectionId = "11111111-1111-1111-1111-111111111111"

func TestConnectionDeleteCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:          "given no connectionID provided when running mesheryctl connection delete then an error message is displayed",
			Args:          []string{"delete"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errNoArgMsg, deleteUsageMsg)),
		},
		{
			Name:          "given multiple arguments provided when running mesheryctl connection delete then an error message is displayed",
			Args:          []string{"delete", "id1", "id2"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errMultiArgMsg, deleteUsageMsg)),
		},
		{
			Name:          "given an invalid connectionID provided when running mesheryctl connection delete invalid-id then an error message is displayed",
			Args:          []string{"delete", "invalid-id"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", "invalid-id")),
		},
		{
			Name: "given a non-existent connectionID provided when running mesheryctl connection delete then an error message is displayed",
			Args: []string{"delete", connectionId},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/integrations/connections/" + connectionId,
					Response:     "delete.connection.notfound.response.golden",
					ResponseCode: 404,
				},
			},
			IsOutputGolden:   true,
			ExpectedResponse: "delete.connection.notfound.output.golden",
		},
		{
			Name: "given a valid connectionID provided when running mesheryctl connection delete then the connection is deleted",
			Args: []string{"delete", connectionId},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/integrations/connections/" + connectionId,
					Response:     "delete.connection.success.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.connection.success.output.golden",
			IsOutputGolden:   true,
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, ConnectionsCmd, tests, currDir, "connection", resetVariables)
}

func resetVariables() {}
