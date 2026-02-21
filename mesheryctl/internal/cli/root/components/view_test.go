package components

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestComponentView(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	currDir := filepath.Dir(filename)

	componentApiPath := "api/meshmodels/components"

	tests := []utils.MesheryListCommandTest{
		{
			Name:           "given no component is provided when running mesheryctl component view then an error message is displayed",
			Args:           []string{"view"},
			URL:            "/api/meshmodels/components",
			Fixture:        "components.empty.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("[component name] is required but not specified\n\n%s", errViewCmdMsg)),
		},
		{
			Name:             "given a non-existent component is provided when running mesheryctl component view non-existent-component then it displays empty output",
			Args:             []string{"view", "foo"},
			URL:              fmt.Sprintf("/%s?pagesize=all&search=foo", componentApiPath),
			Fixture:          "components.empty.api.response.golden",
			ExpectedResponse: "components.view.empty.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
		{
			Name:           "given multiple component name is provided when running mesheryctl component view component-1 component-2 then an error message is displayed",
			Args:           []string{"view", "Test", "ACL"},
			URL:            fmt.Sprintf("/%s?pagesize=all&search=Test&search=ACL", componentApiPath),
			Fixture:        "components.empty.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", errViewCmdMsg)),
		},
		{
			Name:             "given a valid component is provided when running mesheryctl component view valid-component then the detailed output of the component is displayed",
			Args:             []string{"view", "Test"},
			URL:              fmt.Sprintf("/%s?pagesize=all&search=Test", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.view.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
		{
			Name:           "given an invalid argument is provided for --output-format flag when running mesheryctl component view valid-component --output-format invalid-format then an error message is displayed",
			Args:           []string{"view", "Test", "--output-format", "invalid"},
			URL:            fmt.Sprintf("/%s?pagesize=all&search=Test", componentApiPath),
			Fixture:        "components.api.response.golden",
			IsOutputGolden: false,
			ExpectError:    true,
			ExpectedError:  display.ErrInvalidOutputFormat("invalid"),
		},
		{
			Name:             "given a valid argument is provided for --output-format flag when running mesheryctl component view valid-component --output-format valid-format then a detailed output is displayed in specified format",
			Args:             []string{"view", "Test", "--output-format", "json"},
			URL:              fmt.Sprintf("/%s?pagesize=all&search=Test", componentApiPath),
			Fixture:          "components.api.response.golden",
			ExpectedResponse: "components.view.json.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ComponentCmd, tests, currDir, "component")

}
