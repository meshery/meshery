package config

import (
	"os"
	"testing"
)

func TestGetCommitSHA(t *testing.T) {
	version := Version{"", "abcd1234", ""}
	got := version.GetCommitSHA()
	want := "abcd1234"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
func TestGetBuild(t *testing.T) {
	version := Version{"build", "", ""}
	got := version.GetBuild()
	want := "build"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetLocation(t *testing.T) {
	token := Token{}
	token.SetLocation("")
	got := token.GetLocation()
	want, err := os.UserHomeDir()
	want += "/.meshery"
	if err != nil {
		t.Errorf("Fail")
	}
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetVersion(t *testing.T) {
	context := Context{"", "", "", nil, "", ""}
	context.SetVersion("abcd1234")
	got := context.GetVersion()
	want := "abcd1234"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetName(t *testing.T) {
	token := Token{}
	token.SetName("Name")
	got := token.Name
	want := "Name"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetChannel(t *testing.T) {
	context := Context{"", "", "", nil, "", ""}
	context.SetChannel("1234")
	got := context.GetChannel()
	want := "1234"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetAdapters(t *testing.T) {
	context := Context{"", "", "", nil, "", ""}
	got := context.GetAdapters()
	want := []string(nil)
	for i, j := range got {
		if j != want[i] {
			t.Errorf("got %q want %q", got, want)
		}
	}
}

func TestGetPlatform(t *testing.T) {
	context := Context{"", "", "", nil, "", ""}
	context.SetPlatform("windowslinuxmacos")
	got := context.GetPlatform()
	want := "windowslinuxmacos"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetToken(t *testing.T) {
	context := Context{}
	context.SetToken("1234")
	got := context.GetToken()
	want := "1234"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestGetEndpoint(t *testing.T) {
	context := Context{}
	context.SetEndpoint("12asgherdh34")
	got := context.GetEndpoint()
	want := "12asgherdh34"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
