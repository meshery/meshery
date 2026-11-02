package utils

import (
	"errors"
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrParseGithubFileFormat is a regression test for issue #22681: the
// probable-cause string carried a bare "%s" format verb that was never fed
// through Sprintf, so the URL was appended as a stray slice element and the
// literal "%s" surfaced to the user. The probable cause must now contain the
// interpolated URL and no literal format verb.
func TestErrParseGithubFileFormat(t *testing.T) {
	const url = "https://github.com/meshery/meshery/blob/master/README.md"
	underlying := errors.New("unexpected EOF")

	err := ErrParseGithubFile(underlying, url)
	if err == nil {
		t.Fatal("ErrParseGithubFile returned nil, expected a wrapped error")
	}

	cause := meshkiterrors.GetCause(err)
	if !strings.Contains(cause, url) {
		t.Errorf("probable cause %q does not contain interpolated URL %q", cause, url)
	}
	if strings.Contains(cause, "%s") || strings.Contains(cause, "%v") {
		t.Errorf("probable cause %q still contains a literal format verb", cause)
	}
}

// TestErrMarkFlagRequireFormat is a regression test for issue #22681: the
// suggested-remediation string carried a bare "%s" format verb with no
// argument, so the literal "%s" surfaced to the user. The intended value is the
// flag name (as already used in the short description). The remediation must now
// contain the interpolated flagName and no literal format verb.
func TestErrMarkFlagRequireFormat(t *testing.T) {
	const flagName = "output"
	underlying := errors.New("flag not registered")

	err := ErrMarkFlagRequire(flagName, underlying)
	if err == nil {
		t.Fatal("ErrMarkFlagRequire returned nil, expected a wrapped error")
	}

	remedy := meshkiterrors.GetRemedy(err)
	if !strings.Contains(remedy, flagName) {
		t.Errorf("suggested remediation %q does not contain interpolated flagName %q", remedy, flagName)
	}
	if strings.Contains(remedy, "%s") || strings.Contains(remedy, "%v") {
		t.Errorf("suggested remediation %q still contains a literal format verb", remedy)
	}
}
