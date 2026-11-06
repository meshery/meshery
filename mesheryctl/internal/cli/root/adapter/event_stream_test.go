package adapter

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// eventStreamResponder answers like an HTTP endpoint with the given status,
// content type and body.
func eventStreamResponder(status int, contentType, body string) httpmock.Responder {
	return func(*http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(status, body)
		resp.Header.Set("Content-Type", contentType)
		return resp, nil
	}
}

// waitWithDeadline fails the test instead of letting a wait function run until
// its 20 minute timer.
func waitWithDeadline(t *testing.T, wait func() error) error {
	t.Helper()

	done := make(chan error, 1)
	go func() { done <- wait() }()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		require.FailNow(t, "still waiting after the event stream ended")
		return nil
	}
}

// TestWaitForResponseOnEventStream verifies that the wait functions report the
// result of an operation, and fail right away when no result can arrive.
func TestWaitForResponseOnEventStream(t *testing.T) {
	const sse = "text/event-stream"

	type waitFunc func(*config.MesheryCtlConfig, string) (string, error)
	commands := []struct {
		name       string
		wait       waitFunc
		query      string
		success    string
		failure    string
		streamCode string
		failedCode string
	}{
		{
			name:       "deploy",
			wait:       waitForDeployResponse,
			query:      "cli_deploy",
			success:    `data: {"summary":"done","details":"mesh is now installed"}` + "\n\n",
			failure:    `data: {"summary":"failed","details":"Error installing"}` + "\n\n",
			streamCode: ErrCreatingDeployResponseStreamCode,
			failedCode: ErrFailedDeployingMeshCode,
		},
		{
			name:       "validate",
			wait:       waitForValidateResponse,
			query:      "cli_validate",
			success:    `data: {"summary":"mesh is now installed","details":"passed"}` + "\n\n",
			failure:    `data: {"summary":"failed","details":"error running tests"}` + "\n\n",
			streamCode: ErrCreatingValidateResponseStreamCode,
			failedCode: ErrSMIConformanceTestsFailedCode,
		},
	}

	for _, command := range commands {
		tests := []struct {
			name        string
			status      int
			contentType string
			body        string
			wantCode    string
		}{
			{name: "result arrives", status: http.StatusOK, contentType: sse, body: command.success},
			{name: "operation fails", status: http.StatusOK, contentType: sse, body: command.failure, wantCode: command.failedCode},
			{name: "stream ends without a result", status: http.StatusOK, contentType: sse, body: `data: {"summary":"working"}` + "\n\n", wantCode: command.streamCode},
			{name: "response is not an event stream", status: http.StatusOK, contentType: "text/html", body: "<html></html>\n", wantCode: command.streamCode},
			{name: "stream can't be opened", status: http.StatusNotFound, contentType: "text/plain", body: "404\n", wantCode: command.streamCode},
		}

		for _, tt := range tests {
			t.Run(command.name+"/"+tt.name, func(t *testing.T) {
				testContext := utils.InitTestEnvironment(t)
				defer utils.StopMockery(t)
				_ = utils.SetupMeshkitLoggerTesting(t, false)
				utils.TokenFlag = utils.GetToken(t)

				mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
				require.NoError(t, err)

				httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/system/adapter/events?client="+command.query,
					eventStreamResponder(tt.status, tt.contentType, tt.body))

				err = waitWithDeadline(t, func() error {
					_, err := command.wait(mctlCfg, "mesh is now installed")
					return err
				})

				if tt.wantCode == "" {
					assert.NoError(t, err)
					return
				}
				assert.Equal(t, tt.wantCode, errors.GetCode(err))
			})
		}
	}
}
