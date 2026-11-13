package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta3/environment"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestLocalProviderWithEnvironmentDB(t *testing.T) *DefaultLocalProvider {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	dbHandler := &database.Handler{DB: db}
	if err := AutoMigrateSystemTables(dbHandler); err != nil {
		t.Fatalf("failed to auto-migrate system tables: %v", err)
	}

	envPersister := &EnvironmentPersister{DB: dbHandler}
	return &DefaultLocalProvider{
		GenericPersister:     dbHandler,
		EnvironmentPersister: envPersister,
	}
}

// TestDefaultLocalProvider_SaveEnvironmentDefaultsToUserPurpose tests that
// SaveEnvironment always creates environments with Purpose "user" and that
// IsAdministrative() evaluates to false.
func TestDefaultLocalProvider_SaveEnvironmentDefaultsToUserPurpose(t *testing.T) {
	provider := newTestLocalProviderWithEnvironmentDB(t)

	orgID := uuid.Must(uuid.NewV4())
	payload := &environment.EnvironmentPayload{
		Name:        "dev-env",
		Description: "test environment description",
		OrgID:       orgID,
	}

	data, err := provider.SaveEnvironment(nil, payload, "", false)
	if err != nil {
		t.Fatalf("SaveEnvironment failed: %v", err)
	}

	var saved environment.Environment
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("unmarshal saved environment: %v", err)
	}

	if saved.Purpose != environment.EnvironmentPurposeUser {
		t.Fatalf("expected Purpose to be %q, got %q", environment.EnvironmentPurposeUser, saved.Purpose)
	}
	if saved.IsAdministrative() {
		t.Fatal("expected IsAdministrative() to be false for saved environment")
	}
	if got := saved.EffectivePurpose(); got != environment.EnvironmentPurposeUser {
		t.Fatalf("expected EffectivePurpose() to be %q, got %q", environment.EnvironmentPurposeUser, got)
	}
}

// TestDefaultLocalProvider_UpdateEnvironmentPreservesPurpose asserts that
// UpdateEnvironment carries forward the stored Purpose value rather than overwriting
// it with a zero value.
func TestDefaultLocalProvider_UpdateEnvironmentPreservesPurpose(t *testing.T) {
	provider := newTestLocalProviderWithEnvironmentDB(t)

	orgID := uuid.Must(uuid.NewV4())

	// 1. Create an administrative environment directly (simulating provisioning by platform / remote provider)
	adminID := uuid.Must(uuid.NewV4())
	adminEnv := &environment.Environment{
		ID:             adminID,
		Name:           "admin-env",
		OrganizationID: core.Uuid(orgID),
		Purpose:        environment.EnvironmentPurposeAdministrative,
		CreatedAt:      time.Now().Add(-1 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * time.Hour),
	}
	if err := provider.EnvironmentPersister.DB.Create(adminEnv).Error; err != nil {
		t.Fatalf("seed admin environment: %v", err)
	}

	// Update the administrative environment via UpdateEnvironment
	updatePayload := &environment.EnvironmentPayload{
		Name:        "renamed-admin-env",
		Description: "updated description",
		OrgID:       orgID,
	}

	updated, err := provider.UpdateEnvironment(nil, updatePayload, adminID.String())
	if err != nil {
		t.Fatalf("UpdateEnvironment failed: %v", err)
	}

	if updated.Purpose != environment.EnvironmentPurposeAdministrative {
		t.Fatalf("UpdateEnvironment overwrote Purpose! got %q, want %q", updated.Purpose, environment.EnvironmentPurposeAdministrative)
	}
	if !updated.IsAdministrative() {
		t.Fatal("expected updated environment to remain administrative")
	}

	// Verify in database as well
	reloaded, err := provider.EnvironmentPersister.GetEnvironment(core.Uuid(adminID))
	if err != nil {
		t.Fatalf("GetEnvironment: %v", err)
	}
	if reloaded.Purpose != environment.EnvironmentPurposeAdministrative {
		t.Fatalf("database row Purpose was overwritten! got %q, want %q", reloaded.Purpose, environment.EnvironmentPurposeAdministrative)
	}
	if !reloaded.IsAdministrative() {
		t.Fatal("expected reloaded database row to remain administrative")
	}

	// 2. Test that updating an environment with empty Purpose preserves it as empty / user
	legacyID := uuid.Must(uuid.NewV4())
	legacyEnv := &environment.Environment{
		ID:             legacyID,
		Name:           "legacy-env",
		OrganizationID: core.Uuid(orgID),
		Purpose:        "",
		CreatedAt:      time.Now().Add(-2 * time.Hour),
		UpdatedAt:      time.Now().Add(-2 * time.Hour),
	}
	if err := provider.EnvironmentPersister.DB.Create(legacyEnv).Error; err != nil {
		t.Fatalf("seed legacy environment: %v", err)
	}

	legacyUpdated, err := provider.UpdateEnvironment(nil, updatePayload, legacyID.String())
	if err != nil {
		t.Fatalf("UpdateEnvironment for legacy env failed: %v", err)
	}
	if legacyUpdated.IsAdministrative() {
		t.Fatal("legacy environment must not become administrative on update")
	}
	if got := legacyUpdated.EffectivePurpose(); got != environment.EnvironmentPurposeUser {
		t.Fatalf("expected EffectivePurpose() %q, got %q", environment.EnvironmentPurposeUser, got)
	}
}

func TestDefaultLocalProvider_UpdateNonExistentEnvironmentReturnsNotFound(t *testing.T) {
	provider := newTestLocalProviderWithEnvironmentDB(t)
	nonExistentID := uuid.Must(uuid.NewV4()).String()
	payload := &environment.EnvironmentPayload{
		Name:  "updated",
		OrgID: core.Uuid(uuid.Must(uuid.NewV4())),
	}
	_, err := provider.UpdateEnvironment(nil, payload, nonExistentID)
	if err == nil {
		t.Fatal("expected error when updating non-existent environment, got nil")
	}
	if mkerrors.GetCode(err) != ErrResultNotFoundCode {
		t.Fatalf("expected ErrResultNotFoundCode (%s), got: %v", ErrResultNotFoundCode, err)
	}
}
