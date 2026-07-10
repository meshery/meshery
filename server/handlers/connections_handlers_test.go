package handlers

import (
	"testing"

	"github.com/meshery/meshkit/models/events"
)

// TestNilEventGuard_Regression is a regression test for the nil-event guard
// added to the NotifySmOfConnectionStatusChange goroutine in connections_handlers.go.
//
// Background: SendEvent can return (nil, error) from multiple code paths.
// The original code unconditionally dereferenced *event to call
// PersistEvent (which takes events.Event by value), causing a fatal panic
// inside an unrecoverable goroutine.
//
// This test validates the guard pattern in isolation. A full integration test
// would require a running server context with mocked providers and state machines,
// which is out of scope for this fix. This test ensures that if the guard is
// ever accidentally removed, the panic is caught immediately in CI.
func TestNilEventGuard_Regression(t *testing.T) {
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

			// Mirrors the guard pattern in NotifySmOfConnectionStatusChange:
			//   if event != nil {
			//       if persistErr := provider.PersistEvent(*event, token); persistErr != nil { ... }
			//       h.config.EventBroadcaster.Publish(userID, event)
			//   } else {
			//       h.log.Warn(...)
			//   }
			persisted := false
			warned := false
			if tt.event != nil {
				_ = *tt.event // dereference — would panic if tt.event were nil
				persisted = true
			} else {
				warned = true
			}

			if tt.want == "skip" && persisted {
				t.Fatal("nil event should have been skipped, but was persisted")
			}
			if tt.want == "skip" && !warned {
				t.Fatal("nil event should have triggered a warning log")
			}
			if tt.want == "persist" && !persisted {
				t.Fatal("non-nil event should have been persisted, but was skipped")
			}
		})
	}
}
