package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogoutHandler_NilUser(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.LogoutHandler(rec, req, nil, local)
	})

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/provider", resp.Header.Get("Location"))

	// Ensure cookies are cleared with negative MaxAge
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies)
	for _, c := range cookies {
		assert.Equal(t, -1, c.MaxAge)
		assert.Empty(t, c.Value)
	}
}

func TestLogoutHandler_ValidUser(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	userID := uuid.Must(uuid.NewV4())
	user := &models.User{
		ID:        userID,
		FirstName: "Test",
		LastName:  "User",
	}

	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.LogoutHandler(rec, req, user, local)
	})

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/provider", resp.Header.Get("Location"))

	cookies := resp.Cookies()
	require.NotEmpty(t, cookies)
	for _, c := range cookies {
		assert.Equal(t, -1, c.MaxAge)
		assert.Empty(t, c.Value)
	}
}

func TestLogoutHandler_NilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.LogoutHandler(rec, req, nil, nil)
	})

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	require.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/provider", resp.Header.Get("Location"))
}

func TestLogoutHandler_MethodNotAllowed(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	h := newTestHandler(t, map[string]models.Provider{local.Name(): local}, "")

	req := httptest.NewRequest(http.MethodPost, "/user/logout", nil)
	rec := httptest.NewRecorder()

	h.LogoutHandler(rec, req, nil, local)

	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
