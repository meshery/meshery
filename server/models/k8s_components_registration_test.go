package models

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/core"
)

// registrationTestProvider is a partial models.Provider implementation: it embeds
// the interface (so the type satisfies models.Provider) and only overrides
// PersistSystemEvent, the single Provider method RegisterComponents calls before
// spawning each per-context goroutine.
type registrationTestProvider struct {
	Provider
}

func (registrationTestProvider) PersistSystemEvent(events.Event) error { return nil }

func newTestContexts(instanceID *core.Uuid, n int) []*K8sContext {
	ctxs := make([]*K8sContext, 0, n)
	for i := 0; i < n; i++ {
		ctxs = append(ctxs, &K8sContext{
			ID:                uuid.Must(uuid.NewV4()).String(),
			Name:              "ctx",
			MesheryInstanceID: instanceID,
		})
	}
	return ctxs
}

// trackingRegFunc returns a K8sRegistrationFunction that records, via peak/current
// counters, how many invocations are running concurrently, and calls wg.Done() on
// completion.
func trackingRegFunc(wg *sync.WaitGroup, current, peak *int64, mu *sync.Mutex, work time.Duration) K8sRegistrationFunction {
	return func(_ *Provider, _ context.Context, _ []byte, _, _, _ string, _ core.Uuid, _ *meshmodel.RegistryManager, _ *Broadcast, _ logger.Handler, _ string) error {
		defer wg.Done()

		n := atomic.AddInt64(current, 1)
		mu.Lock()
		if n > *peak {
			*peak = n
		}
		mu.Unlock()

		time.Sleep(work)

		atomic.AddInt64(current, -1)
		return nil
	}
}

// TestRegisterComponentsBoundsConcurrency verifies that RegisterComponents never
// runs more than maxConcurrentK8sRegistrations goroutines at once, regardless of
// how many contexts it is given. Before the fix, RegisterComponents spawned one
// unbounded goroutine per context; a kubeconfig with many contexts (which need not
// point at real, reachable clusters) could drive an unbounded goroutine/FD/dial
// spike. See issue #21610.
func TestRegisterComponentsBoundsConcurrency(t *testing.T) {
	const numContexts = 50

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctxs := newTestContexts(&instanceID, numContexts)

	helper := NewComponentsRegistrationHelper(newTestLogger(t))
	helper.UpdateContexts(ctxs)

	var (
		current, peak int64
		mu            sync.Mutex
		wg            sync.WaitGroup
	)
	wg.Add(numContexts)

	regFunc := trackingRegFunc(&wg, &current, &peak, &mu, 20*time.Millisecond)
	provider := registrationTestProvider{}
	broadcaster := NewBroadcaster("test")

	helper.RegisterComponents(ctxs, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)

	waitOrFail(t, &wg, 10*time.Second)

	if peak > maxConcurrentK8sRegistrations {
		t.Fatalf("RegisterComponents ran %d contexts concurrently, want at most %d", peak, maxConcurrentK8sRegistrations)
	}
	if peak == 0 {
		t.Fatal("no registrations observed running concurrently; test did not exercise RegisterComponents")
	}
}

// TestRegisterComponentsBoundsConcurrencyAcrossCalls verifies that the concurrency
// cap is shared across separate RegisterComponents calls on the same helper, not
// reset per call. The production helper is a single long-lived instance reused
// across concurrent HTTP requests (see handlers.K8sCompRegHelper); a per-call
// semaphore would let each concurrent request get its own budget, multiplying the
// effective concurrency by the number of in-flight requests.
func TestRegisterComponentsBoundsConcurrencyAcrossCalls(t *testing.T) {
	const contextsPerCall = 30

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctxsA := newTestContexts(&instanceID, contextsPerCall)
	ctxsB := newTestContexts(&instanceID, contextsPerCall)

	helper := NewComponentsRegistrationHelper(newTestLogger(t))
	helper.UpdateContexts(ctxsA)
	helper.UpdateContexts(ctxsB)

	var (
		current, peak int64
		mu            sync.Mutex
		wg            sync.WaitGroup
	)
	wg.Add(2 * contextsPerCall)

	regFunc := trackingRegFunc(&wg, &current, &peak, &mu, 20*time.Millisecond)
	provider := registrationTestProvider{}
	broadcaster := NewBroadcaster("test")

	// Simulate two concurrent HTTP requests both driving registration through the
	// same shared helper.
	go helper.RegisterComponents(ctxsA, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)
	go helper.RegisterComponents(ctxsB, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)

	waitOrFail(t, &wg, 10*time.Second)

	if peak > maxConcurrentK8sRegistrations {
		t.Fatalf("two concurrent RegisterComponents calls together ran %d contexts concurrently, want at most %d", peak, maxConcurrentK8sRegistrations)
	}
}

