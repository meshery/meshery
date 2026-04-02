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

	mesheryHome, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		mesheryHome = tmpDir
	}
	mesheryHome, _ = filepath.Abs(mesheryHome)

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
	defer func() {
		if err := os.Remove(outsideFile); err != nil {
			t.Logf("failed to remove temporary file: %v", err)
		}
	}()

	evilSymlink := filepath.Join(mesheryHome, "evil_link")
	_ = os.Symlink(outsideFile, evilSymlink)

	// Create a directory that could cause a partial prefix bypass
	// We use the same parent directory as mesheryHome
	parentDir := filepath.Dir(mesheryHome)
	backupDir := filepath.Join(parentDir, filepath.Base(mesheryHome)+"-backup")
	_ = os.MkdirAll(backupDir, 0755)
	backupFile := filepath.Join(backupDir, "secret.txt")
	_ = os.WriteFile(backupFile, []byte("secret"), 0644)
	defer func() {
		_ = os.RemoveAll(backupDir)
	}()

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
			path:     backupFile,
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

	mesheryHome, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		mesheryHome = tmpDir
	}
	mesheryHome, _ = filepath.Abs(mesheryHome)

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
