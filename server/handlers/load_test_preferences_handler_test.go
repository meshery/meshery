package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLoadTestPreferencesProvider struct {
	*models.DefaultLocalProvider
	storeFunc func(req *http.Request, perfConfig *perfprofile.PerformanceTestConfig) (string, error)
	getFunc   func(req *http.Request, testUUID string) (*perfprofile.PerformanceTestConfig, error)
	fetchFunc func(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	delFunc   func(req *http.Request, testUUID string) error
}

func (m *mockLoadTestPreferencesProvider) SMPTestConfigStore(req *http.Request, perfConfig *perfprofile.PerformanceTestConfig) (string, error) {
	if m.storeFunc != nil {
		return m.storeFunc(req, perfConfig)
	}
	return "test-uuid-default", nil
}

func (m *mockLoadTestPreferencesProvider) SMPTestConfigGet(req *http.Request, testUUID string) (*perfprofile.PerformanceTestConfig, error) {
	if m.getFunc != nil {
		return m.getFunc(req, testUUID)
	}
	return &perfprofile.PerformanceTestConfig{ID: testUUID, Name: "TestConfig"}, nil
}

func (m *mockLoadTestPreferencesProvider) SMPTestConfigFetch(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if m.fetchFunc != nil {
		return m.fetchFunc(req, page, pageSize, search, order)
	}
	return []byte(`[{"id":"test-uuid-default","name":"TestConfig"}]`), nil
}

func (m *mockLoadTestPreferencesProvider) SMPTestConfigDelete(req *http.Request, testUUID string) error {
	if m.delFunc != nil {
		return m.delFunc(req, testUUID)
	}
	return nil
}

func validPerformanceTestConfigJSON() string {
	cfg := perfprofile.PerformanceTestConfig{
		Name:     "Sample Load Test",
		Duration: "10s",
		Clients: []perfprofile.PerformanceTestClient{
			{
				Protocol:     "http",
				EndpointUrls: []string{"https://meshery.io/api/health"},
			},
		},
	}
	b, _ := json.Marshal(cfg)
	return string(b)
}

func TestUserTestPreferenceHandler_MethodDispatcher(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{}

	t.Run("POST dispatches to Store", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(validPerformanceTestConfigJSON()))
		rec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(rec, req, nil, nil, provider)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET dispatches to Get", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid=123", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(rec, req, nil, nil, provider)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE dispatches to Delete", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid=123", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(rec, req, nil, nil, provider)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Unsupported method returns 405 MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/user/prefs/perf", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(rec, req, nil, nil, provider)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
		var body map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &body)
		require.NoError(t, err)
		assert.Equal(t, "method not allowed", body["error"])
	})
}

func TestUserTestPreferenceStore_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(validPerformanceTestConfigJSON()))
	rec := httptest.NewRecorder()

	h.UserTestPreferenceStore(rec, req, nil, nil, nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrNilProviderCode, errBody.Code)
}

func TestUserTestPreferenceStore_TypedNilProvider_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(validPerformanceTestConfigJSON()))
	rec := httptest.NewRecorder()
	var typedNil *models.DefaultLocalProvider

	h.UserTestPreferenceStore(rec, req, nil, nil, typedNil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrNilProviderCode, errBody.Code)
}

func TestUserTestPreferenceStore_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{log: nil}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader("invalid-json"))
	rec := httptest.NewRecorder()
	provider := &mockLoadTestPreferencesProvider{}

	assert.NotPanics(t, func() {
		h.UserTestPreferenceStore(rec, req, nil, nil, provider)
	})
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestUserTestPreferenceStore_PayloadTooLarge_ReturnsEntityTooLarge(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	oversized := bytes.Repeat([]byte("a"), maxPerformancePayloadSize+10)
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", bytes.NewReader(oversized))
	rec := httptest.NewRecorder()
	provider := &mockLoadTestPreferencesProvider{}

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrRequestBodyCode, errBody.Code)
	assert.Contains(t, errBody.Error, "request body exceeds maximum allowed size")
}

func TestUserTestPreferenceStore_EmptyBody_ReturnsBadRequest(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(""))
	rec := httptest.NewRecorder()
	provider := &mockLoadTestPreferencesProvider{}

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrRequestBodyCode, errBody.Code)
}

func TestUserTestPreferenceStore_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(`{"name":`))
	rec := httptest.NewRecorder()
	provider := &mockLoadTestPreferencesProvider{}

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, models.ErrUnmarshalCode, errBody.Code)
}

func TestUserTestPreferenceStore_ValidationError_ReturnsBadRequest(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(`{"name":""}`))
	rec := httptest.NewRecorder()
	provider := &mockLoadTestPreferencesProvider{}

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrRecordPreferencesCode, errBody.Code)
}

