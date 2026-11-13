package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// recordingEventsProvider records the events persisted through it. It answers
// with stored when set, as a persister does for an occurrence counted on, or
// the resolution of, an open condition; otherwise with the event as created.
// It correlates conditions unless uncorrelated is set.
type recordingEventsProvider struct {
	models.Provider
	persisted    []events.Event
	stored       *models.PersistedEvent
	err          error
	uncorrelated bool
}

func (p *recordingEventsProvider) CorrelatesConditions() bool { return !p.uncorrelated }

func (p *recordingEventsProvider) PersistEvent(event events.Event, _ string) (models.PersistedEvent, error) {
	p.persisted = append(p.persisted, event)
	if p.err != nil {
		return models.PersistedEvent{}, p.err
	}
	if p.stored != nil {
		return *p.stored, nil
	}
	return models.PersistedEvent{Event: &event, Disposition: models.EventCreated}, nil
}

func (p *recordingEventsProvider) PersistClientEvent(event events.Event, token string) (models.PersistedEvent, error) {
	return p.PersistEvent(event, token)
}

// persisterEventsProvider persists through a real EventsPersister.
type persisterEventsProvider struct {
	models.Provider
	persister *models.EventsPersister
}

func (p *persisterEventsProvider) PersistEvent(event events.Event, token string) (models.PersistedEvent, error) {
	return p.persister.PersistEvent(event, token)
}

func (p *persisterEventsProvider) PersistClientEvent(event events.Event, token string) (models.PersistedEvent, error) {
	return p.persister.PersistClientEvent(event, token)
}

func (p *persisterEventsProvider) CorrelatesConditions() bool {
	return p.persister.CorrelatesConditions()
}

func (p *persisterEventsProvider) GetProviderToken(*http.Request) (string, error) { return "", nil }

// An event the server emits is recorded with source meshery-server even when
// built with a bare events.NewEvent(); an event a client sends without a
// source is recorded without one.
func TestEventSourceIsTheServersOnlyForServerEvents(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&events.Event{}); err != nil {
		t.Fatalf("migrate events: %v", err)
	}
	provider := &persisterEventsProvider{persister: &models.EventsPersister{DB: &database.Handler{DB: db}}}
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	user := uuid.Must(uuid.NewV4())

	h.persistAndPublish(provider, "", user, events.NewEvent().FromOwner(user).WithCategory("connection").
		WithAction("update").WithSeverity(events.Informational).WithDescription("server event").Build())

	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(`{"actedUpon":"`+uuid.Must(uuid.NewV4()).String()+
		`","category":"design","action":"save","severity":"informational","description":"client event"}`))
	rec := httptest.NewRecorder()
	h.ClientEventHandler(rec, req, nil, &models.User{ID: user}, provider)
	if rec.Code != http.StatusCreated {
		t.Fatalf("client event: status %d; body %s", rec.Code, rec.Body.String())
	}

	var rows []events.Event
	if err := db.Find(&rows).Error; err != nil {
		t.Fatalf("read events: %v", err)
	}
	sources := map[string]*string{}
	for _, row := range rows {
		sources[row.Description] = row.Source
	}
	if len(rows) != 2 {
		t.Fatalf("got %d events, want 2", len(rows))
	}
	if s := sources["server event"]; s == nil || *s != models.EventSourceServer {
		t.Errorf("server event source = %v, want %q", s, models.EventSourceServer)
	}
	if s := sources["client event"]; s != nil {
		t.Errorf("client event source = %q, want none", *s)
	}
}

