package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetDeploymentVersion covers the deployment image shapes that previously
// reached unchecked slice indexes (Containers[0], the ":" split [1], and the
// "-" split [1]) and panicked the CLI. A non-standard image must now surface a
// returned error, consistent with the rest of the function.
func TestGetDeploymentVersion(t *testing.T) {
	deployment := func(image string) string {
		return "apiVersion: apps/v1\nkind: Deployment\nspec:\n  template:\n    spec:\n      containers:\n        - name: meshery\n          image: " + image + "\n"
	}
	noContainers := "apiVersion: apps/v1\nkind: Deployment\nspec:\n  template:\n    spec:\n      containers: []\n"

	tests := []struct {
		name    string
		yaml    string
		want    string
		wantErr bool
	}{
		{name: "standard tag", yaml: deployment("layer5/meshery:stable-latest"), want: "latest", wantErr: false},
		{name: "untagged image", yaml: deployment("layer5/meshery"), wantErr: true},
		{name: "tag without dash", yaml: deployment("layer5/meshery:latest"), wantErr: true},
		{name: "no containers", yaml: noContainers, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), MesheryDeployment)
			if err := os.WriteFile(path, []byte(tt.yaml), 0o600); err != nil {
				t.Fatalf("write temp deployment: %v", err)
			}
			got, err := GetDeploymentVersion(path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetDeploymentVersion() err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Fatalf("GetDeploymentVersion() = %q, want %q", got, tt.want)
			}
		})
	}
}
