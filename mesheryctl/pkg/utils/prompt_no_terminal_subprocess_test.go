package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

// The other tests in this package replace IsInteractiveTerminal with a stub, so
// they pin the branch but never the probe that chooses it. This one runs the
// real thing in a real process whose stdin is a pipe rather than a terminal -
// the shape stdin has under CI, cron, and `mesheryctl ... | grep`.
//
// Go's standard re-exec pattern: the test binary runs itself with a marker in
// the environment, and the marked run does the work and reports back.
const subprocessMarker = "MESHERYCTL_NO_TERMINAL_SUBPROCESS"

func TestRunSelectPromptInARealNonTerminalProcess(t *testing.T) {
	if os.Getenv(subprocessMarker) == "1" {
		// Child: IsInteractiveTerminal is untouched here.
		_, err := RunSelectPrompt("Select item", []string{"a", "b", "c"})
		if err == nil {
			fmt.Print("NOERROR")
			return
		}
		fmt.Print(errors.GetCode(err) + "\n" + err.Error())
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRunSelectPromptInARealNonTerminalProcess")
	cmd.Env = append(os.Environ(), subprocessMarker+"=1")

	// A pipe on both ends. Nothing is written to stdin, so a prompt that got
	// past the guard would block rather than pick something - which is the
	// failure this guard exists to prevent.
	stdin, stdinWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	defer func() { _ = stdinWriter.Close() }()
	cmd.Stdin = stdin

	out, err := cmd.Output()
	_ = stdin.Close()
	if err != nil {
		t.Fatalf("subprocess failed: %v (output: %s)", err, out)
	}

	got := string(out)
	if strings.HasPrefix(got, "NOERROR") {
		t.Fatal("a prompt with no terminal returned no error; the guard did not fire")
	}
	if code, _, _ := strings.Cut(got, "\n"); code != ErrNoTerminalForPromptCode {
		t.Errorf("error code = %q, want %q", code, ErrNoTerminalForPromptCode)
	}
	if strings.Contains(strings.ToLower(got), "cancel") {
		t.Errorf("a missing terminal must not be reported as a cancellation, got: %s", got)
	}
}
