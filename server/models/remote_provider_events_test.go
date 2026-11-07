package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
)

// newEventsRemoteProvider returns a RemoteProvider whose PersistEvents
// capability is served by handler, with a local EventsPersister to fall back
// to.
func newEventsRemoteProvider(t *testing.T, handler http.HandlerFunc) *RemoteProvider {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	provider := newTestRemoteProvider(t, server.URL)
	provider.ProviderProperties = ProviderProperties{
		ProviderType: RemoteProviderType,
		ProviderURL:  server.URL,
		ProviderName: "test-remote",
		Capabilities: Capabilities{{Feature: PersistEvents, Endpoint: "/api/events"}},
	}
	provider.EventsPersister = newCorrelationTestPersister(t)
	return provider
}

func localEventCount(t *testing.T, provider *RemoteProvider) int64 {
	t.Helper()
	var n int64
	if err := provider.EventsPersister.DB.Model(&events.Event{}).Count(&n).Error; err != nil {
		t.Fatalf("count local events: %v", err)
	}
	return n
}

// The remote provider receives every correlation, provenance and trace field
// the event carries, under the schemas JSON names, plus the server's source.
func TestRemoteProviderPersistEventSendsCorrelationFields(t *testing.T) {
	var body map[string]interface{}
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("remote provider received invalid JSON: %v", err)
		}
		w.Header().Set(EventDispositionHeader, "created")
		w.WriteHeader(http.StatusCreated)
	})

	firstSeen := time.Date(2026, 11, 6, 10, 0, 0, 0, time.UTC)
	event := events.NewEvent().
		FromOwner(mustUUID(t)).
		WithCategory("connection").
		WithAction("update").
		WithSeverity(events.Error).
		WithCorrelationKey("meshery-server/k8s-connection/abc/unreachable").
		WithOccurrenceCount(2).
		WithFirstSeenAt(firstSeen).
		WithLastSeenAt(firstSeen.Add(time.Minute)).
		WithType("io.meshery.connection.unreachable").
		WithSpecVersion(events.SpecVersion10).
		WithTraceContext("4bf92f3577b34da6a3ce929d0e0e4736", "00f067aa0ba902b7").
		Build()

	if _, err := provider.PersistEvent(*event, "token"); err != nil {
		t.Fatalf("PersistEvent: %v", err)
	}

	want := map[string]interface{}{
		"correlationKey":  "meshery-server/k8s-connection/abc/unreachable",
		"occurrenceCount": float64(2),
		"firstSeenAt":     "2026-11-06T10:00:00Z",
		"lastSeenAt":      "2026-11-06T10:01:00Z",
		"source":          EventSourceServer,
		"type":            "io.meshery.connection.unreachable",
		"specversion":     "1.0",
		"traceId":         "4bf92f3577b34da6a3ce929d0e0e4736",
		"parentId":        "00f067aa0ba902b7",
	}
	for key, value := range want {
		if body[key] != value {
			t.Errorf("remote payload %s = %v, want %v", key, body[key], value)
		}
	}
	if _, ok := body["resolvedAt"]; ok {
		t.Errorf("remote payload carries resolvedAt for an open condition: %v", body["resolvedAt"])
	}
	// 201 means the remote provider recorded it: no local copy.
	if n := localEventCount(t, provider); n != 0 {
		t.Fatalf("a 201 fell back to the local database (%d rows), duplicating the event", n)
	}
}

// An event that sets none of the new fields sends none of them, so a remote
// provider that predates them receives exactly what it did before, apart from
// the server's source. A client's event is sent with its source as sent.
func TestRemoteProviderPersistEventOmitsUnsetFields(t *testing.T) {
	var body map[string]interface{}
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.WriteHeader(http.StatusOK)
	})

	event := events.NewEvent().FromOwner(mustUUID(t)).WithCategory("design").WithAction("save").
		WithSeverity(events.Informational).WithDescription("saved").Build()
	if _, err := provider.PersistEvent(*event, "token"); err != nil {
		t.Fatalf("PersistEvent: %v", err)
	}

	for _, key := range []string{"correlationKey", "occurrenceCount", "firstSeenAt", "lastSeenAt", "resolvedAt", "type", "specversion", "traceId", "parentId"} {
		if _, ok := body[key]; ok {
			t.Errorf("remote payload carries unset %s: %v", key, body[key])
		}
	}
	if body["source"] != EventSourceServer {
		t.Errorf("source = %v, want %q", body["source"], EventSourceServer)
	}

	body = nil
	if _, err := provider.PersistClientEvent(*event, "token"); err != nil {
		t.Fatalf("PersistClientEvent: %v", err)
	}
	if source, ok := body["source"]; ok {
		t.Errorf("client event sent with source %v, want none", source)
	}
	if n := localEventCount(t, provider); n != 0 {
		t.Fatalf("a 200 fell back to the local database (%d rows)", n)
	}
}

