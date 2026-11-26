package environments

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"
)

func TestListEnvironment(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List environments organization ID not provided",
			Args:             []string{"list"},
			URL:              "/api/environments",
			Fixture:          "list.environment.without.orgId.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ orgId ] isn't specified\n\nUsage: mesheryctl environment list --orgId [orgId]\nRun 'mesheryctl environment list --help' to see detailed help message")),
		},
		{
			Name:             "List environments non available",
			Args:             []string{"list", "--orgId", testConstants["orgId"]},
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, testConstants["orgId"]),
			Fixture:          "list.environment.empty.response.golden",
			ExpectedResponse: "list.environment.empty.golden",
			ExpectError:      false,
		},
		{
			Name:             "List environments available",
			Args:             []string{"list", "--orgId", testConstants["orgId"]},
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, testConstants["orgId"]),
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "list.environment.success.golden",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when environment list then throw error",
			Args:             []string{"list", "extra-arg", "--orgId", testConstants["orgId"]},
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
		{
			Name:             "given --page 0 when environment list then throw invalid flag error",
			Args:             []string{"list", "--orgId", testConstants["orgId"], "--page", "0"},
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "given a negative --page when environment list then throw invalid flag error",
			Args:             []string{"list", "--orgId", testConstants["orgId"], "--page", "-1"},
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:             "given --pagesize 0 when environment list then throw invalid flag error",
			Args:             []string{"list", "--orgId", testConstants["orgId"], "--pagesize", "0"},
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "given a negative --pagesize when environment list then throw invalid flag error",
			Args:             []string{"list", "--orgId", testConstants["orgId"], "--pagesize", "-1"},
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '-1'")),
		},
		{
			Name:             "given valid --page and --pagesize when environment list then list environments",
			Args:             []string{"list", "--orgId", testConstants["orgId"], "--page", "2", "--pagesize", "5"},
			URL:              fmt.Sprintf("/%s?orgId=%s&page=1&pagesize=5", environmentApiPath, testConstants["orgId"]),
			Fixture:          "list.environment.response.golden",
			ExpectedResponse: "list.environment.page2.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(EnvironmentCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, EnvironmentCmd, tests, currDir, "environments")
}

func TestProcessEnvironmentData_NilResponse(t *testing.T) {
	rows, total := processEnvironmentData(nil)
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows on nil response, got %d", len(rows))
	}
	if total != 0 {
		t.Fatalf("expected 0 total on nil response, got %d", total)
	}
}

func TestProcessEnvironmentData_EmptyDescriptionFallback(t *testing.T) {
	envPage := &environment.EnvironmentPage{
		TotalCount: 3,
		Environments: []environment.Environment{
			{
				Name:        "empty-desc-env",
				Description: "",
			},
			{
				Name:        "whitespace-desc-env",
				Description: "   \t\n  ",
			},
			{
				Name:        "valid-desc-env",
				Description: "Production Environment",
			},
		},
	}
	rows, total := processEnvironmentData(envPage)
	if total != 3 {
		t.Fatalf("expected total 3, got %d", total)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0][3] != "N/A" {
		t.Fatalf("expected empty description 'N/A', got %q", rows[0][3])
	}
	if rows[1][3] != "N/A" {
		t.Fatalf("expected whitespace-only description 'N/A', got %q", rows[1][3])
	}
	if rows[2][3] != "Production Environment" {
		t.Fatalf("expected preserved description 'Production Environment', got %q", rows[2][3])
	}
}
