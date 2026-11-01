package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

type fakeLoadTestProvider struct {
	*models.DefaultLocalProvider
}

func (f *fakeLoadTestProvider) GetPerformanceProfile(req *http.Request, profileID string) ([]byte, error) {
	profile := models.PerformanceProfile{
		Metadata: map[string]interface{}{},
	}
	return json.Marshal(profile)
}

// TestErrParseBoolNilError verifies that ErrParseBool handles a nil error gracefully
// without panicking on err.Error().
func TestErrParseBoolNilError(t *testing.T) {
	err := ErrParseBool(nil, "sample-bool-field")
	if err == nil {
		t.Fatal("expected non-nil error from ErrParseBool, got nil")
	}
	if meshkiterrors.GetCode(err) != ErrParseBoolCode {
		t.Fatalf("expected code %s, got %s", ErrParseBoolCode, meshkiterrors.GetCode(err))
	}
}

// TestLoadTestHandlerInvalidURL verifies that LoadTestHandler returns HTTP 400 with
// ErrInvalidLTURLCode when supplied with a schemeless or relative URL, without panicking.
func TestLoadTestHandlerInvalidURL(t *testing.T) {
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "Schemeless URL",
			url:  "example.com",
		},
		{
			name: "Relative Path",
			url:  "/relative/endpoint",
		},
		{
			name: "Empty URL",
			url:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{
				log: newTestLogger(t),
			}

			provider := &fakeLoadTestProvider{
				DefaultLocalProvider: &models.DefaultLocalProvider{},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/run?name=test-test&url="+tc.url, nil)
			rec := httptest.NewRecorder()

			h.LoadTestHandler(rec, req, nil, nil, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d. Body: %s", rec.Code, rec.Body.String())
			}

			var body struct {
				Code  string `json:"code"`
				Error string `json:"error"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("failed to decode response JSON: %v", err)
			}

			if body.Code != ErrInvalidLTURLCode {
				t.Fatalf("expected error code %s, got %s", ErrInvalidLTURLCode, body.Code)
			}
		})
	}
}

// TestLoadTestUsingSMPHandlerInvalidURL verifies that LoadTestUsingSMPHandler rejects
// non-absolute or malformed URLs in the SMP performance configuration with HTTP 400.
func TestLoadTestUsingSMPHandlerInvalidURL(t *testing.T) {
	h := &Handler{
		log: newTestLogger(t),
	}

	rawJSON := `{
		"test": {
			"name": "test-smp",
			"duration": "1s",
			"clients": [
				{
					"endpointUrls": ["example.com"]
				}
			]
		}
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/run", bytes.NewBufferString(rawJSON))
	rec := httptest.NewRecorder()

	h.LoadTestUsingSMPHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if body.Code != ErrInvalidLTURLCode {
		t.Fatalf("expected error code %s, got %s", ErrInvalidLTURLCode, body.Code)
	}
}
