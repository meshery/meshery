package models

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
)

// TestBroadcastDeliversToSubscriber is the happy path: a published value reaches
// a live subscriber.
func TestBroadcastDeliversToSubscriber(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	b.Publish(id, 42)

	select {
	case got := <-ch:
		if got != 42 {
			t.Fatalf("got %v, want 42", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for published value")
	}
}

// TestBroadcastBurstWithinBufferNoLoss guards the data-loss regression: a burst
// of events published before the subscriber drains must all be delivered, in
// order. The old Publish used a destructive IsClosed check that consumed a
// buffered event and then skipped the send, losing both.
func TestBroadcastBurstWithinBufferNoLoss(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	const n = 10 // well within listenerBufferSize
	for i := 0; i < n; i++ {
		b.Publish(id, i)
	}

	for i := 0; i < n; i++ {
		select {
		case got := <-ch:
			if got != i {
				t.Fatalf("event %d: got %v, want %d (order/loss)", i, got, i)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out after %d/%d events; earlier events were lost", i, n)
		}
	}
}

// TestBroadcastNonBlockingWhenBufferFull guards the slow-consumer regression:
// Publish must never block, even when a subscriber never drains. Overflow beyond
// the buffer is dropped rather than blocking the publisher.
func TestBroadcastNonBlockingWhenBufferFull(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	done := make(chan struct{})
	go func() {
		// Publish more than the buffer can hold, without anyone draining.
		for i := 0; i < listenerBufferSize+64; i++ {
			b.Publish(id, i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on a full subscriber buffer")
	}

	// The buffer retained exactly listenerBufferSize items; the rest were dropped.
	count := 0
	for {
		select {
		case <-ch:
			count++
		default:
			if count != listenerBufferSize {
				t.Fatalf("buffered %d items, want %d", count, listenerBufferSize)
			}
			return
		}
	}
}

// TestBroadcastConcurrentPublishUnsubscribeNoPanic guards the send-on-closed
// panic: publishers and churning subscribers race for the same id. Run with
// -race. A regression (unlocked send racing a close) panics and fails the test.
func TestBroadcastConcurrentPublishUnsubscribeNoPanic(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					b.Publish(id, 1)
				}
			}
		}()
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ch, unsubscribe := b.Subscribe(id)
					select {
					case <-ch:
					default:
					}
					unsubscribe()
				}
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestBroadcastUnsubscribeIsIdempotent ensures a second unsubscribe is a no-op
// rather than a double-close panic.
func TestBroadcastUnsubscribeIsIdempotent(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	_, unsubscribe := b.Subscribe(id)
	unsubscribe()
	unsubscribe() // must not panic
}

// TestBroadcastUnsubscribeStopsDelivery ensures that after unsubscribe the
// listener channel is closed and receives no further events.
func TestBroadcastUnsubscribeStopsDelivery(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	unsubscribe()

	// Publishing to an id whose only listener has unsubscribed must not panic
	// and must not deliver.
	b.Publish(id, 7)

	select {
	case v, open := <-ch:
		if open {
			t.Fatalf("received %v after unsubscribe; channel should be drained/closed", v)
		}
	case <-time.After(time.Second):
		t.Fatal("closed channel did not yield")
	}
}

// TestBroadcastUnsubscribePrunesEmptyClientsMap ensures that when all subscribers
// for a user unsubscribe, the user entry is deleted from the clients map to prevent
// an unbounded memory leak.
func TestBroadcastUnsubscribePrunesEmptyClientsMap(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	_, unsubscribe1 := b.Subscribe(id)
	_, unsubscribe2 := b.Subscribe(id)

	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("expected client entry to be present in map while subscribed")
	}

	unsubscribe1()
	// Still has one subscriber, entry must remain in map.
	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("client entry pruned prematurely while one subscriber remains")
	}

	unsubscribe2()
	// All subscribers gone, entry must be pruned from map.
	if _, ok := b.clients.Load(id); ok {
		t.Fatal("expected client entry to be pruned from map after all subscribers unsubscribed")
	}
}

