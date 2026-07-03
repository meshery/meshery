package model

import "testing"

// TestImageTag documents the tag parsing used by GetControlPlaneState. The
// container index that feeds it is now guarded against an empty Containers
// slice (a container-less pod spec from meshsync used to panic with index out
// of range inside the control-plane subscription goroutine, crashing the
// server).
func TestImageTag(t *testing.T) {
	tests := []struct {
		image string
		want  string
	}{
		{image: "", want: ""},
		{image: "nginx", want: ""},
		{image: "istio/proxyv2:1.20.0", want: "1.20.0"},
		{image: "docker.io/istio/proxyv2:1.20.0", want: "1.20.0"},
	}
	for _, tt := range tests {
		if got := imageTag(tt.image); got != tt.want {
			t.Errorf("imageTag(%q) = %q, want %q", tt.image, got, tt.want)
		}
	}
}
