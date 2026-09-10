package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	teamv1beta2 "github.com/meshery/schemas/models/v1beta2/team"
	"github.com/meshery/schemas/models/v1beta3/environment"
	workspace "github.com/meshery/schemas/models/v1beta3/workspace"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newWorkspaceFilterTestDB builds an in-memory database carrying the tables the
// workspace-scoped listings query.
func newWorkspaceFilterTestDB(t *testing.T) *database.Handler {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	handler := &database.Handler{DB: db}
	for _, model := range []interface{}{
		&environment.Environment{},
		&teamv1beta2.Team{},
		&workspace.WorkspacesEnvironmentsMapping{},
		&workspace.WorkspacesTeamsMapping{},
	} {
		if err := handler.AutoMigrate(model); err != nil {
			t.Fatalf("failed to auto-migrate %T: %v", model, err)
		}
	}

	return handler
}

func newWorkspaceTestUUID(t *testing.T) core.Uuid {
	t.Helper()

	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}

	return core.Uuid(id)
}

// seedWorkspaceEnvironment persists one environment and assigns it to workspaceID.
func seedWorkspaceEnvironment(t *testing.T, handler *database.Handler, workspaceID core.Uuid, name string, owner core.Uuid) {
	t.Helper()

	envID := newWorkspaceTestUUID(t)
	ownerID := owner
	env := environment.Environment{
		ID:             envID,
		Name:           name,
		OrganizationID: newWorkspaceTestUUID(t),
		Owner:          &ownerID,
	}
	if err := handler.Create(&env).Error; err != nil {
		t.Fatalf("failed to seed environment %q: %v", name, err)
	}

	mapping := workspace.WorkspacesEnvironmentsMapping{
		ID:            newWorkspaceTestUUID(t),
		WorkspaceID:   workspaceID,
		EnvironmentID: envID,
	}
	if err := handler.Create(&mapping).Error; err != nil {
		t.Fatalf("failed to assign environment %q: %v", name, err)
	}
}

// seedWorkspaceTeam persists one team and assigns it to workspaceID.
func seedWorkspaceTeam(t *testing.T, handler *database.Handler, workspaceID core.Uuid, name string, owner core.Uuid) {
	t.Helper()

	teamID := newWorkspaceTestUUID(t)
	ownerID := owner
	team := teamv1beta2.Team{
		ID:    teamID,
		Name:  name,
		Owner: &ownerID,
	}
	if err := handler.Create(&team).Error; err != nil {
		t.Fatalf("failed to seed team %q: %v", name, err)
	}

	mapping := workspace.WorkspacesTeamsMapping{
		ID:          newWorkspaceTestUUID(t),
		WorkspaceID: workspaceID,
		TeamID:      teamID,
	}
	if err := handler.Create(&mapping).Error; err != nil {
		t.Fatalf("failed to assign team %q: %v", name, err)
	}
}

// TestGetWorkspaceEnvironmentsOwnerFilter is the regression for the unreachable
// owner filter. Against the previous code this returned both environments,
// because the JSON filter reached utils.ApplyFilters as a single token that
// matched no dynamic key and was a silent no-op.
func TestGetWorkspaceEnvironmentsOwnerFilter(t *testing.T) {
	handler := newWorkspaceFilterTestDB(t)
	wp := &WorkspacePersister{DB: handler}

	workspaceID := newWorkspaceTestUUID(t)
	ownerA := newWorkspaceTestUUID(t)
	ownerB := newWorkspaceTestUUID(t)

	seedWorkspaceEnvironment(t, handler, workspaceID, "env-owned-by-a", ownerA)
	seedWorkspaceEnvironment(t, handler, workspaceID, "env-owned-by-b", ownerB)

	filter, err := json.Marshal(map[string]interface{}{
		"assigned": true,
		"owner":    ownerA.String(),
	})
	if err != nil {
		t.Fatalf("failed to build filter: %v", err)
	}

	raw, err := wp.GetWorkspaceEnvironments(workspaceID, "", "", "0", "10", string(filter))
	if err != nil {
		t.Fatalf("GetWorkspaceEnvironments returned error: %v", err)
	}

	var page environment.EnvironmentPage
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}

	if len(page.Environments) != 1 {
		names := make([]string, 0, len(page.Environments))
		for _, e := range page.Environments {
			names = append(names, e.Name)
		}
		t.Fatalf("owner filter returned %d environments (%v), want 1", len(page.Environments), names)
	}
	if page.Environments[0].Name != "env-owned-by-a" {
		t.Errorf("owner filter returned %q, want %q", page.Environments[0].Name, "env-owned-by-a")
	}
}

