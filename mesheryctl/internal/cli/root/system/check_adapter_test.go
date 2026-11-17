package system

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TestRunAdapterHealthChecksUnknownAdapter verifies that checking an adapter
// that is not configured fails instead of reporting that nothing is configured,
// and that a configured adapter is still checked.
func TestRunAdapterHealthChecksUnknownAdapter(t *testing.T) {
	setupContextTestEnv(t)
	t.Cleanup(BreakupFunc)
	origToken := utils.TokenFlag
	utils.TokenFlag = utils.GetToken(t)
	t.Cleanup(func() { utils.TokenFlag = origToken })

	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	httpmock.RegisterResponder(http.MethodGet, "http://localhost:9081/api/system/sync",
		httpmock.NewStringResponder(http.StatusOK, `{"meshAdapters":[{"adapterLocation":"meshery-istio:10000"}]}`))
	httpmock.RegisterResponder(http.MethodGet, "http://localhost:9081/api/system/adapters",
		httpmock.NewStringResponder(http.StatusOK, `{}`))
	const adapterCall = "GET http://localhost:9081/api/system/adapters"

	tests := []struct {
		name         string
		adapter      string
		expectError  bool
		adapterCalls int
	}{
		{name: "unknown adapter", adapter: "meshery-unknown", expectError: true, adapterCalls: 0},
		{name: "configured adapter by name", adapter: "meshery-istio", expectError: false, adapterCalls: 1},
		{name: "configured adapter by location", adapter: "meshery-istio:10000", expectError: false, adapterCalls: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.ZeroCallCounters()
			hc := &HealthChecker{Options: &HealthCheckOptions{PrintLogs: false}, mctlCfg: mctlCfg}

			err := hc.runAdapterHealthChecks(tt.adapter)
			if tt.expectError {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.adapter)
					assert.Contains(t, err.Error(), "is not configured in the current context")
				}
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.adapterCalls, httpmock.GetCallCountInfo()[adapterCall])
		})
	}
}
