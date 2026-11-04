package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

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

// TestServeUI_HTMLHasNoCacheAndVersionETag asserts Tier B: HTML documents are
// served with "public, no-cache" and a strong ETag carrying the build version.
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
	if got := rec.Header().Get("ETag"); got != `"v0.8.999"` {
		t.Errorf("ETag: expected %q, got %q", `"v0.8.999"`, got)
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

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("If-None-Match", `"v0.8.999"`)
	rec := httptest.NewRecorder()

	h.ServeUI(rec, req, "", dir)

	if got := rec.Code; got != http.StatusNotModified {
		t.Fatalf("expected 304 for matching If-None-Match, got %d", got)
	}
	if got := rec.Header().Get("ETag"); got != `"v0.8.999"` {
		t.Errorf("ETag on 304: expected %q, got %q", `"v0.8.999"`, got)
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
	if got := rec.Header().Get("ETag"); got == `"v0.8.999"` {
		t.Errorf("immutable asset must not carry the version ETag, got %q", got)
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
			wantETag: `"v0.8.999"`,
		},
		{
			name:     "html route with build",
			build:    "v0.8.999",
			reqURL:   "/settings",
			filePath: "settings.html",
			wantCC:   "public, no-cache",
			wantETag: `"v0.8.999"`,
		},
		{
			name:     "dynamic extension html with build",
			build:    "v0.8.999",
			reqURL:   "/extension/meshmap",
			filePath: "/extension/[...component].html",
			wantCC:   "public, no-cache",
			wantETag: `"v0.8.999"`,
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

			gotCC, gotETag := uiCacheHeaders(tt.reqURL, tt.filePath)
			if gotCC != tt.wantCC {
				t.Errorf("Cache-Control: expected %q, got %q", tt.wantCC, gotCC)
			}
			if gotETag != tt.wantETag {
				t.Errorf("ETag: expected %q, got %q", tt.wantETag, gotETag)
			}
		})
	}
}