// An event constructed without the builder carries no ID; it
// is broadcast as stored, with the stored row's non-zero ID and the server's
// source, never as the zero-ID input.
func TestPersistAndPublishBroadcastsTheStoredRow(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&events.Event{}); err != nil {
		t.Fatalf("migrate events: %v", err)
	}
	provider := &persisterEventsProvider{persister: &models.EventsPersister{DB: &database.Handler{DB: db}}}
	broadcaster := models.NewBroadcaster("test")
	user := uuid.Must(uuid.NewV4())
	ch, unsubscribe := broadcaster.Subscribe(user)
	defer unsubscribe()

	publishers := map[string]func(*events.Event){
		"handler": func(e *events.Event) {
			h := &Handler{config: &models.HandlerConfig{EventBroadcaster: broadcaster}, log: newTestLogger(t)}
			h.persistAndPublish(provider, "", user, e)
		},
		"adapter events": func(e *events.Event) {
			persistAndPublishEvent(newTestLogger(t), provider, "", broadcaster, user, e)
		},
	}
	for name, publish := range publishers {
		sent := events.NewEvent().FromOwner(user).WithCategory("kubernetes").WithAction("register").
			WithSeverity(events.Informational).WithDescription(name).Build()
		// events.NewEvent mints an ID; an event constructed without the
		// builder carries none, which is the case under test.
		sent.ID = uuid.Nil
		publish(sent)

		got := receiveEvent(t, ch)
		if got.ID.IsNil() {
			t.Fatalf("%s: broadcast the zero ID", name)
		}
		var stored events.Event
		if err := db.Where("description = ?", name).First(&stored).Error; err != nil {
			t.Fatalf("%s: read stored event: %v", name, err)
		}
		if got.ID != stored.ID {
			t.Errorf("%s: broadcast ID %s, stored ID %s", name, got.ID, stored.ID)
		}
		if got.Source == nil || *got.Source != models.EventSourceServer {
			t.Errorf("%s: broadcast source %v, want %q", name, got.Source, models.EventSourceServer)
		}
	}
}

// receiveEvent returns the next event published on ch.
func receiveEvent(t *testing.T, ch chan interface{}) *events.Event {
	t.Helper()
	select {
	case data := <-ch:
		event, ok := data.(*events.Event)
		if !ok {
			t.Fatalf("published %T, want *events.Event", data)
		}
		return event
	case <-time.After(5 * time.Second):
		t.Fatal("no event was published")
		return nil
	}
}

// An event that resolves (or counts an occurrence on) an open condition is
// broadcast as the condition's stored event, never as the event that was
// sent, whose ID no stored row has.
func TestReportBackgroundDiscoveryPublishesStoredEvent(t *testing.T) {
	broadcaster := models.NewBroadcaster("test")
	h := &Handler{
		config: &models.HandlerConfig{EventBroadcaster: broadcaster},
		log:    newTestLogger(t),
	}
	user := uuid.Must(uuid.NewV4())
	ch, unsubscribe := broadcaster.Subscribe(user)
	defer unsubscribe()

	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")
	open := events.NewEvent().WithSeverity(events.Success).WithCorrelationKey(key).Build()
	open.ID = uuid.Must(uuid.NewV4())
	for _, disposition := range []models.EventDisposition{models.EventResolved, models.EventIncremented} {
		sent := events.NewEvent().WithSeverity(events.Success).Resolves(key, time.Now()).Build()
		provider := &recordingEventsProvider{stored: &models.PersistedEvent{Event: open, Disposition: disposition}}
		h.reportBackgroundDiscovery(provider, "token", user, sent, nil)
		if got := receiveEvent(t, ch); got.ID != open.ID {
			t.Fatalf("%s: broadcast event %s, want the stored condition %s (sent %s)", disposition, got.ID, open.ID, sent.ID)
		}
	}
}

