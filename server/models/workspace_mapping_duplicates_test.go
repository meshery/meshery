package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

func newWorkspaceMappingDB(t *testing.T) *database.Handler {
	t.Helper()
	db := newMigratedDB(t)
	// The system migration still builds environments from v1beta1 (#21802), so
	// add the v1beta3 columns the environment persister writes.
	if err := db.AutoMigrate(&environment.Environment{}); err != nil {
		t.Fatalf("migrate environments: %v", err)
	}
	return db
}

// workspaceMappingKinds are the workspace assignments that map one member to a
// workspace per row.
var workspaceMappingKinds = []struct {
	name, table, column string
	assign, remove      func(*WorkspacePersister, core.Uuid, core.Uuid) ([]byte, error)
}{
	{"environment", "workspaces_environments_mappings", "environment_id", (*WorkspacePersister).AddEnvironmentToWorkspace, (*WorkspacePersister).DeleteEnvironmentFromWorkspace},
	{"view", "workspaces_views_mappings", "view_id", (*WorkspacePersister).AddViewToWorkspace, (*WorkspacePersister).DeleteViewFromWorkspace},
	{"team", "workspaces_teams_mappings", "team_id", (*WorkspacePersister).AddTeamToWorkspace, (*WorkspacePersister).DeleteTeamFromWorkspace},
	{"design", "workspaces_designs_mappings", "design_id", (*WorkspacePersister).AddDesignToWorkspace, (*WorkspacePersister).DeleteDesignFromWorkspace},
}

func countWorkspaceMappings(t *testing.T, db *database.Handler, table, column string, workspaceID, memberID core.Uuid) int64 {
	t.Helper()
	var n int64
	if err := db.Table(table).Where("workspace_id = ? AND "+column+" = ?", workspaceID, memberID).Count(&n).Error; err != nil {
		t.Fatalf("count %s: %v", table, err)
	}
	return n
}

// Assigning a member twice keeps one mapping, and one removal unassigns it.
func TestWorkspaceAssignTwiceThenRemoveOnce(t *testing.T) {
	for _, kind := range workspaceMappingKinds {
		t.Run(kind.name, func(t *testing.T) {
			db := newWorkspaceMappingDB(t)
			wp := &WorkspacePersister{DB: db}
			workspaceID, memberID := mustUUID(t), mustUUID(t)

			for i := 1; i <= 2; i++ {
				if _, err := kind.assign(wp, workspaceID, memberID); err != nil {
					t.Fatalf("assign #%d: %v", i, err)
				}
			}
			if n := countWorkspaceMappings(t, db, kind.table, kind.column, workspaceID, memberID); n != 1 {
				t.Fatalf("mappings after assigning twice = %d, want 1", n)
			}

			if _, err := kind.remove(wp, workspaceID, memberID); err != nil {
				t.Fatalf("remove: %v", err)
			}
			if n := countWorkspaceMappings(t, db, kind.table, kind.column, workspaceID, memberID); n != 0 {
				t.Errorf("mappings after removing = %d, want 0", n)
			}
		})
	}
}

// Duplicate mappings stored before assignments were checked all go with one
// removal.
func TestWorkspaceRemoveDeletesStoredDuplicates(t *testing.T) {
	for _, kind := range workspaceMappingKinds {
		t.Run(kind.name, func(t *testing.T) {
			db := newWorkspaceMappingDB(t)
			wp := &WorkspacePersister{DB: db}
			workspaceID, memberID := mustUUID(t), mustUUID(t)

			insert := fmt.Sprintf("INSERT INTO %s (id, workspace_id, %s, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", kind.table, kind.column)
			for i := 1; i <= 2; i++ {
				if err := db.Exec(insert, mustUUID(t), workspaceID, memberID, time.Now(), time.Now()).Error; err != nil {
					t.Fatalf("store mapping #%d: %v", i, err)
				}
			}

			if _, err := kind.remove(wp, workspaceID, memberID); err != nil {
				t.Fatalf("remove: %v", err)
			}
			if n := countWorkspaceMappings(t, db, kind.table, kind.column, workspaceID, memberID); n != 0 {
				t.Errorf("mappings after removing = %d, want 0", n)
			}
		})
	}
}

// An environment assigned twice and removed once is no longer listed among
// the workspace's environments.
func TestWorkspaceEnvironmentUnlistedAfterRemoval(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	wp := &WorkspacePersister{DB: db}
	env := &environment.Environment{Name: "dev"}
	if _, err := (&EnvironmentPersister{DB: db}).SaveEnvironment(env); err != nil {
		t.Fatalf("save environment: %v", err)
	}
	workspaceID := mustUUID(t)

	for i := 1; i <= 2; i++ {
		if _, err := wp.AddEnvironmentToWorkspace(workspaceID, env.ID); err != nil {
			t.Fatalf("assign #%d: %v", i, err)
		}
	}
	if _, err := wp.DeleteEnvironmentFromWorkspace(workspaceID, env.ID); err != nil {
		t.Fatalf("remove: %v", err)
	}

	raw, err := wp.GetWorkspaceEnvironments(workspaceID, "", "", "0", "10", "")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	var page struct {
		TotalCount int `json:"totalCount"`
	}
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("decode %s: %v", raw, err)
	}
	if page.TotalCount != 0 {
		t.Errorf("assigned environments after removal = %d, want 0", page.TotalCount)
	}
}
