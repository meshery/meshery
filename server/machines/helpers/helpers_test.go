package helpers

import (
	"testing"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
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
