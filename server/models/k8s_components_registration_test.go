package models

import (
	"context"
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

// TestRegisterComponentsBoundsConcurrency verifies that RegisterComponents never
// runs more than maxConcurrentK8sRegistrations goroutines at once, regardless of
// how many contexts it is given. Before the fix, RegisterComponents spawned one
// unbounded goroutine per context; a kubeconfig with many contexts (which need not
// point at real, reachable clusters) could drive an unbounded goroutine/FD/dial
// spike. See issue #21610.
func TestRegisterComponentsBoundsConcurrency(t *testing.T) {
	const numContexts = 50

	instanceID := core.Uuid(uuid.Must(uuid.NewV4()))

	ctxs := make([]*K8sContext, 0, numContexts)
	for i := 0; i < numContexts; i++ {
		ctxs = append(ctxs, &K8sContext{
			ID:                uuid.Must(uuid.NewV4()).String(),
			Name:              "ctx",
			MesheryInstanceID: &instanceID,
		})
	}

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	helper := NewComponentsRegistrationHelper(log)
	helper.UpdateContexts(ctxs)

	var (
		current int64
		peak    int64
		mu      sync.Mutex
		wg      sync.WaitGroup
	)
	wg.Add(numContexts)

	regFunc := K8sRegistrationFunction(func(_ *Provider, _ context.Context, _ []byte, _, _, _ string, _ core.Uuid, _ *meshmodel.RegistryManager, _ *Broadcast, _ logger.Handler, _ string) error {
		defer wg.Done()

		n := atomic.AddInt64(&current, 1)
		mu.Lock()
		if n > peak {
			peak = n
		}
		mu.Unlock()

		time.Sleep(20 * time.Millisecond)

		atomic.AddInt64(&current, -1)
		return nil
	})

	provider := registrationTestProvider{}
	broadcaster := NewBroadcaster("test")

	helper.RegisterComponents(ctxs, []K8sRegistrationFunction{regFunc}, nil, broadcaster, provider, uuid.Must(uuid.NewV4()).String(), false)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("timed out waiting for registrations to complete")
	}

	if peak > maxConcurrentK8sRegistrations {
		t.Fatalf("RegisterComponents ran %d contexts concurrently, want at most %d", peak, maxConcurrentK8sRegistrations)
	}
	if peak == 0 {
		t.Fatal("no registrations observed running concurrently; test did not exercise RegisterComponents")
	}
}
