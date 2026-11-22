package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLatestVersionForMesheryctl(t *testing.T) {
	t.Run("success returns trimmed release version and closes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		defer server.Close()

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		transport := installTrackingTransport(t)

		version, err := GetLatestVersionForMesheryctl()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if version != "v0.7.50" {
			t.Errorf("GetLatestVersionForMesheryctl() = %q, want %q", version, "v0.7.50")
		}
		assertBodyClosed(t, transport, "expected response body to be closed on 200 response")
	})

	t.Run("non-200 status code returns error and closes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("404 page not found"))
		}))
		defer server.Close()

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		transport := installTrackingTransport(t)

		version, err := GetLatestVersionForMesheryctl()
		if err == nil {
			t.Errorf("expected error on 404 status code, got version %q", version)
		}
		if version != "" {
			t.Errorf("expected empty version on error, got %q", version)
		}
		assertBodyClosed(t, transport, "expected response body to be closed on 404 response")
	})

	t.Run("unreachable endpoint returns error", func(t *testing.T) {
		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = "http://127.0.0.1:0"
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		version, err := GetLatestVersionForMesheryctl()
		if err == nil {
			t.Errorf("expected error on unreachable endpoint, got version %q", version)
		}
		if version != "" {
			t.Errorf("expected empty version on error, got %q", version)
		}
	})
}

func TestCheckMesheryctlClientVersion(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	t.Run("outdated version logs notification", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		defer server.Close()

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		transport := installTrackingTransport(t)

		// Should not panic or crash
		CheckMesheryctlClientVersion("v0.7.40")
		assertBodyClosed(t, transport, "expected response body to be closed")
	})

	t.Run("latest version logs up to date", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		defer server.Close()

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		transport := installTrackingTransport(t)

		// Should not panic or crash
		CheckMesheryctlClientVersion("v0.7.50")
		assertBodyClosed(t, transport, "expected response body to be closed")
	})

	t.Run("non-200 status logs warning instead of using error body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 internal server error"))
		}))
		defer server.Close()

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		defer func() { mesheryctlReleaseVersionURL = origURL }()

		transport := installTrackingTransport(t)

		// Should not panic or treat 500 error body as a new release
		CheckMesheryctlClientVersion("v0.7.50")
		assertBodyClosed(t, transport, "expected response body to be closed on non-200 response")
	})
}
