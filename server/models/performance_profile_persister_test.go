package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

// newMigratedDB returns an in-memory database carrying exactly the schema boot
// creates, so these tests read the same columns a running server does.
func newMigratedDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := AutoMigrateSystemTables(&db); err != nil {
		t.Fatalf("AutoMigrateSystemTables: %v", err)
	}
	return &db
}

func mustUUID(t *testing.T) core.Uuid {
	t.Helper()

	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("new uuid: %v", err)
	}
	return core.Uuid(id)
}

// TestGetPerformanceProfilesReadsMigratedColumns is the regression test for a
// silent total failure of performance-profile listing on the local provider.
//
// PerformanceProfile is the schemas type, and schemas renamed its owner field
// from UserID (`db:"user_id"`) to Owner (`db:"owner"`). AutoMigrate therefore
// creates an `owner` column, but the SELECT here still named `user_id`. Every
// read failed with "no such column: user_id", and because the gorm errors were
// discarded the persister answered HTTP 200 with an empty list - the UI and
// `mesheryctl perf profile` showed zero profiles no matter how many were saved,
// and nothing was logged.
//
// A saved profile MUST come back, and its owner MUST survive the round trip.
func TestGetPerformanceProfilesReadsMigratedColumns(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	owner := mustUUID(t)
	profile := &PerformanceProfile{
		ID:       mustUUID(t),
		Name:     "istio-30s",
		Owner:    owner,
		Duration: "30s",
	}
	if err := ppp.SavePerformanceProfile(profile.ID, profile); err != nil {
		t.Fatalf("SavePerformanceProfile: %v", err)
	}

	raw, err := ppp.GetPerformanceProfiles("", "", "", 0, 10)
	if err != nil {
		t.Fatalf("GetPerformanceProfiles: %v", err)
	}

	var page PerformanceProfilePage
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("unmarshal page: %v", err)
	}

	if page.TotalCount != 1 || len(page.Profiles) != 1 {
		t.Fatalf("saved profile did not come back: totalCount=%d profiles=%d body=%s",
			page.TotalCount, len(page.Profiles), raw)
	}
	if got := page.Profiles[0].Name; got != "istio-30s" {
		t.Errorf("name = %q, want %q", got, "istio-30s")
	}
	if got := page.Profiles[0].Owner; got != owner {
		t.Errorf("owner = %q, want %q; the owner column is not being selected", got, owner)
	}
}

// TestGetPerformanceProfileReadsMigratedColumns covers the single-profile read,
// which carried the same pre-rename column name. Unlike the list path it does
// surface the error, so this one 500'd rather than lying - both are fixed by
// selecting the column AutoMigrate actually creates.
func TestGetPerformanceProfileReadsMigratedColumns(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	owner := mustUUID(t)
	profile := &PerformanceProfile{
		ID:       mustUUID(t),
		Name:     "consul-15s",
		Owner:    owner,
		Duration: "15s",
	}
	if err := ppp.SavePerformanceProfile(profile.ID, profile); err != nil {
		t.Fatalf("SavePerformanceProfile: %v", err)
	}

	got, err := ppp.GetPerformanceProfile(profile.ID)
	if err != nil {
		t.Fatalf("GetPerformanceProfile: %v", err)
	}
	if got.Name != "consul-15s" {
		t.Errorf("name = %q, want %q", got.Name, "consul-15s")
	}
	if got.Owner != owner {
		t.Errorf("owner = %q, want %q; the owner column is not being selected", got.Owner, owner)
	}
}

// TestGetPerformanceProfilesPropagatesQueryError pins the second half of the
// defect: a failing query must be reported, not flattened into an empty page.
// Dropping the gorm errors is what let the column rename go unnoticed.
func TestGetPerformanceProfilesPropagatesQueryError(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	// Deliberately un-migrated: performance_profiles does not exist.
	ppp := &PerformanceProfilePersister{DB: &db}

	raw, err := ppp.GetPerformanceProfiles("", "", "", 0, 10)
	if err == nil {
		t.Fatalf("expected an error querying a missing table, got body %s", raw)
	}
}
