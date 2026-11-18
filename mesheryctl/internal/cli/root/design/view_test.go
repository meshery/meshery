package design

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDesignView(t *testing.T) {
	origPrompt := runDesignSelectPrompt
	t.Cleanup(func() { runDesignSelectPrompt = origPrompt })
	runDesignSelectPrompt = func(_ string, _ []string) (int, error) { return 0, nil }

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
			Name:             "given name provided when design view then design is displayed",
			Args:             []string{"view", "Untitled Design"},
			ExpectedResponse: "view.design.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// GetID paginates: the page-0 fixture reports
					// totalCount 318 with 10 items, so a second page is
					// requested. The empty page lets the lookup conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000&page=1",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=10000&populate=pattern_file&search=Untitled+Design",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// Search paginates: the page-0 fixture reports totalCount 318
					// with 10 items, so a second page is requested. The empty page
					// lets the search conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=1&pagesize=10000&populate=pattern_file&search=Untitled+Design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given id provided when design view then design is displayed",
			Args:             []string{"view", "3817ec9a-1d83-4f6f-9154-0fd4408ba9f0"},
			ExpectedResponse: "view.design.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// GetID paginates: the page-0 fixture reports
					// totalCount 318 with 10 items, so a second page is
					// requested. The empty page lets the lookup conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000&page=1",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/3817ec9a-1d83-4f6f-9154-0fd4408ba9f0",
					Response:     "view.design.by.id.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given no name or ID specified when design view then error is thrown",
			Args:             []string{"view"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    ErrDesignNameOrIDNotSpecified(),
		},
		{
			Name:             "given nonexistent design provided when design view then error is thrown",
			Args:             []string{"view", "nonexistent-design"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// GetID paginates: the page-0 fixture reports
					// totalCount 318 with 10 items, so a second page is
					// requested. The empty page lets the lookup conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000&page=1",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=10000&populate=pattern_file&search=nonexistent-design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("nonexistent-design"),
		},
		{
			Name:             "given invalid output format when design view then error is returned",
			Args:             []string{"view", "Untitled Design", "--output-format", "invalid"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// GetID paginates: the page-0 fixture reports
					// totalCount 318 with 10 items, so a second page is
					// requested. The empty page lets the lookup conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000&page=1",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=10000&populate=pattern_file&search=Untitled+Design",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					// Search paginates: the page-0 fixture reports totalCount 318
					// with 10 items, so a second page is requested. The empty page
					// lets the search conclude.
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=1&pagesize=10000&populate=pattern_file&search=Untitled+Design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  display.ErrUnsupportedFormat("invalid"),
		},
		{
			Name:             "given substring only match when design view then error is returned",
			Args:             []string{"view", "nginx"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=10000&populate=pattern_file&search=nginx",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=1&pagesize=10000&populate=pattern_file&search=nginx",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("nginx"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}
