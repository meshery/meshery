package utils

import (
	"flag"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
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
			expectOutput: "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5",
			fixture:      "manifesturl.api.response.golden",
			url:          "https://api.github.com/repos/meshery/meshery/git/trees/v0.5.45?recursive=1",
			version:      "v0.5.45",
		},
		{
			name:      "Test getting manifest tree url with wrong url",
			expectErr: true,
			fixture:   "manifesturl.err.api.response.golden",
			url:       "https://api.github.com/repos/meshery/meshery/git/trees/v0000000?recursive=1",
			version:   "v0000000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock response
			apiResponse := NewGoldenFile(t, tt.fixture, fixturesDir).Load()
			httpmock.RegisterResponder("GET", tt.url,
				httpmock.NewStringResponder(200, apiResponse))
			url, err := GetManifestTreeURL(tt.version)
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := NewGoldenFile(t, tt.fixture, testdataDir)
			if err != nil {
				if tt.expectErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := tt.expectOutput
					Equals(t, expectedResponse, err.Error())
					return
				}
			}

			Equals(t, tt.expectOutput, url)
		})
	}
	// stop mock server
	StopMockery(t)
}
