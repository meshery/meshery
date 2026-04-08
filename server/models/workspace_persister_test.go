package models

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	schemacore "github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/workspace"
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
		Metadata:       schemacore.Map{"source": "test", "mode": "gorm"},
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
		Metadata:       schemacore.Map{},
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
