package design

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	githubBlobURL = "https://github.com/meshery/meshery/blob/master/designs/sample.yaml"
	githubRawURL  = "https://raw.githubusercontent.com/meshery/meshery/master/designs/sample.yaml"
)

// designRequests records the requests a design command sends to the server.
type designRequests struct {
	importBody map[string]any
	deployBody []byte
	calls      []string
}

// sampleDesign returns the content of the sample design fixture.
func sampleDesign(t *testing.T) []byte {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "not able to get current working directory")
	content, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "fixtures", "sampleDesign.golden"))
	require.NoError(t, err)
	return content
}

// mockDesignServer serves the design endpoints and the GitHub raw file, and
// fails any other request.
func mockDesignServer(t *testing.T, baseURL string, content []byte) *designRequests {
	t.Helper()

	requests := &designRequests{}
	record := func(name string, status int, body string) httpmock.Responder {
		return func(req *http.Request) (*http.Response, error) {
			requests.calls = append(requests.calls, name)
			if req.Body != nil {
				payload, err := io.ReadAll(req.Body)
				if err != nil {
					return nil, err
				}
				switch name {
				case "import":
					if err := json.Unmarshal(payload, &requests.importBody); err != nil {
						return nil, err
					}
				case "deploy", "undeploy":
					requests.deployBody = payload
				}
			}
			return httpmock.NewStringResponse(status, body), nil
		}
	}

	httpmock.RegisterResponder("POST", baseURL+"/api/pattern", record("save", http.StatusOK, `[{"name":""}]`))
	httpmock.RegisterResponder("POST", baseURL+"/api/pattern/import", record("import", http.StatusOK, `[{"name":"saved"}]`))
	httpmock.RegisterResponder("POST", baseURL+"/api/pattern/deploy", record("deploy", http.StatusOK, `{}`))
	httpmock.RegisterResponder("DELETE", baseURL+"/api/pattern/deploy", record("undeploy", http.StatusOK, `{}`))
	httpmock.RegisterResponder("GET", githubRawURL, record("github", http.StatusOK, string(content)))
	return requests
}

// TestApplyRunsAgainstTheServer verifies which requests design apply sends, so
// a design reaches the server through an endpoint that reads it.
func TestApplyRunsAgainstTheServer(t *testing.T) {
	content := sampleDesign(t)
	_, filename, _, _ := runtime.Caller(0)
	localFile := filepath.Join(filepath.Dir(filename), "fixtures", "sampleDesign.golden")

	tests := []struct {
		name      string
		file      string
		skipSave  bool
		calls     []string
		savedName string
	}{
		{name: "local file is saved then deployed", file: localFile, calls: []string{"import", "deploy"}, savedName: "sampleDesign.golden"},
		{name: "local file with --skip-save is only deployed", file: localFile, skipSave: true, calls: []string{"deploy"}},
		{name: "github url is downloaded, saved then deployed", file: githubBlobURL, calls: []string{"github", "import", "deploy"}, savedName: "sample.yaml"},
		{name: "github url with --skip-save is downloaded and deployed", file: githubBlobURL, skipSave: true, calls: []string{"github", "deploy"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			utils.TokenFlag = utils.GetToken(t)
			requests := mockDesignServer(t, testContext.BaseURL, content)

			originalFile, originalSkipSave := file, skipSave
			defer func() { file, skipSave = originalFile, originalSkipSave }()
			file, skipSave = tt.file, tt.skipSave

			require.NoError(t, applyCmd.RunE(applyCmd, nil))

			assert.Equal(t, tt.calls, requests.calls)
			var deployed struct {
				PatternFile string `json:"patternFile"`
			}
			require.NoError(t, json.Unmarshal(requests.deployBody, &deployed))
			assert.Equal(t, string(content), deployed.PatternFile)

			if tt.savedName != "" {
				assert.Equal(t, tt.savedName, requests.importBody["name"])
				assert.Equal(t, tt.savedName, requests.importBody["fileName"])
				assert.NotEmpty(t, requests.importBody["file"])
			}
		})
	}
}

// TestDeleteWithGithubURLUndeploysTheDownloadedDesign verifies that a design
// given by URL is undeployed from its content, without saving anything.
func TestDeleteWithGithubURLUndeploysTheDownloadedDesign(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	utils.TokenFlag = utils.GetToken(t)

	content := sampleDesign(t)
	requests := mockDesignServer(t, testContext.BaseURL, content)

	originalFile := file
	defer func() { file = originalFile }()
	file = githubBlobURL

	require.NoError(t, deleteCmd.RunE(deleteCmd, nil))

	assert.Equal(t, []string{"github", "undeploy"}, requests.calls)
	assertUndeployPayload(t, requests.deployBody, content)
}

// TestDeleteWithLocalFileUndeploysWithDeployPayload verifies that a design
// given as a local file reaches DELETE /api/pattern/deploy as the JSON payload
// the server decodes, whether it is passed with --file or positionally.
func TestDeleteWithLocalFileUndeploysWithDeployPayload(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	localFile := filepath.Join(filepath.Dir(filename), "fixtures", "sampleDesign.golden")

	tests := []struct {
		name string
		file string
		args []string
	}{
		{name: "file flag", file: localFile},
		{name: "positional file argument", args: []string{localFile}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContext := utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			utils.TokenFlag = utils.GetToken(t)

			content := sampleDesign(t)
			requests := mockDesignServer(t, testContext.BaseURL, content)

			originalFile := file
			defer func() { file = originalFile }()
			file = tt.file

			require.NoError(t, deleteCmd.RunE(deleteCmd, tt.args))

			assert.Equal(t, []string{"undeploy"}, requests.calls)
			assertUndeployPayload(t, requests.deployBody, content)
		})
	}
}

// assertUndeployPayload checks that body is the JSON deploy payload the server's
// PatternFileHandler decodes, carrying content as its patternFile. A raw design
// body is rejected by the server with HTTP 400 (meshery-server-1049).
func assertUndeployPayload(t *testing.T, body, content []byte) {
	t.Helper()

	var payload struct {
		PatternFile string `json:"patternFile"`
	}
	require.NoError(t, json.Unmarshal(body, &payload), "undeploy body must be a JSON deploy payload, got: %.80q", body)
	assert.Equal(t, string(content), payload.PatternFile)
}

// TestFetchDesignFromGithub verifies that both GitHub URL forms resolve to the
// raw file, that no Meshery token is sent to GitHub, and that failures surface.
func TestFetchDesignFromGithub(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	var authorization []string
	httpmock.RegisterResponder("GET", githubRawURL, func(req *http.Request) (*http.Response, error) {
		authorization = append(authorization, req.Header.Get("Authorization"), req.Header.Get("Cookie"))
		return httpmock.NewStringResponse(http.StatusOK, "design content"), nil
	})

	for _, fileURL := range []string{githubBlobURL, githubRawURL} {
		content, err := fetchDesignFromGithub(fileURL)
		require.NoError(t, err)
		assert.Equal(t, "design content", string(content))
	}
	assert.Equal(t, []string{"", "", "", ""}, authorization)

	httpmock.RegisterResponder("GET", githubRawURL, httpmock.NewStringResponder(http.StatusNotFound, "404: Not Found"))
	_, err := fetchDesignFromGithub(githubRawURL)
	assert.Error(t, err, "a missing file must not be treated as design content")

	_, err = fetchDesignFromGithub("https://example.com/design.yaml")
	assert.Error(t, err, "only GitHub URLs are supported")
}
