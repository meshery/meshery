//go:build !windows

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// restrictTokenFile must tighten a token file that an older mesheryctl left
// readable by other users down to owner-only.
func TestRestrictTokenFileTightensMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	// os.WriteFile is subject to the umask; chmod sets the mode exactly.
	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()

	if err := restrictTokenFile(f); err != nil {
		t.Fatalf("restrictTokenFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if mode := info.Mode().Perm(); mode != tokenFileMode {
		t.Fatalf("token file mode = %o, want %o", mode, tokenFileMode)
	}
}

// Rewriting an existing token file keeps it owner-only.
func TestWriteTokenFileRewriteKeepsMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	for i := 0; i < 2; i++ {
		if err := WriteTokenFile(path, []byte(`{"token":"t"}`)); err != nil {
			t.Fatalf("WriteTokenFile #%d: %v", i+1, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if mode := info.Mode().Perm(); mode != tokenFileMode {
			t.Fatalf("after write #%d token file mode = %o, want %o", i+1, mode, tokenFileMode)
		}
	}
}
