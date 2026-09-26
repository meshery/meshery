package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/catalog/v1alpha1"
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

func createWorkspaceTestDesign(t *testing.T, dbHandler *database.Handler, name string, now time.Time) core.Uuid {
	t.Helper()

	designID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate design id: %v", err)
	}

	// Stored the way the pattern persister stores it: a YAML design body.
	design := MesheryPattern{
		ID:         &designID,
		Name:       name,
		Visibility: Private,
		PatternFile: "id: " + designID.String() + "\n" +
			"name: " + name + "\n" +
			"schemaVersion: designs.meshery.io/v1beta1\n" +
			"version: 0.0.1\n" +
			"components: []\n" +
			"relationships: []\n",
		CreatedAt: &now,
		UpdatedAt: &now,
	}
	if err := dbHandler.Create(&design).Error; err != nil {
		t.Fatalf("failed to persist design %q: %v", name, err)
	}

	return designID
}

// Regression for meshery/meshery#21948: every non-empty page failed to convert
// patternFile, so the endpoint only ever succeeded for an empty workspace.
func TestWorkspacePersisterGetWorkspaceDesigns_ReturnsNonEmptyPage(t *testing.T) {
	dbHandler := newWorkspaceTestDB(t)
	if err := dbHandler.AutoMigrate(&workspace.WorkspacesDesignsMapping{}, &MesheryPattern{}); err != nil {
		t.Fatalf("failed to auto-migrate: %v", err)
	}

	workspaceID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate workspace id: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	assignedID := createWorkspaceTestDesign(t, dbHandler, "assigned-design", now)
	unassignedID := createWorkspaceTestDesign(t, dbHandler, "unassigned-design", now)

	persister := &WorkspacePersister{DB: dbHandler}
	if _, err := persister.AddDesignToWorkspace(workspaceID, assignedID); err != nil {
		t.Fatalf("failed to assign design: %v", err)
	}

	tests := []struct {
		name     string
		filter   string
		wantID   core.Uuid
		wantName string
	}{
		{name: "assigned", filter: `{"assigned":true}`, wantID: assignedID, wantName: "assigned-design"},
		{name: "unassigned", filter: `{"assigned":false}`, wantID: unassignedID, wantName: "unassigned-design"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := persister.GetWorkspaceDesigns(workspaceID, "", "", "0", "10", tt.filter, nil)
			if err != nil {
				t.Fatalf("expected designs page, got error: %v", err)
			}

			var page workspace.MesheryDesignPage
			if err := json.Unmarshal(raw, &page); err != nil {
				t.Fatalf("failed to decode designs page: %v", err)
			}

			if page.TotalCount != 1 || len(page.Designs) != 1 {
				t.Fatalf("expected exactly one design, got totalCount=%d designs=%d", page.TotalCount, len(page.Designs))
			}

			got := page.Designs[0]
			if got.ID != tt.wantID {
				t.Errorf("expected design id %s, got %s", tt.wantID, got.ID)
			}
			if got.Name != tt.wantName {
				t.Errorf("expected design name %q, got %q", tt.wantName, got.Name)
			}
			if !got.CreatedAt.Equal(now) || !got.UpdatedAt.Equal(now) {
				t.Errorf("expected timestamps %v, got created_at=%v updated_at=%v", now, got.CreatedAt, got.UpdatedAt)
			}
			if got.UserId != LocalProviderUserID {
				t.Errorf("expected local provider owner %s, got %s", LocalProviderUserID, got.UserId)
			}
			if got.PatternFile == nil || got.PatternFile.Name != tt.wantName {
				t.Errorf("expected decoded patternFile named %q, got %+v", tt.wantName, got.PatternFile)
			}
		})
	}
}

func TestSchemaMesheryPatterns_MapsFieldsAndToleratesUnreadableBody(t *testing.T) {
	designID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate design id: %v", err)
	}
	now := time.Now().UTC().Round(time.Second)

	patterns := []*MesheryPattern{
		{
			ID:          &designID,
			Name:        "catalog-design",
			PatternFile: "name: [unterminated",
			Location:    map[string]interface{}{"type": "local", "port": 9081},
			CatalogData: v1alpha1.CatalogData{
				PatternCaveats: "caveats",
				PatternInfo:    "info",
				SnapshotURL:    []string{"https://example.com/snapshot.png"},
			},
			CreatedAt: &now,
		},
		nil,
		{Name: "empty-body"},
	}

	got := schemaMesheryPatterns(patterns)
	if len(got) != 2 {
		t.Fatalf("expected nil entries to be skipped, got %d designs", len(got))
	}

	for _, d := range got {
		if d.PatternFile != nil {
			t.Errorf("expected unreadable or empty body on %q to be omitted, got %+v", d.Name, d.PatternFile)
		}
	}

	first := got[0]
	if first.ID != designID || !first.CreatedAt.Equal(now) {
		t.Errorf("expected id %s and created_at %v, got %s and %v", designID, now, first.ID, first.CreatedAt)
	}
	if first.Location["type"] != "local" || first.Location["port"] != "9081" {
		t.Errorf("expected location to be carried as strings, got %v", first.Location)
	}
	if first.CatalogData == nil || first.CatalogData.PatternCaveats != "caveats" || first.CatalogData.PatternInfo != "info" {
		t.Fatalf("expected catalog caveats and info to be carried, got %+v", first.CatalogData)
	}
	if first.CatalogData.SnapshotURL == nil || len(*first.CatalogData.SnapshotURL) != 1 {
		t.Errorf("expected snapshot URL to be carried, got %v", first.CatalogData.SnapshotURL)
	}
}
