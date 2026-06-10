package models

import (
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestProviderWithCredentialDB(t *testing.T) *DefaultLocalProvider {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&Credential{}); err != nil {
		t.Fatalf("failed to migrate credential schema: %v", err)
	}

	return &DefaultLocalProvider{
		GenericPersister: &database.Handler{DB: db},
	}
}

// TestSaveUserCredentialReturnsPopulatedCredential verifies that SaveUserCredential
// returns the saved credential rather than nil. A nil return caused the HTTP handler
// to panic when accessing the returned credential's fields (e.g. .Name).
func TestSaveUserCredentialReturnsPopulatedCredential(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)

	cred := &Credential{Name: "test-key", Type: "token"}

	got, err := provider.SaveUserCredential("tok", cred)
	if err != nil {
		t.Fatalf("unexpected error saving credential: %v", err)
	}
	if got == nil {
		t.Fatal("SaveUserCredential returned nil credential: callers dereference the return value and will panic")
	}
	if got.Name != cred.Name {
		t.Errorf("got Name=%q, want %q", got.Name, cred.Name)
	}
	if got.Type != cred.Type {
		t.Errorf("got Type=%q, want %q", got.Type, cred.Type)
	}
}

// TestEventsPersisterPersistEventAcceptsStructValue verifies that PersistEvent
// successfully writes events when called with a struct value (the shape
// callers pass via `PersistSystemEvent(*event)`). The function previously
// passed `event` directly to gorm's `Save`, which requires a pointer; gorm
// returned "invalid value, should be pointer to struct or slice" and every
// system event was silently dropped — visible on /management/connections via
// repeated "failed to persist event" logs whenever the controller emitted a
// connection event.
func TestEventsPersisterPersistEventAcceptsStructValue(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&events.Event{}); err != nil {
		t.Fatalf("failed to migrate event schema: %v", err)
	}

	persister := &EventsPersister{DB: &database.Handler{DB: db}}

	eventID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate event id: %v", err)
	}
	systemID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate system id: %v", err)
	}

	built := events.NewEvent().
		FromSystem(systemID).
		WithCategory("connection").
		WithAction("update").
		WithSeverity(events.Informational).
		WithDescription("test event").
		Build()
	built.ID = eventID

	if err := persister.PersistEvent(*built, ""); err != nil {
		t.Fatalf("PersistEvent returned error: %v", err)
	}

	var stored events.Event
	if err := db.First(&stored, "id = ?", eventID).Error; err != nil {
		t.Fatalf("event was not persisted: %v", err)
	}
	if stored.Category != "connection" || stored.Action != "update" {
		t.Errorf("stored event mismatch: got category=%q action=%q",
			stored.Category, stored.Action)
	}
}

// TestEventsPersisterPersistSystemEventAcceptsStructValue mirrors the above
// for PersistSystemEvent — the controller helper path called from
// MesheryControllersHelper.emitEvent that surfaced the original bug.
func TestEventsPersisterPersistSystemEventAcceptsStructValue(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&events.Event{}); err != nil {
		t.Fatalf("failed to migrate event schema: %v", err)
	}

	persister := &EventsPersister{DB: &database.Handler{DB: db}}

	systemID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate system id: %v", err)
	}
	eventID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate event id: %v", err)
	}

	built := events.NewEvent().
		FromSystem(systemID).
		WithCategory("connection").
		WithAction("update").
		WithSeverity(events.Informational).
		Build()
	built.ID = eventID

	if err := persister.PersistSystemEvent(*built); err != nil {
		t.Fatalf("PersistSystemEvent returned error: %v", err)
	}
}

// TestRemoteFilterFileShortGitHubURLReturnsError verifies that RemoteFilterFile
// returns a descriptive error for short GitHub URLs instead of panicking with an
// out-of-bounds index access. Previously parsedPath[3] was accessed before the
// length guard, crashing the server on any URL shorter than 4 path segments.
func TestRemoteFilterFileShortGitHubURLReturnsError(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		wantErrSnip string
	}{
		{
			name:        "host only",
			url:         "https://github.com",
			wantErrSnip: "malformed URL",
		},
		{
			name:        "owner only",
			url:         "https://github.com/owner",
			wantErrSnip: "malformed URL",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider := &DefaultLocalProvider{}
			_, err := provider.RemoteFilterFile(nil, tc.url, "", false, "")
			if err == nil {
				t.Errorf("expected error for URL %q, got nil", tc.url)
				return
			}
			if !strings.Contains(err.Error(), tc.wantErrSnip) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrSnip)
			}
		})
	}
}

