package helpers

import (
	"context"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/grafana"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/machines/prometheus"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
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
