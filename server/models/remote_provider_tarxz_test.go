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

func TestTarXZ_ExtractsArchiveAndCloses(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tarxz-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)

	fileContent := []byte("hello meshery remote provider")
	header := &tar.Header{
		Name:     "test-file.txt",
		Mode:     0600,
		Size:     int64(len(fileContent)),
		Typeflag: tar.TypeReg,
	}
	require.NoError(t, tw.WriteHeader(header))
	_, err = tw.Write(fileContent)
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

func TestTarXZ_InvalidGzipStream(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "tarxz-invalid-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	invalidData := bytes.NewReader([]byte("invalid-non-gzip-stream"))
	err = TarXZ(invalidData, tmpDir)
	assert.Error(t, err)
}
