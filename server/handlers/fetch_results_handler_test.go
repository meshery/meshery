package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshery/server/models"
)

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
