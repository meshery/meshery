package design

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUndeployCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	_, expectedFileReadErr := os.ReadFile(invalidFilePath)

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given valid file when design undeploy then design is undeployed",
			Args:             []string{"undeploy", "--yes", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "undeploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy?contexts=all",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid relative file when design undeploy then design is undeployed",
			Args:             []string{"undeploy", "--yes", "-f", filepath.Join("fixtures", "sampleDesign.golden")},
			ExpectedResponse: "undeploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy?contexts=all",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given invalid file path when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(expectedFileReadErr),
		},
		{
			Name:             "given URL when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", "https://raw.githubusercontent.com/meshery/meshery/master/samples/test.yaml"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("URLs are not currently supported")),
		},
		{
			Name:             "given no args when design undeploy then usage error is returned",
			Args:             []string{"undeploy"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: ErrUndeployDesign(
				fmt.Errorf("provide either a design ID or -f [filepath]"),
			),
		},
		{
			Name:             "given empty file flag when design undeploy then error is thrown",
			Args:             []string{"undeploy", "-f", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError: ErrUndeployDesign(
				fmt.Errorf("Usage: mesheryctl design undeploy -f [filepath]"),
			),
		},
		{
			Name:             "given empty target when design undeploy then error is thrown",
			Args:             []string{"undeploy", "   "},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty")),
		},
		{
			Name:             "given too many positional arguments when design undeploy then error is thrown",
			Args:             []string{"undeploy", "design1", "design2"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("expected at most one positional argument, got 2")),
		},
		{
			Name:             "given both positional argument and file flag when design undeploy then error is thrown",
			Args:             []string{"undeploy", "my-design", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a positional argument and the --file flag")),
		},
		{
			Name:             "given nonexistent design provided when design undeploy then error is thrown",
			Args:             []string{"undeploy", "--yes", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          "/api/pattern/deploy?contexts=all",
					Response:     "undeploy.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          fmt.Sprintf("%s/api/pattern", testContext.BaseURL),
					Response:     "undeploy.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(filepath.Join(fixturesDir, "sampleDesign.golden")),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

func TestUndeployConfirmation(t *testing.T) {
	for _, test := range []struct {
		name, input string
		yes         bool
		wantGets    int
		wantPosts   int
		wantDeletes int
	}{
		{name: "decline", input: "n\n", wantGets: 1},
		{name: "accept", input: "y\n", wantGets: 1, wantPosts: 1, wantDeletes: 1},
		{name: "skip", yes: true, wantPosts: 1, wantDeletes: 1},
		{name: "EOF", wantGets: 1},
	} {
		t.Run(test.name, func(t *testing.T) {
			args := []string{"undeploy", "-f", filepath.Join("fixtures", "sampleDesign.golden")}
			if test.yes {
				args = append(args, "--yes")
			}
			output, prompt, gets, posts, deletes, err := runRemovalPromptTest(t, args, test.input)
			if test.name == "EOF" {
				require.ErrorContains(t, err, "--yes")
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, test.wantGets, gets)
			assert.Equal(t, test.wantPosts, posts)
			assert.Equal(t, test.wantDeletes, deletes)
			if test.yes {
				assert.Empty(t, prompt)
			} else {
				assert.Contains(t, prompt, "Design 'SampleApp' will be undeployed from 2 connected Kubernetes cluster(s): east, west.")
			}
			if test.name == "decline" {
				assert.Contains(t, output, "Undeploy cancelled")
			}
		})
	}
}

func TestUndeployConfirmationListsEveryConnectedClusterPage(t *testing.T) {
	output, prompt, gets, posts, deletes, err := runRemovalPromptTestWithClusterPages(t,
		[]string{"undeploy", "-f", filepath.Join("fixtures", "sampleDesign.golden")}, "n\n",
		"removal.connections.page0.response.golden", "removal.connections.page1.response.golden")

	require.NoError(t, err)
	assert.Equal(t, 2, gets)
	assert.Zero(t, posts)
	assert.Zero(t, deletes)
	assert.Contains(t, prompt, "Design 'SampleApp' will be undeployed from 3 connected Kubernetes cluster(s): east, west, north.")
	assert.Contains(t, output, "Undeploy cancelled")
}

func runRemovalPromptTest(t *testing.T, args []string, input string) (output, prompt string, gets, posts, deletes int, err error) {
	t.Helper()
	return runRemovalPromptTestWithClusterPages(t, args, input, "removal.connections.response.golden")
}

func runRemovalPromptTestWithClusterPages(t *testing.T, args []string, input string, clusterPages ...string) (output, prompt string, gets, posts, deletes int, err error) {
	t.Helper()
	ctx := utils.InitTestEnvironment(t)
	t.Cleanup(func() {
		utils.StopMockery(t)
		utils.ResetCommandFlags(DesignCmd, t)
		resetVariables()
	})
	utils.TokenFlag = utils.GetToken(t)

	clusterURLs := make([]string, 0, len(clusterPages))
	for i, page := range clusterPages {
		pageURL := ctx.BaseURL + "/" + connectedClustersPath
		if i > 0 {
			pageURL = fmt.Sprintf("%s&page=%d", pageURL, i)
		}
		clusterURLs = append(clusterURLs, pageURL)
		httpmock.RegisterResponder("GET", pageURL, httpmock.NewStringResponder(200,
			utils.NewGoldenFile(t, page, "fixtures").Load()))
	}
	patternURL := ctx.BaseURL + "/api/pattern"
	deleteURL := ctx.BaseURL + "/api/pattern/deploy?contexts=all"
	httpmock.RegisterResponder("POST", patternURL, httpmock.NewStringResponder(200,
		utils.NewGoldenFile(t, "apply.designSave.response.golden", "fixtures").Load()))
	httpmock.RegisterResponder("DELETE", deleteURL, httpmock.NewStringResponder(200,
		utils.NewGoldenFile(t, args[0]+".response.golden", "fixtures").Load()))

	buf := utils.SetupMeshkitLoggerTesting(t, false)
	readPrompt := capturePrompt(t, input)
	DesignCmd.SetArgs(args)
	DesignCmd.SetOut(buf)
	err = DesignCmd.Execute()
	calls := httpmock.GetCallCountInfo()
	for _, pageURL := range clusterURLs {
		gets += calls["GET "+pageURL]
	}
	return buf.String(), readPrompt(), gets, calls["POST "+patternURL], calls["DELETE "+deleteURL], err
}
