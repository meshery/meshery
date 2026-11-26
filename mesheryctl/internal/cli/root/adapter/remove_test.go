package adapter

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveMesh(t *testing.T) {
	utils.SetupContextEnv(t)

	// initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		Name             string
		Args             []string
		URLs             []utils.MockURL
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name: "Test Remove Linkerd ",
			Args: []string{"remove", "--yes", "linkerd"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation?contexts=all",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.linkerd.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name: "Test Remove Linkerd with Namespace",
			Args: []string{"remove", "--yes", "linkerd", "--namespace", "linkerd-ns"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation?contexts=all",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.linkerd-namespace.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
		{
			Name: "Test Remove Cilium ",
			Args: []string{"remove", "--yes", "cilium", "service", "mesh"},
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/system/sync",
					Response:     "sync.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/system/adapter/operation?contexts=all",
					Response:     "deploy.golden",
					ResponseCode: 200,
				},
			},
			ExpectedResponse: "remove.cilium.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			t.Cleanup(func() { utils.ResetCommandFlags(AdapterCmd, t) })
			httpmock.RegisterResponder("GET", testContext.BaseURL+"/"+removeClustersPath,
				httpmock.NewStringResponder(200, connectedRemoveClusters))
			for _, url := range tc.URLs {
				apiResponse := utils.NewGoldenFile(t, url.Response, fixturesDir).Load()
				// mock response
				httpmock.RegisterResponder(url.Method, url.URL,
					httpmock.NewStringResponder(url.ResponseCode, apiResponse))
			}
			utils.TokenFlag = tc.Token
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tc.ExpectedResponse, testdataDir)
			buff := utils.SetupMeshkitLoggerTesting(t, false)
			cmd := AdapterCmd
			cmd.SetArgs(tc.Args)
			cmd.SetOut(buff)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tc.ExpectError {
					// write it in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}
			// response being printed in console
			actualResponse := buff.String()

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			assert.Equal(t, expectedResponse, actualResponse)
		})
		t.Log("Mesh Deploy test Passed")
	}
	utils.StopMockery(t)
}

const (
	removeClustersPath      = "api/integrations/connections?kind=kubernetes&status=connected&pagesize=all"
	connectedRemoveClusters = `{"totalCount":2,"connections":[{"name":"east"},{"name":"west"}]}`
)

func TestRemoveMeshConfirmation(t *testing.T) {
	const removalPrompt = "Adapter 'linkerd' will be removed from 2 connected Kubernetes cluster(s): east, west."
	for _, tc := range []struct {
		name, input, clusters, flag, wantError string
		wantPosts                              int
		wantPrompt                             bool
	}{
		{name: "decline", input: "n\n", clusters: connectedRemoveClusters, wantPrompt: true},
		{name: "yes", flag: "--yes", clusters: connectedRemoveClusters, wantPosts: 1},
		{name: "zero clusters", flag: "-y", clusters: `{"totalCount":0,"connections":[]}`, wantError: "No Kubernetes cluster is connected"},
		{name: "EOF", clusters: connectedRemoveClusters, wantError: "--yes", wantPrompt: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := utils.InitTestEnvironment(t)
			t.Cleanup(func() {
				utils.StopMockery(t)
				utils.ResetCommandFlags(AdapterCmd, t)
			})
			utils.TokenFlag = utils.GetToken(t)
			clusterURL := ctx.BaseURL + "/" + removeClustersPath
			operationURL := ctx.BaseURL + "/api/system/adapter/operation?contexts=all"
			httpmock.RegisterResponder("GET", ctx.BaseURL+"/api/system/sync", httpmock.NewStringResponder(200,
				utils.NewGoldenFile(t, "sync.golden", "fixtures").Load()))
			httpmock.RegisterResponder("GET", clusterURL, httpmock.NewStringResponder(200, tc.clusters))
			httpmock.RegisterResponder("POST", operationURL, httpmock.NewStringResponder(200,
				utils.NewGoldenFile(t, "deploy.golden", "fixtures").Load()))

			input, err := os.CreateTemp(t.TempDir(), "stdin")
			require.NoError(t, err)
			originalStdin := os.Stdin
			os.Stdin = input
			t.Cleanup(func() {
				os.Stdin = originalStdin
				_ = input.Close()
			})
			_, err = input.WriteString(tc.input)
			require.NoError(t, err)
			_, err = input.Seek(0, 0)
			require.NoError(t, err)

			args := []string{"remove", "linkerd"}
			if tc.flag != "" {
				args = append(args, tc.flag)
			}
			AdapterCmd.SetArgs(args)
			AdapterCmd.SetOut(utils.SetupMeshkitLoggerTesting(t, false))

			stdoutR, stdoutW, err := os.Pipe()
			require.NoError(t, err)
			originalStdout := os.Stdout
			os.Stdout = stdoutW
			t.Cleanup(func() {
				os.Stdout = originalStdout
				_ = stdoutR.Close()
				_ = stdoutW.Close()
			})

			err = AdapterCmd.Execute()
			os.Stdout = originalStdout
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.wantError)
			}
			require.NoError(t, stdoutW.Close())
			out, err := io.ReadAll(stdoutR)
			require.NoError(t, err)
			if tc.wantPrompt {
				assert.Contains(t, string(out), removalPrompt)
			} else {
				assert.Empty(t, string(out))
			}

			calls := httpmock.GetCallCountInfo()
			assert.Equal(t, 1, calls["GET "+clusterURL])
			assert.Equal(t, tc.wantPosts, calls["POST "+operationURL])
		})
	}
}
