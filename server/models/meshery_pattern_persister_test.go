package models

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	isql "github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshkit/database"
	mkerrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/catalog/v1alpha1"
	"github.com/meshery/meshkit/models/patterns"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const persisterTestDesign = `name: persister-test
id: 00000000-0000-0000-0000-000000000000
schemaVersion: designs.meshery.io/v1beta1
version: 0.0.1
components: []
relationships: []
`

var persisterTestCreatedAt = time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)

func newPatternPersisterTestDB(t *testing.T) (*gorm.DB, *MesheryPatternPersister) {
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
	if err := db.AutoMigrate(&MesheryPattern{}); err != nil {
		t.Fatalf("migrate meshery_patterns: %v", err)
	}
	return db, &MesheryPatternPersister{DB: &database.Handler{DB: db}}
}

// seedImportedPattern stores a design as an import leaves it: original source,
// source type, a creation time, visibility and catalog data.
func seedImportedPattern(t *testing.T, db *gorm.DB) (uuid.UUID, MesheryPattern) {
	t.Helper()
	id := uuid.Must(uuid.NewV4())
	created := persisterTestCreatedAt
	seed := MesheryPattern{
		ID:            &id,
		Name:          "persister-test",
		PatternFile:   persisterTestDesign,
		Visibility:    Public,
		CatalogData:   v1alpha1.CatalogData{PatternInfo: "stored info", PatternCaveats: "stored caveats"},
		Type:          sql.NullString{String: string(K8sManifest), Valid: true},
		SourceContent: []byte("apiVersion: v1\nkind: Namespace\nmetadata:\n  name: demo\n"),
		Location:      isql.Map{"host": "github.com", "path": "demo.yaml"},
		CreatedAt:     &created,
		UpdatedAt:     &created,
		ViewCount:     5,
		CloneCount:    2,
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed pattern: %v", err)
	}
	return id, seed
}

func loadPattern(t *testing.T, db *gorm.DB, id uuid.UUID) MesheryPattern {
	t.Helper()
	var got MesheryPattern
	if err := db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load pattern %s: %v", id, err)
	}
	return got
}

func countPatterns(t *testing.T, db *gorm.DB) int64 {
	t.Helper()
	var n int64
	if err := db.Model(&MesheryPattern{}).Count(&n).Error; err != nil {
		t.Fatalf("count patterns: %v", err)
	}
	return n
}

// postStylePattern is the record handlePatternPOST builds for an existing
// design: only the fields DesignPostPayload carries.
func postStylePattern(id uuid.UUID) *MesheryPattern {
	return &MesheryPattern{
		ID:          &id,
		Name:        "persister-test",
		PatternFile: persisterTestDesign,
	}
}

func TestSaveMesheryPattern_PartialSavePreservesSourceContentAndType(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, seed := seedImportedPattern(t, db)

	if _, err := mpp.SaveMesheryPattern(postStylePattern(id)); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if string(got.SourceContent) != string(seed.SourceContent) {
		t.Errorf("source_content = %q, want %q", got.SourceContent, seed.SourceContent)
	}
	if got.Type != seed.Type {
		t.Errorf("type = %+v, want %+v", got.Type, seed.Type)
	}
	if got.Location["path"] != "demo.yaml" {
		t.Errorf("location = %v, want the stored location", got.Location)
	}
}

func TestSaveMesheryPattern_PartialSavePreservesCreatedAt(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, _ := seedImportedPattern(t, db)

	if _, err := mpp.SaveMesheryPattern(postStylePattern(id)); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if got.CreatedAt == nil || !got.CreatedAt.Equal(persisterTestCreatedAt) {
		t.Errorf("created_at = %v, want %v", got.CreatedAt, persisterTestCreatedAt)
	}
	if got.UpdatedAt == nil || !got.UpdatedAt.After(persisterTestCreatedAt) {
		t.Errorf("updated_at = %v, want it advanced past %v", got.UpdatedAt, persisterTestCreatedAt)
	}
}

