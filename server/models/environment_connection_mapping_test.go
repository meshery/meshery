package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	schemasEnv "github.com/meshery/schemas/models/v1beta3/environment"
)

func countEnvironmentConnectionMappings(t *testing.T, db *database.Handler, envID, connID core.Uuid) int64 {
	t.Helper()
	var n int64
	if err := db.Table("environment_connection_mappings").Where("environment_id = ? AND connection_id = ?", envID, connID).Count(&n).Error; err != nil {
		t.Fatalf("count mappings: %v", err)
	}
	return n
}

// TestEnvironmentAssignConnectionTwiceThenRemoveOnce tests that assigning a connection
// to an environment twice only creates a single mapping, and removing it once completely unassigns it.
func TestEnvironmentAssignConnectionTwiceThenRemoveOnce(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	envID, connID := mustUUID(t), mustUUID(t)

	for i := 1; i <= 2; i++ {
		if _, err := ep.AddConnectionToEnvironment(envID, connID); err != nil {
			t.Fatalf("assign #%d: %v", i, err)
		}
	}
	if n := countEnvironmentConnectionMappings(t, db, envID, connID); n != 1 {
		t.Fatalf("mappings after assigning connection twice = %d, want 1", n)
	}

	if _, err := ep.DeleteConnectionFromEnvironment(envID, connID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if n := countEnvironmentConnectionMappings(t, db, envID, connID); n != 0 {
		t.Errorf("mappings after removing connection = %d, want 0", n)
	}
}

// TestEnvironmentRemoveDeletesStoredDuplicates tests that removing a connection
// cleans up all duplicate mapping rows stored before assignments were checked.
func TestEnvironmentRemoveDeletesStoredDuplicates(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	envID, connID := mustUUID(t), mustUUID(t)

	insert := "INSERT INTO environment_connection_mappings (id, environment_id, connection_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)"
	for i := 1; i <= 2; i++ {
		if err := db.Exec(insert, mustUUID(t), envID, connID, time.Now(), time.Now()).Error; err != nil {
			t.Fatalf("store mapping #%d: %v", i, err)
		}
	}

	if _, err := ep.DeleteConnectionFromEnvironment(envID, connID); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if n := countEnvironmentConnectionMappings(t, db, envID, connID); n != 0 {
		t.Errorf("mappings after removing duplicates = %d, want 0", n)
	}
}

// TestEnvironmentConnectionUnlistedAfterRemoval verifies that an assigned connection
// does not appear duplicated in GetEnvironmentConnections and is completely removed after one unassignment.
func TestEnvironmentConnectionUnlistedAfterRemoval(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	envID := mustUUID(t)
	connID := mustUUID(t)

	savedConn, err := cp.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "test-cluster",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("seed connection: %v", err)
	}

	for i := 1; i <= 2; i++ {
		if _, err := ep.AddConnectionToEnvironment(envID, savedConn.ID); err != nil {
			t.Fatalf("assign #%d: %v", i, err)
		}
	}

	rawBefore, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("list before removal: %v", err)
	}
	var pageBefore connections.ConnectionPage
	if err := json.Unmarshal(rawBefore, &pageBefore); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if pageBefore.TotalCount != 1 {
		t.Errorf("assigned connections before removal = %d, want 1 (got duplicates)", pageBefore.TotalCount)
	}

	if _, err := ep.DeleteConnectionFromEnvironment(envID, savedConn.ID); err != nil {
		t.Fatalf("remove: %v", err)
	}

	rawAfter, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("list after removal: %v", err)
	}
	var pageAfter connections.ConnectionPage
	if err := json.Unmarshal(rawAfter, &pageAfter); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if pageAfter.TotalCount != 0 {
		t.Errorf("assigned connections after removal = %d, want 0", pageAfter.TotalCount)
	}
}

// TestDeleteConnectionCleansUpEnvironmentMappings tests that deleting a connection
// also removes its orphaned environment connection mappings.
func TestDeleteConnectionCleansUpEnvironmentMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	envID := mustUUID(t)
	connID := mustUUID(t)

	savedConn, err := cp.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "test-cluster-to-delete",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("seed connection: %v", err)
	}

	if _, err := ep.AddConnectionToEnvironment(envID, savedConn.ID); err != nil {
		t.Fatalf("assign: %v", err)
	}

	if n := countEnvironmentConnectionMappings(t, db, envID, savedConn.ID); n != 1 {
		t.Fatalf("expected 1 mapping, got %d", n)
	}

	if _, err := cp.DeleteConnectionById(savedConn.ID); err != nil {
		t.Fatalf("delete connection: %v", err)
	}

	if n := countEnvironmentConnectionMappings(t, db, envID, savedConn.ID); n != 0 {
		t.Errorf("orphaned environment mappings after deleting connection = %d, want 0", n)
	}
}

