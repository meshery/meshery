package utils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
)

func TestGetLatestVersionForMesheryctl(t *testing.T) {
	t.Run("success returns trimmed release version and closes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		t.Cleanup(server.Close)

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

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
		t.Cleanup(server.Close)

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

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
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

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
	var logOutput bytes.Buffer
	testLogger, err := logger.New("mesheryctl", logger.Options{
		Format:   logger.TerminalLogFormat,
		LogLevel: int(logrus.InfoLevel),
		Output:   &logOutput,
	})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	prevLog := Log
	Log = testLogger
	t.Cleanup(func() { Log = prevLog })

	t.Run("outdated version logs notification", func(t *testing.T) {
		logOutput.Reset()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		t.Cleanup(server.Close)

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

		transport := installTrackingTransport(t)

		CheckMesheryctlClientVersion("v0.7.40")
		assertBodyClosed(t, transport, "expected response body to be closed")

		out := logOutput.String()
		if !strings.Contains(out, "A new release of mesheryctl is available: v0.7.40 → v0.7.50") {
			t.Errorf("expected upgrade notice in logs, got %q", out)
		}
	})

	t.Run("latest version logs up to date", func(t *testing.T) {
		logOutput.Reset()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("v0.7.50\n"))
		}))
		t.Cleanup(server.Close)

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

		transport := installTrackingTransport(t)

		CheckMesheryctlClientVersion("v0.7.50")
		assertBodyClosed(t, transport, "expected response body to be closed")

		out := logOutput.String()
		if !strings.Contains(out, "is the latest release.") {
			t.Errorf("expected up-to-date message in logs, got %q", out)
		}
	})

	t.Run("non-200 status logs warning instead of using error body", func(t *testing.T) {
		logOutput.Reset()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("500 internal server error"))
		}))
		t.Cleanup(server.Close)

		origURL := mesheryctlReleaseVersionURL
		mesheryctlReleaseVersionURL = server.URL
		t.Cleanup(func() { mesheryctlReleaseVersionURL = origURL })

		transport := installTrackingTransport(t)

		CheckMesheryctlClientVersion("v0.7.50")
		assertBodyClosed(t, transport, "expected response body to be closed on non-200 response")

		out := logOutput.String()
		if !strings.Contains(out, "unable to check for latest version of mesheryctl") {
			t.Errorf("expected warning in logs for non-200 status, got %q", out)
		}
		if strings.Contains(out, "500 internal server error") {
			t.Errorf("did not expect error body to appear as a release version in logs, got %q", out)
		}
	})
}
