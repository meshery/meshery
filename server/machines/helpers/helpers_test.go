package helpers

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/grafana"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/machines/prometheus"
	"github.com/meshery/meshery/server/models"
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
			// Defensive check: an uninitialized or broken instance whose Context
			// was never assigned must be rejected even if encountered.
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
		// Kinds with a reachability probe run their own action on entry to
		// REGISTERED; the rest keep the default machine's no-op action.
		wantRegisterAction machines.Action
	}{
		{kind: "kubernetes"},
		{kind: "grafana", wantRegisterAction: &grafana.RegisterAction{}},
		{kind: "prometheus", wantRegisterAction: &prometheus.RegisterAction{}},
		{kind: "artifacthub", wantRegisterAction: &machines.DefaultAction{}},
		{kind: "github", wantRegisterAction: &machines.DefaultAction{}},
		// Any other definition-backed kind must resolve too, not error.
		{kind: "some-future-kind", wantRegisterAction: &machines.DefaultAction{}},
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

			got := reflect.TypeOf(mch.States[machines.REGISTERED].Action)
			if want := reflect.TypeOf(tt.wantRegisterAction); got != want {
				t.Fatalf("getMachine(%q) REGISTERED action = %v, want %v", tt.kind, got, want)
			}
			// Every non-Kubernetes kind persists the connection and its
			// credential on connect.
			if mch.States[machines.CONNECTED].Action == nil {
				t.Fatalf("getMachine(%q) has no CONNECTED action; the connection would never be persisted", tt.kind)
			}
		})
	}
}

// Kinds without a verification step must actually be able to enter
// REGISTERED from DISCOVERED (the Create Connection wizard's register step),
// not just build a machine: SendEvent rejects a transition into a state with
// no action.
func TestGetMachineRegistersKindsWithoutVerifyAction(t *testing.T) {
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)

	for _, kind := range []string{"github", "artifacthub", "some-future-kind"} {
		t.Run(kind, func(t *testing.T) {
			mch, err := getMachine(machines.DISCOVERED, kind, uuid.Must(uuid.NewV4()).String(), userID, log, nil)
			if err != nil {
				t.Fatalf("getMachine(%q) error = %v", kind, err)
			}

			if _, err := mch.SendEvent(ctx, machines.Register, nil); err != nil {
				t.Fatalf("SendEvent(Register) error = %v, want nil", err)
			}
			if mch.CurrentState != machines.REGISTERED {
				t.Fatalf("CurrentState = %q, want %q", mch.CurrentState, machines.REGISTERED)
			}

			want := machines.Events{machines.Connect: machines.CONNECTED, machines.Ignore: machines.IGNORED}
			if got := mch.States[mch.CurrentState].Events; !reflect.DeepEqual(got, want) {
				t.Fatalf("events available from %q = %v, want %v", mch.CurrentState, got, want)
			}
		})
	}
}

// InitializeMachineWithContext must not add a state machine to smInstanceTracker
// when its initialization/Start fails. Caching an uninitialized instance strands
// a broken machine with a nil Context in the tracker, causing subsequent calls
// to return that cached instance with a nil error.
func TestInitializeMachineWithContext_DoesNotCacheOnStartError(t *testing.T) {
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	ctx := context.Background()

	expectedErr := errors.New("simulated init failure")
	failingInit := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return nil, nil, expectedErr
	}

	// First attempt: initFunc fails
	inst, err := InitializeMachineWithContext(
		nil,
		ctx,
		connID,
		userID,
		tracker,
		log,
		nil,
		machines.InitialState,
		"kubernetes",
		failingInit,
	)
	if err == nil {
		t.Fatal("expected InitializeMachineWithContext to fail, got nil")
	}
	if inst != nil {
		t.Fatalf("expected nil instance on error, got %v", inst)
	}

	// Crucial check: the tracker must NOT contain the failed instance
	if cached, ok := tracker.Get(connID); ok {
		t.Fatalf("tracker cached failed machine instance: %v", cached)
	}

	// Second attempt: initFunc succeeds (simulating recovery)
	successfulInit := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		return &kubernetes.MachineCtx{}, nil, nil
	}

	inst, err = InitializeMachineWithContext(
		nil,
		ctx,
		connID,
		userID,
		tracker,
		log,
		nil,
		machines.InitialState,
		"kubernetes",
		successfulInit,
	)
	if err != nil {
		t.Fatalf("expected InitializeMachineWithContext to succeed on retry, got %v", err)
	}
	if inst == nil {
		t.Fatal("expected non-nil instance on success")
	}

	// Now the tracker SHOULD contain it
	cached, ok := tracker.Get(connID)
	if !ok || cached != inst {
		t.Fatalf("expected tracker to cache successfully initialized machine, got ok=%v, cached=%v", ok, cached)
	}
	if !HasMachineContext(inst) {
		t.Fatal("expected cached instance to have a valid context")
	}
}

// InitializeMachineWithContext must serialize initialization per connection ID
// so that concurrent callers do not spawn redundant machines or race on Start.
func TestInitializeMachineWithContext_SerializesConcurrentCalls(t *testing.T) {
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	userID := core.Uuid(uuid.Must(uuid.NewV4()))
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger.New() error = %v", err)
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine),
	}
	ctx := context.Background()

	var initCount atomic.Int32
	initFunc := func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
		initCount.Add(1)
		time.Sleep(10 * time.Millisecond)
		return &kubernetes.MachineCtx{}, nil, nil
	}

	const goroutines = 10
	var wg sync.WaitGroup
	instances := make([]*machines.StateMachine, goroutines)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			inst, err := InitializeMachineWithContext(
				nil,
				ctx,
				connID,
				userID,
				tracker,
				log,
				nil,
				machines.InitialState,
				"kubernetes",
				initFunc,
			)
			if err != nil {
				t.Errorf("goroutine %d failed: %v", idx, err)
			}
			instances[idx] = inst
		}(i)
	}

	wg.Wait()

	if count := initCount.Load(); count != 1 {
		t.Fatalf("expected initFunc to run exactly once, ran %d times", count)
	}

	for i := 1; i < goroutines; i++ {
		if instances[i] != instances[0] {
			t.Fatalf("expected all goroutines to receive the same instance, goroutine %d received %p != %p", i, instances[i], instances[0])
		}
	}

	// Verify that the per-key mutex was cleanly evicted from initLockMap after completion
	initLockMapMu.Lock()
	remainingLocks := len(initLockMap)
	initLockMapMu.Unlock()
	if remainingLocks != 0 {
		t.Fatalf("expected initLockMap to be empty after completion, found %d entries", remainingLocks)
	}
}
