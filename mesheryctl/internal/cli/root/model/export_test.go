package model

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestExportModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	apiURL := "/api/meshmodels/export"
	testdataDir := filepath.Join(currDir, "testdata")
	fixturesDir := filepath.Join(currDir, "fixtures")

	// Create a custom directory for testing custom output location
	customDir := "customDir"
	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatalf("Failed to create custom directory: %v", err)
	}

	// Setup cleanup to remove exported files generated during testing
	t.Cleanup(func() {
		filesToClean := []string{
			"amd-gpu.tar",
			"amd-gpu.tar.gz",
		}

		for _, file := range filesToClean {
			os.RemoveAll(file)
		}

		// Clean up custom directory
		os.RemoveAll(customDir)
	})

	// setup current context
	utils.SetupContextEnv(t)
	// initialize mock server for handling requests
	utils.StartMockery(t)
	// create test helper to get the base URL for mock requests
	testContext := utils.NewTestHelper(t)

	// test scenarios for exporting models
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		ExpectedResponse string
		ExpectError      bool
		HttpCode         int
		FileName         string
		FileType         string
		FileLocation     string
	}{
		{
			Name:             "Export model without arguments",
			Args:             []string{"export"},
			ExpectedResponse: "export.model.no-args.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Export model with default parameters",
			Args:             []string{"export", "amd-gpu"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.default.output.golden",
			ExpectError:      false,
			HttpCode:         200,
			FileName:         "amd-gpu",
			FileType:         "oci",
			FileLocation:     ".",
		},
		{
			Name:             "Export model with TAR output type",
			Args:             []string{"export", "amd-gpu", "-o", "tar"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=tar&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.tar.output.golden",
			ExpectError:      false,
			HttpCode:         200,
			FileName:         "amd-gpu",
			FileType:         "tar",
			FileLocation:     ".",
		},
		{
			Name:             "Export model with custom location",
			Args:             []string{"export", "amd-gpu", "-l", customDir},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.custom-location.output.golden",
			ExpectError:      false,
			HttpCode:         200,
			FileName:         "amd-gpu",
			FileType:         "oci",
			FileLocation:     customDir,
		},
		{
			Name:             "Export model with version",
			Args:             []string{"export", "amd-gpu", "--version", ""},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.default.output.golden",
			ExpectError:      false,
			HttpCode:         200,
			FileName:         "amd-gpu",
			FileType:         "oci",
			FileLocation:     ".",
		},
		{
			Name:             "Export model with discarded components and relationships",
			Args:             []string{"export", "amd-gpu", "--discard-components", "--discard-relationships"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=false&relationships=false&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.default.output.golden",
			ExpectError:      false,
			HttpCode:         200,
			FileName:         "amd-gpu",
			FileType:         "oci",
			FileLocation:     ".",
		},
		{
			Name:             "Export non-existent model",
			Args:             []string{"export", "non-existent-model"},
			URL:              fmt.Sprintf("%s?name=non-existent-model&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.not-found.api.response.golden",
			ExpectedResponse: "export.model.not-found.output.golden",
			ExpectError:      true,
			HttpCode:         404,
		},
	}

	// A recursive function to reset flags for a command.
	var resetFlags func(*cobra.Command, *testing.T)
	resetFlags = func(c *cobra.Command, t *testing.T) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if err := f.Value.Set(f.DefValue); err != nil {
				t.Fatalf("failed to reset flag %q: %v", f.Name, err)
			}
		})
		for _, sub := range c.Commands() {
			resetFlags(sub, t)
		}
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Reset flags for each test case
			httpmock.Reset()
			resetFlags(ModelCmd, t)

			if tt.URL != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).LoadByte()
				httpmock.RegisterResponder("GET", testContext.BaseURL+tt.URL,
					httpmock.NewBytesResponder(tt.HttpCode, apiResponse))
			}

			utils.TokenFlag = utils.GetToken(t)

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			ModelCmd.SetOut(b)
			ModelCmd.SetArgs(tt.Args)
			err := ModelCmd.Execute()

			if tt.ExpectError {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
				expectedResponse := golden.Load()
				utils.Equals(t, expectedResponse, err.Error())
				return
			}
			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			// Check for successful file creation
			if tt.FileName != "" {
				var expectedFile string
				if tt.FileType == "tar" {
					expectedFile = filepath.Join(tt.FileLocation, tt.FileName+".tar.gz")
				} else {
					expectedFile = filepath.Join(tt.FileLocation, tt.FileName+".tar")
				}
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Fatalf("Expected file %s to be created, but it does not exist", expectedFile)
				}
			}

			// Check for correct output
			actualResponse := b.String()
			expectedResponse := golden.Load()

			utils.Equals(t, expectedResponse, actualResponse)
		})
	}

	// stop mock server
	utils.StopMockery(t)
}
