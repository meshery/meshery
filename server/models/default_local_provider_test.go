package models

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/database"
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