// TestRemotePatternFileShortGitHubURLReturnsError verifies equivalent safe
// behaviour in RemotePatternFile for the same class of short URLs.
func TestRemotePatternFileShortGitHubURLReturnsError(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		wantErrSnip string
	}{
		{
			name:        "host only",
			url:         "https://github.com",
			wantErrSnip: "malformed URL",
		},
		{
			name:        "owner only",
			url:         "https://github.com/owner",
			wantErrSnip: "malformed URL",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider := &DefaultLocalProvider{}
			_, err := provider.RemotePatternFile(nil, tc.url, "", false)
			if err == nil {
				t.Errorf("expected error for URL %q, got nil", tc.url)
				return
			}
			if !strings.Contains(err.Error(), tc.wantErrSnip) {
				t.Errorf("error %q does not contain %q", err.Error(), tc.wantErrSnip)
			}
		})
	}
}

func TestDefaultLocalProviderRemoveExtension_RemovesMatchingNavigatorExtension(t *testing.T) {
	provider := &DefaultLocalProvider{}
	provider.Initialize()
	provider.Extensions.Navigator = NavigatorExtensions{
		{Title: "Kanvas"},
		{Title: "MeshMap Snapshot"},
	}

	if err := provider.RemoveExtension("navigator", "Kanvas"); err != nil {
		t.Fatalf("RemoveExtension returned error: %v", err)
	}

	if len(provider.Extensions.Navigator) != 1 {
		t.Fatalf("expected 1 navigator extension after removal, got %d", len(provider.Extensions.Navigator))
	}
	if provider.Extensions.Navigator[0].Title != "MeshMap Snapshot" {
		t.Fatalf("unexpected extension retained: %+v", provider.Extensions.Navigator[0])
	}
}

func TestDefaultLocalProviderRemoveExtension_ReturnsErrorForMissingExtension(t *testing.T) {
	provider := &DefaultLocalProvider{}
	provider.Initialize()
	provider.Extensions.Navigator = NavigatorExtensions{{Title: "MeshMap Snapshot"}}

	err := provider.RemoveExtension("navigator", "Kanvas")
	if err == nil {
		t.Fatal("expected error removing missing navigator extension, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestDefaultLocalProviderInstallExtension_RequiresPackageWhenAssetsMissing(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	// Exercise the real download path (not skipped) so the missing package URL
	// is surfaced as an error. viper is the source of truth for this flag;
	// capture and restore its original value to avoid leaking state to other tests.
	origSkip := viper.Get(SKIP_DOWNLOAD_EXTENSIONS_ENV)
	defer viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, origSkip)
	viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, false)

	provider := &DefaultLocalProvider{}
	provider.Initialize()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	provider.Log = log

	err = provider.InstallExtension("navigator", "", map[string]interface{}{
		"title":     "Kanvas",
		"component": "/provider/navigator/meshmap/index.js",
		"href": map[string]interface{}{
			"uri": "/meshmap",
		},
	})
	if err == nil {
		t.Fatal("expected install to fail when extension assets are missing and no package URL is provided")
	}
	if !strings.Contains(err.Error(), "package URL is required") {
		t.Fatalf("expected missing package URL error, got %v", err)
	}
	if len(provider.Extensions.Navigator) != 0 {
		t.Fatalf("install should not mutate navigator extensions on failure, got %+v", provider.Extensions.Navigator)
	}
}

func TestDefaultLocalProviderInstallExtension_ReplacesMatchingNavigatorExtension(t *testing.T) {
	// viper, not the OS environment, is the source of truth for this flag, and
	// the test binary never calls viper.AutomaticEnv(); set it via viper so the
	// package download is deterministically skipped without network access.
	// Capture and restore the original value to avoid leaking state to other tests.
	origSkip := viper.Get(SKIP_DOWNLOAD_EXTENSIONS_ENV)
	defer viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, origSkip)
	viper.Set(SKIP_DOWNLOAD_EXTENSIONS_ENV, true)

	provider := &DefaultLocalProvider{}
	provider.Initialize()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	provider.Log = log
	provider.Extensions.Navigator = NavigatorExtensions{
		{Title: "Kanvas", Component: "/provider/navigator/meshmap/index.js?packageVersion=old"},
	}

	err = provider.InstallExtension("navigator", "", map[string]interface{}{
		"title":     "Kanvas",
		"component": "/provider/navigator/meshmap/index.js?packageVersion=new",
		"href": map[string]interface{}{
			"uri": "/meshmap",
		},
	})
	if err != nil {
		t.Fatalf("InstallExtension returned error: %v", err)
	}

	if len(provider.Extensions.Navigator) != 1 {
		t.Fatalf("expected navigator extension to be replaced in place, got %d entries", len(provider.Extensions.Navigator))
	}
	if provider.Extensions.Navigator[0].Component != "/provider/navigator/meshmap/index.js?packageVersion=new" {
		t.Fatalf("expected replacement component to be stored, got %q", provider.Extensions.Navigator[0].Component)
	}
}
