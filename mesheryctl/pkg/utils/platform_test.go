// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// The function are related to download should be test in meshkit package, please do not add test here.

func TestListManifests(t *testing.T) {
	t.Run("ListManifests with empty manifest", func(t *testing.T) {
		url := "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5"

		_, err := ListManifests(url)
		if err != nil {
			t.Errorf("ListManifests failed: %v", err)
		}
	})
}
func TestGetDeploymentVersion(t *testing.T) {
	writeDeploymentFile := func(t *testing.T, contents string) string {
		t.Helper()
		path := filepath.Join(t.TempDir(), "meshery-deployment.yaml")
		if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
			t.Fatalf("failed to write fixture: %v", err)
		}
		return path
	}

	t.Run("valid version-build tag returns the build version", func(t *testing.T) {
		path := writeDeploymentFile(t, `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: meshery
          image: layer5/meshery:v0.7.65-abc123
`)
		version, err := GetDeploymentVersion(path)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if version != "abc123" {
			t.Errorf("expected version %q, got %q", "abc123", version)
		}
	})

	t.Run("plain tag without a hyphen returns an error instead of panicking", func(t *testing.T) {
		path := writeDeploymentFile(t, `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: meshery
          image: layer5/meshery:v0.7.65
`)
		if _, err := GetDeploymentVersion(path); err == nil {
			t.Fatal("expected an error for a tag without a hyphen, got nil")
		}
	})

	t.Run("image with no tag returns an error instead of panicking", func(t *testing.T) {
		path := writeDeploymentFile(t, `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: meshery
          image: layer5/meshery
`)
		if _, err := GetDeploymentVersion(path); err == nil {
			t.Fatal("expected an error for an image with no tag, got nil")
		}
	})

	t.Run("empty containers list returns an error instead of panicking", func(t *testing.T) {
		path := writeDeploymentFile(t, `
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers: []
`)
		if _, err := GetDeploymentVersion(path); err == nil {
			t.Fatal("expected an error for an empty containers list, got nil")
		}
	})
	t.Run("registry host with a port is not mistaken for the tag", func(t *testing.T) {
		content := `
spec:
  template:
    spec:
      containers:
        - image: "registry.example.com:5000/meshery/meshery:stable-v0.7.0"
`
		path := writeDeploymentFile(t, content)
		got, err := GetDeploymentVersion(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "v0.7.0" {
			t.Errorf("got %q, want %q", got, "v0.7.0")
		}
	})

	t.Run("tag with multiple hyphens keeps the full build component", func(t *testing.T) {
		content := `
spec:
  template:
    spec:
      containers:
        - image: "meshery/meshery:stable-v0.7.0-rc1"
`
		path := writeDeploymentFile(t, content)
		got, err := GetDeploymentVersion(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "v0.7.0-rc1" {
			t.Errorf("got %q, want %q", got, "v0.7.0-rc1")
		}
	})

	t.Run("tag with an empty build component returns an error", func(t *testing.T) {
		content := `
spec:
  template:
    spec:
      containers:
        - image: "meshery/meshery:stable-"
`
		path := writeDeploymentFile(t, content)
		_, err := GetDeploymentVersion(path)
		if err == nil {
			t.Fatal("expected an error for empty build component, got nil")
		}
	})

	t.Run("tag with an empty version component returns an error", func(t *testing.T) {
		content := `
spec:
  template:
    spec:
      containers:
        - image: "meshery/meshery:-v0.7.0"
`
		path := writeDeploymentFile(t, content)
		_, err := GetDeploymentVersion(path)
		if err == nil {
			t.Fatal("expected an error for empty version component, got nil")
		}
	})

	t.Run("digest reference returns an error instead of panicking", func(t *testing.T) {
		content := `
spec:
  template:
    spec:
      containers:
        - image: "meshery/meshery@sha256:abc123"
`
		path := writeDeploymentFile(t, content)
		_, err := GetDeploymentVersion(path)
		if err == nil {
			t.Fatal("expected an error for digest reference, got nil")
		}
	})
}