// TestBroadcastResubscribeAfterUnsubscribe ensures that resubscribing after
// full eviction creates a fresh active entry and receives published data normally.
func TestBroadcastResubscribeAfterUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch1, unsub1 := b.Subscribe(id)
	unsub1()

	if _, ok := b.clients.Load(id); ok {
		t.Fatal("entry should be pruned after unsub")
	}

	ch2, unsub2 := b.Subscribe(id)
	defer unsub2()

	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("entry should be present after resubscribe")
	}

	b.Publish(id, "resubscribed-event")

	select {
	case msg := <-ch2:
		if msg != "resubscribed-event" {
			t.Fatalf("got %v, want resubscribed-event", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on resubscribed listener")
	}

	select {
	case _, open := <-ch1:
		if open {
			t.Fatal("old channel should be closed")
		}
	default:
	}
}

// TestBroadcastConcurrentPruneAndSubscribe hammers Subscribe, Publish, and Unsubscribe
// on the same ID to verify that atomic eviction does not race with new subscribers
// or panic with -race.
func TestBroadcastConcurrentPruneAndSubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Publishers
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					b.Publish(id, val)
				}
			}
		}(i)
	}

	// Churning subscribers triggering eviction repeatedly
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ch, unsubscribe := b.Subscribe(id)
					select {
					case <-ch:
					default:
					}
					unsubscribe()
				}
			}
		}()
	}

	time.Sleep(300 * time.Millisecond)
	close(stop)
	wg.Wait()

	// After all churn finishes, resubscribing must work cleanly
	ch, unsub := b.Subscribe(id)
	defer unsub()

	b.Publish(id, 999)
	select {
	case v := <-ch:
		if v != 999 {
			t.Fatalf("got %v, want 999", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on final subscriber")
	}
}

func mustReceiveReplayEvent(t *testing.T, ch <-chan interface{}) ReplayEvent {
	t.Helper()
	select {
	case got, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before receiving ReplayEvent")
		}
		re, ok := got.(ReplayEvent)
		if !ok {
			t.Fatalf("expected ReplayEvent, got %T", got)
		}
		return re
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for ReplayEvent")
		return ReplayEvent{}
	}
}

func newTestEventWithID(userID uuid.UUID, desc string) *events.Event {
	ev := events.NewEvent().
		FromOwner(userID).
		WithCategory("test").
		WithAction("replay").
		WithSeverity(events.Informational).
		WithDescription(desc).
		Build()
	ev.ID = uuid.Must(uuid.NewV4())
	return ev
}

func TestBroadcastSubscribeFromReplaysMissedEvents(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	e1 := newTestEventWithID(userID, "e1")
	e2 := newTestEventWithID(userID, "e2")
	e3 := newTestEventWithID(userID, "e3")

	// Published while no subscriber is connected.
	b.Publish(userID, e1)
	b.Publish(userID, e2)
	b.Publish(userID, e3)

	replay, ch, unsub := b.SubscribeFrom(userID, e1.ID.String())
	defer unsub()

	if len(replay) != 2 {
		t.Fatalf("replay len = %d, want 2", len(replay))
	}
	if replay[0].Data.(*events.Event).ID != e2.ID || replay[1].Data.(*events.Event).ID != e3.ID {
		t.Fatalf("unexpected replay order/IDs: got [%v, %v], want [%v, %v]",
			replay[0].Data.(*events.Event).ID, replay[1].Data.(*events.Event).ID, e2.ID, e3.ID)
	}

	e4 := newTestEventWithID(userID, "e4")
	b.Publish(userID, e4)
	re := mustReceiveReplayEvent(t, ch)
	if re.Data.(*events.Event).ID != e4.ID || re.Cursor == "" {
		t.Fatalf("live event = (%v, %q), want ID %v with non-empty cursor", re.Data.(*events.Event).ID, re.Cursor, e4.ID)
	}
}

func TestBroadcastSubscribeFromUppercaseBareUUID(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	e1 := newTestEventWithID(userID, "e1")
	e2 := newTestEventWithID(userID, "e2")
	e3 := newTestEventWithID(userID, "e3")

	b.Publish(userID, e1)
	b.Publish(userID, e2)
	b.Publish(userID, e3)

	// Uppercase bare UUID must match the lowercase stored event ID.
	replay, _, unsub := b.SubscribeFrom(userID, strings.ToUpper(e1.ID.String()))
	defer unsub()
	if len(replay) != 2 {
		t.Fatalf("uppercase bare UUID replay len = %d, want 2", len(replay))
	}
	if replay[0].Data.(*events.Event).ID != e2.ID || replay[1].Data.(*events.Event).ID != e3.ID {
		t.Fatalf("unexpected replay IDs: got [%v, %v], want [%v, %v]",
			replay[0].Data.(*events.Event).ID, replay[1].Data.(*events.Event).ID, e2.ID, e3.ID)
	}

	// Unknown valid uppercase bare UUID replays full retained buffer.
	unknownUpper := strings.ToUpper(uuid.Must(uuid.NewV4()).String())
	fullReplay, _, unsub2 := b.SubscribeFrom(userID, unknownUpper)
	defer unsub2()
	if len(fullReplay) != 3 {
		t.Fatalf("unknown uppercase bare UUID replay len = %d, want 3", len(fullReplay))
	}
}

