package filter

import (
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	testcase := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "Delete Kuma-Test",
			Args: []string{"delete", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/filter/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84",
					Response:     "delete.kuma.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?pagesize=10000",
					Response:     "filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.kuma.output.golden",
			ExpectError:      false,
		}, {
			Name: "Delete RolloutAndIstio",
			Args: []string{"delete", "RolloutAndIstio"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/filter/d0e09134-acb6-4c71-b051-3d5611653f70",
					Response:     "delete.rollout.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?pagesize=10000",
					Response:     "filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.rollout.output.golden",
			ExpectError:      false,
		}, {
			Name: "Delete exact name that is also a prefix of another name",
			Args: []string{"delete", "envoy"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/filter/1b4e28ba-2fa1-4d3b-a3f5-ef19b5a7633b",
					Response:     "delete.envoy.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?pagesize=10000",
					Response:     "delete.prefix.filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.envoy.output.golden",
			ExpectError:      false,
		}, {
			Name: "Delete with a name prefix shared by more than one filter",
			Args: []string{"delete", "env"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/filter?pagesize=10000",
					Response:     "delete.prefix.filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:   true,
			ExpectedError: utils.ErrInvalidNameOrID(errors.New(`"env" matches more than one filter (envoy, envoy-wasm); use the full name or the ID`)),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, testcase, currDir, "filter", func() {})
}
