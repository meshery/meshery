package filter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testContext := utils.NewTestHelper(t)
	base := testContext.BaseURL

	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name: "Delete Kuma-Test",
			Args: []string{"delete", "c0c6035a-b1b9-412d-aab2-4ed1f1d51f84", "Kuma-Test"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?page_size=10000", base), Response: "filter.name.view.api.response.golden", ResponseCode: 200},
				{Method: "DELETE", URL: fmt.Sprintf("%s/api/filter/c0c6035a-b1b9-412d-aab2-4ed1f1d51f84", base), Response: "delete.kuma.api.response.golden", ResponseCode: 200},
			},
			ExpectedContains: []string{"Filter c0c6035a-b1b9-412d-aab2-4ed1f1d51f84 deleted"},
			ExpectError:      false,
		},
		{
			Name: "Delete RolloutAndIstio",
			Args: []string{"delete", "d0e09134-acb6-4c71-b051-3d5611653f70", "RolloutAndIstio"},
			URLs: []utils.MockURL{
				{Method: "GET", URL: fmt.Sprintf("%s/api/filter?page_size=10000", base), Response: "filter.name.view.api.response.golden", ResponseCode: 200},
				{Method: "DELETE", URL: fmt.Sprintf("%s/api/filter/d0e09134-acb6-4c71-b051-3d5611653f70", base), Response: "delete.rollout.api.response.golden", ResponseCode: 200},
			},
			ExpectedContains: []string{"Filter d0e09134-acb6-4c71-b051-3d5611653f70 deleted"},
			ExpectError:      false,
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, tests, currDir, "filter", func() {})
}
