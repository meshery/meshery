package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"math"
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
	old := providerPackageStallTimeout
	providerPackageStallTimeout = 300 * time.Millisecond
	t.Cleanup(func() { providerPackageStallTimeout = old })

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
		if !strings.Contains(err.Error(), "stalled") {
			t.Errorf("error = %q, want it to say the download stalled", err)
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
		w.WriteHeader(http.StatusInternalServerError)
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

// trickle writes body to w in small pieces with a pause between them, flushing
// each, so a test can model a host that is slow but steady.
func trickle(w http.ResponseWriter, r *http.Request, body []byte, pieces int, pause time.Duration) {
	flusher, _ := w.(http.Flusher)
	size := (len(body) + pieces - 1) / pieces
	for start := 0; start < len(body); start += size {
		end := start + size
		if end > len(body) {
			end = len(body)
		}
		_, _ = w.Write(body[start:end])
		if flusher != nil {
			flusher.Flush()
		}
		select {
		case <-time.After(pause):
		case <-r.Context().Done():
			return
		}
	}
}

// A slow host that keeps sending must not be cut off just because the whole
// transfer takes longer than the stall timeout: the timeout measures silence,
// not total time, so it does not need to scale with the package size.
func TestTarXZF_SlowButSteadyDownloadSucceeds(t *testing.T) {
	oldStall := providerPackageStallTimeout
	providerPackageStallTimeout = 400 * time.Millisecond
	t.Cleanup(func() { providerPackageStallTimeout = oldStall })

	archive := githubArchiveTarGz(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 12 pieces, 100ms apart: about 1.2s in total, three times the stall timeout.
		trickle(w, r, archive, 12, 100*time.Millisecond)
	}))
	t.Cleanup(server.Close)

	start := time.Now()
	destination := filepath.Join(t.TempDir(), "pkg")
	if err := TarXZF(server.URL+"/ext.tar.gz", destination, nil); err != nil {
		t.Fatalf("a slow but steady download was cut off after %s: %v", time.Since(start), err)
	}
	if time.Since(start) <= providerPackageStallTimeout {
		t.Fatalf("the transfer finished in %s, so this test did not exercise a transfer longer than the stall timeout", time.Since(start))
	}
	if _, err := os.Stat(filepath.Join(destination, "ext-v1.0.0", "index.js")); err != nil {
		t.Errorf("package was not installed: %v", err)
	}
}

// The overall bound still holds for a host that keeps trickling data forever,
// which the stall timeout alone would never catch.
func TestTarXZF_OverallTimeoutBoundsAnEndlessTrickle(t *testing.T) {
	oldStall, oldTotal := providerPackageStallTimeout, providerPackageDownloadTimeout
	providerPackageStallTimeout = 300 * time.Millisecond
	providerPackageDownloadTimeout = 800 * time.Millisecond
	t.Cleanup(func() {
		providerPackageStallTimeout, providerPackageDownloadTimeout = oldStall, oldTotal
	})

	// A valid gzip and tar prefix for a file far larger than will ever arrive
	// (but under the extraction cap), then one byte every 50ms. The extractor
	// keeps reading, so only the overall bound can stop it. Raw bytes that are
	// not gzip would be rejected at once and prove nothing about the timeout.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flusher, _ := w.(http.Flusher)
		gz := gzip.NewWriter(w)
		tw := tar.NewWriter(gz)
		if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: "huge.bin", Mode: 0644, Size: 100 << 20}); err != nil {
			return
		}
		for {
			if _, err := tw.Write([]byte("x")); err != nil {
				return
			}
			if err := gz.Flush(); err != nil {
				return
			}
			if flusher != nil {
				flusher.Flush()
			}
			select {
			case <-time.After(50 * time.Millisecond):
			case <-r.Context().Done():
				return
			}
		}
	}))
	t.Cleanup(server.Close)

	done := make(chan error, 1)
	go func() {
		done <- TarXZF(server.URL+"/pkg.tar.gz", filepath.Join(t.TempDir(), "pkg"), nil)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("an endless download must fail")
		}
		if !strings.Contains(err.Error(), "exceeded") {
			t.Errorf("error = %q, want it to say the download exceeded the overall bound", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("the overall timeout did not stop a download that never finishes")
	}
}

