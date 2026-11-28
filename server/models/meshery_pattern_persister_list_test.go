package models

import (
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

// newUnmigratedTestDB returns an in-memory database with no tables, so every
// query a persister runs against it fails.
func newUnmigratedTestDB(t *testing.T) *database.Handler {
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
	t.Cleanup(func() { _ = sqlDB.Close() })
	return &database.Handler{DB: db}
}

// assertDBReadError fails the test unless err is an ErrDBRead and no body was
// returned with it.
func assertDBReadError(t *testing.T, resp []byte, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected an error, got response %s", resp)
	}
	if resp != nil {
		t.Errorf("expected no response body, got %s", resp)
	}
	if code := mkerrors.GetCode(err); code != ErrDBReadCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBReadCode, err)
	}
}

func seedListPattern(t *testing.T, db *gorm.DB, name, visibility string) {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	updated := persisterTestCreatedAt.Add(time.Hour)
	p := MesheryPattern{
		ID:          &id,
		Name:        name,
		PatternFile: persisterTestDesign,
		Visibility:  visibility,
		Location:    isql.Map{"type": "local"},
		CreatedAt:   &updated,
		UpdatedAt:   &updated,
	}
	if err := db.Create(&p).Error; err != nil {
		t.Fatalf("seed pattern %q: %v", name, err)
	}
}

// seedUndecodablePattern stores a row whose location is not JSON. COUNT never
// decodes a row, so only the Find that loads the page fails on it.
func seedUndecodablePattern(t *testing.T, db *gorm.DB, visibility string) {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	if err := db.Exec(
		`INSERT INTO meshery_patterns (id, name, visibility, location, updated_at) VALUES (?, ?, ?, ?, ?)`,
		id.String(), "corrupt", visibility, "not-json", persisterTestCreatedAt,
	).Error; err != nil {
		t.Fatalf("seed undecodable pattern: %v", err)
	}
}

func TestGetMesheryPatterns_ReturnsPage(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	seedListPattern(t, db, "app-b", Private)
	seedListPattern(t, db, "app-a", Private)
	seedListPattern(t, db, "other", Private)

	resp, err := mpp.GetMesheryPatterns("APP", "name asc", 0, 1, "", nil)
	if err != nil {
		t.Fatalf("GetMesheryPatterns: %v", err)
	}
	var page MesheryPatternPage
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if page.TotalCount != 2 {
		t.Errorf("totalCount = %d, want 2", page.TotalCount)
	}
	if len(page.Patterns) != 1 || page.Patterns[0].Name != "app-a" {
		t.Errorf("patterns = %+v, want only app-a", page.Patterns)
	}
}

func TestGetMesheryPatterns_PropagatesCountError(t *testing.T) {
	mpp := &MesheryPatternPersister{DB: newUnmigratedTestDB(t)}

	resp, err := mpp.GetMesheryPatterns("", "", 0, 10, "", nil)
	assertDBReadError(t, resp, err)
}

func TestGetMesheryPatterns_PropagatesFindError(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	seedUndecodablePattern(t, db, Private)

	resp, err := mpp.GetMesheryPatterns("", "", 0, 10, "", nil)
	assertDBReadError(t, resp, err)
}

func TestGetMesheryCatalogPatterns_ReturnsPublishedOnly(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	seedListPattern(t, db, "published", Published)
	seedListPattern(t, db, "private", Private)

	resp, err := mpp.GetMesheryCatalogPatterns("0", "10", "", "")
	if err != nil {
		t.Fatalf("GetMesheryCatalogPatterns: %v", err)
	}
	var page PatternsAPIResponse
	if err := json.Unmarshal(resp, &page); err != nil {
		t.Fatalf("decode page: %v", err)
	}
	if page.TotalCount != 1 || len(page.Patterns) != 1 || page.Patterns[0].Name != "published" {
		t.Errorf("page = %+v, want only the published design", page)
	}
}

func TestGetMesheryCatalogPatterns_PropagatesCountError(t *testing.T) {
	mpp := &MesheryPatternPersister{DB: newUnmigratedTestDB(t)}

	resp, err := mpp.GetMesheryCatalogPatterns("0", "10", "", "")
	assertDBReadError(t, resp, err)
}

func TestGetMesheryCatalogPatterns_PropagatesFindError(t *testing.T) {
	// A page size of 0 loads every record through a separate Find.
	for _, pageSize := range []string{"10", "0"} {
		t.Run("pagesize="+pageSize, func(t *testing.T) {
			db, mpp := newPatternPersisterTestDB(t)
			seedUndecodablePattern(t, db, Published)

			resp, err := mpp.GetMesheryCatalogPatterns("0", pageSize, "", "")
			assertDBReadError(t, resp, err)
		})
	}
}
