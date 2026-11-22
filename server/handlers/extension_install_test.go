package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
// a SIGSEGV runtime panic, returning HTTP 500 with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, "provider not initialized", errBody.Error)
}

// TestRemoveExtensionHandler_NilProvider_ReturnsInternalServerError verifies that
// RemoveExtensionHandler safely guards against a nil provider without triggering
// a SIGSEGV runtime panic, returning HTTP 500 with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, "provider not initialized", errBody.Error)
}

// TestInstallExtensionHandler_RemoteProviderReturnsNotImplemented verifies that
// attempting to install an extension using a remote provider returns HTTP 501 Not Implemented
// with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "install extension is currently supported only for the local provider")
}

// TestRemoveExtensionHandler_RemoteProviderReturnsNotImplemented verifies that
// attempting to remove an extension using a remote provider returns HTTP 501 Not Implemented
// with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusNotImplemented, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "remove extension is currently supported only for the local provider")
}

// TestInstallExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge verifies that
// payloads exceeding the 4MB payload limit are rejected with HTTP 413 and a structured error envelope.
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

	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "request body exceeds maximum allowed size")
}

// TestRemoveExtensionHandler_PayloadTooLarge_ReturnsEntityTooLarge verifies that
// removal payloads exceeding the 4MB payload limit are rejected with HTTP 413 and a structured error envelope.
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

	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "request body exceeds maximum allowed size")
}

// TestInstallExtensionHandler_EmptyBody_ReturnsBadRequest verifies that
// an empty request body returns HTTP 400 Bad Request with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, "empty request body", errBody.Error)
}

// TestRemoveExtensionHandler_EmptyBody_ReturnsBadRequest verifies that
// an empty request body returns HTTP 400 Bad Request with a structured JSON error envelope.
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

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, "empty request body", errBody.Error)
}

// TestInstallExtensionHandler_InvalidJSON_ReturnsBadRequest verifies that
// malformed JSON request bodies return HTTP 400 Bad Request with the standard MeshKit unmarshal error code.
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

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, models.ErrUnmarshalCode, errBody.Code)
	assert.NotEmpty(t, errBody.Error)
	assert.NotEmpty(t, errBody.ProbableCause)
	assert.NotEmpty(t, errBody.SuggestedRemediation)
}

// TestRemoveExtensionHandler_InvalidJSON_ReturnsBadRequest verifies that
// malformed JSON request bodies return HTTP 400 Bad Request with the standard MeshKit unmarshal error code.
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

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, models.ErrUnmarshalCode, errBody.Code)
	assert.NotEmpty(t, errBody.Error)
	assert.NotEmpty(t, errBody.ProbableCause)
	assert.NotEmpty(t, errBody.SuggestedRemediation)
}

// TestInstallExtensionHandler_Success_ReturnsEmptyJSON verifies that installing an
// extension with valid metadata succeeds, returning HTTP 200, Content-Type: application/json,
// and an empty JSON object {}.
func TestInstallExtensionHandler_Success_ReturnsEmptyJSON(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	payload := `{"extType":"account","extensionMetadata":{"title":"Custom Account Extension"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.InstallExtensionHandler(rec, req, nil, nil, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, "{}", strings.TrimSpace(rec.Body.String()))
	require.Len(t, provider.Extensions.Account, 1)
	assert.Equal(t, "Custom Account Extension", provider.Extensions.Account[0].Title)
}

// TestRemoveExtensionHandler_Success_ReturnsEmptyJSON verifies that removing an
// existing extension succeeds, returning HTTP 200, Content-Type: application/json,
// and an empty JSON object {}.
func TestRemoveExtensionHandler_Success_ReturnsEmptyJSON(t *testing.T) {
	provider := newTestLocalProvider()
	err := provider.InstallExtension("account", "", map[string]interface{}{"title": "Custom Account Extension"})
	require.NoError(t, err)
	require.Len(t, provider.Extensions.Account, 1)

	h := newTestHandler(t, map[string]models.Provider{}, "")
	payload := `{"extType":"account","title":"Custom Account Extension"}`
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.RemoveExtensionHandler(rec, req, nil, nil, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, "{}", strings.TrimSpace(rec.Body.String()))
	assert.Empty(t, provider.Extensions.Account)
}

// TestRemoveExtensionHandler_ExtensionNotFound_ReturnsInternalServerError verifies that
// attempting to remove a non-existent extension returns HTTP 500 with a structured error envelope.
func TestRemoveExtensionHandler_ExtensionNotFound_ReturnsInternalServerError(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	payload := `{"extType":"account","title":"NonExistentExtension"}`
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.RemoveExtensionHandler(rec, req, nil, nil, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "failed to remove extension")
	assert.Contains(t, errBody.Error, "not found")
}

// TestInstallExtensionHandler_UnsupportedType_ReturnsInternalServerError verifies that
// attempting to install an unsupported extension type returns HTTP 500 with a structured error envelope.
func TestInstallExtensionHandler_UnsupportedType_ReturnsInternalServerError(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	payload := `{"extType":"unsupported","extensionMetadata":{"title":"MyExt"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.InstallExtensionHandler(rec, req, nil, nil, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "failed to install extension")
	assert.Contains(t, errBody.Error, "unsupported extension type")
}

// TestRemoveExtensionHandler_UnsupportedType_ReturnsInternalServerError verifies that
// attempting to remove an unsupported extension type returns HTTP 500 with a structured error envelope.
func TestRemoveExtensionHandler_UnsupportedType_ReturnsInternalServerError(t *testing.T) {
	provider := newTestLocalProvider()
	h := newTestHandler(t, map[string]models.Provider{}, "")
	payload := `{"extType":"unsupported","title":"MyExt"}`
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/delete", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.RemoveExtensionHandler(rec, req, nil, nil, provider)
	resp := rec.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Contains(t, errBody.Error, "failed to remove extension")
	assert.Contains(t, errBody.Error, "unsupported extension type")
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
