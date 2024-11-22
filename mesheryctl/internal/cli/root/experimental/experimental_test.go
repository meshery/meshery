package experimental

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"flag"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

// Redirect stdin to simulate "Enter" key press
func simulateInput(input string) (*os.File, func()) {
	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()

	// Write the simulated input to the write end of the pipe
	_, _ = w.WriteString(input)

	// Set os.Stdin to the read end of the pipe
	originalStdin := os.Stdin
	os.Stdin = r

	// Return a cleanup function to restore the original os.Stdin
	return w, func() {
		os.Stdin = originalStdin
		r.Close()
		w.Close()
	}
}

func trimLastNLines(s string, n int) string {
	lines := strings.Split(s, "\n")
	if len(lines) <= n {
		return ""
	}
	return strings.Join(lines[:len(lines)-n], "\n")
}

func TestExperimentalList(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	//initialize mock server for handling requests
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
			Name:             "List registered relationships",
			Args:             []string{"relationship", "list"},
			URL:              testContext.BaseURL + "/api/meshmodels/relationships",
			Fixture:          "list.exp.relationship.api.response.golden",
			ExpectedResponse: "list.exp.relationship.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			utils.TokenFlag = tt.Token

			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Redirect stdin to simulate input
			simulateInput("\n")

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			ExpCmd.SetArgs(tt.Args)
			ExpCmd.SetOutput(rescueStdout)
			err := ExpCmd.Execute()
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

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			actualResponse := string(out)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			expectedResponse = trimLastNLines(expectedResponse, 2)

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Log("List Exp test Passed")
	}

	utils.StopMockery(t)
}
