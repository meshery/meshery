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

func TestTruncateTo42Chars(t *testing.T) {
	exact := strings.Repeat("a", 42)
	oneLonger := strings.Repeat("a", 43)
	// 41 single-byte runes followed by a 3-byte rune: byte 42 falls
	// inside the euro sign, which byte slicing would split.
	multibyte := strings.Repeat("a", 41) + "€" + strings.Repeat("b", 10)

	tests := []struct {
		name          string
		input         string
		want          string
		wantTruncated bool
	}{
		{
			name:          "exactly 42 runes is unchanged",
			input:         exact,
			want:          exact,
			wantTruncated: false,
		},
		{
			name:          "43 runes is truncated",
			input:         oneLonger,
			want:          strings.Repeat("a", 42),
			wantTruncated: true,
		},
		{
			name:          "multibyte character split by byte 42 is dropped whole",
			input:         multibyte,
			want:          strings.Repeat("a", 41) + "€",
			wantTruncated: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, truncated := truncateTo42Chars(tt.input)
			if got != tt.want {
				t.Fatalf("truncateTo42Chars() = %q, want %q", got, tt.want)
			}
			if truncated != tt.wantTruncated {
				t.Fatalf("truncated = %v, want %v", truncated, tt.wantTruncated)
			}
			if !utf8.ValidString(got) {
				t.Fatalf("truncateTo42Chars() returned invalid UTF-8: %q", got)
			}
		})
	}
}