// When the remote provider fails, the local fallback applies the same
// correlation contract: repeated occurrences count on one open event.
func TestRemoteProviderPersistEventFallbackCorrelates(t *testing.T) {
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	owner := mustUUID(t)
	for i := 0; i < 3; i++ {
		if _, err := provider.PersistEvent(*conditionEvent(owner, "meshery-server/k8s-connection/xyz/unreachable", "down").Build(), "token"); err != nil {
			t.Fatalf("PersistEvent: %v", err)
		}
	}

	rows := storedEvents(t, provider.EventsPersister, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != 3 {
		t.Fatalf("got %+v, want one local event counting 3 occurrences", rows)
	}
}

// Every 2xx answer means the remote provider has the event; anything else
// falls back to the local database.
func TestRemoteProviderPersistEventAcceptsAny2xx(t *testing.T) {
	cases := []struct {
		status    int
		wantLocal int64
	}{
		{http.StatusOK, 0},
		{http.StatusCreated, 0},
		{http.StatusAccepted, 0},
		{http.StatusNoContent, 0},
		{http.StatusMultipleChoices, 1},
		{http.StatusBadRequest, 1},
		{http.StatusInternalServerError, 1},
	}
	for _, c := range cases {
		status := c.status
		provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set(EventDispositionHeader, "incremented")
			w.WriteHeader(status)
		})
		event := events.NewEvent().FromOwner(mustUUID(t)).WithCategory("connection").WithAction("update").
			WithSeverity(events.Error).Build()
		persisted, err := provider.PersistEvent(*event, "token")
		if err != nil {
			t.Fatalf("status %d: PersistEvent: %v", status, err)
		}
		if n := localEventCount(t, provider); n != c.wantLocal {
			t.Errorf("status %d: %d local rows, want %d", status, n, c.wantLocal)
		}
		// A 2xx that reports an occurrence counted but carries no stored event
		// must not present the sent event, whose ID no row has, as stored.
		if c.wantLocal == 0 && persisted.Event != nil {
			t.Errorf("status %d: returned event %s for an incremented create with no body, want none", status, persisted.Event.ID)
		}
	}
}

// A body-less or unparseable answer for a create that recorded nothing new
// returns no event, whatever the disposition; one for a created event returns
// the event as sent.
func TestRemoteProviderPersistEventWithoutTheStoredEvent(t *testing.T) {
	for _, c := range []struct {
		disposition EventDisposition
		body        string
		wantSent    bool
	}{
		{EventIncremented, "", false},
		{EventResolved, "not json", false},
		{EventReplayed, "{}", false},
		{EventCreated, "", true},
		{"", "", true},
	} {
		c := c
		provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, _ *http.Request) {
			if c.disposition != "" {
				w.Header().Set(EventDispositionHeader, string(c.disposition))
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(c.body))
		})
		sent := events.NewEvent().FromOwner(mustUUID(t)).WithCategory("connection").WithAction("update").
			WithSeverity(events.Error).Build()
		persisted, err := provider.PersistEvent(*sent, "token")
		if err != nil {
			t.Fatalf("%q: PersistEvent: %v", c.disposition, err)
		}
		if c.wantSent {
			if persisted.Event == nil || persisted.Event.ID != sent.ID {
				t.Errorf("%q: got %v, want the event as sent", c.disposition, persisted.Event)
			}
		} else if persisted.Event != nil {
			t.Errorf("%q: got event %s, want none", c.disposition, persisted.Event.ID)
		}
	}
}

