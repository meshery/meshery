package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
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

func TestTruncateEventDescription(t *testing.T) {
	long := strings.Repeat("a", 50)
	accented := "Imported model " + strings.Repeat("é", 40)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "empty", in: "", want: ""},
		{name: "shorter than the limit", in: "model imported", want: "model imported"},
		{name: "exactly the limit is not shortened", in: strings.Repeat("a", 42), want: strings.Repeat("a", 42)},
		{name: "one over the limit", in: strings.Repeat("a", 43), want: strings.Repeat("a", 42) + "..."},
		{name: "long ASCII", in: long, want: long[:42] + "..."},
		{name: "multibyte characters count as one", in: strings.Repeat("é", 42), want: strings.Repeat("é", 42)},
		{name: "multibyte cut on a character boundary", in: accented, want: string([]rune(accented)[:42]) + "..."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateEventDescription(tt.in)
			if got != tt.want {
				t.Errorf("truncateEventDescription(%q) = %q, want %q", tt.in, got, tt.want)
			}
			if !utf8.ValidString(got) {
				t.Errorf("truncateEventDescription(%q) = %q is not valid UTF-8", tt.in, got)
			}
		})
	}
}
