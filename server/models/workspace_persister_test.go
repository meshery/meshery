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

func TestSchemasWorkspaceAutoMigrateAndPersistMetadata(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	dbHandler := &database.Handler{DB: db}
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
		Metadata:       schemacore.MapObject{"source": "test", "mode": "gorm"},
		OrganizationID: &organizationID,
		Owner:          "meshery",
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