package design

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var (
	update                 = flag.Bool("update", false, "update golden files")
	invalidFilePath        = "/invalid/path/design.yaml"
	validDesignSourceTypes = []string{"Helm Chart", "Kubernetes Manifest", "Docker Compose", "Meshery Design"}
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
			ExpectedResponse: "design.apply.output.golden",
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
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "design delete",
			Args:             []string{"delete", "-f", filepath.Join(fixturesDir, "design.golden")},
			ExpectedResponse: "design.delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "design.delete.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "design view",
			Args:             []string{"view", "kumatest"},
			ExpectedResponse: "design.view.kuma.output.golden",
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
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "design view with ID",
			Args:             []string{"view", "4o7fbc9b-708d-4396-84b8-e2ba37c1adcc"},
			ExpectedResponse: "design.id.view.output.golden",
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
			Token:       filepath.Join(fixturesDir, "token.golden"),
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
			Token:          filepath.Join(fixturesDir, "token.golden"),
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