// TestDeleteEnvironmentCleansUpConnectionMappings tests that deleting an environment
// cascades and removes its environment connection mappings.
func TestDeleteEnvironmentCleansUpConnectionMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	env := &schemasEnv.Environment{Name: "test-env-to-delete"}
	if _, err := ep.SaveEnvironment(env); err != nil {
		t.Fatalf("save environment: %v", err)
	}

	savedConn, err := cp.SaveConnection(&connections.Connection{
		ID:             mustUUID(t),
		Name:           "test-cluster-env-delete",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("seed connection: %v", err)
	}

	if _, err := ep.AddConnectionToEnvironment(env.ID, savedConn.ID); err != nil {
		t.Fatalf("assign: %v", err)
	}

	if n := countEnvironmentConnectionMappings(t, db, env.ID, savedConn.ID); n != 1 {
		t.Fatalf("expected 1 mapping, got %d", n)
	}

	if _, err := ep.DeleteEnvironmentByID(env.ID); err != nil {
		t.Fatalf("delete environment: %v", err)
	}

	if n := countEnvironmentConnectionMappings(t, db, env.ID, savedConn.ID); n != 0 {
		t.Errorf("orphaned environment mappings after deleting environment = %d, want 0", n)
	}
}

// TestGetEnvironmentConnectionsExcludesSoftDeletedConnections verifies that soft-deleted
// connections are not returned in assigned or unassigned connection queries.
func TestGetEnvironmentConnectionsExcludesSoftDeletedConnections(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	envID := mustUUID(t)
	connAssignedID := mustUUID(t)
	connUnassignedID := mustUUID(t)

	// Create and assign connection 1
	savedAssigned, err := cp.SaveConnection(&connections.Connection{
		ID:             connAssignedID,
		Name:           "assigned-cluster",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save assigned conn: %v", err)
	}
	if _, err := ep.AddConnectionToEnvironment(envID, savedAssigned.ID); err != nil {
		t.Fatalf("add assigned conn: %v", err)
	}

	// Create unassigned connection 2
	savedUnassigned, err := cp.SaveConnection(&connections.Connection{
		ID:             connUnassignedID,
		Name:           "unassigned-cluster",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save unassigned conn: %v", err)
	}

	// Verify before soft-deletion: 1 assigned, 1 unassigned
	rawAssigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("get assigned before: %v", err)
	}
	var pageAssigned connections.ConnectionPage
	if err := json.Unmarshal(rawAssigned, &pageAssigned); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageAssigned.TotalCount != 1 {
		t.Fatalf("expected 1 assigned conn before soft-delete, got %d", pageAssigned.TotalCount)
	}

	rawUnassigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":false}`)
	if err != nil {
		t.Fatalf("get unassigned before: %v", err)
	}
	var pageUnassigned connections.ConnectionPage
	if err := json.Unmarshal(rawUnassigned, &pageUnassigned); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageUnassigned.TotalCount != 1 {
		t.Fatalf("expected 1 unassigned conn before soft-delete, got %d", pageUnassigned.TotalCount)
	}

	// Soft-delete both connections
	now := time.Now()
	if err := db.Table("connections").Where("id IN ?", []core.Uuid{savedAssigned.ID, savedUnassigned.ID}).Update("deleted_at", now).Error; err != nil {
		t.Fatalf("soft-delete connections: %v", err)
	}

	// Both queries should now return 0 connections by default
	rawAssignedAfter, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("get assigned after: %v", err)
	}
	var pageAssignedAfter connections.ConnectionPage
	if err := json.Unmarshal(rawAssignedAfter, &pageAssignedAfter); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageAssignedAfter.TotalCount != 0 {
		t.Errorf("expected 0 assigned conns after soft-delete, got %d", pageAssignedAfter.TotalCount)
	}

	rawUnassignedAfter, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":false}`)
	if err != nil {
		t.Fatalf("get unassigned after: %v", err)
	}
	var pageUnassignedAfter connections.ConnectionPage
	if err := json.Unmarshal(rawUnassignedAfter, &pageUnassignedAfter); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageUnassignedAfter.TotalCount != 0 {
		t.Errorf("expected 0 unassigned conns after soft-delete, got %d", pageUnassignedAfter.TotalCount)
	}

	// When explicitly querying for deleted connections, they should be found
	rawDeleted, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true,"deletedAt":true}`)
	if err != nil {
		t.Fatalf("get deleted assigned: %v", err)
	}
	var pageDeleted connections.ConnectionPage
	if err := json.Unmarshal(rawDeleted, &pageDeleted); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageDeleted.TotalCount != 1 {
		t.Errorf("expected 1 deleted assigned conn, got %d", pageDeleted.TotalCount)
	}
}

// TestGetEnvironmentConnectionsFilterSafeStringParsing verifies that string-encoded
// boolean filter parameters parse correctly and take effect on the query results.
func TestGetEnvironmentConnectionsFilterSafeStringParsing(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}
	envID := mustUUID(t)

	// Seed 1 active assigned connection
	connAssigned, err := cp.SaveConnection(&connections.Connection{
		ID:             mustUUID(t),
		Name:           "assigned-conn",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save assigned conn: %v", err)
	}
	if _, err := ep.AddConnectionToEnvironment(envID, connAssigned.ID); err != nil {
		t.Fatalf("add assigned conn: %v", err)
	}

	// Seed 1 active unassigned connection
	connUnassigned, err := cp.SaveConnection(&connections.Connection{
		ID:             mustUUID(t),
		Name:           "unassigned-conn",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save unassigned conn: %v", err)
	}

	// Seed 1 soft-deleted assigned connection
	connDeletedAssigned, err := cp.SaveConnection(&connections.Connection{
		ID:             mustUUID(t),
		Name:           "deleted-assigned-conn",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save deleted conn: %v", err)
	}
	if _, err := ep.AddConnectionToEnvironment(envID, connDeletedAssigned.ID); err != nil {
		t.Fatalf("add deleted assigned conn: %v", err)
	}
	if err := db.Table("connections").Where("id = ?", connDeletedAssigned.ID).Update("deleted_at", time.Now()).Error; err != nil {
		t.Fatalf("soft-delete connection: %v", err)
	}

	cases := []struct {
		filter        string
		wantCount     int
		wantConnID    core.Uuid
		shouldCheckID bool
	}{
		{filter: `{"assigned":"true"}`, wantCount: 1, wantConnID: connAssigned.ID, shouldCheckID: true},
		{filter: `{"assigned":"false"}`, wantCount: 1, wantConnID: connUnassigned.ID, shouldCheckID: true},
		{filter: `{"assigned":"true","deletedAt":"true"}`, wantCount: 1, wantConnID: connDeletedAssigned.ID, shouldCheckID: true},
		{filter: `{"assigned":"true","deleted_at":"false"}`, wantCount: 1, wantConnID: connAssigned.ID, shouldCheckID: true},
		{filter: `{"assigned":"false","deletedAt":"false"}`, wantCount: 1, wantConnID: connUnassigned.ID, shouldCheckID: true},
		{filter: `{"assigned":"false","deleted_at":"true"}`, wantCount: 0, shouldCheckID: false},
	}

	for _, tc := range cases {
		raw, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", tc.filter)
		if err != nil {
			t.Fatalf("filter %q returned error: %v", tc.filter, err)
		}
		var page connections.ConnectionPage
		if err := json.Unmarshal(raw, &page); err != nil {
			t.Fatalf("filter %q returned invalid JSON: %v", tc.filter, err)
		}
		if page.TotalCount != tc.wantCount {
			t.Errorf("filter %q TotalCount = %d, want %d", tc.filter, page.TotalCount, tc.wantCount)
		}
		if tc.shouldCheckID {
			if len(page.Connections) != tc.wantCount {
				t.Fatalf("filter %q returned %d connections in page, want %d", tc.filter, len(page.Connections), tc.wantCount)
			}
			if len(page.Connections) > 0 && page.Connections[0].ID != tc.wantConnID {
				t.Errorf("filter %q returned connection ID %s, want %s", tc.filter, page.Connections[0].ID, tc.wantConnID)
			}
		}
	}
}

// TestEnvironmentReassignRestoresSoftDeletedMapping verifies that assigning a connection
// whose mapping was previously soft-deleted clears deleted_at and restores the assignment.
func TestEnvironmentReassignRestoresSoftDeletedMapping(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	envID := mustUUID(t)
	connID := mustUUID(t)

	savedConn, err := cp.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "reassign-test-cluster",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save conn: %v", err)
	}

	// 1. Initial assignment
	if _, err := ep.AddConnectionToEnvironment(envID, savedConn.ID); err != nil {
		t.Fatalf("initial assign: %v", err)
	}

	// 2. Soft-delete the mapping
	now := time.Now()
	if err := db.Table("environment_connection_mappings").Where("environment_id = ? AND connection_id = ?", envID, savedConn.ID).Update("deleted_at", now).Error; err != nil {
		t.Fatalf("soft-delete mapping: %v", err)
	}

	// Connection should show as unassigned
	rawUnassigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":false}`)
	if err != nil {
		t.Fatalf("get unassigned: %v", err)
	}
	var unassignedPage connections.ConnectionPage
	if err := json.Unmarshal(rawUnassigned, &unassignedPage); err != nil {
		t.Fatalf("unmarshal unassigned: %v", err)
	}
	if unassignedPage.TotalCount != 1 {
		t.Fatalf("expected 1 unassigned conn after mapping soft-delete, got %d", unassignedPage.TotalCount)
	}

	// 3. Reassign the connection
	if _, err := ep.AddConnectionToEnvironment(envID, savedConn.ID); err != nil {
		t.Fatalf("reassign conn: %v", err)
	}

	// 4. Verify mapping is active again (deleted_at is NULL) and shows as assigned
	var activeCount, deletedCount int64
	if err := db.Table("environment_connection_mappings").Where("environment_id = ? AND connection_id = ? AND deleted_at IS NULL", envID, savedConn.ID).Count(&activeCount).Error; err != nil {
		t.Fatalf("count active: %v", err)
	}
	if activeCount != 1 {
		t.Errorf("expected 1 active mapping row in db, got %d", activeCount)
	}
	if err := db.Table("environment_connection_mappings").Where("environment_id = ? AND connection_id = ? AND deleted_at IS NOT NULL", envID, savedConn.ID).Count(&deletedCount).Error; err != nil {
		t.Fatalf("count deleted: %v", err)
	}
	if deletedCount != 0 {
		t.Errorf("expected 0 soft-deleted mapping rows in db, got %d", deletedCount)
	}

	rawAssigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("get assigned after reassign: %v", err)
	}
	var assignedPage connections.ConnectionPage
	if err := json.Unmarshal(rawAssigned, &assignedPage); err != nil {
		t.Fatalf("unmarshal assigned: %v", err)
	}
	if assignedPage.TotalCount != 1 {
		t.Errorf("expected 1 assigned conn after reassign, got %d", assignedPage.TotalCount)
	}
}

