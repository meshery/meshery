package handlers

import (
	"net/http"
	"testing"
	"time"

	serverutils "github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models/telemetry/grafana"
	"github.com/meshery/meshery/server/models/telemetry/prometheus"
	"github.com/meshery/meshkit/logger"
)

// telemetryTestTimeout mirrors the client packages' own request timeout.
const telemetryTestTimeout = 30 * time.Second

// allowPrivateHTTPClient returns a transport permitted to reach loopback, for
// tests that stand up a local httptest server in place of a Grafana or
// Prometheus instance. Production code uses the guarded client that New
// installs; nothing outside _test.go may call this.
func allowPrivateHTTPClient() *http.Client {
	return serverutils.NewSafeHTTPClientWithAllowPrivate(telemetryTestTimeout, true)
}

// useLoopbackTelemetryClients points the handler's telemetry client
// constructors at a loopback-capable transport for the duration of one test,
// so tests that exercise the handler path (rather than calling the client
// directly) can still drive a local stand-in server. The originals are restored
// on cleanup.
func useLoopbackTelemetryClients(t *testing.T) {
	t.Helper()

	origGrafana, origPrometheus := newGrafanaTelemetryClient, newPrometheusTelemetryClient
	newGrafanaTelemetryClient = func(baseURL, secret string, log logger.Handler) *grafana.Client {
		return grafana.NewWithHTTPClient(baseURL, secret, log, allowPrivateHTTPClient())
	}
	newPrometheusTelemetryClient = func(baseURL, secret string, log logger.Handler) *prometheus.Client {
		return prometheus.NewWithHTTPClient(baseURL, secret, log, allowPrivateHTTPClient())
	}
	t.Cleanup(func() {
		newGrafanaTelemetryClient, newPrometheusTelemetryClient = origGrafana, origPrometheus
	})
}
