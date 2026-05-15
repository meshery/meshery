package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// expectedViewFlagError generates the expected error for an invalid view flag
func expectedViewFlagError(outputFormat string) error {
	fv := mesheryctlflags.GetFlagValidator()
	return fv.Validate(&componentViewFlags{OutputFormat: outputFormat})
}

func TestComponentView(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:           "given no argument provided when component view then throw error",
			Args:           []string{"view"},
			URL:            "/api/meshmodels/components",
			Fixture:        "components.empty.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errInvalidArg, viewUsageMsg)),
		},
		{
			Name:             "given non-existent component provided when component view then display no component found",
			Args:             []string{"view", "foo"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10&search=foo", componentApiPath),
			Fixture:          "components.empty.api.response.golden",
			ExpectedResponse: "",
			IsOutputGolden:   false,
			ExpectError:      true,
			ExpectedError:    utils.ErrNotFound(fmt.Errorf("%s%s", errNoComponentFound, "foo")),
		},
		{
			Name:           "given multiple component name is provided when running component view then throw error",
			Args:           []string{"view", "Test", "ACL"},
			URL:            fmt.Sprintf("/%s?pagesize=all&search=Test&search=ACL", componentApiPath),
			Fixture:        "components.empty.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errInvalidArg, viewUsageMsg)),
		},
		{
			Name:             "given valid component provided when component view then display detailed information",
			Args:             []string{"view", "Test"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10&search=Test", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.view.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
		{
			Name:           "given invalid argument value provided for --output-format flag when component view then throw error",
			Args:           []string{"view", "Test", "--output-format", "invalid"},
			URL:            fmt.Sprintf("/%s?pagesize=all&search=Test", componentApiPath),
			Fixture:        "components.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  expectedViewFlagError("invalid"),
		},
		{
			Name:             "given valid component provided and --output-format flag when component view then display detail information in the specified output format",
			Args:             []string{"view", "Test", "--output-format", "json"},
			URL:              fmt.Sprintf("/%s?page=0&pagesize=10&search=Test", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.view.json.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
		{
			Name:             "given a valid UUID provided when component view then display detailed information of the component",
			Args:             []string{"view", "fda1c4e7-14ae-4435-8236-adfb9cea0395"},
			URL:              fmt.Sprintf("/%s?id=fda1c4e7-14ae-4435-8236-adfb9cea0395&page=0&pagesize=10", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.view.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
		{
			Name:          "given a non-existent UUID is provided when running mesheryctl component view non-existent-uuid then an error message is displayed",
			Args:          []string{"view", "a12b3c4d-5e6f-4890-abcd-ef1234567890"},
			URL:           fmt.Sprintf("/%s?id=a12b3c4d-5e6f-4890-abcd-ef1234567890&page=0&pagesize=10", componentApiPath),
			Fixture:       "components.empty.api.response.golden",
			ExpectError:   true,
			ExpectedError: utils.ErrNotFound(fmt.Errorf("%s%s", errNoComponentFound, "a12b3c4d-5e6f-4890-abcd-ef1234567890")),
		},
	}

	mesheryctlflags.InitValidators(ComponentCmd)
	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currDir, "component")
}
