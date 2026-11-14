package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// githubArchiveTarGz builds the shape of tarball GitHub serves from
// /archive/refs/tags/<tag>.tar.gz and that `git archive` produces: it begins
// with a PAX global header carrying the commit id, then a single top-level
// directory.
func githubArchiveTarGz(t *testing.T) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)

	write := func(h *tar.Header, body string) {
		t.Helper()
		if err := tw.WriteHeader(h); err != nil {
			t.Fatalf("write header %q: %v", h.Name, err)
		}
		if body != "" {
			if _, err := tw.Write([]byte(body)); err != nil {
				t.Fatalf("write body %q: %v", h.Name, err)
			}
		}
	}

	write(&tar.Header{
		Typeflag:   tar.TypeXGlobalHeader,
		Name:       "pax_global_header",
		PAXRecords: map[string]string{"comment": "0123456789abcdef0123456789abcdef01234567"},
	}, "")
	write(&tar.Header{Typeflag: tar.TypeDir, Name: "ext-v1.0.0/", Mode: 0755}, "")
	write(&tar.Header{Typeflag: tar.TypeReg, Name: "ext-v1.0.0/index.js", Mode: 0644, Size: int64(len("main"))}, "main")

	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

// manyFilesTarGz builds an archive of incompressible files, large enough that
// cutting it in half leaves several complete entries before the break.
func manyFilesTarGz(t *testing.T, files, size int) []byte {
	t.Helper()

	rng := rand.New(rand.NewSource(1))
	entries := make([]tarEntry, 0, files)
	for i := 0; i < files; i++ {
		body := make([]byte, size)
		_, _ = rng.Read(body)
		entries = append(entries, tarEntry{name: "f" + string(rune('a'+i%26)) + string(rune('a'+i/26)) + ".bin", body: string(body)})
	}
	return tarGzBytes(t, entries)
}