func TestReportBackgroundDiscoveryResolvesSharedConnection(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	broadcaster := h.config.EventBroadcaster
	observer := uuid.Must(uuid.NewV4())
	observerEvents, stopObserver := broadcaster.Subscribe(observer)
	defer stopObserver()
	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")
	otherKey := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")

	open := func(owner uuid.UUID, key string) *events.Event {
		t.Helper()
		event := events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
			WithSeverity(events.Error).WithDescription("connection unreachable").WithCorrelationKey(key).Build()
		persisted, err := provider.PersistEvent(*event, "")
		if err != nil {
			t.Fatalf("open condition: %v", err)
		}
		return persisted.Event
	}

	firstOwner, secondOwner := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	open(firstOwner, key)
	first := open(firstOwner, key) // Keep the recorded recurrence count.
	second := open(secondOwner, key)
	unrelated := open(firstOwner, otherKey)

	firstEvents, stopFirst := broadcaster.Subscribe(firstOwner)
	defer stopFirst()
	secondEvents, stopSecond := broadcaster.Subscribe(secondOwner)
	defer stopSecond()
	recovery := events.NewEvent().FromOwner(observer).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("connection reachable").Resolves(key, time.Now().UTC()).Build()
	h.reportBackgroundDiscovery(provider, "", observer, recovery, nil)

	for _, owner := range []struct {
		before *events.Event
		stream chan interface{}
	}{
		{first, firstEvents},
		{second, secondEvents},
	} {
		got := receiveEvent(t, owner.stream)
		if got.ID != owner.before.ID || got.Owner == nil || *got.Owner != *owner.before.Owner ||
			got.ResolvedAt == nil || !got.ResolvedAt.Equal(*recovery.ResolvedAt) ||
			got.OccurrenceCount == nil || *got.OccurrenceCount != *owner.before.OccurrenceCount ||
			got.Severity != recovery.Severity || got.Description != recovery.Description {
			t.Errorf("published event %+v does not match the resolved condition %+v", got, owner.before)
		}
		var stored events.Event
		if err := provider.EventsPersister.DB.Where("id = ?", owner.before.ID).Take(&stored).Error; err != nil {
			t.Fatalf("read condition: %v", err)
		}
		if stored.ResolvedAt == nil {
			t.Errorf("condition %s remains open", stored.ID)
		}
	}

	h.reportBackgroundDiscovery(provider, "", observer, recovery, nil)
	for _, stream := range []chan interface{}{firstEvents, secondEvents, observerEvents} {
		select {
		case unexpected := <-stream:
			t.Errorf("unexpected notification: %+v", unexpected)
		default:
		}
	}
	var rows []events.Event
	if err := provider.EventsPersister.DB.Find(&rows).Error; err != nil {
		t.Fatalf("read events: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d event rows, want 3 original rows", len(rows))
	}
	for _, row := range rows {
		if row.ID == unrelated.ID && row.ResolvedAt != nil {
			t.Errorf("unrelated condition %s was resolved", row.ID)
		}
	}
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
		{"suppressed", nil, nil, false},
		{"recurrence", events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build(), nil, true},
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
		"source":"meshery-ui","traceId":"4bf92f3577b34da6a3ce929d0e0e4736","parentId":"00f067aa0ba902b7","idempotencyKey":"ui/retry-1"}`)
	if rec.Code != http.StatusCreated || len(provider.persisted) != 1 {
		t.Fatalf("status %d, %d persisted; body %s", rec.Code, len(provider.persisted), rec.Body.String())
	}
	got := provider.persisted[0]
	if got.CorrelationKey == nil || *got.CorrelationKey != "meshery-ui/k8s-connection/abc/unreachable" ||
		got.OccurrenceCount == nil || *got.OccurrenceCount != 2 ||
		got.Source == nil || *got.Source != "meshery-ui" ||
		got.TraceID == nil || *got.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" || got.ParentID == nil ||
		got.IdempotencyKey == nil || *got.IdempotencyKey != "ui/retry-1" {
		t.Fatalf("contract fields were not forwarded: %+v", got)
	}

	// A client that names no source is recorded without one.
	provider = &tokenEventsProvider{}
	rec = postClientEvent(t, h, provider, `{"actedUpon":"`+actedUpon+`","category":"design","action":"save","severity":"informational"}`)
	if rec.Code != http.StatusCreated || provider.persisted[0].Source != nil {
		t.Fatalf("status %d, source %v; want none", rec.Code, provider.persisted[0].Source)
	}

	// A field that breaks the contract is refused before anything is recorded.
	provider = &tokenEventsProvider{}
	rec = postClientEvent(t, h, provider, `{"actedUpon":"`+actedUpon+`","category":"connection","action":"update","severity":"error","traceId":"not-hex"}`)
	if rec.Code != http.StatusBadRequest || len(provider.persisted) != 0 {
		t.Fatalf("invalid traceId: status %d, %d persisted; want 400 and none", rec.Code, len(provider.persisted))
	}
}

// POST /api/events answers, and broadcasts, the event as stored: an occurrence
// a client reports on an open condition comes back as that condition's event,
// with the disposition the persister reported.
func TestClientEventHandlerAnswersStoredEvent(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	broadcaster := models.NewBroadcaster("test")
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: broadcaster},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ch, unsubscribe := broadcaster.Subscribe(user.ID)
	defer unsubscribe()

	open := events.NewEvent().WithSeverity(events.Error).WithCorrelationKey("meshery-ui/k8s-connection/abc/unreachable").WithOccurrenceCount(3).Build()
	open.ID = uuid.Must(uuid.NewV4())
	provider := &tokenEventsProvider{recordingEventsProvider{stored: &models.PersistedEvent{Event: open, Disposition: models.EventIncremented}}}
	req := httptest.NewRequest(http.MethodPost, "/api/events", strings.NewReader(`{"actedUpon":"`+uuid.Must(uuid.NewV4()).String()+
		`","category":"connection","action":"update","severity":"error","correlationKey":"meshery-ui/k8s-connection/abc/unreachable","source":"meshery-ui"}`))
	rec := httptest.NewRecorder()
	h.ClientEventHandler(rec, req, nil, user, provider)

	// An occurrence counted on an open condition recorded nothing new.
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, want 200; body %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get(models.EventDispositionHeader); got != string(models.EventIncremented) {
		t.Errorf("%s = %q, want %q", models.EventDispositionHeader, got, models.EventIncremented)
	}
	var answered events.Event
	if err := json.Unmarshal(rec.Body.Bytes(), &answered); err != nil {
		t.Fatalf("response is not an event: %v", err)
	}
	if answered.ID != open.ID || answered.OccurrenceCount == nil || *answered.OccurrenceCount != 3 {
		t.Errorf("answered event %s (count %v), want the stored condition %s", answered.ID, answered.OccurrenceCount, open.ID)
	}
	if got := receiveEvent(t, ch); got.ID != open.ID {
		t.Errorf("broadcast event %s, want the stored condition %s", got.ID, open.ID)
	}
}

