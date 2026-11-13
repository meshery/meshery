package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCreateTemp pins the contract the model import handler relies on: the
// returned path holds exactly the given bytes, the file is already closed, and
// the caller owns only its removal. The handler used to receive an *os.File it
// closed a second time, logging "file already closed" at error severity on
// every model import.
func TestCreateTemp(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	data := []byte("model archive bytes")

	path, err := CreateTemp("kubecon-logo.tar", data)
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })

	if !strings.HasPrefix(filepath.Base(path), "kubecon-logo.tar") {
		t.Errorf("temp file %q does not carry the requested name prefix", path)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading temp file: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("temp file content = %q, want %q", got, data)
	}
	if err := os.Remove(path); err != nil {
		t.Errorf("caller must be able to remove the temp file: %v", err)
	}
}

func TestCreateTempFailureLeavesNoFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("TMPDIR", dir)

	// A path separator in the pattern is rejected by os.CreateTemp.
	if _, err := CreateTemp("bad/name", []byte("x")); err == nil {
		t.Fatal("CreateTemp with a path separator in the pattern must fail")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("failed CreateTemp left %d file(s) behind", len(entries))
	}
}
