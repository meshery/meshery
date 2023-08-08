package mesh

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestDeployCmd(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)
	// initialize mock server for handling requests
	utils.StartMockery(t)
	// create a test helper
	testcontext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	//utils.SetupCustomContextEnv(t, currDir+"/fixtures/.meshery/TestContext.yaml")
	fixturesDir := filepath.Join(currDir, "fixtures")

	testcases := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		URls             []utils.MockURL
		Token            string
		ExpectedErr      bool
	}{
		{
			Name:             "Deploy Istio",
			Args:             []string{"deploy", "istio"},
			ExpectedResponse: "deploy.istio.output.golden",
			URls: []utils.MockURL{
				{
					Method:   "GET",
					URL:      testcontext.BaseURL + "/api/system/sync",
					Response: "deploy.istio.api.golden",
				},
			},
			Token:       filepath.Join(fixturesDir, "token.golden"),
			ExpectedErr: false,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, url := range tc.URls {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()

				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}
			utils.TokenFlag = tc.Token
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			b := utils.SetupMeshkitLoggerTesting(t, false)
			MeshCmd.SetOutput(b)
			MeshCmd.SetArgs(tc.Args)
			err := MeshCmd.Execute()
			if err != nil {
				if tc.ExpectedErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					utils.Equals(t, expectedResponse, err.Error())
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
	}
	// stop mock server
	utils.StopMockery(t)
}
