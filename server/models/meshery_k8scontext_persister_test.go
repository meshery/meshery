package models

import (
	"encoding/json"
	"testing"
)

func TestGetMesheryK8sContexts(t *testing.T) {
	db := newMigratedDB(t)
	mkcp := &MesheryK8sContextPersister{DB: db}

	// 1. Initially empty
	data, err := mkcp.GetMesheryK8sContexts("", "", 0, 10)
	if err != nil {
		t.Fatalf("unexpected error on empty table: %v", err)
	}
	var page MesheryK8sContextPage
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if page.TotalCount != 0 || len(page.Contexts) != 0 {
		t.Fatalf("expected 0 contexts, got totalCount=%d len=%d", page.TotalCount, len(page.Contexts))
	}

	// 2. Add contexts with explicit IDs
	ctx1 := K8sContext{
		ID:     "ctx-alpha",
		Name:   "cluster-alpha",
		Server: "https://k8s-alpha.example.com",
	}
	ctx2 := K8sContext{
		ID:     "ctx-beta",
		Name:   "cluster-beta",
		Server: "https://k8s-beta.example.com",
	}
	if _, err := mkcp.SaveMesheryK8sContext(ctx1); err != nil {
		t.Fatalf("save context 1: %v", err)
	}
	if _, err := mkcp.SaveMesheryK8sContext(ctx2); err != nil {
		t.Fatalf("save context 2: %v", err)
	}

	// 3. Query all
	data, err = mkcp.GetMesheryK8sContexts("", "", 0, 10)
	if err != nil {
		t.Fatalf("unexpected error fetching contexts: %v", err)
	}
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if page.TotalCount != 2 || len(page.Contexts) != 2 {
		t.Fatalf("expected 2 contexts, got totalCount=%d len=%d", page.TotalCount, len(page.Contexts))
	}

	// 4. Search filter
	data, err = mkcp.GetMesheryK8sContexts("alpha", "", 0, 10)
	if err != nil {
		t.Fatalf("unexpected error searching contexts: %v", err)
	}
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if page.TotalCount != 1 || len(page.Contexts) != 1 || page.Contexts[0].Name != "cluster-alpha" {
		t.Fatalf("expected 1 context matching 'alpha', got %+v", page)
	}
}

func TestGetMesheryK8sContexts_DBError(t *testing.T) {
	db := newMigratedDB(t)
	mkcp := &MesheryK8sContextPersister{DB: db}

	// Drop the k8s_contexts table to simulate a database query failure
	if err := db.Migrator().DropTable(&K8sContext{}); err != nil {
		t.Fatalf("drop table: %v", err)
	}

	// Calling GetMesheryK8sContexts must now return the database error rather than swallowing it
	_, err := mkcp.GetMesheryK8sContexts("", "", 0, 10)
	if err == nil {
		t.Fatal("expected error from GetMesheryK8sContexts when table does not exist, got nil")
	}
}
