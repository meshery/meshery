package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshkit/logger"
)

func TestIsValidPath(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	tmpDir := t.TempDir()
	h := &Handler{
		log:         log,
		mesheryHome: tmpDir,
	}

	mesheryHome, _ := filepath.Abs(tmpDir)

	// Create some files for EvalSymlinks to work with
	testFile := filepath.Join(mesheryHome, "test.log")
	_ = os.WriteFile(testFile, []byte("test"), 0644)

	configDir := filepath.Join(mesheryHome, "config")
	_ = os.Mkdir(configDir, 0755)
	configFile := filepath.Join(configDir, "mesherydb.sql")
	_ = os.WriteFile(configFile, []byte("sensitive"), 0644)

	// Create a symlink that points outside
	outsideFile := filepath.Join(os.TempDir(), "outside.txt")
	_ = os.WriteFile(outsideFile, []byte("outside"), 0644)
	defer os.Remove(outsideFile)

	evilSymlink := filepath.Join(mesheryHome, "evil_link")
	_ = os.Symlink(outsideFile, evilSymlink)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Valid path in mesheryHome",
			path:     testFile,
			expected: true,
		},
		{
			name:     "Invalid path - sensitive config",
			path:     configFile,
			expected: false,
		},
		{
			name:     "Invalid path - etc passwd",
			path:     "/etc/passwd",
			expected: false,
		},
		{
			name:     "Invalid path - symlink traversal",
			path:     evilSymlink,
			expected: false,
		},
		{
			name:     "Invalid path - partial prefix match bypass",
			path:     mesheryHome + "-backup/secret.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.isValidPath(tt.path); got != tt.expected {
				t.Errorf("isValidPath(%v) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}

func TestViewHandler_Validation(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	tmpDir := t.TempDir()
	h := &Handler{
		log:         log,
		mesheryHome: tmpDir,
	}

	mesheryHome, _ := filepath.Abs(tmpDir)

	// Create a dummy log file for testing
	testLogDir := filepath.Join(mesheryHome, "test_logs")
	_ = os.MkdirAll(testLogDir, 0755)
	testLogFile := filepath.Join(testLogDir, "access.log")
	_ = os.WriteFile(testLogFile, []byte("test log content"), 0644)

	tests := []struct {
		name           string
		fileParam      string
		expectedStatus int
	}{
		{
			name:           "Allowed file",
			fileParam:      testLogFile,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Forbidden file - outside home",
			fileParam:      "/etc/passwd",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Forbidden file - sensitive config",
			fileParam:      filepath.Join(mesheryHome, "config", "mesherydb.sql"),
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file="+url.QueryEscape(tt.fileParam), nil)
			rec := httptest.NewRecorder()

			h.ViewHandler(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
