package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

			handlers := []struct {
				name   string
				path   string
				invoke func(w http.ResponseWriter, req *http.Request)
			}{
				{"Metrics", "/metrics", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryMetricsHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Labels", "/labels", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryLabelsHandler(w, req, nil, &models.User{}, provider)
				}},
				{"LabelValues", "/label_values", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryLabelValuesHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Metadata", "/metadata", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryMetadataHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Query", "/query", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryQueryHandler(w, req, nil, &models.User{}, provider)
				}},
				{"QueryRange", "/query_range", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryQueryRangeHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Panels", "/panels", func(w http.ResponseWriter, req *http.Request) {
					h.PrometheusTelemetryPanelsHandler(w, req, nil, &models.User{}, provider)
				}},
			}

			for _, target := range handlers {
				req := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+tc.id+target.path, nil)
				if tc.id != "" {
					req = mux.SetURLVars(req, map[string]string{"connectionID": tc.id})
				}
				rec := httptest.NewRecorder()
				target.invoke(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Errorf("Prometheus %s expected 400, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
				}
				if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
					t.Errorf("Prometheus %s expected code %q, got %q", target.name, models.ErrInvalidUUIDValueCode, decoded.Code)
				}
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Prometheus handlers unexpectedly persisted events for invalid ID: %+v", provider.persistedEvents)
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

			handlers := []struct {
				name   string
				path   string
				vars   map[string]string
				invoke func(w http.ResponseWriter, req *http.Request)
			}{
				{"Boards", "/boards", nil, func(w http.ResponseWriter, req *http.Request) {
					h.GrafanaTelemetryBoardsHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Board", "/boards/test-uid", map[string]string{"uid": "test-uid"}, func(w http.ResponseWriter, req *http.Request) {
					h.GrafanaTelemetryBoardHandler(w, req, nil, &models.User{}, provider)
				}},
				{"Datasources", "/datasources", nil, func(w http.ResponseWriter, req *http.Request) {
					h.GrafanaTelemetryDatasourcesHandler(w, req, nil, &models.User{}, provider)
				}},
				{"QueryRange", "/query_range", nil, func(w http.ResponseWriter, req *http.Request) {
					h.GrafanaTelemetryQueryRangeHandler(w, req, nil, &models.User{}, provider)
				}},
				{"PinnedBoards", "/pinned", nil, func(w http.ResponseWriter, req *http.Request) {
					h.GrafanaTelemetryPinnedBoardsHandler(w, req, nil, &models.User{}, provider)
				}},
			}

			for _, target := range handlers {
				req := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+tc.id+target.path, nil)
				reqVars := map[string]string{}
				if tc.id != "" {
					reqVars["connectionID"] = tc.id
				}
				for k, v := range target.vars {
					reqVars[k] = v
				}
				if len(reqVars) > 0 {
					req = mux.SetURLVars(req, reqVars)
				}
				rec := httptest.NewRecorder()
				target.invoke(rec, req)

				if rec.Code != http.StatusBadRequest {
					t.Errorf("Grafana %s expected 400, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
				}
				if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
					t.Errorf("Grafana %s expected code %q, got %q", target.name, models.ErrInvalidUUIDValueCode, decoded.Code)
				}
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Grafana handlers unexpectedly persisted events for invalid ID: %+v", provider.persistedEvents)
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

// TestPrometheusTelemetryQueryRangeBatchHandler_RejectsInvalidConnectionID pins
// that the batch handler validates the connection ID before reading the request
// body: an invalid ID with a body that would fail to unmarshal still returns
// 400 with ErrInvalidUUIDValueCode, so a regression that moves validation after
// the body read fails this test.
func TestPrometheusTelemetryQueryRangeBatchHandler_RejectsInvalidConnectionID(t *testing.T) {
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

			batchReq := httptest.NewRequest(http.MethodPost, "/api/telemetry/prometheus/"+tc.id+"/query_range_batch", strings.NewReader("this is not json"))
			if tc.id != "" {
				batchReq = mux.SetURLVars(batchReq, map[string]string{"connectionID": tc.id})
			}
			batchRec := httptest.NewRecorder()
			h.PrometheusTelemetryQueryRangeBatchHandler(batchRec, batchReq, nil, &models.User{}, provider)

			if batchRec.Code != http.StatusBadRequest {
				t.Errorf("Prometheus Batch expected 400, got %d (body=%q)", batchRec.Code, batchRec.Body.String())
			}
			if decoded := decodeErrorBody(t, batchRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Prometheus Batch expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Prometheus Batch unexpectedly persisted failure event for invalid ID: %+v", provider.persistedEvents)
			}
		})
	}
}

// TestGrafanaTelemetryQueryRangeBatchHandler_RejectsInvalidConnectionID pins
// that the batch handler validates the connection ID before reading the request
// body: an invalid ID with a body that would fail to unmarshal still returns
// 400 with ErrInvalidUUIDValueCode, so a regression that moves validation after
// the body read fails this test.
func TestGrafanaTelemetryQueryRangeBatchHandler_RejectsInvalidConnectionID(t *testing.T) {
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

			batchReq := httptest.NewRequest(http.MethodPost, "/api/telemetry/grafana/"+tc.id+"/query_range_batch", strings.NewReader("this is not json"))
			if tc.id != "" {
				batchReq = mux.SetURLVars(batchReq, map[string]string{"connectionID": tc.id})
			}
			batchRec := httptest.NewRecorder()
			h.GrafanaTelemetryQueryRangeBatchHandler(batchRec, batchReq, nil, &models.User{}, provider)

			if batchRec.Code != http.StatusBadRequest {
				t.Errorf("Grafana Batch expected 400, got %d (body=%q)", batchRec.Code, batchRec.Body.String())
			}
			if decoded := decodeErrorBody(t, batchRec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("Grafana Batch expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
			if len(provider.persistedEvents) > 0 {
				t.Errorf("Grafana Batch unexpectedly persisted failure event for invalid ID: %+v", provider.persistedEvents)
			}
		})
	}
}
