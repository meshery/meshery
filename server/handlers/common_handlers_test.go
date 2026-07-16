package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// mustCanonical returns the fully symlink-resolved absolute path, matching what
// resolveFileWithinDir compares against. Tests build their expectations through
// it so they pass on macOS (where /var and t.TempDir resolve through symlinks).
func mustCanonical(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("EvalSymlinks(%q): %v", path, err)
	}
	return resolved
}

// TestResolveFileWithinDir exercises the single validation boundary shared by
// ViewHandler and DownloadHandler. It covers the concrete bypasses called out in
// issues #18442, #18375 and #14193: "../" traversal, absolute paths outside the
// base, symlink traversal, and the sibling-prefix ("<base>-backup") bypass.
func TestResolveFileWithinDir(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, "logs")
	if err := os.MkdirAll(filepath.Join(base, "registry"), 0o755); err != nil {
		t.Fatalf("mkdir base: %v", err)
	}

	// A legitimate log file at the base root and a nested one.
	rootLog := filepath.Join(base, "registry-logs.log")
	if err := os.WriteFile(rootLog, []byte("log line\n"), 0o644); err != nil {
		t.Fatalf("write rootLog: %v", err)
	}
	nestedLog := filepath.Join(base, "registry", "model-generation.log")
	if err := os.WriteFile(nestedLog, []byte("nested\n"), 0o644); err != nil {
		t.Fatalf("write nestedLog: %v", err)
	}

	// A sensitive file OUTSIDE the base, e.g. ~/.meshery/config/mesherydb.sql.
	outsideDir := filepath.Join(root, "config")
	if err := os.MkdirAll(outsideDir, 0o755); err != nil {
		t.Fatalf("mkdir outsideDir: %v", err)
	}
	outsideFile := filepath.Join(outsideDir, "mesherydb.sql")
	if err := os.WriteFile(outsideFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("write outsideFile: %v", err)
	}

	// A sibling directory sharing the base name as a prefix ("<base>-backup").
	siblingDir := base + "-backup"
	if err := os.MkdirAll(siblingDir, 0o755); err != nil {
		t.Fatalf("mkdir siblingDir: %v", err)
	}
	siblingFile := filepath.Join(siblingDir, "secret.txt")
	if err := os.WriteFile(siblingFile, []byte("secret"), 0o644); err != nil {
		t.Fatalf("write siblingFile: %v", err)
	}

	// A relative "../" escape whose target actually exists, so the rejection is
	// exercised by the containment check rather than by non-existence.
	escapeExisting := filepath.Join(root, "escape.txt")
	if err := os.WriteFile(escapeExisting, []byte("escape"), 0o644); err != nil {
		t.Fatalf("write escapeExisting: %v", err)
	}

	// A symlink INSIDE the base that points OUTSIDE it (symlink traversal).
	symlinkOut := filepath.Join(base, "evil_link")
	if err := os.Symlink(outsideFile, symlinkOut); err != nil {
		t.Fatalf("symlink out: %v", err)
	}
	// A symlink INSIDE the base that points to another file INSIDE it (allowed).
	symlinkIn := filepath.Join(base, "good_link")
	if err := os.Symlink(rootLog, symlinkIn); err != nil {
		t.Fatalf("symlink in: %v", err)
	}

	tests := []struct {
		name         string
		requested    string
		wantOK       bool
		wantResolved string // only checked when wantOK
		reasonSubstr string // only checked when !wantOK && non-empty
	}{
		{name: "legit file at base root", requested: rootLog, wantOK: true, wantResolved: mustCanonical(t, rootLog)},
		{name: "legit nested file", requested: nestedLog, wantOK: true, wantResolved: mustCanonical(t, nestedLog)},
		{name: "relative path within base", requested: "registry/model-generation.log", wantOK: true, wantResolved: mustCanonical(t, nestedLog)},
		{name: "symlink within base is allowed", requested: symlinkIn, wantOK: true, wantResolved: mustCanonical(t, rootLog)},

		{name: "absolute path outside base", requested: outsideFile, wantOK: false, reasonSubstr: "escapes"},
		{name: "relative dotdot traversal", requested: "../escape.txt", wantOK: false, reasonSubstr: "escapes"},
		{name: "deep dotdot traversal to etc", requested: "../../../../../../etc/hosts", wantOK: false},
		{name: "symlink traversal outside base", requested: symlinkOut, wantOK: false, reasonSubstr: "escapes"},
		{name: "sibling prefix bypass", requested: siblingFile, wantOK: false, reasonSubstr: "escapes"},
		{name: "base directory itself", requested: base, wantOK: false, reasonSubstr: "escapes"},
		{name: "empty path", requested: "", wantOK: false, reasonSubstr: "empty"},
		{name: "nul byte injection", requested: "registry/model.log\x00.png", wantOK: false, reasonSubstr: "NUL"},
		{name: "non-existent file within base fails closed", requested: filepath.Join(base, "does-not-exist.log"), wantOK: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resolved, ok, reason := resolveFileWithinDir(tc.requested, base)
			if ok != tc.wantOK {
				t.Fatalf("resolveFileWithinDir(%q) ok=%v reason=%q, want ok=%v", tc.requested, ok, reason, tc.wantOK)
			}
			if tc.wantOK {
				if resolved != tc.wantResolved {
					t.Fatalf("resolved=%q, want %q", resolved, tc.wantResolved)
				}
				return
			}
			if resolved != "" {
				t.Fatalf("rejected request returned a non-empty path %q", resolved)
			}
			if tc.reasonSubstr != "" && !strings.Contains(reason, tc.reasonSubstr) {
				t.Fatalf("reason=%q, want it to contain %q", reason, tc.reasonSubstr)
			}
		})
	}
}