// The design editor's save sends neither visibility nor catalogData.
func TestSaveMesheryPattern_PartialSavePreservesVisibilityAndCatalogData(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, seed := seedImportedPattern(t, db)

	if _, err := mpp.SaveMesheryPattern(postStylePattern(id)); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if got.Visibility != seed.Visibility {
		t.Errorf("visibility = %q, want %q", got.Visibility, seed.Visibility)
	}
	if got.CatalogData.PatternInfo != seed.CatalogData.PatternInfo ||
		got.CatalogData.PatternCaveats != seed.CatalogData.PatternCaveats {
		t.Errorf("catalog_data = %+v, want %+v", got.CatalogData, seed.CatalogData)
	}
}

func TestSaveMesheryPattern_AppliesSuppliedFields(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, seed := seedImportedPattern(t, db)

	update := postStylePattern(id)
	update.Name = "renamed"
	update.Visibility = Private
	update.CatalogData = v1alpha1.CatalogData{PatternInfo: "new info"}
	if _, err := mpp.SaveMesheryPattern(update); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if got.Name != "renamed" {
		t.Errorf("name = %q, want %q", got.Name, "renamed")
	}
	if got.Visibility != Private {
		t.Errorf("visibility = %q, want %q", got.Visibility, Private)
	}
	if got.CatalogData.PatternInfo != "new info" || got.CatalogData.PatternCaveats != "" {
		t.Errorf("catalog_data = %+v, want the supplied catalog data written whole", got.CatalogData)
	}
	if got.PatternFile == seed.PatternFile {
		t.Errorf("pattern_file was not rewritten with the next version")
	}
	if string(got.SourceContent) != string(seed.SourceContent) {
		t.Errorf("source_content changed on a save that did not carry it")
	}
}

// Conversion re-saves a full record read earlier. Usage counters belong to the
// row, so a stale copy in the struct must not overwrite a newer stored value.
func TestSaveMesheryPattern_FullRecordResaveKeepsStoredCounters(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, _ := seedImportedPattern(t, db)

	stale := loadPattern(t, db, id)
	if err := db.Model(&MesheryPattern{}).Where("id = ?", id).Update("view_count", 9).Error; err != nil {
		t.Fatalf("bump view_count: %v", err)
	}
	stale.PatternFile = persisterTestDesign
	if _, err := mpp.SaveMesheryPattern(&stale); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if got.ViewCount != 9 {
		t.Errorf("view_count = %d, want 9 (a save must not write counters)", got.ViewCount)
	}
	if string(got.SourceContent) != string(stale.SourceContent) {
		t.Errorf("source_content changed on a full-record resave")
	}
}