// POST /api/events answers 201 only when the create recorded a new event, and
// 200 when it recorded nothing new, as the createEvent contract requires. An
// answer that names no stored event is sent with no body and broadcasts
// nothing.
func TestClientEventHandlerStatusFollowsTheDisposition(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	body := `{"actedUpon":"` + uuid.Must(uuid.NewV4()).String() + `","category":"connection","action":"update","severity":"error"}`
	stored := events.NewEvent().WithSeverity(events.Error).Build()
	for _, c := range []struct {
		disposition models.EventDisposition
		event       *events.Event
		want        int
	}{
		{models.EventCreated, stored, http.StatusCreated},
		{"", stored, http.StatusCreated},
		{models.EventIncremented, stored, http.StatusOK},
		{models.EventResolved, stored, http.StatusOK},
		{models.EventReplayed, stored, http.StatusOK},
		{models.EventIncremented, nil, http.StatusOK},
	} {
		provider := &tokenEventsProvider{recordingEventsProvider{stored: &models.PersistedEvent{Event: c.event, Disposition: c.disposition}}}
		rec := postClientEvent(t, h, provider, body)
		if rec.Code != c.want {
			t.Errorf("%q: status %d, want %d", c.disposition, rec.Code, c.want)
		}
		if c.event == nil && rec.Body.Len() != 0 {
			t.Errorf("%q with no stored event: body %q, want none", c.disposition, rec.Body.String())
		}
	}
}

