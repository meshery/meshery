package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshkit/logger"
)

// newTestLogger returns a meshkit logger suitable for use in tests.
// Mirrors the init pattern in newTestRemoteProvider so test files do not
// have to reach into meshkit/logger directly.
func newTestLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	return log
}

// TestProviderTracker_LocalReportsOnlineWithoutNetwork asserts the local
// provider is always reported as ProviderStatusOnline without a network
// probe, satisfying the "Local provider is ever-present and perpetually
// available" requirement.
func TestProviderTracker_LocalReportsOnlineWithoutNetwork(t *testing.T) {
	lp := newTestLocalProviderForTracker(t)
	providers := map[string]Provider{LocalProviderName: lp}
	tracker := NewProviderTracker(providers, lp.Log)

	snap := tracker.Snapshot()
	got, ok := snap[LocalProviderName]
	if !ok {
		t.Fatalf("expected %q in tracker snapshot, got keys %v", LocalProviderName, mapKeys(snap))
	}
	if got.Status != ProviderStatusOnline {
		t.Errorf("local provider initial status = %q, want %q", got.Status, ProviderStatusOnline)
	}

	tracker.VerifyAll(context.Background())
	snap = tracker.Snapshot()
	if snap[LocalProviderName].Status != ProviderStatusOnline {
		t.Errorf("local provider post-verify status = %q, want %q",
			snap[LocalProviderName].Status, ProviderStatusOnline)
	}
}

// TestProviderTracker_RemoteFailureReportsOfflineWithoutBlockingPeers
// asserts each remote is probed concurrently, so an unreachable remote
// does not delay the verdict for a healthy peer. This is the core
// guarantee the SSE chooser relies on.
func TestProviderTracker_RemoteFailureReportsOfflineWithoutBlockingPeers(t *testing.T) {
	log := newTestLogger(t)

	healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providerName":"healthy-cloud","providerType":"remote","capabilities":[{"feature":"sync-prefs","endpoint":"/sync"}]}`))
	}))
	defer healthy.Close()

	good := newTestRemoteProvider(t, healthy.URL)
	good.Initialize()
	// Deliberately unreachable address (RFC5737 documentation block) so the
	// probe always fails; the http.Client timeout caps wait time.
	bad := newTestRemoteProvider(t, "http://203.0.113.1:65530")
	bad.Initialize()

	providers := map[string]Provider{
		"healthy": good,
		"bad":     bad,
	}
	tracker := NewProviderTracker(providers, log)

	done := make(chan struct{})
	go func() {
		tracker.VerifyAll(context.Background())
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("VerifyAll did not complete within 30s; concurrent probe is blocking on the unreachable peer")
	}

	snap := tracker.Snapshot()
	if snap["healthy"].Status != ProviderStatusOnline {
		t.Errorf("healthy.Status = %q, want %q (err=%q)",
			snap["healthy"].Status, ProviderStatusOnline, snap["healthy"].Error)
	}
	if snap["healthy"].Properties.ProviderName != "healthy-cloud" {
		t.Errorf("healthy.ProviderName = %q, want %q",
			snap["healthy"].Properties.ProviderName, "healthy-cloud")
	}
	if snap["bad"].Status != ProviderStatusOffline {
		t.Errorf("bad.Status = %q, want %q", snap["bad"].Status, ProviderStatusOffline)
	}
}

