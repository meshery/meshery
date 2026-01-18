package relationships

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestSearch_WithoutFlags(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		ExpectedContains []string
		ExpectError      bool
	}{
		{
			Name:             "Search with missing arguments",
			Args:             []string{"search"},
			ExpectedContains: []string{"[--kind, --subtype or --type or --model] and [query-text] are required"},
			ExpectError:      true,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()
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
			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}
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
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Search registered relationships matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.matching.result.golden",
			ExpectedContains: []string{"Total number of relationships: 1", "KIND", "Mount"},
			ExpectError:      false,
		},
		{
			Name:             "Search registered relationships no matching result(s) found",
			Args:             []string{"search", "--model", "kubernetes"},
			URL:              "/api/meshmodels/models/kubernetes/relationships",
			Fixture:          "search.relationship.api.response.no.matching.result.golden",
			ExpectedContains: []string{"No relationships found"},
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