func TestBroadcastSubscribeFromRepeatedEventIDs(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	_, liveCh, liveUnsub := b.SubscribeFrom(userID, "")
	defer liveUnsub()

	sharedID := uuid.Must(uuid.NewV4())
	firstA := newTestEventWithID(userID, "first-A")
	firstA.ID = sharedID
	eventB := newTestEventWithID(userID, "event-B")
	secondA := newTestEventWithID(userID, "second-A")
	secondA.ID = sharedID

	b.Publish(userID, firstA)
	firstMsg := mustReceiveReplayEvent(t, liveCh)

	b.Publish(userID, eventB)
	mustReceiveReplayEvent(t, liveCh)
	b.Publish(userID, secondA)
	mustReceiveReplayEvent(t, liveCh)

	if firstMsg.Cursor == "" || firstMsg.Cursor == sharedID.String() {
		t.Fatalf("expected sequence-qualified cursor, got %q", firstMsg.Cursor)
	}

	// Reconnecting with the first A's cursor must replay B and the second A.
	replay, _, unsub := b.SubscribeFrom(userID, firstMsg.Cursor)
	defer unsub()

	if len(replay) != 2 {
		t.Fatalf("replay len = %d, want 2 (B and second A)", len(replay))
	}
	if replay[0].Data.(*events.Event).Description != "event-B" ||
		replay[1].Data.(*events.Event).Description != "second-A" {
		t.Fatalf("unexpected replay sequence: got [%q, %q], want [event-B, second-A]",
			replay[0].Data.(*events.Event).Description,
			replay[1].Data.(*events.Event).Description)
	}
}

func TestBroadcastSubscribeFromEdgeCases(t *testing.T) {
	// Zero-value Broadcast must work without panicking.
	var b Broadcast
	userA := uuid.Must(uuid.NewV4())
	userB := uuid.Must(uuid.NewV4())

	eA := newTestEventWithID(userA, "a1")
	eB := newTestEventWithID(userB, "b1")
	nilIDEvent := events.NewEvent().FromOwner(userA).WithDescription("nil-id").Build()
	nilIDEvent.ID = uuid.Nil

	b.Publish(userA, eA)
	b.Publish(userA, nilIDEvent)
	b.Publish(userA, "non-event")
	b.Publish(userB, eB)

	// Empty lastEventID -> nil replay.
	replay, ch, unsub := b.SubscribeFrom(userA, "")
	if replay != nil {
		t.Fatalf("empty lastEventID: got %v, want nil", replay)
	}
	unsub()
	_ = ch

	// Malformed lastEventID -> nil replay.
	replay, _, unsub = b.SubscribeFrom(userA, "not-a-valid-uuid")
	if replay != nil {
		t.Fatalf("malformed lastEventID: got %v, want nil", replay)
	}
	unsub()

	// Cross-user isolation: userB's event ID must never return userB's event,
	// and as an unknown valid UUID on userA it replays only userA's buffer.
	replay, _, unsub = b.SubscribeFrom(userA, eB.ID.String())
	if len(replay) != 1 || replay[0].Data.(*events.Event).ID != eA.ID {
		t.Fatalf("cross-user lastEventID: got %v, want [eA]", replay)
	}
	unsub()

	// Latest event ID -> nil replay (no missed events).
	replay, _, unsub = b.SubscribeFrom(userA, eA.ID.String())
	if replay != nil {
		t.Fatalf("latest lastEventID: got %v, want nil", replay)
	}
	unsub()
}

func TestBroadcastReplayBufferRingWrap(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	total := replayBufferSize + 10
	published := make([]*events.Event, total)
	for i := 0; i < total; i++ {
		published[i] = newTestEventWithID(userID, "ring")
		b.Publish(userID, published[i])
	}

	// Evicted oldest valid UUID replays the full retained buffer as best effort.
	replay, _, unsub := b.SubscribeFrom(userID, published[0].ID.String())
	unsub()
	if len(replay) != replayBufferSize {
		t.Fatalf("evicted event ID returned %d replay items, want %d", len(replay), replayBufferSize)
	}
	if replay[0].Data.(*events.Event).ID != published[10].ID {
		t.Fatalf("replay[0] ID = %v, want %v", replay[0].Data.(*events.Event).ID, published[10].ID)
	}

	// Retained event near tail replays only the events strictly after it.
	anchorIdx := total - 5
	replay, _, unsub = b.SubscribeFrom(userID, published[anchorIdx].ID.String())
	defer unsub()
	if len(replay) != 4 {
		t.Fatalf("replay len = %d, want 4", len(replay))
	}
	for i, item := range replay {
		want := published[anchorIdx+1+i].ID
		got := item.Data.(*events.Event).ID
		if got != want {
			t.Fatalf("replay[%d] ID = %v, want %v", i, got, want)
		}
	}
}