// PersistEvent returns the event the remote provider answered with, which for
// an incremented or resolved condition is the condition's stored event, and
// the disposition it reported. A legacy answer with no body means the event
// was stored as sent.
func TestRemoteProviderPersistEventReturnsTheStoredEvent(t *testing.T) {
	owner := mustUUID(t)
	stored := conditionEvent(owner, "meshery-server/k8s-connection/s/unreachable", "down").WithOccurrenceCount(4).Build()
	stored.ID = mustUUID(t)
	answer, err := json.Marshal(stored)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(EventDispositionHeader, string(EventIncremented))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(answer)
	})
	sent := conditionEvent(owner, "meshery-server/k8s-connection/s/unreachable", "down").Build()
	persisted, err := provider.PersistEvent(*sent, "token")
	if err != nil {
		t.Fatalf("PersistEvent: %v", err)
	}
	if persisted.Disposition != EventIncremented || persisted.Event.ID != stored.ID || count(persisted.Event.OccurrenceCount) != 4 {
		t.Fatalf("got %s %s count %d, want incremented %s count 4 (sent %s)",
			persisted.Disposition, persisted.Event.ID, count(persisted.Event.OccurrenceCount), stored.ID, sent.ID)
	}

	for _, c := range []struct {
		status int
		header string
		want   EventDisposition
	}{
		{http.StatusOK, "created", EventCreated},
		{http.StatusCreated, "", EventCreated},
		{http.StatusOK, "", ""},
	} {
		c := c
		provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, _ *http.Request) {
			if c.header != "" {
				w.Header().Set(EventDispositionHeader, c.header)
			}
			w.WriteHeader(c.status)
		})
		sent := events.NewEvent().FromOwner(owner).WithCategory("design").WithAction("save").WithSeverity(events.Informational).Build()
		sent.ID = mustUUID(t)
		persisted, err := provider.PersistEvent(*sent, "token")
		if err != nil {
			t.Fatalf("status %d: PersistEvent: %v", c.status, err)
		}
		if persisted.Event.ID != sent.ID || persisted.Disposition != c.want {
			t.Errorf("status %d %q: got %s %q, want the sent %s and %q", c.status, c.header, persisted.Event.ID, persisted.Disposition, sent.ID, c.want)
		}
	}
}

// An event constructed without the builder carries no ID. It is given one
// before it is sent, so a legacy remote provider that answers with no body,
// and therefore stores the event as sent, returns it with that ID rather than
// the zero ID.
func TestRemoteProviderPersistEventAssignsAnIDToAnEventWithout(t *testing.T) {
	var sentID string
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		sentID, _ = body["id"].(string)
		w.WriteHeader(http.StatusOK)
	})
	sent := events.NewEvent().FromOwner(mustUUID(t)).WithCategory("design").WithAction("save").WithSeverity(events.Informational).Build()
	// events.NewEvent mints an ID; an event constructed without the
	// builder carries none, which is the case under test.
	sent.ID = uuid.Nil
	persisted, err := provider.PersistEvent(*sent, "token")
	if err != nil {
		t.Fatalf("PersistEvent: %v", err)
	}
	if persisted.Event.ID.IsNil() || persisted.Event.ID.String() != sentID {
		t.Fatalf("returned ID %s, sent ID %q; want the same non-zero ID", persisted.Event.ID, sentID)
	}
}

// A remote provider is known to correlate conditions only once a create
// carrying a correlationKey is answered with an Event-Disposition, and stops
// being known to once such a create is answered without one. A create without
// a correlationKey says nothing either way.
func TestRemoteProviderCorrelatesConditionsFollowsCorrelatedCreates(t *testing.T) {
	var disposition string
	provider := newEventsRemoteProvider(t, func(w http.ResponseWriter, r *http.Request) {
		if disposition != "" {
			w.Header().Set(EventDispositionHeader, disposition)
		}
		w.WriteHeader(http.StatusCreated)
	})
	owner := mustUUID(t)
	correlated := func() events.Event {
		return *events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
			WithSeverity(events.Error).WithCorrelationKey("meshery-server/k8s-connection/abc/unreachable").Build()
	}
	plain := func() events.Event {
		return *events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
			WithSeverity(events.Informational).Build()
	}
	persist := func(event events.Event) {
		t.Helper()
		if _, err := provider.PersistEvent(event, "token"); err != nil {
			t.Fatalf("PersistEvent: %v", err)
		}
	}

	if provider.CorrelatesConditions() {
		t.Fatal("a remote provider no correlated create has reached is known to correlate")
	}
	disposition = string(EventCreated)
	persist(plain())
	if provider.CorrelatesConditions() {
		t.Fatal("an uncorrelated create made the remote provider known to correlate")
	}
	persist(correlated())
	if !provider.CorrelatesConditions() {
		t.Fatal("a correlated create answered with a disposition did not make the remote provider known to correlate")
	}
	disposition = ""
	persist(plain())
	if !provider.CorrelatesConditions() {
		t.Fatal("an uncorrelated create answered without a disposition changed what is known")
	}
	persist(correlated())
	if provider.CorrelatesConditions() {
		t.Fatal("a correlated create answered without a disposition left the remote provider known to correlate")
	}
}
