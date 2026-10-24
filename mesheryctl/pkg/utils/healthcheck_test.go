package utils

import (
	"testing"

	"k8s.io/apimachinery/pkg/version"
)

var versionCheck = []struct {
	version  string // input
	expected [3]int // expected output
}{
	{"1.12.0", [3]int{1, 12, 0}},
	{"1.12.1", [3]int{1, 12, 1}},
	{"1.12.2", [3]int{1, 12, 2}},
	{"1.12.3", [3]int{1, 12, 3}},
	{"1.12.4", [3]int{1, 12, 4}},
	{"1.12.5", [3]int{1, 12, 5}},
	{"1.12.6", [3]int{1, 12, 6}},
	{"1.12.7", [3]int{1, 12, 7}},
	{"1.12.8", [3]int{1, 12, 8}},
	{"1.12.9", [3]int{1, 12, 9}},
	{"1.12.10", [3]int{1, 12, 10}},
}

func TestGetK8sVersion(t *testing.T) {
	for _, tt := range versionCheck {
		t.Run(tt.version, func(t *testing.T) {
			if got, _ := getK8sVersion(tt.version); got != tt.expected {
				t.Errorf("getK8sVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckKubectlVersion(t *testing.T) {
	tests := []struct {
		version  string // input
		expected [3]int // expected output
	}{
		{
			version:  "v1.12.0",
			expected: [3]int{1, 12, 0},
		},
		{
			version:  "v1.12.1",
			expected: [3]int{1, 12, 1},
		},
		{
			version:  "v1.12.2",
			expected: [3]int{1, 12, 2},
		},
		{
			version:  "v1.12.3",
			expected: [3]int{1, 12, 3},
		},
		{
			version:  "v1.12.4",
			expected: [3]int{1, 12, 4},
		},
		{
			version:  "v1.12.5",
			expected: [3]int{1, 12, 5},
		},
		{
			version:  "v1.12.6",
			expected: [3]int{1, 12, 6},
		},
		{
			version:  "v1.12.7",
			expected: [3]int{1, 12, 7},
		},
		{
			version:  "v1.12.8",
			expected: [3]int{1, 12, 8},
		},
		{
			version:  "v1.12.9",
			expected: [3]int{1, 12, 9},
		},
		{
			version:  "v1.12.10",
			expected: [3]int{1, 12, 10},
		},
	}

	_, err := parseKubectlShortVersion("Not really a version")
	if err == nil {
		t.Fatalf("Expected to get an error")
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got, _ := parseKubectlShortVersion(tt.version); got != tt.expected {
				t.Errorf("checkKubectlVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckK8sVersion(t *testing.T) {
	tests := []struct {
		version  *version.Info
		expected [3]int
	}{
		{
			version: &version.Info{
				Major:        "1",
				Minor:        "2",
				GitVersion:   "v1.12.0",
				GitCommit:    "abcdefg",
				GitTreeState: "clean",
				BuildDate:    "2023-06-20",
				GoVersion:    "go1.16",
				Compiler:     "gc",
				Platform:     "linux/amd64",
			},
			expected: [3]int{1, 12, 0},
		},
	}

	for _, tt := range tests {
		t.Run("Check K8s Version", func(t *testing.T) {
			if got := CheckK8sVersion(tt.version); got != nil {
				t.Errorf("getK8sVersion() = %v, want %v", got, tt.expected)
			}
		})
	}

}

func TestIsCompatibleVersion(t *testing.T) {
	tests := []struct {
		name     string
		minimum  [3]int
		current  [3]int
		expected bool
	}{
		{"equal", [3]int{1, 12, 0}, [3]int{1, 12, 0}, true},
		{"higher major", [3]int{1, 11, 0}, [3]int{2, 13, 2}, true},
		{"lower major", [3]int{2, 13, 2}, [3]int{1, 11, 0}, false},
		{"higher minor", [3]int{1, 28, 3}, [3]int{1, 29, 0}, true},
		{"lower minor", [3]int{1, 28, 3}, [3]int{1, 27, 9}, false},
		{"higher patch", [3]int{1, 28, 3}, [3]int{1, 28, 4}, true},
		{"lower patch", [3]int{1, 28, 3}, [3]int{1, 28, 2}, false},
		{"higher major overrides lower minor", [3]int{1, 28, 3}, [3]int{2, 0, 0}, true},
		{"lower major overrides higher minor", [3]int{2, 0, 0}, [3]int{1, 99, 99}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCompatibleVersion(tt.minimum, tt.current)
			if got != tt.expected {
				t.Errorf("minimum=%v current=%v: got %v, want %v",
					tt.minimum, tt.current, got, tt.expected)
			}
		})
	}
}

func TestAreMesheryComponentsRunning(t *testing.T) {
	tests := []struct {
		platform string
		expected bool
	}{
		{
			platform: "docker",
			expected: false,
		},
		{
			platform: "kubernetes",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("Meshery Components Running", func(t *testing.T) {
			if got, _ := AreMesheryComponentsRunning(tt.platform); got != tt.expected {
				t.Errorf("Meshery Components Running got %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAreAllPodsRunning(t *testing.T) {
	tests := []struct {
		expected bool
	}{
		{
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run("Are All Pods Running", func(t *testing.T) {
			if got, _ := AreAllPodsRunning(); got != tt.expected {
				t.Errorf("Are all pods running got %v want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckMesheryNsDelete(t *testing.T) {
	expected := false
	t.Run("Test Existing Meshery Namespace", func(t *testing.T) {
		if got, _ := CheckMesheryNsDelete(); got != expected {
			t.Errorf("Namespace not deleted got %v want %v ", got, expected)
		}
	})
}
