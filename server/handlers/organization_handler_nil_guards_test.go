package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrganizations_NilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	t.Run("untyped nil provider returns 500", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetOrganizations(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})

	t.Run("typed nil provider returns 500", func(t *testing.T) {
		var typedNilProv *models.DefaultLocalProvider
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetOrganizations(rec, req, nil, nil, typedNilProv)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), ErrNilProviderCode)
	})
}

func TestGetOrganizations_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    nil,
	}

	t.Run("nil logger does not panic on nil provider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetOrganizations(rec, req, nil, nil, nil)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("nil logger does not panic on missing token", func(t *testing.T) {
		prov := &models.DefaultLocalProvider{}
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetOrganizations(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("nil logger does not panic on provider error", func(t *testing.T) {
		prov := newKeysOrgsFailingProvider(errors.New("provider failure"))
		req := httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil)
		ctx := context.WithValue(req.Context(), models.TokenCtxKey, "test-token")
		req = req.WithContext(ctx)
		rec := httptest.NewRecorder()
		require.NotPanics(t, func() {
			h.GetOrganizations(rec, req, nil, nil, prov)
		})
		assert.Equal(t, http.StatusBadGateway, rec.Code)
	})
}
