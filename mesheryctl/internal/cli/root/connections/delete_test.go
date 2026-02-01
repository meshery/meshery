package connections

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

const ConnectionID = "11111111-1111-1111-1111-111111111111"

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
			Name:          "Delete connection without ID",
			Args:          []string{"delete"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidArgument(errors.New(errMsg)),
		},
		{
			Name:          "Delete connection with invalid UUID",
			Args:          []string{"delete", "invalid-id"},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", "invalid-id")),
		},
		{
			Name: "Delete connection not found",
			Args: []string{"delete", ConnectionID},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/integrations/connections/" + ConnectionID,
					Response:     "delete.connection.notfound.response.golden",
					ResponseCode: 500,
				},
			},
			ExpectError:   true,
			ExpectedError: errConnectionNotFound(fmt.Errorf("No connection with id %q found", ConnectionID)),
		},
		{
			Name: "Delete connection successfully",
			Args: []string{"delete", ConnectionID},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/integrations/connections/" + ConnectionID,
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
