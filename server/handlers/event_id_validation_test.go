package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newEventIDFixture wires the handler over DefaultLocalProvider and an in-memory
// database, so the test sees what gorm does with the parsed event ID.
func newEventIDFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(events.Event{}))

	h := &Handler{log: newTestLogger(t)}
	provider := &models.DefaultLocalProvider{EventsPersister: &models.EventsPersister{DB: &db}}
	return h, provider
}

func deleteEvent(h *Handler, provider models.Provider, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, "/api/system/events/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"id": id})
	rec := httptest.NewRecorder()
	h.DeleteEvent(rec, req, nil, &models.User{}, provider)
	return rec
}

func updateEventStatus(h *Handler, provider models.Provider, id, status string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/"+id, strings.NewReader(`{"status":"`+status+`"}`))
	req = mux.SetURLVars(req, map[string]string{"id": id})
	rec := httptest.NewRecorder()
	h.UpdateEventStatus(rec, req, nil, &models.User{}, provider)
	return rec
}

// TestEventHandlersRejectInvalidID verifies that an ID that is not a UUID, or is
// the nil UUID, is a client error and not a 500 from the database.
func TestEventHandlersRejectInvalidID(t *testing.T) {
	for _, id := range []string{"not-a-uuid", uuid.Nil.String()} {
		t.Run("delete/"+id, func(t *testing.T) {
			h, provider := newEventIDFixture(t)
			rec := deleteEvent(h, provider, id)
			assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
		})
		t.Run("update_status/"+id, func(t *testing.T) {
			h, provider := newEventIDFixture(t)
			rec := updateEventStatus(h, provider, id, "read")
			assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
		})
	}
}

// TestEventHandlersAcceptStoredEvent verifies that a valid ID still updates and
// deletes the stored event.
func TestEventHandlersAcceptStoredEvent(t *testing.T) {
	h, provider := newEventIDFixture(t)

	event := events.Event{
		ID:          uuid.Must(uuid.NewV4()),
		Category:    "test",
		Action:      "test",
		Description: "stored event",
		Status:      events.EventStatus("unread"),
	}
	require.NoError(t, provider.EventsPersister.DB.Create(&event).Error)
	id := event.ID.String()

	rec := updateEventStatus(h, provider, id, "read")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var stored events.Event
	require.NoError(t, provider.EventsPersister.DB.First(&stored, "id = ?", event.ID).Error)
	assert.Equal(t, events.EventStatus("read"), stored.Status)

	rec = deleteEvent(h, provider, id)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var count int64
	require.NoError(t, provider.EventsPersister.DB.Model(&events.Event{}).Where("id = ?", event.ID).Count(&count).Error)
	assert.Zero(t, count, "the event must be deleted")
}
