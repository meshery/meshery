package handlers

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models/telemetry"
	"github.com/meshery/meshery/server/models/telemetry/grafana"
	"github.com/meshery/meshery/server/models/telemetry/prometheus"
	meshkiterrors "github.com/meshery/meshkit/errors"
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
			t.Fatalf("%s missing size diagnosis or remedy: status %d, body %+v", tc.name, response.Code, body)
		}
	}
}

func TestTelemetryAuthPrecedesResponseSize(t *testing.T) {
	for _, code := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(code)
			_, _ = w.Write([]byte(strings.Repeat("x", 16<<20+1)))
		}))
		_, grafanaErr := grafana.New(srv.URL, "", nil).Health(context.Background())
		_, prometheusErr := prometheus.New(srv.URL, "", nil).Health(context.Background())
		srv.Close()
		if grafanaErr == nil || prometheusErr == nil {
			t.Fatalf("HTTP %d: expected both clients to reject the oversized response", code)
		}
		if got := meshkiterrors.GetCode(ErrTelemetryGrafana(grafanaErr, "ping")); !stderrors.Is(grafanaErr, telemetry.ErrResponseTooLarge) || got != ErrTelemetryGrafanaAuthCode {
			t.Errorf("grafana HTTP %d: error %v mapped to %s, want %s", code, grafanaErr, got, ErrTelemetryGrafanaAuthCode)
		}
		if got := meshkiterrors.GetCode(ErrTelemetryPrometheus(prometheusErr, "ping")); !stderrors.Is(prometheusErr, telemetry.ErrResponseTooLarge) || got != ErrTelemetryPrometheusAuthCode {
			t.Errorf("prometheus HTTP %d: error %v mapped to %s, want %s", code, prometheusErr, got, ErrTelemetryPrometheusAuthCode)
		}
	}
}
