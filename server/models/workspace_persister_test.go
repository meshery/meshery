package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	schemacore "github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/organization"
	viewv1beta1 "github.com/meshery/schemas/models/v1beta1/view"
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

func TestWorkspacePersisterGetWorkspaceViews_EnrichesLocation(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(
		&organization.Organization{},
		&workspace.Workspace{},
		&viewv1beta1.MesheryView{},
		&workspace.WorkspacesViewsMapping{},
	); err != nil {
		t.Fatalf("failed to auto-migrate schemas for workspace views: %v", err)
	}

	organizationID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate organization id: %v", err)
	}

	workspaceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate workspace id: %v", err)
	}

	viewID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate view id: %v", err)
	}

	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	mappingID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate mapping id: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	org := organization.Organization{
		ID:          organizationID,
		Name:        "Test Organization",
		Description: "organization for workspace views test",
		Country:     "US",
		Region:      "NA",
		Owner:       userID,
		Metadata:    organization.OrgMetadata{},
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := dbHandler.Create(&org).Error; err != nil {
		t.Fatalf("failed to persist organization: %v", err)
	}

	ws := workspace.Workspace{
		ID:             workspaceID,
		CreatedAt:      now,
		UpdatedAt:      now,
		Name:           "Test Workspace",
		Description:    "workspace for view enrichment test",
		Metadata:       schemacore.Map{},
		OrganizationID: organizationID,
	}
	if err := dbHandler.Create(&ws).Error; err != nil {
		t.Fatalf("failed to persist workspace: %v", err)
	}

	mesheryView := viewv1beta1.MesheryView{
		ID:         viewID,
		Name:       "Test View",
		Visibility: "private",
		Filters:    schemacore.Map{"kind": "deployment"},
		Metadata:   schemacore.Map{"source": "test"},
		UserID:     userID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := dbHandler.Table("meshery_views").Create(&mesheryView).Error; err != nil {
		t.Fatalf("failed to persist view: %v", err)
	}

	mapping := workspace.WorkspacesViewsMapping{
		ID:          mappingID,
		ViewId:      viewID,
		WorkspaceId: workspaceID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := dbHandler.Create(&mapping).Error; err != nil {
		t.Fatalf("failed to persist workspace-view mapping: %v", err)
	}

	persister := &WorkspacePersister{DB: dbHandler}
	resp, err := persister.GetWorkspaceViews(workspaceID, "", "", "0", "10", "")
	if err != nil {
		t.Fatalf("expected GetWorkspaceViews to succeed, got %v", err)
	}

	var page viewv1beta1.MesheryViewPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("failed to unmarshal workspace views page: %v", err)
	}

	if len(page.Views) != 1 {
		t.Fatalf("expected 1 view, got %d", len(page.Views))
	}

	got := page.Views[0]
	if got.ID != viewID {
		t.Fatalf("expected view id %s, got %s", viewID, got.ID)
	}
	if got.WorkspaceID != workspaceID {
		t.Fatalf("expected workspace id %s, got %s", workspaceID, got.WorkspaceID)
	}
	if got.WorkspaceName != ws.Name {
		t.Fatalf("expected workspace name %q, got %q", ws.Name, got.WorkspaceName)
	}
	if got.OrganizationID != organizationID {
		t.Fatalf("expected organization id %s, got %s", organizationID, got.OrganizationID)
	}
	if got.OrganizationName != org.Name {
		t.Fatalf("expected organization name %q, got %q", org.Name, got.OrganizationName)
	}
}
