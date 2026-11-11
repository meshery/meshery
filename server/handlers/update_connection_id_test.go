package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// putConnection sends PUT /api/integrations/connections/{connectionId} with a
// body that carries no id, as the UI's partial updates do.
func putConnection(h *Handler, provider models.Provider, id string) *httptest.ResponseRecorder {
	return putConnectionBody(h, provider, id, `{"name":"renamed","kind":"prometheus"}`)
}

// putConnectionBody sends PUT /api/integrations/connections/{connectionId} with
// an explicit JSON body.
func putConnectionBody(h *Handler, provider models.Provider, id, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+id, strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"connectionId": id})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

func countConnections(t *testing.T, provider *models.DefaultLocalProvider) int64 {
	t.Helper()
	var n int64
	require.NoError(t, provider.ConnectionPersister.DB.Model(&connections.Connection{}).Count(&n).Error)
	return n
}

// TestUpdateConnectionByIDMalformedIDIsBadRequest checks that an ID that is not
// a UUID is a client error and does not reach the database.
func TestUpdateConnectionByIDMalformedIDIsBadRequest(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := putConnection(h, provider, "not-a-uuid")

	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
	assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDNilIDIsBadRequest checks that the nil ID is a client
// error: issue #22999 requires a malformed or nil ID to answer 400, and
// nothing is written.
func TestUpdateConnectionByIDNilIDIsBadRequest(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := putConnection(h, provider, uuid.Nil.String())

	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
	assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDUnknownIDIsNotFound checks that updating an ID that
// was never issued answers 404 instead of creating a connection under it.
func TestUpdateConnectionByIDUnknownIDIsNotFound(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := putConnection(h, provider, uuid.Must(uuid.NewV4()).String())

	assert.Equal(t, http.StatusNotFound, rec.Code, rec.Body.String())
	assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
}

// failingLookupProvider wraps a real provider and fails the existence lookup
// with a non-404 error, as a database outage would.
type failingLookupProvider struct {
	models.Provider
}

func (p failingLookupProvider) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	return nil, http.StatusInternalServerError, errors.New("database unreachable")
}

// TestUpdateConnectionByIDLookupErrorIsServerError checks that an existence
// lookup failing for a reason other than "not found" stops the update: the
// provider backfills omitted fields from the lookup, so proceeding would let
// a partial record clobber the stored one. The 500 names a failed lookup,
// not a missing query parameter, and the provider's own error stays out of
// the response.
func TestUpdateConnectionByIDLookupErrorIsServerError(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := putConnection(h, failingLookupProvider{Provider: provider}, uuid.Must(uuid.NewV4()).String())

	assert.Equal(t, http.StatusInternalServerError, rec.Code, rec.Body.String())
	assert.Contains(t, rec.Body.String(), "Unable to fetch data from the Provider")
	assert.NotContains(t, rec.Body.String(), "query parameter")
	assert.NotContains(t, rec.Body.String(), "database unreachable")
	assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDExistingIsUpdated guards the fix's blast radius: the
// lookup added before the save must not stop a stored connection from updating.
func TestUpdateConnectionByIDExistingIsUpdated(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "original", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)

	rec := putConnection(h, provider, saved.ID.String())

	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	got, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.Name)
	assert.Equal(t, int64(1), countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDBodyIDDoesNotRedirect checks that the URL id is the
// authority for the update: a body carrying a different non-nil id must not
// redirect the update onto that other connection.
func TestUpdateConnectionByIDBodyIDDoesNotRedirect(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	target, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "target", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)
	other, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "other", Kind: "kubernetes", ConnectionType: "platform",
	})
	require.NoError(t, err)

	body := `{"id":"` + other.ID.String() + `","name":"renamed","kind":"prometheus"}`
	rec := putConnectionBody(h, provider, target.ID.String(), body)

	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	gotTarget, err := provider.ConnectionPersister.GetConnection(target.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "renamed", gotTarget.Name)
	gotOther, err := provider.ConnectionPersister.GetConnection(other.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "other", gotOther.Name, "a body id must not redirect the update")
	assert.Equal(t, int64(2), countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDStatusOnlyPreservesFields checks that a status-only
// PUT, as the UI's connect action sends, updates the status without clobbering
// the fields the payload omits.
func TestUpdateConnectionByIDStatusOnlyPreservesFields(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "original", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)

	rec := putConnectionBody(h, provider, saved.ID.String(), `{"status":"connected"}`)

	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	got, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "original", got.Name)
	assert.Equal(t, "prometheus", got.Kind)
	assert.Equal(t, connections.CONNECTED, got.Status)
	assert.Equal(t, int64(1), countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDPreservesCreatedAtAndOwner checks that updating an
// existing connection preserves its original created_at timestamp and owner
// (user_id) both in the database row and in the HTTP response body.
func TestUpdateConnectionByIDPreservesCreatedAtAndOwner(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	ownerID := uuid.Must(uuid.NewV4())
	created := time.Date(2026, 11, 1, 9, 30, 0, 0, time.UTC)

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "original",
		Kind:           "prometheus",
		ConnectionType: "observability",
		SubType:        "monitoring",
		Status:         connections.DISCOVERED,
		Owner:          &ownerID,
		CreatedAt:      created,
		UpdatedAt:      created,
	})
	require.NoError(t, err)

	rec := putConnectionBody(h, provider, saved.ID.String(), `{"name":"renamed","status":"connected"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	got, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.Name)
	assert.Equal(t, connections.CONNECTED, got.Status)
	assert.True(t, got.CreatedAt.Equal(created), "stored created_at = %v, want %v to survive the update", got.CreatedAt, created)
	if assert.NotNil(t, got.Owner, "stored owner (user_id) must survive the update") {
		assert.Equal(t, ownerID, *got.Owner)
	}

	// Direct provider call (used by FSM transitions and action handlers that
	// serialize the returned *connections.Connection into the HTTP response)
	// must also return the preserved created_at and owner.
	updated, err := provider.UpdateConnectionById("", &connections.ConnectionPayload{
		Status: connections.MAINTENANCE,
	}, saved.ID.String())
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, connections.MAINTENANCE, updated.Status)
	assert.True(t, updated.CreatedAt.Equal(created), "returned created_at = %v, want %v", updated.CreatedAt, created)
	if assert.NotNil(t, updated.Owner, "returned owner (user_id) must survive the update") {
		assert.Equal(t, ownerID, *updated.Owner)
	}

	// Direct UpdateConnection call omitting CreatedAt must not wipe created_at
	// in the database and must populate the returned struct from the stored row.
	directUpdate := &connections.Connection{
		ID:             saved.ID,
		Name:           "direct-update",
		Kind:           saved.Kind,
		ConnectionType: saved.ConnectionType,
		SubType:        saved.SubType,
		Status:         connections.CONNECTED,
		Owner:          &ownerID,
		UpdatedAt:      time.Now(),
	}
	returnedDirect, err := provider.UpdateConnection(nil, directUpdate)
	require.NoError(t, err)
	require.NotNil(t, returnedDirect)
	assert.Equal(t, "direct-update", returnedDirect.Name)
	assert.True(t, returnedDirect.CreatedAt.Equal(created), "direct UpdateConnection returned created_at = %v, want %v", returnedDirect.CreatedAt, created)
}

// TestUpdateConnectionByIDNonKubernetesStatusRecordsNoErrorEvent: a status
// change for a non-Kubernetes connection must not run the Kubernetes machine
// path, which fails the K8s context lookup and persisted an error event while
// the client still got 200.
func TestUpdateConnectionByIDNonKubernetesStatusRecordsNoErrorEvent(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "original", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)

	rec := putConnectionBody(h, provider, saved.ID.String(), `{"status":"ignored"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	got, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	require.NoError(t, err)
	assert.Equal(t, connections.IGNORED, got.Status)

	var errorEvents int64
	require.NoError(t, provider.EventsPersister.DB.Model(&events.Event{}).Where("severity = ?", string(events.Error)).Count(&errorEvents).Error)
	assert.Zero(t, errorEvents, "a non-Kubernetes status change must not persist an error event")
	assert.Empty(t, h.ConnectionToStateMachineInstanceTracker.ConnectToInstanceMap, "no Kubernetes machine may be created for it")
}

// TestNotifySmOfConnectionStatusChangeByKind pins the routing: a non-Kubernetes
// kind returns without touching the Kubernetes machinery, while a Kubernetes
// (or unknown) kind still resolves its K8s context first.
func TestNotifySmOfConnectionStatusChangeByKind(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	userID := uuid.Must(uuid.NewV4())
	newPayload := func() *connections.ConnectionPayload {
		return &connections.ConnectionPayload{ID: uuid.Must(uuid.NewV4()), Status: connections.CONNECTED}
	}

	for _, kind := range []string{"prometheus", "grafana", "anthropic-claude"} {
		t.Run(kind+" skips the Kubernetes path", func(t *testing.T) {
			event, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, provider, "", newPayload(), kind)
			require.NoError(t, err)
			assert.NotEqual(t, events.Error, event.Severity)
			assert.Empty(t, h.ConnectionToStateMachineInstanceTracker.ConnectToInstanceMap)
		})
	}

	for _, kind := range []string{"kubernetes", "Kubernetes", ""} {
		t.Run("kind '"+kind+"' keeps the Kubernetes path", func(t *testing.T) {
			event, err := h.NotifySmOfConnectionStatusChange(context.Background(), userID, provider, "", newPayload(), kind)
			require.Error(t, err, "the K8s context lookup must still run for this kind")
			assert.Equal(t, events.Error, event.Severity)
		})
	}
}

// blankKindUpdateProvider mimics a provider whose update response omits Kind,
// as a remote provider's PUT response may.
type blankKindUpdateProvider struct {
	models.Provider
}

func (p blankKindUpdateProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	updated, err := p.Provider.UpdateConnectionById(token, conn, connID)
	if updated != nil {
		updated.Kind = ""
	}
	return updated, err
}

// TestUpdateConnectionByIDFallsBackToPersistedKind checks that when the update
// response carries no kind, the handler routes on the persisted connection's
// kind instead of falling through to the Kubernetes path.
func TestUpdateConnectionByIDFallsBackToPersistedKind(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "original", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)

	rec := putConnectionBody(h, blankKindUpdateProvider{Provider: provider}, saved.ID.String(), `{"status":"ignored"}`)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var errorEvents int64
	require.NoError(t, provider.EventsPersister.DB.Model(&events.Event{}).Where("severity = ?", string(events.Error)).Count(&errorEvents).Error)
	assert.Zero(t, errorEvents, "an empty kind in the update response must fall back to the persisted kind")
}
