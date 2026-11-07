package models

import (
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// newCorrelationTestPersister returns an EventsPersister over a database built
// by AutoMigrateSystemTables, so the open-condition index under test is the
// one Meshery Server boots with.
func newCorrelationTestPersister(t *testing.T) *EventsPersister {
	t.Helper()
	db := newMigratedDB(t)
	// Every pooled connection to ":memory:" opens its own empty database, so
	// the test pins the pool to one connection.
	sqlDB, err := db.DB.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	return &EventsPersister{DB: db}
}

func conditionEvent(owner core.Uuid, key, description string) *events.EventBuilder {
	return events.NewEvent().
		FromOwner(owner).
		WithCategory("connection").
		WithAction("update").
		WithSeverity(events.Error).
		WithDescription(description).
		WithCorrelationKey(key)
}

func storedEvents(t *testing.T, p *EventsPersister, owner core.Uuid) []events.Event {
	t.Helper()
	var rows []events.Event
	if err := p.DB.Where("owner = ?", owner).Order("created_at asc").Find(&rows).Error; err != nil {
		t.Fatalf("read events: %v", err)
	}
	return rows
}

// count and timeOf read optional fields without dereferencing nil, so a failing
// assertion reports instead of panicking and hiding the tests after it.
func count(n *int) int {
	if n == nil {
		return 0
	}
	return *n
}

func timeOf(at *time.Time) time.Time {
	if at == nil {
		return time.Time{}
	}
	return *at
}

func mustPersist(t *testing.T, p *EventsPersister, event *events.Event) PersistedEvent {
	t.Helper()
	persisted, err := p.PersistEvent(*event, "")
	if err != nil {
		t.Fatalf("PersistEvent: %v", err)
	}
	return persisted
}

// An event that sets none of the new fields is recorded as a new row every
// time, exactly as before correlation existed; the only field it gains is the
// server's source.
func TestPersistEventWithoutCorrelationKeyRecordsEveryEvent(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	for i := 0; i < 2; i++ {
		mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
			WithSeverity(events.Error).WithDescription("same failure").Build())
	}

	rows := storedEvents(t, p, owner)
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2: uncorrelated events must never be rolled up", len(rows))
	}
	for _, row := range rows {
		if row.CorrelationKey != nil || row.OccurrenceCount != nil || row.FirstSeenAt != nil ||
			row.LastSeenAt != nil || row.ResolvedAt != nil || row.TraceID != nil || row.ParentID != nil {
			t.Errorf("uncorrelated event gained correlation fields: %+v", row)
		}
		if row.Source == nil || *row.Source != EventSourceServer {
			t.Errorf("source = %v, want %q", row.Source, EventSourceServer)
		}
	}
}

// A producer's own source is kept.
func TestPersistEventKeepsProducerSource(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("design").WithAction("save").
		WithSeverity(events.Informational).WithSource("meshery-adapter-istio").Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || rows[0].Source == nil || *rows[0].Source != "meshery-adapter-istio" {
		t.Fatalf("got %+v, want one row with the producer's source", rows)
	}
}

// A client's event is recorded with its source as sent, none when it sent none.
func TestPersistClientEventKeepsTheClientsSource(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	for _, source := range []string{"", "meshery-ui"} {
		builder := events.NewEvent().FromOwner(owner).WithCategory("design").WithAction("save").
			WithSeverity(events.Informational).WithDescription("client " + source)
		if source != "" {
			builder = builder.WithSource(source)
		}
		if _, err := p.PersistClientEvent(*builder.Build(), ""); err != nil {
			t.Fatalf("PersistClientEvent: %v", err)
		}
	}

	for _, row := range storedEvents(t, p, owner) {
		switch row.Description {
		case "client ":
			if row.Source != nil {
				t.Errorf("source = %q, want none for a client that sent none", *row.Source)
			}
		case "client meshery-ui":
			if row.Source == nil || *row.Source != "meshery-ui" {
				t.Errorf("source = %v, want the client's", row.Source)
			}
		}
	}
}

