package models

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestSMIResultsPersister(t *testing.T) *SMIResultsPersister {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}

	if err := db.AutoMigrate(&SmiResultWithID{}); err != nil {
		t.Fatalf("failed to migrate smi result schema: %v", err)
	}

	return &SMIResultsPersister{
		DB: &database.Handler{DB: db},
	}
}

func mustCoreUUID() core.Uuid {
	return core.Uuid(uuid.Must(uuid.NewV4()))
}

func TestSMIResultsPersister_WriteAndGetResult(t *testing.T) {
	persister := newTestSMIResultsPersister(t)

	key := mustCoreUUID()
	result := &SmiResult{
		MeshName: "istio",
		Status:   "passed",
	}
	resultBytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("failed to marshal smi result: %v", err)
	}

	if err := persister.WriteResult(key, resultBytes); err != nil {
		t.Fatalf("failed to write smi result: %v", err)
	}

	gotBytes, err := persister.GetResult(key)
	if err != nil {
		t.Fatalf("failed to get smi result: %v", err)
	}

	var got SmiResult
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("failed to unmarshal retrieved smi result: %v", err)
	}

	if got.MeshName != "istio" || got.Status != "passed" {
		t.Errorf("got SmiResult = %+v, want meshName=istio, status=passed", got)
	}

	// Update the result
	result.Status = "failed"
	updatedBytes, _ := json.Marshal(result)
	if err := persister.WriteResult(key, updatedBytes); err != nil {
		t.Fatalf("failed to update smi result: %v", err)
	}

	gotBytes, err = persister.GetResult(key)
	if err != nil {
		t.Fatalf("failed to get updated smi result: %v", err)
	}
	if err := json.Unmarshal(gotBytes, &got); err != nil {
		t.Fatalf("failed to unmarshal updated smi result: %v", err)
	}
	if got.Status != "failed" {
		t.Errorf("expected updated status failed, got %s", got.Status)
	}
}

func TestSMIResultsPersister_GetResults_Pagination(t *testing.T) {
	persister := newTestSMIResultsPersister(t)

	for i := 0; i < 5; i++ {
		key := mustCoreUUID()
		res := &SmiResult{MeshName: "linkerd"}
		bytes, _ := json.Marshal(res)
		if err := persister.WriteResult(key, bytes); err != nil {
			t.Fatalf("failed to seed smi result: %v", err)
		}
	}

	pageBytes, err := persister.GetResults(0, 3)
	if err != nil {
		t.Fatalf("failed to get results page: %v", err)
	}

	var page SmiResultPage
	if err := json.Unmarshal(pageBytes, &page); err != nil {
		t.Fatalf("failed to unmarshal result page: %v", err)
	}

	if page.TotalCount != 5 {
		t.Errorf("got TotalCount=%d, want 5", page.TotalCount)
	}
	if len(page.Results) != 3 {
		t.Errorf("got %d results in page, want 3", len(page.Results))
	}
}

func TestSMIResultsPersister_DeleteResult_SuccessAndNotFound(t *testing.T) {
	persister := newTestSMIResultsPersister(t)

	key := mustCoreUUID()
	res := &SmiResult{MeshName: "consul"}
	bytes, _ := json.Marshal(res)
	if err := persister.WriteResult(key, bytes); err != nil {
		t.Fatalf("failed to write smi result: %v", err)
	}

	// Delete existing
	if err := persister.DeleteResult(key); err != nil {
		t.Fatalf("failed to delete smi result: %v", err)
	}

	// Verify get after delete returns ErrResultNotFound
	_, err := persister.GetResult(key)
	if err == nil {
		t.Fatal("expected ErrResultNotFound after delete, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("expected ErrResultNotFoundCode (%s), got: %s", ErrResultNotFoundCode, code)
	}

	// Deleting again should return ErrResultNotFound
	err = persister.DeleteResult(key)
	if err == nil {
		t.Fatal("expected ErrResultNotFound when deleting nonexistent record, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrResultNotFoundCode {
		t.Errorf("expected ErrResultNotFoundCode (%s), got: %s", ErrResultNotFoundCode, code)
	}
}

func TestSMIResultsPersister_NilDB(t *testing.T) {
	persister := &SMIResultsPersister{DB: nil}
	key := mustCoreUUID()

	if _, err := persister.GetResults(0, 10); !errors.Is(err, ErrDBConnection) {
		t.Errorf("expected ErrDBConnection on GetResults, got: %v", err)
	}

	if err := persister.WriteResult(key, []byte("{}")); !errors.Is(err, ErrDBConnection) {
		t.Errorf("expected ErrDBConnection on WriteResult, got: %v", err)
	}

	if err := persister.DeleteResult(key); !errors.Is(err, ErrDBConnection) {
		t.Errorf("expected ErrDBConnection on DeleteResult, got: %v", err)
	}

	if _, err := persister.GetResult(key); !errors.Is(err, ErrDBConnection) {
		t.Errorf("expected ErrDBConnection on GetResult, got: %v", err)
	}
}
