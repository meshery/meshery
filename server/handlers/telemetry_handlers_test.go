package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	meshkiterrors "github.com/meshery/meshkit/errors"
	_events "github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// telemetrySpyProvider is a mock provider used for testing telemetry handlers.
type telemetrySpyProvider struct {
	models.DefaultLocalProvider
	persistedEvents []*_events.Event
	getConnectionFn func(token string, connectionID core.Uuid) (*connections.Connection, int, error)
}

// newTelemetrySpyProvider constructs an initialized telemetrySpyProvider.
func newTelemetrySpyProvider() *telemetrySpyProvider {
	return &telemetrySpyProvider{
		persistedEvents: make([]*_events.Event, 0),
	}
}

// PersistEvent records the event in persistedEvents and returns a default PersistedEvent.
func (s *telemetrySpyProvider) PersistEvent(event _events.Event, _ string) (models.PersistedEvent, error) {
	s.persistedEvents = append(s.persistedEvents, &event)
	return models.PersistedEvent{}, nil
}

// GetConnectionByID dispatches to getConnectionFn when non-nil, or delegates to DefaultLocalProvider.
func (s *telemetrySpyProvider) GetConnectionByID(token string, connectionID core.Uuid) (*connections.Connection, int, error) {
	if s.getConnectionFn != nil {
		return s.getConnectionFn(token, connectionID)
	}
	return s.DefaultLocalProvider.GetConnectionByID(token, connectionID)
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

// TestTelemetryPingHandlers_NilUser_ReturnsUnauthorized verifies that both Grafana
// and Prometheus ping handlers return 401 Unauthorized when user is nil, preventing
// nil pointer dereferences.
func TestTelemetryPingHandlers_NilUser_ReturnsUnauthorized(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := newTelemetrySpyProvider()

	// Test Grafana Ping with nil user
	gReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+connID+"/ping", nil)
	gReq = mux.SetURLVars(gReq, map[string]string{"connectionID": connID})
	gRec := httptest.NewRecorder()
	h.GrafanaTelemetryPingHandler(gRec, gReq, nil, nil, provider)

	if gRec.Code != http.StatusUnauthorized {
		t.Errorf("Grafana Ping expected 401 Unauthorized for nil user, got %d (body=%q)", gRec.Code, gRec.Body.String())
	}
	if !strings.Contains(gRec.Body.String(), "user unauthorized") {
		t.Errorf("Grafana Ping expected body to contain 'user unauthorized', got %q", gRec.Body.String())
	}

	// Test Prometheus Ping with nil user
	pReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+connID+"/ping", nil)
	pReq = mux.SetURLVars(pReq, map[string]string{"connectionID": connID})
	pRec := httptest.NewRecorder()
	h.PrometheusTelemetryPingHandler(pRec, pReq, nil, nil, provider)

	if pRec.Code != http.StatusUnauthorized {
		t.Errorf("Prometheus Ping expected 401 Unauthorized for nil user, got %d (body=%q)", pRec.Code, pRec.Body.String())
	}
	if !strings.Contains(pRec.Body.String(), "user unauthorized") {
		t.Errorf("Prometheus Ping expected body to contain 'user unauthorized', got %q", pRec.Body.String())
	}
}

// TestTelemetryPingHandlers_NilSystemID_DoesNotPanic verifies that both Grafana and
// Prometheus ping handlers safely fall back to uuid.Nil when h.SystemID is nil.
func TestTelemetryPingHandlers_NilSystemID_DoesNotPanic(t *testing.T) {
	connUUID := uuid.Must(uuid.NewV4())
	connID := connUUID.String()
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: nil,
	}
	provider := newTelemetrySpyProvider()
	provider.getConnectionFn = func(_ string, _ core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{
			ID:   connUUID,
			Kind: "grafana",
			Name: "test-grafana",
		}, http.StatusOK, nil
	}

	gReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+connID+"/ping", nil)
	gReq = mux.SetURLVars(gReq, map[string]string{"connectionID": connID})
	gRec := httptest.NewRecorder()
	h.GrafanaTelemetryPingHandler(gRec, gReq, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	// Since upstream URL is empty, client.Health will fail, but handler must not panic on nil SystemID
	if len(provider.persistedEvents) > 0 {
		event := provider.persistedEvents[0]
		if event.SystemID != uuid.Nil {
			t.Errorf("expected event SystemID to be uuid.Nil when h.SystemID is nil, got %v", event.SystemID)
		}
	}

	provider.getConnectionFn = func(_ string, _ core.Uuid) (*connections.Connection, int, error) {
		return &connections.Connection{
			ID:   connUUID,
			Kind: "prometheus",
			Name: "test-prometheus",
		}, http.StatusOK, nil
	}
	pReq := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+connID+"/ping", nil)
	pReq = mux.SetURLVars(pReq, map[string]string{"connectionID": connID})
	pRec := httptest.NewRecorder()
	h.PrometheusTelemetryPingHandler(pRec, pReq, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
}

