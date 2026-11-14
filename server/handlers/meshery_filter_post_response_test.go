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
	// The remote provider reads its token from this cookie; the built-in one ignores it.
	req.AddCookie(&http.Cookie{Name: models.TokenCookieName, Value: "test-token"})
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

// TestFilterPOSTResponseCarriesURLImportedFilter pins the body of an import
// from a URL: it must carry the imported filter as stored.
func TestFilterPOSTResponseCarriesURLImportedFilter(t *testing.T) {
	h, provider, db := newFilterPOSTFixture(t)
	provider.Log = h.log
	source := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("name: from-url\n"))
	}))
	t.Cleanup(source.Close)

	rec := postFilter(t, h, provider, models.MesheryFilterRequestBody{Save: true, URL: source.URL + "/filter.yaml"})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var imported []models.MesheryFilter
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &imported))
	require.Len(t, imported, 1, "the response must carry the imported filter, got %s", rec.Body.String())
	assert.Equal(t, "from-url", imported[0].Name)
	require.NotNil(t, imported[0].ID)

	var stored models.MesheryFilter
	require.NoError(t, db.Where("id = ?", *imported[0].ID).First(&stored).Error,
		"the returned id must be the stored row's id")
}

// TestFilterPOSTResponseCarriesRemoteProviderFilters pins the save path with
// the remote provider: the body carries the filters Cloud answered with,
// re-encoded through MesheryFilter, so fields outside that model are dropped.
func TestFilterPOSTResponseCarriesRemoteProviderFilters(t *testing.T) {
	h, _, _ := newFilterPOSTFixture(t)
	savedID := uuid.Must(uuid.NewV4()).String()
	var gotMethod, gotPath, gotAuth string
	cloud := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath, gotAuth = r.Method, r.URL.Path, r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"` + savedID + `","name":"rate-limit","visibility":"private","internalNote":"cloud-only"}]`))
	}))
	t.Cleanup(cloud.Close)

	remote := &models.RemoteProvider{
		ProviderProperties: models.ProviderProperties{
			ProviderType: models.RemoteProviderType,
			Capabilities: models.Capabilities{{Feature: models.PersistMesheryFilters, Endpoint: "/api/content/filters"}},
		},
		RemoteProviderURL: cloud.URL,
		Log:               h.log,
	}

	rec := postFilter(t, h, remote, models.MesheryFilterRequestBody{
		Save:       true,
		FilterData: &models.MesheryFilterPayload{Name: "rate-limit", FilterFile: []byte("wasm-bytes")},
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/content/filters", gotPath)
	assert.Equal(t, "bearer test-token", gotAuth)

	var saved []models.MesheryFilter
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &saved))
	require.Len(t, saved, 1, "the response must carry the filter Cloud saved, got %s", rec.Body.String())
	require.NotNil(t, saved[0].ID)
	assert.Equal(t, savedID, saved[0].ID.String())
	assert.Equal(t, "rate-limit", saved[0].Name)
	assert.NotContains(t, rec.Body.String(), "internalNote", "fields outside MesheryFilter must not be echoed")
}
