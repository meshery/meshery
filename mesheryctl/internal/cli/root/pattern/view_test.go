package pattern

import (
	"bytes"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func setupDesign(t *testing.T, testContext *utils.TestHelper, name string, shouldExist bool) {
	design := `{ "name": "` + name + `" }` // Replace with actual design JSON
	url := testContext.BaseURL + "/api/pattern"

	if shouldExist {
		// Mock uploading the design
		httpmock.RegisterResponder("POST", url,
			httpmock.NewStringResponder(201, design))
		
		// Upload design
		_, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(design)))
		if err != nil {
			t.Fatalf("Failed to upload design: %v", err)
		}
	} else {
		// Mock deleting the design
		httpmock.RegisterResponder("DELETE", url+"?name="+name,
			httpmock.NewStringResponder(200, ""))

		// Delete design
		req, err := http.NewRequest("DELETE", url+"?name="+name, nil)
		if err != nil {
			t.Fatalf("Failed to create delete request: %v", err)
		}

		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			t.Fatalf("Failed to delete design: %v", err)
		}
	}
}

func TestPatternView(t *testing.T) {
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

	// Ensure design exists
	setupDesign(t, testContext, "Untitled Design", true)

	// Ensure design does not exist
	setupDesign(t, testContext, "NonExistentDesign", false)

	// test scenarios for fetching data
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
			Name:             "Fetch Pattern View",
			Args:             []string{"view", "design"},
			ExpectedResponse: "view.pattern.output.golden",
			Fixture:          "view.pattern.api.response.golden",
			URL:              testContext.BaseURL + "/api/pattern",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Specific Pattern by Name",
			Args:             []string{"view", "design"},
			ExpectedResponse: "view.pattern.output.golden",
			Fixture:          "view.pattern.api.response.golden",
			URL:              testContext.BaseURL + "/api/pattern?name=Untitled Design",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "View Non-Existent Pattern by Name",
			Args:             []string{"view", "design"},
			ExpectedResponse: "view.pattern.output.golden",
			Fixture:          "view.pattern.api.response.golden",
			URL:              testContext.BaseURL + "/api/pattern?name=NonExistentDesign",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// View api response from golden files
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			// set token
			utils.TokenFlag = tt.Token

			// mock response
			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			PatternCmd.SetOutput(b)
			PatternCmd.SetArgs(tt.Args)
			err := PatternCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
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
			output := b.String()
			actualResponse := output

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)

			if tt.Name == "View Specific Pattern by Name" {
				t.Logf("Output for '%s': %s", tt.Name, actualResponse)
			}
		})
		t.Log("View Pattern test Passed")
	}

	// stop mock server
	utils.StopMockery(t)
}