// Saving a record whose stored values already match must update in place, not
// fall through to an insert.
func TestSaveMesheryPattern_UnchangedResaveUpdatesInPlace(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, _ := seedImportedPattern(t, db)

	for i := 0; i < 2; i++ {
		record := loadPattern(t, db, id)
		if _, err := mpp.SaveMesheryPattern(&record); err != nil {
			t.Fatalf("save %d: %v", i, err)
		}
	}
	if n := countPatterns(t, db); n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

func TestSaveMesheryPattern_UnknownIDInserts(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id := uuid.Must(uuid.NewV4())

	if _, err := mpp.SaveMesheryPattern(postStylePattern(id)); err != nil {
		t.Fatalf("save: %v", err)
	}

	got := loadPattern(t, db, id)
	if got.Visibility != Private {
		t.Errorf("visibility = %q, want %q for a new design", got.Visibility, Private)
	}
	if got.CreatedAt == nil {
		t.Errorf("created_at not set on insert")
	}
}

// Import saves a new design (no ID) carrying its source type and content.
func TestSaveMesheryPattern_NewDesignKeepsImportedSource(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	record := &MesheryPattern{
		Name:          "imported",
		PatternFile:   persisterTestDesign,
		Type:          sql.NullString{String: string(HelmChart), Valid: true},
		SourceContent: []byte("chart-bytes"),
	}

	if _, err := mpp.SaveMesheryPattern(record); err != nil {
		t.Fatalf("save: %v", err)
	}
	if record.ID == nil {
		t.Fatalf("save did not assign an ID")
	}

	got := loadPattern(t, db, *record.ID)
	if got.Type.String != string(HelmChart) || string(got.SourceContent) != "chart-bytes" {
		t.Errorf("type/source = %q/%q, want %q/%q", got.Type.String, got.SourceContent, HelmChart, "chart-bytes")
	}
}

// A clone is a new design: it takes the source's content, but not its creation
// time, usage counters, catalog data or design-file identity.
func TestCloneMesheryPattern_InsertsNewRow(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, seed := seedImportedPattern(t, db)
	if err := db.Model(&MesheryPattern{}).Where("id = ?", id).
		Updates(map[string]interface{}{"share_count": 3, "download_count": 7, "deployment_count": 4}).Error; err != nil {
		t.Fatalf("set counters: %v", err)
	}

	if _, err := mpp.CloneMesheryPattern(id.String(), &MesheryClonePatternRequestBody{Name: "copy"}); err != nil {
		t.Fatalf("clone: %v", err)
	}

	if n := countPatterns(t, db); n != 2 {
		t.Fatalf("rows = %d, want 2", n)
	}
	var clone MesheryPattern
	if err := db.First(&clone, "name = ?", "copy").Error; err != nil {
		t.Fatalf("load clone: %v", err)
	}
	if clone.Visibility != Private {
		t.Errorf("clone visibility = %q, want %q", clone.Visibility, Private)
	}
	if clone.CreatedAt == nil || clone.CreatedAt.Equal(persisterTestCreatedAt) {
		t.Errorf("clone created_at = %v, want the clone's own creation time", clone.CreatedAt)
	}
	if counts := [5]int{clone.ViewCount, clone.ShareCount, clone.DownloadCount, clone.CloneCount, clone.DeploymentCount}; counts != [5]int{} {
		t.Errorf("clone counters (view, share, download, clone, deploy) = %v, want all zero", counts)
	}
	if clone.CatalogData.PatternInfo != "" || clone.CatalogData.PatternCaveats != "" {
		t.Errorf("clone catalog_data = %+v, want empty", clone.CatalogData)
	}
	if string(clone.SourceContent) != string(seed.SourceContent) || clone.Type != seed.Type {
		t.Errorf("clone type/source = %q/%q, want the source's", clone.Type.String, clone.SourceContent)
	}
	pf, err := patterns.GetPatternFormat(clone.PatternFile)
	if err != nil {
		t.Fatalf("parse clone design file: %v", err)
	}
	if pf.ID != *clone.ID || pf.Name != "copy" || pf.Version != "0.0.1" {
		t.Errorf("clone design file id/name/version = %s/%q/%s, want %s/%q/0.0.1", pf.ID, pf.Name, pf.Version, clone.ID, "copy")
	}
	if original := loadPattern(t, db, id); original.Name != "persister-test" || original.PatternFile != seed.PatternFile {
		t.Errorf("original changed by the clone: name %q, design file %q", original.Name, original.PatternFile)
	}
}

func TestReplaceSeededPatterns_InsertsSeededDesigns(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)

	seeded := make([]*MesheryPattern, 0, 2)
	for _, name := range []string{"seed-a", "seed-b"} {
		id := uuid.Must(uuid.NewV4())
		seeded = append(seeded, &MesheryPattern{ID: &id, Name: name, PatternFile: persisterTestDesign, Visibility: Published})
	}
	for pass := 0; pass < 2; pass++ {
		if err := mpp.ReplaceSeededPatterns(seeded); err != nil {
			t.Fatalf("replace seeded patterns (pass %d): %v", pass, err)
		}
	}

	if n := countPatterns(t, db); n != 2 {
		t.Errorf("rows = %d, want 2", n)
	}
}