// TestRegisterComponentsDoesNotBlockCaller verifies that RegisterComponents returns
// promptly even when the shared semaphore is already saturated by a prior call.
// RegisterComponents is invoked synchronously from the HTTP handler right before it
// writes a 202 Accepted response; if acquiring the semaphore happened on the
// caller's goroutine, a full semaphore would stall that response.
func TestRegisterComponentsDoesNotBlockCaller(t *testing.T) {
	const numContexts = maxConcurrentK8sRegistrations + 5

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctxs := newTestContexts(&instanceID, numContexts)

	helper := NewComponentsRegistrationHelper(newTestLogger(t))
	helper.UpdateContexts(ctxs)

	var wg sync.WaitGroup
	wg.Add(numContexts)

	// Slow enough that, if RegisterComponents blocked the caller while the
	// semaphore is full, the call below would still be blocked when we check.
	regFunc := K8sRegistrationFunction(func(_ *Provider, _ context.Context, _ []byte, _, _, _ string, _ core.Uuid, _ *meshmodel.RegistryManager, _ *Broadcast, _ logger.Handler, _ string) error {
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	provider := registrationTestProvider{}
	broadcaster := NewBroadcaster("test")

	callReturned := make(chan struct{})
	go func() {
		helper.RegisterComponents(ctxs, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)
		close(callReturned)
	}()

	select {
	case <-callReturned:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("RegisterComponents did not return promptly; it appears to block the caller once the semaphore is saturated")
	}

	waitOrFail(t, &wg, 10*time.Second)
}

// TestRegisterComponentsDoesNotSpawnGoroutinePerContext verifies that
// RegisterComponents does not create a goroutine for every admitted context up
// front. Before this fix, the shared semaphore (cg.regSem) was acquired inside
// each per-context goroutine, so a call with many contexts still spawned one
// goroutine per context immediately -- all but maxConcurrentK8sRegistrations of
// them just parked blocked on the semaphore acquire. A kubeconfig with a large
// number of contexts (which need not point at real, reachable clusters) could
// still drive an unbounded goroutine/memory/scheduler spike before any
// registration work started. See PR #21611 review feedback.
func TestRegisterComponentsDoesNotSpawnGoroutinePerContext(t *testing.T) {
	const numContexts = 500

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))
	ctxs := newTestContexts(&instanceID, numContexts)

	helper := NewComponentsRegistrationHelper(newTestLogger(t))
	helper.UpdateContexts(ctxs)

	var wg sync.WaitGroup
	wg.Add(numContexts)

	// Block every worker so admitted-but-not-yet-running goroutines stay alive
	// long enough to be observed by runtime.NumGoroutine().
	release := make(chan struct{})
	regFunc := K8sRegistrationFunction(func(_ *Provider, _ context.Context, _ []byte, _, _, _ string, _ core.Uuid, _ *meshmodel.RegistryManager, _ *Broadcast, _ logger.Handler, _ string) error {
		defer wg.Done()
		<-release
		return nil
	})

	provider := registrationTestProvider{}
	broadcaster := NewBroadcaster("test")

	before := runtime.NumGoroutine()

	helper.RegisterComponents(ctxs, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)

	// Give any (incorrectly) eagerly-spawned goroutines time to start and park
	// on the semaphore acquire.
	time.Sleep(200 * time.Millisecond)

	after := runtime.NumGoroutine()
	spawned := after - before

	close(release)
	waitOrFail(t, &wg, 10*time.Second)

	// Allow generous headroom above the concurrency cap for the dispatcher
	// goroutine itself and any incidental runtime goroutines, but reject
	// anything on the order of numContexts, which is what an eager per-context
	// spawn would produce.
	const maxAllowedGoroutines = maxConcurrentK8sRegistrations + 10
	if spawned > maxAllowedGoroutines {
		t.Fatalf("RegisterComponents spawned %d goroutines for %d contexts, want at most %d", spawned, numContexts, maxAllowedGoroutines)
	}
}

func waitOrFail(t *testing.T, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatal("timed out waiting for registrations to complete")
	}
}
