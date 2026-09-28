package utils

import (
	"strings"
	"testing"

	"github.com/meshery/meshkit/errors"
)

func TestRunSelectPrompt(t *testing.T) {
	stubTerminal := func(t *testing.T, present bool) {
		original := IsInteractiveTerminal
		IsInteractiveTerminal = func() bool { return present }
		t.Cleanup(func() { IsInteractiveTerminal = original })
	}

	t.Run("Given_No_Terminal_When_RunSelectPrompt_Then_ErrNoTerminalForPrompt_Not_ErrPromptCancelled", func(t *testing.T) {
		stubTerminal(t, false)

		index, err := RunSelectPrompt("Select item", []string{"a", "b", "c"})

		if err == nil {
			t.Fatalf("expected an error, got index %d and nil", index)
		}
		if index != 0 {
			t.Errorf("index = %d, want 0 alongside an error", index)
		}
		if code := errors.GetCode(err); code != ErrNoTerminalForPromptCode {
			t.Errorf("error code = %s, want %s", code, ErrNoTerminalForPromptCode)
		}
		if code := errors.GetCode(err); code == ErrPromptCancelledCode {
			t.Errorf("error code = %s, must not be %s", code, ErrPromptCancelledCode)
		}
	})

	t.Run("Given_A_Terminal_When_RunSelectPrompt_Then_The_Guard_Does_Not_Fire", func(t *testing.T) {
		stubTerminal(t, true)

		_, err := RunSelectPrompt("Select item", []string{"a", "b"})

		if err != nil && errors.GetCode(err) == ErrNoTerminalForPromptCode {
			t.Errorf("got %s with a terminal present", ErrNoTerminalForPromptCode)
		}
	})
}

func TestErrNoTerminalForPrompt(t *testing.T) {
	t.Run("Given_A_Choice_And_A_Hint_When_ErrNoTerminalForPrompt_Then_The_Message_Names_Both_And_Not_Cancel", func(t *testing.T) {
		err := ErrNoTerminalForPrompt("which design to deploy", "Pass the design ID instead")

		got := err.Error()
		for _, want := range []string{"which design to deploy", "Pass the design ID instead"} {
			if !strings.Contains(got, want) {
				t.Errorf("error should mention %q, got: %s", want, got)
			}
		}
		if strings.Contains(strings.ToLower(got), "cancel") {
			t.Errorf("error should not mention a cancellation, got: %s", got)
		}
	})

	t.Run("Given_Any_Choice_When_ErrNoTerminalForPrompt_Then_The_Short_Description_Is_A_Fixed_String", func(t *testing.T) {
		first := errors.GetSDescription(ErrNoTerminalForPrompt("a design to deploy", "hint"))
		second := errors.GetSDescription(ErrNoTerminalForPrompt("a context to switch to", "hint"))

		if first == "" || first != second {
			t.Errorf("short descriptions = %q and %q, want the same non-empty literal", first, second)
		}
	})
}