// TestGetEnvironmentConnectionsExcludesSoftDeletedMappings verifies that when an assignment
// mapping is soft-deleted, the connection is treated as unassigned.
func TestGetEnvironmentConnectionsExcludesSoftDeletedMappings(t *testing.T) {
	db := newWorkspaceMappingDB(t)
	ep := &EnvironmentPersister{DB: db}
	cp := &ConnectionPersister{DB: db}

	envID := mustUUID(t)
	connID := mustUUID(t)

	savedConn, err := cp.SaveConnection(&connections.Connection{
		ID:             connID,
		Name:           "mapping-soft-delete-cluster",
		Kind:           "kubernetes",
		ConnectionType: "orchestrator",
		SubType:        "k8s",
		Status:         connections.CONNECTED,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("save conn: %v", err)
	}

	if _, err := ep.AddConnectionToEnvironment(envID, savedConn.ID); err != nil {
		t.Fatalf("assign: %v", err)
	}

	// Soft-delete the mapping in environment_connection_mappings
	now := time.Now()
	if err := db.Table("environment_connection_mappings").Where("environment_id = ? AND connection_id = ?", envID, savedConn.ID).Update("deleted_at", now).Error; err != nil {
		t.Fatalf("soft-delete mapping: %v", err)
	}

	// Assigned query should return 0
	rawAssigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":true}`)
	if err != nil {
		t.Fatalf("get assigned: %v", err)
	}
	var pageAssigned connections.ConnectionPage
	if err := json.Unmarshal(rawAssigned, &pageAssigned); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageAssigned.TotalCount != 0 {
		t.Errorf("expected 0 assigned conns when mapping is soft-deleted, got %d", pageAssigned.TotalCount)
	}

	// Unassigned query should return 1 (it is available to assign again)
	rawUnassigned, err := ep.GetEnvironmentConnections(envID, "", "", "0", "10", `{"assigned":false}`)
	if err != nil {
		t.Fatalf("get unassigned: %v", err)
	}
	var pageUnassigned connections.ConnectionPage
	if err := json.Unmarshal(rawUnassigned, &pageUnassigned); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if pageUnassigned.TotalCount != 1 {
		t.Errorf("expected 1 unassigned conn when mapping is soft-deleted, got %d", pageUnassigned.TotalCount)
	}
}
