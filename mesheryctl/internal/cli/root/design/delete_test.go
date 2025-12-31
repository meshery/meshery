package design

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterr "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

const (
	invalidFilePath = "/invalid/path/design.yaml"
	nonExistentID   = "a12b3c4d-5e6f-4890-abcd-ef1234567890"
)

func TestDeleteCmd(t *testing.T) {
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
			Name:             "Delete Design",
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
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectError: false,
		},
		{
			Name:             "Delete design with invalid file path",
			Args:             []string{"delete", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "Delete design by ID with API error",
			Args:             []string{"delete", nonExistentID},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&page_size=10000",
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
			Token:          filepath.Join(fixturesDir, "token.golden"),
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError: func() error {
				body := "{\"error\": \"design not found\"}\n"
				innerErr := utils.ErrFailReqStatus(404, body)
				return ErrDeleteDesign(innerErr, nonExistentID)
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

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

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
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("Delete Design test Passed")
	}
	// stop mock server
	utils.StopMockery(t)
}
