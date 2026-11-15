package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	mkerrors "github.com/meshery/meshkit/errors"
	"gorm.io/gorm"
)

func TestSavePerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	body := []byte(`{"name":"test-profile"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.SavePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestDeletePerformanceProfileHandler_MissingID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/", nil)
	rec := httptest.NewRecorder()

	h.DeletePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "performanceProfileId is required") {
		t.Errorf("expected 'performanceProfileId is required', got %s", rec.Body.String())
	}
}

func TestDeletePerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodDelete, "/api/perf/profile/profile-123", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
	rec := httptest.NewRecorder()

	h.DeletePerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestDeletePerformanceProfileHandler_InvalidUUID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()

	req := httptest.NewRequest(http.MethodDelete, "/api/user/performance/profiles/not-a-valid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "not-a-valid-uuid"})
	rec := httptest.NewRecorder()

	h.DeletePerformanceProfileHandler(rec, req, nil, nil, base)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request on invalid UUID, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

// deletePerfProfileErrProvider fails performance profile deletes with err.
type deletePerfProfileErrProvider struct {
	*models.DefaultLocalProvider
	err error
}

func (p *deletePerfProfileErrProvider) DeletePerformanceProfile(_ *http.Request, _ string) ([]byte, error) {
	return nil, p.err
}

func TestDeletePerformanceProfileHandler_NotFoundIs404(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"local provider not found", models.ErrResultNotFound(fmt.Errorf("performance profile not found")), http.StatusNotFound},
		{"remote provider 404", models.ErrFetch(fmt.Errorf("failed to delete the profile"), "not found", http.StatusNotFound), http.StatusNotFound},
		{"other failure", fmt.Errorf("database is locked"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &deletePerfProfileErrProvider{DefaultLocalProvider: base, err: tc.err}

			id := "abcdefab-abcd-abcd-abcd-abcdefabcdef"
			req := httptest.NewRequest(http.MethodDelete, "/api/user/performance/profiles/"+id, nil)
			req = mux.SetURLVars(req, map[string]string{"performanceProfileId": id})
			rec := httptest.NewRecorder()

			h.DeletePerformanceProfileHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}

			// The body must be a MeshKit error, not just carry the right status.
			// A not-found error is written as is; anything else is wrapped in
			// ErrFailToDelete.
			wantErrCode := mkerrors.GetCode(tc.err)
			if tc.wantCode == http.StatusInternalServerError {
				wantErrCode = mkerrors.GetCode(ErrFailToDelete(tc.err, "performance profile"))
			}
			var body struct {
				Error string `json:"error"`
				Code  string `json:"code"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("response body is not a JSON error: %v (body=%q)", err, rec.Body.String())
			}
			if body.Code != wantErrCode {
				t.Errorf("error code = %q, want %q", body.Code, wantErrCode)
			}
			if body.Error == "" {
				t.Errorf("error message is empty (body=%q)", rec.Body.String())
			}
		})
	}
}

// A failure the remote provider answered with an HTTP status must reach the
// client with that status, not a flattened 500. Only a failure that carries no
// status (unreachable provider, local error) falls back to 500.
func TestDeletePerformanceProfileHandler_PreservesProviderStatus(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"remote provider 403", models.ErrDelete(fmt.Errorf("forbidden"), "Perf Profile :id", http.StatusForbidden), http.StatusForbidden},
		{"remote provider 409", models.ErrDelete(fmt.Errorf("conflict"), "Perf Profile :id", http.StatusConflict), http.StatusConflict},
		{"remote provider 500", models.ErrDelete(fmt.Errorf("boom"), "Perf Profile :id", http.StatusInternalServerError), http.StatusInternalServerError},
		{"provider unreachable, no status", models.ErrUnreachableRemoteProvider(fmt.Errorf("dial tcp: connection refused")), http.StatusInternalServerError},
		{"plain error, no status", fmt.Errorf("database is locked"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &deletePerfProfileErrProvider{DefaultLocalProvider: base, err: tc.err}

			id := "abcdefab-abcd-abcd-abcd-abcdefabcdef"
			req := httptest.NewRequest(http.MethodDelete, "/api/user/performance/profiles/"+id, nil)
			req = mux.SetURLVars(req, map[string]string{"performanceProfileId": id})
			rec := httptest.NewRecorder()

			h.DeletePerformanceProfileHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}

			// Whatever the status, the body is the ErrFailToDelete MeshKit error.
			var body struct {
				Error string `json:"error"`
				Code  string `json:"code"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("response body is not a JSON error: %v (body=%q)", err, rec.Body.String())
			}
			if want := mkerrors.GetCode(ErrFailToDelete(tc.err, "performance profile")); body.Code != want {
				t.Errorf("error code = %q, want %q", body.Code, want)
			}
		})
	}
}

func TestGetPerformanceProfileHandler_MissingID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/", nil)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "performanceProfileId is required") {
		t.Errorf("expected 'performanceProfileId is required', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfileHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/profile-123", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "profile-123"})
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "provider not initialized") {
		t.Errorf("expected 'provider not initialized', got %s", rec.Body.String())
	}
}

func TestGetPerformanceProfileHandler_InvalidUUID_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()

	req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/not-a-valid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": "not-a-valid-uuid"})
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, base)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request on invalid UUID, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
		t.Errorf("expected error body to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
	}
}

// getPerfProfileMockProvider stubs GetPerformanceProfile.
type getPerfProfileMockProvider struct {
	*models.DefaultLocalProvider
	resp []byte
	err  error
}

func (p *getPerfProfileMockProvider) GetPerformanceProfile(_ *http.Request, _ string) ([]byte, error) {
	return p.resp, p.err
}

func TestGetPerformanceProfileHandler_NotFoundIs404(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"local provider ErrResultNotFound", models.ErrResultNotFound(fmt.Errorf("performance profile not found")), http.StatusNotFound},
		{"gorm record not found", gorm.ErrRecordNotFound, http.StatusNotFound},
		{"remote provider 404", models.ErrFetch(fmt.Errorf("profile not found"), "not found", http.StatusNotFound), http.StatusNotFound},
		{"other failure", fmt.Errorf("database connection failed"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &getPerfProfileMockProvider{DefaultLocalProvider: base, err: tc.err}

			id := "abcdefab-abcd-abcd-abcd-abcdefabcdef"
			req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/"+id, nil)
			req = mux.SetURLVars(req, map[string]string{"performanceProfileId": id})
			rec := httptest.NewRecorder()

			h.GetPerformanceProfileHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}

			var body struct {
				Error string `json:"error"`
				Code  string `json:"code"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("response body is not a JSON error: %v (body=%q)", err, rec.Body.String())
			}
			if body.Error == "" {
				t.Errorf("error message is empty (body=%q)", rec.Body.String())
			}
		})
	}
}

