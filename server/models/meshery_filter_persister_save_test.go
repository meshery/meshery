package models

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	isql "github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	filterTestCreatedAt = time.Date(2026, 11, 5, 12, 32, 10, 0, time.UTC)
	filterTestFile      = []byte("\x00asm wasm binary")
)

func newFilterPersisterTestDB(t *testing.T) (*gorm.DB, *MesheryFilterPersister) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	// Every pooled connection to ":memory:" is a separate database.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&MesheryFilter{}); err != nil {
		t.Fatalf("migrate meshery_filters: %v", err)
	}
	return db, &MesheryFilterPersister{DB: &database.Handler{DB: db}}
}

// localFilterLocation is the location POST /api/filter assigns when a request
// carries none.
func localFilterLocation() isql.Map {
	return isql.Map{"host": "", "path": "", "type": "local", "branch": ""}
}

// seedPublishedFilter stores a filter as an upload followed by publishing
// leaves it: a WASM file, published visibility, catalog data and a creation time.
func seedPublishedFilter(t *testing.T, db *gorm.DB) uuid.UUID {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	created := filterTestCreatedAt
	seed := MesheryFilter{
		ID:          &id,
		Name:        "stored",
		FilterFile:  filterTestFile,
		Visibility:  Published,
		CatalogData: isql.Map{"info": "stored info"},
		Location:    localFilterLocation(),
		CreatedAt:   &created,
		UpdatedAt:   &created,
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed filter: %v", err)
	}
	return id
}

// renameRequest is the record POST /api/filter builds when the UI renames a
// filter: the ID, the new name, empty catalog data and the default location.
func renameRequest(id uuid.UUID, name string) *MesheryFilter {
	return &MesheryFilter{ID: &id, Name: name, CatalogData: isql.Map{}, Location: localFilterLocation()}
}

func loadFilter(t *testing.T, db *gorm.DB, id uuid.UUID) MesheryFilter {
	t.Helper()
	var got MesheryFilter
	if err := db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load filter %s: %v", id, err)
	}
	return got
}

func countFilters(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&MesheryFilter{}).Count(&n).Error; err != nil {
		t.Fatalf("count filters: %v", err)
	}
	return n
}