func TestPersistEventCountsOccurrencesOfOpenCondition(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := ConditionKey("k8s-connection", mustUUID(t).String(), "unreachable")

	before := time.Now().UTC().Add(-time.Second)
	mustPersist(t, p, conditionEvent(owner, key, "first occurrence").Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 {
		t.Fatalf("got %d rows after opening the condition, want 1", len(rows))
	}
	opened := rows[0]
	if opened.OccurrenceCount == nil || *opened.OccurrenceCount != 1 {
		t.Fatalf("occurrenceCount = %v, want 1 on opening", opened.OccurrenceCount)
	}
	if opened.FirstSeenAt == nil || opened.LastSeenAt == nil || opened.FirstSeenAt.Before(before) || !opened.FirstSeenAt.Equal(*opened.LastSeenAt) {
		t.Fatalf("firstSeenAt/lastSeenAt = %v/%v, want both the time of receipt", opened.FirstSeenAt, opened.LastSeenAt)
	}

	later := time.Now().UTC().Add(time.Hour)
	mustPersist(t, p, conditionEvent(owner, key, "second occurrence").WithOccurrenceCount(3).WithLastSeenAt(later).Build())

	rows = storedEvents(t, p, owner)
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1: an occurrence of an open condition records nothing new", len(rows))
	}
	counted := rows[0]
	if counted.ID != opened.ID {
		t.Fatalf("occurrence counted on %s, want the open event %s", counted.ID, opened.ID)
	}
	if counted.OccurrenceCount == nil || *counted.OccurrenceCount != 4 {
		t.Fatalf("occurrenceCount = %v, want 1 + 3", counted.OccurrenceCount)
	}
	if counted.LastSeenAt == nil || !counted.LastSeenAt.Equal(later) {
		t.Fatalf("lastSeenAt = %v, want the occurrence's %v", counted.LastSeenAt, later)
	}
	if !timeOf(counted.FirstSeenAt).Equal(timeOf(opened.FirstSeenAt)) || counted.Description != "first occurrence" {
		t.Fatalf("an occurrence changed fields it must keep: firstSeenAt %v, description %q", counted.FirstSeenAt, counted.Description)
	}
	if !counted.UpdatedAt.After(opened.UpdatedAt) && !counted.UpdatedAt.Equal(opened.UpdatedAt) {
		t.Fatalf("updatedAt moved backwards: %v -> %v", opened.UpdatedAt, counted.UpdatedAt)
	}

	// An occurrence reported as seen earlier than the condition's latest one
	// leaves lastSeenAt at the later time.
	mustPersist(t, p, conditionEvent(owner, key, "late report").WithLastSeenAt(later.Add(-30*time.Minute)).Build())
	rows = storedEvents(t, p, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != 5 || !timeOf(rows[0].LastSeenAt).Equal(later) {
		t.Fatalf("got %+v, want one row counting 5 with lastSeenAt %v", rows, later)
	}
}

// PersistEvent returns the event as stored and what it did: the new row when
// one is created, and the condition's open row, as updated, when an occurrence
// is counted on it or it is resolved. That is the event a caller broadcasts.
func TestPersistEventReturnsTheStoredEvent(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/r/unreachable"

	opened := mustPersist(t, p, conditionEvent(owner, key, "unreachable").Build())
	if rows := storedEvents(t, p, owner); opened.Disposition != EventCreated || len(rows) != 1 || opened.Event.ID != rows[0].ID {
		t.Fatalf("opening: %s %s, want created and the stored row %+v", opened.Disposition, opened.Event.ID, rows)
	}

	sent := conditionEvent(owner, key, "still unreachable").Build()
	counted := mustPersist(t, p, sent)
	if counted.Disposition != EventIncremented || counted.Event.ID != opened.Event.ID || count(counted.Event.OccurrenceCount) != 2 {
		t.Fatalf("occurrence: %s %s count %d, want incremented %s count 2 (sent %s)",
			counted.Disposition, counted.Event.ID, count(counted.Event.OccurrenceCount), opened.Event.ID, sent.ID)
	}

	resolvedAt := time.Now().UTC().Add(time.Minute).Truncate(time.Second)
	sent = events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("reachable again").Resolves(key, resolvedAt).Build()
	resolved := mustPersist(t, p, sent)
	if resolved.Disposition != EventResolved || resolved.Event.ID != opened.Event.ID ||
		!timeOf(resolved.Event.ResolvedAt).Equal(resolvedAt) || resolved.Event.Description != "reachable again" ||
		count(resolved.Event.OccurrenceCount) != 2 {
		t.Fatalf("resolution: %s %+v, want resolved %s with the resolution applied", resolved.Disposition, resolved.Event, opened.Event.ID)
	}

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || rows[0].ID != resolved.Event.ID || !timeOf(rows[0].ResolvedAt).Equal(resolvedAt) {
		t.Fatalf("got %+v, want the one stored row the persister returned", rows)
	}
}

func TestPersistEventOpeningTakesTheProducersTimes(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	lastSeen := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	mustPersist(t, p, conditionEvent(owner, "meshery-server/k8s-connection/a/unreachable", "seen once").WithLastSeenAt(lastSeen).Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || !timeOf(rows[0].FirstSeenAt).Equal(lastSeen) || !timeOf(rows[0].LastSeenAt).Equal(lastSeen) {
		t.Fatalf("got %+v, want firstSeenAt defaulted to the payload's lastSeenAt %v", rows, lastSeen)
	}
}

func TestPersistEventResolutionClosesTheOpenCondition(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/b/unreachable"

	mustPersist(t, p, conditionEvent(owner, key, "unreachable").Build())
	resolvedAt := time.Now().UTC().Add(time.Minute).Truncate(time.Second)
	mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("reachable again").WithOccurrenceCount(5).Resolves(key, resolvedAt).Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1: a resolution of an open condition records nothing new", len(rows))
	}
	if rows[0].ResolvedAt == nil || !rows[0].ResolvedAt.Equal(resolvedAt) {
		t.Fatalf("resolvedAt = %v, want %v", rows[0].ResolvedAt, resolvedAt)
	}
	// The closed event says what cleared and when it was last seen.
	if rows[0].Severity != events.Success || rows[0].Description != "reachable again" {
		t.Fatalf("closed event severity %q description %q, want the resolution's", rows[0].Severity, rows[0].Description)
	}
	if !timeOf(rows[0].LastSeenAt).Equal(resolvedAt) {
		t.Fatalf("lastSeenAt = %v, want the resolution time %v", rows[0].LastSeenAt, resolvedAt)
	}
	if rows[0].Category != "connection" || count(rows[0].OccurrenceCount) != 1 {
		t.Fatalf("a resolution changed fields it must keep (it is a close, not an occurrence): %+v", rows[0])
	}

	// Once closed, the next occurrence opens a new condition.
	mustPersist(t, p, conditionEvent(owner, key, "unreachable again").Build())
	rows = storedEvents(t, p, owner)
	if len(rows) != 2 || rows[1].ResolvedAt != nil || count(rows[1].OccurrenceCount) != 1 || rows[1].Description != "unreachable again" {
		t.Fatalf("got %+v, want a second, open event with one occurrence", rows)
	}
}

// A resolution reported as seen before the condition's latest occurrence
// leaves lastSeenAt at the later time.
func TestPersistEventResolutionKeepsTheLaterLastSeenAt(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/h/unreachable"

	latest := time.Now().UTC().Add(time.Hour).Truncate(time.Second)
	mustPersist(t, p, conditionEvent(owner, key, "unreachable").WithLastSeenAt(latest).Build())
	mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("restored").Resolves(key, latest.Add(-time.Minute)).Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || !timeOf(rows[0].LastSeenAt).Equal(latest) {
		t.Fatalf("got %+v, want lastSeenAt kept at %v", rows, latest)
	}
}

func TestPersistEventResolutionWithNoOpenConditionIsRecorded(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("reachable").Resolves("meshery-server/k8s-connection/c/unreachable", time.Now()).Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || rows[0].ResolvedAt == nil {
		t.Fatalf("got %+v, want one already resolved event", rows)
	}
	resolvedAt := timeOf(rows[0].ResolvedAt)
	if count(rows[0].OccurrenceCount) != 1 || !timeOf(rows[0].FirstSeenAt).Equal(resolvedAt) || !timeOf(rows[0].LastSeenAt).Equal(resolvedAt) {
		t.Fatalf("got count %d firstSeenAt %v lastSeenAt %v, want 1 and both at resolvedAt %v",
			count(rows[0].OccurrenceCount), rows[0].FirstSeenAt, rows[0].LastSeenAt, resolvedAt)
	}
}

// A resolution recorded with no open condition keeps the occurrence fields the
// producer sent.
func TestPersistEventResolutionWithNoOpenConditionKeepsTheProducersTimes(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	firstSeen := time.Date(2026, 11, 6, 9, 0, 0, 0, time.UTC)
	lastSeen := firstSeen.Add(10 * time.Minute)
	resolvedAt := lastSeen.Add(time.Minute)
	mustPersist(t, p, events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("restored").WithOccurrenceCount(4).
		WithFirstSeenAt(firstSeen).WithLastSeenAt(lastSeen).
		Resolves("meshery-server/k8s-connection/i/unreachable", resolvedAt).Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != 4 ||
		!timeOf(rows[0].FirstSeenAt).Equal(firstSeen) || !timeOf(rows[0].LastSeenAt).Equal(lastSeen) {
		t.Fatalf("got %+v, want the producer's count 4, firstSeenAt %v and lastSeenAt %v", rows, firstSeen, lastSeen)
	}
}

func TestPersistEventScopesConditionsToTheOwner(t *testing.T) {
	p := newCorrelationTestPersister(t)
	alice, bob := mustUUID(t), mustUUID(t)
	key := "meshery-server/k8s-connection/d/unreachable"

	mustPersist(t, p, conditionEvent(alice, key, "alice").Build())
	mustPersist(t, p, conditionEvent(bob, key, "bob").Build())

	if a, b := storedEvents(t, p, alice), storedEvents(t, p, bob); len(a) != 1 || len(b) != 1 {
		t.Fatalf("got %d and %d rows, want one open condition per owner", len(a), len(b))
	}
}

// An event that identifies no user is recorded without being correlated, as
// createEvent does for a call that identifies no user.
func TestPersistEventWithoutOwnerIsNotCorrelated(t *testing.T) {
	p := newCorrelationTestPersister(t)
	system := mustUUID(t)
	key := "meshery-server/registry/seed/failed"

	for i := 0; i < 2; i++ {
		mustPersist(t, p, events.NewEvent().FromSystem(system).WithCategory("registry").WithAction("seed").
			WithSeverity(events.Error).WithCorrelationKey(key).Build())
	}

	var count int64
	if err := p.DB.Model(&events.Event{}).Where("correlation_key = ?", key).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("got %d rows, want 2 uncorrelated system events", count)
	}
}

// An event owned by the zero UUID (what ResolveUserAndSystemID yields when no
// user is on the request) identifies no user, so every open occurrence of its
// condition is recorded uncorrelated rather than rejected by the open-condition
// index.
func TestPersistEventWithZeroOwnerIsNotCorrelated(t *testing.T) {
	p := newCorrelationTestPersister(t)
	key := "meshery-server/k8s-connection/z/unreachable"

	for i := 0; i < 2; i++ {
		mustPersist(t, p, conditionEvent(core.Uuid{}, key, "zero owner").Build())
	}

	var count int64
	if err := p.DB.Model(&events.Event{}).Where("correlation_key = ?", key).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 2 {
		t.Fatalf("got %d rows, want 2 uncorrelated zero-owner events", count)
	}
}

// A row an older server recorded open without a count counts as one
// occurrence.
func TestPersistEventCountsAnOpenRowWithoutACount(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/e/unreachable"

	legacy := conditionEvent(owner, key, "recorded without a count").Build()
	if err := p.DB.Create(legacy).Error; err != nil {
		t.Fatalf("seed open row: %v", err)
	}
	mustPersist(t, p, conditionEvent(owner, key, "occurrence").Build())

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != 2 {
		t.Fatalf("got %+v, want the open row counted to 2", rows)
	}
}

func TestOpenEventCorrelationIndexAllowsOneOpenEventPerOwnerAndKey(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/f/unreachable"

	if err := p.DB.Create(conditionEvent(owner, key, "open").Build()).Error; err != nil {
		t.Fatalf("first open row: %v", err)
	}
	err := p.DB.Create(conditionEvent(owner, key, "second open").Build()).Error
	if err == nil {
		t.Fatal("a second open event with the same owner and key was accepted")
	}
	if !isUniqueViolation(err) {
		t.Fatalf("isUniqueViolation(%v) = false, want true so PersistEvent retries", err)
	}

	// Resolved rows and rows without a key stay outside the index.
	resolved := conditionEvent(owner, key, "resolved").WithResolvedAt(time.Now()).Build()
	if err := p.DB.Create(resolved).Error; err != nil {
		t.Fatalf("resolved row with the same key: %v", err)
	}
	for i := 0; i < 2; i++ {
		plain := events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").WithSeverity(events.Error).Build()
		if err := p.DB.Create(plain).Error; err != nil {
			t.Fatalf("row without a key: %v", err)
		}
	}
}

// Concurrent occurrences of one condition (the K8s middleware discovers each
// connection in its own goroutine, on every request) open it once and count
// every occurrence.
func TestPersistEventCountsConcurrentOccurrencesOnce(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/g/unreachable"

	const n = 16
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.PersistEvent(*conditionEvent(owner, key, "occurrence").Build(), "")
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("PersistEvent: %v", err)
		}
	}

	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != n {
		t.Fatalf("got %+v, want one row counting %d occurrences", rows, n)
	}
}

