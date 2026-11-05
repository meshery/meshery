package design

import (
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// TestDesignRemediationFormat is a regression test for issue #22681: the
// suggested-remediation strings carried a bare "%v" format verb that was never
// fed through Sprintf, so the usage message was appended as a stray slice
// element and the literal "%v" surfaced to the user. Each remediation must now
// contain the interpolated usage text and no literal format verb.
func TestDesignRemediationFormat(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string // a distinctive substring of the interpolated usage message
	}{
		{"ErrDesignManifest", ErrDesignManifest(), "mesheryctl design import"},
		{"ErrDesignFileNotProvided", ErrDesignFileNotProvided(), "mesheryctl design import"},
		{"ErrDesignSourceType", ErrDesignSourceType(), "mesheryctl design import"},
		{"ErrUndeployDesign", ErrUndeployDesign(nil), "mesheryctl design undeploy"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Fatalf("%s returned nil, expected a wrapped error", tc.name)
			}

			remedy := meshkiterrors.GetRemedy(tc.err)
			if !strings.Contains(remedy, tc.want) {
				t.Errorf("suggested remediation %q does not contain interpolated usage text %q", remedy, tc.want)
			}
			if strings.Contains(remedy, "%s") || strings.Contains(remedy, "%v") {
				t.Errorf("suggested remediation %q still contains a literal format verb", remedy)
			}
		})
	}
}
