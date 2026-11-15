package models

import (
	"errors"
	"testing"

	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestSessionPreferencePersister(t *testing.T) *SessionPreferencePersister {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&UserPreference{}); err != nil {
		t.Fatalf("failed to migrate user preference schema: %v", err)
	}

	return &SessionPreferencePersister{
		DB: &database.Handler{DB: db},
	}
}

func TestSessionPreferencePersisterReadFromPersisterReturnsDefaultsWhenMissing(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	pref, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("expected missing preference to return defaults, got error: %v", err)
	}

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

func TestSessionPreferencePersisterReadFromPersisterReturnsPersistedPreference(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	want := &Preference{
		AnonymousUsageStats:    false,
		AnonymousPerfResults:   false,
		SelectedOrganizationID: "org-1",
		DashboardPreferences: map[string]interface{}{
			"theme": "dark",
		},
	}

	if err := persister.WriteToPersister("user-1", want); err != nil {
		t.Fatalf("failed to persist preference: %v", err)
	}

	got, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("failed to read persisted preference: %v", err)
	}

	if got == nil {
		t.Fatal("expected persisted preference, got nil")
	}

	if got.AnonymousUsageStats != want.AnonymousUsageStats {
		t.Errorf("got AnonymousUsageStats=%v, want %v", got.AnonymousUsageStats, want.AnonymousUsageStats)
	}

	if got.AnonymousPerfResults != want.AnonymousPerfResults {
		t.Errorf("got AnonymousPerfResults=%v, want %v", got.AnonymousPerfResults, want.AnonymousPerfResults)
	}

	if got.SelectedOrganizationID != want.SelectedOrganizationID {
		t.Errorf("got SelectedOrganizationID=%q, want %q", got.SelectedOrganizationID, want.SelectedOrganizationID)
	}

	if got.DashboardPreferences["theme"] != want.DashboardPreferences["theme"] {
		t.Errorf("got DashboardPreferences[theme]=%v, want %v", got.DashboardPreferences["theme"], want.DashboardPreferences["theme"])
	}
}

func TestSessionPreferencePersisterReadFromPersisterReturnsDefaultsForStoredNull(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	if err := persister.DB.Model(&UserPreference{}).Create(&UserPreference{
		ID:              "user-1",
		PreferenceBytes: []byte("null"),
	}).Error; err != nil {
		t.Fatalf("failed to seed null preference: %v", err)
	}

	pref, err := persister.ReadFromPersister("user-1")
	if err != nil {
		t.Fatalf("expected stored null preference to return defaults, got error: %v", err)
	}

	if pref == nil {
		t.Fatal("expected default preference for stored null, got nil")
	}

	if !pref.AnonymousUsageStats {
		t.Error("expected AnonymousUsageStats to default to true for stored null preference")
	}

	if !pref.AnonymousPerfResults {
		t.Error("expected AnonymousPerfResults to default to true for stored null preference")
	}
}

func TestSessionPreferencePersisterNilDBAndEmptyUserGuards(t *testing.T) {
	nilPersister := &SessionPreferencePersister{DB: nil}

	// Nil DB tests
	if _, err := nilPersister.ReadFromPersister("user-1"); !errors.Is(err, ErrDBConnection) {
		t.Errorf("ReadFromPersister on nil DB: want ErrDBConnection, got %v", err)
	}

	if err := nilPersister.WriteToPersister("user-1", NewDefaultPreference()); !errors.Is(err, ErrDBConnection) {
		t.Errorf("WriteToPersister on nil DB: want ErrDBConnection, got %v", err)
	}

	if err := nilPersister.DeleteFromPersister("user-1"); !errors.Is(err, ErrDBConnection) {
		t.Errorf("DeleteFromPersister on nil DB: want ErrDBConnection, got %v", err)
	}

	// Valid DB but empty user ID tests
	persister := newTestSessionPreferencePersister(t)

	if _, err := persister.ReadFromPersister(""); !errors.Is(err, ErrUserID) {
		t.Errorf("ReadFromPersister with empty userID: want ErrUserID, got %v", err)
	}

	if err := persister.WriteToPersister("", NewDefaultPreference()); !errors.Is(err, ErrUserID) {
		t.Errorf("WriteToPersister with empty userID: want ErrUserID, got %v", err)
	}

	if err := persister.DeleteFromPersister(""); !errors.Is(err, ErrUserID) {
		t.Errorf("DeleteFromPersister with empty userID: want ErrUserID, got %v", err)
	}

	// Nil preference data test
	if err := persister.WriteToPersister("user-1", nil); !errors.Is(err, ErrNilConfigData) {
		t.Errorf("WriteToPersister with nil data: want ErrNilConfigData, got %v", err)
	}
}

func TestSessionPreferencePersisterDeleteFromPersister(t *testing.T) {
	persister := newTestSessionPreferencePersister(t)

	// Deleting nonexistent user preference should return ErrResultNotFound
	err := persister.DeleteFromPersister("nonexistent-user")
	if err == nil {
		t.Fatal("expected error deleting nonexistent user preference, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}

	// Write and then delete successfully
	pref := NewDefaultPreference()
	pref.SelectedOrganizationID = "org-test"
	if err := persister.WriteToPersister("user-to-delete", pref); err != nil {
		t.Fatalf("failed to write preference: %v", err)
	}

	if err := persister.DeleteFromPersister("user-to-delete"); err != nil {
		t.Fatalf("failed to delete existing preference: %v", err)
	}

	// Verifying it was deleted
	got, err := persister.ReadFromPersister("user-to-delete")
	if err != nil {
		t.Fatalf("unexpected error after delete: %v", err)
	}
	if got.SelectedOrganizationID == "org-test" {
		t.Errorf("expected default preference after deletion, got org-test")
	}
}
