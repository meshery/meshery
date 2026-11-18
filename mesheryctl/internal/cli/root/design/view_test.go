package design

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
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
			Name:             "given name provided when design view then the design with that name is displayed",
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
					Response:     "view.design.search.api.response.golden",
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
			Name:             "given a name that designs only contain when design view then error is thrown",
			Args:             []string{"view", "Untitled"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "view.design.api.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000&page=1",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=10000&populate=pattern_file&search=Untitled",
					Response:     "view.design.search.api.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("Untitled"),
		},
		{
			Name:             "given invalid output format when design view then error is returned before any request",
			Args:             []string{"view", "Untitled Design", "--output-format", "invalid"},
			ExpectedResponse: "",
			// No responders: a request, or a prompt for one of several designs
			// with the name, would surface as a different error.
			URLs:           []utils.MockURL{},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  display.ErrUnsupportedFormat("invalid"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

// TestDesignViewRejectsUnsupportedFormatBeforeAnyRequest pins the order: the
// output format is checked before the name lookup, the search and the prompt
// that several designs with one name would raise. The table above cannot pin
// it, because RunMesheryctlMultiURLTests leaves each case's responders
// registered for the next one, so an empty URL list does not prove that no
// request was served.
func TestDesignViewRejectsUnsupportedFormatBeforeAnyRequest(t *testing.T) {
	_ = utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(DesignCmd, t)
	defer resetVariables()
	utils.TokenFlag = utils.GetToken(t)

	// Drop the responders and call counters left by earlier tests.
	httpmock.Reset()

	DesignCmd.SetArgs([]string{"view", "Untitled Design", "--output-format", "invalid"})
	err := DesignCmd.Execute()
	if err == nil {
		t.Fatal("expected an error for --output-format invalid, got nil")
	}
	utils.AssertMeshkitErrorsEqual(t, err, display.ErrUnsupportedFormat("invalid"))
	if n := httpmock.GetTotalCallCount(); n != 0 {
		t.Fatalf("expected no requests before rejecting the output format, got %d", n)
	}
}
