package organizations

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func expectedOrgViewFlagError(outputFormat string) error {
	fv := mesheryctlflags.GetFlagValidator()
	return fv.Validate(&orgViewFlags{OutputFormat: outputFormat})
}

func TestViewOrganization(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	mesheryctlflags.InitValidators(OrgCmd)

	testOrgID := "3a4b5c6d-7e8f-9a0b-1c2d-3e4f5a6b7c8d"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no arguments when running organization view then return error",
			Args:             []string{"view"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", organizationsApiPath),
			Fixture:          "view.organization.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one organization name or ID\n\nUsage: mesheryctl organization view [organization-name|organization-id]\nRun 'mesheryctl organization view --help' to see detailed help message")),
		},
		{
			Name:             "given too many arguments when running organization view then return error",
			Args:             []string{"view", testOrgID, "extra-arg"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", organizationsApiPath),
			Fixture:          "view.organization.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one organization name or ID\n\nUsage: mesheryctl organization view [organization-name|organization-id]\nRun 'mesheryctl organization view --help' to see detailed help message")),
		},
		{
			Name:             "given valid organization ID when running organization view then display organization details",
			Args:             []string{"view", testOrgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", organizationsApiPath, testOrgID),
			Fixture:          "view.organization.api.response.golden",
			ExpectedResponse: "view.organization.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given invalid output format flag when running organization view then return error",
			Args:             []string{"view", testOrgID, "--output-format", "invalid"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s", organizationsApiPath, testOrgID),
			Fixture:          "view.organization.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    expectedOrgViewFlagError("invalid"),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, OrgCmd, tests, currDir, "organization")
}
