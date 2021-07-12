package utils

import "testing"

var versionTests = []struct {
	v        string // input
	expected [3]int // expected
}{
	{"Client Version: v1.10.1", [3]int{1, 10, 1}},
	{"client Version: v1.10.1", [3]int{1, 10, 1}},
	{"Client Version - v1.10.1", [3]int{1, 10, 1}},
	{"v1.10.1", [3]int{1, 10, 1}},
	{"Client Version: v1.10.1 beta:2", [3]int{1, 10, 1}},
	{"Client Version: v2.1348.1", [3]int{2, 1348, 1}},
}

func TestParseKubectlShortVersion(t *testing.T) {
	for _, tt := range versionTests {
		actual, err := parseKubectlShortVersion(tt.v)
		if err != nil {
			t.Fatalf("Unexpected error while parsing kubectl short version: %v", err)
		}
		if actual != tt.expected {
			t.Fatalf("Expected to get %v but got %v", tt.expected, actual)
		}
	}
}

func TestParseKubectlShortVersionIncorrectVersion(t *testing.T) {
	_, err := parseKubectlShortVersion("Not really a version")
	if err == nil {
		t.Fatalf("Expected to get an error")
	}
}
