package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/spf13/viper"
)

// newUITestDir creates a temporary directory that mimics the Next.js static
// export layout (`ui/out`) just enough for ServeUI: an index.html document and
// a content-hashed immutable asset under /_next/static/chunks.
func newUITestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>meshery</title>"), 0o600); err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}

	chunkDir := filepath.Join(dir, "_next", "static", "chunks")
	if err := os.MkdirAll(chunkDir, 0o750); err != nil {
		t.Fatalf("failed to create chunk dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(chunkDir, "x.js"), []byte("console.log('x')"), 0o600); err != nil {
		t.Fatalf("failed to write chunk: %v", err)
	}

	return dir
}

// wantHTMLETag returns the ETag ServeUI should send for the HTML file at path
// when BUILD is build.
func wantHTMLETag(t *testing.T, path, build string) string {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat %s: %v", path, err)
	}
	return fmt.Sprintf(`"%s-%x-%x"`, build, info.ModTime().UnixNano(), info.Size())
}

// TestServeUI_HTMLHasNoCacheAndVersionETag asserts Tier B: HTML documents are
// served with "public, no-cache" and a strong ETag carrying the build version
// and the file's modification time and size.
func TestServeUI_HTMLHasNoCacheAndVersionETag(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusOK {
		t.Fatalf("expected 200 for /, got %d", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "public, no-cache" {
		t.Errorf("Cache-Control: expected %q, got %q", "public, no-cache", got)
	}
	if got, want := rec.Header().Get("ETag"), wantHTMLETag(t, filepath.Join(dir, "index.html"), "v0.8.999"); got != want {
		t.Errorf("ETag: expected %q, got %q", want, got)
	}
}

// TestServeUI_HTMLConditionalRequestReturns304 asserts the revalidation path:
// when the request's If-None-Match matches the release ETag, the origin returns
// 304 Not Modified (Go's http.ServeContent does this when the ETag header is set
// before the body is written - the headers-before-ServeFile ordering is what
// makes this work).
func TestServeUI_HTMLConditionalRequestReturns304(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	h := newTestHandler(t, map[string]models.Provider{}, "")
	etag := wantHTMLETag(t, filepath.Join(dir, "index.html"), "v0.8.999")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", etag)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusNotModified {
		t.Fatalf("expected 304 for matching If-None-Match, got %d", got)
	}
	if got := rec.Header().Get("ETag"); got != etag {
		t.Errorf("ETag on 304: expected %q, got %q", etag, got)
	}
}

// TestServeUI_RebuiltHTMLWithSameBuildReturnsNewDocument covers a rebuilt UI
// that keeps the same BUILD. Edge images are all built with the latest stable
// tag as their version, and a local `make` build uses that tag too, so BUILD
// stays the same while index.html changes. A browser or CDN revalidating with
// the previous ETag must get the new document, not 304 with the old one.
//
// The rebuilt document has the same length as the old one, as it does when
// only content hashes in chunk URLs change, so only its mtime tells them apart.
func TestServeUI_RebuiltHTMLWithSameBuildReturnsNewDocument(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	indexPath := filepath.Join(dir, "index.html")
	oldDoc := `<!doctype html><script src="/_next/static/chunks/aaaa.js"></script>`
	newDoc := `<!doctype html><script src="/_next/static/chunks/bbbb.js"></script>`
	if err := os.WriteFile(indexPath, []byte(oldDoc), 0o600); err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}
	built := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	if err := os.Chtimes(indexPath, built, built); err != nil {
		t.Fatalf("failed to set index.html mtime: %v", err)
	}

	h := newTestHandler(t, map[string]models.Provider{}, "")

	first := httptest.NewRecorder()
	h.ServeUI(first, httptest.NewRequest(http.MethodGet, "/", nil), "", dir)
	oldETag := first.Header().Get("ETag")
	if oldETag == "" {
		t.Fatal("expected an ETag on the first response")
	}

	// The next image ships a rebuilt index.html: new content, same length,
	// written later, same BUILD.
	if err := os.WriteFile(indexPath, []byte(newDoc), 0o600); err != nil {
		t.Fatalf("failed to rewrite index.html: %v", err)
	}
	rebuilt := built.Add(time.Hour)
	if err := os.Chtimes(indexPath, rebuilt, rebuilt); err != nil {
		t.Fatalf("failed to set index.html mtime: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", oldETag)
	rec := httptest.NewRecorder()
	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusOK {
		t.Fatalf("revalidating the rebuilt index.html with the old ETag %s: expected 200, got %d", oldETag, got)
	}
	if got := rec.Header().Get("ETag"); got == oldETag {
		t.Errorf("ETag did not change after index.html was rebuilt: %s", got)
	}
	if got := rec.Body.String(); got != newDoc {
		t.Errorf("expected the rebuilt document, got %q", got)
	}
}

// TestServeUI_ImmutableAssetCacheControl asserts Tier A: content-hashed assets
// under /_next/static/chunks get a one-year immutable Cache-Control and are NOT
// marked no-cache and carry no version ETag.
func TestServeUI_ImmutableAssetCacheControl(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/_next/static/chunks/x.js", nil)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusOK {
		t.Fatalf("expected 200 for asset, got %d", got)
	}
	if got, want := rec.Header().Get("Cache-Control"), "public, max-age=31536000, immutable"; got != want {
		t.Errorf("Cache-Control: expected %q, got %q", want, got)
	}
	if got := rec.Header().Get("Cache-Control"); got == "public, no-cache" {
		t.Errorf("immutable asset must not be no-cache, got %q", got)
	}
	if got := rec.Header().Get("ETag"); got != "" {
		t.Errorf("immutable asset must not carry an ETag, got %q", got)
	}
}

