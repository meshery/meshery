package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// HasMachineContext gates every site that drives a state machine returned by
// InitializeMachineWithContext. It has to reject both shapes of a failed
// initialization: the nil instance returned on the first attempt, and the
// cached, non-nil-but-Context-less instance every later attempt gets back.
func TestHasMachineContext(t *testing.T) {
	tests := []struct {
		name string
		inst *machines.StateMachine
		want bool
	}{
		{
			name: "nil instance",
			inst: nil,
			want: false,
		},
		{
			// InitializeMachineWithContext Adds to the tracker before checking the
			// Start error, so this is what every call after the first one sees for
			// a connection whose cluster was unreachable.
			name: "cached instance whose Start failed",
			inst: &machines.StateMachine{Context: nil},
			want: false,
		},
		{
			// Latent today - no InitFunc returns a typed-nil - but a bare
			// `Context == nil` check would wrongly accept this, and the cast that
			// follows would yield a nil pointer with a nil error.
			name: "boxed typed-nil context",
			inst: &machines.StateMachine{Context: (*kubernetes.MachineCtx)(nil)},
			want: false,
		},
		{
			// Readiness is only about the Context having been assigned; whether it
			// is the *right* type is the caller's cast to report as a real error.
			name: "non-kubernetes context is still assigned",
			inst: &machines.StateMachine{Context: &struct{ notAMachineCtx bool }{}},
			want: true,
		},
		{
			name: "fully initialized instance",
			inst: &machines.StateMachine{Context: &kubernetes.MachineCtx{}},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasMachineContext(tt.inst); got != tt.want {
				t.Fatalf("HasMachineContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// getMachine must not carry a hardcoded allowlist of connection kinds. The
// registerable set is the set of registered connection definitions
// (models/.../connections/*.json) — that is what the Create Connection wizard
// builds its kind list from — so a kind that ships a definition but has no
// bespoke verify action (artifacthub, github) has to get the default machine
// rather than meshery-server-1218.
func TestGetMachineCoversEveryDefinitionKind(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}

	tests := []struct {
		kind string
		// Kinds with a reachability probe run an action on entry to REGISTERED;
		// the rest register with no extra action.
		wantRegisterAction bool
	}{
		{kind: "kubernetes", wantRegisterAction: true},
		{kind: "grafana", wantRegisterAction: true},
		{kind: "prometheus", wantRegisterAction: true},
		{kind: "artifacthub", wantRegisterAction: false},
		{kind: "github", wantRegisterAction: false},
		// Any other definition-backed kind must resolve too, not error.
		{kind: "some-future-kind", wantRegisterAction: false},
	}

	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			mch, err := getMachine(machines.DISCOVERED, tt.kind, connID, userID, log, nil)
			if err != nil {
				t.Fatalf("getMachine(%q) error = %v, want nil", tt.kind, err)
			}
			if mch == nil {
				t.Fatalf("getMachine(%q) returned a nil machine", tt.kind)
			}

			if tt.kind == "kubernetes" {
				// The bespoke Kubernetes machine owns its own state set.
				return
			}

			gotRegisterAction := mch.States[machines.REGISTERED].Action != nil
			if gotRegisterAction != tt.wantRegisterAction {
				t.Fatalf("getMachine(%q) REGISTERED action present = %v, want %v", tt.kind, gotRegisterAction, tt.wantRegisterAction)
			}
			// Every non-Kubernetes kind persists the connection and its
			// credential on connect.
			if mch.States[machines.CONNECTED].Action == nil {
				t.Fatalf("getMachine(%q) has no CONNECTED action; the connection would never be persisted", tt.kind)
			}
		})
	}
}

// fakeProvider satisfies models.Provider via a nil embedded interface, and
// overrides only the one method InitializeMachineWithContext's cache-miss
// path actually calls. Any other Provider method being invoked would panic
// on the nil embedded interface — which is fine, since none of these tests
// exercise those paths.
type fakeProvider struct {
	models.Provider
}

func (fakeProvider) GetGenericPersister() *database.Handler {
	return nil
}

func newHelpersTestLogger(t *testing.T) logger.Handler {
	t.Helper()
	l, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}
	return l
}

func newEmptyTracker() *machines.ConnectionToStateMachineInstanceTracker {
	return &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: map[core.Uuid]*machines.StateMachine{},
	}
}

