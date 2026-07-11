package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewKubeClientMissingConfigPathFallsBack(t *testing.T) {
	orig := ConfigPath
	t.Cleanup(func() { ConfigPath = orig })

	// Point at a path that does not exist; ambient resolution may still fail
	// if the machine has no kubeconfig — we only assert we do not error *because*
	// of a missing ConfigPath before ambient lookup.
	ConfigPath = filepath.Join(t.TempDir(), "does-not-exist.yaml")

	_, err := NewKubeClient()
	// Ambient may succeed or fail depending on environment; missing ConfigPath
	// must not produce the "invalid Meshery kubeconfig" / "unable to read" errors.
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "invalid Meshery kubeconfig") ||
			strings.Contains(msg, "unable to read Meshery kubeconfig") ||
			strings.Contains(msg, "is empty") {
			t.Fatalf("missing ConfigPath should fall back to ambient, got Meshery-path error: %v", err)
		}
	}
}

func TestNewKubeClientEmptyFileFailsLoudly(t *testing.T) {
	orig := ConfigPath
	t.Cleanup(func() { ConfigPath = orig })

	dir := t.TempDir()
	path := filepath.Join(dir, "kubeconfig.yaml")
	if err := os.WriteFile(path, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	ConfigPath = path

	_, err := NewKubeClient()
	if err == nil {
		t.Fatal("expected error for empty Meshery kubeconfig")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Fatalf("expected empty-file error, got: %v", err)
	}
}

func TestNewKubeClientInvalidFileFailsLoudly(t *testing.T) {
	orig := ConfigPath
	t.Cleanup(func() { ConfigPath = orig })

	dir := t.TempDir()
	path := filepath.Join(dir, "kubeconfig.yaml")
	if err := os.WriteFile(path, []byte("this-is-not-valid-kubeconfig"), 0o600); err != nil {
		t.Fatal(err)
	}
	ConfigPath = path

	_, err := NewKubeClient()
	if err == nil {
		t.Fatal("expected error for invalid Meshery kubeconfig")
	}
	if !strings.Contains(err.Error(), "invalid Meshery kubeconfig") {
		t.Fatalf("expected invalid-file error, got: %v", err)
	}
}

func TestNewKubeClientValidFileSucceeds(t *testing.T) {
	orig := ConfigPath
	t.Cleanup(func() { ConfigPath = orig })

	// Minimal syntactically valid kubeconfig. Client construction may still
	// fail later if the API server is unreachable; meshkit.New only needs
	// a parseable config with a usable rest.Config.
	valid := []byte(`apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1:1
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-context
current-context: test-context
users:
- name: test-user
  user:
    token: test-token
`)
	dir := t.TempDir()
	path := filepath.Join(dir, "kubeconfig.yaml")
	if err := os.WriteFile(path, valid, 0o600); err != nil {
		t.Fatal(err)
	}
	ConfigPath = path

	client, err := NewKubeClient()
	if err != nil {
		t.Fatalf("expected valid Meshery kubeconfig to build a client: %v", err)
	}
	if client == nil || client.KubeClient == nil {
		t.Fatal("expected non-nil client")
	}
}
