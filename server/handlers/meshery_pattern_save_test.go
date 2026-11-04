package handlers

import (
	"database/sql"
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
)

// TestHandlePatternPOST_ResaveKeepsStoredFields saves an existing imported
// design through POST /api/pattern with the body the design editor sends (id,
// name, designFile) and checks that the columns the body does not carry
// survive in the database.
func TestHandlePatternPOST_ResaveKeepsStoredFields(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(&models.MesheryPattern{}, &events.Event{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

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

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t), SystemID: &systemID}
	provider := &models.DefaultLocalProvider{
		MesheryPatternPersister: &models.MesheryPatternPersister{DB: &db},
		EventsPersister:         &models.EventsPersister{DB: &db},
	}

	body := `{"id":"` + id.String() + `","name":"demo","designFile":{"id":"` + id.String() +
		`","name":"demo","schemaVersion":"designs.meshery.io/v1beta1","version":"0.0.1","components":[],"relationships":[]}}`
	req := httptest.NewRequest(http.MethodPost, "/api/pattern", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.handlePatternPOST(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var got models.MesheryPattern
	if err := db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load: %v", err)
	}
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
