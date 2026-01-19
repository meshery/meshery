package relationships

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestView(t *testing.T) {
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
		ExpectedContains []string
		ExpectError      bool
	}{
		{
			Name:             "View relationship without model name",
			Args:             []string{"view"},
			URL:              testContext.BaseURL + "/api/meshmodels/models/kubernetes/relationships?pagesize=all",
			Fixture:          "",
			ExpectedContains: []string{"[model-name] isn't specified"},
			Token:            utils.GetToken(t),
			ExpectError:      true,
		},
		{
			Name:             "View registered relationship",
			Args:             []string{"view", "kubernetes"},
			URL:              testContext.BaseURL + "/api/meshmodels/models/kubernetes/relationships?pagesize=all",
			Fixture:          "view.relationship.api.response.golden",
			Token:            utils.GetToken(t),
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Fixture != "" {
				apiResponse := utils.ReadTestFixture(t, fixturesDir, tt.Fixture)

				utils.TokenFlag = tt.Token

				httpmock.RegisterResponder("GET", tt.URL,
					httpmock.NewStringResponder(200, apiResponse))
			}

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
			_ = w.Close()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					for _, s := range tt.ExpectedContains {
						if !strings.Contains(err.Error(), s) {
							t.Fatalf("expected error to contain %q, got %q", s, err.Error())
						}
					}
					return
				}
				t.Fatal(err)
			}

			// Even when the command succeeds, it may print nothing.
			// Consume and discard stdout to avoid blocking.
			_, _ = io.Copy(io.Discard, r)
		})
		t.Log("View experimental relationship test passed")
	}

	utils.StopMockery(t)
}
