package utils

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/meshery/meshkit/logger"
)

// RecoveredPanic describes a panic that Guard or SafeGo contained.
type RecoveredPanic struct {
	// Name identifies the call site that was running, for example
	// "controlplane-subscription". It is a short, low-cardinality label.
	Name string
	// Value is what was passed to panic.
	Value any
	// Stack is the goroutine stack at the point of the panic. It is meant for
	// server-side logs; it must never be returned to an API client.
	Stack []byte
}

// PanicHook observes every contained panic, after it has been logged. A hook
// is how a recovered panic becomes visible beyond a log line: a system event,
// a metric, an alert. Hooks run on the goroutine that recovered, so they must
// be quick and must not block. A hook that panics is itself contained.
type PanicHook func(RecoveredPanic)

var (
	hooksMu sync.RWMutex
	hooks   []PanicHook

	countsMu sync.Mutex
	counts   = map[string]int64{}
)

// OnRecoveredPanic registers a hook that is called for every panic contained by
// Guard, SafeGo or ReportRecoveredPanic. Register hooks once at start-up.
func OnRecoveredPanic(h PanicHook) {
	if h == nil {
		return
	}
	hooksMu.Lock()
	defer hooksMu.Unlock()
	hooks = append(hooks, h)
}

// RecoveredPanicCounts returns how many panics have been contained so far, per
// call-site name. It is a snapshot; callers may keep or modify the result.
func RecoveredPanicCounts() map[string]int64 {
	countsMu.Lock()
	defer countsMu.Unlock()
	out := make(map[string]int64, len(counts))
	for name, n := range counts {
		out[name] = n
	}
	return out
}

// Guard runs fn and contains a panic raised by it, so one bad input costs one
// logged and reported failure instead of the whole process. It reports whether
// fn panicked.
//
// Use Guard around the body of ONE iteration of a long-lived loop that handles
// external input (a cluster object, a broker event), so that the loop survives
// a malformed item and the next item is still handled:
//
//	for ev := range events {
//		utils.Guard("controlplane-subscription", log, func() { handle(ev) })
//	}
//
// Do not wrap a whole process-level goroutine in Guard: a goroutine that must
// not keep running after a defect should still crash loudly.
func Guard(name string, log logger.Handler, fn func()) (panicked bool) {
	defer func() {
		if recovered := recover(); recovered != nil {
			panicked = true
			ReportRecoveredPanic(name, recovered, log)
		}
	}()
	fn()
	return false
}

// SafeGo runs fn in a new goroutine under Guard. Use it for a request-scoped
// goroutine, one started on behalf of a single request or operation, whose
// failure must not take the server down. A bare `go` statement is still right
// for process-level goroutines.
func SafeGo(name string, log logger.Handler, fn func()) {
	go Guard(name, log, fn)
}

// ReportRecoveredPanic logs a panic that the caller has just recovered, counts
// it and runs the registered hooks. It must be called from the deferred
// function that called recover(), so that the stack still shows the panic.
// Guard calls it; the HTTP recovery middleware calls it too, so both feed the
// same counters and hooks. The stack goes to the log and the hooks only.
//
// log may be nil, in which case nothing is logged but the count and hooks still
// run.
func ReportRecoveredPanic(name string, recovered any, log logger.Handler) {
	p := RecoveredPanic{Name: name, Value: recovered, Stack: debug.Stack()}

	countsMu.Lock()
	counts[name]++
	countsMu.Unlock()

	if log != nil {
		log.Error(fmt.Errorf("recovered panic in %s: %v\n%s", name, recovered, p.Stack))
	}

	hooksMu.RLock()
	snapshot := append([]PanicHook(nil), hooks...)
	hooksMu.RUnlock()
	for _, h := range snapshot {
		runHook(h, p, log)
	}
}

// runHook calls one hook and contains a panic from it, so a faulty observer
// cannot undo the containment it is observing. The panic is counted under
// "hook-panic:<site>" even without a logger, so a faulty hook stays visible
// instead of failing silently forever.
func runHook(h PanicHook, p RecoveredPanic, log logger.Handler) {
	defer func() {
		if recovered := recover(); recovered != nil {
			countsMu.Lock()
			counts["hook-panic:"+p.Name]++
			countsMu.Unlock()
			if log != nil {
				log.Error(fmt.Errorf("panic hook for %s panicked: %v", p.Name, recovered))
			}
		}
	}()
	h(p)
}
