package design

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
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
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "delete.idList.response.golden",
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
		{
			Name:             "given ID prefix shared by two designs when design delete then ambiguity error is thrown and nothing is deleted",
			Args:             []string{"delete", "a12"},
			ExpectedResponse: "",
			// DELETE responders are registered so that any delete request would
			// succeed and fail this test, which expects the ambiguity error.
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "delete.ambiguousIdList.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/a12b3c4d-5e6f-4890-abcd-ef1234567890",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/a12f9e8d-7c6b-4a59-8d4c-3b2a19087654",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  utils.ErrAmbiguousID("a12", []string{"a12b3c4d-5e6f-4890-abcd-ef1234567890", "a12f9e8d-7c6b-4a59-8d4c-3b2a19087654"}),
		},
		{
			Name:             "given generic remote URL with -f flag when design delete then design is deleted",
			Args:             []string{"delete", "-f", "https://example.com/designs/sample.yaml"},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          "https://example.com/designs/sample.yaml",
					Response:     "sampleDesign.golden",
					ResponseCode: 200,
				},
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
			Name:             "given generic remote URL as positional argument when design delete then design is deleted",
			Args:             []string{"delete", "https://example.com/designs/sample.yaml"},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          "https://example.com/designs/sample.yaml",
					Response:     "sampleDesign.golden",
					ResponseCode: 200,
				},
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
			Name:             "given github URL with -f flag when design delete then design is deleted",
			Args:             []string{"delete", "-f", "https://github.com/meshery/meshery/blob/master/mesheryctl/internal/cli/root/design/fixtures/sampleDesign.golden"},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          "https://raw.githubusercontent.com/meshery/meshery/master/mesheryctl/internal/cli/root/design/fixtures/sampleDesign.golden",
					Response:     "sampleDesign.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

// TestDeleteWithGenericURL verifies that mesheryctl design delete -f <generic-url>
// downloads the design via fetchDesignFile and sends the undeploy payload to Meshery server.
func TestDeleteWithGenericURL(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	const genericURL = "https://example.com/designs/sample.yaml"
	const expectedContent = "name: sampleDesign\nservices:\n  app:\n    type: Application\n"

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "delete with -f flag",
			args: []string{"delete", "-f", genericURL},
		},
		{
			name: "delete with positional URL",
			args: []string{"delete", genericURL},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetVariables()
			defer utils.ResetCommandFlags(DesignCmd, t)

			// ResetCommandFlags restores the --token flag to its default after
			// each run, so point it at the fixture token for every subtest.
			utils.TokenFlag = utils.GetToken(t)

			var getCalled, deleteCalled bool
			var undeployPayload struct {
				PatternFile string `json:"patternFile"`
			}

			httpmock.Reset()
			httpmock.RegisterResponder("GET", genericURL, func(req *http.Request) (*http.Response, error) {
				getCalled = true
				return httpmock.NewStringResponse(http.StatusOK, expectedContent), nil
			})
			httpmock.RegisterResponder("DELETE", testContext.BaseURL+"/api/pattern/deploy",
				func(req *http.Request) (*http.Response, error) {
					deleteCalled = true
					if req.Body != nil {
						defer req.Body.Close()
						bodyBytes, err := io.ReadAll(req.Body)
						assert.NoError(t, err)
						assert.NoError(t, json.Unmarshal(bodyBytes, &undeployPayload))
					}
					return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
				},
			)

			DesignCmd.SetArgs(tt.args)
			err := DesignCmd.Execute()
			assert.NoError(t, err)
			assert.True(t, getCalled, "expected GET request to download generic URL via fetchDesignFile")
			assert.True(t, deleteCalled, "expected DELETE request to /api/pattern/deploy")
			assert.Equal(t, expectedContent, undeployPayload.PatternFile)
		})
	}
}
