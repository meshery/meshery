package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestListComponent(t *testing.T) {

	mesheryctlflags.InitValidators(ComponentCmd)
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	testUrl := fmt.Sprintf("/%s?page=0&pagesize=10", componentApiPath)
	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "list components with page number",
			Args:             []string{"list", "--page", "1"},
			URL:              testUrl,
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components non empty count",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.list.count.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "list components empty count",
			Args:             []string{"list", "--count"},
			URL:              testUrl,
			Fixture:          "components.empty.api.response.golden",
			ExpectedResponse: "components.list.count.empty.ouput.golden",
			ExpectError:      false,
		},
	}

	mesheryctlflags.InitValidators(ComponentCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currentDirectory, "component")

}

// TestComponentFlagValidation verifies that component flag validation errors report the expected flag name.
func TestComponentFlagValidation(t *testing.T) {
	mesheryctlflags.InitValidators(listComponentCmd)

	// Validate list flag struct with invalid pagesize
	err := mesheryctlflags.ValidateCmdFlags(listComponentCmd, &componentListFlag{PageSize: 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid value for --pagesize '0'")

	// Validate search flag struct with invalid pagesize
	mesheryctlflags.InitValidators(searchComponentsCmd)
	err = mesheryctlflags.ValidateCmdFlags(searchComponentsCmd, &cmdComponentSearchFlags{PageSize: 0})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid value for --pagesize '0'")
}
