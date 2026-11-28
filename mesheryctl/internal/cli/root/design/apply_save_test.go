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
	githubBlobURL    = "https://github.com/meshery/meshery/blob/master/designs/sample.yaml"
	githubRawLinkURL = "https://github.com/meshery/meshery/raw/master/designs/sample.yaml"
	githubRawURL     = "https://raw.githubusercontent.com/meshery/meshery/master/designs/sample.yaml"
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

// TestFetchDesignFromGithub verifies that the blob, raw-link and
// raw.githubusercontent.com URL forms resolve to the raw file, that no Meshery
// token is sent to GitHub, and that failures surface.
func TestFetchDesignFromGithub(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	// A raw link whose file name has encoded characters must reach the
	// encoded raw.githubusercontent.com URL with %20 and %23 intact.
	const encodedRawLinkURL = "https://github.com/meshery/meshery/raw/master/designs/my%20design%23v1.yaml"
	const encodedRawURL = "https://raw.githubusercontent.com/meshery/meshery/master/designs/my%20design%23v1.yaml"

	var authorization []string
	responder := func(req *http.Request) (*http.Response, error) {
		authorization = append(authorization, req.Header.Get("Authorization"), req.Header.Get("Cookie"))
		return httpmock.NewStringResponse(http.StatusOK, "design content"), nil
	}
	httpmock.RegisterResponder("GET", githubRawURL, responder)
	httpmock.RegisterResponder("GET", encodedRawURL, responder)

	for _, fileURL := range []string{githubBlobURL, githubRawLinkURL, encodedRawLinkURL, githubRawURL} {
		content, err := fetchDesignFromGithub(fileURL)
		require.NoError(t, err, fileURL)
		assert.Equal(t, "design content", string(content))
	}
	require.NotEmpty(t, authorization)
	for _, header := range authorization {
		assert.Empty(t, header)
	}

	httpmock.RegisterResponder("GET", githubRawURL, httpmock.NewStringResponder(http.StatusNotFound, "404: Not Found"))
	_, err := fetchDesignFromGithub(githubRawURL)
	assert.Error(t, err, "a missing file must not be treated as design content")

	_, err = fetchDesignFromGithub("https://example.com/design.yaml")
	assert.Error(t, err, "only GitHub URLs are supported")
}

func TestFetchDesignFile_GenericURL(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	genericURL := "https://example.com/designs/sample.yaml"
	httpmock.RegisterResponder("GET", genericURL, httpmock.NewStringResponder(http.StatusOK, "generic design content"))

	content, err := fetchDesignFile(genericURL)
	require.NoError(t, err)
	assert.Equal(t, "generic design content", string(content))

	// 404 response returns error
	httpmock.RegisterResponder("GET", genericURL, httpmock.NewStringResponder(http.StatusNotFound, "404 Not Found"))
	_, err = fetchDesignFile(genericURL)
	assert.Error(t, err)

	// Non-HTTP(S) schemes are rejected cleanly
	_, err = fetchDesignFile("ftp://example.com/designs/sample.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported URL scheme")
}

func TestFetchDesignFile_GithubURL(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("GET", githubRawURL, httpmock.NewStringResponder(http.StatusOK, "github design content"))

	content, err := fetchDesignFile(githubBlobURL)
	require.NoError(t, err)
	assert.Equal(t, "github design content", string(content))

	authedRawURL := githubRawURL + "?token=secret123"
	httpmock.RegisterResponder("GET", authedRawURL, httpmock.NewStringResponder(http.StatusOK, "authed raw content"))
	authedContent, err := fetchDesignFile(authedRawURL)
	require.NoError(t, err)
	assert.Equal(t, "authed raw content", string(authedContent))
}

func TestIsGithubURL(t *testing.T) {
	assert.True(t, isGithubURL(githubBlobURL))
	assert.True(t, isGithubURL(githubRawURL))
	assert.False(t, isGithubURL("https://github.com.evil.com/meshery/design.yaml"))
	assert.False(t, isGithubURL("https://gitlab.com/meshery/design.yaml"))
	assert.False(t, isGithubURL("https://example.com/design.yaml"))
	assert.False(t, isGithubURL("not-a-url"))
}

func TestDesignFileName(t *testing.T) {
	assert.Equal(t, "design.yaml", designFileName("https://example.com/designs/design.yaml?token=secret&sig=123"))
	assert.Equal(t, "design.yaml", designFileName("https://example.com/?token=secret"))
	assert.Equal(t, "design.yaml", designFileName("https://example.com"))
	assert.Equal(t, "sample.yaml", designFileName("https://example.com/sample.yaml"))
	assert.Equal(t, "local.yaml", designFileName("designs/local.yaml"))
}

func TestGithubRawLinkToRawURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "given github raw link when converted then raw.githubusercontent.com URL is returned",
			input: githubRawLinkURL,
			want:  githubRawURL,
		},
		{
			name:  "given github raw link with full ref when converted then ref is kept",
			input: "https://github.com/meshery/meshery/raw/refs/heads/master/sample.yaml",
			want:  "https://raw.githubusercontent.com/meshery/meshery/refs/heads/master/sample.yaml",
		},
		{
			name:  "given github raw link with a raw directory in the path when converted then the directory is kept",
			input: "https://github.com/meshery/meshery/raw/master/raw/sample.yaml",
			want:  "https://raw.githubusercontent.com/meshery/meshery/master/raw/sample.yaml",
		},
		{
			name:  "given github raw link with encoded characters when converted then they stay encoded",
			input: "https://github.com/meshery/meshery/raw/master/my%20design%23v1.yaml",
			want:  "https://raw.githubusercontent.com/meshery/meshery/master/my%20design%23v1.yaml",
		},
		{
			name:  "given github raw link with query when converted then query is kept",
			input: "https://github.com/meshery/meshery/raw/master/sample.yaml?token=abc",
			want:  "https://raw.githubusercontent.com/meshery/meshery/master/sample.yaml?token=abc",
		},
		{
			name:  "given github blob link when converted then it is unchanged",
			input: githubBlobURL,
			want:  githubBlobURL,
		},
		{
			name:  "given repository named raw when converted then it is unchanged",
			input: "https://github.com/meshery/raw/blob/master/sample.yaml",
			want:  "https://github.com/meshery/raw/blob/master/sample.yaml",
		},
		{
			name:  "given raw.githubusercontent.com URL when converted then it is unchanged",
			input: githubRawURL,
			want:  githubRawURL,
		},
		{
			name:  "given github raw link without a file path when converted then it is unchanged",
			input: "https://github.com/meshery/meshery/raw/master",
			want:  "https://github.com/meshery/meshery/raw/master",
		},
		{
			name:  "given non-github URL with a raw segment when converted then it is unchanged",
			input: "https://example.com/meshery/meshery/raw/master/sample.yaml",
			want:  "https://example.com/meshery/meshery/raw/master/sample.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, githubRawLinkToRawURL(tt.input))
		})
	}
}
