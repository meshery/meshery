package organizations

import (
	"bytes"
	"encoding/json"
	"fmt"
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
// harness: the response envelope's custom MarshalJSON (see
// server/models/organization.go) makes an exact byte-for-byte golden fragile
// to hand-maintain, so this instead round-trips the command's output back
// into models.OrganizationsPage and asserts on the decoded fields, plus
// asserts on the raw serialized key names (see the yaml case below).
//
// Output is captured via cmd.SetOut(&buf) (listOrgsAsStructuredOutput writes
// to cmd.OutOrStdout(), see list.go) rather than redirecting the
// process-global os.Stdout, so this test is safe to run alongside others
// even if a future change adds t.Parallel() to this package.
func TestListOrganizationsStructuredOutput(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	fixturesDir := filepath.Join(currentDirectory, "fixtures")

	// Each case registers the mock responder at the exact URL
	// listOrgsAsStructuredOutput should build for the given --page/--pagesize
	// flags. If the 1-based-to-0-based page conversion (or the pagesize
	// default/pass-through) regresses, the request will miss this exact URL,
	// httpmock will return a "no responder found" error, and the command
	// will fail - so a passing test here is itself the assertion that the
	// pagination math produced the right request. (We don't additionally
	// assert the decoded response's Page/PageSize fields: the fixture below
	// is a fixed canned response - it always reports page=0/pageSize=10
	// regardless of what was requested - so asserting those against the
	// requested values would be asserting something the fixture can't
	// actually reflect.)
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

				var page models.OrganizationsPage
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
					// Regression guard: models.OrganizationsPage previously had
					// no yaml struct tags, so yaml.v3 fell back to lowercasing
					// Go field names (totalcount, pagesize) instead of the
					// canonical camelCase used by JSON and the API. A plain
					// yaml.Unmarshal round-trip doesn't catch this, since
					// decoding back into the same struct accepts either
					// casing - so assert on the raw serialized text instead.
					if !bytes.Contains(out, []byte("totalCount:")) {
						t.Fatalf("yaml output is missing the canonical \"totalCount:\" key (got lowercased/legacy key instead):\n%s", out)
					}
					if bytes.Contains(out, []byte("totalcount:")) {
						t.Fatalf("yaml output contains lowercased \"totalcount:\" key, want camelCase \"totalCount:\":\n%s", out)
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
}
