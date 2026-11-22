package connections

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta3/connection"
)

func TestConnectionViewCmd(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []utils.MesheryListCommandTest{
		{
			Name:           "given no argument provided when connection view then throw error",
			Args:           []string{"view"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errNoArgMsg, viewUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given multiple arguments provided when connection view then throw error",
			Args:           []string{"view", "foo", "bar"},
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errMultiArgMsg, viewUsageMsg)),
			IsOutputGolden: false,
		},
		{
			Name:           "given an invalid argument for --output-format flag provided when connection view then throw error",
			Args:           []string{"view", connectionId, "--output-format", "foo"},
			URL:            "/api/integrations/connections/" + connectionId,
			Fixture:        "view.connection.api.empty.response.golden",
			ExpectError:    true,
			ExpectedError:  display.ErrInvalidOutputFormat("foo"),
			IsOutputGolden: false,
		},
		{
			Name:             "given a valid connection-id provided when connection view then display detailed information",
			Args:             []string{"view", connectionId},
			URL:              "/api/integrations/connections/" + connectionId,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given a valid --output-format argument provided when connection view then display detailed information in the specified format",
			Args:             []string{"view", connectionId, "--output-format", "yaml"},
			URL:              "/api/integrations/connections/" + connectionId,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
		{
			Name:             "given an uppercase --output-format argument provided when connection view then display detailed information in the specified format",
			Args:             []string{"view", connectionId, "--output-format", "YAML"},
			URL:              "/api/integrations/connections/" + connectionId,
			Fixture:          "view.connection.api.response.golden",
			ExpectedResponse: "view.connection.yaml.output.golden",
			ExpectError:      false,
			IsOutputGolden:   true,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ConnectionsCmd, tests, currDir, "connection")
}

// connectionViewRun describes the view invocation shared by the --save tests
// below, which assert on the file written rather than on golden output.
func connectionViewRun(commandDir string, args []string) utils.MesheryCommandRun {
	return utils.MesheryCommandRun{
		Cmd:        ConnectionsCmd,
		Args:       args,
		CommandDir: commandDir,
		Mocks: []utils.MockURL{
			{
				URL:      "/api/integrations/connections/" + connectionId,
				Response: "view.connection.api.response.golden",
			},
		},
	}
}

// TestConnectionViewSaveCreatesFile verifies that --save writes a file named
// for the connection and the output format. The home directory is isolated so
// the real ~/.meshery is never touched.
func TestConnectionViewSaveCreatesFile(t *testing.T) {
	mesheryDir := utils.IsolateMesheryHome(t)

	run := connectionViewRun(utils.CallerDir(t), []string{"view", connectionId, "--save"})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "connection_minikube.yaml"))
}

// TestConnectionViewSaveWithUppercaseOutputFormat verifies that an uppercase
// --output-format saves the file with a lowercase extension.
func TestConnectionViewSaveWithUppercaseOutputFormat(t *testing.T) {
	mesheryDir := utils.IsolateMesheryHome(t)

	run := connectionViewRun(utils.CallerDir(t), []string{"view", connectionId, "--output-format", "JSON", "--save"})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "connection_minikube.json"))
}

// TestConnectionViewNoSaveWithBrokenHome verifies that view without --save
// succeeds when no home directory can be resolved, proving os.UserHomeDir is
// not called on the display-only path.
func TestConnectionViewNoSaveWithBrokenHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	run := connectionViewRun(utils.CallerDir(t), []string{"view", connectionId})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", err)
	}
}

