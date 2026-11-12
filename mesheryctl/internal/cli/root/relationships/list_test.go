package relationships

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestList(t *testing.T) {
	mesheryctlflags.InitValidators(RelationshipCmd)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List registered relationships empty",
			Args:             []string{"list"},
			URL:              "/api/registry/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "list.relationship.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "List registered relationships",
			Args:             []string{"list"},
			URL:              "/api/registry/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "list.relationship.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display count of registered relationships empty result",
			Args:             []string{"list", "--count"},
			URL:              "/api/registry/relationships",
			Fixture:          "list.relationship.api.empty.response.golden",
			ExpectedResponse: "list.count.relationship.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display count of registered relationships",
			Args:             []string{"list", "--count"},
			URL:              "/api/registry/relationships",
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "list.count.relationship.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "List registered relationships with --page 0",
			Args:             []string{"list", "--page", "0"},
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '0'")),
		},
		{
			Name:             "List registered relationships with a negative --page",
			Args:             []string{"list", "--page", "-1"},
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --page '-1'")),
		},
		{
			Name:             "List registered relationships with --pagesize 0",
			Args:             []string{"list", "--pagesize", "0"},
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '0'")),
		},
		{
			Name:             "List registered relationships with a negative --pagesize",
			Args:             []string{"list", "--pagesize", "-1"},
			Fixture:          "list.relationship.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --pagesize '-1'")),
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, RelationshipCmd, tests, currDir, "relationships")
}
