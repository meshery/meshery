package model

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestImportModel(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	testdataDir := filepath.Join(currDir, "testdata")
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Import model providing invalid template file path",
			Args:             []string{"import", "-f", "invalidTemplatePath"},
			ExpectedResponse: "import.invalid.file.output.golden",
			Fixture:          "import.invalid.file.api.response.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/register",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Import model providing empty template file",
			Args:             []string{"import", "-f", filepath.Join(fixturesDir, "import.empty.file.api.response.golden")},
			ExpectedResponse: "import.empty.file.output.golden",
			Fixture:          "import.empty.file.api.response.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/register",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
			utils.TokenFlag = tt.Token
			httpmock.RegisterResponder("POST", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			ModelCmd.SetOut(b)
			ModelCmd.SetArgs(tt.Args)
			err := ModelCmd.Execute()
			if err != nil {
				if tt.ExpectError {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}

			output := b.String()
			actualResponse := output

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("model import test Passed")
	}

	utils.StopMockery(t)
}
