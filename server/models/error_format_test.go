package models

import (
	"errors"
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrMarshalYAML is a regression test for issue #22681: the short
// description carried a bare "%s" format verb that was never fed through
// Sprintf, so the obj value was dropped and the literal "%s" surfaced to the
// user. The rendered message must now contain the interpolated obj value and
// no literal format verb.
func TestErrMarshalYAML(t *testing.T) {
	underlying := errors.New("yaml: line 3: mapping values are not allowed")
	const obj = "my-design.yaml"

	err := ErrMarshalYAML(underlying, obj)
	if err == nil {
		t.Fatal("ErrMarshalYAML returned nil, expected a wrapped error")
	}

	if meshkiterrors.GetCode(err) != ErrMarshalYAMLCode {
		t.Errorf("expected error code %q, got %q", ErrMarshalYAMLCode, meshkiterrors.GetCode(err))
	}

	sdesc := meshkiterrors.GetSDescription(err)
	if !strings.Contains(sdesc, obj) {
		t.Errorf("short description %q does not contain interpolated obj %q", sdesc, obj)
	}
	if strings.Contains(sdesc, "%s") || strings.Contains(sdesc, "%v") {
		t.Errorf("short description %q still contains a literal format verb", sdesc)
	}

	if msg := err.Error(); !strings.Contains(msg, obj) {
		t.Errorf("Error() %q does not contain interpolated obj %q", msg, obj)
	}
}
