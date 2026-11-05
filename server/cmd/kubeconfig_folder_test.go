package main

import (
	"errors"
	"path"
	"strings"
	"testing"
)

func TestDefaultKubeconfigFolder(t *testing.T) {
	got, err := defaultKubeconfigFolder("/home/meshery", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := path.Join("/home/meshery", ".kube"); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	got, err = defaultKubeconfigFolder("", errors.New("$HOME is not defined"))
	if err == nil {
		t.Fatalf("expected an error when the home directory lookup fails, got folder %q", got)
	}
	if got != "" {
		t.Fatalf("expected no folder on error, got %q", got)
	}
	if !strings.Contains(err.Error(), "$HOME is not defined") {
		t.Fatalf("error should wrap the home directory error, got %v", err)
	}
}
