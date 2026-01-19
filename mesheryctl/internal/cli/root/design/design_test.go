package design

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var (
	invalidFilePath = "/invalid/path/design.yaml"
)

func TestDesignCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "design apply",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "design.golden")},
			ExpectedContains: []string{"design applied", "deployed application meshapp"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "design.apply.save.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "design.apply.deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "design delete",
			Args:             []string{"delete", "-f", filepath.Join(fixturesDir, "design.golden")},
			ExpectedContains: []string{"deleted application meshapp"},
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "design.delete.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "design view",
			Args:             []string{"view", "kumatest"},
			ExpectedContains: []string{"name: kumatest", "id: 957fbc9b-708d-4396-84b8-e2ba37c1adcc"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=kumatest",
					Response:     "design.view.kuma.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "design view with ID",
			Args:             []string{"view", "4o7fbc9b-708d-4396-84b8-e2ba37c1adcc"},
			ExpectedContains: []string{"name: kumatest", "id: 957fbc9b-708d-4396-84b8-e2ba37c1adcc"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/4o7fbc9b-708d-4396-84b8-e2ba37c1adcc",
					Response:     "design.id.view.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=4o7fbc9b-708d-4396-84b8-e2ba37c1adcc",
					Response:     "design.id.view.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "design invalid view",
			Args:             []string{"view", "test-view"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=test-view",
					Response:     "design.view.invalid.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
			},
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  ErrDesignNotFound(),
		},
	}

	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

// reset other flags if needed
func resetVariables() {
	skipSave = false
	patternFile = ""
	file = ""
}
