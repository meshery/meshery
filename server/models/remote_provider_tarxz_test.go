package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type tarEntry struct {
	name string
	body string
	dir  bool
	link bool
}

// tarGzBytes builds an in-memory gzip-compressed tar archive. Keeping the
// fixture in memory means the tests do not depend on the remote provider or on
// a checked-in binary blob.
func tarGzBytes(t *testing.T, entries []tarEntry) []byte {
	t.Helper()

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for _, entry := range entries {
		header := &tar.Header{
			Name: entry.name,
			Mode: 0644,
			Size: int64(len(entry.body)),
		}

		switch {
		case entry.dir:
			header.Typeflag = tar.TypeDir
			header.Mode = 0755
			header.Size = 0
		case entry.link:
			header.Typeflag = tar.TypeSymlink
			header.Linkname = entry.body
			header.Size = 0
		default:
			header.Typeflag = tar.TypeReg
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("write tar header %q: %v", entry.name, err)
		}

		if !entry.dir && !entry.link {
			if _, err := tarWriter.Write([]byte(entry.body)); err != nil {
				t.Fatalf("write tar body %q: %v", entry.name, err)
			}
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("close tar writer: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("close gzip writer: %v", err)
	}

	return buf.Bytes()
}

func TestTarXZExtractsFilesAndDirectories(t *testing.T) {
	archive := tarGzBytes(t, []tarEntry{
		{name: "top.txt", body: "top level"},
		{name: "nested", dir: true},
		{name: "nested/deep.txt", body: "nested file"},
	})

	destination := t.TempDir()
	if err := TarXZ(bytes.NewReader(archive), destination); err != nil {
		t.Fatalf("TarXZ returned an error for a valid archive: %v", err)
	}

	for _, tc := range []struct {
		path string
		want string
	}{
		{"top.txt", "top level"},
		{"nested/deep.txt", "nested file"},
	} {
		got, err := os.ReadFile(filepath.Join(destination, filepath.FromSlash(tc.path)))
		if err != nil {
			t.Fatalf("read extracted file %q: %v", tc.path, err)
		}
		if string(got) != tc.want {
			t.Errorf("content of %q = %q, want %q", tc.path, string(got), tc.want)
		}
	}
}

func TestTarXZRejectsPathTraversal(t *testing.T) {
	archive := tarGzBytes(t, []tarEntry{
		{name: "../escaped.txt", body: "should not be written"},
	})

	parent := t.TempDir()
	destination := filepath.Join(parent, "destination")
	escapedPath := filepath.Join(parent, "escaped.txt")

	err := TarXZ(bytes.NewReader(archive), destination)
	if err == nil {
		t.Fatal("TarXZ accepted an archive with a path traversal entry")
	}
	if !strings.Contains(err.Error(), "illegal file path") {
		t.Errorf("error = %q, want it to mention an illegal file path", err)
	}
	if _, statErr := os.Stat(escapedPath); !os.IsNotExist(statErr) {
		t.Errorf("path traversal entry wrote %q outside the destination", escapedPath)
	}
}

func TestTarXZRejectsUnsupportedEntryType(t *testing.T) {
	archive := tarGzBytes(t, []tarEntry{
		{name: "link", body: "/etc/passwd", link: true},
	})

	err := TarXZ(bytes.NewReader(archive), t.TempDir())
	if err == nil {
		t.Fatal("TarXZ accepted a symlink entry")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("error = %q, want it to report an unknown entry type", err)
	}
}

func TestTarXZRejectsNonGzipStream(t *testing.T) {
	err := TarXZ(strings.NewReader("this is not a gzip stream"), t.TempDir())
	if err == nil {
		t.Fatal("TarXZ accepted a stream that is not gzip-compressed")
	}
}
