package root

import (
	"errors"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestFormatCLIError(t *testing.T) {
	t.Run("uses meshkit long description in non verbose mode", func(t *testing.T) {
		err := utils.ErrInvalidArgument(errors.New("this command takes no arguments"))
		msg := formatCLIError(err, false)
		expected := "this command takes no arguments\nPlease check the arguments passed"
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

	t.Run("joins multi-element long description with space not period", func(t *testing.T) {
		// ErrCreateFile has a 2-element LongDescription; without this fix GetLDescription joins with "." (no space).
		err := utils.ErrCreateFile("/tmp/test", errors.New("permission denied"))
		msg := formatCLIError(err, false)
		if strings.Contains(msg, ".permission") {
			t.Fatalf("expected space-separated join, got period-concatenated: %q", msg)
		}
		if !strings.Contains(msg, "permission denied") {
			t.Fatalf("expected error detail in message, got: %q", msg)
		}
	})

	t.Run("strips pipe-delimited metadata from plain errors", func(t *testing.T) {
		// Simulate what MeshKit Error.Error() produces for errors that go through fallback path.
		raw := errors.New("model not found | Short Description: Internal Server Error | Probable Cause: something | Suggested Remediation: try again")
		msg := formatCLIError(raw, false)
		if msg != "model not found" {
			t.Fatalf("expected stripped message, got %q", msg)
		}
	})
}
