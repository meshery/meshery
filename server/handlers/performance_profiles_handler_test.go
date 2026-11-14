package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	mkerrors "github.com/meshery/meshkit/errors"
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
