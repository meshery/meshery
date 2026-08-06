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
}