// TestServeUI_ProviderBasePathStripped asserts the Tier A prefix check runs on
// the POST-STRIP reqURL: a /provider/_next/... request (reqBasePath="/provider")
// must still be detected as an immutable asset.
func TestServeUI_ProviderBasePathStripped(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/provider/_next/static/chunks/x.js", nil)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "/provider", dir)

	if got := rec.Code; got != http.StatusOK {
		t.Fatalf("expected 200 for provider asset, got %d", got)
	}
	if got, want := rec.Header().Get("Cache-Control"), "public, max-age=31536000, immutable"; got != want {
		t.Errorf("Cache-Control: expected %q, got %q", want, got)
	}
}

// TestServeUI_MissingAssetDoesNotCache404 guards against caching error responses:
// a request for a content-hashed asset that does not exist on disk must 404
// WITHOUT any Cache-Control/ETag header. Otherwise a transient 404 (e.g. an old
// hash requested mid-deploy against a pod that doesn't have it) would be cached
// "immutable" for a year and permanently break the URL until a manual CDN purge.
func TestServeUI_MissingAssetDoesNotCache404(t *testing.T) {
	dir := newUITestDir(t)
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodGet, "/_next/static/chunks/missing.js", nil)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusNotFound {
		t.Fatalf("expected 404 for missing asset, got %d", got)
	}
	if got := rec.Header().Get("Cache-Control"); got != "" {
		t.Errorf("404 must not carry a Cache-Control header, got %q", got)
	}
	if got := rec.Header().Get("ETag"); got != "" {
		t.Errorf("404 must not carry an ETag header, got %q", got)
	}
}

