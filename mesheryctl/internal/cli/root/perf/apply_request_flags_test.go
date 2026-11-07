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

const requestFlagsProfileID = "906f8876-33b5-4a97-906e-7a409d3b8ae9"

// requestFlagCalls records what perf apply sent to the server.
type requestFlagCalls struct {
	saved map[string]any
	run   url.Values
}

// mockRequestFlagServer serves the given profile page and records the requests
// that save a profile and run a test.
func mockRequestFlagServer(t *testing.T, baseURL, profilePage string) *requestFlagCalls {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "not able to get current working directory")
	origToken := utils.TokenFlag
	utils.TokenFlag = filepath.Join(filepath.Dir(filename), "fixtures", "auth.json")
	t.Cleanup(func() { utils.TokenFlag = origToken })

	calls := &requestFlagCalls{}
	profilesURL := baseURL + "/api/performance/profiles"
	httpmock.RegisterResponder(http.MethodGet, profilesURL, httpmock.NewStringResponder(http.StatusOK, profilePage))
	httpmock.RegisterResponder(http.MethodPost, profilesURL, func(req *http.Request) (*http.Response, error) {
		if err := json.NewDecoder(req.Body).Decode(&calls.saved); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(http.StatusOK, `{"id":"`+requestFlagsProfileID+`","name":"new-profile"}`), nil
	})
	httpmock.RegisterResponder(http.MethodGet, baseURL+"/api/user/performance/profiles/"+requestFlagsProfileID+"/run",
		func(req *http.Request) (*http.Response, error) {
			calls.run = req.URL.Query()
			return httpmock.NewStringResponse(http.StatusOK, `{}`), nil
		})

	return calls
}

// setRequestFlags sets the flags through cobra, so the test also covers how they
// are registered and parsed.
func setRequestFlags(t *testing.T) {
	t.Helper()

	require.NoError(t, applyCmd.Flags().Set("header", "Authorization: Bearer a:b"))
	require.NoError(t, applyCmd.Flags().Set("header", "X-Trace: 1"))
	require.NoError(t, applyCmd.Flags().Set("cookie", "session=abc=def"))
	require.NoError(t, applyCmd.Flags().Set("content-type", "application/json"))
}

// TestApplyRequestFlagsReplaceSavedProfileValues verifies that the flags win over
// the saved profile for this run and, like --body, leave the profile unchanged.
func TestApplyRequestFlagsReplaceSavedProfileValues(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	profilePage := `{"profiles":[{"id":"` + requestFlagsProfileID + `","name":"existing","qps":7,"duration":"30s",` +
		`"endpoints":["https://meshery.io"],"loadGenerators":["fortio"],` +
		`"requestHeaders":"{\"X-Saved\":\"1\"}","requestCookies":"{\"saved\":\"1\"}","contentType":"text/plain"}],"totalCount":1}`
	calls := mockRequestFlagServer(t, testContext.BaseURL, profilePage)

	resetVariables()
	defer resetVariables()
	setRequestFlags(t)

	require.NoError(t, applyCmd.RunE(applyCmd, []string{"existing"}))

	require.NotNil(t, calls.run, "the run request was never sent")
	assert.JSONEq(t, `{"Authorization":"Bearer a:b","X-Trace":"1"}`, calls.run.Get("headers"))
	assert.JSONEq(t, `{"session":"abc=def"}`, calls.run.Get("cookies"))
	assert.Equal(t, "application/json", calls.run.Get("contentType"))
	assert.Nil(t, calls.saved, "the saved profile must not be changed")
}

// TestApplySavesRequestFlagsOnNewProfile verifies that a new profile keeps the
// flag values, so later runs of it send them too.
func TestApplySavesRequestFlagsOnNewProfile(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	defer utils.StopMockery(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	calls := mockRequestFlagServer(t, testContext.BaseURL, `{"profiles":[],"totalCount":0}`)

	resetVariables()
	defer resetVariables()
	origSilent := utils.SilentFlag
	utils.SilentFlag = true
	defer func() { utils.SilentFlag = origSilent }()
	testURL = "https://meshery.io"
	setRequestFlags(t)

	require.NoError(t, applyCmd.RunE(applyCmd, []string{"new-profile"}))

	require.NotNil(t, calls.saved, "the profile was never saved")
	assert.JSONEq(t, `{"Authorization":"Bearer a:b","X-Trace":"1"}`, calls.saved["requestHeaders"].(string))
	assert.JSONEq(t, `{"session":"abc=def"}`, calls.saved["requestCookies"].(string))
	assert.Equal(t, "application/json", calls.saved["contentType"])

	require.NotNil(t, calls.run, "the run request was never sent")
	assert.JSONEq(t, `{"Authorization":"Bearer a:b","X-Trace":"1"}`, calls.run.Get("headers"))
	assert.JSONEq(t, `{"session":"abc=def"}`, calls.run.Get("cookies"))
	assert.Equal(t, "application/json", calls.run.Get("contentType"))
}

// TestApplyRejectsMalformedRequestFlags verifies that a bad value stops the
// command before it sends anything.
func TestApplyRejectsMalformedRequestFlags(t *testing.T) {
	tests := []struct {
		name  string
		flag  string
		value string
	}{
		{name: "header without a colon", flag: "header", value: "Authorization"},
		{name: "header without a name", flag: "header", value: ": value"},
		{name: "cookie without an equals sign", flag: "cookie", value: "session"},
		{name: "cookie without a name", flag: "cookie", value: "=value"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_ = utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			resetVariables()
			defer resetVariables()
			require.NoError(t, applyCmd.Flags().Set(tc.flag, tc.value))

			err := applyCmd.RunE(applyCmd, []string{"existing"})

			require.Error(t, err)
			assert.Contains(t, err.Error(), "--"+tc.flag)
			assert.Zero(t, httpmock.GetTotalCallCount(), "nothing may be sent for a bad value")
		})
	}
}

// TestApplyRejectsDuplicateRequestFlags verifies that a repeated name stops the
// command before it sends anything, instead of silently keeping the last value.
func TestApplyRejectsDuplicateRequestFlags(t *testing.T) {
	tests := []struct {
		name   string
		flag   string
		first  string
		second string
	}{
		{name: "duplicate header", flag: "header", first: "X-Trace: 1", second: "X-Trace: 2"},
		{name: "duplicate header in a different case", flag: "header", first: "Authorization: one", second: "authorization: two"},
		{name: "duplicate cookie", flag: "cookie", first: "session=abc", second: "session=def"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_ = utils.InitTestEnvironment(t)
			defer utils.StopMockery(t)
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			resetVariables()
			defer resetVariables()
			require.NoError(t, applyCmd.Flags().Set(tc.flag, tc.first))
			require.NoError(t, applyCmd.Flags().Set(tc.flag, tc.second))

			err := applyCmd.RunE(applyCmd, []string{"existing"})

			require.Error(t, err)
			assert.Contains(t, err.Error(), "--"+tc.flag)
			assert.Contains(t, err.Error(), "duplicate")
			assert.Zero(t, httpmock.GetTotalCallCount(), "nothing may be sent for a duplicate value")
		})
	}
}

// TestFlagValuesToJSONKeepsCookieNamesCaseSensitive verifies that only header
// names are compared without regard to case. Cookie names are case-sensitive, so
// two that differ only in case are different cookies.
func TestFlagValuesToJSONKeepsCookieNamesCaseSensitive(t *testing.T) {
	encoded, err := flagValuesToJSON("cookie", "=", []string{"session=abc", "Session=def"})

	require.NoError(t, err)
	assert.JSONEq(t, `{"session":"abc","Session":"def"}`, encoded)
}
