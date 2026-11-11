package design

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDesignListCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "given no arguments provided when design list then design list is displayed",
			Args:             []string{"list", "--page", "1"},
			ExpectedResponse: "list.design.output.golden",
			Fixture:          "list.design.api.response.golden",
			URL:              "/api/pattern?page=0&pagesize=10",
			ExpectError:      false,
		},
		{
			Name:             "given extra positional arguments when design list then throw error",
			Args:             []string{"list", "extra-arg"},
			ExpectedResponse: "",
			Fixture:          "list.design.api.response.golden",
			URL:              "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg)),
		},
		{
			Name:             "given --page 0 when design list then throw invalid flag error",
			Args:             []string{"list", "--page", "0"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "given a negative --page when design list then throw invalid flag error",
			Args:             []string{"list", "--page", "-1"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:             "given --pagesize 0 when design list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "0"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "given a negative --pagesize when design list then throw invalid flag error",
			Args:             []string{"list", "--pagesize", "-1"},
			ExpectedResponse: "",
			URL:              "",
			Fixture:          "list.design.api.response.golden",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '-1'")),
		},
	}

	mesheryctlflags.InitValidators(DesignCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, DesignCmd, tests, currDir, "design")
}
