package design

import (
	"encoding/json"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	meshkiterr "github.com/meshery/meshkit/errors"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestApplyCmd(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// initialize mock server for handling requests
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

	// test scenrios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		URLs             []utils.MockURL
		Token            string
		ExpectError      bool
		IsOutputGolden   bool  `default:"true"`
		ExpectedError    error `default:"nil"`
	}{
		{
			Name:             "Apply Designs",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Apply Designs with --skip-save",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--skip-save"},
			ExpectedResponse: "apply.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "apply.designDeploy.response.golden",
					ResponseCode: 200,
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Apply design with invalid file path",
			Args:             []string{"apply", "-f", "/invalid/path/design.yaml"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(errors.Errorf(errInvalidPathMsg, "/invalid/path/design.yaml")),
		},
		{
			Name:             "Apply design with invalid server response",
			Args:             []string{"apply", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.invalidJSON.response.golden",
					ResponseCode: 200,
				},
			},
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError: func() error {
				// Replicate the exact JSON unmarshal error
				var response []*models.MesheryPattern
				innerErr := json.Unmarshal([]byte(`{ "patterns": [ { "id": "123", "name": "incomplete-json"`), &response)

				return utils.ErrUnmarshal(innerErr)
			}(),
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for _, url := range tt.URLs {
				// View api response from golden files
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()

				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}

			// set token
			utils.TokenFlag = tt.Token

			testdataDir := filepath.Join(currDir, "testdata")

			// Skip golden file creation for error tests that use ExpectedError instead
			var golden *utils.GoldenFile
			if tt.ExpectedResponse != "" {
				golden = utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			}

			// setting up log to grab logs
			b := utils.SetupMeshkitLoggerTesting(t, false)
			DesignCmd.SetOut(b)
			DesignCmd.SetArgs(tt.Args)
			err := DesignCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					if tt.IsOutputGolden {

						// write it in file
						if *update {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						utils.Equals(t, expectedResponse, err.Error())
						resetVariables()
						return
					}
					assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(tt.ExpectedError), "error type mismatch")
					assert.Equal(t, meshkiterr.GetCode(err), meshkiterr.GetCode(tt.ExpectedError), "error code mismatch")
					assert.Equal(t, meshkiterr.GetLDescription(err), meshkiterr.GetLDescription(tt.ExpectedError), "long description mismatch")
					resetVariables()
					return

				}
				t.Error(err)
			}

			// response being printed in console
			actualResponse := b.String()

			// write it in file
			if *update {
				if golden != nil {
					golden.Write(actualResponse)
				}
			}
			var expectedResponse string
			expectedResponse = golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("Apply Design Test Passed")
	}

	// stop mock server
	utils.StopMockery(t)
}

// reset other flags if needed
func resetVariables() {
	skipSave = false
	patternFile = ""
	file = ""
}
