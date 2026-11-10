package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	_events "github.com/meshery/meshkit/models/events"
)

type telemetrySpyProvider struct {
	models.DefaultLocalProvider
	persistedEvents []*_events.Event
}

func newTelemetrySpyProvider() *telemetrySpyProvider {
	return &telemetrySpyProvider{
		persistedEvents: make([]*_events.Event, 0),
	}
}

func (s *telemetrySpyProvider) PersistEvent(event _events.Event, _ string) (models.PersistedEvent, error) {
	s.persistedEvents = append(s.persistedEvents, &event)
	return models.PersistedEvent{}, nil
}

// TestPrometheusTelemetryHandlers_RejectsInvalidConnectionID pins that malformed,
// nil, or empty connection ID path parameters return 400 Bad Request with
// ErrInvalidUUIDValueCode and never trigger provider lookups or ping failure events.
func TestPrometheusTelemetryHandlers_RejectsInvalidConnectionID(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   string
	}{
		{"absent", ""},
		{"malformed", "not-a-uuid"},
		{"explicit nil uuid", "00000000-0000-0000-0000-000000000000"},
		{"too short", "1234"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			systemID := uuid.Must(uuid.NewV4())
			h := &Handler{
				config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
				log:      newTestLogger(t),
				SystemID: &systemID,
			}
			provider := newTelemetrySpyProvider()

			// Test Ping Handler
			pingReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+tc.id+"/ping", nil)
			if tc.id != "" {
				pingReq = mux.SetURLVars(pingReq, map[string]string{"connectionID": tc.id})
			}
			pingRec := httptest.NewRecorder()
			h.PrometheusTelemetryPingHandler(pingRec, pingReq, nil, &models.User{}, provider)

			if pingRec.Code != http.StatusBadRequest {
				t.Errorf("Prometheus Ping expected 400, got %d (body=%q)", pingRec.Code, pingRec.Body.String())
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Prometheus Ping unexpectedly persisted failure event for invalid ID: %+v", provider.persistedEvents)
			}
			if decoded := decodeErrorBody(t, pingRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Prometheus Ping expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}

			// Test Metrics Handler
			metricsReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+tc.id+"/metrics", nil)
			if tc.id != "" {
				metricsReq = mux.SetURLVars(metricsReq, map[string]string{"connectionID": tc.id})
			}
			metricsRec := httptest.NewRecorder()
			h.PrometheusTelemetryMetricsHandler(metricsRec, metricsReq, nil, &models.User{}, provider)

			if metricsRec.Code != http.StatusBadRequest {
				t.Errorf("Prometheus Metrics expected 400, got %d (body=%q)", metricsRec.Code, metricsRec.Body.String())
			}
			if decoded := decodeErrorBody(t, metricsRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Prometheus Metrics expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}

			// Test Panels Handler
			panelsReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+tc.id+"/panels", nil)
			if tc.id != "" {
				panelsReq = mux.SetURLVars(panelsReq, map[string]string{"connectionID": tc.id})
			}
			panelsRec := httptest.NewRecorder()
			h.PrometheusTelemetryPanelsHandler(panelsRec, panelsReq, nil, &models.User{}, provider)

			if panelsRec.Code != http.StatusBadRequest {
				t.Errorf("Prometheus Panels expected 400, got %d (body=%q)", panelsRec.Code, panelsRec.Body.String())
			}
			if decoded := decodeErrorBody(t, panelsRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Prometheus Panels expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
		})
	}
}

// TestGrafanaTelemetryHandlers_RejectsInvalidConnectionID pins that malformed,
// nil, or empty connection ID path parameters return 400 Bad Request with
// ErrInvalidUUIDValueCode and never trigger provider lookups or ping failure events.
func TestGrafanaTelemetryHandlers_RejectsInvalidConnectionID(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   string
	}{
		{"absent", ""},
		{"malformed", "not-a-uuid"},
		{"explicit nil uuid", "00000000-0000-0000-0000-000000000000"},
		{"too short", "1234"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			systemID := uuid.Must(uuid.NewV4())
			h := &Handler{
				config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
				log:      newTestLogger(t),
				SystemID: &systemID,
			}
			provider := newTelemetrySpyProvider()

			// Test Ping Handler
			pingReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+tc.id+"/ping", nil)
			if tc.id != "" {
				pingReq = mux.SetURLVars(pingReq, map[string]string{"connectionID": tc.id})
			}
			pingRec := httptest.NewRecorder()
			h.GrafanaTelemetryPingHandler(pingRec, pingReq, nil, &models.User{}, provider)

			if pingRec.Code != http.StatusBadRequest {
				t.Errorf("Grafana Ping expected 400, got %d (body=%q)", pingRec.Code, pingRec.Body.String())
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Grafana Ping unexpectedly persisted failure event for invalid ID: %+v", provider.persistedEvents)
			}
			if decoded := decodeErrorBody(t, pingRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Grafana Ping expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}

			// Test Boards Handler
			boardsReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+tc.id+"/boards", nil)
			if tc.id != "" {
				boardsReq = mux.SetURLVars(boardsReq, map[string]string{"connectionID": tc.id})
			}
			boardsRec := httptest.NewRecorder()
			h.GrafanaTelemetryBoardsHandler(boardsRec, boardsReq, nil, &models.User{}, provider)

			if boardsRec.Code != http.StatusBadRequest {
				t.Errorf("Grafana Boards expected 400, got %d (body=%q)", boardsRec.Code, boardsRec.Body.String())
			}
			if decoded := decodeErrorBody(t, boardsRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Grafana Boards expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}

			// Test Pinned Boards Handler
			pinnedReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+tc.id+"/pinned", nil)
			if tc.id != "" {
				pinnedReq = mux.SetURLVars(pinnedReq, map[string]string{"connectionID": tc.id})
			}
			pinnedRec := httptest.NewRecorder()
			h.GrafanaTelemetryPinnedBoardsHandler(pinnedRec, pinnedReq, nil, &models.User{}, provider)

			if pinnedRec.Code != http.StatusBadRequest {
				t.Errorf("Grafana Pinned Boards expected 400, got %d (body=%q)", pinnedRec.Code, pinnedRec.Body.String())
			}
			if decoded := decodeErrorBody(t, pinnedRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Grafana Pinned Boards expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
		})
	}
}

// TestPrometheusTelemetryClientForConnection_NilUUID verifies that
// prometheusTelemetryClientForConnection rejects nil UUID with ErrInvalidUUID.
func TestPrometheusTelemetryClientForConnection_NilUUID(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	client, conn, code, err := h.prometheusTelemetryClientForConnection("", uuid.Nil, &models.DefaultLocalProvider{})
	if client != nil || conn != nil {
		t.Errorf("expected nil client and connection, got client=%v, conn=%v", client, conn)
	}
	if code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, code)
	}
	if err == nil || err.Error() == "" {
		t.Fatal("expected non-nil error")
	}
}

// TestGrafanaClientForConnection_NilUUID verifies that
// grafanaClientForConnection rejects nil UUID with ErrInvalidUUID.
func TestGrafanaClientForConnection_NilUUID(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	client, conn, code, err := h.grafanaClientForConnection("", uuid.Nil, &models.DefaultLocalProvider{})
	if client != nil || conn != nil {
		t.Errorf("expected nil client and connection, got client=%v, conn=%v", client, conn)
	}
	if code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, code)
	}
	if err == nil || err.Error() == "" {
		t.Fatal("expected non-nil error")
	}
}
