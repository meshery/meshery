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
