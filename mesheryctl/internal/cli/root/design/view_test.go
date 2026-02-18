package design

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDesignView(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given name provider when view command is executed then design view is displayed",
			Args:             []string{"view", "desgin"},
			ExpectedResponse: "view.design.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=desgin",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       "",
			ExpectError: false,
		},
		{
			Name:             "given id provider when view command is executed then design view is displayed",
			Args:             []string{"view", "3817ec9a-1d83-4f6f-9154-0fd4408ba9f0"},
			ExpectedResponse: "view.design.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/3817ec9a-1d83-4f6f-9154-0fd4408ba9f0",
					Response:     "view.design.by.id.api.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       "",
			ExpectError: false,
		},
		{
			Name:             "given no name or ID specified when view command is executed then error is thrown",
			Args:             []string{"view"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    ErrDesignNameOrIDNotSpecified(),
		},
		{
			Name:             "given nonexistent design when view command is executed then error is thrown",
			Args:             []string{"view", "nonexistent-design"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=nonexistent-design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          "",
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
