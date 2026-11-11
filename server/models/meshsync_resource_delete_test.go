package models

import (
	"testing"

	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

// countMeshSyncRows returns how many rows each MeshSync table holds for id.
func countMeshSyncRows(t *testing.T, db *database.Handler, id string) map[string]int64 {
	t.Helper()
	tables := map[string]interface{}{
		"resource":    &meshsyncmodel.KubernetesResource{},
		"object meta": &meshsyncmodel.KubernetesResourceObjectMeta{},
		"spec":        &meshsyncmodel.KubernetesResourceSpec{},
		"status":      &meshsyncmodel.KubernetesResourceStatus{},
		"key value":   &meshsyncmodel.KubernetesKeyValue{},
	}
	counts := make(map[string]int64, len(tables))
	for name, model := range tables {
		var n int64
		if err := db.Model(model).Where("id = ?", id).Count(&n).Error; err != nil {
			t.Fatalf("counting %s rows for %s: %v", name, id, err)
		}
		counts[name] = n
	}
	return counts
}

func assertMeshSyncRows(t *testing.T, db *database.Handler, id string, want int64) {
	t.Helper()
	for name, got := range countMeshSyncRows(t, db, id) {
		if got != want {
			t.Errorf("%s rows for %s: got %d, want %d", name, id, got, want)
		}
	}
}

// TestDeleteMeshSyncResourceByID guards the single-resource delete behind
// DELETE /api/system/meshsync/resources/{id}: it must remove the resource and
// every child row sharing its id (object meta, spec, status, labels and
// annotations), and leave other resources untouched. Before the fix only the
// parent row was deleted, so the children were orphaned and reattached to the
// resource when MeshSync discovered it again under the same id.
func TestDeleteMeshSyncResourceByID(t *testing.T) {
	db := newMeshSyncTestDB(t)

	const (
		deletedID = "res-deleted"
		keptID    = "res-kept"
	)
	seedMeshSyncResource(t, db, deletedID, "cluster-a")
	seedMeshSyncResource(t, db, keptID, "cluster-a")

	deleted, err := DeleteMeshSyncResourceByID(db, deletedID)
	if err != nil {
		t.Fatalf("DeleteMeshSyncResourceByID returned error: %v", err)
	}
	if !deleted {
		t.Fatal("expected the resource to be reported as deleted")
	}

	assertMeshSyncRows(t, db, deletedID, 0)
	assertMeshSyncRows(t, db, keptID, 1)
}

// TestDeleteMeshSyncResourceByIDNotFound verifies that an unknown id is reported
// as not deleted, so the handler can answer 404, and that nothing else is
// removed.
func TestDeleteMeshSyncResourceByIDNotFound(t *testing.T) {
	db := newMeshSyncTestDB(t)

	const keptID = "res-kept"
	seedMeshSyncResource(t, db, keptID, "cluster-a")

	deleted, err := DeleteMeshSyncResourceByID(db, "res-missing")
	if err != nil {
		t.Fatalf("DeleteMeshSyncResourceByID returned error: %v", err)
	}
	if deleted {
		t.Error("expected an unknown id to be reported as not deleted")
	}

	assertMeshSyncRows(t, db, keptID, 1)
}

// TestDeleteMeshSyncResourceByIDNilHandler verifies the helper reports the
// structured MeshKit error instead of panicking on an unusable handler.
func TestDeleteMeshSyncResourceByIDNilHandler(t *testing.T) {
	cases := map[string]*database.Handler{
		"nil handler":     nil,
		"nil embedded DB": {DB: nil},
	}
	for name, db := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := DeleteMeshSyncResourceByID(db, "any-id")
			if err == nil {
				t.Fatal("expected an error for an unusable database handler, got nil")
			}
			if code := mkerrors.GetCode(err); code != ErrEmptyMeshSyncHandlerCode {
				t.Errorf("expected error code %s, got %s", ErrEmptyMeshSyncHandlerCode, code)
			}
		})
	}
}

// TestDeleteMeshSyncResourceByIDRollsBackOnError verifies the delete is atomic:
// if a child delete fails, the resource and its other child rows are kept
// rather than leaving a partially deleted resource behind.
func TestDeleteMeshSyncResourceByIDRollsBackOnError(t *testing.T) {
	db := newMeshSyncTestDB(t)

	const id = "res-atomic"
	seedMeshSyncResource(t, db, id, "cluster-a")

	if err := db.Migrator().DropTable(&meshsyncmodel.KubernetesResourceObjectMeta{}); err != nil {
		t.Fatalf("failed to drop object-meta table: %v", err)
	}

	if _, err := DeleteMeshSyncResourceByID(db, id); err == nil {
		t.Fatal("expected an error when a child delete fails, got nil")
	}

	for name, model := range map[string]interface{}{
		"resource":  &meshsyncmodel.KubernetesResource{},
		"spec":      &meshsyncmodel.KubernetesResourceSpec{},
		"status":    &meshsyncmodel.KubernetesResourceStatus{},
		"key value": &meshsyncmodel.KubernetesKeyValue{},
	} {
		var n int64
		if err := db.Model(model).Where("id = ?", id).Count(&n).Error; err != nil {
			t.Fatalf("counting %s rows: %v", name, err)
		}
		if n != 1 {
			t.Errorf("%s rows after a failed delete: got %d, want 1 (rolled back)", name, n)
		}
	}
}
