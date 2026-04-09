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
	mesheryHome, err = filepath.Abs(mesheryHome)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Create some files for EvalSymlinks to work with
	testFile := filepath.Join(mesheryHome, "test.log")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	configDir := filepath.Join(mesheryHome, "config")
	if err := os.Mkdir(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}
	configFile := filepath.Join(configDir, "mesherydb.sql")
	if err := os.WriteFile(configFile, []byte("sensitive"), 0644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	// Create a symlink that points outside
	outsideFile := filepath.Join(os.TempDir(), "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("outside"), 0644); err != nil {
		t.Fatalf("failed to create outside file: %v", err)
	}
	defer func() {
		if err := os.Remove(outsideFile); err != nil {
			t.Logf("failed to remove temporary file: %v", err)
		}
	}()

	evilSymlink := filepath.Join(mesheryHome, "evil_link")
	if err := os.Symlink(outsideFile, evilSymlink); err != nil {
		t.Fatalf("failed to create symlink: %v", err)
	}

	// Create a directory that could cause a partial prefix bypass
	// We use the same parent directory as mesheryHome
	parentDir := filepath.Dir(mesheryHome)
	backupDir := filepath.Join(parentDir, filepath.Base(mesheryHome)+"-backup")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		t.Fatalf("failed to create backup dir: %v", err)
	}
	backupFile := filepath.Join(backupDir, "secret.txt")
	if err := os.WriteFile(backupFile, []byte("secret"), 0644); err != nil {
		t.Fatalf("failed to create backup file: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(backupDir); err != nil {
			t.Logf("failed to remove backup dir: %v", err)
		}
	}()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Given a valid path in mesheryHome, When isValidPath is called, Then it should return true",
			path:     testFile,
			expected: true,
		},
		{
			name:     "Given an invalid path to a sensitive config file, When isValidPath is called, Then it should return false",
			path:     configFile,
			expected: false,
		},
		{
			name:     "Given an invalid path to an external file like etc/passwd, When isValidPath is called, Then it should return false",
			path:     "/etc/passwd",
			expected: false,
		},
		{
			name:     "Given an invalid path with symlink traversal, When isValidPath is called, Then it should return false",
			path:     evilSymlink,
			expected: false,
		},
		{
			name:     "Given an invalid path with a partial prefix match bypass, When isValidPath is called, Then it should return false",
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
	mesheryHome, err = filepath.Abs(mesheryHome)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Create a dummy log file for testing
	testLogDir := filepath.Join(mesheryHome, "test_logs")
	if err := os.MkdirAll(testLogDir, 0755); err != nil {
		t.Fatalf("failed to create log dir: %v", err)
	}
	testLogFile := filepath.Join(testLogDir, "access.log")
	if err := os.WriteFile(testLogFile, []byte("test log content"), 0644); err != nil {
		t.Fatalf("failed to create log file: %v", err)
	}

	tests := []struct {
		name           string
		fileParam      string
		expectedStatus int
	}{
		{
			name:           "Given an allowed file, When ViewHandler is called, Then it should return StatusOK",
			fileParam:      testLogFile,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Given a forbidden file outside mesheryHome, When ViewHandler is called, Then it should return StatusForbidden",
			fileParam:      "/etc/passwd",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Given a forbidden sensitive config file, When ViewHandler is called, Then it should return StatusForbidden",
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
