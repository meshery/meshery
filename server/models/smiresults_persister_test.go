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

func marshalSMIResult(t *testing.T, meshName string) []byte {
	t.Helper()

	payload, err := json.Marshal(&SmiResult{MeshName: meshName})
	if err != nil {
		t.Fatalf("marshal SMI result: %v", err)
	}
	return payload
}

func TestSMIResultsPersisterWriteResultUpdatesExistingRow(t *testing.T) {
	persister := newSMIResultsTestPersister(t)
	existingID := uuid.Must(uuid.NewV4())

	if err := persister.WriteResult(existingID, marshalSMIResult(t, "old")); err != nil {
		t.Fatalf("write initial result: %v", err)
	}
	if err := persister.WriteResult(existingID, marshalSMIResult(t, "new")); err != nil {
		t.Fatalf("update result: %v", err)
	}

	var updated SmiResultWithID
	if err := persister.DB.Where("id = ?", existingID).First(&updated).Error; err != nil {
		t.Fatalf("read updated result: %v", err)
	}
	if updated.MeshName != "new" {
		t.Fatalf("updated mesh name = %q, want %q", updated.MeshName, "new")
	}

	newID := uuid.Must(uuid.NewV4())
	if err := persister.WriteResult(newID, marshalSMIResult(t, "inserted")); err != nil {
		t.Fatalf("write new result: %v", err)
	}
	var inserted SmiResultWithID
	if err := persister.DB.Where("id = ?", newID).First(&inserted).Error; err != nil {
		t.Fatalf("read inserted result: %v", err)
	}
	if inserted.MeshName != "inserted" {
		t.Fatalf("inserted mesh name = %q, want %q", inserted.MeshName, "inserted")
	}
}
