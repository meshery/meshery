package models

import (
	"testing"

	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// TestLocalGetEventsReportsReadCount checks that the local provider fills
// readCount, which the Notification Center shows on its "Read" chip. The
// count covers the user's and the system's read events, whatever filter the
// page itself uses, and leaves out other users' events.
func TestLocalGetEventsReportsReadCount(t *testing.T) {
	persister := newCorrelationTestPersister(t)
	provider := &DefaultLocalProvider{EventsPersister: persister}

	userID, sysID, otherID := mustUUID(t), mustUUID(t), mustUUID(t)
	persistWithStatus := func(owner core.Uuid, status events.EventStatus) {
		t.Helper()
		ev := events.NewEvent().
			FromOwner(owner).
			WithCategory("connection").
			WithAction("update").
			WithSeverity(events.Informational).
			WithDescription("test event").
			Build()
		ev.Status = status
		if _, err := persister.PersistEvent(*ev, ""); err != nil {
			t.Fatalf("persist event: %v", err)
		}
	}

	// Before anything is read, readCount is zero.
	persistWithStatus(userID, events.Unread)
	resp, err := provider.GetEvents("", &events.EventsFilter{Status: events.Unread, Limit: 10}, 0, userID, sysID)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if resp.ReadCount != 0 {
		t.Errorf("readCount = %d, want 0 before any event is read", resp.ReadCount)
	}

	persistWithStatus(userID, events.Read)
	persistWithStatus(userID, events.Read)
	persistWithStatus(sysID, events.Read)
	// An unread system event must not be counted: the owner check is
	// grouped, so the status condition applies to both owners.
	persistWithStatus(sysID, events.Unread)
	persistWithStatus(otherID, events.Read)

	// The Notification Center header asks for unread events and reads
	// readCount from the same response.
	resp, err = provider.GetEvents("", &events.EventsFilter{Status: events.Unread, Limit: 10}, 0, userID, sysID)
	if err != nil {
		t.Fatalf("GetEvents: %v", err)
	}
	if resp.TotalCount != 2 {
		t.Errorf("totalCount = %d, want 2 unread events (one user, one system)", resp.TotalCount)
	}
	if resp.ReadCount != 3 {
		t.Errorf("readCount = %d, want 3 (two user events and one system event)", resp.ReadCount)
	}
}
