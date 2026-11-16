package environments

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
)

func environmentViewRun(commandDir string, args []string) utils.MesheryCommandRun {
	return utils.MesheryCommandRun{
		Cmd:        EnvironmentCmd,
		Args:       args,
		CommandDir: commandDir,
		Mocks: []utils.MockURL{
			{
				URL:      "/api/environments?orgId=" + testConstants["orgId"],
				Response: "view.environment.api.response.golden",
			},
		},
	}
}

func TestEnvironmentViewNoSaveWithBrokenHome(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("USERPROFILE", "")

	run := environmentViewRun(utils.CallerDir(t), []string{"view", "--orgId", testConstants["orgId"]})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("view without --save should succeed even with no HOME: %v", err)
	}
}

func TestEnvironmentViewSaveCreatesFile(t *testing.T) {
	mesheryDir := utils.IsolateMesheryHome(t)

	run := environmentViewRun(utils.CallerDir(t), []string{"view", "--orgId", testConstants["orgId"], "--save"})
	if _, err := utils.RunMesheryctlCommand(t, run); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.AssertFileExists(t, filepath.Join(mesheryDir, "environment_test-environment.yaml"))
}

func TestEnvironmentViewExtraArgsFails(t *testing.T) {
	run := environmentViewRun(utils.CallerDir(t), []string{"view", "extra-arg", "--orgId", testConstants["orgId"]})
	if _, err := utils.RunMesheryctlCommand(t, run); err == nil {
		t.Fatalf("expected error when extra positional argument is passed, got nil")
	}
}

func TestEnvironmentViewInconsistentTotalCount(t *testing.T) {
	orgID := uuid.Must(uuid.FromString(testConstants["orgId"]))
	env1 := environment.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "prod-env",
		OrganizationID: orgID,
	}
	env2 := environment.Environment{
		ID:             uuid.Must(uuid.NewV4()),
		Name:           "staging-env",
		OrganizationID: orgID,
	}

	tests := []struct {
		name            string
		rawBody         []byte
		page            *environment.EnvironmentPage
		promptChoice    int
		wantPromptCall  bool
		wantPromptItems []string
	}{
		{
			name:           "nil environmentResponse returns nil without panicking",
			rawBody:        []byte("null"),
			wantPromptCall: false,
		},
		{
			name: "TotalCount > 0 with empty Environments slice returns nil without panicking",
			page: &environment.EnvironmentPage{
				TotalCount:   1,
				Environments: []environment.Environment{},
			},
			wantPromptCall: false,
		},
		{
			name: "TotalCount > 1 with empty Environments slice returns nil without prompting",
			page: &environment.EnvironmentPage{
				TotalCount:   2,
				Environments: []environment.Environment{},
			},
			wantPromptCall: false,
		},
		{
			name: "inconsistent TotalCount > 1 with single environment skips prompt",
			page: &environment.EnvironmentPage{
				TotalCount:   2,
				Environments: []environment.Environment{env1},
			},
			wantPromptCall: false,
		},
		{
			name: "multiple environments prompts user and selects chosen environment",
			page: &environment.EnvironmentPage{
				TotalCount:   2,
				Environments: []environment.Environment{env1, env2},
			},
			promptChoice:   1,
			wantPromptCall: true,
			wantPromptItems: []string{
				fmt.Sprintf("ID: %s, Name: %s, Owner: %s, Organization: %s", env1.ID, env1.Name, env1.Owner, env1.OrganizationID),
				fmt.Sprintf("ID: %s, Name: %s, Owner: %s, Organization: %s", env2.ID, env2.Name, env2.Owner, env2.OrganizationID),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			utils.TokenFlag = utils.GetToken(t)

			origFlags := environmentViewFlagsProvided
			origRunSelectPrompt := runSelectPrompt
			defer func() {
				environmentViewFlagsProvided = origFlags
				runSelectPrompt = origRunSelectPrompt
			}()
			environmentViewFlagsProvided = environmentViewFlags{
				outputFormat: "yaml",
				save:         false,
			}

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

			body := tt.rawBody
			if body == nil {
				var err error
				body, err = json.Marshal(tt.page)
				if err != nil {
					t.Fatalf("failed to marshal environment page: %v", err)
				}
			}

			reqURL := fmt.Sprintf("%s/%s?orgId=%s", testContext.BaseURL, environmentApiPath, testConstants["orgId"])
			httpmock.RegisterResponder(http.MethodGet, reqURL, httpmock.NewBytesResponder(http.StatusOK, body))

			EnvironmentCmd.SetArgs([]string{"view", "--orgId", testConstants["orgId"], "--save=false"})
			if err := EnvironmentCmd.Execute(); err != nil {
				t.Fatalf("EnvironmentCmd.Execute() returned unexpected error: %v", err)
			}
			if promptCalled != tt.wantPromptCall {
				t.Fatalf("promptCalled = %v, want %v", promptCalled, tt.wantPromptCall)
			}
		})
	}
}
