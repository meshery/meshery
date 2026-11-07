package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
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

// Deleting a Kubernetes connection while it is NOTFOUND resolves its open
// unreachable condition instead of leaving it open forever; deleting one that
// is not NOTFOUND resolves nothing.
func TestDeleteConnectionResolvesTheUnreachableCondition(t *testing.T) {
	for _, c := range []struct {
		status      connections.ConnectionStatus
		wantResolve bool
	}{
		{connections.ConnectionStatus(machines.NOTFOUND), true},
		{connections.CONNECTED, false},
	} {
		h, provider := newDeleteConnectionFixture(t)
		user := &models.User{ID: uuid.Must(uuid.NewV4())}
		saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
			Name: "test-connection", Kind: "kubernetes", ConnectionType: "platform", Status: c.status,
		})
		if err != nil {
			t.Fatalf("save connection: %v", err)
		}
		key := kubernetes.UnreachableConditionKey(saved.ID)
		if _, err := provider.EventsPersister.PersistEvent(*events.NewEvent().FromOwner(user.ID).WithCategory("connection").
			WithAction("update").WithSeverity(events.Error).WithCorrelationKey(key).Build(), ""); err != nil {
			t.Fatalf("open the condition: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/"+saved.ID.String(), nil)
		req = mux.SetURLVars(req, map[string]string{"connectionId": saved.ID.String()})
		rec := httptest.NewRecorder()
		h.DeleteConnection(rec, req, nil, user, provider)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d. body: %s", c.status, rec.Code, rec.Body.String())
		}

		var open int64
		if err := provider.EventsPersister.DB.Model(&events.Event{}).
			Where("correlation_key = ? AND resolved_at IS NULL", key).Count(&open).Error; err != nil {
			t.Fatalf("count open conditions: %v", err)
		}
		if got := open == 0; got != c.wantResolve {
			t.Errorf("%s: %d open conditions after the delete, want resolved=%v", c.status, open, c.wantResolve)
		}
	}
}
