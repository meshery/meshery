package filter

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenrios for fetching data
	tests := []utils.MesheryCommandTest{
		{
			Name:             "Fetch Filter List",
			Args:             []string{"list"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
		{
			Name:             "Fetch Filter List with pagesize",
			Args:             []string{"list", "--page", "1", "--pagesize", "10"},
			ExpectedResponse: "list.filter.output.golden",
			Fixture:          "list.filter.api.response.golden",
			URL:              "/api/filter?pagesize=10&page=0",
			ExpectError:      false,
			HttpMethod:       "GET",
			HttpStatusCode:   200,
		},
		{
			Name:             "List filters with invalid page number",
			Args:             []string{"list", "--page", "0"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "List filters with invalid page size",
			Args:             []string{"list", "--pagesize", "0"},
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pageSize '0'")),
		},
	}
	mesheryctlflags.InitValidators(FilterCmd)
	// Run tests
	utils.InvokeMesheryctlTestCommand(t, update, FilterCmd, tests, currDir, "filter")
}

func TestListCmdCount(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:             "Fetch Filter Count",
			Args:             []string{"list", "--count"},
			ExpectedResponse: "list.filter.count.output.golden",
			Fixture:          "filter.list.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
		{
			Name:             "Fetch Filter Count with no matching filters",
			Args:             []string{"list", "nonexistent", "--count"},
			ExpectedResponse: "list.filter.count.empty.output.golden",
			Fixture:          "filter.list.empty.api.response.golden",
			URL:              "/api/filter",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(FilterCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, FilterCmd, tests, currDir, "filter")
}
