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

	stored, ok := persister.db.Load("user-1")
	if !ok {
		t.Fatal("expected missing preference read to seed the default preference")
	}

	seededPref, ok := stored.(*Preference)
	if !ok {
		t.Fatalf("expected stored preference to be *Preference, got %T", stored)
	}

	if !reflect.DeepEqual(seededPref, pref) {
		t.Errorf("expected seeded preference to match returned preference, got %+v want %+v", seededPref, pref)
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

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadFromPersister() got = %+v, want %+v", got, want)
	}
}
