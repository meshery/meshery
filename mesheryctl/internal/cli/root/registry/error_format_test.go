package registry

import (
	"errors"
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrGenerateComponentFormat is a regression test for issue #22681: the
// short description carried bare "%s" format verbs that were never fed through
// Sprintf, so compName and modelName were appended as stray slice elements and
// the literal "%s" surfaced to the user. The short description must now contain
// both interpolated values and no literal format verb.
func TestErrGenerateComponentFormat(t *testing.T) {
	const (
		compName  = "my-component"
		modelName = "my-model"
	)
	underlying := errors.New("schema validation failed")

	err := ErrGenerateComponent(underlying, modelName, compName)
	if err == nil {
		t.Fatal("ErrGenerateComponent returned nil, expected a wrapped error")
	}

	sdesc := meshkiterrors.GetSDescription(err)
	if !strings.Contains(sdesc, compName) {
		t.Errorf("short description %q does not contain interpolated compName %q", sdesc, compName)
	}
	if !strings.Contains(sdesc, modelName) {
		t.Errorf("short description %q does not contain interpolated modelName %q", sdesc, modelName)
	}
	if strings.Contains(sdesc, "%s") || strings.Contains(sdesc, "%v") {
		t.Errorf("short description %q still contains a literal format verb", sdesc)
	}
}

// TestErrParsingSheetFormat is a regression test for issue #22753: the probable
// cause carried a bare "%s" that was never fed through Sprintf, so the literal
// verb surfaced to the user. The probable cause must now contain the sheet name
// and no literal format verb. The short and long descriptions already contain
// the sheet name, so the assertion targets the probable cause specifically.
func TestErrParsingSheetFormat(t *testing.T) {
	err := ErrParsingSheet(errors.New("boom"), "Models")
	cause := meshkiterrors.GetCause(err)

	if !strings.Contains(cause, "The provided sheet id might be incorrect. Sheet: Models") {
		t.Errorf("expected the sheet name in the probable cause, got: %s", cause)
	}
	if strings.Contains(cause, "%s") || strings.Contains(cause, "%v") {
		t.Errorf("literal format verb leaked into the probable cause: %s", cause)
	}
}