func TestBroadcastHistoryIdleEviction(t *testing.T) {
	b := &Broadcast{
		historyIdleTTL:       20 * time.Millisecond,
		historySweepInterval: time.Millisecond,
	}
	idleUser := uuid.Must(uuid.NewV4())
	activeUser := uuid.Must(uuid.NewV4())
	triggerUser := uuid.Must(uuid.NewV4())

	b.Publish(idleUser, newTestEventWithID(idleUser, "idle"))

	// activeUser has a connected listener, so its history must not be evicted.
	_, unsubActive := b.Subscribe(activeUser)
	defer unsubActive()
	b.Publish(activeUser, newTestEventWithID(activeUser, "active"))

	time.Sleep(35 * time.Millisecond)

	// Trigger lazy sweep via another publish.
	b.Publish(triggerUser, newTestEventWithID(triggerUser, "trigger"))

	if _, ok := b.history.Load(idleUser); ok {
		t.Fatal("expected idle user history entry to be evicted after TTL")
	}
	if _, ok := b.history.Load(activeUser); !ok {
		t.Fatal("expected active user history entry to be retained while listener is connected")
	}
}

func TestBroadcastCursorsUniqueAcrossIdleEviction(t *testing.T) {
	b := &Broadcast{
		historyIdleTTL:       20 * time.Millisecond,
		historySweepInterval: time.Millisecond,
	}
	userID := uuid.Must(uuid.NewV4())
	otherUser := uuid.Must(uuid.NewV4())

	_, liveCh, liveUnsub := b.SubscribeFrom(userID, "")
	sharedID := uuid.Must(uuid.NewV4())
	first := newTestEventWithID(userID, "first")
	first.ID = sharedID

	b.Publish(userID, first)
	firstCursor := mustReceiveReplayEvent(t, liveCh).Cursor
	liveUnsub()

	time.Sleep(35 * time.Millisecond)

	// Trigger eviction of userID's idle buffer.
	b.Publish(otherUser, newTestEventWithID(otherUser, "sweep"))
	if _, ok := b.history.Load(userID); ok {
		t.Fatal("expected idle user history to be evicted")
	}

	// Publish the same event UUID again after buffer recreation.
	_, liveCh2, liveUnsub2 := b.SubscribeFrom(userID, "")
	defer liveUnsub2()
	second := newTestEventWithID(userID, "second")
	second.ID = sharedID
	b.Publish(userID, second)
	secondCursor := mustReceiveReplayEvent(t, liveCh2).Cursor

	if secondCursor == firstCursor {
		t.Fatalf("cursor reused across buffer eviction: %q", secondCursor)
	}

	// Reconnecting with firstCursor must not collide with secondCursor and must replay second.
	replay, _, unsub := b.SubscribeFrom(userID, firstCursor)
	defer unsub()
	if len(replay) != 1 || replay[0].Data.(*events.Event).Description != "second" {
		t.Fatalf("expected [second] on reconnect with old cursor, got %v", replay)
	}
}