// TestSaveMesheryFilterUpdateKeepsOmittedColumns renames a stored filter the way
// the UI does and checks that the columns the request left out keep their values.
func TestSaveMesheryFilterUpdateKeepsOmittedColumns(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := seedPublishedFilter(t, db)

	if _, err := mfp.SaveMesheryFilter(renameRequest(id, "renamed")); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadFilter(t, db, id)
	if got.Name != "renamed" {
		t.Errorf("name = %q, want %q", got.Name, "renamed")
	}
	if !bytes.Equal(got.FilterFile, filterTestFile) {
		t.Errorf("filter_file = %q, want the stored binary %q to survive the rename", got.FilterFile, filterTestFile)
	}
	if got.Visibility != Published {
		t.Errorf("visibility = %q, want the filter to stay %q", got.Visibility, Published)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(filterTestCreatedAt) {
		t.Errorf("created_at = %v, want %v", got.CreatedAt, filterTestCreatedAt)
	}
	if got.CatalogData["info"] != "stored info" {
		t.Errorf("catalog_data = %v, want the stored catalog data", got.CatalogData)
	}
	if n := countFilters(t, db); n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

// TestSaveMesheryFilterResponseShowsStoredFilter checks that the save response
// describes the stored filter, not the partial request.
func TestSaveMesheryFilterResponseShowsStoredFilter(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := seedPublishedFilter(t, db)

	resp, err := mfp.SaveMesheryFilter(renameRequest(id, "renamed"))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	var saved []MesheryFilter
	if err := json.Unmarshal(resp, &saved); err != nil || len(saved) != 1 {
		t.Fatalf("decode response %s: %v", resp, err)
	}
	got := saved[0]
	if got.Name != "renamed" {
		t.Errorf("response name = %q, want %q", got.Name, "renamed")
	}
	if !bytes.Equal(got.FilterFile, filterTestFile) {
		t.Errorf("response filterFile = %q, want the stored binary %q", got.FilterFile, filterTestFile)
	}
	if got.Visibility != Published {
		t.Errorf("response visibility = %q, want %q", got.Visibility, Published)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(filterTestCreatedAt) {
		t.Errorf("response createdAt = %v, want %v", got.CreatedAt, filterTestCreatedAt)
	}
}

// TestSaveMesheryFilterUpdateAppliesSentFields checks that fields the request
// does carry still replace the stored values.
func TestSaveMesheryFilterUpdateAppliesSentFields(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := seedPublishedFilter(t, db)
	newFile := []byte("\x00asm new binary")

	record := renameRequest(id, "renamed")
	record.FilterFile = newFile
	record.Visibility = Private
	record.CatalogData = isql.Map{"info": "new info"}
	if _, err := mfp.SaveMesheryFilter(record); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadFilter(t, db, id)
	if !bytes.Equal(got.FilterFile, newFile) {
		t.Errorf("filter_file = %q, want %q", got.FilterFile, newFile)
	}
	if got.Visibility != Private {
		t.Errorf("visibility = %q, want %q", got.Visibility, Private)
	}
	if got.CatalogData["info"] != "new info" {
		t.Errorf("catalog_data = %v, want the new catalog data", got.CatalogData)
	}
}

// TestSaveMesheryFilterUnknownIDCreatesPrivateFilter checks that a save under an
// ID the caller chose, which is not stored yet, creates a private filter.
func TestSaveMesheryFilterUnknownIDCreatesPrivateFilter(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := uuid.Must(uuid.NewV4())

	record := renameRequest(id, "new")
	record.FilterFile = filterTestFile
	if _, err := mfp.SaveMesheryFilter(record); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadFilter(t, db, id)
	if got.Visibility != Private {
		t.Errorf("visibility = %q, want %q", got.Visibility, Private)
	}
	if !bytes.Equal(got.FilterFile, filterTestFile) {
		t.Errorf("filter_file = %q, want %q", got.FilterFile, filterTestFile)
	}
	if n := countFilters(t, db); n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

// A concurrent first save can insert the same ID between the existence check
// and the insert. The losing save must then update the winner's row like any
// existing row: no zero-valued overwrite and no visibility default.
func TestSaveMesheryFilterConcurrentFirstSaveKeepsWinnerFields(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := uuid.Must(uuid.NewV4())

	fired := false
	if err := db.Callback().Create().Before("gorm:create").Register("test:concurrent_insert", func(tx *gorm.DB) {
		if fired {
			return
		}
		fired = true
		tx.Session(&gorm.Session{NewDB: true}).Exec(
			"INSERT INTO meshery_filters (id, name, filter_file, visibility, created_at) VALUES (?, ?, ?, ?, ?)",
			id, "winner", filterTestFile, Published, filterTestCreatedAt)
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	if _, err := mfp.SaveMesheryFilter(renameRequest(id, "renamed")); err != nil {
		t.Fatalf("save: %v", err)
	}
	if !fired {
		t.Fatalf("insert path was not taken")
	}

	got := loadFilter(t, db, id)
	if got.Name != "renamed" {
		t.Errorf("name = %q, want the losing save's name applied", got.Name)
	}
	if !bytes.Equal(got.FilterFile, filterTestFile) {
		t.Errorf("filter_file = %q, want the winner's file", got.FilterFile)
	}
	if got.Visibility != Published {
		t.Errorf("visibility = %q, want %q", got.Visibility, Published)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(filterTestCreatedAt) {
		t.Errorf("created_at = %v, want %v", got.CreatedAt, filterTestCreatedAt)
	}
	if n := countFilters(t, db); n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

func TestSaveMesheryFilterDeletedDuringSaveIsError(t *testing.T) {
	db, mfp := newFilterPersisterTestDB(t)
	id := seedPublishedFilter(t, db)

	if err := db.Callback().Update().Before("gorm:update").Register("test:concurrent_delete", func(tx *gorm.DB) {
		tx.Session(&gorm.Session{NewDB: true}).Exec("DELETE FROM meshery_filters WHERE id = ?", id)
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	_, err := mfp.SaveMesheryFilter(renameRequest(id, "renamed"))
	if code := mkerrors.GetCode(err); code != ErrDBPutCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBPutCode, err)
	}
	if n := countFilters(t, db); n != 0 {
		t.Errorf("rows = %d, want the deleted filter to stay deleted", n)
	}
}
