package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

// writeFile writes content into t's temp dir and returns the absolute path.
func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("writing fixture %s: %v", name, err)
	}
	return path
}

// clusterKubeconfig builds a kubeconfig whose single cluster stanza carries key
// with the given value, so a test can place a crafted path under any field.
func clusterKubeconfig(key, value string) []byte {
	return []byte(fmt.Sprintf(`apiVersion: v1
kind: Config
current-context: c1
clusters:
- name: c1
  cluster:
    server: https://127.0.0.1:6443
    %s: %s
contexts:
- name: c1
  context:
    cluster: c1
    user: u1
users:
- name: u1
  user:
    token: abc
`, key, value))
}

// clusterStanza parses a flattened kubeconfig and returns its single cluster's
// inner `cluster` map. Asserting on the parsed tree rather than on the marshalled
// text keeps these tests independent of yaml.v2's line folding, which wraps a
// long base64 scalar across lines.
func clusterStanza(t *testing.T, out []byte) map[interface{}]interface{} {
	t.Helper()

	cfg := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(out, &cfg); err != nil {
		t.Fatalf("flattened kubeconfig is not valid yaml: %v\n%s", err, out)
	}

	clusters, ok := cfg["clusters"].([]interface{})
	if !ok || len(clusters) != 1 {
		t.Fatalf("expected exactly one cluster in the flattened kubeconfig, got: %#v", cfg["clusters"])
	}
	entry, ok := clusters[0].(map[interface{}]interface{})
	if !ok {
		t.Fatalf("cluster entry is not a map: %#v", clusters[0])
	}
	stanza, ok := entry["cluster"].(map[interface{}]interface{})
	if !ok {
		t.Fatalf("cluster stanza is not a map: %#v", entry["cluster"])
	}
	return stanza
}

// TestFlattenMinifyKubeConfigOnlyResolvesCredentialFields pins the security
// boundary of the kubeconfig flattener.
//
// The kubeconfig reaching FlattenMinifyKubeConfig is uploaded by a user (the
// `k8sfile` part of POST /api/system/kubernetes and
// POST /api/system/kubernetes/contexts). Resolving every value that merely
// looked like an existing path made the function an arbitrary-file-read
// primitive: a crafted kubeconfig naming a server-side file under any key had
// that file's bytes inlined, and the result was echoed back in the contexts
// those handlers return and persisted into the connection's credential.
//
// Only the three fields a kubeconfig legitimately points at a file with may be
// resolved. The deny cases are the ones that matter: an allow-only test passes
// against the vulnerable implementation too.
func TestFlattenMinifyKubeConfigOnlyResolvesCredentialFields(t *testing.T) {
	const secret = "eyJhbGciOiJSUzI1NiJ9.super-secret-service-account-token"
	secretPath := writeFile(t, "token", secret)
	encodedSecret := base64.StdEncoding.EncodeToString([]byte(secret))

	tests := []struct {
		name string
		// key is the field the crafted path is placed under.
		key string
		// wantResolved is true when the field is one the flattener may inline.
		wantResolved bool
	}{
		{name: "arbitrary key is never resolved", key: "leak", wantResolved: false},
		{name: "tls-server-name is never resolved", key: "tls-server-name", wantResolved: false},
		{name: "tokenFile is never resolved", key: "tokenFile", wantResolved: false},
		{name: "certificate-authority is resolved", key: "certificate-authority", wantResolved: true},
		{name: "client-certificate is resolved", key: "client-certificate", wantResolved: true},
		{name: "client-key is resolved", key: "client-key", wantResolved: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := FlattenMinifyKubeConfig(clusterKubeconfig(tt.key, secretPath))
			if err != nil {
				t.Fatalf("FlattenMinifyKubeConfig returned an error: %v", err)
			}
			stanza := clusterStanza(t, out)

			dataKey := tt.key + "-data"

			if !tt.wantResolved {
				if got, present := stanza[dataKey]; present {
					t.Fatalf("field %q was resolved against the server's filesystem and inlined as %q (%v): any uploaded kubeconfig can now read an arbitrary server-side file", tt.key, dataKey, got)
				}
				if got := stanza[tt.key]; got != secretPath {
					t.Errorf("a field the flattener does not own must pass through untouched, want %q, got %#v", secretPath, got)
				}
				return
			}

			got, present := stanza[dataKey]
			if !present {
				t.Fatalf("field %q was not inlined; flattening no longer works for a real credential field. stanza: %#v", tt.key, stanza)
			}
			if got != encodedSecret {
				t.Errorf("inlined value under %q: want the base64 of the referenced file, got %#v", dataKey, got)
			}
			if _, stillThere := stanza[tt.key]; stillThere {
				t.Errorf("the original path reference %q should have been replaced by %q, stanza: %#v", tt.key, dataKey, stanza)
			}
		})
	}
}

// TestFlattenMinifyKubeConfigRejectsNonRegularAndOversizedFiles covers the reads
// os.Stat alone would have allowed. A directory and a device node both stat
// successfully, and reading a device node such as /dev/zero never returns; an
// oversized file would be read into memory in full.
func TestFlattenMinifyKubeConfigRejectsNonRegularAndOversizedFiles(t *testing.T) {
	dir := t.TempDir()
	oversized := filepath.Join(dir, "huge.crt")
	if err := os.WriteFile(oversized, make([]byte, maxInlinedCredentialSize+1), 0o600); err != nil {
		t.Fatalf("writing oversized fixture: %v", err)
	}

	tests := []struct {
		name string
		path string
	}{
		{name: "directory", path: dir},
		{name: "oversized file", path: oversized},
		{name: "missing file", path: filepath.Join(dir, "does-not-exist.crt")},
		{name: "bare filename is not resolved against the working directory", path: "helpers.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := FlattenMinifyKubeConfig(clusterKubeconfig("certificate-authority", tt.path))
			if err != nil {
				t.Fatalf("FlattenMinifyKubeConfig returned an error: %v", err)
			}
			stanza := clusterStanza(t, out)

			if got, present := stanza["certificate-authority-data"]; present {
				t.Errorf("%s must not be inlined, got %#v", tt.name, got)
			}
			if got := stanza["certificate-authority"]; got != tt.path {
				t.Errorf("an unresolvable reference must be left untouched, want %q, got %#v", tt.path, got)
			}
		})
	}
}
