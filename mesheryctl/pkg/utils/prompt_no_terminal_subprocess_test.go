package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

// The test binary re-runs itself with this marker set; the marked run is the
// child whose stdin is a pipe.
const subprocessMarker = "MESHERYCTL_NO_TERMINAL_SUBPROCESS"

const subprocessTestName = "TestRunSelectPrompt_Given_Stdin_Is_A_Pipe_When_Prompting_Then_ErrNoTerminalForPrompt"

func TestRunSelectPrompt_Given_Stdin_Is_A_Pipe_When_Prompting_Then_ErrNoTerminalForPrompt(t *testing.T) {
	if os.Getenv(subprocessMarker) == "1" {
		_, err := RunSelectPrompt("Select item", []string{"a", "b", "c"})
		if err == nil {
			fmt.Print("NOERROR")
			return
		}
		fmt.Print(errors.GetCode(err) + "\n" + err.Error())
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=^"+subprocessTestName+"$")
	cmd.Env = append(os.Environ(), subprocessMarker+"=1")

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
		t.Fatal("expected an error with stdin on a pipe, got none")
	}
	if code, _, _ := strings.Cut(got, "\n"); code != ErrNoTerminalForPromptCode {
		t.Errorf("error code = %q, want %q", code, ErrNoTerminalForPromptCode)
	}
	if strings.Contains(strings.ToLower(got), "cancel") {
		t.Errorf("error should not mention a cancellation, got: %s", got)
	}
}
