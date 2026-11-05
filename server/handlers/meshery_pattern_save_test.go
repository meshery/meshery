package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/catalog/v1alpha1"
	"github.com/meshery/meshkit/models/events"
	"gorm.io/gorm"
)

// storedDesignFile is a design file with a component, so a test can tell it
// apart from the empty design a request without designFile would produce.
const storedDesignFile = "name: demo\nschemaVersion: designs.meshery.io/v1beta1\nversion: 0.0.3\ncomponents:\n- id: 11111111-1111-1111-1111-111111111111\n  displayName: stored-component\n"

type patternPOSTFixture struct {
	h        *Handler
	provider *models.DefaultLocalProvider
	db       *database.Handler
}

func newPatternPOSTFixture(t *testing.T) patternPOSTFixture {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.MesheryPattern{}, &events.Event{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	return patternPOSTFixture{
		h: &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t), SystemID: &systemID},
		provider: &models.DefaultLocalProvider{
			MesheryPatternPersister: &models.MesheryPatternPersister{DB: &db},
			EventsPersister:         &models.EventsPersister{DB: &db},
		},
		db: &db,
	}
}

func (f patternPOSTFixture) post(t *testing.T, body string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(body))
	rec := httptest.NewRecorder()
	f.h.handlePatternPOST(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, f.provider)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func (f patternPOSTFixture) load(t *testing.T, id uuid.UUID) models.MesheryPattern {
	t.Helper()
	var got models.MesheryPattern
	if err := f.db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load %s: %v", id, err)
	}
	return got
}

func (f patternPOSTFixture) seed(t *testing.T, id uuid.UUID) {
	t.Helper()
	if err := f.db.Create(&models.MesheryPattern{ID: &id, Name: "demo", PatternFile: storedDesignFile, Visibility: models.Private}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
}

// TestHandlePatternPOST_ResaveKeepsStoredFields saves an existing imported
// design through POST /api/pattern with the body the design editor sends (id,
// name, designFile) and checks that the columns the body does not carry
// survive in the database.
func TestHandlePatternPOST_ResaveKeepsStoredFields(t *testing.T) {
	f := newPatternPOSTFixture(t)
	db := f.db

	id := uuid.Must(uuid.NewV4())
	created := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	source := []byte("apiVersion: v1\nkind: Namespace\nmetadata:\n  name: demo\n")
	if err := db.Create(&models.MesheryPattern{
		ID:            &id,
		Name:          "demo",
		PatternFile:   "name: demo\nschemaVersion: designs.meshery.io/v1beta1\nversion: 0.0.1\ncomponents: []\n",
		Visibility:    models.Public,
		CatalogData:   v1alpha1.CatalogData{PatternInfo: "stored info"},
		Type:          sql.NullString{String: string(models.K8sManifest), Valid: true},
		SourceContent: source,
		CreatedAt:     &created,
	}).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}

	f.post(t, `{"id":"`+id.String()+`","name":"demo","designFile":{"id":"`+id.String()+
		`","name":"demo","schemaVersion":"designs.meshery.io/v1beta1","version":"0.0.1","components":[],"relationships":[]}}`)

	got := f.load(t, id)
	if string(got.SourceContent) != string(source) {
		t.Errorf("source_content = %q, want %q", got.SourceContent, source)
	}
	if got.Type.String != string(models.K8sManifest) {
		t.Errorf("type = %+v, want %q", got.Type, models.K8sManifest)
	}
	if got.CreatedAt == nil || !got.CreatedAt.Equal(created) {
		t.Errorf("created_at = %v, want %v", got.CreatedAt, created)
	}
	if got.Visibility != models.Public {
		t.Errorf("visibility = %q, want %q", got.Visibility, models.Public)
	}
	if got.CatalogData.PatternInfo != "stored info" {
		t.Errorf("catalog_data.patternInfo = %q, want %q", got.CatalogData.PatternInfo, "stored info")
	}
}

func TestHandlePatternPOST_OmittedDesignFileKeepsStoredDesign(t *testing.T) {
	f := newPatternPOSTFixture(t)
	id := uuid.Must(uuid.NewV4())
	f.seed(t, id)

	f.post(t, `{"id":"`+id.String()+`","name":"renamed","visibility":"public"}`)

	got := f.load(t, id)
	if got.PatternFile != storedDesignFile {
		t.Errorf("pattern_file = %q, want the stored design %q", got.PatternFile, storedDesignFile)
	}
	if got.Name != "renamed" || got.Visibility != models.Public {
		t.Errorf("name/visibility = %q/%q, want the supplied renamed/public", got.Name, got.Visibility)
	}
}

