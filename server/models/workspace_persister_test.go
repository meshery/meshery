package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	environment "github.com/meshery/schemas/models/v1beta3/environment"
	workspace "github.com/meshery/schemas/models/v1beta3/workspace"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newWorkspaceTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	return &database.Handler{DB: db}
}

func TestSchemasWorkspaceAutoMigrateAndPersistMetadata(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.Workspace{}); err != nil {
		t.Fatalf("failed to auto-migrate schemas workspace: %v", err)
	}

	workspaceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate workspace id: %v", err)
	}

	organizationID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate organization id: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	ws := workspace.Workspace{
		ID:             workspaceID,
		CreatedAt:      now,
		UpdatedAt:      now,
		Name:           "Regression Workspace",
		Description:    "verifies schemas workspace metadata storage",
		Metadata:       core.Map{"source": "test", "mode": "gorm"},
		OrganizationID: organizationID,
	}

	if err := dbHandler.Create(&ws).Error; err != nil {
		t.Fatalf("failed to persist schemas workspace: %v", err)
	}

	stored := workspace.Workspace{}
	if err := dbHandler.First(&stored, "id = ?", workspaceID).Error; err != nil {
		t.Fatalf("failed to read persisted schemas workspace: %v", err)
	}

	if got := stored.Metadata["source"]; got != "test" {
		t.Fatalf("unexpected metadata[source]: got %q, want %q", got, "test")
	}

	if got := stored.Metadata["mode"]; got != "gorm" {
		t.Fatalf("unexpected metadata[mode]: got %q, want %q", got, "gorm")
	}
}

func TestWorkspacePersisterUpdateWorkspace_PreservesOrganizationIDWhenOmitted(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.Workspace{}); err != nil {
		t.Fatalf("failed to auto-migrate schemas workspace: %v", err)
	}

	workspaceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate workspace id: %v", err)
	}

	organizationID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate organization id: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	ws := workspace.Workspace{
		ID:             workspaceID,
		CreatedAt:      now,
		UpdatedAt:      now,
		Name:           "Original Workspace",
		Description:    "before update",
		Metadata:       core.Map{},
		OrganizationID: organizationID,
	}

	if err := dbHandler.Create(&ws).Error; err != nil {
		t.Fatalf("failed to persist schemas workspace: %v", err)
	}

	persister := &WorkspacePersister{DB: dbHandler}
	updated, err := persister.UpdateWorkspace(workspaceID, &workspace.WorkspaceUpdatePayload{
		Name: "Updated Workspace",
	})
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	if updated.Name != "Updated Workspace" {
		t.Fatalf("expected updated name, got %q", updated.Name)
	}
	if updated.OrganizationID != organizationID {
		t.Fatalf("expected organization ID %s to be preserved, got %s", organizationID, updated.OrganizationID)
	}
}

func TestDefaultLocalProviderUpdateWorkspace_ReturnsErrorForInvalidUUID(t *testing.T) {
	provider := &DefaultLocalProvider{}

	_, err := provider.UpdateWorkspace(nil, &workspace.WorkspaceUpdatePayload{}, "not-a-uuid")
	if err == nil {
		t.Fatal("expected invalid workspace ID to return an error")
	}
}

// TestGetWorkspaceEnvironments_SoftDeletedMappingReturnsToUnassigned is a
// regression test for https://github.com/meshery/meshery/issues/21947.
//
// Removing an environment from a workspace soft-deletes the mapping row
// (deleted_at is set) rather than removing it. The "assigned" query already
// ignores soft-deleted mappings via `... AND wem.deleted_at IS NULL` inside its
// EXISTS clause, but the "unassigned" query's LEFT JOIN previously omitted that
// predicate. The join therefore matched the soft-deleted mapping,
// `wem.workspace_id IS NULL` evaluated false, and the environment vanished from
// BOTH the assigned and the unassigned lists — leaving it unreachable from the
// UI. This test pins the fixed behavior: once the mapping is soft-deleted the
// environment must reappear in the unassigned list and stay out of the assigned
// list.
func TestGetWorkspaceEnvironments_SoftDeletedMappingReturnsToUnassigned(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(
		&workspace.Workspace{},
		&workspace.WorkspacesEnvironmentsMapping{},
		&environment.Environment{},
	); err != nil {
		t.Fatalf("failed to auto-migrate schemas models: %v", err)
	}

	persister := &WorkspacePersister{DB: dbHandler}

	workspaceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate workspace id: %v", err)
	}
	environmentID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate environment id: %v", err)
	}

	env := environment.Environment{ID: environmentID, Name: "Regression Environment"}
	if err := dbHandler.Create(&env).Error; err != nil {
		t.Fatalf("failed to persist environment: %v", err)
	}

	// Assign, then remove (soft-delete) the environment mapping.
	if _, err := persister.AddEnvironmentToWorkspace(workspaceID, environmentID); err != nil {
		t.Fatalf("failed to add environment to workspace: %v", err)
	}
	if _, err := persister.DeleteEnvironmentFromWorkspace(workspaceID, environmentID); err != nil {
		t.Fatalf("failed to remove environment from workspace: %v", err)
	}

	// After the mapping is soft-deleted, the environment must be selectable
	// again from the unassigned list. This is the regression being guarded:
	// before the fix the soft-deleted mapping matched the LEFT JOIN and hid the
	// environment here.
	if got := environmentIDsForWorkspace(t, persister, workspaceID, `{"assigned":false}`); !containsUUID(got, environmentID) {
		t.Fatalf("expected soft-deleted environment %s to appear in the unassigned list, got %v", environmentID, got)
	}

	// ...and it must not appear in the assigned list.
	if got := environmentIDsForWorkspace(t, persister, workspaceID, `{"assigned":true}`); containsUUID(got, environmentID) {
		t.Fatalf("expected soft-deleted environment %s to be absent from the assigned list, got %v", environmentID, got)
	}
}

func environmentIDsForWorkspace(t *testing.T, persister *WorkspacePersister, workspaceID core.Uuid, filter string) []core.Uuid {
	t.Helper()

	raw, err := persister.GetWorkspaceEnvironments(workspaceID, "", "", "0", "all", filter)
	if err != nil {
		t.Fatalf("GetWorkspaceEnvironments(filter=%s) failed: %v", filter, err)
	}

	var page environment.EnvironmentPage
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("failed to unmarshal environment page: %v", err)
	}

	ids := make([]core.Uuid, 0, len(page.Environments))
	for _, e := range page.Environments {
		ids = append(ids, e.ID)
	}
	return ids
}

func containsUUID(ids []core.Uuid, target core.Uuid) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}
