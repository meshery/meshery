package perf

import (
	"os"
	"path/filepath"
	"testing"

	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func TestApplyConfigResolution(t *testing.T) {
	for _, tc := range []struct {
		url   string
		valid bool
	}{{"", false}, {"not a url", false}, {"https://meshery.io", true}, {"http://localhost:9081", true}} {
		if err := validateTestURL(tc.url); (err == nil) != tc.valid {
			t.Errorf("validateTestURL(%q): %v", tc.url, err)
		}
	}
	dir := t.TempDir()
	for _, tc := range []struct {
		name, content string
		valid         bool
	}{
		{"bad", "test: [", false}, {"empty", "test:\n  clients: []\n", false},
		{"valid", "test:\n  name: from-file\n  duration: 30s\n  clients:\n    - rps: 7\n      connections: 2\n      endpointUrls: [https://meshery.io]\nmesh: istio\n", true},
	} {
		path := filepath.Join(dir, tc.name)
		if err := os.WriteFile(path, []byte(tc.content), 0600); err != nil {
			t.Fatal(err)
		}
		cfg, err := loadSMPConfigFile(path)
		if (err == nil) != tc.valid {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if !tc.valid {
			continue
		}
		oldName, oldURL, oldMesh, oldQPS, oldConcurrent, oldDuration, oldGenerator, oldBody, oldProfile := testName, testURL, testMesh, qps, concurrentRequests, testDuration, loadGenerator, loadTestBody, profileName
		defer func() {
			testName, testURL, testMesh, qps, concurrentRequests, testDuration, loadGenerator, loadTestBody, profileName = oldName, oldURL, oldMesh, oldQPS, oldConcurrent, oldDuration, oldGenerator, oldBody, oldProfile
		}()
		testName, testURL, testMesh, qps, concurrentRequests, testDuration, loadGenerator, loadTestBody = "flag", "", "", "", "", "", "", ""
		consolidateSMPConfig(cfg)
		if testName != "flag" || qps != "7" || concurrentRequests != "2" || testMesh != "istio" {
			t.Fatalf("file precedence failed: %q %q %q %q", testName, qps, concurrentRequests, testMesh)
		}
		testURL = ""
		profile := &perfprofile.PerformanceProfile{Endpoints: []string{"https://example.com"}, QPS: 100, RequestHeaders: "headers"}
		headers, _, _ := consolidateExistingProfileConfig(profile, true)
		if testURL != "https://example.com" || qps != "7" || headers != "headers" {
			t.Fatalf("profile precedence failed: %q %q %q", testURL, qps, headers)
		}
	}
	if _, err := loadSMPConfigFile(filepath.Join(dir, "missing")); err == nil {
		t.Fatal("expected missing file error")
	}
}