// TestTelemetryHandlers_NilProvider_ReturnsServerError verifies that all Grafana
// and Prometheus telemetry endpoints and client resolvers return 500 Internal Server
// Error when provider is nil, preventing panics.
func TestTelemetryHandlers_NilProvider_ReturnsServerError(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	connUUID := uuid.Must(uuid.NewV4())
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	grafanaHandlers := []struct {
		name   string
		path   string
		vars   map[string]string
		invoke func(w http.ResponseWriter, req *http.Request)
	}{
		{"Ping", "/ping", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryPingHandler(w, req, nil, user, nil)
		}},
		{"Boards", "/boards", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryBoardsHandler(w, req, nil, user, nil)
		}},
		{"Board", "/boards/test-uid", map[string]string{"uid": "test-uid"}, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryBoardHandler(w, req, nil, user, nil)
		}},
		{"Datasources", "/datasources", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryDatasourcesHandler(w, req, nil, user, nil)
		}},
		{"QueryRange", "/query_range", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryQueryRangeHandler(w, req, nil, user, nil)
		}},
		{"PinnedBoards", "/pinned", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryPinnedBoardsHandler(w, req, nil, user, nil)
		}},
	}

	for _, target := range grafanaHandlers {
		t.Run("Grafana_"+target.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+connID+target.path, nil)
			reqVars := map[string]string{"connectionID": connID}
			for k, v := range target.vars {
				reqVars[k] = v
			}
			req = mux.SetURLVars(req, reqVars)
			rec := httptest.NewRecorder()
			target.invoke(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Errorf("Grafana %s expected 500 for nil provider, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
			}
		})
	}

	prometheusHandlers := []struct {
		name   string
		path   string
		invoke func(w http.ResponseWriter, req *http.Request)
	}{
		{"Ping", "/ping", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryPingHandler(w, req, nil, user, nil)
		}},
		{"Metrics", "/metrics", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryMetricsHandler(w, req, nil, user, nil)
		}},
		{"Labels", "/labels", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryLabelsHandler(w, req, nil, user, nil)
		}},
		{"LabelValues", "/label_values", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryLabelValuesHandler(w, req, nil, user, nil)
		}},
		{"Metadata", "/metadata", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryMetadataHandler(w, req, nil, user, nil)
		}},
		{"Query", "/query", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryQueryHandler(w, req, nil, user, nil)
		}},
		{"QueryRange", "/query_range", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryQueryRangeHandler(w, req, nil, user, nil)
		}},
		{"Panels", "/panels", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryPanelsHandler(w, req, nil, user, nil)
		}},
	}

	for _, target := range prometheusHandlers {
		t.Run("Prometheus_"+target.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+connID+target.path, nil)
			req = mux.SetURLVars(req, map[string]string{"connectionID": connID})
			rec := httptest.NewRecorder()
			target.invoke(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Errorf("Prometheus %s expected 500 for nil provider, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
			}
		})
	}

	// Client constructor helpers with nil provider
	_, _, gCode, gErr := h.grafanaClientForConnection("token", connUUID, nil)
	if gCode != http.StatusInternalServerError || gErr == nil {
		t.Errorf("grafanaClientForConnection expected 500 and error for nil provider, got code=%d, err=%v", gCode, gErr)
	}

	_, _, prCode, prErr := h.prometheusTelemetryClientForConnection("token", connUUID, nil)
	if prCode != http.StatusInternalServerError || prErr == nil {
		t.Errorf("prometheusTelemetryClientForConnection expected 500 and error for nil provider, got code=%d, err=%v", prCode, prErr)
	}
}