func TestFetchConnectionByName(t *testing.T) {
	conn1 := &connection.Connection{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "minikube-1",
		ConnectionType: "platform",
	}
	conn2 := &connection.Connection{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "minikube-2",
		ConnectionType: "platform",
	}

	tests := []struct {
		name            string
		page            connection.ConnectionPage
		promptChoice    int
		wantPromptCall  bool
		wantPromptItems []string
		wantConn        *connection.Connection
	}{
		{
			name: "TotalCount > 0 with empty Connections slice returns nil without panicking",
			page: connection.ConnectionPage{
				TotalCount:  1,
				Connections: []*connection.Connection{},
			},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name: "TotalCount > 1 with empty Connections slice returns nil without prompting",
			page: connection.ConnectionPage{
				TotalCount:  2,
				Connections: []*connection.Connection{},
			},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name: "zero results returns nil",
			page: connection.ConnectionPage{
				TotalCount:  0,
				Connections: []*connection.Connection{},
			},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name: "single nil entry returns nil without panicking",
			page: connection.ConnectionPage{
				TotalCount:  1,
				Connections: []*connection.Connection{nil},
			},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name: "all nil entries return nil without prompting",
			page: connection.ConnectionPage{
				TotalCount:  2,
				Connections: []*connection.Connection{nil, nil},
			},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name: "mixed nil and single valid connection skips prompt and returns valid connection",
			page: connection.ConnectionPage{
				TotalCount:  2,
				Connections: []*connection.Connection{nil, conn1},
			},
			wantPromptCall: false,
			wantConn:       conn1,
		},
		{
			name: "exactly one connection returns it directly",
			page: connection.ConnectionPage{
				TotalCount:  1,
				Connections: []*connection.Connection{conn1},
			},
			wantPromptCall: false,
			wantConn:       conn1,
		},
		{
			name: "inconsistent TotalCount > 1 with single connection skips prompt",
			page: connection.ConnectionPage{
				TotalCount:  2,
				Connections: []*connection.Connection{conn1},
			},
			wantPromptCall: false,
			wantConn:       conn1,
		},
		{
			name: "multiple connections prompts user and returns selection",
			page: connection.ConnectionPage{
				TotalCount:  2,
				Connections: []*connection.Connection{conn1, conn2},
			},
			promptChoice:   1,
			wantPromptCall: true,
			wantPromptItems: []string{
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn1.ID.String(), conn1.Name, conn1.ConnectionType),
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn2.ID.String(), conn2.Name, conn2.ConnectionType),
			},
			wantConn: conn2,
		},
		{
			name: "mixed nil and multiple valid connections filters nil entries and preserves order",
			page: connection.ConnectionPage{
				TotalCount:  3,
				Connections: []*connection.Connection{conn1, nil, conn2},
			},
			promptChoice:   1,
			wantPromptCall: true,
			wantPromptItems: []string{
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn1.ID.String(), conn1.Name, conn1.ConnectionType),
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn2.ID.String(), conn2.Name, conn2.ConnectionType),
			},
			wantConn: conn2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			utils.TokenFlag = utils.GetToken(t)

			origRunSelectPrompt := runSelectPrompt
			defer func() {
				runSelectPrompt = origRunSelectPrompt
			}()

			promptCalled := false
			runSelectPrompt = func(label string, items []string) (int, error) {
				promptCalled = true
				if len(items) != len(tt.wantPromptItems) {
					t.Fatalf("prompt items len = %d, want %d", len(items), len(tt.wantPromptItems))
				}
				for i := range tt.wantPromptItems {
					if items[i] != tt.wantPromptItems[i] {
						t.Fatalf("prompt item[%d] = %q, want %q", i, items[i], tt.wantPromptItems[i])
					}
				}
				return tt.promptChoice, nil
			}

			body, err := json.Marshal(tt.page)
			if err != nil {
				t.Fatalf("failed to marshal connection page: %v", err)
			}

			query := url.Values{}
			query.Add("search", "minikube")
			query.Add("pagesize", "all")
			reqURL := fmt.Sprintf("%s/%s?%s", testContext.BaseURL, connectionApiPath, query.Encode())
			httpmock.RegisterResponder(http.MethodGet, reqURL, httpmock.NewBytesResponder(http.StatusOK, body))

			got, err := fetchConnectionByName("minikube")
			if err != nil {
				t.Fatalf("fetchConnectionByName returned unexpected error: %v", err)
			}
			if promptCalled != tt.wantPromptCall {
				t.Fatalf("promptCalled = %v, want %v", promptCalled, tt.wantPromptCall)
			}
			if tt.wantConn == nil {
				if got != nil {
					t.Fatalf("got connection %+v, want nil", got)
				}
				return
			}
			if got == nil || got.ID != tt.wantConn.ID {
				t.Fatalf("got connection = %+v, want ID %s", got, tt.wantConn.ID)
			}
		})
	}
}

