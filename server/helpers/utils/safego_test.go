package utils

import (
	"bytes"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
)

// lockedBuffer is a bytes.Buffer that is safe to write from the goroutines the
// tests start and read from the test goroutine under -race.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// resetPanicState gives a test a clean set of hooks and counters and restores
// the previous state afterwards, so tests do not leak into one another.
func resetPanicState(t *testing.T) {
	t.Helper()
	hooksMu.Lock()
	savedHooks := hooks
	hooks = nil
	hooksMu.Unlock()
	countsMu.Lock()
	savedCounts := counts
	counts = map[string]int64{}
	countsMu.Unlock()

	t.Cleanup(func() {
		hooksMu.Lock()
		hooks = savedHooks
		hooksMu.Unlock()
		countsMu.Lock()
		counts = savedCounts
		countsMu.Unlock()
	})
}

func newCapturingLogger(t *testing.T) (logger.Handler, *lockedBuffer) {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	// A zero-value logger sits above Error level and would drop what the test
	// wants to see.
	log.SetLevel(logrus.ErrorLevel)
	out := &lockedBuffer{}
	log.UpdateErrorLogOutput(out)
	return log, out
}

func TestGuard_NoPanic(t *testing.T) {
	resetPanicState(t)
	log, out := newCapturingLogger(t)

	hookCalls := 0
	OnRecoveredPanic(func(RecoveredPanic) { hookCalls++ })

	ran := false
	if Guard("no-panic", log, func() { ran = true }) {
		t.Fatal("Guard reported a panic for a function that did not panic")
	}
	if !ran {
		t.Error("the guarded function did not run")
	}
	if hookCalls != 0 || len(RecoveredPanicCounts()) != 0 || out.String() != "" {
		t.Errorf("nothing should be reported without a panic: hooks=%d counts=%v log=%q",
			hookCalls, RecoveredPanicCounts(), out.String())
	}
}

// setKey writes into m. Tests call it with a nil map to trigger the runtime
// panic; going through a parameter keeps staticcheck (SA5000) from rejecting a
// write it can see is aimed at a nil map.
func setKey(m map[string]string, k, v string) { m[k] = v }

// The panics Meshery has actually hit are runtime errors from unchecked input:
// an out-of-range index (#20434) and a write to a nil map. Both must be contained.
func TestGuard_ContainsRealRuntimePanics(t *testing.T) {
	cases := map[string]func(){
		"index out of range": func() {
			parts := []string{"only-one"}
			_ = parts[1]
		},
		"nil map write": func() { setKey(nil, "boom", "x") },
		"nil pointer dereference": func() {
			var p *struct{ n int }
			_ = p.n
		},
		"explicit panic": func() { panic("bad input") },
	}

	for name, fn := range cases {
		t.Run(name, func(t *testing.T) {
			resetPanicState(t)
			log, out := newCapturingLogger(t)

			var got RecoveredPanic
			OnRecoveredPanic(func(p RecoveredPanic) { got = p })

			if !Guard("site-"+name, log, fn) {
				t.Fatal("Guard did not report the panic")
			}
			if got.Name != "site-"+name {
				t.Errorf("hook name = %q, want %q", got.Name, "site-"+name)
			}
			if got.Value == nil {
				t.Error("hook received no panic value")
			}
			if !strings.Contains(string(got.Stack), "goroutine") {
				t.Errorf("hook stack does not look like a stack: %q", got.Stack)
			}
			logged := out.String()
			if !strings.Contains(logged, "recovered panic in site-"+name) || !strings.Contains(logged, "goroutine") {
				t.Errorf("log lacks the name or the stack: %q", logged)
			}
			if n := RecoveredPanicCounts()["site-"+name]; n != 1 {
				t.Errorf("count = %d, want 1", n)
			}
		})
	}
}