// A process killed mid-extraction leaves its staging directory behind. The next
// install of the same package sweeps the stale ones, and only those.
func TestTarXZ_SweepsStaleStagingDirectories(t *testing.T) {
	parent := t.TempDir()
	mk := func(name string, age time.Duration) string {
		dir := filepath.Join(parent, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("mkdir %s: %v", name, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "half.bin"), []byte("partial"), 0644); err != nil {
			t.Fatalf("write in %s: %v", name, err)
		}
		when := time.Now().Add(-age)
		if err := os.Chtimes(dir, when, when); err != nil {
			t.Fatalf("chtimes %s: %v", name, err)
		}
		return dir
	}

	staleOfThisPackage := mk(".pkg.partial-crashed", 3*time.Hour)
	freshOfThisPackage := mk(".pkg.partial-running", time.Minute) // another login mid-install
	staleOfAnotherPackage := mk(".other.partial-crashed", 3*time.Hour)
	if err := os.WriteFile(filepath.Join(parent, "keep.txt"), []byte("unrelated"), 0644); err != nil {
		t.Fatalf("write keep.txt: %v", err)
	}

	destination := filepath.Join(parent, "pkg")
	if err := TarXZ(bytes.NewReader(githubArchiveTarGz(t)), destination); err != nil {
		t.Fatalf("TarXZ failed: %v", err)
	}

	if _, err := os.Stat(staleOfThisPackage); !os.IsNotExist(err) {
		t.Error("a stale staging directory of this package was not swept")
	}
	if _, err := os.Stat(freshOfThisPackage); err != nil {
		t.Error("a fresh staging directory was removed; it may belong to an install still running")
	}
	if _, err := os.Stat(staleOfAnotherPackage); err != nil {
		t.Error("the sweep removed another package's staging directory")
	}
	if _, err := os.Stat(filepath.Join(parent, "keep.txt")); err != nil {
		t.Error("the sweep removed an unrelated file")
	}
	if _, err := os.Stat(filepath.Join(destination, "ext-v1.0.0", "index.js")); err != nil {
		t.Errorf("package was not installed: %v", err)
	}
}

// rawTarGz builds a gzip-compressed tar from headers, writing only the bodies
// given. It ignores the tar writer's complaint about unwritten body bytes, so a
// test can declare a size it never sends, which is what a hostile archive does.
func rawTarGz(t *testing.T, entries []rawEntry) []byte {
	t.Helper()

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		if err := tw.WriteHeader(e.header); err != nil {
			t.Fatalf("write header %q: %v", e.header.Name, err)
		}
		if len(e.body) > 0 {
			if _, err := tw.Write(e.body); err != nil {
				t.Fatalf("write body %q: %v", e.header.Name, err)
			}
		}
	}
	_ = tw.Flush()
	_ = gz.Close()
	return buf.Bytes()
}

type rawEntry struct {
	header *tar.Header
	body   []byte
}

// Two logins can both see a missing package and download it. The one that
// finishes second must not touch the first one's completed install.
func TestTarXZ_NeverOverwritesAnInstallThatAppearedMeanwhile(t *testing.T) {
	parent := t.TempDir()
	destination := filepath.Join(parent, "pkg")

	// Another login finished first.
	if err := os.MkdirAll(destination, 0755); err != nil {
		t.Fatal(err)
	}
	installed := filepath.Join(destination, "index.js")
	if err := os.WriteFile(installed, []byte("A: the completed install"), 0644); err != nil {
		t.Fatal(err)
	}

	t.Run("a complete second download leaves the first install alone", func(t *testing.T) {
		second := tarGzBytes(t, []tarEntry{{name: "index.js", body: "B: the second download"}})
		if err := TarXZ(bytes.NewReader(second), destination); err != nil {
			t.Fatalf("TarXZ returned an error although the package is already installed: %v", err)
		}
		got, _ := os.ReadFile(installed)
		if string(got) != "A: the completed install" {
			t.Errorf("the completed install was overwritten: %q", got)
		}
	})

	t.Run("a truncated second download cannot damage the first install", func(t *testing.T) {
		big := manyFilesTarGz(t, 20, 4096)
		_ = TarXZ(bytes.NewReader(big[:len(big)/2]), destination)
		got, _ := os.ReadFile(installed)
		if string(got) != "A: the completed install" {
			t.Errorf("a truncated download damaged the completed install: %q", got)
		}
		if names := dirEntries(t, destination); len(names) != 1 {
			t.Errorf("the completed install gained or lost files: %v", names)
		}
	})
}

