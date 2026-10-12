package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
)

func newSMIResultsTestPersister(t *testing.T) *SMIResultsPersister {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&SmiResultWithID{}); err != nil {
		t.Fatalf("migrate SMI results: %v", err)
	}

	return &SMIResultsPersister{DB: &db}
}

func TestSMIResultsPersisterGetResultReadsWrittenRow(t *testing.T) {
	persister := newSMIResultsTestPersister(t)
	id := uuid.Must(uuid.NewV4())
	payload, err := json.Marshal(&SmiResult{ID: id, MeshName: "test-mesh"})
	if err != nil {
		t.Fatalf("marshal SMI result: %v", err)
	}
	if err := persister.WriteResult(id, payload); err != nil {
		t.Fatalf("write result: %v", err)
	}

	gotJSON, err := persister.GetResult(id)
	if err != nil {
		t.Fatalf("get result: %v", err)
	}
	var got SmiResult
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if got.ID != id || got.MeshName != "test-mesh" {
		t.Fatalf("result = {ID: %s, MeshName: %q}, want {ID: %s, MeshName: %q}", got.ID, got.MeshName, id, "test-mesh")
	}

	var listed []SmiResultWithID
	if err := persister.DB.Find(&listed).Error; err != nil {
		t.Fatalf("list persisted results: %v", err)
	}
	if len(listed) != 1 {
		t.Fatalf("listed result count = %d, want 1", len(listed))
	}
	if listed[0].ID != id || listed[0].MeshName != got.MeshName {
		t.Fatalf("listed result = {ID: %s, MeshName: %q}, single result = {ID: %s, MeshName: %q}", listed[0].ID, listed[0].MeshName, got.ID, got.MeshName)
	}
}
