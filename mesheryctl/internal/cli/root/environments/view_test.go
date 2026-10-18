package environments

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestViewEnvironment(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	orgID := testConstants["orgId"]
	environmentName := testConstants["environmentName"]
	environmentID := "5d8f6f1a-9c2b-4d7e-8a41-3f0e2b7c9d10"

	tests := []utils.MesheryCommandTest{
		{
			Name:             "given too many arguments when running environment view then return error",
			Args:             []string{"view", environmentID, "extra-arg", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("please provide at most one environment name or ID\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message")),
		},
		{
			Name:             "given no orgId when running environment view then return error",
			Args:             []string{"view", environmentID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidArgument(errors.New("[ orgId ] isn't specified\n\nUsage: mesheryctl environment view [environment-name|environment-id] --orgId [orgId]\nRun 'mesheryctl environment view --help' to see detailed help message")),
		},
		{
			Name:             "given invalid orgId when running environment view then return error",
			Args:             []string{"view", environmentID, "--orgId", "not-a-uuid"},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s", environmentApiPath),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrInvalidUUID(fmt.Errorf("invalid orgId: %s", "not-a-uuid")),
		},
		{
			Name:             "given environment ID when running environment view then display environment",
			Args:             []string{"view", environmentID, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s?orgId=%s", environmentApiPath, environmentID, orgID),
			Fixture:          "view.environment.api.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given environment ID from another organization when running environment view then return not found error",
			Args:             []string{"view", environmentID, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s/%s?orgId=%s", environmentApiPath, environmentID, orgID),
			Fixture:          "view.environment.other.org.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrNotFound(fmt.Errorf("No environment found with ID %s in organization: %s", environmentID, orgID)),
		},
		{
			Name:             "given environment name with a single match when running environment view then display environment",
			Args:             []string{"view", environmentName, "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10&search=%s", environmentApiPath, orgID, environmentName),
			Fixture:          "view.environment.search.single.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "given no environment argument and a single environment when running environment view then display environment",
			Args:             []string{"view", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, orgID),
			Fixture:          "view.environment.search.single.response.golden",
			ExpectedResponse: "view.environment.output.golden",
			ExpectError:      false,
		},
		{
			// The server reports the total as camelCase totalCount (v1beta3). It must be
			// read as 11, not 0, so the command pages instead of reporting nothing found.
			Name:             "given more environments than one page when running environment view without a terminal then report every match",
			Args:             []string{"view", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, orgID),
			Fixture:          "view.environment.first.page.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    display.ErrAmbiguousSelection(11),
		},
		{
			Name:             "given unknown environment name when running environment view then return not found error",
			Args:             []string{"view", "does-not-exist", "--orgId", orgID},
			HttpMethod:       "GET",
			HttpStatusCode:   200,
			URL:              fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10&search=%s", environmentApiPath, orgID, "does-not-exist"),
			Fixture:          "view.environment.search.empty.response.golden",
			ExpectedResponse: "",
			ExpectError:      true,
			ExpectedError:    utils.ErrNotFound(errors.New("No environment(s) found with name: does-not-exist")),
		},
	}

	utils.InvokeMesheryctlTestCommand(t, update, EnvironmentCmd, tests, currDir, "environment")
}

func TestViewEnvironment_PaginationSecondPage(t *testing.T) {
	orgID := testConstants["orgId"]

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	testdataDir := filepath.Join(currDir, "testdata")

	originalTerminal := utils.IsInteractiveTerminal
	utils.IsInteractiveTerminal = func() bool { return true }
	t.Cleanup(func() { utils.IsInteractiveTerminal = originalTerminal })

	testContext := utils.InitTestEnvironment(t)
	utils.TokenFlag = utils.GetToken(t)
	defer func() {
		utils.TokenFlag = "Not Set"
	}()

	page0Resp := utils.NewGoldenFile(t, "view.environment.first.page.response.golden", fixturesDir).Load()
	page1Resp := utils.NewGoldenFile(t, "view.environment.second.page.response.golden", fixturesDir).Load()
	expectedOutput := utils.NewGoldenFile(t, "view.environment.second.page.output.golden", testdataDir).Load()

	urlPage0 := testContext.BaseURL + fmt.Sprintf("/%s?orgId=%s&page=0&pagesize=10", environmentApiPath, orgID)
	urlPage1 := testContext.BaseURL + fmt.Sprintf("/%s?orgId=%s&page=1&pagesize=10", environmentApiPath, orgID)

	page1Requested := false
	httpmock.RegisterResponder("GET", urlPage0, httpmock.NewStringResponder(200, page0Resp))
	httpmock.RegisterResponder("GET", urlPage1, func(req *http.Request) (*http.Response, error) {
		page1Requested = true
		return httpmock.NewStringResponse(200, page1Resp), nil
	})

	promptCount := 0
	origRunPrompt := display.RunSelectPrompt
	display.RunSelectPrompt = func(p promptui.Select) (int, string, error) {
		promptCount++
		if promptCount == 1 {
			// On page 0, select "Load More....." (last item, index 10)
			return 10, "Load More.....", nil
		}
		// On page 1, select env-1 (index 0)
		return 0, "env-1", nil
	}
	t.Cleanup(func() { display.RunSelectPrompt = origRunPrompt })

	cmd := EnvironmentCmd
	defer utils.ResetCommandFlags(cmd, t)

	originalStdout := os.Stdout
	b := utils.SetupMeshkitLoggerTesting(t, false)
	defer func() {
		os.Stdout = originalStdout
	}()

	cmd.SetArgs([]string{"view", "--orgId", orgID})
	cmd.SetOut(b)

	err := cmd.Execute()
	assert.NoError(t, err)
	assert.True(t, page1Requested, "expected page 1 to be requested from server")
	assert.Equal(t, 2, promptCount, "expected 2 prompts to be presented (page 0 and page 1)")
	utils.Equals(t, expectedOutput, b.String())
}
