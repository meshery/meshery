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
