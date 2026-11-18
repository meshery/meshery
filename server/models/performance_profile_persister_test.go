package models

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
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

// TestDeletePerformanceProfileRemovesStoredProfile pins the success path: a
// stored profile is deleted and the persister reports the id it removed.
func TestDeletePerformanceProfileRemovesStoredProfile(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	profile := &PerformanceProfile{ID: mustUUID(t), Name: "linkerd-10s", Owner: mustUUID(t), Duration: "10s"}
	if err := ppp.SavePerformanceProfile(profile.ID, profile); err != nil {
		t.Fatalf("SavePerformanceProfile: %v", err)
	}

	raw, err := ppp.DeletePerformanceProfile(profile.ID)
	if err != nil {
		t.Fatalf("DeletePerformanceProfile: %v", err)
	}

	var deleted PerformanceProfile
	if err := json.Unmarshal(raw, &deleted); err != nil {
		t.Fatalf("unmarshal deleted profile: %v", err)
	}
	if deleted.ID != profile.ID {
		t.Errorf("deleted id = %s, want %s", deleted.ID, profile.ID)
	}

	page, err := ppp.GetPerformanceProfiles("", "", "", 0, 10)
	if err != nil {
		t.Fatalf("GetPerformanceProfiles: %v", err)
	}
	var after PerformanceProfilePage
	if err := json.Unmarshal(page, &after); err != nil {
		t.Fatalf("unmarshal page: %v", err)
	}
	if after.TotalCount != 0 {
		t.Errorf("profile still stored after delete: totalCount=%d", after.TotalCount)
	}
}

// TestDeletePerformanceProfileReportsMissingProfile is the regression test for
// a silent success: deleting an id that was never stored affects no row, which
// gorm does not treat as an error, so the persister marshalled the stub profile
// it had just built (carrying only the caller's id) and the handler answered 200.
// A client could not tell "deleted" from "never existed". It must be not-found.
func TestDeletePerformanceProfileReportsMissingProfile(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	raw, err := ppp.DeletePerformanceProfile(mustUUID(t))
	if err == nil {
		t.Fatalf("expected a not-found error, got body %s", raw)
	}
	if got := mkerrors.GetCode(err); got != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", got, ErrResultNotFoundCode, err)
	}
	if raw != nil {
		t.Errorf("expected no body for a missing profile, got %s", raw)
	}
}

// TestPerformanceProfilePersister_NilGuards verifies that calling persister methods
// on a nil receiver or an uninitialized DB returns ErrDBConnection rather than panicking.
func TestPerformanceProfilePersister_NilGuards(t *testing.T) {
	cases := []struct {
		name string
		ppp  *PerformanceProfilePersister
	}{
		{name: "nil receiver", ppp: nil},
		{name: "nil db handler", ppp: &PerformanceProfilePersister{DB: nil}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.ppp.GetPerformanceProfiles("", "", "", 0, 10); !errors.Is(err, ErrDBConnection) {
				t.Errorf("GetPerformanceProfiles: expected ErrDBConnection, got %v", err)
			}
			if got, err := tc.ppp.GetPerformanceProfile(mustUUID(t)); !errors.Is(err, ErrDBConnection) || got != nil {
				t.Errorf("GetPerformanceProfile: expected (nil, ErrDBConnection), got (%v, %v)", got, err)
			}
			if _, err := tc.ppp.DeletePerformanceProfile(mustUUID(t)); !errors.Is(err, ErrDBConnection) {
				t.Errorf("DeletePerformanceProfile: expected ErrDBConnection, got %v", err)
			}
			if err := tc.ppp.SavePerformanceProfile(mustUUID(t), &PerformanceProfile{}); !errors.Is(err, ErrDBConnection) {
				t.Errorf("SavePerformanceProfile: expected ErrDBConnection, got %v", err)
			}
		})
	}
}

// TestGetPerformanceProfile_NotFound verifies that requesting a non-existent profile ID
// returns nil (not a pointer to an empty struct) and an error with ErrResultNotFoundCode.
func TestGetPerformanceProfile_NotFound(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	id := mustUUID(t)
	got, err := ppp.GetPerformanceProfile(id)
	if got != nil {
		t.Errorf("expected nil profile for missing id, got %+v", got)
	}
	if err == nil {
		t.Fatalf("expected ErrResultNotFound, got nil error")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrResultNotFoundCode, err)
	}
}

// TestGetPerformanceProfile_DBError verifies that when a database read fails,
// ErrDBRead is returned with a nil profile pointer.
func TestGetPerformanceProfile_DBError(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	// Deliberately un-migrated: performance_profiles table does not exist.
	ppp := &PerformanceProfilePersister{DB: &db}

	got, err := ppp.GetPerformanceProfile(mustUUID(t))
	if got != nil {
		t.Errorf("expected nil profile on DB error, got %+v", got)
	}
	if err == nil {
		t.Fatalf("expected ErrDBRead on unmigrated table, got nil error")
	}
	if code := mkerrors.GetCode(err); code != ErrDBReadCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBReadCode, err)
	}
}

// TestDeletePerformanceProfile_DBError verifies that when a database delete fails,
// ErrDBDelete is returned with nil body.
func TestDeletePerformanceProfile_DBError(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	// Deliberately un-migrated: performance_profiles table does not exist.
	ppp := &PerformanceProfilePersister{DB: &db}

	raw, err := ppp.DeletePerformanceProfile(mustUUID(t))
	if raw != nil {
		t.Errorf("expected nil raw body on DB error, got %s", raw)
	}
	if err == nil {
		t.Fatalf("expected ErrDBDelete on unmigrated table, got nil error")
	}
	if code := mkerrors.GetCode(err); code != ErrDBDeleteCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBDeleteCode, err)
	}
}

// TestSavePerformanceProfile_NilGuard verifies that passing a nil profile to
// SavePerformanceProfile returns an ErrInvalidPerformanceProfile coded error
// rather than panicking or returning an uncoded error.
func TestSavePerformanceProfile_NilGuard(t *testing.T) {
	db := newMigratedDB(t)
	ppp := &PerformanceProfilePersister{DB: db}

	err := ppp.SavePerformanceProfile(mustUUID(t), nil)
	if err == nil {
		t.Fatalf("expected error when saving nil profile, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrInvalidPerformanceProfileCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrInvalidPerformanceProfileCode, err)
	}
}
