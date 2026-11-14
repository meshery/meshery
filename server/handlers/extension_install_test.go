package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

// newTestLocalProvider returns an initialized DefaultLocalProvider with
// LocalProviderType set, representing a real local provider instance for tests.
func newTestLocalProvider() *models.DefaultLocalProvider {
	lp := &models.DefaultLocalProvider{}
	lp.Initialize()
	return lp
}

// TestInstallExtensionHandler_NilProvider_ReturnsInternalServerError verifies that
// InstallExtensionHandler safely guards against a nil provider without triggering
// a SIGSEGV runtime panic, returning HTTP 500.
func TestInstallExtensionHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(`{"extType":"navigator"}`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, nil)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d for nil provider, got %d", http.StatusInternalServerError, rec.Code)
	}
}

// TestRemoveExtensionHandler_NilProvider_ReturnsInternalServerError verifies that
// RemoveExtensionHandler safely guards against a nil provider without triggering
// a SIGSEGV runtime panic, returning HTTP 500.
func TestRemoveExtensionHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(`{"extType":"navigator","title":"MeshMap"}`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.RemoveExtensionHandler(rec, req, nil, &models.User{ID: userID}, nil)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d for nil provider, got %d", http.StatusInternalServerError, rec.Code)
	}
}

// TestInstallExtensionHandler_RemoteProviderReturnsNotImplemented verifies that
// attempting to install an extension using a remote provider returns HTTP 501 Not Implemented.
func TestInstallExtensionHandler_RemoteProviderReturnsNotImplemented(t *testing.T) {
	remote := &models.RemoteProvider{
		ProviderProperties: models.ProviderProperties{
			ProviderType: models.RemoteProviderType,
		},
	}
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(`{"extType":"navigator"}`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, remote)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected status %d, got %d", http.StatusNotImplemented, rec.Code)
	}
}

// TestRemoveExtensionHandler_RemoteProviderReturnsNotImplemented verifies that
// attempting to remove an extension using a remote provider returns HTTP 501 Not Implemented.
func TestRemoveExtensionHandler_RemoteProviderReturnsNotImplemented(t *testing.T) {
	remote := &models.RemoteProvider{
		ProviderProperties: models.ProviderProperties{
			ProviderType: models.RemoteProviderType,
		},
	}
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(`{"extType":"navigator","title":"MeshMap"}`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.RemoveExtensionHandler(rec, req, nil, &models.User{ID: userID}, remote)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected status %d, got %d", http.StatusNotImplemented, rec.Code)
	}
}

// TestInstallExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge verifies that
// payloads exceeding the 4MB payload limit are rejected with HTTP 413.
func TestInstallExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	oversizedBody := bytes.Repeat([]byte("a"), maxExtensionPayloadSize+10)
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d for oversized payload, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}

// TestRemoveExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge verifies that
// removal payloads exceeding the 4MB payload limit are rejected with HTTP 413.
func TestRemoveExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	oversizedBody := bytes.Repeat([]byte("a"), maxExtensionPayloadSize+10)
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.RemoveExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status %d for oversized payload, got %d", http.StatusRequestEntityTooLarge, rec.Code)
	}
}

// TestInstallExtensionHandler_EmptyBody_ReturnsBadRequest verifies that
// an empty request body returns HTTP 400 Bad Request.
func TestInstallExtensionHandler_EmptyBody_ReturnsBadRequest(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(""))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for empty body, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestRemoveExtensionHandler_EmptyBody_ReturnsBadRequest verifies that
// an empty request body returns HTTP 400 Bad Request.
func TestRemoveExtensionHandler_EmptyBody_ReturnsBadRequest(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(""))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.RemoveExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for empty body, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestInstallExtensionHandler_InvalidJSON_ReturnsBadRequest verifies that
// malformed JSON request bodies return HTTP 400 Bad Request.
func TestInstallExtensionHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(`{"extType":`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for invalid json, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestRemoveExtensionHandler_InvalidJSON_ReturnsBadRequest verifies that
// malformed JSON request bodies return HTTP 400 Bad Request.
func TestRemoveExtensionHandler_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(`{"extType":`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.RemoveExtensionHandler(rec, req, nil, &models.User{ID: userID}, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d for invalid json, got %d", http.StatusBadRequest, rec.Code)
	}
}

// TestPersistLocalProviderCapabilitiesForUser_NilSafety verifies that
// persistLocalProviderCapabilitiesForUser handles nil user and nil provider
// without panicking or failing.
func TestPersistLocalProviderCapabilitiesForUser_NilSafety(t *testing.T) {
	if err := persistLocalProviderCapabilitiesForUser(nil, nil); err != nil {
		t.Fatalf("expected nil error for nil inputs, got %v", err)
	}

	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}
	user := &models.User{ID: userID}
	if err := persistLocalProviderCapabilitiesForUser(user, nil); err != nil {
		t.Fatalf("expected nil error for nil provider, got %v", err)
	}

	nilUser := &models.User{ID: uuid.Nil}
	if err := persistLocalProviderCapabilitiesForUser(nilUser, nil); err != nil {
		t.Fatalf("expected nil error for zero-UUID user, got %v", err)
	}
}