// setupHandlerHome points os.UserHomeDir at a temp home, creates
// ~/.meshery/logs, and returns the Handler under test plus the logs dir.
func setupHandlerHome(t *testing.T) (*Handler, string) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", home)
	}
	logsDir := filepath.Join(home, ".meshery", "logs")
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		t.Fatalf("mkdir logs: %v", err)
	}
	return &Handler{log: newTestLogger(t)}, logsDir
}

func newFileRequest(t *testing.T, endpoint, filePath string) *http.Request {
	t.Helper()
	return httptest.NewRequest(http.MethodGet, endpoint+"?file="+url.QueryEscape(filePath), nil)
}

func TestViewHandler_ServesFileWithinMesheryLogs(t *testing.T) {
	h, logsDir := setupHandlerHome(t)
	logFile := filepath.Join(logsDir, "registry", "registry-logs.log")
	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	const want = "registration completed\n"
	if err := os.WriteFile(logFile, []byte(want), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	rec := httptest.NewRecorder()
	h.ViewHandler(rec, newFileRequest(t, "/api/system/fileView", logFile))

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	if got := rec.Body.String(); got != want {
		t.Fatalf("body=%q, want %q", got, want)
	}
}

func TestDownloadHandler_ServesFileWithinMesheryLogs(t *testing.T) {
	h, logsDir := setupHandlerHome(t)
	logFile := filepath.Join(logsDir, "registry-logs.log")
	if err := os.WriteFile(logFile, []byte("payload"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	rec := httptest.NewRecorder()
	h.DownloadHandler(rec, newFileRequest(t, "/api/system/fileDownload", logFile))

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	if cd := rec.Header().Get("Content-Disposition"); !strings.Contains(cd, `filename="registry-logs.log"`) {
		t.Fatalf("Content-Disposition=%q, want it to name registry-logs.log", cd)
	}
}

// TestFileHandlers_RejectTraversal reproduces the PoCs from #18442/#18375/#14193
// against both handlers and asserts a 403 with no file content leaked.
func TestFileHandlers_RejectTraversal(t *testing.T) {
	handlers := []struct {
		name     string
		endpoint string
		call     func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request)
	}{
		{"view", "/api/system/fileView", func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request) { h.ViewHandler(rec, req) }},
		{"download", "/api/system/fileDownload", func(h *Handler, rec *httptest.ResponseRecorder, req *http.Request) { h.DownloadHandler(rec, req) }},
	}

	for _, hc := range handlers {
		t.Run(hc.name, func(t *testing.T) {
			h, logsDir := setupHandlerHome(t)

			// A secret sitting in ~/.meshery/config, a sibling of ~/.meshery/logs.
			configDir := filepath.Join(filepath.Dir(logsDir), "config")
			if err := os.MkdirAll(configDir, 0o755); err != nil {
				t.Fatalf("mkdir config: %v", err)
			}
			secret := filepath.Join(configDir, "mesherydb.sql")
			if err := os.WriteFile(secret, []byte("TOP SECRET"), 0o644); err != nil {
				t.Fatalf("write secret: %v", err)
			}

			// A symlink planted inside the logs dir pointing at the secret.
			planted := filepath.Join(logsDir, "evil_link")
			if err := os.Symlink(secret, planted); err != nil {
				t.Fatalf("symlink: %v", err)
			}

			cases := []struct {
				name string
				file string
			}{
				{"absolute etc passwd", "/etc/passwd"},
				{"relative traversal", "../../../../../../etc/passwd"},
				{"meshery config secret", secret},
				{"symlink escape", planted},
				{"sibling prefix", logsDir + "-backup/secret.txt"},
			}
			for _, c := range cases {
				t.Run(c.name, func(t *testing.T) {
					rec := httptest.NewRecorder()
					hc.call(h, rec, newFileRequest(t, hc.endpoint, c.file))

					if rec.Code != http.StatusForbidden {
						t.Fatalf("status=%d, want 403", rec.Code)
					}
					if strings.Contains(rec.Body.String(), "TOP SECRET") {
						t.Fatalf("handler leaked secret content: %q", rec.Body.String())
					}
				})
			}
		})
	}
}

// TestViewHandler_InvalidEscaping ensures a malformed percent-encoding is a 400,
// distinct from the 403 for a well-formed but forbidden path. The raw value
// "%25zz" survives net/url's first (query) decode as "%zz", which then fails the
// handler's url.QueryUnescape and exercises the ErrInvalidFileRequest branch.
func TestViewHandler_InvalidEscaping(t *testing.T) {
	h, _ := setupHandlerHome(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/fileView?file=%25zz", nil)
	rec := httptest.NewRecorder()
	h.ViewHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d, want 400 (body=%q)", rec.Code, rec.Body.String())
	}
	// Drain the body so the recorder is fully exercised.
	_, _ = io.Copy(io.Discard, rec.Body)
}
