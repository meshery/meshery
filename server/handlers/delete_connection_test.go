package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
)

// newDeleteConnectionFixture wires the handler the way the router does, over an
// in-memory database. It deliberately goes through DefaultLocalProvider rather
// than a stub: the behaviour under test is what the handler makes of a
// provider's error, and the local provider is the one whose "not found" the
// remote provider was taught to match.
func newDeleteConnectionFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, events.Event{}); err != nil {
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
		EventsPersister:     &models.EventsPersister{DB: &db},
	}
	return h, provider
}

func deleteConnection(t *testing.T, h *Handler, provider models.Provider, connectionID uuid.UUID) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+connectionID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.DeleteConnection(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// TestDeleteConnectionNotFound pins the status `mesheryctl connection delete`
// keys off. A syntactically valid id that was never issued must come back 404,
// which the CLI turns into a "No connection with ID ..." warning and exit 0.
// It used to come back 500 "Failed to Save: .connection"
// (meshery-server-1051), because every provider error - not-found included -
// was wrapped in ErrFailToSave.
func TestDeleteConnectionNotFound(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := deleteConnection(t, h, provider, uuid.Must(uuid.NewV4()))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

// TestDeleteConnectionExisting is the regression guard for the fix's blast
// radius: the not-found path must not swallow real deletes. An earlier attempt
// short-circuited on the handler's own database before calling the provider,
// which 404s every connection under a remote provider - those rows live at the
// provider, not in Meshery Server's database.
func TestDeleteConnectionExisting(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "test-connection",
		Kind:           "kubernetes",
		ConnectionType: "platform",
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}

	rec := deleteConnection(t, h, provider, saved.ID)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	if _, err := provider.ConnectionPersister.GetConnection(saved.ID, ""); err == nil {
		t.Error("connection is still readable after a successful delete")
	}
}