// TestGetWorkspaceTeamsOwnerFilter pins the same predicate on the other table
// that carries a real `owner` column.
func TestGetWorkspaceTeamsOwnerFilter(t *testing.T) {
	handler := newWorkspaceFilterTestDB(t)
	wp := &WorkspacePersister{DB: handler}

	workspaceID := newWorkspaceTestUUID(t)
	ownerA := newWorkspaceTestUUID(t)
	ownerB := newWorkspaceTestUUID(t)

	seedWorkspaceTeam(t, handler, workspaceID, "team-owned-by-a", ownerA)
	seedWorkspaceTeam(t, handler, workspaceID, "team-owned-by-b", ownerB)

	filter, err := json.Marshal(map[string]interface{}{
		"assigned": true,
		"owner":    ownerA.String(),
	})
	if err != nil {
		t.Fatalf("failed to build filter: %v", err)
	}

	raw, err := wp.GetWorkspaceTeams(workspaceID, "", "", "0", "10", string(filter))
	if err != nil {
		t.Fatalf("GetWorkspaceTeams returned error: %v", err)
	}

	var page struct {
		Teams []teamv1beta2.Team `json:"teams"`
	}
	if err := json.Unmarshal(raw, &page); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}

	if len(page.Teams) != 1 {
		names := make([]string, 0, len(page.Teams))
		for _, tm := range page.Teams {
			names = append(names, tm.Name)
		}
		t.Fatalf("owner filter returned %d teams (%v), want 1", len(page.Teams), names)
	}
	if page.Teams[0].Name != "team-owned-by-a" {
		t.Errorf("owner filter returned %q, want %q", page.Teams[0].Name, "team-owned-by-a")
	}
}

// TestGetWorkspaceEnvironmentsRejectsBadFilter pins that a malformed filter
// comes back as the structured caller error the handler renders as 400,
// instead of the raw encoding/json message the previous code leaked.
func TestGetWorkspaceEnvironmentsRejectsBadFilter(t *testing.T) {
	handler := newWorkspaceFilterTestDB(t)
	wp := &WorkspacePersister{DB: handler}
	workspaceID := newWorkspaceTestUUID(t)

	for _, filter := range []string{"owner abc", `{"assigned":"false"}`, `null`} {
		t.Run(filter, func(t *testing.T) {
			_, err := wp.GetWorkspaceEnvironments(workspaceID, "", "", "0", "10", filter)
			if err == nil {
				t.Fatalf("GetWorkspaceEnvironments(%q) returned no error, want one", filter)
			}
			if code := meshkiterrors.GetCode(err); code != ErrInvalidWorkspaceFilterCode {
				t.Errorf("error code = %q, want %q", code, ErrInvalidWorkspaceFilterCode)
			}
		})
	}
}

// TestGetWorkspaceEnvironmentsAssignedFilter pins that the `assigned` filter,
// the one grammar that did work before, still does.
func TestGetWorkspaceEnvironmentsAssignedFilter(t *testing.T) {
	handler := newWorkspaceFilterTestDB(t)
	wp := &WorkspacePersister{DB: handler}

	workspaceID := newWorkspaceTestUUID(t)
	owner := newWorkspaceTestUUID(t)

	seedWorkspaceEnvironment(t, handler, workspaceID, "assigned-env", owner)

	// An environment that exists but is assigned to no workspace.
	unassignedID := newWorkspaceTestUUID(t)
	ownerID := owner
	if err := handler.Create(&environment.Environment{
		ID:             unassignedID,
		Name:           "unassigned-env",
		OrganizationID: newWorkspaceTestUUID(t),
		Owner:          &ownerID,
	}).Error; err != nil {
		t.Fatalf("failed to seed unassigned environment: %v", err)
	}

	for _, tc := range []struct {
		filter string
		want   string
	}{
		{`{"assigned":true}`, "assigned-env"},
		{`{"assigned":false}`, "unassigned-env"},
	} {
		t.Run(tc.filter, func(t *testing.T) {
			raw, err := wp.GetWorkspaceEnvironments(workspaceID, "", "", "0", "10", tc.filter)
			if err != nil {
				t.Fatalf("GetWorkspaceEnvironments returned error: %v", err)
			}

			var page environment.EnvironmentPage
			if err := json.Unmarshal(raw, &page); err != nil {
				t.Fatalf("failed to decode page: %v", err)
			}

			if len(page.Environments) != 1 || page.Environments[0].Name != tc.want {
				t.Errorf("filter %s returned %d environments, want exactly %q", tc.filter, len(page.Environments), tc.want)
			}
		})
	}
}
