package models

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
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

func TestDefaultLocalProviderWorkspaceMethods_ReturnErrorForInvalidUUIDs(t *testing.T) {
	provider := &DefaultLocalProvider{}
	validID := "11111111-1111-1111-1111-111111111111"
	invalidID := "not-a-uuid"

	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "get workspace rejects invalid workspace ID",
			call: func() error {
				_, err := provider.GetWorkspaceByID(nil, invalidID, "")
				return err
			},
		},
		{
			name: "delete workspace rejects invalid workspace ID",
			call: func() error {
				_, err := provider.DeleteWorkspace(nil, invalidID)
				return err
			},
		},
		{
			name: "add environment rejects invalid workspace ID",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(nil, invalidID, validID)
				return err
			},
		},
		{
			name: "add environment rejects invalid environment ID",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(nil, validID, invalidID)
				return err
			},
		},
		{
			name: "remove environment rejects invalid workspace ID",
			call: func() error {
				_, err := provider.RemoveEnvironmentFromWorkspace(nil, invalidID, validID)
				return err
			},
		},
		{
			name: "remove environment rejects invalid environment ID",
			call: func() error {
				_, err := provider.RemoveEnvironmentFromWorkspace(nil, validID, invalidID)
				return err
			},
		},
		{
			name: "get workspace environments rejects invalid workspace ID",
			call: func() error {
				_, err := provider.GetEnvironmentsOfWorkspace(nil, invalidID, "", "", "", "", "")
				return err
			},
		},
		{
			name: "add design rejects invalid workspace ID",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(nil, invalidID, validID)
				return err
			},
		},
		{
			name: "add design rejects invalid design ID",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(nil, validID, invalidID)
				return err
			},
		},
		{
			name: "get workspace designs rejects invalid workspace ID",
			call: func() error {
				_, err := provider.GetDesignsOfWorkspace(nil, invalidID, "", "", "", "", "", nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil {
				t.Fatal("expected invalid UUID to return an error")
			}
		})
	}
}
