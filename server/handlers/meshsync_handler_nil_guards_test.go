package handlers

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockNilPersisterProvider struct {
	*models.DefaultLocalProvider
}

func (m *mockNilPersisterProvider) GetGenericPersister() *database.Handler {
	return nil
}

func newMeshSyncTestProvider(t *testing.T) *models.DefaultLocalProvider {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "meshsync_test.db")})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.KubernetesResource{}, &model.KubernetesResourceObjectMeta{},
		&model.KubernetesKeyValue{}, &model.KubernetesResourceSpec{}, &model.KubernetesResourceStatus{}))
	return &models.DefaultLocalProvider{GenericPersister: &db}
}

func TestMeshSyncHandlers_NilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	t.Run("GetMeshSyncResources with untyped nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResources(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResources with typed nil provider returns 500", func(t *testing.T) {
		var typedNil *models.DefaultLocalProvider
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResources(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourceByID with untyped nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourceByID(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourceByID with typed nil provider returns 500", func(t *testing.T) {
		var typedNil *models.DefaultLocalProvider
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourceByID(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourcesSummary with untyped nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/summary?clusterId=c-1", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourcesSummary(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourcesSummary with typed nil provider returns 500", func(t *testing.T) {
		var typedNil *models.DefaultLocalProvider
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/summary?clusterId=c-1", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourcesSummary(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("DeleteMeshSyncResource with untyped nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteMeshSyncResource(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("DeleteMeshSyncResource with typed nil provider returns 500", func(t *testing.T) {
		var typedNil *models.DefaultLocalProvider
		req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteMeshSyncResource(rec, req, nil, nil, typedNil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})
}

func TestMeshSyncHandlers_NilPersister(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	prov := &mockNilPersisterProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}

	t.Run("GetMeshSyncResources with nil persister returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResources(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourceByID with nil persister returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourceByID(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("GetMeshSyncResourcesSummary with nil persister returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/summary?clusterId=c-1", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourcesSummary(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("DeleteMeshSyncResource with nil persister returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/00000000-0000-0000-0000-000000000001", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "00000000-0000-0000-0000-000000000001"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteMeshSyncResource(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})
}

func TestMeshSyncHandlers_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    nil,
	}

	t.Run("GetMeshSyncResourcesSummary with nil logger does not panic on missing clusterIds", func(t *testing.T) {
		prov := newMeshSyncTestProvider(t)
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/summary", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourcesSummary(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetMeshSyncResourceByID with nil logger does not panic on invalid ID", func(t *testing.T) {
		prov := newMeshSyncTestProvider(t)
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources/invalid-id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResourceByID(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DeleteMeshSyncResource with nil logger does not panic on invalid ID", func(t *testing.T) {
		prov := newMeshSyncTestProvider(t)
		req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/invalid-id", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid-id"})
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.DeleteMeshSyncResource(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetMeshSyncResources with nil logger does not panic on malformed clusterIds", func(t *testing.T) {
		prov := newMeshSyncTestProvider(t)
		req := httptest.NewRequest(http.MethodGet, "/api/system/meshsync/resources?clusterIds=invalid-json", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetMeshSyncResources(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
