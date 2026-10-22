package utils

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshkit/logger"
)

// mockLogger implements logger.Handler for tests without external dependencies.
type mockLogger struct {
	logger.Handler
	lastErr atomic.Pointer[error]
}

func (m *mockLogger) Error(err error) {
	m.lastErr.Store(&err)
}

func TestSafeGo_NormalExecution(t *testing.T) {
	log := &mockLogger{}
	var executed atomic.Bool

	var wg sync.WaitGroup
	wg.Add(1)

	SafeGo(log, "test-normal", func() {
		defer wg.Done()
		executed.Store(true)
	})

	wg.Wait()

	if !executed.Load() {
		t.Fatal("expected function to execute")
	}
	if log.lastErr.Load() != nil {
		t.Fatalf("unexpected error logged: %v", *log.lastErr.Load())
	}
}

func TestSafeGo_RecoversPanicWithoutCrash(t *testing.T) {
	log := &mockLogger{}
	var wg sync.WaitGroup
	wg.Add(1)

	var panicCaught atomic.Bool

	SafeGo(log, "test-panic", func() {
		panic("simulated fatal nil dereference in background routine")
	}, WithOnPanic(func(r any, stack []byte) {
		defer wg.Done()
		if r == "simulated fatal nil dereference in background routine" {
			panicCaught.Store(true)
		}
		if len(stack) == 0 {
			t.Error("expected non-empty stack trace")
		}
	}))

	wg.Wait()

	if !panicCaught.Load() {
		t.Fatal("expected panic callback to be invoked with the panic value")
	}
	if log.lastErr.Load() == nil {
		t.Fatal("expected error to be logged by mockLogger on panic")
	}
}

func TestSafeGo_RecoversPanicInOnPanicCallback(t *testing.T) {
	log := &mockLogger{}
	var wg sync.WaitGroup
	wg.Add(1)

	SafeGo(log, "test-panic-in-callback", func() {
		panic("initial worker panic")
	}, WithOnPanic(func(r any, stack []byte) {
		defer wg.Done()
		panic("panic inside onPanic callback")
	}))

	wg.Wait()

	if log.lastErr.Load() == nil {
		t.Fatal("expected error to be logged when onPanic callback panics")
	}
}

func TestSafeGo_RespectsCancelledContext(t *testing.T) {
	log := &mockLogger{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	var executed atomic.Bool
	done := make(chan struct{})

	SafeGo(log, "test-cancelled-context", func() {
		executed.Store(true)
		close(done)
	}, WithContext(ctx))

	select {
	case <-done:
		t.Fatal("fn should not execute when context is already cancelled")
	case <-time.After(50 * time.Millisecond):
		// Success: routine terminated early
	}

	if executed.Load() {
		t.Fatal("fn was executed despite cancelled context")
	}
}

func TestRunSafe_NormalExecution(t *testing.T) {
	log := &mockLogger{}
	expectedErr := errors.New("business logic error")

	err := RunSafe(log, "test-runsafe-err", func() error {
		return expectedErr
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if log.lastErr.Load() != nil {
		t.Fatal("unexpected log on normal error return")
	}
}

func TestRunSafe_RecoversPanic(t *testing.T) {
	log := &mockLogger{}

	err := RunSafe(log, "test-runsafe-panic", func() error {
		panic("runtime panic")
	})

	if err == nil {
		t.Fatal("expected RunSafe to return an error when panic is recovered")
	}
	if log.lastErr.Load() == nil {
		t.Fatal("expected panic to be logged via logger")
	}
}
