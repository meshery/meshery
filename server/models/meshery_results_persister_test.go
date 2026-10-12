package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
)

func newMesheryResultsTestPersister(t *testing.T) *MesheryResultsPersister {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&PerformanceProfile{}, &MesheryResult{}); err != nil {
		t.Fatalf("migrate meshery results: %v", err)
	}

	return &MesheryResultsPersister{DB: &db}
}

func writeMesheryResultFixture(t *testing.T, persister *MesheryResultsPersister, id uuid.UUID, name string) {
	t.Helper()

	payload, err := json.Marshal(&MesheryResult{
		Name:   name,
		Result: map[string]interface{}{"name": name},
	})
	if err != nil {
		t.Fatalf("marshal result fixture: %v", err)
	}
	if err := persister.WriteResult(id, payload); err != nil {
		t.Fatalf("write result fixture: %v", err)
	}
}

func TestMesheryResultsPersisterGetResultFiltersByID(t *testing.T) {
	persister := newMesheryResultsTestPersister(t)
	firstID := uuid.Must(uuid.NewV4())
	secondID := uuid.Must(uuid.NewV4())
	writeMesheryResultFixture(t, persister, firstID, "first")
	writeMesheryResultFixture(t, persister, secondID, "second")

	for _, test := range []struct {
		name string
		id   uuid.UUID
	}{
		{name: "first result", id: firstID},
		{name: "second result", id: secondID},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := persister.GetResult(test.id, nil)
			if err != nil {
				t.Fatalf("GetResult(%s): %v", test.id, err)
			}
			if result.ID != test.id {
				t.Fatalf("result ID = %s, want %s", result.ID, test.id)
			}
		})
	}

	missingID := uuid.Must(uuid.NewV4())
	result, err := persister.GetResult(missingID, nil)
	if err == nil {
		t.Fatalf("GetResult(%s) returned result %s without an error", missingID, result.ID)
	}
	if result != nil {
		t.Fatalf("GetResult(%s) result = %#v, want nil", missingID, result)
	}
}