func TestUserTestPreferenceStore_Success(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		storeFunc: func(req *http.Request, perfConfig *perfprofile.PerformanceTestConfig) (string, error) {
			assert.Equal(t, "Sample Load Test", perfConfig.Name)
			return "test-uuid-777", nil
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(validPerformanceTestConfigJSON()))
	rec := httptest.NewRecorder()

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	var res map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)
	assert.Equal(t, "test-uuid-777", res["test_uuid"])
}

func TestUserTestPreferenceStore_StoreError_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		storeFunc: func(req *http.Request, perfConfig *perfprofile.PerformanceTestConfig) (string, error) {
			return "", errors.New("db error")
		},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", strings.NewReader(validPerformanceTestConfigJSON()))
	rec := httptest.NewRecorder()

	h.UserTestPreferenceStore(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrFailToSaveCode, errBody.Code)
}

func TestUserTestPreferenceGet_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrNilProviderCode, errBody.Code)
}

func TestUserTestPreferenceGet_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{log: nil}
	provider := &mockLoadTestPreferencesProvider{}

	assert.NotPanics(t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid=123", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceGet(rec, req, nil, nil, provider)
	})

	assert.NotPanics(t, func() {
		req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?page=1&pagesize=10", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceGet(rec, req, nil, nil, provider)
	})
}

func TestUserTestPreferenceGet_AllConfigs_Success(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		fetchFunc: func(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
			assert.Equal(t, "1", page)
			assert.Equal(t, "10", pageSize)
			return []byte(`[{"id":"uuid-1","name":"Test 1"}]`), nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?page=1&pagesize=10", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	assert.JSONEq(t, `[{"id":"uuid-1","name":"Test 1"}]`, rec.Body.String())
}

func TestUserTestPreferenceGet_AllConfigs_FetchError_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		fetchFunc: func(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
			return nil, errors.New("fetch error")
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrTestConfigsCode, errBody.Code)
}

func TestUserTestPreferenceGet_SingleConfig_Success(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	expected := &perfprofile.PerformanceTestConfig{
		ID:   "cfg-uuid-1",
		Name: "Specific Test",
	}
	provider := &mockLoadTestPreferencesProvider{
		getFunc: func(req *http.Request, testUUID string) (*perfprofile.PerformanceTestConfig, error) {
			assert.Equal(t, "cfg-uuid-1", testUUID)
			return expected, nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid=cfg-uuid-1", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	var result perfprofile.PerformanceTestConfig
	err := json.Unmarshal(rec.Body.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.Name, result.Name)
}

func TestUserTestPreferenceGet_SingleConfig_NotFound(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		getFunc: func(req *http.Request, testUUID string) (*perfprofile.PerformanceTestConfig, error) {
			return nil, nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid=missing-uuid", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrTestConfigsCode, errBody.Code)
}

func TestUserTestPreferenceGet_SingleConfig_GetError_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		getFunc: func(req *http.Request, testUUID string) (*perfprofile.PerformanceTestConfig, error) {
			return nil, errors.New("get failed")
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid=err-uuid", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceGet(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrTestConfigsCode, errBody.Code)
}

func TestUserTestPreferenceDelete_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid=123", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceDelete(rec, req, nil, nil, nil)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrNilProviderCode, errBody.Code)
}

func TestUserTestPreferenceDelete_NilLogger_DoesNotPanic(t *testing.T) {
	h := &Handler{log: nil}
	provider := &mockLoadTestPreferencesProvider{}

	assert.NotPanics(t, func() {
		req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf", nil)
		rec := httptest.NewRecorder()
		h.UserTestPreferenceDelete(rec, req, nil, nil, provider)
	})
}

func TestUserTestPreferenceDelete_MissingUUID_ReturnsBadRequest(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{}
	req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceDelete(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrQueryGetCode, errBody.Code)
}

func TestUserTestPreferenceDelete_Success(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		delFunc: func(req *http.Request, testUUID string) error {
			assert.Equal(t, "uuid-to-delete", testUUID)
			return nil
		},
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid=uuid-to-delete", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceDelete(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	assert.Equal(t, "{}", strings.TrimSpace(rec.Body.String()))
}

func TestUserTestPreferenceDelete_DeleteError_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	provider := &mockLoadTestPreferencesProvider{
		delFunc: func(req *http.Request, testUUID string) error {
			return errors.New("delete failed in db")
		},
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid=err-uuid", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceDelete(rec, req, nil, nil, provider)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	errBody := decodeErrorBody(t, rec.Body.Bytes())
	assert.Equal(t, ErrFailToDeleteCode, errBody.Code)
}