func TestHandlePatternPOST_SuppliedDesignFileReplacesStoredDesign(t *testing.T) {
	f := newPatternPOSTFixture(t)
	id := uuid.Must(uuid.NewV4())
	f.seed(t, id)

	f.post(t, `{"id":"`+id.String()+`","name":"demo","designFile":{"id":"`+id.String()+
		`","name":"supplied-design","schemaVersion":"designs.meshery.io/v1beta1","version":"0.0.3","components":[],"relationships":[]}}`)

	got := f.load(t, id)
	if !strings.Contains(got.PatternFile, "supplied-design") || strings.Contains(got.PatternFile, "stored-component") {
		t.Errorf("pattern_file = %q, want the supplied design", got.PatternFile)
	}
	if !strings.Contains(got.PatternFile, "0.0.4") {
		t.Errorf("pattern_file = %q, want the version bumped to 0.0.4", got.PatternFile)
	}
}

// A request without designFile can lose an insert race to another save of the
// same ID; the fallback update must keep the winner's design file too.
func TestHandlePatternPOST_OmittedDesignFileKeepsWinnerDesignOnInsertConflict(t *testing.T) {
	f := newPatternPOSTFixture(t)
	id := uuid.Must(uuid.NewV4())

	fired := false
	if err := f.db.Callback().Create().Before("gorm:create").Register("test:concurrent_insert", func(tx *gorm.DB) {
		if fired || tx.Statement.Table != "meshery_patterns" {
			return
		}
		fired = true
		tx.Session(&gorm.Session{NewDB: true}).Exec(
			"INSERT INTO meshery_patterns (id, name, pattern_file, visibility) VALUES (?, ?, ?, ?)",
			id, "winner", storedDesignFile, models.Public)
	}); err != nil {
		t.Fatalf("register callback: %v", err)
	}

	f.post(t, `{"id":"`+id.String()+`","name":"loser"}`)

	if !fired {
		t.Fatalf("insert path was not taken")
	}
	got := f.load(t, id)
	if got.PatternFile != storedDesignFile {
		t.Errorf("pattern_file = %q, want the winner's design %q", got.PatternFile, storedDesignFile)
	}
	if got.Visibility != models.Public {
		t.Errorf("visibility = %q, want the winner's %q", got.Visibility, models.Public)
	}
}

// New designs saved without designFile still get a row and a design file, as
// before field presence was tracked.
func TestHandlePatternPOST_OmittedDesignFileStillInsertsNewDesign(t *testing.T) {
	for name, idField := range map[string]string{
		"unknown ID":   `"id":"` + uuid.Must(uuid.NewV4()).String() + `",`,
		"generated ID": ``,
	} {
		t.Run(name, func(t *testing.T) {
			f := newPatternPOSTFixture(t)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(`{`+idField+`"name":"new-design"}`))
			rec := httptest.NewRecorder()
			f.h.handlePatternPOST(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, f.provider)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}

			var saved []models.MesheryPattern
			if err := json.Unmarshal(rec.Body.Bytes(), &saved); err != nil || len(saved) != 1 || saved[0].ID == nil {
				t.Fatalf("decode response %s: %v", rec.Body.String(), err)
			}
			got := f.load(t, *saved[0].ID)
			if got.Name != "new-design" || got.Visibility != models.Private {
				t.Errorf("name/visibility = %q/%q, want new-design/private", got.Name, got.Visibility)
			}
			if got.PatternFile == "" {
				t.Errorf("pattern_file is empty, want the empty design stored as before")
			}
		})
	}
}

func TestDesignPostPayload_TracksDesignFilePresence(t *testing.T) {
	cases := map[string]struct {
		body string
		want bool
	}{
		"designFile":       {`{"designFile":{"name":"d"}}`, true},
		"patternFile":      {`{"patternFile":{"name":"d"}}`, true},
		"design_file":      {`{"design_file":{"name":"d"}}`, true},
		"pattern_file":     {`{"pattern_file":{"name":"d"}}`, true},
		"empty designFile": {`{"designFile":{}}`, true},
		"absent":           {`{"name":"d"}`, false},
		"null designFile":  {`{"designFile":null}`, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var p DesignPostPayload
			if err := json.Unmarshal([]byte(tc.body), &p); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if p.designFilePresent != tc.want {
				t.Errorf("designFilePresent = %v, want %v", p.designFilePresent, tc.want)
			}
		})
	}

	var reused DesignPostPayload
	_ = json.Unmarshal([]byte(`{"designFile":{"name":"d"}}`), &reused)
	if err := json.Unmarshal([]byte(`{"name":"d"}`), &reused); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if reused.designFilePresent {
		t.Errorf("designFilePresent stayed true on a reused payload without a design file")
	}
}