// The reason Guard exists: a loop that handles external input survives a
// malformed item, and the items after it are still handled.
func TestGuard_LoopSurvivesPoisonItem(t *testing.T) {
	resetPanicState(t)
	log, _ := newCapturingLogger(t)

	var handled []string
	for _, ev := range []string{"a", "poison", "c"} {
		Guard("loop", log, func() {
			if ev == "poison" {
				setKey(nil, ev, ev)
			}
			handled = append(handled, ev)
		})
	}

	if got := strings.Join(handled, ","); got != "a,c" {
		t.Errorf("handled = %q, want %q: the loop must skip only the poison item", got, "a,c")
	}
	if n := RecoveredPanicCounts()["loop"]; n != 1 {
		t.Errorf("count = %d, want 1", n)
	}
}

func TestGuard_NilLoggerStillReports(t *testing.T) {
	resetPanicState(t)

	called := false
	OnRecoveredPanic(func(RecoveredPanic) { called = true })

	if !Guard("nil-logger", nil, func() { panic("x") }) {
		t.Fatal("Guard did not report the panic")
	}
	if !called {
		t.Error("the hook must still run when there is no logger")
	}
}

func TestGuard_HookPanicIsContained(t *testing.T) {
	resetPanicState(t)
	log, out := newCapturingLogger(t)

	OnRecoveredPanic(func(RecoveredPanic) { panic("faulty observer") })
	secondRan := false
	OnRecoveredPanic(func(RecoveredPanic) { secondRan = true })

	if !Guard("hook-panic", log, func() { panic("x") }) {
		t.Fatal("Guard did not report the panic")
	}
	if !secondRan {
		t.Error("a hook that panics must not stop the hooks registered after it")
	}
	if !strings.Contains(out.String(), "panic hook for hook-panic panicked") {
		t.Errorf("the hook's own panic should be logged: %q", out.String())
	}
}

func TestSafeGo_ContainsPanicInNewGoroutine(t *testing.T) {
	resetPanicState(t)
	log, _ := newCapturingLogger(t)

	done := make(chan RecoveredPanic, 1)
	OnRecoveredPanic(func(p RecoveredPanic) { done <- p })

	SafeGo("request-scoped", log, func() {
		var s []int
		_ = s[3]
	})

	select {
	case p := <-done:
		if p.Name != "request-scoped" {
			t.Errorf("name = %q, want %q", p.Name, "request-scoped")
		}
		if _, ok := p.Value.(runtime.Error); !ok {
			t.Errorf("value = %#v, want a runtime.Error", p.Value)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("SafeGo did not contain the panic")
	}
}

func TestSafeGo_RunsTheFunction(t *testing.T) {
	resetPanicState(t)

	ran := make(chan struct{})
	SafeGo("ok", nil, func() { close(ran) })

	select {
	case <-ran:
	case <-time.After(5 * time.Second):
		t.Fatal("SafeGo never ran the function")
	}
	if len(RecoveredPanicCounts()) != 0 {
		t.Errorf("no panic occurred, counts = %v", RecoveredPanicCounts())
	}
}

// Hooks, counters and Guard are used from many goroutines at once. This is the
// test that fails under -race if the shared state is not synchronised.
func TestGuard_ConcurrentUseIsRaceFree(t *testing.T) {
	resetPanicState(t)
	log, _ := newCapturingLogger(t)

	const workers = 32
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			OnRecoveredPanic(func(RecoveredPanic) {})
			Guard("concurrent", log, func() { panic("x") })
			_ = RecoveredPanicCounts()
		}()
	}
	wg.Wait()

	if n := RecoveredPanicCounts()["concurrent"]; n != workers {
		t.Errorf("count = %d, want %d", n, workers)
	}
}

func TestRecoveredPanicCounts_ReturnsACopy(t *testing.T) {
	resetPanicState(t)

	Guard("copy", nil, func() { panic("x") })
	snapshot := RecoveredPanicCounts()
	snapshot["copy"] = 999

	if n := RecoveredPanicCounts()["copy"]; n != 1 {
		t.Errorf("mutating the snapshot changed the real count to %d", n)
	}
}
