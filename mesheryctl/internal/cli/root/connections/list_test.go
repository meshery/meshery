package connections

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestConnectionListCmd(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []utils.MesheryListCommandTest{
		{
			Name:             "List connections empty",
			Args:             []string{"list"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.empty.response.golden",
			ExpectedResponse: "list.connection.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "List connections",
			Args:             []string{"list"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.response.golden",
			ExpectedResponse: "list.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display count of connections empty result",
			Args:             []string{"list", "--count"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.empty.response.golden",
			ExpectedResponse: "list.count.connection.empty.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display count of connections",
			Args:             []string{"list", "--count"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.response.golden",
			ExpectedResponse: "list.count.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --kind flag",
			Args:             []string{"list", "--kind", "kubernetes"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.kind.api.response.golden",
			ExpectedResponse: "list.kind.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --status flag",
			Args:             []string{"list", "--status", "connected"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.status.api.response.golden",
			ExpectedResponse: "list.status.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --page flag",
			Args:             []string{"list", "--page", "1"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.page.api.response.golden",
			ExpectedResponse: "list.page.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --pagesize flag",
			Args:             []string{"list", "--pagesize", "3"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.pagesize.api.response.golden",
			ExpectedResponse: "list.pagesize.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --status and --kind flag",
			Args:             []string{"list", "--status", "connected", "--kind", "kubernetes"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.status.kind.api.response.golden",
			ExpectedResponse: "list.status.kind.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "Display list with --count and --kind flag",
			Args:             []string{"list", "--kind", "meshery", "--count"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.kind.count.api.response.golden",
			ExpectedResponse: "list.kind.count.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}
