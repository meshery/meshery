package handlers

import (
	"sync"
	"testing"

	corepolicies "github.com/meshery/meshkit/models/meshmodel/core/policies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetRego_ConcurrentWithRegoEngine pins the publication contract of the
// OPA/Rego engine handoff: SetRego writers and regoEngine readers may run
// concurrently, so the field must be an atomic pointer - a plain pointer or
// struct field fails this test under -race.
func TestSetRego_ConcurrentWithRegoEngine(t *testing.T) {
	h := &Handler{}

	first, second := &corepolicies.Rego{}, &corepolicies.Rego{}

	var writerWG, readerWG sync.WaitGroup

	writerWG.Add(1)
	go func() {
		defer writerWG.Done()
		for i := 0; i < 1000; i++ {
			if i%2 == 0 {
				h.SetRego(first)
			} else {
				h.SetRego(second)
			}
		}
	}()

	readerWG.Add(1)
	go func() {
		defer readerWG.Done()
		for i := 0; i < 1000; i++ {
			r, err := h.regoEngine()
			if err == nil {
				// While the writer runs, the published engine is one of
				// the two instances, never a mix of them.
				require.True(t, r == first || r == second)
			}
		}
	}()

	writerWG.Wait()
	readerWG.Wait()

	r, err := h.regoEngine()
	require.NoError(t, err)
	require.True(t, r == first || r == second)
}

// TestPublishBootRego_DoesNotOverwriteReloadedEngine pins the stale-overwrite
// fix: a runtime reload (SetRego) that published before the boot seeding
// goroutine finishes keeps its engine.
func TestPublishBootRego_DoesNotOverwriteReloadedEngine(t *testing.T) {
	h := &Handler{}

	reloaded := &corepolicies.Rego{}
	boot := &corepolicies.Rego{}

	h.SetRego(reloaded)
	h.PublishBootRego(boot)

	r, err := h.regoEngine()
	require.NoError(t, err)
	assert.Same(t, reloaded, r)
}

// TestPublishBootRego_PublishesOnEmptyHandler is the positive boot path: a
// handler with no published engine takes the boot seeding engine.
func TestPublishBootRego_PublishesOnEmptyHandler(t *testing.T) {
	h := &Handler{}

	boot := &corepolicies.Rego{}
	h.PublishBootRego(boot)

	r, err := h.regoEngine()
	require.NoError(t, err)
	assert.Same(t, boot, r)
}
