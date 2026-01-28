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
	fixturesDir := filepath.Join(currDir, "fixtures")

	listURL := testContext.BaseURL + "/api/filter?page_size=10000"
	listResponse := `{"filters":[{"id":"c0c6035a-b1b9-412d-aab2-4ed1f1d51f84","name":"Kuma-Test"},{"id":"d0e09134-acb6-4c71-b051-3d5611653f70","name":"RolloutAndIstio"}]}`
	// Mock the filter list endpoint (used by ValidId/GetID)
	filterListResponse := utils.NewGoldenFile(t, "delete.filter.list.api.response.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/filter?page_size=10000",
		httpmock.NewStringResponder(200, filterListResponse))

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
	for _, tt := range testcase {
		t.Run(tt.Name, func(t *testing.T) {
			httpmock.RegisterResponder("GET", listURL,
				httpmock.NewStringResponder(200, listResponse))

			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
			// set token
			utils.TokenFlag = tt.Token
			// mock response
			httpmock.RegisterResponder(tt.Method, tt.URL,
				httpmock.NewStringResponder(tt.ResponseCode, apiResponse))
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, FilterCmd, testcase, currDir, "filter", func() {})
}
