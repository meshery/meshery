package filter

import (
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
					URL:          testContext.BaseURL + "/api/filter?page_size=10000",
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
					URL:          testContext.BaseURL + "/api/filter?page_size=10000",
					Response:     "filter.list.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "delete.rollout.output.golden",
			ExpectError:      false,
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, testcase, currDir, "filter", func() {})
}