// An existing, empty destination (the directory a caller made in advance) is
// still a valid place to install into.
func TestTarXZ_InstallsIntoAnExistingEmptyDirectory(t *testing.T) {
	destination := filepath.Join(t.TempDir(), "pkg")
	if err := os.MkdirAll(destination, 0755); err != nil {
		t.Fatal(err)
	}
	if err := TarXZ(bytes.NewReader(githubArchiveTarGz(t)), destination); err != nil {
		t.Fatalf("TarXZ failed for an empty existing directory: %v", err)
	}
	if _, err := os.Stat(filepath.Join(destination, "ext-v1.0.0", "index.js")); err != nil {
		t.Errorf("package was not installed into the existing directory: %v", err)
	}
}

// filepath.Dir("/x/pkg/") is "/x/pkg", so a destination with a trailing
// separator would stage inside itself.
func TestTarXZ_TrailingSeparatorDestination(t *testing.T) {
	parent := t.TempDir()
	destination := filepath.Join(parent, "pkg") + string(os.PathSeparator)

	if err := TarXZ(bytes.NewReader(githubArchiveTarGz(t)), destination); err != nil {
		t.Fatalf("TarXZ failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(parent, "pkg", "ext-v1.0.0", "index.js")); err != nil {
		t.Errorf("package was not installed at the destination: %v", err)
	}
	if names := dirEntries(t, parent); len(names) != 1 || names[0] != "pkg" {
		t.Errorf("unexpected entries beside the destination: %v", names)
	}
}

// Only a directory may stand for the root. A link or file named "." is not
// something to skip: links and special entries are rejected.
func TestTarXZ_RejectsNonDirectoryRootEntries(t *testing.T) {
	cases := map[string]*tar.Header{
		"symlink":  {Typeflag: tar.TypeSymlink, Name: ".", Linkname: "/etc", Mode: 0777},
		"hardlink": {Typeflag: tar.TypeLink, Name: "./", Linkname: "somewhere", Mode: 0644},
		"file":     {Typeflag: tar.TypeReg, Name: ".", Mode: 0644, Size: 0},
	}
	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			destination := filepath.Join(t.TempDir(), "pkg")
			err := TarXZ(bytes.NewReader(rawTarGz(t, []rawEntry{{header: header}})), destination)
			if err == nil {
				t.Fatalf("a %s entry named %q was accepted", name, header.Name)
			}
			if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
				t.Error("a rejected package left an install behind")
			}
		})
	}
}

// Adding a declared size to the running total wraps negative for a size near
// MaxInt64, which would let the next comparison pass.
func TestTarXZ_SizeCapSurvivesOverflow(t *testing.T) {
	archive := rawTarGz(t, []rawEntry{
		{header: &tar.Header{Typeflag: tar.TypeReg, Name: "small.txt", Mode: 0644, Size: 1}, body: []byte("x")},
		{header: &tar.Header{Typeflag: tar.TypeReg, Name: "huge.bin", Mode: 0644, Size: math.MaxInt64}},
	})
	destination := filepath.Join(t.TempDir(), "pkg")

	err := TarXZ(bytes.NewReader(archive), destination)
	if err == nil {
		t.Fatal("an entry declaring a size near MaxInt64 was accepted")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Errorf("error = %q, want it to report the size cap, not an I/O failure after writing", err)
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Error("an oversized package left an install behind")
	}
}

// Go does not apply global PAX records to later entries, but the format says a
// global path or size applies to every following entry. Skipping such a record
// would extract the wrong layout, so only the inert comment is accepted.
func TestTarXZ_GlobalPaxRecords(t *testing.T) {
	build := func(records map[string]string) []byte {
		return rawTarGz(t, []rawEntry{
			{header: &tar.Header{Typeflag: tar.TypeXGlobalHeader, Name: "pax_global_header", PAXRecords: records}},
			{header: &tar.Header{Typeflag: tar.TypeReg, Name: "index.js", Mode: 0644, Size: 4}, body: []byte("main")},
		})
	}

	t.Run("comment only is accepted", func(t *testing.T) {
		destination := filepath.Join(t.TempDir(), "pkg")
		if err := TarXZ(bytes.NewReader(build(map[string]string{"comment": "abc123"})), destination); err != nil {
			t.Fatalf("a comment-only global header was rejected: %v", err)
		}
	})

	for _, key := range []string{"path", "size"} {
		t.Run(key+" is rejected", func(t *testing.T) {
			destination := filepath.Join(t.TempDir(), "pkg")
			err := TarXZ(bytes.NewReader(build(map[string]string{key: "1"})), destination)
			if err == nil {
				t.Fatalf("a global %q record was ignored, which would extract the wrong layout", key)
			}
			if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
				t.Error("a rejected package left an install behind")
			}
		})
	}
}