func TestGetPerformanceProfileHandler_PreservesProviderStatus(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"remote provider 403", models.ErrFetch(fmt.Errorf("forbidden"), "profile", http.StatusForbidden), http.StatusForbidden},
		{"remote provider 409", models.ErrFetch(fmt.Errorf("conflict"), "profile", http.StatusConflict), http.StatusConflict},
		{"remote provider 500", models.ErrFetch(fmt.Errorf("internal"), "profile", http.StatusInternalServerError), http.StatusInternalServerError},
		{"plain error", fmt.Errorf("database timeout"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &getPerfProfileMockProvider{DefaultLocalProvider: base, err: tc.err}

			id := "abcdefab-abcd-abcd-abcd-abcdefabcdef"
			req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/"+id, nil)
			req = mux.SetURLVars(req, map[string]string{"performanceProfileId": id})
			rec := httptest.NewRecorder()

			h.GetPerformanceProfileHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}
		})
	}
}

func TestGetPerformanceProfileHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	mockResp := []byte(`{"id":"abcdefab-abcd-abcd-abcd-abcdefabcdef","name":"profile-test"}`)
	provider := &getPerfProfileMockProvider{DefaultLocalProvider: base, resp: mockResp}

	id := "abcdefab-abcd-abcd-abcd-abcdefabcdef"
	req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles/"+id, nil)
	req = mux.SetURLVars(req, map[string]string{"performanceProfileId": id})
	rec := httptest.NewRecorder()

	h.GetPerformanceProfileHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != string(mockResp) {
		t.Errorf("expected body %s, got %s", string(mockResp), rec.Body.String())
	}
}

type perfProfilesMockProvider struct {
	*models.DefaultLocalProvider
	getProfilesFn func(tokenVal, page, pageSize, search, order string) ([]byte, error)
	saveProfileFn func(tokenVal string, profile *models.PerformanceProfile) ([]byte, error)
}

func (p *perfProfilesMockProvider) GetPerformanceProfiles(tokenVal, page, pageSize, search, order string) ([]byte, error) {
	if p.getProfilesFn != nil {
		return p.getProfilesFn(tokenVal, page, pageSize, search, order)
	}
	return p.DefaultLocalProvider.GetPerformanceProfiles(tokenVal, page, pageSize, search, order)
}

func (p *perfProfilesMockProvider) SavePerformanceProfile(tokenVal string, profile *models.PerformanceProfile) ([]byte, error) {
	if p.saveProfileFn != nil {
		return p.saveProfileFn(tokenVal, profile)
	}
	return p.DefaultLocalProvider.SavePerformanceProfile(tokenVal, profile)
}

type remotePerfProfilesMockProvider struct {
	*models.DefaultLocalProvider
	getProfilesFn func(tokenVal, page, pageSize, search, order string) ([]byte, error)
}

func (p *remotePerfProfilesMockProvider) GetProviderType() models.ProviderType {
	return models.RemoteProviderType
}

func (p *remotePerfProfilesMockProvider) GetPerformanceProfiles(tokenVal, page, pageSize, search, order string) ([]byte, error) {
	if p.getProfilesFn != nil {
		return p.getProfilesFn(tokenVal, page, pageSize, search, order)
	}
	return []byte(`{"profiles":[],"total_count":0}`), nil
}

