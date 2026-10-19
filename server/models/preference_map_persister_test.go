package models

import (
	"reflect"
	"testing"
)

func newTestMapPreferencePersister(t *testing.T) *MapPreferencePersister {
	t.Helper()

	persister, err := NewMapPreferencePersister()
	if err != nil {
		t.Fatalf("failed to create map preference persister: %v", err)
	}

	return persister
}

func TestMapPreferencePersisterReadFromPersisterReturnsDefaultsWhenMissing(t *testing.T) {
	t.Skip("TODO: https://github.com/meshery/meshery/issues/18238 - implement MapPreferencePersister tests")
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	// When
	pref, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("expected missing preference to return defaults, got error: %v", err)
	}

	// Then
	if pref == nil {
		t.Fatal("expected default preference, got nil")
	}

	if !pref.AnonymousUsageStats {
		t.Error("expected AnonymousUsageStats to default to true")
	}

	if !pref.AnonymousPerfResults {
		t.Error("expected AnonymousPerfResults to default to true")
	}

}

func TestMapPreferencePersisterReadFromPersisterReturnsPersistedPreference(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	want := &Preference{
		AnonymousUsageStats:    false,
		AnonymousPerfResults:   false,
		SelectedOrganizationID: "org-1",
		DashboardPreferences: map[string]interface{}{
			"theme": "dark",
		},
	}

	// When
	if err := persister.WriteToPersister("user-1", want); err != nil {
		t.Fatalf("failed to persist preference: %v", err)
	}

	// Then
	got, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read persisted preference: %v", err)
	}

	if got == nil {
		t.Fatal("expected persisted preference, got nil")
	}

	if got.AnonymousUsageStats != want.AnonymousUsageStats {
		t.Errorf("AnonymousUsageStats = %v, want %v", got.AnonymousUsageStats, want.AnonymousUsageStats)
	}
	if got.AnonymousPerfResults != want.AnonymousPerfResults {
		t.Errorf("AnonymousPerfResults = %v, want %v", got.AnonymousPerfResults, want.AnonymousPerfResults)
	}
	if got.SelectedOrganizationID != want.SelectedOrganizationID {
		t.Errorf("SelectedOrganizationID = %q, want %q", got.SelectedOrganizationID, want.SelectedOrganizationID)
	}
	if !reflect.DeepEqual(got.DashboardPreferences, want.DashboardPreferences) {
		t.Errorf("DashboardPreferences = %+v, want %+v", got.DashboardPreferences, want.DashboardPreferences)
	}
	if got.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestMapPreferencePersisterReadReturnsDefensiveCopy(t *testing.T) {
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	initial := &Preference{
		AnonymousUsageStats:  false,
		AnonymousPerfResults: false,
		LoadTestPreferences: &LoadTestPreferences{
			QueriesPerSecond: 10,
		},
		DashboardPreferences: map[string]interface{}{
			"theme": "dark",
		},
		SelectedWorkspaceForOrganizations: map[string]string{
			"org-1": "ws-1",
		},
	}

	if err := persister.WriteToPersister("user-1", initial); err != nil {
		t.Fatalf("failed to write preference: %v", err)
	}

	firstRead, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read preference: %v", err)
	}

	// Mutate the returned preference object
	firstRead.AnonymousUsageStats = true
	firstRead.LoadTestPreferences.QueriesPerSecond = 999
	firstRead.DashboardPreferences["theme"] = "light"
	firstRead.DashboardPreferences["newKey"] = 123
	firstRead.SelectedWorkspaceForOrganizations["org-1"] = "ws-mutated"

	// Read again and ensure the store was not affected
	secondRead, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read preference again: %v", err)
	}

	if secondRead.AnonymousUsageStats != false {
		t.Errorf("stored AnonymousUsageStats was mutated: got true, want false")
	}
	if secondRead.LoadTestPreferences == nil || secondRead.LoadTestPreferences.QueriesPerSecond != 10 {
		t.Errorf("stored QueriesPerSecond was mutated: got %v, want 10", secondRead.LoadTestPreferences)
	}
	if secondRead.DashboardPreferences["theme"] != "dark" || secondRead.DashboardPreferences["newKey"] != nil {
		t.Errorf("stored DashboardPreferences was mutated: got %+v, want theme=dark only", secondRead.DashboardPreferences)
	}
	if secondRead.SelectedWorkspaceForOrganizations["org-1"] != "ws-1" {
		t.Errorf("stored SelectedWorkspaceForOrganizations was mutated: got %v, want ws-1", secondRead.SelectedWorkspaceForOrganizations["org-1"])
	}
}

func TestMapPreferencePersisterWriteStoresDefensiveCopy(t *testing.T) {
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	initial := &Preference{
		AnonymousUsageStats:  false,
		AnonymousPerfResults: false,
		LoadTestPreferences: &LoadTestPreferences{
			QueriesPerSecond: 10,
		},
		DashboardPreferences: map[string]interface{}{
			"theme": "dark",
		},
	}

	if err := persister.WriteToPersister("user-1", initial); err != nil {
		t.Fatalf("failed to write preference: %v", err)
	}

	// Mutate the initial object after write
	initial.AnonymousUsageStats = true
	initial.LoadTestPreferences.QueriesPerSecond = 888
	initial.DashboardPreferences["theme"] = "mutated-after-write"

	got, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read preference: %v", err)
	}

	if got.AnonymousUsageStats != false {
		t.Errorf("persisted AnonymousUsageStats mutated after write: got true, want false")
	}
	if got.LoadTestPreferences == nil || got.LoadTestPreferences.QueriesPerSecond != 10 {
		t.Errorf("persisted QueriesPerSecond mutated after write: got %v, want 10", got.LoadTestPreferences)
	}
	if got.DashboardPreferences["theme"] != "dark" {
		t.Errorf("persisted DashboardPreferences mutated after write: got %+v, want theme=dark", got.DashboardPreferences)
	}
}
