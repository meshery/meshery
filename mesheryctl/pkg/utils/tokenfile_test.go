package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// A token refresh writes a shorter token over a longer one. The file must end
// up holding exactly the new bytes, with nothing left over from the old token,
// whichever platform restricts it.
func TestWriteTokenFileReplacesLongerContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	long := []byte(`{"meshery-provider":"Meshery","token":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`)
	short := []byte(`{"token":"b"}`)

	if err := WriteTokenFile(path, long); err != nil {
		t.Fatalf("first WriteTokenFile: %v", err)
	}
	if err := WriteTokenFile(path, short); err != nil {
		t.Fatalf("second WriteTokenFile: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(short) {
		t.Fatalf("token file content = %q, want %q", got, short)
	}
}

// Logout writes an empty token file. It must succeed and leave an empty file.
func TestWriteTokenFileEmptyData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	if err := WriteTokenFile(path, []byte("secret")); err != nil {
		t.Fatalf("WriteTokenFile: %v", err)
	}
	if err := WriteTokenFile(path, []byte{}); err != nil {
		t.Fatalf("WriteTokenFile with empty data: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("token file content = %q, want empty", got)
	}
}

// A path whose directory does not exist must fail without creating anything.
func TestWriteTokenFileMissingDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "auth.json")

	if err := WriteTokenFile(path, []byte("x")); err == nil {
		t.Fatal("WriteTokenFile succeeded for a path in a missing directory")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("token file should not exist, Stat error = %v", err)
	}
}
