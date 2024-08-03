package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestParseKubectlShortVersion(t *testing.T) {
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
		name        string
		versionInfo *version.Info
		expectError bool
	}{
		{
			name: "Compatible version",
			versionInfo: &version.Info{
				Major:      "1",
				Minor:      "18",
				GitVersion: "v1.18.0",
			},
			expectError: false,
		},
		{
			name: "Incompatible major version",
			versionInfo: &version.Info{
				Major:      "0",
				Minor:      "18",
				GitVersion: "v0.18.0",
			},
			expectError: true,
		},
		{
			name: "Invalid version string",
			versionInfo: &version.Info{
				Major:      "1",
				Minor:      "12",
				GitVersion: "invalid-version",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckK8sVersion(tt.versionInfo)
			if (err != nil) != tt.expectError {
				t.Errorf("CheckK8sVersion() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
func TestIsCompatibleVersion(t *testing.T) {
	tests := []struct {
		minimum  [3]int
		current  [3]int
		expected bool
	}{
		{
			minimum:  [3]int{1, 12, 0},
			current:  [3]int{1, 12, 0},
			expected: true,
		},
		{
			minimum:  [3]int{1, 11, 0},
			current:  [3]int{2, 13, 2},
			expected: true,
		},
		{
			minimum:  [3]int{2, 13, 2},
			current:  [3]int{1, 11, 0},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("Check compatible Version", func(t *testing.T) {
			if got := isCompatibleVersion(tt.minimum, tt.current); got != tt.expected {
				t.Errorf("Version not compatible got %v, want %v", got, tt.expected)
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

func TestCheckKubectlVersion(t *testing.T) {
	tests := []struct {
		name           string
		kubectlOutput  string
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "Valid kubectl version",
			kubectlOutput: "Client Version: v1.20.0\n",
			expectedError: false,
		},
		{
			name:           "Invalid kubectl version",
			kubectlOutput:  "Client Version: v1.10.0\n",
			expectedError:  true,
			expectedErrMsg: "kubectl is on version [1.10.0], but version [1.12.0] or more recent is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory to hold the fake kubectl binary
			tempDir := t.TempDir()
			fakeKubectlPath := filepath.Join(tempDir, "kubectl")

			// Write the fake kubectl binary
			err := os.WriteFile(fakeKubectlPath, []byte(fmt.Sprintf("#!/bin/sh\necho '%s'", tt.kubectlOutput)), 0755)
			assert.NoError(t, err)

			// Prepend the temporary directory to the PATH
			originalPath := os.Getenv("PATH")
			defer os.Setenv("PATH", originalPath)
			os.Setenv("PATH", fmt.Sprintf("%s:%s", tempDir, originalPath))

			err = CheckKubectlVersion()

			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