func TestGetPerformanceProfilesHandler_LocalProvider_EmptyToken_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	mockResp := []byte(`{"profiles":[],"total_count":0}`)
	provider := &perfProfilesMockProvider{
		DefaultLocalProvider: base,
		getProfilesFn: func(tokenVal, page, pageSize, search, order string) ([]byte, error) {
			return mockResp, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK for local provider with empty token, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != string(mockResp) {
		t.Errorf("expected body %s, got %s", string(mockResp), rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_RemoteProvider_MissingToken(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	provider := &remotePerfProfilesMockProvider{DefaultLocalProvider: base}

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for remote provider without token, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchTokenCode, rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_RemoteProvider_EmptyToken(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	provider := &remotePerfProfilesMockProvider{DefaultLocalProvider: base}

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for remote provider with empty token, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrFetchTokenCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrFetchTokenCode, rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	mockResp := []byte(`{"profiles":[],"total_count":0}`)
	provider := &perfProfilesMockProvider{
		DefaultLocalProvider: base,
		getProfilesFn: func(tokenVal, page, pageSize, search, order string) ([]byte, error) {
			return mockResp, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile", nil)
	ctx := context.WithValue(req.Context(), models.TokenCtxKey, "valid-token")
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if rec.Body.String() != string(mockResp) {
		t.Errorf("expected body %s, got %s", string(mockResp), rec.Body.String())
	}
}

func TestGetPerformanceProfilesHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	cases := []struct {
		name     string
		query    string
		provider error
	}{
		{"page", "page=-1", models.ErrPageNumber(fmt.Errorf("bad page"))},
		{"pagesize", "pagesize=abc", models.ErrPageSize(fmt.Errorf("bad pagesize"))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &perfProfilesMockProvider{
				DefaultLocalProvider: base,
				getProfilesFn: func(tokenVal, page, pageSize, search, order string) ([]byte, error) {
					return nil, tc.provider
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/user/performance/profiles?"+tc.query, nil)
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "valid-token"))
			rec := httptest.NewRecorder()

			h.GetPerformanceProfilesHandler(rec, req, nil, nil, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestSavePerformanceProfileHandler_MalformedJSON_ReturnsBadRequest(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()

	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader([]byte("{invalid-json")))
	rec := httptest.NewRecorder()

	h.SavePerformanceProfileHandler(rec, req, nil, nil, base)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected error body to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestSavePerformanceProfileHandler_WithUserOwner(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	mockResp := []byte(`{"id":"abcdefab-abcd-abcd-abcd-abcdefabcdef","name":"profile-test"}`)
	var capturedProfile *models.PerformanceProfile
	provider := &perfProfilesMockProvider{
		DefaultLocalProvider: base,
		saveProfileFn: func(tokenVal string, profile *models.PerformanceProfile) ([]byte, error) {
			capturedProfile = profile
			return mockResp, nil
		},
	}

	userID, _ := uuid.FromString("11111111-2222-3333-4444-555555555555")
	user := &models.User{ID: userID}
	body := []byte(`{"name":"profile-test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.SavePerformanceProfileHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if capturedProfile == nil || capturedProfile.Owner != userID {
		t.Errorf("expected parsedBody.Owner to match user.ID (%v), got %v", userID, capturedProfile)
	}
}

func TestSavePerformanceProfileHandler_PreservesProviderStatus(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
	}{
		{"remote provider 403", models.ErrPost(fmt.Errorf("forbidden"), "perf profile", http.StatusForbidden), http.StatusForbidden},
		{"remote provider 409", models.ErrPost(fmt.Errorf("conflict"), "perf profile", http.StatusConflict), http.StatusConflict},
		{"remote provider 500", models.ErrPost(fmt.Errorf("internal"), "perf profile", http.StatusInternalServerError), http.StatusInternalServerError},
		{"plain error", fmt.Errorf("storage failure"), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			base := &models.DefaultLocalProvider{}
			base.Initialize()
			provider := &perfProfilesMockProvider{
				DefaultLocalProvider: base,
				saveProfileFn: func(tokenVal string, profile *models.PerformanceProfile) ([]byte, error) {
					return nil, tc.err
				},
			}

			body := []byte(`{"name":"profile-test"}`)
			req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			h.SavePerformanceProfileHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}
		})
	}
}

func TestSavePerformanceProfileHandler_BodyExceedsLimit_ReturnsRequestEntityTooLarge(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	called := false
	provider := &perfProfilesMockProvider{
		DefaultLocalProvider: base,
		saveProfileFn: func(tokenVal string, profile *models.PerformanceProfile) ([]byte, error) {
			called = true
			return []byte(`{}`), nil
		},
	}

	const maxBodyLimit = 4 << 20 // 4MB
	oversizedBody := bytes.Repeat([]byte("a"), maxBodyLimit+10)
	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.SavePerformanceProfileHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if called {
		t.Error("provider.SavePerformanceProfile must NOT be called when body exceeds maxBodyLimit")
	}
}
