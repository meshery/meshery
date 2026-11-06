package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models/telemetry/grafana"
	"github.com/meshery/meshery/server/models/telemetry/prometheus"
)

func TestTelemetryResponseSizeErrorIsUserFacing(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		wrap func(error, string) error
	}{
		{"grafana", grafana.ErrResponseTooLarge, ErrTelemetryGrafana},
		{"prometheus", prometheus.ErrResponseTooLarge, ErrTelemetryPrometheus},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			writeMeshkitError(response, tc.wrap(fmt.Errorf("upstream: %w", tc.err), "datasource query"), http.StatusBadGateway)
			var body struct {
				Error                string   `json:"error"`
				Code                 string   `json:"code"`
				SuggestedRemediation []string `json:"suggestedRemediation"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != http.StatusBadGateway || body.Code != ErrTelemetrySizeCode || !strings.Contains(body.Error, "16 MiB") ||
				len(body.SuggestedRemediation) != 2 || !strings.Contains(body.SuggestedRemediation[1], "shorten the time range") {
				t.Fatalf("missing size diagnosis or remedy: status %d, body %+v", response.Code, body)
			}
		})
	}
}
