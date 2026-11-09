package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta3/environment"
	workspace "github.com/meshery/schemas/models/v1beta3/workspace"
)

// TestWorkspaceDeleteCascadesToAllMappings verifies that deleting a workspace
// transactionally deletes all of its mapping rows across environments, views,
// teams, and designs.
func TestWorkspaceDeleteCascadesToAllMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}

	wsID := mustUUID(t)
	ws := &workspace.Workspace{
		ID:   wsID,
		Name: "test-workspace-cascade",
	}
	if _, err := wp.SaveWorkspace(ws); err != nil {
		t.Fatalf("save workspace: %v", err)
	}

	envID, viewID, teamID, designID := mustUUID(t), mustUUID(t), mustUUID(t), mustUUID(t)

	if _, err := wp.AddEnvironmentToWorkspace(wsID, envID); err != nil {
		t.Fatalf("add env: %v", err)
	}
	if _, err := wp.AddViewToWorkspace(wsID, viewID); err != nil {
		t.Fatalf("add view: %v", err)
	}
	if _, err := wp.AddTeamToWorkspace(wsID, teamID); err != nil {
		t.Fatalf("add team: %v", err)
	}
	if _, err := wp.AddDesignToWorkspace(wsID, designID); err != nil {
		t.Fatalf("add design: %v", err)
	}

	if _, err := wp.DeleteWorkspace(ws); err != nil {
		t.Fatalf("delete workspace: %v", err)
	}

	tables := []struct {
		table  string
		column string
		member core.Uuid
	}{
		{"workspaces_environments_mappings", "environment_id", envID},
		{"workspaces_views_mappings", "view_id", viewID},
		{"workspaces_teams_mappings", "team_id", teamID},
		{"workspaces_designs_mappings", "design_id", designID},
	}

	for _, tbl := range tables {
		var count int64
		if err := db.Table(tbl.table).Where("workspace_id = ?", wsID).Count(&count).Error; err != nil {
			t.Fatalf("count %s: %v", tbl.table, err)
		}
		if count != 0 {
			t.Errorf("table %s still has %d orphaned mapping(s) after workspace deletion; want 0", tbl.table, count)
		}
	}
}

// TestWorkspaceRemoveDesignDeletesStoredDuplicates verifies that unassigning a
// design from a workspace removes all duplicate mappings stored previously.
func TestWorkspaceRemoveDesignDeletesStoredDuplicates(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}

	wsID, designID := mustUUID(t), mustUUID(t)

	insert := "INSERT INTO workspaces_designs_mappings (id, workspace_id, design_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	for i := 1; i <= 2; i++ {
		if err := db.Exec(insert, mustUUID(t), wsID, designID, time.Now(), time.Now()).Error; err != nil {
			t.Fatalf("seed duplicate mapping #%d: %v", i, err)
		}
	}

	if _, err := wp.DeleteDesignFromWorkspace(wsID, designID); err != nil {
		t.Fatalf("remove design: %v", err)
	}

	var count int64
	if err := db.Table("workspaces_designs_mappings").Where("workspace_id = ? AND design_id = ?", wsID, designID).Count(&count).Error; err != nil {
		t.Fatalf("count mappings: %v", err)
	}
	if count != 0 {
		t.Errorf("workspaces_designs_mappings has %d mapping(s) remaining after removal; want 0", count)
	}
}

// TestEnvironmentDeleteCascadesToWorkspaceMappings verifies that deleting an
// environment cleans up all references in workspaces_environments_mappings.
func TestEnvironmentDeleteCascadesToWorkspaceMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}
	ep := &EnvironmentPersister{DB: db}

	env := &environment.Environment{
		ID:   mustUUID(t),
		Name: "test-environment-cascade",
	}
	if _, err := ep.SaveEnvironment(env); err != nil {
		t.Fatalf("save env: %v", err)
	}

	wsID := mustUUID(t)
	if _, err := wp.AddEnvironmentToWorkspace(wsID, env.ID); err != nil {
		t.Fatalf("add env to ws: %v", err)
	}

	if _, err := ep.DeleteEnvironment(env); err != nil {
		t.Fatalf("delete env: %v", err)
	}

	var count int64
	if err := db.Table("workspaces_environments_mappings").Where("environment_id = ?", env.ID).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Errorf("workspaces_environments_mappings has %d orphaned mapping(s) after environment deletion; want 0", count)
	}
}

// TestMesheryPatternDeleteCascadesToWorkspaceMappings verifies that deleting a
// design cleans up all of its references in workspaces_designs_mappings.
func TestMesheryPatternDeleteCascadesToWorkspaceMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}
	mpp := &MesheryPatternPersister{DB: db}

	designID := mustUUID(t)
	pattern := &MesheryPattern{
		ID:          &designID,
		Name:        "test-cascade-pattern",
		PatternFile: "name: test\nservices: {}\n",
	}
	if _, err := mpp.SaveMesheryPattern(pattern); err != nil {
		t.Fatalf("save pattern: %v", err)
	}

	wsID := mustUUID(t)
	if _, err := wp.AddDesignToWorkspace(wsID, designID); err != nil {
		t.Fatalf("add design to ws: %v", err)
	}

	if _, err := mpp.DeleteMesheryPattern(designID); err != nil {
		t.Fatalf("delete pattern: %v", err)
	}

	var count int64
	if err := db.Table("workspaces_designs_mappings").Where("design_id = ?", designID).Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 0 {
		t.Errorf("workspaces_designs_mappings has %d orphaned mapping(s) after design deletion; want 0", count)
	}
}

// TestWorkspaceAddDesignPreservesExistingMapping verifies that adding an already
// assigned design returns the existing mapping without deleting or recreating IDs.
func TestWorkspaceAddDesignPreservesExistingMapping(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}

	wsID, designID := mustUUID(t), mustUUID(t)

	raw1, err := wp.AddDesignToWorkspace(wsID, designID)
	if err != nil {
		t.Fatalf("add design #1: %v", err)
	}

	var m1 workspace.WorkspacesDesignsMapping
	if err := json.Unmarshal(raw1, &m1); err != nil {
		t.Fatalf("decode #1: %v", err)
	}

	raw2, err := wp.AddDesignToWorkspace(wsID, designID)
	if err != nil {
		t.Fatalf("add design #2: %v", err)
	}

	var m2 workspace.WorkspacesDesignsMapping
	if err := json.Unmarshal(raw2, &m2); err != nil {
		t.Fatalf("decode #2: %v", err)
	}

	if m1.ID != m2.ID {
		t.Errorf("mapping ID changed on second add: %v != %v", m1.ID, m2.ID)
	}
}

// TestWorkspaceDeleteNotFound verifies that deleting a non-existent workspace
// returns ErrResultNotFound instead of nil error.
func TestWorkspaceDeleteNotFound(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}

	nonExistentWS := &workspace.Workspace{
		ID: mustUUID(t),
	}
	_, err := wp.DeleteWorkspace(nonExistentWS)
	if err == nil {
		t.Fatalf("DeleteWorkspace returned nil error for non-existent workspace; want ErrResultNotFound")
	}
}
