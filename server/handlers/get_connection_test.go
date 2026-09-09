package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
)

func newGetConnectionFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
	}
	return h, provider
}

func getConnection(t *testing.T, h *Handler, provider models.Provider, connectionID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections/"+connectionID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.GetConnectionByID(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

func TestGetConnectionNotFound(t *testing.T) {
	h, provider := newGetConnectionFixture(t)

	rec := getConnection(t, h, provider, uuid.Must(uuid.NewV4()))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}

	if code, ok := response["code"].(string); !ok || code != ErrGetConnectionCode {
		t.Fatalf("error code = %v, want %q", response["code"], ErrGetConnectionCode)
	}
}

func TestGetConnectionExisting(t *testing.T) {
	h, provider := newGetConnectionFixture(t)

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "test-connection",
		Kind:           "kubernetes",
		ConnectionType: "platform",
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}

	rec := getConnection(t, h, provider, saved.ID)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var fetched connections.Connection
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("unmarshal connection: %v", err)
	}

	if fetched.ID != saved.ID {
		t.Fatalf("connection ID = %s, want %s", fetched.ID, saved.ID)
	}
}
