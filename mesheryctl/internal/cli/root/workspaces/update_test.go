package workspaces

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta3/workspace"
)

func TestUpdateWorkspaceValidation(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	workspaceID := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given no workspace ID when workspace update then throw error",
			Args:             []string{"update", "--orgId", testOrgId, "-n", validWorkspaceName},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace ID\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given an invalid workspace ID when workspace update then throw error",
			Args:             []string{"update", "not-a-uuid", "--orgId", testOrgId, "-n", validWorkspaceName},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %s\n\n%v", "not-a-uuid",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given no --name or --description when workspace update then throw error",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("at least one of --name or --description must be provided\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given missing orgId when workspace update then throw error",
			Args:             []string{"update", workspaceID, "-n", validWorkspaceName},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --orgId ''")),
		},
		{
			Name:             "given --description with an empty value when workspace update then throw error",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId, "--description", ""},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("clearing the description is not currently supported by the server - provide a non-empty --description\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
		{
			Name:             "given --name with an empty value when workspace update then throw error",
			Args:             []string{"update", workspaceID, "--orgId", testOrgId, "--name", ""},
			Fixture:          "",
			ExpectedResponse: "",
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: utils.ErrInvalidArgument(fmt.Errorf("--name cannot be empty\n\n%v",
				"Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message")),
		},
	}

	mesheryctlflags.InitValidators(WorkSpaceCmd)
	utils.InvokeMesheryctlTestCommand(t, update, WorkSpaceCmd, tests, currDir, "workspaces")
}

func TestUpdateWorkspaceFlow(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	workspaceID := "d56fb25b-f92c-4cd6-821b-2cfd6bb87259"
	mesheryctlflags.InitValidators(WorkSpaceCmd)

	run := func(t *testing.T, args []string, registerPut func(getURL, putURL string)) (string, error) {
		t.Helper()

		testContext := utils.InitTestEnvironment(t)
		defer utils.StopMockery(t)
		defer utils.ResetCommandFlags(WorkSpaceCmd, t)
		out := utils.SetupMeshkitLoggerTesting(t, false)
		utils.TokenFlag = utils.GetToken(t)

		getURL := fmt.Sprintf("%s/api/workspaces/%s?orgId=%s", testContext.BaseURL, workspaceID, testOrgId)
		putURL := fmt.Sprintf("%s/api/workspaces/%s", testContext.BaseURL, workspaceID)
		if registerPut != nil {
			registerPut(getURL, putURL)
		}

		WorkSpaceCmd.SetArgs(args)
		WorkSpaceCmd.SetOut(out)
		defer WorkSpaceCmd.SetOut(nil)

		err := WorkSpaceCmd.Execute()
		return out.String(), err
	}

	t.Run("given a valid workspace ID, orgId and name when workspace update then workspace is updated and the response body is closed", func(t *testing.T) {
		getFixture := utils.NewGoldenFile(t, "get.workspace.for.update.golden", fixturesDir).Load()
		putFixture := utils.NewGoldenFile(t, "update.workspace.api.response.golden", fixturesDir).Load()

		out, err := run(t, []string{"update", workspaceID, "--orgId", testOrgId, "-n", "workspace-test-renamed"}, func(getURL, putURL string) {
			httpmock.RegisterResponder(http.MethodGet, getURL, httpmock.NewStringResponder(200, getFixture))
			httpmock.RegisterResponder(http.MethodPut, putURL, func(req *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(req.Body)
				var payload workspace.WorkspaceUpdatePayload
				if unmarshalErr := json.Unmarshal(body, &payload); unmarshalErr != nil {
					t.Fatalf("PUT body is not valid JSON: %v", unmarshalErr)
				}
				if payload.Description != "" {
					t.Fatalf("PUT body Description = %q, want empty/omitted (unchanged field should not be sent)", payload.Description)
				}
				if payload.Name != "workspace-test-renamed" {
					t.Fatalf("PUT body Name = %q, want %q", payload.Name, "workspace-test-renamed")
				}
				return httpmock.NewStringResponse(200, putFixture), nil
			})
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := fmt.Sprintf("Workspace with ID %s has been updated\n", workspaceID)
		if out != want {
			t.Fatalf("output = %q, want %q", out, want)
		}
	})

	t.Run("given an --orgId that does not match the workspace's actual organization when workspace update then it is refused without ever issuing the PUT", func(t *testing.T) {
		getFixture := utils.NewGoldenFile(t, "get.workspace.for.update.other.org.golden", fixturesDir).Load()

		_, err := run(t, []string{"update", workspaceID, "--orgId", testOrgId, "-n", "workspace-test-renamed"}, func(getURL, _ string) {
			httpmock.RegisterResponder(http.MethodGet, getURL, httpmock.NewStringResponder(200, getFixture))
		})

		utils.AssertMeshkitErrorsEqual(t, err, ErrOrganizationMismatch(workspaceID, "9a8f5c10-1111-2222-3333-444455556666", testOrgId))
	})

	t.Run("given a workspace ID that no longer exists when workspace update then it is reported as a hard failure, not a silent success", func(t *testing.T) {
		_, err := run(t, []string{"update", workspaceID, "--orgId", testOrgId, "-n", "workspace-test-renamed"}, func(getURL, _ string) {
			httpmock.RegisterResponder(http.MethodGet, getURL, httpmock.NewStringResponder(404, ""))
		})

		if err == nil {
			t.Fatal("expected an error for a not-found workspace, got nil (this must not exit 0)")
		}
		utils.AssertMeshkitErrorsEqual(t, err, utils.ErrNotFound(fmt.Errorf("workspace with ID %s not found", workspaceID)))
	})
}
