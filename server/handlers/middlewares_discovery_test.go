package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// recordingEventsProvider records the events persisted through it.
type recordingEventsProvider struct {
	models.Provider
	persisted []events.Event
}

func (p *recordingEventsProvider) PersistEvent(event events.Event, _ string) error {
	p.persisted = append(p.persisted, event)
	return nil
}

// A background Discovery reports its failure and the resolution of the
// connection's unreachable condition, and nothing for a routine success.
func TestReportBackgroundDiscovery(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:    newTestLogger(t),
	}
	user := uuid.Must(uuid.NewV4())
	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")

	routine := events.NewEvent().WithSeverity(events.Informational).WithDescription("kubernetes connection changed to connected").Build()
	failure := events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build()
	resolution := events.NewEvent().WithSeverity(events.Success).Resolves(key, time.Now()).Build()

	cases := []struct {
		name        string
		event       *events.Event
		err         error
		wantPersist bool
	}{
		{"routine success", routine, nil, false},
		{"suppressed repeat", nil, nil, false},
		{"failure", failure, errors.New("unreachable"), true},
		{"resolution", resolution, nil, true},
	}
	for _, c := range cases {
		provider := &recordingEventsProvider{}
		h.reportBackgroundDiscovery(provider, "token", user, c.event, c.err)
		if got := len(provider.persisted) == 1; got != c.wantPersist {
			t.Errorf("%s: persisted %d events, want persisted=%v", c.name, len(provider.persisted), c.wantPersist)
		}
	}
}

// tokenEventsProvider is a recordingEventsProvider that also hands out a token.
type tokenEventsProvider struct{ recordingEventsProvider }

func (p *tokenEventsProvider) GetProviderToken(*http.Request) (string, error) { return "token", nil }

func postClientEvent(t *testing.T, h *Handler, provider models.Provider, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ClientEventHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// POST /api/events records the schemas EventPayload's correlation, provenance
// and trace fields, so a client's condition correlates like the server's own.
func TestClientEventHandlerForwardsContractFields(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	actedUpon := uuid.Must(uuid.NewV4()).String()

	provider := &tokenEventsProvider{}
	rec := postClientEvent(t, h, provider, `{"actedUpon":"`+actedUpon+`","category":"connection","action":"update","severity":"error",
		"description":"cluster unreachable","correlationKey":"meshery-ui/k8s-connection/abc/unreachable","occurrenceCount":2,
		"source":"meshery-ui","traceId":"4bf92f3577b34da6a3ce929d0e0e4736","parentId":"00f067aa0ba902b7"}`)
	if rec.Code != http.StatusCreated || len(provider.persisted) != 1 {
		t.Fatalf("status %d, %d persisted; body %s", rec.Code, len(provider.persisted), rec.Body.String())
	}
	got := provider.persisted[0]
	if got.CorrelationKey == nil || *got.CorrelationKey != "meshery-ui/k8s-connection/abc/unreachable" ||
		got.OccurrenceCount == nil || *got.OccurrenceCount != 2 ||
		got.Source == nil || *got.Source != "meshery-ui" ||
		got.TraceID == nil || *got.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" || got.ParentID == nil {
		t.Fatalf("contract fields were not forwarded: %+v", got)
	}

	// A client that names no source is not attributed to the server.
	provider = &tokenEventsProvider{}
	rec = postClientEvent(t, h, provider, `{"actedUpon":"`+actedUpon+`","category":"design","action":"save","severity":"informational"}`)
	if rec.Code != http.StatusCreated || provider.persisted[0].Source == nil || *provider.persisted[0].Source != models.EventSourceUnknownClient {
		t.Fatalf("status %d, source %v; want %q", rec.Code, provider.persisted[0].Source, models.EventSourceUnknownClient)
	}

	// A field that breaks the contract is refused before anything is recorded.
	provider = &tokenEventsProvider{}
	rec = postClientEvent(t, h, provider, `{"actedUpon":"`+actedUpon+`","category":"connection","action":"update","severity":"error","traceId":"not-hex"}`)
	if rec.Code != http.StatusBadRequest || len(provider.persisted) != 0 {
		t.Fatalf("invalid traceId: status %d, %d persisted; want 400 and none", rec.Code, len(provider.persisted))
	}
}