func TestSelectConnectionPromptNilEntries(t *testing.T) {
	conn1 := &connection.Connection{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "minikube-1",
		ConnectionType: "platform",
	}
	conn2 := &connection.Connection{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "minikube-2",
		ConnectionType: "platform",
	}

	tests := []struct {
		name            string
		input           []*connection.Connection
		promptChoice    int
		wantPromptCall  bool
		wantPromptItems []string
		wantConn        *connection.Connection
	}{
		{
			name:           "all nil entries return nil without prompting",
			input:          []*connection.Connection{nil, nil},
			wantPromptCall: false,
			wantConn:       nil,
		},
		{
			name:           "mixed nil and valid entries filter nil entries without mutating input",
			input:          []*connection.Connection{nil, conn1, nil, conn2},
			promptChoice:   1,
			wantPromptCall: true,
			wantPromptItems: []string{
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn1.ID.String(), conn1.Name, conn1.ConnectionType),
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn2.ID.String(), conn2.Name, conn2.ConnectionType),
			},
			wantConn: conn2,
		},
		{
			name:           "multiple valid entries preserve order and selection",
			input:          []*connection.Connection{conn1, conn2},
			promptChoice:   0,
			wantPromptCall: true,
			wantPromptItems: []string{
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn1.ID.String(), conn1.Name, conn1.ConnectionType),
				fmt.Sprintf("ID: %s, Name: %s, Type: %s", conn2.ID.String(), conn2.Name, conn2.ConnectionType),
			},
			wantConn: conn1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := append([]*connection.Connection(nil), tt.input...)

			origRunSelectPrompt := runSelectPrompt
			defer func() {
				runSelectPrompt = origRunSelectPrompt
			}()

			promptCalled := false
			runSelectPrompt = func(label string, items []string) (int, error) {
				promptCalled = true
				if len(items) != len(tt.wantPromptItems) {
					t.Fatalf("prompt items len = %d, want %d", len(items), len(tt.wantPromptItems))
				}
				for i := range tt.wantPromptItems {
					if items[i] != tt.wantPromptItems[i] {
						t.Fatalf("prompt item[%d] = %q, want %q", i, items[i], tt.wantPromptItems[i])
					}
				}
				return tt.promptChoice, nil
			}

			got, err := selectConnectionPrompt(tt.input)
			if err != nil {
				t.Fatalf("selectConnectionPrompt returned unexpected error: %v", err)
			}
			if promptCalled != tt.wantPromptCall {
				t.Fatalf("promptCalled = %v, want %v", promptCalled, tt.wantPromptCall)
			}
			if len(tt.input) != len(snapshot) {
				t.Fatalf("input slice len mutated: got %d, want %d", len(tt.input), len(snapshot))
			}
			for i := range snapshot {
				if tt.input[i] != snapshot[i] {
					t.Fatalf("input slice element[%d] mutated: got %v, want %v", i, tt.input[i], snapshot[i])
				}
			}
			if tt.wantConn == nil {
				if got != nil {
					t.Fatalf("got connection %+v, want nil", got)
				}
				return
			}
			if got == nil || got.ID != tt.wantConn.ID {
				t.Fatalf("got connection = %+v, want ID %s", got, tt.wantConn.ID)
			}
		})
	}
}
