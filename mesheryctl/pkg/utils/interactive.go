package utils

import (
	"os"

	"golang.org/x/term"
)

// IsInteractiveTerminal reports whether mesheryctl is attached to a terminal a
// human can answer from, i.e. whether interactive affordances (paging on a
// keypress, selection prompts) are possible at all.
//
// Both streams have to be terminals. Stdin is where the keystroke would come
// from and stdout is where the prompt would be shown, and either being
// redirected means there is nobody on the other end: `mesheryctl connection
// list | grep ...` and a bats `run` both take stdout; a CI job and a cron entry
// both take stdin.
//
// This matters because the keyboard library used for interactive paging reads
// the *controlling terminal* directly (/dev/tty). Where there isn't one it
// fails with "open /dev/tty: no such device or address", and a list command
// that had already printed its first page exited non-zero on that error -
// making mesheryctl unusable in scripts, pipelines and CI for any list longer
// than a page. Callers use this to fall back to non-interactive behaviour
// instead of failing.
//
// It is a var so tests can state which environment they are describing.
var IsInteractiveTerminal = func() bool {
	return isTerminal(os.Stdin) && isTerminal(os.Stdout)
}

// isTerminal reports whether f is a terminal. It asks the terminal driver
// rather than inspecting the file mode, so that /dev/null - a character device
// that is emphatically not a terminal, and what stdin usually is in CI - is
// answered correctly.
func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