func TestBroadcastConcurrentPublishAndSubscribeFromNoRace(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	seed := newTestEventWithID(userID, "seed")
	b.Publish(userID, seed)

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					b.Publish(userID, newTestEventWithID(userID, "concurrent"))
				}
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_, ch, unsub := b.SubscribeFrom(userID, seed.ID.String())
					select {
					case <-ch:
					default:
					}
					unsub()
				}
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestBroadcastSubscribeSynchronizesWithIdleEviction verifies that Subscribe and
// SubscribeFrom(id, "") participate in replay-buffer locking so concurrent idle
// sweeps cannot evict an active subscriber's replay buffer, while stale buffers
// without active listeners are still evicted.
func TestBroadcastSubscribeSynchronizesWithIdleEviction(t *testing.T) {
	b := &Broadcast{
		historyIdleTTL:       15 * time.Millisecond,
		historySweepInterval: time.Millisecond,
	}
	userID := uuid.Must(uuid.NewV4())
	idleUser := uuid.Must(uuid.NewV4())
	sweepUser := uuid.Must(uuid.NewV4())

	// Seed a buffer for userID, let it go past TTL, then subscribe without Last-Event-ID.
	b.Publish(userID, newTestEventWithID(userID, "before-subscribe"))
	b.Publish(idleUser, newTestEventWithID(idleUser, "idle"))
	time.Sleep(25 * time.Millisecond)

	_, ch, unsub := b.SubscribeFrom(userID, "")
	_ = ch

	// Sweep now: idleUser has no listener and is stale -> evicted;
	// userID has an active listener -> must remain retained.
	b.Publish(sweepUser, newTestEventWithID(sweepUser, "sweep"))
	if _, ok := b.history.Load(idleUser); ok {
		t.Fatal("expected stale idle user history without active listener to be evicted")
	}
	if _, ok := b.history.Load(userID); !ok {
		t.Fatal("expected active user history to be retained after SubscribeFrom without Last-Event-ID")
	}
	unsub()

	// Concurrent stress: race Subscribe / SubscribeFrom(id, "") against frequent
	// maybePruneIdleHistory sweeps; whenever a listener is active, its buffer must
	// not be evicted underneath it.
	var wg sync.WaitGroup
	stop := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
				b.Publish(sweepUser, newTestEventWithID(sweepUser, "sweep-concurrent"))
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			b.Publish(userID, newTestEventWithID(userID, "active"))
			var lch chan interface{}
			var lunsub func()
			if i%2 == 0 {
				lch, lunsub = b.Subscribe(userID)
			} else {
				_, lch, lunsub = b.SubscribeFrom(userID, "")
			}
			_ = lch
			if _, ok := b.history.Load(userID); !ok {
				lunsub()
				t.Errorf("iteration %d: active subscriber's history buffer was evicted", i)
				return
			}
			lunsub()
		}
	}()

	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestBroadcastCursorListenerOverflowClosesChannel verifies that when a
// cursor-enabled listener's channel fills up, publishToClients removes it and
// closes the channel after the queued items so the SSE handler can drain and
// close the response to trigger EventSource reconnect, while legacy Subscribe
// listeners preserve non-blocking drop behavior.
func TestBroadcastCursorListenerOverflowClosesChannel(t *testing.T) {
	b := NewBroadcaster("test")
	userID := uuid.Must(uuid.NewV4())

	_, cursorCh, cursorUnsub := b.SubscribeFrom(userID, "")
	defer cursorUnsub()

	legacyCh, legacyUnsub := b.Subscribe(userID)
	defer legacyUnsub()

	total := listenerBufferSize + 5
	published := make([]*events.Event, total)
	done := make(chan struct{})
	go func() {
		for i := 0; i < total; i++ {
			published[i] = newTestEventWithID(userID, "overflow")
			b.Publish(userID, published[i])
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on full cursor-enabled listener channel")
	}

	// Drain the listenerBufferSize queued events from cursorCh in order.
	for i := 0; i < listenerBufferSize; i++ {
		re := mustReceiveReplayEvent(t, cursorCh)
		if re.Data.(*events.Event).ID != published[i].ID {
			t.Fatalf("cursorCh[%d] ID = %v, want %v", i, re.Data.(*events.Event).ID, published[i].ID)
		}
	}

	// After draining queued items, cursorCh must be closed.
	select {
	case v, open := <-cursorCh:
		if open {
			t.Fatalf("expected cursorCh to be closed after queued items, got %v", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for overflowed cursorCh to close")
	}

	// Calling unsubscribe on an already-closed cursor listener must be safe.
	cursorUnsub()

	// Legacy Subscribe listener must still be open and have retained listenerBufferSize items.
	for i := 0; i < listenerBufferSize; i++ {
		select {
		case got := <-legacyCh:
			if got.(*events.Event).ID != published[i].ID {
				t.Fatalf("legacyCh[%d] ID = %v, want %v", i, got.(*events.Event).ID, published[i].ID)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out draining legacyCh item %d", i)
		}
	}
	extra := newTestEventWithID(userID, "after-overflow")
	b.Publish(userID, extra)
	select {
	case got, open := <-legacyCh:
		if !open {
			t.Fatal("legacyCh should remain open after overflow")
		}
		if got.(*events.Event).ID != extra.ID {
			t.Fatalf("legacyCh got %v, want %v", got.(*events.Event).ID, extra.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on legacyCh after overflow")
	}
}