// TestProviderTracker_SubscriberReceivesSnapshotAndUpdates asserts a new
// subscriber sees the current snapshot first (one event per registered
// provider) and then receives subsequent Publish updates, which is the
// contract /api/providers/stream depends on.
func TestProviderTracker_SubscriberReceivesSnapshotAndUpdates(t *testing.T) {
	lp := newTestLocalProviderForTracker(t)
	providers := map[string]Provider{LocalProviderName: lp}
	tracker := NewProviderTracker(providers, lp.Log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events, unsubscribe := tracker.Subscribe(ctx)
	defer unsubscribe()

	// First event: snapshot replay for the local provider.
	select {
	case evt := <-events:
		if evt.Key != LocalProviderName {
			t.Fatalf("snapshot event Key = %q, want %q", evt.Key, LocalProviderName)
		}
		if evt.Status != ProviderStatusOnline {
			t.Fatalf("snapshot event Status = %q, want %q", evt.Status, ProviderStatusOnline)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive initial snapshot event within 2s")
	}

	// Second event: a manual Publish should reach the subscriber.
	tracker.Publish(ProviderStatusEvent{
		Key:    "synthetic",
		Status: ProviderStatusOffline,
		Error:  "test injected",
	})

	select {
	case evt := <-events:
		if evt.Key != "synthetic" {
			t.Fatalf("published event Key = %q, want %q", evt.Key, "synthetic")
		}
		if evt.Status != ProviderStatusOffline {
			t.Errorf("published event Status = %q, want %q", evt.Status, ProviderStatusOffline)
		}
		if evt.Error != "test injected" {
			t.Errorf("published event Error = %q, want %q", evt.Error, "test injected")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive published event within 2s")
	}
}

// TestProviderTracker_OverlappingVerifyAllCoalesces asserts the refresh
// mutex collapses concurrent VerifyAll calls so a burst of SSE
// subscriptions cannot stampede remotes with overlapping capability
// probes.
func TestProviderTracker_OverlappingVerifyAllCoalesces(t *testing.T) {
	var hits int64
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&hits, 1)
		// Hold the response just long enough for an overlapping
		// VerifyAll call to start (and be coalesced away).
		time.Sleep(200 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providerName":"slow","providerType":"remote","capabilities":[]}`))
	}))
	defer slow.Close()

	rp := newTestRemoteProvider(t, slow.URL)
	rp.Initialize()
	providers := map[string]Provider{"slow": rp}
	tracker := NewProviderTracker(providers, rp.Log)

	// Fire a primary VerifyAll and, while it's in flight, fire a second
	// one that should be coalesced away by the refresh mutex.
	go tracker.VerifyAll(context.Background())
	time.Sleep(50 * time.Millisecond)
	tracker.VerifyAll(context.Background())

	// Give the primary plenty of time to finish.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt64(&hits) >= 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}

	// The coalesced call must have returned without issuing a probe. The
	// primary's probe may complete at any time; with the version fallback
	// path enabled it will issue one request for the versioned path
	// (success on first try) plus zero for the unversioned (skipped on
	// success), so we expect exactly one hit.
	if got := atomic.LoadInt64(&hits); got != 1 {
		t.Errorf("expected exactly 1 capability probe hit, got %d (coalesce broken)", got)
	}
}

// TestRemoteProvider_VersionedCapabilitiesFallback asserts the boot-time
// /capabilities probe falls back from /<version>/capabilities to plain
// /capabilities when the versioned path returns a non-2xx. This protects
// older remotes that do not version their capability manifest.
func TestRemoteProvider_VersionedCapabilitiesFallback(t *testing.T) {
	var versionedHit, unversionedHit int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/capabilities":
			atomic.AddInt64(&unversionedHit, 1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"providerName":"legacy","providerType":"remote","capabilities":[]}`))
		default:
			atomic.AddInt64(&versionedHit, 1)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	rp := newTestRemoteProvider(t, server.URL)
	rp.Initialize()

	props, err := rp.VerifyAvailability(context.Background())
	if err != nil {
		t.Fatalf("VerifyAvailability err = %v", err)
	}
	if props.ProviderName != "legacy" {
		t.Errorf("ProviderName = %q, want %q", props.ProviderName, "legacy")
	}
	if v, u := atomic.LoadInt64(&versionedHit), atomic.LoadInt64(&unversionedHit); v == 0 || u == 0 {
		t.Errorf("expected at least one hit on each path; versioned=%d unversioned=%d", v, u)
	}
}

// TestRemoteProvider_VersionedCapabilitiesNoFallbackWhenVersionedSucceeds
// asserts the unversioned fallback is NOT attempted when the versioned
// path returns 2xx. Keeps boot probe traffic bounded for modern remotes.
func TestRemoteProvider_VersionedCapabilitiesNoFallbackWhenVersionedSucceeds(t *testing.T) {
	var unversionedHit int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/capabilities" {
			atomic.AddInt64(&unversionedHit, 1)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providerName":"modern","providerType":"remote","capabilities":[]}`))
	}))
	defer server.Close()

	rp := newTestRemoteProvider(t, server.URL)
	rp.Initialize()

	if _, err := rp.VerifyAvailability(context.Background()); err != nil {
		t.Fatalf("VerifyAvailability err = %v", err)
	}
	if u := atomic.LoadInt64(&unversionedHit); u != 0 {
		t.Errorf("unversioned fallback was attempted %d times despite versioned success", u)
	}
}

