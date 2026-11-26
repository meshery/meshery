package organizations

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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
			Name:             "List organizations without providing organization ID",
			Args:             []string{"list"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10", organizationsApiPath),
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "list.organization.golden",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when organization list then throw error",
			Args:             []string{"list", "extra-arg"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
		{
			Name:             "given --page 0 when organization list then throw invalid flag error",
			Args:             []string{"list", "--page", "0"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "given a negative --page when organization list then throw invalid flag error",
			Args:             []string{"list", "--page", "-1"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:             "given --pagesize 0 when organization list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "0"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "given a negative --pagesize when organization list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "-1"},
			URL:              "",
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '-1'")),
		},
		{
			Name:             "given valid --page and --pagesize when organization list then list organizations",
			Args:             []string{"list", "--page", "2", "--pagesize", "5"},
			URL:              fmt.Sprintf("/%s?page=1&pagesize=5", organizationsApiPath),
			Fixture:          "list.organization.response.golden",
			ExpectedResponse: "list.organization.page2.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(OrgCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, OrgCmd, tests, currentDirectory, "organization")
}
