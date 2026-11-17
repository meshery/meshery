package models

import (
	"slices"
	"testing"

	"github.com/meshery/meshkit/models/events"
)

// The event types offer the categories and actions of the events GetEvents
// lists: the user's own and the Meshery instance's, not another user's.
func TestGetEventTypesCoversUserAndSystemEvents(t *testing.T) {
	db := newMigratedDB(t)
	provider := &DefaultLocalProvider{EventsPersister: &EventsPersister{DB: db}}
	userID, sysID, otherUserID := mustUUID(t), mustUUID(t), mustUUID(t)

	for _, event := range []*events.Event{
		events.NewEvent().FromOwner(userID).FromSystem(sysID).WithCategory("connection").WithAction("update").
			WithSeverity(events.Informational).WithDescription("the user's event").Build(),
		events.NewEvent().FromOwner(sysID).FromSystem(sysID).WithCategory("meshsync").WithAction("flush").
			WithSeverity(events.Informational).WithDescription("a system event").Build(),
		events.NewEvent().FromOwner(otherUserID).FromSystem(sysID).WithCategory("design").WithAction("deploy").
			WithSeverity(events.Informational).WithDescription("another user's event").Build(),
	} {
		if _, err := provider.PersistEvent(*event, ""); err != nil {
			t.Fatalf("persist %q: %v", event.Description, err)
		}
	}

	types, err := provider.GetEventTypes("", userID, sysID)
	if err != nil {
		t.Fatalf("event types: %v", err)
	}
	slices.Sort(types.Category)
	slices.Sort(types.Action)

	if want := []string{"connection", "meshsync"}; !slices.Equal(types.Category, want) {
		t.Errorf("categories = %v, want %v", types.Category, want)
	}
	if want := []string{"flush", "update"}; !slices.Equal(types.Action, want) {
		t.Errorf("actions = %v, want %v", types.Action, want)
	}
}