// TestUICacheHeaders unit-tests the pure decision helper directly across the
// full prefix/HTML/version matrix, including the dev-mode BUILD sentinels that
// must suppress the ETag.
func TestUICacheHeaders(t *testing.T) {
	modTime := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	const size = 1234
	htmlETag := fmt.Sprintf(`"v0.8.999-%x-%x"`, modTime.UnixNano(), size)

	tests := []struct {
		name     string
		build    string
		reqURL   string
		filePath string
		wantCC   string
		wantETag string
	}{
		{
			name:     "chunks asset",
			build:    "v0.8.999",
			reqURL:   "/_next/static/chunks/x.js",
			filePath: "/_next/static/chunks/x.js",
			wantCC:   "public, max-age=31536000, immutable",
			wantETag: "",
		},
		{
			name:     "css asset",
			build:    "v0.8.999",
			reqURL:   "/_next/static/css/app.css",
			filePath: "/_next/static/css/app.css",
			wantCC:   "public, max-age=31536000, immutable",
			wantETag: "",
		},
		{
			name:     "media asset",
			build:    "v0.8.999",
			reqURL:   "/_next/static/media/logo.svg",
			filePath: "/_next/static/media/logo.svg",
			wantCC:   "public, max-age=31536000, immutable",
			wantETag: "",
		},
		{
			name:     "html index with build",
			build:    "v0.8.999",
			reqURL:   "/",
			filePath: "index.html",
			wantCC:   "public, no-cache",
			wantETag: htmlETag,
		},
		{
			name:     "html route with build",
			build:    "v0.8.999",
			reqURL:   "/settings",
			filePath: "settings.html",
			wantCC:   "public, no-cache",
			wantETag: htmlETag,
		},
		{
			name:     "dynamic extension html with build",
			build:    "v0.8.999",
			reqURL:   "/extension/meshmap",
			filePath: "/extension/[...component].html",
			wantCC:   "public, no-cache",
			wantETag: htmlETag,
		},
		{
			name:     "html with empty build suppresses etag",
			build:    "",
			reqURL:   "/",
			filePath: "index.html",
			wantCC:   "public, no-cache",
			wantETag: "",
		},
		{
			name:     "html with Not set build suppresses etag",
			build:    "Not set",
			reqURL:   "/",
			filePath: "index.html",
			wantCC:   "public, no-cache",
			wantETag: "",
		},
		{
			// main.go's default version when no -ldflags set it.
			name:     "html with Not Set build suppresses etag",
			build:    "Not Set",
			reqURL:   "/",
			filePath: "index.html",
			wantCC:   "public, no-cache",
			wantETag: "",
		},
		{
			name:     "non-immutable _next path is not cached forever",
			build:    "v0.8.999",
			reqURL:   "/_next/static/x.js",
			filePath: "/_next/static/x.js",
			wantCC:   "",
			wantETag: "",
		},
		{
			name:     "other static file gets defaults",
			build:    "v0.8.999",
			reqURL:   "/favicon.ico",
			filePath: "favicon.ico",
			wantCC:   "",
			wantETag: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("BUILD", tt.build)
			t.Cleanup(func() { viper.Set("BUILD", "") })

			gotCC, gotETag := uiCacheHeaders(tt.reqURL, tt.filePath, modTime, size)
			if gotCC != tt.wantCC {
				t.Errorf("Cache-Control: expected %q, got %q", tt.wantCC, gotCC)
			}
			if gotETag != tt.wantETag {
				t.Errorf("ETag: expected %q, got %q", tt.wantETag, gotETag)
			}
		})
	}
}

// TestUICacheHeaders_ETagTracksFileIdentity pins that, with BUILD unchanged, the
// HTML ETag changes when the file's mtime or size changes and stays the same
// for the same file.
func TestUICacheHeaders_ETagTracksFileIdentity(t *testing.T) {
	viper.Set("BUILD", "v0.8.999")
	t.Cleanup(func() { viper.Set("BUILD", "") })

	built := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	_, base := uiCacheHeaders("/", "index.html", built, 1234)

	if _, again := uiCacheHeaders("/", "index.html", built, 1234); again != base {
		t.Errorf("same file: expected ETag %q, got %q", base, again)
	}
	if _, rebuilt := uiCacheHeaders("/", "index.html", built.Add(time.Second), 1234); rebuilt == base {
		t.Errorf("newer mtime, same size: ETag stayed %q", rebuilt)
	}
	if _, resized := uiCacheHeaders("/", "index.html", built, 1235); resized == base {
		t.Errorf("same mtime, new size: ETag stayed %q", resized)
	}
}