func dirEntries(t *testing.T, dir string) []string {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// GitHub archive tarballs start with a PAX global header. TarXZ treated it as
// an unknown entry type and refused the whole package.
func TestTarXZ_AcceptsGitHubArchiveWithPaxGlobalHeader(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "pkg")

	if err := TarXZ(bytes.NewReader(githubArchiveTarGz(t)), destination); err != nil {
		t.Fatalf("TarXZ rejected a GitHub-style archive: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(destination, "ext-v1.0.0", "index.js"))
	if err != nil {
		t.Fatalf("expected file was not extracted: %v", err)
	}
	if string(got) != "main" {
		t.Errorf("content = %q, want %q", got, "main")
	}
	if _, err := os.Stat(filepath.Join(destination, "pax_global_header")); err == nil {
		t.Error("the PAX global header must not be written out as a file")
	}
}

// `tar czf x.tgz -C dir .` is the most common way to build a package and begins
// with a "./" entry. validateExtractPath rejected it because the cleaned path
// equals the destination itself.
func TestTarXZ_AcceptsDotSlashEntries(t *testing.T) {
	archive := tarGzBytes(t, []tarEntry{
		{name: "./", dir: true},
		{name: "./a.txt", body: "a"},
		{name: "./sub/", dir: true},
		{name: "./sub/b.txt", body: "b"},
	})
	destination := filepath.Join(t.TempDir(), "pkg")

	if err := TarXZ(bytes.NewReader(archive), destination); err != nil {
		t.Fatalf("TarXZ rejected an archive built with `tar -C dir .`: %v", err)
	}
	for path, want := range map[string]string{"a.txt": "a", "sub/b.txt": "b"} {
		got, err := os.ReadFile(filepath.Join(destination, filepath.FromSlash(path)))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q (err %v), want %q", path, got, err, want)
		}
	}
}

// The path guard must stay strict: a bare ".." entry still escapes.
func TestTarXZ_StillRejectsParentDirectoryEntries(t *testing.T) {
	for _, name := range []string{"..", "../", "a/../../x.txt"} {
		archive := tarGzBytes(t, []tarEntry{{name: name, dir: strings.HasSuffix(name, "/") || name == "..", body: "x"}})
		err := TarXZ(bytes.NewReader(archive), filepath.Join(t.TempDir(), "pkg"))
		if err == nil {
			t.Errorf("TarXZ accepted the escaping entry %q", name)
		}
	}
}

// The callers skip the download whenever the destination exists. A failed
// extraction that leaves a half-populated directory therefore becomes a
// permanent, silently broken install.
func TestTarXZ_FailedExtractionLeavesNothingBehind(t *testing.T) {
	archive := manyFilesTarGz(t, 40, 4096)
	truncated := archive[:len(archive)/2]

	parent := t.TempDir()
	destination := filepath.Join(parent, "provider", "v1")

	if err := TarXZ(bytes.NewReader(truncated), destination); err == nil {
		t.Fatal("TarXZ reported success for a truncated archive")
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Errorf("a failed extraction left %s behind, so the next start would skip the download", destination)
	}
	// And no temporary directory is left next to it.
	if leftovers := dirEntries(t, filepath.Join(parent, "provider")); len(leftovers) != 0 {
		t.Errorf("failed extraction left temporary files behind: %v", leftovers)
	}
}

func TestTarXZ_CanRetryAfterAFailure(t *testing.T) {
	good := manyFilesTarGz(t, 10, 1024)
	destination := filepath.Join(t.TempDir(), "pkg")

	if err := TarXZ(bytes.NewReader(good[:len(good)/2]), destination); err == nil {
		t.Fatal("truncated archive unexpectedly succeeded")
	}
	if err := TarXZ(bytes.NewReader(good), destination); err != nil {
		t.Fatalf("a retry with the full archive failed: %v", err)
	}
	if got := len(dirEntries(t, destination)); got != 10 {
		t.Errorf("retry extracted %d files, want 10", got)
	}
}

func TestTarXZ_RejectsPackagesOverTheSizeCap(t *testing.T) {
	old := maxProviderPackageExtractedBytes
	maxProviderPackageExtractedBytes = 100
	t.Cleanup(func() { maxProviderPackageExtractedBytes = old })

	archive := tarGzBytes(t, []tarEntry{{name: "big.bin", body: strings.Repeat("x", 500)}})
	destination := filepath.Join(t.TempDir(), "pkg")

	err := TarXZ(bytes.NewReader(archive), destination)
	if err == nil {
		t.Fatal("TarXZ accepted a package larger than the cap")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("error = %q, want it to say the package exceeds the cap", err)
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Error("an oversized package left a partial install behind")
	}
}

// A package host that accepts the connection and never answers held the login
// request open indefinitely, because the download used http.Get.
func TestTarXZF_StalledServerTimesOut(t *testing.T) {
	old := providerPackageDownloadTimeout
	providerPackageDownloadTimeout = 300 * time.Millisecond
	t.Cleanup(func() { providerPackageDownloadTimeout = old })

	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-r.Context().Done():
		}
	}))
	t.Cleanup(func() { close(release); server.Close() })

	done := make(chan error, 1)
	go func() {
		done <- TarXZF(server.URL+"/pkg.tar.gz", filepath.Join(t.TempDir(), "pkg"), nil)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("a stalled download must fail")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("TarXZF did not give up on a stalled server; login would hang")
	}
}

func TestTarXZF_DownloadsAndInstalls(t *testing.T) {
	archive := githubArchiveTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(archive)
	}))
	t.Cleanup(server.Close)

	destination := filepath.Join(t.TempDir(), "provider", "v1")
	if err := TarXZF(server.URL+"/ext-v1.0.0.tar.gz", destination, nil); err != nil {
		t.Fatalf("TarXZF failed for a well-formed package: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destination, "ext-v1.0.0", "index.js")); err != nil {
		t.Errorf("package was not installed: %v", err)
	}
}

func TestTarXZF_ServerErrorLeavesNoInstall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	t.Cleanup(server.Close)

	destination := filepath.Join(t.TempDir(), "pkg")
	if err := TarXZF(server.URL+"/pkg.tar.gz", destination, nil); err == nil {
		t.Fatal("TarXZF reported success for an HTTP 500")
	}
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Error("an HTTP error left an install directory behind")
	}
}
