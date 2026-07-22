package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestHandlerWrappers_ForwardToHttputil is a smoke test for the four thin
// wrappers in utils.go that forward to server/models/httputil. The full
// behavioral surface (headers, MeshKit serialization, empty-object body)
// is covered by server/models/httputil/httputil_test.go; this test only
// guards the wrapper layer itself against accidental signature drift or
// a missed rename during a future refactor.
func TestHandlerWrappers_ForwardToHttputil(t *testing.T) {
	t.Run("writeJSONError", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONError(rec, "some error", http.StatusBadRequest)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("expected JSON Content-Type, got %q", ct)
		}

		var decoded map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("body did not parse as JSON: %v", err)
		}
		if decoded["error"] != "some error" {
			t.Errorf("expected wrapper to pass message through, got %q", decoded["error"])
		}
	})

	t.Run("writeMeshkitError", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeMeshkitError(rec, nil, http.StatusInternalServerError)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("expected JSON Content-Type, got %q", ct)
		}
	})

	t.Run("writeJSONMessage", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONMessage(rec, map[string]string{"message": "ok"}, http.StatusOK)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var decoded map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("body did not parse as JSON: %v", err)
		}
		if decoded["message"] != "ok" {
			t.Errorf("expected wrapper to pass payload through, got %q", decoded["message"])
		}
	})

	t.Run("writeJSONEmptyObject", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONEmptyObject(rec, http.StatusOK)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if body := rec.Body.String(); body != "{}" {
			t.Errorf("expected body %q, got %q", "{}", body)
		}
	})
}

func TestSafeOpenFile(t *testing.T) {
	// Redirect HOME so os.UserHomeDir() (and therefore the allowed directory)
	// is a hermetic temp dir instead of the developer's real home. t.TempDir()
	// is removed automatically, so the test leaves nothing behind. os.UserHomeDir
	// reads $HOME on unix/darwin; on other platforms it reads a different var and
	// this redirect is a no-op, so the hermetic assertions below are skipped.
	if _, ok := os.LookupEnv("HOME"); !ok {
		t.Skip("HOME is not the home-dir source on this platform")
	}
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	logsDir := filepath.Join(tmpHome, ".meshery", "logs", "registry")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// A real, existing log file: exercises the os.Root open on a path production
	// actually serves and lets us assert the returned content — a wrong return
	// (e.g. an empty or unresolved path) no longer passes silently.
	const wantContent = "registry log line\n"
	validFile := filepath.Join(logsDir, "registry-logs.log")
	if err := os.WriteFile(validFile, []byte(wantContent), 0o644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	// A legitimate file whose name contains "..": must be served, not rejected.
	dotsFile := filepath.Join(logsDir, "registry..2026-07-16.log")
	if err := os.WriteFile(dotsFile, []byte("rotated\n"), 0o644); err != nil {
		t.Fatalf("failed to write dotted log file: %v", err)
	}

	// A sensitive file inside ~/.meshery but OUTSIDE the logs dir: the point of
	// the fix is that these are no longer reachable through the endpoints.
	sensitiveFile := filepath.Join(tmpHome, ".meshery", "mesherydb.sql")
	if err := os.WriteFile(sensitiveFile, []byte("secret\n"), 0o644); err != nil {
		t.Fatalf("failed to write sensitive file: %v", err)
	}

	t.Run("valid existing log file is served with its content", func(t *testing.T) {
		f, err := SafeOpenFile(validFile)
		if err != nil {
			t.Fatalf("SafeOpenFile(%q) unexpected error: %v", validFile, err)
		}
		defer func() { _ = f.Close() }()
		got, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("reading returned file: %v", err)
		}
		if string(got) != wantContent {
			t.Errorf("content = %q, want %q", got, wantContent)
		}
	})

	t.Run("filename containing .. is allowed", func(t *testing.T) {
		f, err := SafeOpenFile(dotsFile)
		if err != nil {
			t.Fatalf("SafeOpenFile(%q) unexpected error: %v", dotsFile, err)
		}
		_ = f.Close()
	})

	// wantOutside distinguishes the out-of-allowed-dirs sentinel (handler → 400)
	// from any other open failure such as a missing file (handler → 500), so a
	// benign not-found is never reported to the client as an unsafe-path attack.
	errorCases := []struct {
		name        string
		inputPath   string
		wantOutside bool
	}{
		{"empty path", "", true},
		{"absolute path outside allowed scope", "/etc/passwd", true},
		{"sensitive file inside ~/.meshery but outside logs", sensitiveFile, true},
		{"lexical traversal escaping the logs dir", logsDir + "/../../../../etc/passwd", true},
		{"missing file inside allowed dir", filepath.Join(logsDir, "does-not-exist.log"), false},
	}
	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := SafeOpenFile(tc.inputPath)
			if err == nil {
				_ = f.Close()
				t.Fatalf("SafeOpenFile(%q) = nil error, want error", tc.inputPath)
			}
			if got := errors.Is(err, errFileOutsideAllowedDirs); got != tc.wantOutside {
				t.Errorf("errors.Is(err, errFileOutsideAllowedDirs) = %v, want %v (err = %v)", got, tc.wantOutside, err)
			}
		})
	}

	t.Run("symlink escaping the allowed dir is rejected", func(t *testing.T) {
		linkPath := filepath.Join(logsDir, "evil.log")
		if err := os.Symlink("/etc/passwd", linkPath); err != nil {
			t.Skipf("symlinks unsupported in this environment: %v", err)
		}
		t.Cleanup(func() { _ = os.Remove(linkPath) })
		f, err := SafeOpenFile(linkPath)
		if err == nil {
			_ = f.Close()
			t.Fatalf("SafeOpenFile(%q) followed a symlink out of the allowed dir", linkPath)
		}
	})
}
