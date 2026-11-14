package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	_events "github.com/meshery/meshkit/utils/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFilterPOSTFixture extends the delete fixture with what the POST path
// needs: a registry for the filter component lookup, an events table for the
// event the unsaved-conversion path records, and an events buffer.
func newFilterPOSTFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider, database.Handler) {
	t.Helper()

	h, provider, db := newDeleteFilterFixture(t)
	t.Cleanup(func() { _ = db.DBClose() })
	require.NoError(t, db.AutoMigrate(events.Event{}))
	provider.EventsPersister = &models.EventsPersister{DB: &db}

	rm, _ := newTestRegistryDB(t)
	h.registryManager = rm
	h.EventsBuffer = _events.NewEventStreamer()
	return h, provider, db
}

func postFilter(t *testing.T, h *Handler, provider models.Provider, body models.MesheryFilterRequestBody) *httptest.ResponseRecorder {
	t.Helper()

	payload, err := json.Marshal(body)
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/filter", bytes.NewReader(payload))
	rec := httptest.NewRecorder()

	h.FilterFileRequestHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

// TestFilterPOSTResponseCarriesSavedFilter pins the body of a save: it must
// carry the filter the provider stored, including the id the server assigned,
// so a caller can address the filter it just created.
func TestFilterPOSTResponseCarriesSavedFilter(t *testing.T) {
	h, provider, db := newFilterPOSTFixture(t)

	rec := postFilter(t, h, provider, models.MesheryFilterRequestBody{
		Save: true,
		FilterData: &models.MesheryFilterPayload{
			Name:       "rate-limit",
			FilterFile: []byte("wasm-bytes"),
		},
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var saved []models.MesheryFilter
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &saved))
	require.Len(t, saved, 1, "the response must carry the saved filter, got %s", rec.Body.String())
	assert.Equal(t, "rate-limit", saved[0].Name)
	require.NotNil(t, saved[0].ID, "the response must carry the id the server assigned")

	var stored models.MesheryFilter
	require.NoError(t, db.Where("id = ?", *saved[0].ID).First(&stored).Error,
		"the returned id must be the stored row's id")
	assert.Equal(t, "rate-limit", stored.Name)
}

// TestFilterPOSTResponseCarriesUnsavedFilter pins the body when save is false:
// the handler only converts the filter, so the response is its only output.
func TestFilterPOSTResponseCarriesUnsavedFilter(t *testing.T) {
	h, provider, db := newFilterPOSTFixture(t)

	rec := postFilter(t, h, provider, models.MesheryFilterRequestBody{
		FilterData: &models.MesheryFilterPayload{
			Name:       "header-rewrite",
			FilterFile: []byte("wasm-bytes"),
		},
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var converted []models.MesheryFilter
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &converted))
	require.Len(t, converted, 1, "the response must carry the converted filter, got %s", rec.Body.String())
	assert.Equal(t, "header-rewrite", converted[0].Name)
	assert.Equal(t, []byte("wasm-bytes"), converted[0].FilterFile)

	var count int64
	require.NoError(t, db.Model(&models.MesheryFilter{}).Count(&count).Error)
	assert.Zero(t, count, "a request without save must not persist the filter")
}

// TestFormatFilterOutputKeepsEmptyListForNull pins the wire shape when a
// provider answers null: the body stays an empty list, never null.
func TestFormatFilterOutputKeepsEmptyListForNull(t *testing.T) {
	h, _, _ := newFilterPOSTFixture(t)
	rec := httptest.NewRecorder()

	h.formatFilterOutput(rec, []byte("null"), "", &meshes.EventsResponse{}, events.NewEvent())

	assert.JSONEq(t, "[]", rec.Body.String())
}
