package experimental

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestExperimentalList(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	//initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectedError    bool
	}{
		{
			Name:             "List registered relationships",
			Args:             []string{"relationship", "list"},
			URL:              testContext.BaseURL + "/api/meshmodels/relationships",
			Fixture:          "list.exp.api.response.golden",
			ExpectedResponse: "list.exp.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectedError:    false,
		},
	}

	// run tests
	/* for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			utils.TokenFlag = tt.Token

			httpmock.RegisterResponder("GET", tt.URL, httpmock.NewStringResponder(200, apiResponse))
		})
	} */
}