func TestConditionKeyQualifiesBySourceAndComponent(t *testing.T) {
	cases := []struct {
		component, subject, condition, want string
	}{
		{"k8s-connection", "5b0b1c2d", "unreachable", "meshery-server/k8s-connection/5b0b1c2d/unreachable"},
		{" k8s-connection ", "a/b", "", "meshery-server/k8s-connection/a-b/unknown"},
	}
	for _, c := range cases {
		if got := ConditionKey(c.component, c.subject, c.condition); got != c.want {
			t.Errorf("ConditionKey(%q, %q, %q) = %q, want %q", c.component, c.subject, c.condition, got, c.want)
		}
	}
}

// receiveSystemEvent returns the next event published on ch.
func receiveSystemEvent(t *testing.T, ch chan interface{}) *events.Event {
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

// A system event constructed without the builder carries no ID; it is
// broadcast as stored, with the stored row's ID and the server's source.
func TestPersistAndPublishSystemEventBroadcastsTheStoredRow(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	broadcaster := NewBroadcaster("test")
	ch, unsubscribe := broadcaster.Subscribe(owner)
	defer unsubscribe()

	sent := events.NewEvent().FromOwner(owner).WithCategory("registration").WithAction("get_summary").
		WithSeverity(events.Informational).WithDescription("system event").Build()
	// events.NewEvent mints an ID; an event constructed without the
	// builder carries none, which is the case under test.
	sent.ID = uuid.Nil
	if err := PersistAndPublishSystemEvent(p, broadcaster, owner, sent); err != nil {
		t.Fatalf("PersistAndPublishSystemEvent: %v", err)
	}

	got := receiveSystemEvent(t, ch)
	rows := storedEvents(t, p, owner)
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	if got.ID.IsNil() || got.ID != rows[0].ID {
		t.Fatalf("broadcast ID %s, stored ID %s; want the stored, non-zero ID", got.ID, rows[0].ID)
	}
	if got.Source == nil || *got.Source != EventSourceServer {
		t.Fatalf("broadcast source %v, want %q", got.Source, EventSourceServer)
	}
}

// A system event that counts an occurrence on, or resolves, an open condition
// is broadcast as the condition's stored event, never as the event sent.
func TestPersistAndPublishSystemEventBroadcastsTheCondition(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/c/unreachable"
	broadcaster := NewBroadcaster("test")
	ch, unsubscribe := broadcaster.Subscribe(owner)
	defer unsubscribe()

	if err := PersistAndPublishSystemEvent(p, broadcaster, owner, conditionEvent(owner, key, "unreachable").Build()); err != nil {
		t.Fatalf("open: %v", err)
	}
	opened := receiveSystemEvent(t, ch)

	occurrence := conditionEvent(owner, key, "still unreachable").Build()
	if err := PersistAndPublishSystemEvent(p, broadcaster, owner, occurrence); err != nil {
		t.Fatalf("occurrence: %v", err)
	}
	if got := receiveSystemEvent(t, ch); got.ID != opened.ID || count(got.OccurrenceCount) != 2 {
		t.Fatalf("occurrence broadcast ID %s count %d, want the condition %s with 2", got.ID, count(got.OccurrenceCount), opened.ID)
	}

	resolution := events.NewEvent().FromOwner(owner).WithCategory("connection").WithAction("update").
		WithSeverity(events.Success).WithDescription("reachable again").Resolves(key, time.Now().UTC()).Build()
	if err := PersistAndPublishSystemEvent(p, broadcaster, owner, resolution); err != nil {
		t.Fatalf("resolution: %v", err)
	}
	if got := receiveSystemEvent(t, ch); got.ID != opened.ID || got.ResolvedAt == nil {
		t.Fatalf("resolution broadcast ID %s resolvedAt %v, want the condition %s, resolved", got.ID, got.ResolvedAt, opened.ID)
	}
	if rows := storedEvents(t, p, owner); len(rows) != 1 {
		t.Fatalf("got %d rows, want the one condition", len(rows))
	}
}

// A new event is inserted once, with the ID its builder minted, and a retry of
// the same event value records nothing new: it is answered as a replay with
// the stored event.
func TestPersistEventInsertsANewEventOnceWithItsBuilderID(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	sent := events.NewEvent().FromOwner(owner).WithCategory("design").WithAction("save").
		WithSeverity(events.Informational).WithDescription("saved").Build()

	first := mustPersist(t, p, sent)
	if first.Disposition != EventCreated || first.Event.ID != sent.ID {
		t.Fatalf("first persist = %s with ID %s, want created with the builder's ID %s", first.Disposition, first.Event.ID, sent.ID)
	}
	retry := mustPersist(t, p, sent)
	if retry.Disposition != EventReplayed || retry.Event.ID != sent.ID {
		t.Fatalf("retry = %s with ID %s, want replayed with ID %s", retry.Disposition, retry.Event.ID, sent.ID)
	}
	if rows := storedEvents(t, p, owner); len(rows) != 1 || rows[0].ID != sent.ID {
		t.Fatalf("got %+v, want one row with the builder's ID", rows)
	}
}

// A retry of the event that opened a condition is a replay, not a second
// occurrence; a retry of a resolution recorded with no open condition is a
// replay, not a second resolved row.
func TestPersistEventRetryOfACorrelatedEventIsAReplay(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	opening := conditionEvent(owner, "k", "down").Build()
	mustPersist(t, p, opening)
	retry := mustPersist(t, p, opening)
	if retry.Disposition != EventReplayed || retry.Event.ID != opening.ID {
		t.Fatalf("retry of the opening event = %s with ID %s, want replayed with ID %s", retry.Disposition, retry.Event.ID, opening.ID)
	}
	if rows := storedEvents(t, p, owner); len(rows) != 1 || count(rows[0].OccurrenceCount) != 1 {
		t.Fatalf("got %+v, want the one open row still at count 1", rows)
	}

	other := mustUUID(t)
	resolvedAt := time.Now().UTC()
	resolution := conditionEvent(other, "k", "cleared").WithResolvedAt(resolvedAt).Build()
	mustPersist(t, p, resolution)
	again := mustPersist(t, p, resolution)
	if again.Disposition != EventReplayed || again.Event.ID != resolution.ID {
		t.Fatalf("retry of the resolution = %s with ID %s, want replayed with ID %s", again.Disposition, again.Event.ID, resolution.ID)
	}
	if rows := storedEvents(t, p, other); len(rows) != 1 {
		t.Fatalf("got %d rows, want the one resolved row", len(rows))
	}
}

// An occurrenceCount below 1 counts as one occurrence, both for the event that
// opens a condition and for one counted on it: a server producer's zero or
// negative count never leaves the count unchanged or decrements it.
func TestPersistEventClampsTheOccurrenceCountToOne(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	for _, n := range []int{0, -3} {
		n := n
		event := conditionEvent(owner, "k", "down").Build()
		event.OccurrenceCount = &n
		mustPersist(t, p, event)
	}
	rows := storedEvents(t, p, owner)
	if len(rows) != 1 || count(rows[0].OccurrenceCount) != 2 {
		t.Fatalf("got %+v, want one open row at count 2", rows)
	}
}

func withKey(event *events.Event, key string) *events.Event {
	event.IdempotencyKey = &key
	return event
}

// A create carrying an idempotencyKey its caller already used records and
// changes nothing, whatever the rest of the payload says, and is answered with
// the event the key is bound to: the recorded event, or the condition's event
// for a key that counted an occurrence or resolved it.
func TestPersistEventReplaysAnIdempotencyKey(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)

	created := mustPersist(t, p, withKey(events.NewEvent().FromOwner(owner).WithCategory("perf").WithAction("run").
		WithSeverity(events.Informational).WithDescription("first").Build(), "run-1"))
	retry := mustPersist(t, p, withKey(events.NewEvent().FromOwner(owner).WithCategory("perf").WithAction("run").
		WithSeverity(events.Error).WithDescription("different payload").Build(), "run-1"))
	if retry.Disposition != EventReplayed || retry.Event.ID != created.Event.ID || retry.Event.Description != "first" {
		t.Fatalf("retry = %s %s %q, want replayed %s %q", retry.Disposition, retry.Event.ID, retry.Event.Description, created.Event.ID, "first")
	}

	opened := mustPersist(t, p, conditionEvent(owner, "k", "down").Build())
	counted := mustPersist(t, p, withKey(conditionEvent(owner, "k", "down").Build(), "occ-1"))
	if counted.Disposition != EventIncremented {
		t.Fatalf("keyed occurrence = %s, want incremented", counted.Disposition)
	}
	again := mustPersist(t, p, withKey(conditionEvent(owner, "k", "down").Build(), "occ-1"))
	if again.Disposition != EventReplayed || again.Event.ID != opened.Event.ID || count(again.Event.OccurrenceCount) != 2 {
		t.Fatalf("retried occurrence = %s %s count %d, want replayed %s count 2",
			again.Disposition, again.Event.ID, count(again.Event.OccurrenceCount), opened.Event.ID)
	}

	resolution := withKey(conditionEvent(owner, "k", "up").WithResolvedAt(time.Now().UTC()).Build(), "res-1")
	if got := mustPersist(t, p, resolution); got.Disposition != EventResolved {
		t.Fatalf("keyed resolution = %s, want resolved", got.Disposition)
	}
	retried := withKey(conditionEvent(owner, "k", "up").WithResolvedAt(time.Now().UTC()).Build(), "res-1")
	if got := mustPersist(t, p, retried); got.Disposition != EventReplayed || got.Event.ID != opened.Event.ID {
		t.Fatalf("retried resolution = %s %s, want replayed %s", got.Disposition, got.Event.ID, opened.Event.ID)
	}
	if rows := storedEvents(t, p, owner); len(rows) != 2 {
		t.Fatalf("got %d rows, want the created event and the condition", len(rows))
	}

	// Keys are scoped to the caller: another owner's key is its own.
	other := mustUUID(t)
	if got := mustPersist(t, p, withKey(events.NewEvent().FromOwner(other).WithCategory("perf").WithAction("run").
		WithSeverity(events.Informational).Build(), "run-1")); got.Disposition != EventCreated {
		t.Fatalf("another caller's key = %s, want created", got.Disposition)
	}
}

