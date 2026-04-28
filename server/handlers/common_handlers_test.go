package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestViewHandlerServesAllowedMesheryLogFile(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	logDir := filepath.Join(homeDir, ".meshery", "logs", "registry")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("failed to create log directory: %v", err)
	}

	filePath := filepath.Join(logDir, "registry-logs.log")
	fileContents := "registry import failed"
	if err := os.WriteFile(filePath, []byte(fileContents), 0o600); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	h := newTestHandler(t, nil, "")
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file="+url.QueryEscape(filePath), nil)
	rec := httptest.NewRecorder()

	h.ViewHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}

	if body := rec.Body.String(); body != fileContents {
		t.Fatalf("expected response body %q, got %q", fileContents, body)
	}

	if contentType := resp.Header.Get("Content-Type"); contentType != "text/plain" {
		t.Fatalf("expected text/plain content type, got %q", contentType)
	}
}

func TestDownloadHandlerRejectsFileOutsideAllowedRoots(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	outsideDir := t.TempDir()
	filePath := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(filePath, []byte("secret"), 0o600); err != nil {
		t.Fatalf("failed to write outside file: %v", err)
	}

	h := newTestHandler(t, nil, "")
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileDownload?file="+url.QueryEscape(filePath), nil)
	rec := httptest.NewRecorder()

	h.DownloadHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", resp.StatusCode)
	}

	if !strings.Contains(rec.Body.String(), errFileAccessDenied.Error()) {
		t.Fatalf("expected forbidden response body to mention %q, got %q", errFileAccessDenied.Error(), rec.Body.String())
	}
}

func TestViewHandlerRejectsSymlinkEscapingAllowedRoots(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	logDir := filepath.Join(homeDir, ".meshery", "logs", "registry")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		t.Fatalf("failed to create log directory: %v", err)
	}

	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "secret.txt")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0o600); err != nil {
		t.Fatalf("failed to write outside file: %v", err)
	}

	linkPath := filepath.Join(logDir, "registry-link.log")
	if err := os.Symlink(outsideFile, linkPath); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	h := newTestHandler(t, nil, "")
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file="+url.QueryEscape(linkPath), nil)
	rec := httptest.NewRecorder()

	h.ViewHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 Forbidden, got %d", resp.StatusCode)
	}
}

func TestDownloadHandlerRejectsMissingFileQuery(t *testing.T) {
	h := newTestHandler(t, nil, "")
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileDownload", nil)
	rec := httptest.NewRecorder()

	h.DownloadHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.StatusCode)
	}
}
