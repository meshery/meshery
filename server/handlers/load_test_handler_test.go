package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

type loadTestSpyProvider struct {
	*models.DefaultLocalProvider
	getProfileCalled bool
	getProfileResp   []byte
	getProfileErr    error
}

func (m *loadTestSpyProvider) GetPerformanceProfile(_ *http.Request, _ string) ([]byte, error) {
	m.getProfileCalled = true
	return m.getProfileResp, m.getProfileErr
}

// TestLoadTestHandler_WhitespaceBodyIgnored verifies that a request body
// containing only whitespace (such as trailing newlines or spaces) is not
// treated as an SMP configuration payload, and instead proceeds to the
// regular performance profile flow (calling provider.GetPerformanceProfile).
func TestLoadTestHandler_WhitespaceBodyIgnored(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &loadTestSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		getProfileResp:       []byte(`{"id":"test-profile-id","name":"test","metadata":{}}`),
	}

	testCases := []struct {
		name string
		body string
	}{
		{
			name: "empty body",
			body: "",
		},
		{
			name: "spaces only",
			body: "   ",
		},
		{
			name: "newline only",
			body: "\n",
		},
		{
			name: "carriage return and newline",
			body: "\r\n",
		},
		{
			name: "mixed whitespace tabs and newlines",
			body: " \t \r\n \n\t ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider.getProfileCalled = false
			req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/test-profile-id/run", bytes.NewBufferString(tc.body))
			req = mux.SetURLVars(req, map[string]string{"id": "test-profile-id"})
			rec := httptest.NewRecorder()

			h.LoadTestHandler(rec, req, &models.Preference{}, &models.User{}, provider)

			if !provider.getProfileCalled {
				t.Fatalf("expected GetPerformanceProfile to be called for whitespace body %q, but it was not (handler likely short-circuited to SMP handler)", tc.body)
			}
		})
	}
}

// TestLoadTestHandler_NonWhitespaceBodyDispatchesToSMP verifies that a request
// with actual non-whitespace body content is dispatched to the SMP handler.
func TestLoadTestHandler_NonWhitespaceBodyDispatchesToSMP(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := &loadTestSpyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/test-profile-id/run", bytes.NewBufferString("not-empty-smp-content"))
	req = mux.SetURLVars(req, map[string]string{"id": "test-profile-id"})
	rec := httptest.NewRecorder()

	h.LoadTestHandler(rec, req, &models.Preference{}, &models.User{}, provider)

	if provider.getProfileCalled {
		t.Fatalf("expected GetPerformanceProfile NOT to be called when body has content, but it was called")
	}
}
