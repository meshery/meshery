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
			Name:             "given all requirements met when running mesheryctl connection list then a list of available connections are displayed",
			Args:             []string{"list"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.response.golden",
			ExpectedResponse: "list.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --count flag provided when running mesheryctl connection list --count then the total count of connections is displayed",
			Args:             []string{"list", "--count"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.api.response.golden",
			ExpectedResponse: "list.count.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --kind flag provided when running mesheryctl connection list --kind valid-kind then the list of connections with specified kind is displayed",
			Args:             []string{"list", "--kind", "kubernetes"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.kind.api.response.golden",
			ExpectedResponse: "list.kind.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --status flag provided when running mesheryctl connection list --status valid-status then the list of connections with specified status is displayed",
			Args:             []string{"list", "--status", "connected"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.status.api.response.golden",
			ExpectedResponse: "list.status.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --page flag provided when running mesheryctl connection list --page page-number then the list of connections with specified page is displayed",
			Args:             []string{"list", "--page", "1"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.page.api.response.golden",
			ExpectedResponse: "list.page.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --pagesize flag provided when running mesheryctl connection list --pagesize page-size-number then the list of connections with specified pagesize is displayed",
			Args:             []string{"list", "--pagesize", "3"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.pagesize.api.response.golden",
			ExpectedResponse: "list.pagesize.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --status flag and --kind flag provided when running mesheryctl connection list --kind valid-kind --status valid-status then the list of connections with specified kind and status is displayed",
			Args:             []string{"list", "--status", "connected", "--kind", "kubernetes"},
			URL:              "/api/integrations/connections",
			Fixture:          "list.connection.status.kind.api.response.golden",
			ExpectedResponse: "list.status.kind.connection.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a --kind flag and --count flag provided when running mesheryctl connection list --kind valid-kind --count then the total count of connections with specified kind is displayed",
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
