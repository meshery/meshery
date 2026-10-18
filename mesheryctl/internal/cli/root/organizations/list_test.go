package organizations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta2/organization"
	"gopkg.in/yaml.v3"
)

func TestListOrganizations(t *testing.T) {
	mesheryctlflags.InitValidators(OrgCmd)

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
			ExpectedError: utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --output-format 'foo': valid values are json yaml")),
		},
		{
			Name:          "List organizations with --count and --output-format together returns an error",
			Args:          []string{"list", "--count", "--output-format", "json"},
			Fixture:       "list.organization.response.golden",
			ExpectError:   true,
			ExpectedError: ErrCountWithOutputFormat(),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, OrgCmd, tests, currentDirectory, "organization")
}

// TestListOrganizationsStructuredOutput exercises `organization list
// --output-format json|yaml` directly rather than through the golden-file
// harness, round-tripping the command's output back into the schema's
// organization.OrganizationsPage and asserting on the raw serialized keys.
//
// Each case registers the mock responder at the exact URL the command should
// build for the given --page/--pagesize. If the 1-based-to-0-based page
// conversion regresses, the request misses this URL, httpmock returns an
// error, and the command fails. (The decoded response's Page/PageSize are not
// asserted against the request: the fixture is a fixed canned response.)
func TestListOrganizationsStructuredOutput(t *testing.T) {
	mesheryctlflags.InitValidators(OrgCmd)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	fixturesDir := filepath.Join(currentDirectory, "fixtures")

	cases := []struct {
		name         string
		extraArgs    []string
		expectedPage int
		expectedSize int
	}{
		{
			name:         "default page and pagesize",
			expectedPage: 0,
			expectedSize: 10,
		},
		{
			name:         "custom page and pagesize",
			extraArgs:    []string{"--page", "2", "--pagesize", "3"},
			expectedPage: 1,
			expectedSize: 3,
		},
	}

	for _, tc := range cases {
		for _, format := range []string{"json", "yaml"} {
			t.Run(tc.name+"/"+format, func(t *testing.T) {
				testContext := utils.InitTestEnvironment(t)
				defer utils.StopMockery(t)
				defer utils.ResetCommandFlags(OrgCmd, t)
				_ = utils.SetupMeshkitLoggerTesting(t, false)

				apiResponse := utils.NewGoldenFile(t, "list.organization.response.golden", fixturesDir).Load()
				utils.TokenFlag = utils.GetToken(t)

				url := fmt.Sprintf("%s/%s?page=%d&pagesize=%d", testContext.BaseURL, organizationsApiPath, tc.expectedPage, tc.expectedSize)
				httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, apiResponse))

				args := append([]string{"list", "--output-format", format}, tc.extraArgs...)

				var buf bytes.Buffer
				OrgCmd.SetArgs(args)
				OrgCmd.SetOut(&buf)
				defer OrgCmd.SetOut(nil)

				if err := OrgCmd.Execute(); err != nil {
					t.Fatalf("unexpected error (likely means the request did not hit %s - check the page/pagesize conversion): %v", url, err)
				}

				out := buf.Bytes()

				var page organization.OrganizationsPage
				switch format {
				case "json":
					if err := json.Unmarshal(out, &page); err != nil {
						t.Fatalf("output is not valid JSON: %v\noutput:\n%s", err, out)
					}
					if !bytes.Contains(out, []byte(`"totalCount"`)) {
						t.Fatalf("json output is missing the canonical \"totalCount\" key:\n%s", out)
					}
				case "yaml":
					if err := yaml.Unmarshal(out, &page); err != nil {
						t.Fatalf("output is not valid YAML: %v\noutput:\n%s", err, out)
					}
					// A plain round-trip decode accepts either key casing, so
					// assert on the raw serialized text instead.
					if !bytes.Contains(out, []byte("totalCount:")) {
						t.Fatalf("yaml output is missing the canonical \"totalCount:\" key:\n%s", out)
					}
					if bytes.Contains(out, []byte("totalcount:")) {
						t.Fatalf("yaml output contains lowercased \"totalcount:\" key, want camelCase:\n%s", out)
					}
				}

				gotTotal := -1
				if page.TotalCount != nil {
					gotTotal = *page.TotalCount
				}
				if gotTotal != 1 {
					t.Fatalf("TotalCount = %d, want 1", gotTotal)
				}
				if page.Organizations == nil || len(*page.Organizations) != 1 {
					t.Fatalf("Organizations = %v, want exactly 1 entry", page.Organizations)
				}
				if got := (*page.Organizations)[0].Name; got != "name" {
					t.Fatalf("Organizations[0].Name = %q, want %q", got, "name")
				}
			})
		}
	}
}
