package utils

import (
	"os"
	"testing"
)

// TestIsTerminal pins what the interactive-terminal probe decides on, using the
// shapes stdin/stdout take when nobody is at the keyboard: a pipe
// (`mesheryctl connection list | grep ...`, a bats `run`), /dev/null (a CI job,
// a cron entry) and a regular file (output redirected to a log).
//
// /dev/null is the case a file-mode check gets wrong - it is a character device
// but not a terminal - which is why this asks the terminal driver instead.
func TestIsTerminal(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("open pipe: %v", err)
	}
	t.Cleanup(func() {
		_ = r.Close()
		_ = w.Close()
	})

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatalf("open %s: %v", os.DevNull, err)
	}
	t.Cleanup(func() { _ = devNull.Close() })

	regular, err := os.CreateTemp(t.TempDir(), "interactive")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Cleanup(func() { _ = regular.Close() })

	cases := []struct {
		name string
		file *os.File
		want bool
	}{
		{name: "pipe read end", file: r, want: false},
		{name: "pipe write end", file: w, want: false},
		{name: "dev null", file: devNull, want: false},
		{name: "regular file", file: regular, want: false},
		{name: "nil file", file: nil, want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTerminal(tc.file); got != tc.want {
				t.Errorf("isTerminal(%s) = %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