func TestSaveMesheryPattern_DatabaseErrorIsMeshKitError(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, _ := seedImportedPattern(t, db)
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	_ = sqlDB.Close()

	_, err = mpp.SaveMesheryPattern(postStylePattern(id))
	if code := mkerrors.GetCode(err); code != ErrDBReadCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBReadCode, err)
	}
}

// A concurrent first save can insert the same ID between the existence check
// and the insert. The losing save must then update the winner's row like any
// existing row: no zero-valued overwrite and no visibility default.
func TestSaveMesheryPattern_ConcurrentFirstSaveKeepsWinnerFields(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id := uuid.Must(uuid.NewV4())
	source := "apiVersion: v1\nkind: Namespace\n"

	fired := false
	if err := db.Callback().Create().Before("gorm:create").Register("test:concurrent_insert", func(tx *gorm.DB) {
		if fired {
			return
		}
		fired = true
		tx.Session(&gorm.Session{NewDB: true}).Exec(
			"INSERT INTO meshery_patterns (id, name, pattern_file, visibility, type, source_content, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			id, "winner", persisterTestDesign, Public, string(K8sManifest), []byte(source), persisterTestCreatedAt)
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	if _, err := mpp.SaveMesheryPattern(postStylePattern(id)); err != nil {
		t.Fatalf("save: %v", err)
	}
	if !fired {
		t.Fatalf("insert path was not taken")
	}

	got := loadPattern(t, db, id)
	if got.Name != "persister-test" {
		t.Errorf("name = %q, want the losing save's name applied", got.Name)
	}
	if got.Visibility != Public {
		t.Errorf("visibility = %q, want %q", got.Visibility, Public)
	}
	if string(got.SourceContent) != source || got.Type.String != string(K8sManifest) {
		t.Errorf("type/source = %q/%q, want the winner's values", got.Type.String, got.SourceContent)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(persisterTestCreatedAt) {
		t.Errorf("created_at = %v, want %v", got.CreatedAt, persisterTestCreatedAt)
	}
	if n := countPatterns(t, db); n != 1 {
		t.Errorf("rows = %d, want 1", n)
	}
}

func TestSaveMesheryPattern_DeletedDuringSaveIsError(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, _ := seedImportedPattern(t, db)

	if err := db.Callback().Update().Before("gorm:update").Register("test:concurrent_delete", func(tx *gorm.DB) {
		tx.Session(&gorm.Session{NewDB: true}).Exec("DELETE FROM meshery_patterns WHERE id = ?", id)
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	_, err := mpp.SaveMesheryPattern(postStylePattern(id))
	if code := mkerrors.GetCode(err); code != ErrDBPutCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrDBPutCode, err)
	}
	if n := countPatterns(t, db); n != 0 {
		t.Errorf("rows = %d, want the deleted design to stay deleted", n)
	}
}

// The save response must show the stored row, not the sparse request.
func TestSaveMesheryPattern_ResponseShowsStoredRow(t *testing.T) {
	db, mpp := newPatternPersisterTestDB(t)
	id, seed := seedImportedPattern(t, db)

	resp, err := mpp.SaveMesheryPattern(postStylePattern(id))
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	var saved []MesheryPattern
	if err := json.Unmarshal(resp, &saved); err != nil || len(saved) != 1 {
		t.Fatalf("decode response %s: %v", resp, err)
	}
	got := saved[0]
	if string(got.SourceContent) != string(seed.SourceContent) || got.Type != seed.Type {
		t.Errorf("response type/source = %+v/%q, want the stored values", got.Type, got.SourceContent)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(persisterTestCreatedAt) {
		t.Errorf("response createdAt = %v, want %v", got.CreatedAt, persisterTestCreatedAt)
	}
	if got.Visibility != seed.Visibility || got.ViewCount != seed.ViewCount {
		t.Errorf("response visibility/viewCount = %q/%d, want %q/%d", got.Visibility, got.ViewCount, seed.Visibility, seed.ViewCount)
	}
	if got.UserID == nil {
		t.Errorf("response is missing the stamped owner")
	}
}
