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
	profile *models.PerformanceProfile
}

func (f *fakeLoadTestProvider) GetPerformanceProfile(req *http.Request, profileID string) ([]byte, error) {
	profile := models.PerformanceProfile{
		Metadata: map[string]interface{}{},
	}
	if f.profile != nil {
		profile = *f.profile
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
			name: "Host and port without scheme",
			url:  "example.com:8080",
		},
		{
			name: "Relative Path",
			url:  "/relative/endpoint",
		},
		{
			name: "Empty URL",
			url:  "",
		},
		{
			name: "Malformed URL with parse error",
			url:  "http://[::1]:namedport",
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
	testCases := []struct {
		name string
		url  string
	}{
		{
			name: "Schemeless URL",
			url:  "example.com",
		},
		{
			name: "Host and port without scheme",
			url:  "example.com:8080",
		},
		{
			name: "Malformed URL with parse error",
			url:  "http://[::1]:namedport",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{
				log: newTestLogger(t),
			}

			rawJSON := map[string]interface{}{
				"test": map[string]interface{}{
					"name":     "test-smp",
					"duration": "1s",
					"clients": []map[string]interface{}{
						{
							"endpointUrls": []string{tc.url},
						},
					},
				},
			}

			payloadBytes, err := json.Marshal(rawJSON)
			if err != nil {
				t.Fatalf("failed to marshal SMP config: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/run", bytes.NewReader(payloadBytes))
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
		})
	}
}

// TestLoadTestHandlerSSLCertificateTypeAssertion verifies that LoadTestHandler does not panic
// when ca_certificate in performance profile metadata has missing or invalid types for name or file.
func TestLoadTestHandlerSSLCertificateTypeAssertion(t *testing.T) {
	testCases := []struct {
		name     string
		metadata map[string]interface{}
	}{
		{
			name: "Missing file field in ca_certificate",
			metadata: map[string]interface{}{
				"ca_certificate": map[string]interface{}{
					"name": "ca.crt",
				},
			},
		},
		{
			name: "Nil file field in ca_certificate",
			metadata: map[string]interface{}{
				"ca_certificate": map[string]interface{}{
					"name": "ca.crt",
					"file": nil,
				},
			},
		},
		{
			name: "Non-string file field in ca_certificate",
			metadata: map[string]interface{}{
				"ca_certificate": map[string]interface{}{
					"name": "ca.crt",
					"file": 12345,
				},
			},
		},
		{
			name: "Nil name field in ca_certificate",
			metadata: map[string]interface{}{
				"ca_certificate": map[string]interface{}{
					"name": nil,
					"file": "some-cert-data",
				},
			},
		},
		{
			name: "Non-map ca_certificate",
			metadata: map[string]interface{}{
				"ca_certificate": "invalid-cert-structure",
			},
		},
		{
			name: "Valid ca_certificate with name and file",
			metadata: map[string]interface{}{
				"ca_certificate": map[string]interface{}{
					"name": "ca.crt",
					"file": "-----BEGIN CERTIFICATE-----\nMIIB...\n-----END CERTIFICATE-----",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := &Handler{
				log: newTestLogger(t),
			}

			provider := &fakeLoadTestProvider{
				DefaultLocalProvider: &models.DefaultLocalProvider{},
				profile: &models.PerformanceProfile{
					Metadata: tc.metadata,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/run?cert=true&name=test-test&url=example.com", nil)
			rec := httptest.NewRecorder()

			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("LoadTestHandler panicked on %q: %v", tc.name, r)
				}
			}()

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
