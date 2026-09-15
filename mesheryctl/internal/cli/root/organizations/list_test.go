package organizations

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"gopkg.in/yaml.v3"
)

func TestListOrganizations(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List organizations whithout providing organization ID",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10", organizationsApiPath),
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "list.organization.golden",
			ExpectError:      false,
		},
		{
			Name:          "List organizations with an invalid --output-format value returns an error",
			Args:          []string{"list", "--output-format", "foo"},
			Fixture:       "list.organization.response.golden",
			ExpectError:   true,
			ExpectedError: display.ErrInvalidOutputFormat("foo"),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, OrgCmd, tests, currentDirectory, "organization")
}

// TestListOrganizationsStructuredOutput exercises `organization list
// --output-format json|yaml` directly rather than through the golden-file
// harness: the response envelope's custom MarshalJSON (see
// server/models/organization.go) makes an exact byte-for-byte golden fragile
// to hand-maintain, so this instead round-trips stdout back into
// models.OrganizationsPage and asserts on the decoded fields.
func TestListOrganizationsStructuredOutput(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	fixturesDir := filepath.Join(currentDirectory, "fixtures")

	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(OrgCmd, t)

	apiResponse := utils.NewGoldenFile(t, "list.organization.response.golden", fixturesDir).Load()
	utils.TokenFlag = utils.GetToken(t)

	url := fmt.Sprintf("%s/%s?page=0&pagesize=10", testContext.BaseURL, organizationsApiPath)
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, apiResponse))

	for _, format := range []string{"json", "yaml"} {
		t.Run(format, func(t *testing.T) {
			defer utils.ResetCommandFlags(OrgCmd, t)

			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			defer func() { os.Stdout = originalStdout }()

			_ = utils.SetupMeshkitLoggerTesting(t, false)
			OrgCmd.SetArgs([]string{"list", "--output-format", format})
			OrgCmd.SetOut(w)
			err := OrgCmd.Execute()
			_ = w.Close()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			out, readErr := io.ReadAll(r)
			if readErr != nil {
				t.Fatalf("failed to read command output: %v", readErr)
			}

			var page models.OrganizationsPage
			switch format {
			case "json":
				if err := json.Unmarshal(out, &page); err != nil {
					t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, out)
				}
			case "yaml":
				if err := yaml.Unmarshal(out, &page); err != nil {
					t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, out)
				}
			}

			if page.TotalCount != 1 {
				t.Fatalf("TotalCount = %d, want 1", page.TotalCount)
			}
			if len(page.Organizations) != 1 {
				t.Fatalf("len(Organizations) = %d, want 1", len(page.Organizations))
			}
			if page.Organizations[0].Name != "name" {
				t.Fatalf("Organizations[0].Name = %q, want %q", page.Organizations[0].Name, "name")
			}
		})
	}
}
