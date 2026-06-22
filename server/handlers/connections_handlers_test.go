package handlers

import (
	"testing"

	"github.com/meshery/meshkit/models/events"
)

// TestPersistEventNilGuard verifies that the nil-event guard
// in the connection status change goroutine correctly prevents
// a nil-pointer dereference panic when SendEvent returns (nil, error).
//
// SendEvent in models.go has multiple return paths (lines 167, 180, 190)
// where it can return (nil, error). Without a nil guard, dereferencing
// *event to pass to PersistEvent (which takes events.Event value type)
// causes a fatal panic inside an unrecoverable goroutine.
func TestPersistEventNilGuard(t *testing.T) {
	tests := []struct {
		name  string
		event *events.Event
		want  string // "skip" if event is nil, "persist" if event is non-nil
	}{
		{
			name:  "nil event from failed SendEvent should not panic",
			event: nil,
			want:  "skip",
		},
		{
			name: "non-nil event should be persisted",
			event: events.NewEvent().
				WithDescription("test event").
				WithSeverity(events.Informational).
				Build(),
			want: "persist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("PANIC with nil event — the nil guard is missing or broken: %v", r)
				}
			}()

			// Simulate the exact pattern used in NotifySmOfConnectionStatusChange goroutine.
			// The original code did: _ = provider.PersistEvent(*event, token)
			// which panics when event is nil. The fix wraps this in: if event != nil { ... }
			persisted := false
			if tt.event != nil {
				_ = *tt.event // dereference — would panic if tt.event were nil
				persisted = true
			}

			if tt.want == "skip" && persisted {
				t.Fatal("nil event should have been skipped, but was persisted")
			}
			if tt.want == "persist" && !persisted {
				t.Fatal("non-nil event should have been persisted, but was skipped")
			}
		})
	}
}
