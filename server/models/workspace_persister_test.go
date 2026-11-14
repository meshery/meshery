package models

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
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

func TestWorkspacePersisterAddEnvironmentToWorkspace_Deduplication(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.WorkspacesEnvironmentsMapping{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	wsID := core.Uuid(uuid.Must(uuid.NewV4()))
	envID := core.Uuid(uuid.Must(uuid.NewV4()))

	persister := &WorkspacePersister{DB: dbHandler}
	_, err := persister.AddEnvironmentToWorkspace(wsID, envID)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	_, err = persister.AddEnvironmentToWorkspace(wsID, envID)
	if err != nil {
		t.Fatalf("second add failed: %v", err)
	}

	var count int64
	dbHandler.Model(&workspace.WorkspacesEnvironmentsMapping{}).Where("workspace_id = ? AND environment_id = ?", wsID, envID).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 mapping after adding twice, got %d", count)
	}
}

func TestWorkspacePersisterAddViewToWorkspace_Deduplication(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.WorkspacesViewsMapping{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	wsID := core.Uuid(uuid.Must(uuid.NewV4()))
	viewID := core.Uuid(uuid.Must(uuid.NewV4()))

	persister := &WorkspacePersister{DB: dbHandler}
	_, err := persister.AddViewToWorkspace(wsID, viewID)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	_, err = persister.AddViewToWorkspace(wsID, viewID)
	if err != nil {
		t.Fatalf("second add failed: %v", err)
	}

	var count int64
	dbHandler.Model(&workspace.WorkspacesViewsMapping{}).Where("workspace_id = ? AND view_id = ?", wsID, viewID).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 mapping after adding twice, got %d", count)
	}
}

func TestWorkspacePersisterAddTeamToWorkspace_Deduplication(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.WorkspacesTeamsMapping{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	wsID := core.Uuid(uuid.Must(uuid.NewV4()))
	teamID := core.Uuid(uuid.Must(uuid.NewV4()))

	persister := &WorkspacePersister{DB: dbHandler}
	_, err := persister.AddTeamToWorkspace(wsID, teamID)
	if err != nil {
		t.Fatalf("first add failed: %v", err)
	}

	_, err = persister.AddTeamToWorkspace(wsID, teamID)
	if err != nil {
		t.Fatalf("second add failed: %v", err)
	}

	var count int64
	dbHandler.Model(&workspace.WorkspacesTeamsMapping{}).Where("workspace_id = ? AND team_id = ?", wsID, teamID).Count(&count)
	if count != 1 {
		t.Fatalf("expected 1 mapping after adding twice, got %d", count)
	}
}

func TestWorkspacePersisterGetWorkspaces_InvalidPagination(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.Workspace{}); err != nil {
		t.Fatalf("failed to auto-migrate schemas workspace: %v", err)
	}
	persister := &WorkspacePersister{DB: dbHandler}
	orgID := uuid.Must(uuid.NewV4()).String()

	tests := []struct {
		name     string
		page     string
		pageSize string
		wantCode string
	}{
		{"invalid page", "-1", "10", ErrPageNumberCode},
		{"invalid pagesize", "0", "abc", ErrPageSizeCode},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := persister.GetWorkspaces(orgID, "", "", tc.page, tc.pageSize, "")
			if err == nil {
				t.Fatalf("expected an error for page=%q pagesize=%q, got nil", tc.page, tc.pageSize)
			}
			if got := meshkiterrors.GetCode(err); got != tc.wantCode {
				t.Errorf("error code = %q, want %q (err: %v)", got, tc.wantCode, err)
			}
		})
	}
}
