package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

type mockFetchResultsProvider struct {
	*models.DefaultLocalProvider
	fetchResultsFn    func(tokenVal, page, pageSize, search, order, profileID string) ([]byte, error)
	fetchAllResultsFn func(tokenVal, page, pageSize, search, order, from, to string) ([]byte, error)
	getResultFn       func(tokenVal string, resultID core.Uuid) (*models.MesheryResult, error)
	fetchSmiResultsFn func(req *http.Request, page, pageSize, search, order string) ([]byte, error)
	fetchSmiResultFn  func(req *http.Request, page, pageSize, search, order string, resultID core.Uuid) ([]byte, error)
}

func (m *mockFetchResultsProvider) FetchResults(tokenVal, page, pageSize, search, order, profileID string) ([]byte, error) {
	if m.fetchResultsFn != nil {
		return m.fetchResultsFn(tokenVal, page, pageSize, search, order, profileID)
	}
	return m.DefaultLocalProvider.FetchResults(tokenVal, page, pageSize, search, order, profileID)
}

func (m *mockFetchResultsProvider) FetchAllResults(tokenVal, page, pageSize, search, order, from, to string) ([]byte, error) {
	if m.fetchAllResultsFn != nil {
		return m.fetchAllResultsFn(tokenVal, page, pageSize, search, order, from, to)
	}
	return m.DefaultLocalProvider.FetchAllResults(tokenVal, page, pageSize, search, order, from, to)
}

func (m *mockFetchResultsProvider) GetResult(tokenVal string, resultID core.Uuid) (*models.MesheryResult, error) {
	if m.getResultFn != nil {
		return m.getResultFn(tokenVal, resultID)
	}
	return m.DefaultLocalProvider.GetResult(tokenVal, resultID)
}

func (m *mockFetchResultsProvider) FetchSmiResults(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
	if m.fetchSmiResultsFn != nil {
		return m.fetchSmiResultsFn(req, page, pageSize, search, order)
	}
	return m.DefaultLocalProvider.FetchSmiResults(req, page, pageSize, search, order)
}

func (m *mockFetchResultsProvider) FetchSmiResult(req *http.Request, page, pageSize, search, order string, resultID core.Uuid) ([]byte, error) {
	if m.fetchSmiResultFn != nil {
		return m.fetchSmiResultFn(req, page, pageSize, search, order, resultID)
	}
	return m.DefaultLocalProvider.FetchSmiResult(req, page, pageSize, search, order, resultID)
}

type resultsTestProvider struct {
	models.Provider
}

func TestErrMissingResultID_Remediation(t *testing.T) {
	err := ErrMissingResultID()
	if err == nil {
		t.Fatal("expected non-nil error from ErrMissingResultID")
	}

	code := errors.GetCode(err)
	if code != ErrMissingResultIDCode {
		t.Errorf("expected error code %q, got %q", ErrMissingResultIDCode, code)
	}

	remedy := errors.GetRemedy(err)
	expectedRoute := "/api/perf/profile/result/{id}"
	staleRoute := "/api/user/performance/results/{id}"

	if !strings.Contains(remedy, expectedRoute) {
		t.Errorf("expected remediation to reference registered route %q, got %q", expectedRoute, remedy)
	}
	if strings.Contains(remedy, staleRoute) {
		t.Errorf("remediation must not reference stale route %q, got %q", staleRoute, remedy)
	}
}

func TestGetResultHandler_MissingResultID(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/", nil)
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}

	body := rec.Body.String()
	if !strings.Contains(body, ErrMissingResultIDCode) {
		t.Errorf("expected response body to contain error code %q, got %s", ErrMissingResultIDCode, body)
	}
	if !strings.Contains(body, "/api/perf/profile/result/{id}") {
		t.Errorf("expected response body to contain updated remediation route, got %s", body)
	}

	var jsonBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &jsonBody); err != nil {
		t.Fatalf("failed to decode response as JSON: %v", err)
	}
}

func TestGetResultHandler_InvalidUUID(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/invalid-uuid-string", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid-string"})
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestGetResultHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	validID := "11111111-1111-1111-1111-111111111111"
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/"+validID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validID})
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestFetchAllResultsHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/results", nil)
	rec := httptest.NewRecorder()

	h.FetchAllResultsHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestFetchSingleSmiResultHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	validID := "11111111-1111-1111-1111-111111111111"
	req := httptest.NewRequest(http.MethodGet, "/api/perf/service-mesh-interface/result/"+validID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validID})
	rec := httptest.NewRecorder()

	h.FetchSingleSmiResultHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestFetchSingleSmiResultHandler_InvalidUUID_ReturnsBadRequest(t *testing.T) {
	fakeProvider := &resultsTestProvider{}
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/service-mesh-interface/result/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})
	rec := httptest.NewRecorder()

	h.FetchSingleSmiResultHandler(rec, req, nil, nil, fakeProvider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestFetchResultsHandler_NilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result", nil)
	rec := httptest.NewRecorder()

	h.FetchResultsHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestFetchResultsHandler_MissingToken(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result", nil)
	rec := httptest.NewRecorder()

	h.FetchResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchTokenCode, rec.Body.String())
	}
}

