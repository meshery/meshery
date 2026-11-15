package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newK8sContextPersister returns a persister backed by an in-memory SQLite database,
// with the k8s_contexts table migrated only when migrate is true.
func newK8sContextPersister(t *testing.T, migrate bool) *MesheryK8sContextPersister {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	// Each pooled connection to ":memory:" gets its own empty database.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to access underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })

	if migrate {
		if err := db.AutoMigrate(&K8sContext{}); err != nil {
			t.Fatalf("failed to migrate k8s contexts: %v", err)
		}
	}
	return &MesheryK8sContextPersister{DB: &database.Handler{DB: db}}
}

// TestGetMesheryK8sContexts_ReturnsPage verifies case-insensitive search, ordering,
// pagination and the total count on a successful read.
func TestGetMesheryK8sContexts_ReturnsPage(t *testing.T) {
	persister := newK8sContextPersister(t, true)
	for _, c := range []K8sContext{{ID: "1", Name: "alpha"}, {ID: "2", Name: "beta"}, {ID: "3", Name: "alphabet"}} {
		if err := persister.DB.Create(&c).Error; err != nil {
			t.Fatalf("failed to seed context %s: %v", c.Name, err)
		}
	}

	resp, err := persister.GetMesheryK8sContexts("ALPHA", "name asc", 0, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var page MesheryK8sContextPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}
	if page.TotalCount != 2 {
		t.Fatalf("expected total count 2, got %d", page.TotalCount)
	}
	if len(page.Contexts) != 1 || page.Contexts[0].Name != "alpha" {
		t.Fatalf("expected only context alpha on the page, got %+v", page.Contexts)
	}
}

// TestGetMesheryK8sContexts_PropagatesCountError verifies that a failed Count, here a
// missing table, is returned as ErrDBRead with no response body.
func TestGetMesheryK8sContexts_PropagatesCountError(t *testing.T) {
	persister := newK8sContextPersister(t, false)

	resp, err := persister.GetMesheryK8sContexts("", "", 0, 10)
	if err == nil {
		t.Fatalf("expected an error, got response %s", resp)
	}
	if resp != nil {
		t.Fatalf("expected no response body, got %s", resp)
	}
	if code := meshkiterrors.GetCode(err); code != ErrDBReadCode {
		t.Fatalf("expected error code %s, got %s", ErrDBReadCode, code)
	}
}

// TestGetMesheryK8sContexts_PropagatesFindError verifies that a failed Find, here a row
// whose auth column cannot be decoded, is returned as ErrDBRead with no response body.
func TestGetMesheryK8sContexts_PropagatesFindError(t *testing.T) {
	persister := newK8sContextPersister(t, true)
	// COUNT never decodes a row, so an undecodable column fails only in Find.
	if err := persister.DB.Exec(`INSERT INTO k8s_contexts (id, name, auth) VALUES ('1', 'corrupt', 'not-json')`).Error; err != nil {
		t.Fatalf("failed to seed corrupt context: %v", err)
	}

	resp, err := persister.GetMesheryK8sContexts("", "", 0, 10)
	if err == nil {
		t.Fatalf("expected an error, got response %s", resp)
	}
	if resp != nil {
		t.Fatalf("expected no response body, got %s", resp)
	}
	if code := meshkiterrors.GetCode(err); code != ErrDBReadCode {
		t.Fatalf("expected error code %s, got %s", ErrDBReadCode, code)
	}
	if !strings.Contains(err.Error(), "Unable to read data") {
		t.Fatalf("expected a read error, got %v", err)
	}
}

// TestGetMesheryK8sContextsWithStatus_FiltersByConnectionStatus verifies that
// GetMesheryK8sContextsWithStatus filters contexts by their associated connection status.
func TestGetMesheryK8sContextsWithStatus_FiltersByConnectionStatus(t *testing.T) {
	persister := newK8sContextPersister(t, true)
	if err := persister.DB.AutoMigrate(&connections.Connection{}); err != nil {
		t.Fatalf("failed to migrate connections: %v", err)
	}

	conn1ID := "11111111-1111-1111-1111-111111111111"
	conn2ID := "22222222-2222-2222-2222-222222222222"

	c1 := connections.Connection{ID: uuid.FromStringOrNil(conn1ID), Status: connections.CONNECTED}
	c2 := connections.Connection{ID: uuid.FromStringOrNil(conn2ID), Status: connections.DISCONNECTED}

	if err := persister.DB.Create(&c1).Error; err != nil {
		t.Fatalf("failed to seed connection 1: %v", err)
	}
	if err := persister.DB.Create(&c2).Error; err != nil {
		t.Fatalf("failed to seed connection 2: %v", err)
	}

	ctx1 := K8sContext{ID: "ctx-1", Name: "cluster-connected", ConnectionID: conn1ID}
	ctx2 := K8sContext{ID: "ctx-2", Name: "cluster-disconnected", ConnectionID: conn2ID}

	if err := persister.DB.Create(&ctx1).Error; err != nil {
		t.Fatalf("failed to seed context 1: %v", err)
	}
	if err := persister.DB.Create(&ctx2).Error; err != nil {
		t.Fatalf("failed to seed context 2: %v", err)
	}

	// Status = CONNECTED
	resp, err := persister.GetMesheryK8sContextsWithStatus("", "", 0, 10, string(connections.CONNECTED))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	var page MesheryK8sContextPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}
	if page.TotalCount != 1 || len(page.Contexts) != 1 {
		t.Fatalf("expected 1 connected context, got total %d, contexts %d", page.TotalCount, len(page.Contexts))
	}
	if page.Contexts[0].ID != "ctx-1" {
		t.Errorf("expected ctx-1, got %s", page.Contexts[0].ID)
	}

	// Status = "" (all)
	respAll, err := persister.GetMesheryK8sContexts("", "", 0, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	var pageAll MesheryK8sContextPage
	if err := json.Unmarshal(respAll, &pageAll); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}
	if pageAll.TotalCount != 2 {
		t.Fatalf("expected 2 total contexts without status filter, got %d", pageAll.TotalCount)
	}

	// Status = CONNECTED + Search = "connected" + Order = "name asc"
	respSearch, err := persister.GetMesheryK8sContextsWithStatus("connected", "name asc", 0, 10, string(connections.CONNECTED))
	if err != nil {
		t.Fatalf("expected no error for search with status and order, got %v", err)
	}
	var pageSearch MesheryK8sContextPage
	if err := json.Unmarshal(respSearch, &pageSearch); err != nil {
		t.Fatalf("failed to decode page: %v", err)
	}
	if pageSearch.TotalCount != 1 || len(pageSearch.Contexts) != 1 {
		t.Fatalf("expected 1 context matching search and status, got total %d, contexts %d", pageSearch.TotalCount, len(pageSearch.Contexts))
	}
	if pageSearch.Contexts[0].ID != "ctx-1" {
		t.Errorf("expected ctx-1, got %s", pageSearch.Contexts[0].ID)
	}
}
