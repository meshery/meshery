package model

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestListModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testUrl := fmt.Sprintf("/%s?page=0&pagesize=10", modelsApiPath)
	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List model",
			Args:             []string{"list", "--page", "1"},
			URL:              testUrl,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "list.model.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "List model no data",
			Args:             []string{"list"},
			URL:              testUrl,
			Fixture:          "list.model.empty.api.response.golden",
			ExpectedResponse: "list.model.empty.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Display model count no data",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "list.model.empty.api.response.golden",
			ExpectedResponse: "list.model.count.empty.ouput.golden",
			ExpectError:      false,
		},
		{
			Name:             "Display count model",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "list.model.api.response.golden",
			ExpectedResponse: "list.model.count.ouput.golden",
			ExpectError:      false,
		},
		// The list flags live on the subcommand, which ResetCommandFlags does not
		// reach, so each invalid-flag case sets the other flag to a valid value.
		{
			Name:           "given negative page when model list then throw error",
			Args:           []string{"list", "--page", "-1", "--pagesize", "10"},
			URL:            testUrl,
			Fixture:        "list.model.api.response.golden",
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:           "given zero page size when model list then throw error",
			Args:           []string{"list", "--page", "1", "--pagesize", "0"},
			URL:            testUrl,
			Fixture:        "list.model.api.response.golden",
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
	}

	mesheryctlflags.InitValidators(ModelCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
