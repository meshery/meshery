package utils

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

// RunSelectPrompt asks the user to choose from a list. With no terminal there
// is nobody to choose, and promptui returns promptui.ErrEOF - which every call
// site used to report as something else. Four of them said "cancelled", which
// sends an operator looking for a Ctrl+C that never happened.
func TestRunSelectPromptWithoutTerminal(t *testing.T) {
	original := IsInteractiveTerminal
	IsInteractiveTerminal = func() bool { return false }
	t.Cleanup(func() { IsInteractiveTerminal = original })

	index, err := RunSelectPrompt("Select item", []string{"a", "b", "c"})

	if err == nil {
		t.Fatalf("expected an error without a terminal, got index %d and nil", index)
	}
	if index != 0 {
		t.Errorf("index should be the zero value alongside an error, got %d", index)
	}
	if code := errors.GetCode(err); code != ErrNoTerminalForPromptCode {
		t.Errorf("error code = %s, want %s", code, ErrNoTerminalForPromptCode)
	}
	// The distinction this whole change exists for.
	if code := errors.GetCode(err); code == ErrPromptCancelledCode {
		t.Errorf("a missing terminal must not be reported as a cancelled prompt")
	}
}

// The message has to name the missing choice and a way to avoid it, or a CI log
// says only that something failed.
func TestErrNoTerminalForPromptNamesTheChoiceAndTheWayOut(t *testing.T) {
	err := ErrNoTerminalForPrompt("which design to deploy", "Pass the design ID instead")

	got := err.Error()
	for _, want := range []string{"which design to deploy", "Pass the design ID instead"} {
		if !strings.Contains(got, want) {
			t.Errorf("error should mention %q, got: %s", want, got)
		}
	}
	if strings.Contains(strings.ToLower(got), "cancel") {
		t.Errorf("must not describe a missing terminal as a cancellation, got: %s", got)
	}
}

// A terminal being present must leave the existing path alone.
func TestRunSelectPromptStillPromptsWithATerminal(t *testing.T) {
	original := IsInteractiveTerminal
	IsInteractiveTerminal = func() bool { return true }
	t.Cleanup(func() { IsInteractiveTerminal = original })

	// promptui will fail on the test's non-tty stdin, but the point is that it
	// gets that far: the error must come from the prompt, not from the guard.
	_, err := RunSelectPrompt("Select item", []string{"a", "b"})

	if err != nil && errors.GetCode(err) == ErrNoTerminalForPromptCode {
		t.Error("the guard fired even though a terminal was reported present")
	}
}