// TestInitializeMachineWithContext covers meshery#20820: a cache hit must
// respect what the CALLER's initFunc requires, not silently hand back a
// machine initialized for a different purpose.
func TestInitializeMachineWithContext(t *testing.T) {
	t.Run("cache hit, nil initFunc: returns cached instance untouched", func(t *testing.T) {
		id := core.Uuid(uuid.Must(uuid.NewV4()))
		userID := core.Uuid(uuid.Must(uuid.NewV4()))
		tracker := newEmptyTracker()
		cached := &machines.StateMachine{ID: id, Context: nil}
		tracker.Add(id, cached)

		// initFunc == nil mirrors the non-kubernetes registration callers
		// (auto_register.go, connections_handlers.go:67) — they never require
		// a Context, so a nil-Context cache hit is exactly what they expect.
		got, err := InitializeMachineWithContext(nil, context.Background(), id, userID, tracker, newHelpersTestLogger(t), fakeProvider{}, machines.DISCOVERED, "artifacthub", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != cached {
			t.Fatalf("expected the cached instance to be returned as-is")
		}
	})

	t.Run("cache hit, instance already has a context: initFunc is NOT re-run", func(t *testing.T) {
		id := core.Uuid(uuid.Must(uuid.NewV4()))
		userID := core.Uuid(uuid.Must(uuid.NewV4()))
		tracker := newEmptyTracker()
		existingCtx := &struct{ marker string }{marker: "already-set"}
		cached := &machines.StateMachine{ID: id, Context: existingCtx}
		tracker.Add(id, cached)

		called := false
		spyInit := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
			called = true
			return &struct{ marker string }{marker: "should-never-be-used"}, nil, nil
		}

		got, err := InitializeMachineWithContext(nil, context.Background(), id, userID, tracker, newHelpersTestLogger(t), fakeProvider{}, machines.DISCOVERED, "artifacthub", spyInit)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if called {
			t.Fatalf("initFunc must not run again when the cached instance already has a context — unnecessary re-init")
		}
		if got.Context != existingCtx {
			t.Fatalf("expected the existing context to be preserved untouched, got %#v", got.Context)
		}
	})

	// This is the exact scenario from the issue: a connection first cached by
	// a nil-initFunc caller (or one whose Start failed) must NOT stay stuck
	// forever once a caller that actually needs a Context comes along.
	t.Run("cache hit, instance lacks context: initFunc re-runs and Context gets assigned", func(t *testing.T) {
		id := core.Uuid(uuid.Must(uuid.NewV4()))
		userID := core.Uuid(uuid.Must(uuid.NewV4()))
		tracker := newEmptyTracker()
		cached := &machines.StateMachine{ID: id, Context: nil}
		tracker.Add(id, cached)

		newCtx := &struct{ marker string }{marker: "freshly-initialized"}
		called := false
		initFn := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
			called = true
			return newCtx, nil, nil
		}

		got, err := InitializeMachineWithContext(nil, context.Background(), id, userID, tracker, newHelpersTestLogger(t), fakeProvider{}, machines.DISCOVERED, "artifacthub", initFn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !called {
			t.Fatalf("expected initFunc to be re-run for a cached instance missing its context")
		}
		if got != cached {
			t.Fatalf("expected the SAME *StateMachine to be reused (mutated in place), not a new instance")
		}
		if got.Context != newCtx {
			t.Fatalf("expected Context to be updated to the freshly-initialized value, got %#v", got.Context)
		}

		// The tracker must reflect the update too — no stale second copy.
		trackedInst, ok := tracker.Get(id)
		if !ok || trackedInst.Context != newCtx {
			t.Fatalf("tracker entry was not updated in place")
		}
	})

	t.Run("cache hit, instance lacks context: re-init failure surfaces the error", func(t *testing.T) {
		id := core.Uuid(uuid.Must(uuid.NewV4()))
		userID := core.Uuid(uuid.Must(uuid.NewV4()))
		tracker := newEmptyTracker()
		cached := &machines.StateMachine{ID: id, Context: nil}
		tracker.Add(id, cached)

		wantErr := errors.New("cluster unreachable")
		initFn := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
			return nil, nil, wantErr
		}

		got, err := InitializeMachineWithContext(nil, context.Background(), id, userID, tracker, newHelpersTestLogger(t), fakeProvider{}, machines.DISCOVERED, "artifacthub", initFn)
		if !errors.Is(err, wantErr) {
			t.Fatalf("error = %v, want %v", err, wantErr)
		}
		if got != nil {
			t.Fatalf("expected a nil instance on re-init failure, got %#v", got)
		}
	})

	t.Run("cache miss: builds, initializes, and caches a new instance", func(t *testing.T) {
		id := core.Uuid(uuid.Must(uuid.NewV4()))
		userID := core.Uuid(uuid.Must(uuid.NewV4()))
		tracker := newEmptyTracker()

		newCtx := &struct{ marker string }{marker: "built-fresh"}
		initFn := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
			return newCtx, nil, nil
		}

		got, err := InitializeMachineWithContext(nil, context.Background(), id, userID, tracker, newHelpersTestLogger(t), fakeProvider{}, machines.DISCOVERED, "artifacthub", initFn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatalf("expected a non-nil instance")
		}
		if got.Context != newCtx {
			t.Fatalf("expected Context to be set from initFunc, got %#v", got.Context)
		}

		trackedInst, ok := tracker.Get(id)
		if !ok || trackedInst != got {
			t.Fatalf("expected the new instance to be cached under ID")
		}
	})
}