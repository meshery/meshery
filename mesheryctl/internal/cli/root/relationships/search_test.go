package relationships

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearch_WithoutFlags(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Search with missing arguments",
			Args:             []string{"search"},
			ExpectedResponse: "search.missing.args.output.golden",
			ExpectError:      true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints with proper cleanup
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Ensure stdout is always restored
			defer func() {
				os.Stdout = originalStdout
			}()

			_ = utils.SetupMeshkitLoggerTesting(t, false)
			RelationshipCmd.SetArgs(tt.Args)
			RelationshipCmd.SetOut(originalStdout)
			err := RelationshipCmd.Execute()

			// Close write end before reading
			w.Close()

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

			out, _ := io.ReadAll(r)
			actualResponse := string(out)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Log("Search experimental relationship test passed")
	}

	utils.StopMockery(t)
}

func TestSearch_WithFlags(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Search registered relationships matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.matching.result.golden",
			ExpectedResponse: "search.relationship.output.matching.result.golden",
			ExpectError:      false,
		},
		{
			Name:             "Search registered relationships no matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.no.matching.result.golden",
			ExpectedResponse: "search.relationship.output.no.matching.result.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
