package model

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

var update = flag.Bool("update", false, "update golden files")

func TestIsPathValidUrl(t *testing.T) {
	var tests = []struct {
		name     string
		path     string
		expected bool
	}{
		{"valid http URL", "http://example.com", true},
		{"valid https URL", "https://example.com", true},
		{"valid git URL", "git://example.com/repo.git", true},
		{"invalid URL missing scheme", "example.com", false},
		{"invalid ftp URL", "ftp://example.com", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPathValidUrl(tt.path)
			if result != tt.expected {
				t.Errorf("isPathValidUrl(%s) = %v; want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsPathDirectory(t *testing.T) {
	tempDir := t.TempDir()
	f, _ := os.CreateTemp("", "tempFile")
	defer os.Remove(f.Name())
	tests := []struct {
		name        string
		path        string
		expected    bool
		expectError bool
	}{
		{"valid directory", tempDir, true, false},
		{"valid file", f.Name(), false, false},
		{"non-existent path", "./nonexistent", false, true},
		{"empty path", "", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := isPathDirectory(tt.path)
			if (err != nil) != tt.expectError {
				t.Errorf("isPathDirectory(%s) unexpected error state: got %v, want error %v", tt.path, err, tt.expectError)
				return
			}
			if result != tt.expected {
				t.Errorf("isPathDirectory(%s) = %v; want %v", tt.path, result, tt.expected)
			}
		})
	}
}


func TestImportModel(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures", "import")
	testdataDir := filepath.Join(currDir, "testdata", "import")
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Import model providing invalid template file path",
			Args:             []string{"import", "-f", "invalidTemplatePath"},
			ExpectedResponse: "invalid.template.file.path.golden",
			Fixture:          "invalid.template.file.path.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/register",
			Token:            filepath.Join(fixturesDir, "..", "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Import model providing empty template file",
			Args:             []string{"import", "-f", filepath.Join(fixturesDir, "empty.template.file.golden"), "-t", "https://example.com"},
			ExpectedResponse: "empty.template.file.golden",
			Fixture:          "empty.template.file.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/register",
			Token:            filepath.Join(fixturesDir, "..", "token.golden"),
			ExpectError:      false,
		},
		{
			Name:             "Import model providing empty template file",
			Args:             []string{"import", "-f", filepath.Join(fixturesDir, "empty.template.file.golden"), "-t", "https://example.com"},
			ExpectedResponse: "empty.template.file.golden",
			Fixture:          "empty.template.file.golden",
			URL:              testContext.BaseURL + "/api/meshmodels/register",
			Token:            filepath.Join(fixturesDir, "..", "token.golden"),
			ExpectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()
			utils.TokenFlag = tt.Token
			httpmock.RegisterResponder("POST", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			ModelCmd.SetOut(b)
			ModelCmd.SetArgs(tt.Args)
			err := ModelCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// write it in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}

			// response being printed in console
			output := b.String()
			actualResponse := output

			// write it in file
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
		t.Log("model import test Passed")
	}

	utils.StopMockery(t)
}