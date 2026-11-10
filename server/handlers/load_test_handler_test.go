package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// mockLoadTestProvider provides configurable performance profile responses for unit tests.
type mockLoadTestProvider struct {
	models.Provider
	profileData []byte
	profileErr  error
}

// GetPerformanceProfile returns mock performance profile bytes or configured error.
func (m *mockLoadTestProvider) GetPerformanceProfile(_ *http.Request, _ string) ([]byte, error) {
	if m.profileErr != nil {
		return nil, m.profileErr
	}
	return m.profileData, nil
}

// TestLoadTestHandler_CertificateTempFileCleanedUpOnExit verifies that temporary certificate
// files created when cert=true are properly closed and unlinked upon handler exit.
func TestLoadTestHandler_CertificateTempFileCleanedUpOnExit(t *testing.T) {
	certContent := "-----BEGIN CERTIFICATE-----\ntest-cert-data\n-----END CERTIFICATE-----"
	profile := models.PerformanceProfile{
		Name: "test-profile",
		Metadata: map[string]interface{}{
			"ca_certificate": map[string]interface{}{
				"name": "meshery-test-ca.crt",
				"file": certContent,
			},
		},
	}
	profileBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	provider := &mockLoadTestProvider{profileData: profileBytes}
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/test-profile-id/run?cert=true&name=test-run&url=:", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "test-profile-id",
	})
	rec := httptest.NewRecorder()

	h.LoadTestHandler(rec, req, nil, nil, provider)

	// Verify that any temp certificate created with pattern "meshery-test-ca-*.crt" was deleted
	tempDir := os.TempDir()
	matches, err := filepath.Glob(filepath.Join(tempDir, "meshery-test-ca-*.crt"))
	require.NoError(t, err)
	assert.Empty(t, matches, "expected temporary certificate file to be removed upon handler exit")
}

// TestLoadTestHandler_CertificateInvalidFile_NoPanic ensures LoadTestHandler returns HTTP 400
// without panicking when certificate file data is not a string, and cleans up any temp files.
func TestLoadTestHandler_CertificateInvalidFile_NoPanic(t *testing.T) {
	testCases := []struct {
		name     string
		fileData interface{}
	}{
		{
			name:     "nil file data",
			fileData: nil,
		},
		{
			name:     "integer file data",
			fileData: 12345,
		},
		{
			name:     "map file data",
			fileData: map[string]interface{}{"key": "value"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			profile := models.PerformanceProfile{
				Name: "test-profile",
				Metadata: map[string]interface{}{
					"ca_certificate": map[string]interface{}{
						"name": "meshery-panic-test.crt",
						"file": tc.fileData,
					},
				},
			}
			profileBytes, err := json.Marshal(profile)
			require.NoError(t, err)

			provider := &mockLoadTestProvider{profileData: profileBytes}
			h := &Handler{
				config: &models.HandlerConfig{},
				log:    newTestLogger(t),
			}

			req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/test-profile-id/run?cert=true&name=test-run&url=:", nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": "test-profile-id",
			})
			rec := httptest.NewRecorder()

			// Must not panic on invalid type assertion
			assert.NotPanics(t, func() {
				h.LoadTestHandler(rec, req, nil, nil, provider)
			})

			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Verify temp file is cleaned up
			tempDir := os.TempDir()
			matches, err := filepath.Glob(filepath.Join(tempDir, "meshery-panic-test-*.crt"))
			require.NoError(t, err)
			assert.Empty(t, matches, "expected temporary certificate file to be removed even on error")
		})
	}
}

// TestLoadTestHandler_CertificateUnsetParam_DoesNotCreateTempFile confirms that requests
// without the cert=true query parameter allocate zero temporary certificate files.
func TestLoadTestHandler_CertificateUnsetParam_DoesNotCreateTempFile(t *testing.T) {
	profile := models.PerformanceProfile{
		Name: "test-profile",
		Metadata: map[string]interface{}{
			"ca_certificate": map[string]interface{}{
				"name": "meshery-nocert-test.crt",
				"file": "some-cert",
			},
		},
	}
	profileBytes, err := json.Marshal(profile)
	require.NoError(t, err)

	provider := &mockLoadTestProvider{profileData: profileBytes}
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	// Request without cert=true param
	req := httptest.NewRequest(http.MethodPost, "/api/perf/profile/test-profile-id/run?name=test-run&url=:", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "test-profile-id",
	})
	rec := httptest.NewRecorder()

	assert.NotPanics(t, func() {
		h.LoadTestHandler(rec, req, nil, nil, provider)
	})

	tempDir := os.TempDir()
	matches, err := filepath.Glob(filepath.Join(tempDir, "meshery-nocert-test-*.crt"))
	require.NoError(t, err)
	assert.Empty(t, matches)
}