// TestRemoteProvider_InitializeStampsFallbackName asserts Initialize sets
// a recognizable fallback ProviderName (the URL host) without performing a
// network call, so the chooser can render the entry immediately and never
// has to display a raw URL.
func TestRemoteProvider_InitializeStampsFallbackName(t *testing.T) {
	rp := newTestRemoteProvider(t, "https://cloud.example.org:8443")
	rp.Initialize()

	if rp.GetProviderType() != RemoteProviderType {
		t.Errorf("ProviderType = %q, want %q", rp.GetProviderType(), RemoteProviderType)
	}
	if rp.GetProviderURL() != "https://cloud.example.org:8443" {
		t.Errorf("ProviderURL = %q, want %q", rp.GetProviderURL(), "https://cloud.example.org:8443")
	}
	if rp.Name() != "cloud.example.org:8443" {
		t.Errorf("Name() = %q, want %q (URL host fallback)", rp.Name(), "cloud.example.org:8443")
	}
}

func newTestLocalProviderForTracker(t *testing.T) *DefaultLocalProvider {
	t.Helper()
	lp := &DefaultLocalProvider{Log: newTestLogger(t)}
	lp.Initialize()
	return lp
}

// TestProviderTracker_SnapshotDeliveredBeforeLiveUpdates asserts the
// snapshot is always delivered before live updates, in full: a subscriber
// sees one event per registered provider (the snapshot taken at Subscribe
// time) before any Publish that happens after Subscribe returns, even when
// the snapshot is larger than the old 32-event channel buffer.
func TestProviderTracker_SnapshotDeliveredBeforeLiveUpdates(t *testing.T) {
	log := newTestLogger(t)
	providers := map[string]Provider{
		LocalProviderName: newTestLocalProviderForTracker(t),
	}
	tracker := NewProviderTracker(providers, log)

	// 40 snapshot entries: more than the 32-event buffer the old replay
	// goroutine used, so a drop-prone replay would lose events here.
	const extra = 40
	tracker.mu.Lock()
	for i := 0; i < extra; i++ {
		key := "provider-" + string(rune('a'+i))
		tracker.statuses[key] = ProviderStatusEvent{
			Key:    key,
			Status: ProviderStatusOnline,
		}
	}
	tracker.mu.Unlock()
	total := extra + 1 // plus the local provider

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events, unsubscribe := tracker.Subscribe(ctx)
	defer unsubscribe()

	// A live update for an existing key, published after Subscribe
	// returns, must arrive after the whole snapshot, never before it.
	tracker.Publish(ProviderStatusEvent{
		Key:    "provider-a",
		Status: ProviderStatusOffline,
		Error:  "test injected",
	})

	// Read the snapshot plus the one live event.
	got := make([]ProviderStatusEvent, 0, total+1)
	timeout := time.After(5 * time.Second)
	for len(got) < total+1 {
		select {
		case evt, ok := <-events:
			if !ok {
				t.Fatalf("channel closed after %d events, want %d", len(got), total+1)
			}
			got = append(got, evt)
		case <-timeout:
			t.Fatalf("received %d events, want %d (snapshot incomplete or live update lost)", len(got), total+1)
		}
	}

	// The live event must be last; everything before it is the snapshot.
	live := got[len(got)-1]
	if live.Key != "provider-a" || live.Status != ProviderStatusOffline {
		t.Fatalf("last event = %+v, want live update for provider-a Offline", live)
	}
	seen := make(map[string]ProviderStatusEvent, total)
	for _, evt := range got[:len(got)-1] {
		if _, dup := seen[evt.Key]; dup {
			t.Fatalf("duplicate snapshot event for key %q", evt.Key)
		}
		if evt.Key == "provider-a" && evt.Status != ProviderStatusOnline {
			t.Fatalf("snapshot event for provider-a has status %q, want Online (live update overtook snapshot)", evt.Status)
		}
		seen[evt.Key] = evt
	}
	if len(seen) != total {
		t.Fatalf("snapshot delivered %d distinct keys, want %d", len(seen), total)
	}
}

// TestProviderTracker_ConcurrentUnsubscribeAndPublishDoesNotPanic asserts that
// concurrent Publish calls and unsubscribe calls do not race into a send on a
// closed channel.
func TestProviderTracker_ConcurrentUnsubscribeAndPublishDoesNotPanic(t *testing.T) {
	log := newTestLogger(t)
	providers := map[string]Provider{
		LocalProviderName: newTestLocalProviderForTracker(t),
	}
	tracker := NewProviderTracker(providers, log)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 500; i++ {
			tracker.Publish(ProviderStatusEvent{
				Key:    "churn",
				Status: ProviderStatusOnline,
			})
			time.Sleep(10 * time.Microsecond)
		}
		close(done)
	}()

	for i := 0; i < 200; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		ch, unsubscribe := tracker.Subscribe(ctx)

		go func() {
			time.Sleep(5 * time.Microsecond)
			unsubscribe()
			cancel()
		}()

		for range ch {
		}
	}

	<-done
}

func mapKeys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
