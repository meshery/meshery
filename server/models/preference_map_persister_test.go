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

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadFromPersister() got = %+v, want %+v", got, want)
	}
}

func TestMapPreferencePersisterDeleteFromPersister(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	pref := &Preference{
		SelectedOrganizationID: "org-delete-test",
		AnonymousUsageStats:    false,
		AnonymousPerfResults:   false,
	}

	if err := persister.WriteToPersister("user-to-delete", pref); err != nil {
		t.Fatalf("failed to write preference: %v", err)
	}

	// When
	if err := persister.DeleteFromPersister("user-to-delete"); err != nil {
		t.Fatalf("expected successful delete, got error: %v", err)
	}

	// Then - reading again should return defaults
	got, err := persister.ReadFromPersister("user-to-delete")
	if err != nil {
		t.Fatalf("expected read after delete to return defaults, got error: %v", err)
	}

	if got == nil {
		t.Fatal("expected default preference after deletion, got nil")
	}

	if got.SelectedOrganizationID != "" {
		t.Errorf("expected empty SelectedOrganizationID after deletion, got %q", got.SelectedOrganizationID)
	}

	if !got.AnonymousUsageStats || !got.AnonymousPerfResults {
		t.Error("expected default anonymous flags to be true after deletion")
	}
}

func TestMapPreferencePersisterClosePersister(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	persister.ClosePersister()

	// Then - all operations should return ErrDBConnection
	if _, err := persister.ReadFromPersister("user-1"); err != ErrDBConnection {
		t.Errorf("ReadFromPersister on closed persister: expected ErrDBConnection, got %v", err)
	}

	if err := persister.WriteToPersister("user-1", NewDefaultPreference()); err != ErrDBConnection {
		t.Errorf("WriteToPersister on closed persister: expected ErrDBConnection, got %v", err)
	}

	if err := persister.DeleteFromPersister("user-1"); err != ErrDBConnection {
		t.Errorf("DeleteFromPersister on closed persister: expected ErrDBConnection, got %v", err)
	}
}

func TestMapPreferencePersisterEmptyUserIDValidation(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	// Then - empty userID should return ErrUserID
	if _, err := persister.ReadFromPersister(""); err != ErrUserID {
		t.Errorf("ReadFromPersister with empty userID: expected ErrUserID, got %v", err)
	}

	if err := persister.WriteToPersister("", NewDefaultPreference()); err != ErrUserID {
		t.Errorf("WriteToPersister with empty userID: expected ErrUserID, got %v", err)
	}

	if err := persister.DeleteFromPersister(""); err != ErrUserID {
		t.Errorf("DeleteFromPersister with empty userID: expected ErrUserID, got %v", err)
	}
}

func TestMapPreferencePersisterWriteNilData(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	// Then - writing nil preference data should return ErrNilConfigData
	if err := persister.WriteToPersister("user-1", nil); err != ErrNilConfigData {
		t.Errorf("WriteToPersister with nil data: expected ErrNilConfigData, got %v", err)
	}
}

func TestMapPreferencePersisterCorruptedSessionFallback(t *testing.T) {
	// Given
	persister := newTestMapPreferencePersister(t)
	defer persister.db.Clear()

	// Store corrupted/invalid type in sync.Map directly
	persister.db.Store("corrupted-user", "not-a-preference-struct")

	// When
	got, err := persister.ReadFromPersister("corrupted-user")
	if err != nil {
		t.Fatalf("expected corrupted session fallback to not return an error, got: %v", err)
	}

	// Then - should return fresh default preference
	if got == nil {
		t.Fatal("expected default preference on corrupted session fallback, got nil")
	}

	if !got.AnonymousUsageStats || !got.AnonymousPerfResults {
		t.Error("expected default flags to be true for corrupted session fallback")
	}
}