func TestFetchResultsHandler_FormParseError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result?invalid=%zz", nil)
	rec := httptest.NewRecorder()

	h.FetchResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on form parse failure, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrParseFormCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrParseFormCode, rec.Body.String())
	}
}

func TestFetchResultsHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchResultsFn: func(tokenVal, page, pageSize, search, order, profileID string) ([]byte, error) {
			return []byte(`{"results":[],"total_count":0}`), nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result", nil)
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.FetchResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != `{"results":[],"total_count":0}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestFetchResultsHandler_EmptyToken_AllowedForLocalProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchResultsFn: func(tokenVal, page, pageSize, search, order, profileID string) ([]byte, error) {
			if tokenVal != "" {
				t.Errorf("expected empty token for local provider, got %q", tokenVal)
			}
			return []byte(`{"results":[],"total_count":0}`), nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result", nil)
	// Local provider context has empty string token injected by SessionInjectorMiddleware
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.FetchResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != `{"results":[],"total_count":0}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestFetchAllResultsHandler_MissingToken(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/results", nil)
	rec := httptest.NewRecorder()

	h.FetchAllResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchTokenCode, rec.Body.String())
	}
}

func TestFetchAllResultsHandler_FormParseError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/results?invalid=%zz", nil)
	rec := httptest.NewRecorder()

	h.FetchAllResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on form parse failure, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrParseFormCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrParseFormCode, rec.Body.String())
	}
}

func TestFetchAllResultsHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchAllResultsFn: func(tokenVal, page, pageSize, search, order, from, to string) ([]byte, error) {
			return []byte(`{"results":[],"page":1}`), nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/results", nil)
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.FetchAllResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != `{"results":[],"page":1}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestGetResultHandler_MissingToken(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	validUUID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/"+validUUID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID})
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchTokenCode, rec.Body.String())
	}
}

func TestGetResultHandler_ResultNotFound_NilResult(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getResultFn: func(tokenVal string, resultID core.Uuid) (*models.MesheryResult, error) {
			return nil, nil
		},
	}
	validUUID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/"+validUUID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrGetResultCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrGetResultCode, rec.Body.String())
	}
}

func TestGetResultHandler_ResultNotFound_GormRecordNotFound(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getResultFn: func(tokenVal string, resultID core.Uuid) (*models.MesheryResult, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	validUUID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/"+validUUID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrGetResultCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrGetResultCode, rec.Body.String())
	}
}

func TestGetResultHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	validUUID := uuid.Must(uuid.NewV4())
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getResultFn: func(tokenVal string, resultID core.Uuid) (*models.MesheryResult, error) {
			return &models.MesheryResult{
				ID:   validUUID,
				Name: "test-result",
				Mesh: "Istio",
				Result: map[string]interface{}{
					"RunType": "HTTP",
					"DurationHistogram": map[string]interface{}{
						"Min": 0.001,
						"Max": 0.05,
						"Avg": 0.02,
					},
				},
			}, nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/result/"+validUUID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID.String()})
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Header().Get("content-type") != "application/yaml" {
		t.Errorf("expected application/yaml, got %s", rec.Header().Get("content-type"))
	}
}

func TestFetchSmiResultsHandler_NilProvider(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results", nil)
	rec := httptest.NewRecorder()

	h.FetchSmiResultsHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestFetchSmiResultsHandler_FormParseError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results?invalid=%zz", nil)
	rec := httptest.NewRecorder()

	h.FetchSmiResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on form parse failure, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrParseFormCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrParseFormCode, rec.Body.String())
	}
}

func TestFetchSmiResultsHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchSmiResultsFn: func(req *http.Request, page, pageSize, search, order string) ([]byte, error) {
			return []byte(`{"results":[],"page":1}`), nil
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results", nil)
	rec := httptest.NewRecorder()

	h.FetchSmiResultsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != `{"results":[],"page":1}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestFetchSingleSmiResultHandler_FormParseError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results/some-id?invalid=%zz", nil)
	rec := httptest.NewRecorder()

	h.FetchSingleSmiResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request on form parse failure, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrParseFormCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrParseFormCode, rec.Body.String())
	}
}

func TestFetchSingleSmiResultHandler_NotFound(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchSmiResultFn: func(req *http.Request, page, pageSize, search, order string, resultID core.Uuid) ([]byte, error) {
			return []byte{}, nil
		},
	}
	validUUID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results/"+validUUID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID})
	rec := httptest.NewRecorder()

	h.FetchSingleSmiResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchSMIResultsCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchSMIResultsCode, rec.Body.String())
	}
}

func TestFetchSingleSmiResultHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &mockFetchResultsProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		fetchSmiResultFn: func(req *http.Request, page, pageSize, search, order string, resultID core.Uuid) ([]byte, error) {
			return []byte(`{"status":"passed"}`), nil
		},
	}
	validUUID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodGet, "/api/perf/smi/results/"+validUUID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": validUUID})
	rec := httptest.NewRecorder()

	h.FetchSingleSmiResultHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != `{"status":"passed"}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}
