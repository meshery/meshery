package design

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

const nonExistentID = "a12b3c4d-5e6f-4890-abcd-ef1234567890"

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

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given valid file when design delete then design is deleted",
			Args:             []string{"delete", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file as positional argument when design delete then design is deleted",
			Args:             []string{"delete", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid relative file as positional argument when design delete then design is deleted",
			Args:             []string{"delete", filepath.Join("fixtures", "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given invalid file path when design delete then error is thrown",
			Args:             []string{"delete", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "given mistyped file path positionally when server unreachable then file read error is thrown",
			Args:             []string{"delete", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "given nonexistent design ID when design delete then error is thrown",
			Args:             []string{"delete", nonExistentID},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?page=0&pagesize=100&populate=pattern_file",
					Response:     "delete-valid-id-list.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/" + nonExistentID,
					Response:     "delete.error.response.golden",
					ResponseCode: 404,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(nonExistentID),
		},
		{
			Name:             "given no arguments when design delete then error is thrown",
			Args:             []string{"delete"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide a design name, ID, or file path")),
		},
		{
			Name:             "given empty positional argument when design delete then error is thrown",
			Args:             []string{"delete", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty")),
		},
		{
			Name:             "given whitespace positional argument when design delete then error is thrown",
			Args:             []string{"delete", "   "},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty")),
		},
		{
			Name:             "given multiple positional arguments when design delete then error is thrown",
			Args:             []string{"delete", "a.yaml", "b.yaml"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("expected at most one positional argument, got 2")),
		},
		{
			Name:             "given both positional argument and -f flag when design delete then error is thrown",
			Args:             []string{"delete", filepath.Join(fixturesDir, "sampleDesign.golden"), "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a positional argument and the --file flag")),
		},
		{
			Name:             "given empty file flag when design delete then error is thrown",
			Args:             []string{"delete", "-f", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("file path cannot be empty")),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

func TestDeleteRejectsAmbiguousIDWithoutSendingDelete(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	utils.TokenFlag = utils.GetToken(t)
	utils.SetupMeshkitLoggerTesting(t, false)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)

	first := "dc761694-5e61-4dd4-83c9-bca45667ef70"
	second := "d7ccca72-adea-42e1-bb19-b28701db050d"
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/pattern?page=0&pagesize=100&populate=pattern_file",
		httpmock.NewStringResponder(200, fmt.Sprintf(`{"totalCount":2,"patterns":[{"id":%q},{"id":%q}]}`, first, second)))
	deleted := false
	for _, id := range []string{first, second} {
		httpmock.RegisterResponder("DELETE", testContext.BaseURL+"/api/pattern/"+id, func(_ *http.Request) (*http.Response, error) {
			deleted = true
			return httpmock.NewStringResponse(200, ""), nil
		})
	}
	DesignCmd.SetArgs([]string{"delete", "d"})
	err := DesignCmd.Execute()
	if meshkiterrors.GetCode(err) != utils.ErrAmbiguousIDCode || !strings.Contains(err.Error(), first) || !strings.Contains(err.Error(), second) {
		t.Fatalf("expected ambiguity with both IDs, got %v", err)
	}
	if deleted {
		t.Fatal("ambiguous ID must not send a DELETE request")
	}
}
