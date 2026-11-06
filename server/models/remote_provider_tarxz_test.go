package models

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTarXZ_ExtractsArchive(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	fileContent := []byte("hello meshery remote provider")
	header := &tar.Header{
		Name:     "test-file.txt",
		Mode:     0o600,
		Size:     int64(len(fileContent)),
		Typeflag: tar.TypeReg,
	}
	require.NoError(t, tw.WriteHeader(header))
	_, err := tw.Write(fileContent)
	require.NoError(t, err)

	require.NoError(t, tw.Close())
	require.NoError(t, gw.Close())

	err = TarXZ(&buf, tmpDir)
	assert.NoError(t, err)

	extractedPath := filepath.Join(tmpDir, "test-file.txt")
	assert.FileExists(t, extractedPath)

	content, err := os.ReadFile(extractedPath)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, content)
}

func TestTarXZ_PathTraversalBlocked(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	fileContent := []byte("malicious content")
	header := &tar.Header{
		Name:     "../escape.txt",
		Mode:     0o600,
		Size:     int64(len(fileContent)),
		Typeflag: tar.TypeReg,
	}
	require.NoError(t, tw.WriteHeader(header))
	_, err := tw.Write(fileContent)
	require.NoError(t, err)

	require.NoError(t, tw.Close())
	require.NoError(t, gw.Close())

	err = TarXZ(&buf, tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "illegal file path")
}

func TestTarXZ_RejectsUnsupportedType(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	header := &tar.Header{
		Name:     "test-symlink",
		Typeflag: tar.TypeSymlink,
		Linkname: "/etc/passwd",
	}
	require.NoError(t, tw.WriteHeader(header))

	require.NoError(t, tw.Close())
	require.NoError(t, gw.Close())

	err := TarXZ(&buf, tmpDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown type")
}

func TestTarXZ_InvalidGzipStream(t *testing.T) {
	tmpDir := t.TempDir()

	invalidData := bytes.NewReader([]byte("invalid-non-gzip-stream"))
	err := TarXZ(invalidData, tmpDir)
	assert.Error(t, err)
}
