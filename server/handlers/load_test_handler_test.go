package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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

// certObservingRecorder wraps httptest.ResponseRecorder to inspect the temp directory
// during response write, before LoadTestHandler's deferred cleanUpFiles loop deletes the file.
type certObservingRecorder struct {
	*httptest.ResponseRecorder
	tempDir         string
	certObserved    bool
	observedContent string
	observedPath    string
}

func (r *certObservingRecorder) WriteHeader(code int) {
	r.inspectDir()
	r.ResponseRecorder.WriteHeader(code)
}

func (r *certObservingRecorder) Write(b []byte) (int, error) {
	r.inspectDir()
	return r.ResponseRecorder.Write(b)
}

func (r *certObservingRecorder) inspectDir() {
	if r.certObserved {
		return
	}
	entries, err := os.ReadDir(r.tempDir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "ca-") && strings.HasSuffix(entry.Name(), ".crt") {
			fullPath := filepath.Join(r.tempDir, entry.Name())
			data, err := os.ReadFile(fullPath)
			if err == nil {
				r.certObserved = true
				r.observedContent = string(data)
				r.observedPath = fullPath
			}
		}
	}
}

// TestLoadTestHandlerSSLCertificateHappyPath verifies that a valid ca_certificate in performance
// profile metadata is written to disk with the expected contents, flushed/closed, and cleaned up
// on handler completion.
func TestLoadTestHandlerSSLCertificateHappyPath(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)

	expectedCert := "-----BEGIN CERTIFICATE-----\nMIIB_TEST_CERT_DATA\n-----END CERTIFICATE-----"
	profile := &models.PerformanceProfile{
		Metadata: map[string]interface{}{
			"ca_certificate": map[string]interface{}{
				"name": "ca.crt",
				"file": expectedCert,
			},
		},
	}

	h := &Handler{
		log: newTestLogger(t),
	}

	provider := &fakeLoadTestProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		profile:              profile,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/perf/profile/run?cert=true&name=test-test&url=example.com", nil)
	rec := &certObservingRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		tempDir:          tempDir,
	}

	h.LoadTestHandler(rec, req, nil, nil, provider)

	// 1. Verify that the certificate file was written to disk before handler exit
	if !rec.certObserved {
		t.Fatal("expected certificate file matching ca-*.crt to be written to disk in temp directory, but none was observed")
	}

	// 2. Verify that the file content matches the certificate data
	if rec.observedContent != expectedCert {
		t.Fatalf("expected written certificate content %q, got %q", expectedCert, rec.observedContent)
	}

	// 3. Verify that cleanUpFiles deferred cleanup removed the temp file upon handler exit
	if _, err := os.Stat(rec.observedPath); !os.IsNotExist(err) {
		t.Fatalf("expected certificate temp file %s to be cleaned up on handler exit, but it still exists", rec.observedPath)
	}
}

// TestAssignCertificatePath verifies that assignCertificatePath properly wires the certificate
// file path into the CACert field of LoadTestOptions.
func TestAssignCertificatePath(t *testing.T) {
	opts := &models.LoadTestOptions{}
	testPath := "/tmp/ca-12345.crt"
	assignCertificatePath("ca_certificate", testPath, opts)
	if opts.CACert != testPath {
		t.Fatalf("expected CACert to be %q, got %q", testPath, opts.CACert)
	}
}
