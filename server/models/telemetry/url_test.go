package telemetry

import "testing"

func TestNormalizeBaseURL(t *testing.T) {
	cases := map[string]string{
		"  http://h/p/  ": "http://h/p",
		"http://h/p#frag": "http://h/p",
		"http://h/p?a=b":  "http://h/p",
		"http://h/p#?a=b": "http://h/p",
		"http://h":        "http://h",
		"":                "",
	}
	for in, want := range cases {
		if got := NormalizeBaseURL(in); got != want {
			t.Errorf("NormalizeBaseURL(%q) = %q, want %q", in, got, want)
		}
	}
}
