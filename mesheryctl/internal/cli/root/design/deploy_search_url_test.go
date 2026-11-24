package design

import "testing"

// Test_getDesignDeploySearchURLPath_EncodesNameOnce pins the single-encoding
// contract for the design lookup query.
//
// deploy.go used to pre-escape the name with url.QueryEscape and then hand it to
// url.Values.Encode(), which encodes again. "My Design" went out as
// search=My%2BDesign and the server read back the literal "My+Design", so the
// LIKE matched nothing and the command reported `design not found`. Every name
// that percent-encodes was affected, not only ones with spaces.
//
// Wire values here are what url.Values.Encode() produces; the server recovers
// the original name from them via r.URL.Query() (which decodes "+" to a space).
func Test_getDesignDeploySearchURLPath_EncodesNameOnce(t *testing.T) {
	const baseAPIPath = "api/pattern"

	tests := []struct {
		name       string
		designName string
		want       string
	}{
		{
			// The reported case. Double-encoding sent search=My%2BDesign.
			name:       "space is encoded once",
			designName: "My Design",
			want:       "api/pattern?populate=pattern_file&search=My+Design",
		},
		{
			// Guards the path that worked before: nothing to encode, so a
			// regression here would mean the fix broke the common case.
			name:       "single token is unchanged",
			designName: "nginx",
			want:       "api/pattern?populate=pattern_file&search=nginx",
		},
		{
			// "&" must not terminate the parameter.
			name:       "ampersand is escaped, not a separator",
			designName: "A&B",
			want:       "api/pattern?populate=pattern_file&search=A%26B",
		},
		{
			// A literal "+" has to survive as %2B, or the server decodes it to
			// a space and searches for the wrong name.
			name:       "literal plus stays distinct from a space",
			designName: "C+D",
			want:       "api/pattern?populate=pattern_file&search=C%2BD",
		},
		{
			// Non-ASCII was corrupted too: the double-encoded wire was
			// search=caf%25C3%25A9, which the server read back as the literal
			// "caf%C3%A9" rather than "café". The value below is the correct
			// single-encoded wire form.
			name:       "non-ascii is encoded once",
			designName: "café",
			want:       "api/pattern?populate=pattern_file&search=caf%C3%A9",
		},
		{
			// A bare "%" is the encoder's own escape character, so an
			// unescaped one would make the server's decode fail rather than
			// return the wrong design.
			name:       "percent is escaped, not read as an escape",
			designName: "100% cotton",
			want:       "api/pattern?populate=pattern_file&search=100%25+cotton",
		},
		{
			// "#" would start a fragment and "?" a second query string, either
			// of which truncates the value in transit if left unescaped.
			name:       "fragment and query delimiters do not truncate",
			designName: "v1.0#beta?test",
			want:       "api/pattern?populate=pattern_file&search=v1.0%23beta%3Ftest",
		},
		{
			// The name is encoded exactly as the caller built it; this helper
			// does not normalise. Trimming would widen the server's LIKE and
			// change which designs match, which is out of scope for an
			// encoding fix - see the discussion on #21847.
			name:       "surrounding whitespace is preserved, not trimmed",
			designName: "  My Design  ",
			want:       "api/pattern?populate=pattern_file&search=++My+Design++",
		},
		{
			// No args reduce to an empty search rather than a malformed URL.
			name:       "empty name yields an empty search",
			designName: "",
			want:       "api/pattern?populate=pattern_file&search=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Safe without rebinding tt: range variables are per-iteration
			// from Go 1.22, and this module is go 1.26.
			t.Parallel()

			if got := getDesignDeploySearchURLPath(baseAPIPath, tt.designName); got != tt.want {
				t.Errorf("getDesignDeploySearchURLPath(%q, %q)\n got: %s\nwant: %s", baseAPIPath, tt.designName, got, tt.want)
			}
		})
	}
}
