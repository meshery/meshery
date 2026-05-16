package root

import (
	"errors"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestFormatCLIError(t *testing.T) {
	t.Run("uses meshkit long description in non verbose mode", func(t *testing.T) {
		err := utils.ErrInvalidArgument(errors.New("this command takes no arguments"))
		msg := formatCLIError(err, false)
		expected := "this command takes no arguments"
		if msg != expected {
			t.Fatalf("expected %q, got %q", expected, msg)
		}
	})

	t.Run("returns full error in verbose mode", func(t *testing.T) {
		err := utils.ErrInvalidArgument(errors.New("this command takes no arguments"))
		msg := formatCLIError(err, true)
		if msg != err.Error() {
			t.Fatalf("expected verbose message to match original error")
		}
	})
}