// TestTelemetryHandlers_NilConnection_ReturnsNotFound verifies that all Grafana
// and Prometheus telemetry endpoints and client resolvers return 404 Not Found when
// the provider returns a nil connection without an error, preventing nil pointer
// dereferences on connection.Kind.
func TestTelemetryHandlers_NilConnection_ReturnsNotFound(t *testing.T) {
	connUUID := uuid.Must(uuid.NewV4())
	connID := connUUID.String()
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	provider := newTelemetrySpyProvider()
	provider.getConnectionFn = func(_ string, _ core.Uuid) (*connections.Connection, int, error) {
		return nil, http.StatusOK, nil
	}

	grafanaHandlers := []struct {
		name   string
		path   string
		vars   map[string]string
		invoke func(w http.ResponseWriter, req *http.Request)
	}{
		{"Ping", "/ping", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryPingHandler(w, req, nil, user, provider)
		}},
		{"Boards", "/boards", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryBoardsHandler(w, req, nil, user, provider)
		}},
		{"Board", "/boards/test-uid", map[string]string{"uid": "test-uid"}, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryBoardHandler(w, req, nil, user, provider)
		}},
		{"Datasources", "/datasources", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryDatasourcesHandler(w, req, nil, user, provider)
		}},
		{"QueryRange", "/query_range", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryQueryRangeHandler(w, req, nil, user, provider)
		}},
		{"PinnedBoards", "/pinned", nil, func(w http.ResponseWriter, req *http.Request) {
			h.GrafanaTelemetryPinnedBoardsHandler(w, req, nil, user, provider)
		}},
	}

	for _, target := range grafanaHandlers {
		t.Run("Grafana_"+target.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/telemetry/grafana/"+connID+target.path, nil)
			reqVars := map[string]string{"connectionID": connID}
			for k, v := range target.vars {
				reqVars[k] = v
			}
			req = mux.SetURLVars(req, reqVars)
			rec := httptest.NewRecorder()
			target.invoke(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Errorf("Grafana %s expected 404 for nil connection, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
			}
		})
	}

	prometheusHandlers := []struct {
		name   string
		path   string
		invoke func(w http.ResponseWriter, req *http.Request)
	}{
		{"Ping", "/ping", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryPingHandler(w, req, nil, user, provider)
		}},
		{"Metrics", "/metrics", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryMetricsHandler(w, req, nil, user, provider)
		}},
		{"Labels", "/labels", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryLabelsHandler(w, req, nil, user, provider)
		}},
		{"LabelValues", "/label_values", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryLabelValuesHandler(w, req, nil, user, provider)
		}},
		{"Metadata", "/metadata", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryMetadataHandler(w, req, nil, user, provider)
		}},
		{"Query", "/query", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryQueryHandler(w, req, nil, user, provider)
		}},
		{"QueryRange", "/query_range", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryQueryRangeHandler(w, req, nil, user, provider)
		}},
		{"Panels", "/panels", func(w http.ResponseWriter, req *http.Request) {
			h.PrometheusTelemetryPanelsHandler(w, req, nil, user, provider)
		}},
	}

	for _, target := range prometheusHandlers {
		t.Run("Prometheus_"+target.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/telemetry/prometheus/"+connID+target.path, nil)
			req = mux.SetURLVars(req, map[string]string{"connectionID": connID})
			rec := httptest.NewRecorder()
			target.invoke(rec, req)

			if rec.Code != http.StatusNotFound {
				t.Errorf("Prometheus %s expected 404 for nil connection, got %d (body=%q)", target.name, rec.Code, rec.Body.String())
			}
		})
	}

	// Client constructor helpers
	_, _, gCode, gErr := h.grafanaClientForConnection("token", connUUID, provider)
	if gCode != http.StatusNotFound || meshkiterrors.GetCode(gErr) != models.ErrResultNotFoundCode {
		t.Errorf("grafanaClientForConnection expected 404 and ErrResultNotFoundCode for nil connection, got code=%d, err=%v", gCode, gErr)
	}

	_, _, prCode, prErr := h.prometheusTelemetryClientForConnection("token", connUUID, provider)
	if prCode != http.StatusNotFound || meshkiterrors.GetCode(prErr) != models.ErrResultNotFoundCode {
		t.Errorf("prometheusTelemetryClientForConnection expected 404 and ErrResultNotFoundCode for nil connection, got code=%d, err=%v", prCode, prErr)
	}
}
