package utils

import (
	"io/ioutil"
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

func TestGetK8sVersionInfo(t *testing.T) {
	tests := []struct {
		version  string        // input
		expected *version.Info // expected output
		fixture  string        // fixture file
		golden   string        // golden fileō
	}{
		{
			version:  "v1.12.0",
			expected: &version.Info{Major: "1", Minor: "12", GitVersion: "1.12.0-alpha.1+1"},
			fixture:  "k8s_version_info.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			_, err := ioutil.ReadFile("fixtures/healthcheck/" + tt.fixture)
			if err != nil {
				t.Fatalf("Error reading fixture file: %v", err)
			}
		})
	}
}

func TestCheckK8sVersion(t *testing.T) {
	for _, tt := range versionCheck {
		t.Run(tt.version, func(t *testing.T) {
			_, err := getK8sVersion(tt.version)
			if err != nil {
				t.Errorf("Error checking k8s version: %v", err)
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
