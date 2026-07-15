package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/utils"
)

func TestViewHandler_SafePath(t *testing.T) {
	// Create a temporary file inside the allowed directory (~/.meshery)
	homeDir := utils.GetHome()
	mesheryDir := filepath.Join(homeDir, ".meshery")
	
	// Ensure directory exists for testing
	if err := os.MkdirAll(mesheryDir, 0o755); err != nil {
		t.Fatalf("failed to create .meshery directory: %v", err)
	}

	tmpFile, err := os.CreateTemp(mesheryDir, "safe-file-*.txt")
	if err != nil {
		t.Fatalf("failed to create safe temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	expectedContent := "safe logs or files content"
	if _, err := tmpFile.WriteString(expectedContent); err != nil {
		t.Fatalf("failed to write to safe temp file: %v", err)
	}
	tmpFile.Close()

	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Construct request with the filepath of the safe file
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file="+url.QueryEscape(tmpFile.Name()), nil)
	rec := httptest.NewRecorder()

	h.ViewHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(body))
	}
}

func TestViewHandler_PathTraversalBlocked(t *testing.T) {
	// Create a temporary file outside the allowed directory
	tmpFile, err := os.CreateTemp("", "unsafe-file-*.txt")
	if err != nil {
		t.Fatalf("failed to create unsafe temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Request access to the file located outside of ~/.meshery
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file="+url.QueryEscape(tmpFile.Name()), nil)
	rec := httptest.NewRecorder()

	h.ViewHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	// Should be blocked and return 400 Bad Request
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for path traversal attempt, got %d", resp.StatusCode)
	}
}

func TestDownloadHandler_SafePath(t *testing.T) {
	homeDir := utils.GetHome()
	mesheryDir := filepath.Join(homeDir, ".meshery")

	if err := os.MkdirAll(mesheryDir, 0o755); err != nil {
		t.Fatalf("failed to create .meshery directory: %v", err)
	}

	tmpFile, err := os.CreateTemp(mesheryDir, "safe-download-*.txt")
	if err != nil {
		t.Fatalf("failed to create safe download file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	expectedContent := "safe downloadable logs content"
	if _, err := tmpFile.WriteString(expectedContent); err != nil {
		t.Fatalf("failed to write to safe download file: %v", err)
	}
	tmpFile.Close()

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/system/fileDownload?file="+url.QueryEscape(tmpFile.Name()), nil)
	rec := httptest.NewRecorder()

	h.DownloadHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if string(body) != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, string(body))
	}
}

func TestDownloadHandler_PathTraversalBlocked(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "unsafe-download-*.txt")
	if err != nil {
		t.Fatalf("failed to create unsafe download file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/api/system/fileDownload?file="+url.QueryEscape(tmpFile.Name()), nil)
	rec := httptest.NewRecorder()

	h.DownloadHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for path traversal attempt, got %d", resp.StatusCode)
	}
}
