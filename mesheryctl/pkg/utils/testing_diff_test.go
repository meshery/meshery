package utils

import (
	"strings"
	"testing"
)

// The whole value of the golden diff is that it names what differs. A trailing
// newline is the case most likely to be invisible in an editor and therefore
// the one the diagnostic most needs to state, so it is pinned here.
func TestStringDiffReportsTrailingNewlineDifferences(t *testing.T) {
	cases := []struct {
		name     string
		exp, act string
		wantLine string
		// wantDetail pins WHICH side lacks the line. Asserting only the line
		// number would pass even if the diagnostic never said that, which is
		// the whole point of reporting a trailing-newline difference.
		wantDetail string
	}{
		{
			name:       "actual is missing the final newline",
			exp:        "a\n",
			act:        "a",
			wantLine:   "first difference at line 2",
			wantDetail: "actual:   <no line 2; actual ends at line 1>",
		},
		{
			name:       "actual has an extra final newline",
			exp:        "a",
			act:        "a\n",
			wantLine:   "first difference at line 2",
			wantDetail: "expected: <no line 2; expected ends at line 1>",
		},
		{
			name:       "actual has two trailing newlines where one is expected",
			exp:        "a\n",
			act:        "a\n\n",
			wantLine:   "first difference at line 3",
			wantDetail: "expected: <no line 3; expected ends at line 2>",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := stringDiff(tc.exp, tc.act)
			if !strings.Contains(got, tc.wantLine) {
				t.Fatalf("diff must locate the difference (%q); got:\n%s", tc.wantLine, got)
			}
			if !strings.Contains(got, "differ only in leading/trailing whitespace") {
				t.Fatalf("diff must say the difference is whitespace-only; got:\n%s", got)
			}
			if !strings.Contains(got, tc.wantDetail) {
				t.Fatalf("diff must say which side lacks the line (%q); got:\n%s", tc.wantDetail, got)
			}
			if !strings.Contains(got, "byte counts:") {
				t.Fatalf("diff must report byte counts; got:\n%s", got)
			}
		})
	}
}

// A difference in the middle of otherwise-identical content must still be
// located by line, and must not be reported as whitespace-only.
func TestStringDiffLocatesAContentDifference(t *testing.T) {
	got := stringDiff("one\ntwo\nthree\n", "one\nTWO\nthree\n")
	if !strings.Contains(got, "first difference at line 2") {
		t.Fatalf("expected line 2 to be named; got:\n%s", got)
	}
	if strings.Contains(got, "differ only in leading/trailing whitespace") {
		t.Fatalf("a content difference must not be reported as whitespace-only; got:\n%s", got)
	}
}

// A single-line golden can be kilobytes wide (rendered JSON), where printing
// both sides whole names the line but not the difference.
func TestStringDiffWindowsAVeryLongLine(t *testing.T) {
	exp := strings.Repeat("x", 400) + "userId" + strings.Repeat("y", 400)
	act := strings.Repeat("x", 400) + "owner!" + strings.Repeat("y", 400)
	got := stringDiff(exp, act)
	if !strings.Contains(got, "first differing column 401") {
		t.Fatalf("expected the differing column to be named; got:\n%s", got)
	}
	if !strings.Contains(got, "...") {
		t.Fatalf("expected the long line to be windowed; got:\n%s", got)
	}
}
