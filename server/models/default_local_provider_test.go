package models

import (
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
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

// TestGetUserCredentialsRejectsUnsanitizedOrderInput verifies that GetUserCredentials
// no longer passes the caller-supplied "order" query parameter straight into gorm's
// Order(), which treats a plain string as raw, unescaped SQL. Without sanitization,
// arbitrary attacker-controlled text (e.g. "foobar asc", or a crafted SQL injection
// payload) reaches the database driver directly -- an unpatched sibling of the
// SQL injection vulnerabilities already fixed elsewhere (CVE-2024-35181, CVE-2024-35182)
// via models.SanitizeOrderInput, which this endpoint had been missed by.
func TestGetUserCredentialsRejectsUnsanitizedOrderInput(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)

	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	for _, name := range []string{"b-cred", "a-cred"} {
		cred := &Credential{Name: name, Type: "token", UserId: userID}
		if _, err := provider.SaveUserCredential("tok", cred); err != nil {
			t.Fatalf("failed to seed credential %q: %v", name, err)
		}
	}

	t.Run("rejects unsanitized order input", func(t *testing.T) {
		// "foobar" is not a real column. Pre-fix, this string was handed straight to
		// gorm's Order(), which built `ORDER BY foobar asc` and the query failed with
		// a driver-level "no such column" error -- proof the raw input reached SQL
		// construction unsanitized. Post-fix, SanitizeOrderInput rejects it outright
		// and Order() becomes a no-op, so the query must still succeed.
		page, err := provider.GetUserCredentials(nil, userID.String(), 0, 10, "", "foobar asc")
		if err != nil {
			t.Fatalf("unsanitized order value reached the query and caused an error: %v", err)
		}
		if page.TotalCount != 2 {
			t.Errorf("got TotalCount=%d, want 2", page.TotalCount)
		}
	})

	t.Run("accepts legitimate order input", func(t *testing.T) {
		// A legitimate, allowlisted order value must still work after sanitization.
		page, err := provider.GetUserCredentials(nil, userID.String(), 0, 10, "", "name asc")
		if err != nil {
			t.Fatalf("unexpected error with legitimate order value: %v", err)
		}
		if len(page.Credentials) != 2 || page.Credentials[0].Name != "a-cred" || page.Credentials[1].Name != "b-cred" {
			t.Fatalf("got order %v, want [a-cred b-cred]", credentialNames(page.Credentials))
		}
	})
}

func credentialNames(creds []Credential) []string {
	names := make([]string, len(creds))
	for i, c := range creds {
		names[i] = c.Name
	}
	return names
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

func TestDefaultLocalProviderSaveK8sContextPropagatesPersistenceFailure(t *testing.T) {
	newProvider := func(t *testing.T) (*DefaultLocalProvider, *database.Handler) {
		t.Helper()

		connectionDB, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
		if err != nil {
			t.Fatalf("open connection database: %v", err)
		}
		if err := connectionDB.AutoMigrate(&connections.Connection{}); err != nil {
			t.Fatalf("migrate connections: %v", err)
		}

		contextDB, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
		if err != nil {
			t.Fatalf("open context database: %v", err)
		}
		if err := contextDB.AutoMigrate(&K8sContext{}); err != nil {
			t.Fatalf("migrate Kubernetes contexts: %v", err)
		}

		return &DefaultLocalProvider{
			ConnectionPersister:        &ConnectionPersister{DB: &connectionDB},
			MesheryK8sContextPersister: &MesheryK8sContextPersister{DB: &contextDB},
		}, &contextDB
	}

	newContext := func() K8sContext {
		instanceID := uuid.Must(uuid.NewV4())
		return K8sContext{
			Name:              "test-context",
			Server:            "https://cluster.example.test",
			MesheryInstanceID: &instanceID,
			Auth:              map[string]interface{}{},
			Cluster:           map[string]interface{}{"server": "https://cluster.example.test"},
		}
	}

	t.Run("context persistence failure", func(t *testing.T) {
		provider, contextDB := newProvider(t)
		if err := contextDB.DBClose(); err != nil {
			t.Fatalf("close context database: %v", err)
		}

		if _, err := provider.SaveK8sContext("", newContext(), nil); err == nil {
			t.Fatal("SaveK8sContext returned nil after context persistence failed")
		}
	})

	t.Run("success", func(t *testing.T) {
		provider, _ := newProvider(t)
		connection, err := provider.SaveK8sContext("", newContext(), nil)
		if err != nil {
			t.Fatalf("SaveK8sContext returned error: %v", err)
		}
		if connection.ID == uuid.Nil {
			t.Fatal("SaveK8sContext returned a connection without an ID")
		}
	})
}
