package perf

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApplySendsSavedRequestSettings runs a saved profile and checks the run
// request. The server reads headers, cookies and the content type only from this
// query string, so anything left out is dropped from the load test.
func TestApplySendsSavedRequestSettings(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "not able to get current working directory")
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(filepath.Dir(filename), "fixtures", "auth.json")
	defer func() { utils.TokenFlag = origToken }()

	const (
		profileID   = "906f8876-33b5-4a97-906e-7a409d3b8ae9"
		headers     = `{"Authorization":"Bearer token"}`
		cookies     = `{"session":"abc"}`
		contentType = "application/json"
		requestBody = `{"item":1}`
	)

	page, err := json.Marshal(map[string]any{
		"totalCount": 1,
		"profiles": []map[string]any{{
			"id":             profileID,
			"name":           "existing",
			"qps":            7,
			"duration":       "30s",
			"endpoints":      []string{"https://meshery.io"},
			"loadGenerators": []string{"fortio"},
			"requestHeaders": headers,
			"requestCookies": cookies,
			"contentType":    contentType,
			"requestBody":    requestBody,
		}},
	})
	require.NoError(t, err)

	var query url.Values
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/performance/profiles",
		httpmock.NewStringResponder(http.StatusOK, string(page)))
	httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/user/performance/profiles/"+profileID+"/run",
		func(req *http.Request) (*http.Response, error) {
			query = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
		})

	resetVariables()
	defer resetVariables()

	require.NoError(t, applyCmd.RunE(applyCmd, []string{"existing"}))

	require.NotNil(t, query, "the run request was never sent")
	assert.Equal(t, requestBody, query.Get("reqBody"))
	assert.Equal(t, headers, query.Get("headers"))
	assert.Equal(t, cookies, query.Get("cookies"))
	assert.Equal(t, contentType, query.Get("contentType"))
}