// Of concurrent creates carrying one key, exactly one is recorded.
func TestPersistEventConcurrentIdempotencyKeyRecordsOnce(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := p.PersistEvent(*withKey(events.NewEvent().FromOwner(owner).WithCategory("perf").WithAction("run").
				WithSeverity(events.Informational).Build(), "run-1"), ""); err != nil {
				t.Errorf("PersistEvent: %v", err)
			}
		}()
	}
	wg.Wait()
	if rows := storedEvents(t, p, owner); len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
}

// A retry of an event resolved since, sent while a newer event has opened the
// same condition, is a replay of the resolved event: the newer condition is
// neither counted on nor resolved by it.
func TestPersistEventRetryOfAResolvedEventLeavesTheNewerConditionAlone(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	key := "meshery-server/k8s-connection/c/unreachable"

	// A opens the condition and is resolved; B then opens it again.
	opening := conditionEvent(owner, key, "down").Build()
	mustPersist(t, p, opening)
	mustPersist(t, p, conditionEvent(owner, key, "up").WithResolvedAt(time.Now().UTC()).Build())
	newer := mustPersist(t, p, conditionEvent(owner, key, "down again").Build())
	if newer.Disposition != EventCreated || newer.Event.ID == opening.ID {
		t.Fatalf("newer condition = %s %s, want a new event", newer.Disposition, newer.Event.ID)
	}

	retry := mustPersist(t, p, opening)
	if retry.Disposition != EventReplayed || retry.Event.ID != opening.ID || retry.Event.ResolvedAt == nil {
		t.Fatalf("retry of the resolved event = %s %s resolvedAt %v, want replayed %s, resolved",
			retry.Disposition, retry.Event.ID, retry.Event.ResolvedAt, opening.ID)
	}
	assertOpenConditionUntouched(t, p, owner, newer.Event.ID)

	// A resolution recorded with no condition open, retried once a newer
	// condition is open, is a replay of the resolution, never resolving it.
	other := mustUUID(t)
	resolution := conditionEvent(other, key, "cleared").WithResolvedAt(time.Now().UTC()).Build()
	if got := mustPersist(t, p, resolution); got.Disposition != EventCreated {
		t.Fatalf("resolution with none open = %s, want created", got.Disposition)
	}
	opened := mustPersist(t, p, conditionEvent(other, key, "down").Build())
	again := mustPersist(t, p, resolution)
	if again.Disposition != EventReplayed || again.Event.ID != resolution.ID {
		t.Fatalf("retry of the resolution = %s %s, want replayed %s", again.Disposition, again.Event.ID, resolution.ID)
	}
	assertOpenConditionUntouched(t, p, other, opened.Event.ID)
}

