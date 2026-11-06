package adapter

import (
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestErrValidMeshNameFormat is a regression test for issue #22681: the
// probable-cause string carried a bare "%s" format verb that was never fed
// through Sprintf, so meshName was appended as a stray slice element and the
// literal "%s" surfaced to the user. The probable cause must now contain the
// interpolated meshName and no literal format verb.
func TestErrValidMeshNameFormat(t *testing.T) {
	const meshName = "bogus-mesh"

	err := ErrValidMeshName(meshName)
	if err == nil {
		t.Fatal("ErrValidMeshName returned nil, expected a wrapped error")
	}

	cause := meshkiterrors.GetCause(err)
	if !strings.Contains(cause, meshName) {
		t.Errorf("probable cause %q does not contain interpolated meshName %q", cause, meshName)
	}
	if strings.Contains(cause, "%s") || strings.Contains(cause, "%v") {
		t.Errorf("probable cause %q still contains a literal format verb", cause)
	}
}
