package models

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// syncPrefsTestProvider returns a RemoteProvider whose preference-sync PUTs land
// on a local test server, plus a channel that reports each one that arrives.
// Capabilities start empty, which is the state Initialize leaves behind: the
// /capabilities fetch happens later, in VerifyAvailability.
func syncPrefsTestProvider(t *testing.T) (*RemoteProvider, <-chan struct{}) {
	t.Helper()

	synced := make(chan struct{}, 32)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusCreated)
		select {
		case synced <- struct{}{}:
		default:
		}
	}))
	t.Cleanup(srv.Close)

	l := &RemoteProvider{
		RemoteProviderURL: srv.URL,
		Log:               newTestLogger(t),
	}
	l.Initialize()
	t.Cleanup(l.StopSyncPreferences)

	return l, synced
}

// syncPrefsCapabilities is the capability set a reachable remote reports, as
// VerifyAvailability would hand it to SetProviderProperties.
func syncPrefsCapabilities(providerURL string) ProviderProperties {
	return ProviderProperties{
		ProviderType: RemoteProviderType,
		ProviderName: "Meshery",
		ProviderURL:  providerURL,
		Capabilities: Capabilities{{Feature: SyncPrefs, Endpoint: "/api/user/preferences"}},
	}
}

// syncWorkerRunning reads the worker flag under its own lock so these
// assertions stay clean under -race.
func syncWorkerRunning(l *RemoteProvider) bool {
	l.syncMu.Lock()
	defer l.syncMu.Unlock()
	return l.syncActive
}

func syncQueue(l *RemoteProvider) chan *userSession {
	l.syncMu.Lock()
	defer l.syncMu.Unlock()
	return l.syncChan
}

// enqueueWithin runs the enqueue the request path performs and fails the test
// if it does not return within d. Blocking here is the defect being guarded
// against: it parks the goroutine serving the user's request permanently.
func enqueueWithin(t *testing.T, l *RemoteProvider, d time.Duration) error {
	t.Helper()

	errCh := make(chan error, 1)
	go func() {
		errCh <- l.enqueuePrefSync("test-token", NewDefaultPreference())
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(d):
		t.Fatalf("enqueuePrefSync blocked for %s; the goroutine serving the request would be wedged", d)
		return nil
	}
}

// TestSyncPreferencesActivatesWhenCapabilitiesLoadAfterBoot is the regression
// test for #21260.
//
// A remote that is unreachable when the boot probe runs has no capabilities, so
// the boot-time SyncPreferences call cannot start the worker. It used to be the
// only call site, which left the queue nil for the rest of the process and hung
// every preference write on a nil-channel send. Capabilities arriving later must
// start the worker instead.
func TestSyncPreferencesActivatesWhenCapabilitiesLoadAfterBoot(t *testing.T) {
	l, synced := syncPrefsTestProvider(t)

	// Boot: the probe failed, so main.go's activation finds no capabilities.
	l.SyncPreferences()
	if syncWorkerRunning(l) {
		t.Fatal("worker started without SyncPrefs in the capability set")
	}

	// A write in this window must fail fast rather than hang.
	if err := enqueueWithin(t, l, 2*time.Second); err == nil {
		t.Error("expected an error while no worker is running, got nil")
	}

	// The remote comes back and a capability set lands. This is exactly what
	// VerifyAvailability does after a successful /capabilities fetch, whether
	// it runs from the boot probe or from an /api/providers/stream refresh.
	l.SetProviderProperties(syncPrefsCapabilities(l.RemoteProviderURL))

	if !syncWorkerRunning(l) {
		t.Fatal("capabilities loaded but the sync worker was not started")
	}

	if err := enqueueWithin(t, l, 2*time.Second); err != nil {
		t.Fatalf("enqueue failed after activation: %v", err)
	}

	select {
	case <-synced:
	case <-time.After(10 * time.Second):
		t.Fatal("the preference never reached the remote provider")
	}
}

// TestSyncPreferencesIsIdempotent covers activation now happening from two
// places: main.go after the boot probe, and every SetProviderProperties. A
// second activation must not replace the queue, which would abandon whatever is
// still buffered in it and leave two workers competing for the same writes.
func TestSyncPreferencesIsIdempotent(t *testing.T) {
	l, _ := syncPrefsTestProvider(t)
	l.SetProviderProperties(syncPrefsCapabilities(l.RemoteProviderURL))

	first := syncQueue(l)
	if first == nil {
		t.Fatal("expected a queue after activation")
	}

	for i := 0; i < 5; i++ {
		l.SyncPreferences()
	}

	if syncQueue(l) != first {
		t.Error("SyncPreferences replaced the queue; a duplicate worker is running")
	}
}

// TestEnqueuePrefSyncDoesNotBlockOnFullQueue pins the second way the send could
// park a request goroutine. The preference is already persisted locally by this
// point, so a saturated best-effort queue is logged and accepted rather than
// failing the user's request.
//
// The state is built directly because a queue that is genuinely full while a
// worker drains it cannot be produced deterministically.
func TestEnqueuePrefSyncDoesNotBlockOnFullQueue(t *testing.T) {
	l := &RemoteProvider{Log: newTestLogger(t)}

	l.syncMu.Lock()
	l.syncChan = make(chan *userSession, 1)
	l.syncChan <- &userSession{session: NewDefaultPreference()}
	l.syncActive = true
	l.syncMu.Unlock()

	if err := enqueueWithin(t, l, 2*time.Second); err != nil {
		t.Errorf("a full queue should be accepted and logged, got %v", err)
	}
}

// TestStopSyncPreferencesIsSafeWhenNeverStarted covers main.go, which defers
// StopSyncPreferences for every remote whether or not activation happened.
func TestStopSyncPreferencesIsSafeWhenNeverStarted(t *testing.T) {
	l, _ := syncPrefsTestProvider(t)

	stopWithin(t, l, 2*time.Second)
}

// TestStopSyncPreferencesIsIdempotent guards the shutdown path against a second
// call. Signalling the worker used to be a send on an unbuffered channel, which
// blocks once the worker has returned.
func TestStopSyncPreferencesIsIdempotent(t *testing.T) {
	l, _ := syncPrefsTestProvider(t)
	l.SetProviderProperties(syncPrefsCapabilities(l.RemoteProviderURL))

	stopWithin(t, l, 2*time.Second)
	if syncWorkerRunning(l) {
		t.Error("worker still marked running after Stop")
	}
	stopWithin(t, l, 2*time.Second)
}

func stopWithin(t *testing.T, l *RemoteProvider, d time.Duration) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		defer close(done)
		l.StopSyncPreferences()
	}()

	select {
	case <-done:
	case <-time.After(d):
		t.Fatalf("StopSyncPreferences blocked for %s", d)
	}
}

// TestSyncPrefsLifecycleConcurrent exercises the three entry points together.
// Activation is reachable from every probe goroutine and from main.go, enqueue
// from every request goroutine, and Stop from the shutdown defer. Run with
// -race.
func TestSyncPrefsLifecycleConcurrent(t *testing.T) {
	l, _ := syncPrefsTestProvider(t)
	l.SetProviderProperties(syncPrefsCapabilities(l.RemoteProviderURL))

	const workers = 8
	const iterations = 40

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				switch (id + j) % 4 {
				case 0:
					l.SyncPreferences()
				case 1, 2:
					_ = l.enqueuePrefSync("test-token", NewDefaultPreference())
				case 3:
					l.StopSyncPreferences()
				}
			}
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		t.Fatal("preference-sync lifecycle deadlocked under concurrent access")
	}
}
