package adapter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRemoveMesh(t *testing.T) {
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
	tests := []struct {
		Name             string
		Args             []string
		URLs             []utils.MockURL
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name: "Test Remove Linkerd ",
			Args: []string{"remove", "linkerd"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.linkerd.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name: "Test Remove Linkerd with Namespace",
			Args: []string{"remove", "linkerd", " --namespace", " linkerd-ns"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.linkerd-namespace.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name: "Test Remove Cilium ",
			Args: []string{"remove", "cilium"},
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.cilium.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			for _, url := range tc.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}
			utils.TokenFlag = tc.Token
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			cmd := AdapterCmd
			cmd.SetArgs(tc.Args)
			cmd.SetOutput(buff)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tc.ExpectError {
					// write it in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}
			// response being printed in console
			actualResponse := buff.String()

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("Mesh Deploy test Passed")
	}
	utils.StopMockery(t)
}