// assertOpenConditionUntouched fails unless the owner's event id is still
// open with the one occurrence that opened it.
func assertOpenConditionUntouched(t *testing.T, p *EventsPersister, owner core.Uuid, id uuid.UUID) {
	t.Helper()
	for _, row := range storedEvents(t, p, owner) {
		if row.ID != id {
			continue
		}
		if row.ResolvedAt != nil || count(row.OccurrenceCount) != 1 {
			t.Fatalf("newer condition %s resolvedAt %v count %d, want open at count 1", id, row.ResolvedAt, count(row.OccurrenceCount))
		}
		return
	}
	t.Fatalf("newer condition %s not stored", id)
}

// Two creates carrying one key whose binding is stale (its event was deleted)
// record one event between them. The callback holds the second create's
// removal of the stale binding until the first has replaced it, the
// interleaving in which removing the binding by owner and key alone would
// delete the first create's binding and record a second event for the key.
func TestPersistEventConcurrentCreatesOverAStaleBindingRecordOnce(t *testing.T) {
	p := newCorrelationTestPersister(t)
	owner := mustUUID(t)
	const key = "run-1"

	stale := mustPersist(t, p, withKey(events.NewEvent().FromOwner(owner).WithCategory("perf").WithAction("run").
		WithSeverity(events.Informational).Build(), key))
	if err := p.DB.Where("id = ?", stale.Event.ID).Delete(&events.Event{}).Error; err != nil {
		t.Fatalf("delete the bound event: %v", err)
	}

	const wait = 5 * time.Second
	var arrivals int32
	var mu sync.Mutex
	bothRead := make(chan struct{})
	rebound := func() bool {
		var binding EventIdempotencyBinding
		found := p.DB.Session(&gorm.Session{NewDB: true}).
			Where("owner = ? AND idempotency_key = ? AND event_id <> ?", owner, key, stale.Event.ID).
			Limit(1).Find(&binding)
		return found.Error == nil && found.RowsAffected > 0
	}
	// Registered before the delete's own transaction begins, so a held
	// create holds no connection of the one-connection pool.
	if err := p.DB.Callback().Delete().Before("gorm:begin_transaction").Register("test:stale-binding-race", func(tx *gorm.DB) {
		if tx.Statement.Table != (EventIdempotencyBinding{}).TableName() {
			return
		}
		mu.Lock()
		arrivals++
		arrival := arrivals
		mu.Unlock()
		switch arrival {
		case 1:
			// Both creates read the stale binding before either removes it.
			select {
			case <-bothRead:
			case <-time.After(wait):
				t.Errorf("the second create never read the stale binding")
			}
		case 2:
			close(bothRead)
			deadline := time.Now().Add(wait)
			for !rebound() {
				if time.Now().After(deadline) {
					t.Errorf("the first create never rebound the key")
					return
				}
				time.Sleep(5 * time.Millisecond)
			}
		}
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	results := make(chan PersistedEvent, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			persisted, err := p.PersistEvent(*withKey(events.NewEvent().FromOwner(owner).WithCategory("perf").WithAction("run").
				WithSeverity(events.Informational).Build(), key), "")
			if err != nil {
				t.Errorf("PersistEvent: %v", err)
				return
			}
			results <- persisted
		}()
	}
	wg.Wait()
	close(results)

	dispositions := map[EventDisposition]int{}
	var ids []uuid.UUID
	for persisted := range results {
		dispositions[persisted.Disposition]++
		ids = append(ids, persisted.Event.ID)
	}
	if dispositions[EventCreated] != 1 || dispositions[EventReplayed] != 1 || len(ids) != 2 || ids[0] != ids[1] {
		t.Fatalf("dispositions %v ids %v, want one created and one replay of it", dispositions, ids)
	}
	if rows := storedEvents(t, p, owner); len(rows) != 1 || rows[0].ID != ids[0] {
		t.Fatalf("got %d rows, want the one created event", len(rows))
	}
}