// persistAndPublishEvent reports nothing for a nil event (a suppressed
// Discovery failure), broadcasts nothing when the stored event is unknown, and
// suppresses broadcast when persistence fails to avoid ghost notifications.
func TestPersistAndPublishEventEdgeCases(t *testing.T) {
	broadcaster := models.NewBroadcaster("test")
	user := uuid.Must(uuid.NewV4())
	ch, unsubscribe := broadcaster.Subscribe(user)
	defer unsubscribe()
	log := newTestLogger(t)

	provider := &recordingEventsProvider{}
	persistAndPublishEvent(log, provider, "", broadcaster, user, nil)
	if len(provider.persisted) != 0 {
		t.Fatalf("persisted %d events for a nil event", len(provider.persisted))
	}

	provider = &recordingEventsProvider{stored: &models.PersistedEvent{Disposition: models.EventIncremented}}
	persistAndPublishEvent(log, provider, "", broadcaster, user, events.NewEvent().WithSeverity(events.Error).Build())
	select {
	case got := <-ch:
		t.Fatalf("broadcast %v for an occurrence whose stored event is unknown", got)
	case <-time.After(100 * time.Millisecond):
	}

	provider = &recordingEventsProvider{err: errors.New("database is down")}
	sent := events.NewEvent().WithSeverity(events.Error).Build()
	sent.ID = uuid.Nil
	persistAndPublishEvent(log, provider, "", broadcaster, user, sent)
	select {
	case got := <-ch:
		t.Fatalf("unexpected broadcast %v when persistence failed", got)
	case <-time.After(100 * time.Millisecond):
	}
}

// A Recurrence is counted on its open condition without notifying anyone.
func TestReportBackgroundDiscoveryCountsARecurrenceSilently(t *testing.T) {
	broadcaster := models.NewBroadcaster("test")
	h := &Handler{config: &models.HandlerConfig{EventBroadcaster: broadcaster}, log: newTestLogger(t)}
	user := uuid.Must(uuid.NewV4())
	ch, unsubscribe := broadcaster.Subscribe(user)
	defer unsubscribe()

	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")
	provider := &recordingEventsProvider{}
	h.reportBackgroundDiscovery(provider, "token", user, events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build(), nil)
	if len(provider.persisted) != 1 {
		t.Fatalf("persisted %d events, want the recurrence counted", len(provider.persisted))
	}
	select {
	case got := <-ch:
		t.Fatalf("broadcast %v for a recurrence", got)
	case <-time.After(100 * time.Millisecond):
	}
}

// A provider not known to correlate conditions is sent no Recurrence, which it
// would record as a new event, but is still sent a condition's first
// occurrence and its resolution.
func TestReportMachineOutcomeSendsRecurrencesOnlyToACorrelatingProvider(t *testing.T) {
	key := models.ConditionKey("k8s-connection", uuid.Must(uuid.NewV4()).String(), "unreachable")
	user := uuid.Must(uuid.NewV4())
	resolvedAt := time.Now()
	resolution := events.NewEvent().WithSeverity(events.Informational).WithCorrelationKey(key).WithResolvedAt(resolvedAt).Build()
	firstOccurrence := events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build()
	recurrence := events.NewEvent().WithSeverity(events.Error).WithCorrelationKey(key).Build()

	for _, tc := range []struct {
		name         string
		uncorrelated bool
		want         []*events.Event
	}{
		{name: "correlating", want: []*events.Event{firstOccurrence, recurrence, resolution}},
		{name: "not correlating", uncorrelated: true, want: []*events.Event{firstOccurrence, resolution}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{config: &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")}, log: newTestLogger(t)}
			provider := &recordingEventsProvider{uncorrelated: tc.uncorrelated}
			h.reportMachineOutcome(provider, "token", user, firstOccurrence, errors.New("cluster unreachable"))
			h.reportMachineOutcome(provider, "token", user, recurrence, nil)
			h.reportMachineOutcome(provider, "token", user, resolution, nil)

			if len(provider.persisted) != len(tc.want) {
				t.Fatalf("persisted %d events, want %d", len(provider.persisted), len(tc.want))
			}
			for i, want := range tc.want {
				if provider.persisted[i].ID != want.ID {
					t.Errorf("persisted[%d] = %s, want %s", i, provider.persisted[i].ID, want.ID)
				}
			}
		})
	}
}
