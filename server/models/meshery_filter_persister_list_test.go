package models

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

func seedListFilter(t *testing.T, db *gorm.DB, name, visibility string) {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	created := filterTestCreatedAt
	f := MesheryFilter{
		ID:         &id,
		Name:       name,
		FilterFile: filterTestFile,
		Visibility: visibility,
		Location:   localFilterLocation(),
		CreatedAt:  &created,
		UpdatedAt:  &created,
	}
	if err := db.Create(&f).Error; err != nil {
		t.Fatalf("seed filter %q: %v", name, err)
	}
}

// seedUndecodableFilter stores a row whose location is not JSON. COUNT never
// decodes a row, so only the Find that loads the page fails on it.
func seedUndecodableFilter(t *testing.T, db *gorm.DB, visibility string) {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	if err := db.Exec(
		`INSERT INTO meshery_filters (id, name, visibility, location) VALUES (?, ?, ?, ?)`,
		id.String(), "corrupt", visibility, "not-json",
	).Error; err != nil {
		t.Fatalf("seed undecodable filter: %v", err)
	}
}

func TestGetMesheryFilters_ReturnsPage(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	seedListFilter(t, db, "wasm-b", Private)
	seedListFilter(t, db, "wasm-a", Private)
	seedListFilter(t, db, "other", Private)

	resp, err := mfp.GetMesheryFilters("WASM", "name asc", 0, 1, nil)
	if err != nil {
		t.Fatalf("GetMesheryFilters: %v", err)
	}
	var page MesheryFilterPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if page.TotalCount != 2 {
		t.Errorf("totalCount = %d, want 2", page.TotalCount)
	}
	if len(page.Filters) != 1 || page.Filters[0].Name != "wasm-a" {
		t.Errorf("filters = %+v, want only wasm-a", page.Filters)
	}
}

func TestGetMesheryFilters_PropagatesCountError(t *testing.T) {
	mfp := &MesheryFilterPersister{DB: newUnmigratedTestDB(t)}

	resp, err := mfp.GetMesheryFilters("", "", 0, 10, nil)
	assertDBReadError(t, resp, err)
}

func TestGetMesheryFilters_PropagatesFindError(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	seedUndecodableFilter(t, db, Private)

	resp, err := mfp.GetMesheryFilters("", "", 0, 10, nil)
	assertDBReadError(t, resp, err)
}

func TestGetMesheryCatalogFilters_ReturnsPublishedOnly(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	seedListFilter(t, db, "published", Published)
	seedListFilter(t, db, "private", Private)

	resp, err := mfp.GetMesheryCatalogFilters("0", "10", "", "")
	if err != nil {
		t.Fatalf("GetMesheryCatalogFilters: %v", err)
	}
	var page FiltersAPIResponse
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if page.TotalCount != 1 || len(page.Filters) != 1 || page.Filters[0].Name != "published" {
		t.Errorf("page = %+v, want only the published filter", page)
	}
}

func TestGetMesheryCatalogFilters_PropagatesCountError(t *testing.T) {
	mfp := &MesheryFilterPersister{DB: newUnmigratedTestDB(t)}

	resp, err := mfp.GetMesheryCatalogFilters("0", "10", "", "")
	assertDBReadError(t, resp, err)
}

func TestGetMesheryCatalogFilters_PropagatesFindError(t *testing.T) {
	// A page size of 0 loads every record through a separate Find.
	for _, pageSize := range []string{"10", "0"} {
		t.Run("pagesize="+pageSize, func(t *testing.T) {
			db, mfp := newFilterPersisterTestDB(t)
			seedUndecodableFilter(t, db, Published)

			resp, err := mfp.GetMesheryCatalogFilters("0", pageSize, "", "")
			assertDBReadError(t, resp, err)
		})
	}
}
