package utils

import (
	"encoding/json"
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/constants"
)

var update = flag.Bool("update", false, "update golden files")

func TestGetManifestTreeURL(t *testing.T) {
	// initialize mock server for handling requests
	StartMockery(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		name         string
		version      string
		url          string
		expectOutput string
		fixture      string
		expectErr    bool
	}{
		{
			name:         "Test getting manifest tree url",
			expectErr:    false,
			expectOutput: "manifesturl.expect.golden",
			fixture:      "manifesturl.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/v0.5.45?recursive=1",
			version:      "v0.5.45",
		},
		{
			name:         "Test getting manifest tree url with wrong version",
			expectErr:    true,
			expectOutput: "manifesturl.out.err.golden",
			fixture:      "manifesturl.err.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/v0000000?recursive=1",
			version:      "v0000000",
		},
		{
			name:         "Test getting manifest tree url with empty version",
			expectErr:    true,
			expectOutput: "manifesturlempty.out.err.golden",
			fixture:      "manifesturl.empty.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/?recursive=1",
			version:      "",
		},
		{
			name:         "Test getting manifest tree url with wrong version",
			expectErr:    false,
			expectOutput: "manifesturl.out.latest.golden",
			fixture:      "manifesturl.latest.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/" + "master" + "?recursive=1",
			version:      "master",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := NewGoldenFile(t, tt.expectOutput, testdataDir)
			// mock response
			apiResponse := NewGoldenFile(t, tt.fixture, fixturesDir).Load()
			httpmock.RegisterResponder("GET", tt.url,
				httpmock.NewStringResponder(200, apiResponse))
			actualurl, err := GetManifestTreeURL(tt.version)
			if err != nil {
				if tt.expectErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					Equals(t, expectedResponse, err.Error())
					return
				}
			}
			if *update {
				golden.Write(actualurl)
			}
			expectedResponse := golden.Load()
			Equals(t, expectedResponse, actualurl)
		})
	}
	// stop mock server
	StopMockery(t)
}

func TestListManifests(t *testing.T) {
	// initialize mock server for handling requests
	StartMockery(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		name         string
		url          string
		expectOutput string
		fixture      string
		expectErr    bool
	}{
		{
			name:         "Test listing manifests",
			expectErr:    false,
			expectOutput: "listmanifest.expect.golden",
			fixture:      "listmanifest.api.response.golden",
			url:          "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5",
		},
		{
			name:         "Test listing manifests with wrong url",
			expectErr:    true,
			expectOutput: "listmanifest.expect.err.golden",
			fixture:      "listmanifest.api.err.response.golden",
			url:          "https://api.github.com/repos/meshery/meshery/git/trees/gibberish",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := NewGoldenFile(t, tt.expectOutput, testdataDir)
			// mock response
			apiResponse := NewGoldenFile(t, tt.fixture, fixturesDir).Load()
			httpmock.RegisterResponder("GET", tt.url,
				httpmock.NewStringResponder(200, apiResponse))
			manifests, err := ListManifests(tt.url)
			if err != nil {
				if tt.expectErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					Equals(t, expectedResponse, err.Error())
					return
				}
			}
			manifestsactual, err := json.Marshal(&manifests)
			if err != nil {
				t.Error("Could not unmarshall manifests from response")
			}
			if *update {
				golden.Write(string(manifestsactual))
			}
			expectedResponse := golden.Load()
			Equals(t, expectedResponse, string(manifestsactual))
		})
	}
	// stop mock server
	StopMockery(t)
}
