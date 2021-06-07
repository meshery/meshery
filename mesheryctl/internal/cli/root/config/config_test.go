package config

import "testing"

func TestGetCommitSHA(t *testing.T) {
	version := Version{"", "abcd1234", ""}
	got := version.GetCommitSHA()
	want := "abcd1234"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
